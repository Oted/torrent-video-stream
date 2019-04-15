package client

import (
	"errors"
	"fmt"
	"github.com/Oted/torrent-video-stream/lib/logger"
	"github.com/Oted/torrent-video-stream/lib/peer"
	"github.com/Oted/torrent-video-stream/lib/torrent"
	"github.com/Oted/torrent-video-stream/lib/tracker"
	"net"
	"runtime"
	"sync"
	"time"
)

//if piece empty or null then done
type result struct {
	piece                *torrent.Piece
	chunkPositionInPiece int64
	bytes                []byte
}

type Client struct {
	torrent            *torrent.Torrent
	MaxPeers           int
	Ip                 string
	Id                 [20]byte
	listener           net.Listener
	Peers              map[string]*peer.Peer //[ip:port]Peer
	Errors             chan error
	Tracker            tracker.Tracker
	AvailablePeers     chan string
	Results            chan result
	Jobs               chan *torrent.Piece
	done               bool
	seek               int64
	finished           []bool
	mapLock            *sync.Mutex
	latestChunk        int64
	selectedFileOffset int64
}

const workers = 1
const chunkWaitTimeout = 10
const peerCheckInterval = 3
const ChunkSize = int64(1024)
const MaxCurrentJobs = 1

func New(ip string, startPort int, t *torrent.Torrent) (error, *Client) {
	err, listener, port := findAvailPort(startPort)
	if err != nil {
		return err, nil
	}

	id := generateId()

	err, tracker := tracker.Create(t, ip, port, id)
	if err != nil {
		return err, nil
	}

	c := Client{
		latestChunk:        t.SelectedFile().Start/ChunkSize,
		torrent:            t,
		MaxPeers:           10,
		Ip:                 ip,
		Id:                 id,
		Tracker:            tracker,
		listener:           listener,
		Peers:              make(map[string]*peer.Peer),
		Errors:             make(chan error, 1),
		Jobs:               make(chan *torrent.Piece, len(t.Meta.SelectedPieces)),
		Results:            make(chan result, len(t.Meta.SelectedPieces)),
		done:               false,
		seek:               0,
		finished:           make([]bool, len(t.Meta.SelectedPieces)),
		mapLock:            &sync.Mutex{},
		AvailablePeers:     make(chan string),
		selectedFileOffset: t.SelectedFile().Start % t.Info.PieceLength,
	}

	return nil, &c
}

func (c *Client) StartDownload() error {
	if !IsLocal(c.Ip) {
		go c.listen()
	}

	err, res := c.Tracker.Announce(nil)
	if err != nil {
		return err
	}

	err = c.addPeers(*res)
	if err != nil {
		return err
	}

	err = c.Interested()
	if err != nil {
		return err
	}

	c.startJobs(workers)

	return nil
}

func (c *Client) startJobs(count int) {
	for w := 1; w <= count; w++ {
		go func(wor int) {
			for piece := range c.Jobs {
				err, res := c.getPiece(piece)

				if err != nil {
					c.fatal(err)
					return
				}

				go c.GotPiece(piece, res)
			}
		}(w)
	}

	targetPieces := c.torrent.Meta.SelectedPieces

	for _, p := range targetPieces {
		c.Jobs <- p
	}

	close(c.Jobs)
}

//have to return full byte slice of piece
//which should match the torrent.info.piece.length
//except in the case of the last piece
func (c *Client) getPiece(p *torrent.Piece) (error, []byte) {
	res := make([]byte, c.torrent.Info.PieceLength)
	offset := int64(0)

	if p.Index ==  c.torrent.Meta.SelectedPieces[0].Index {
		offset = c.selectedFileOffset
	}

	runtime.Gosched()
	for peerId := range c.AvailablePeers {
		targetPeer := c.Peers[peerId]

		if targetPeer == nil {
			logger.Error(errors.New("non existing peer"))
			continue
		}

		if offset >= c.torrent.Info.PieceLength {
			//here done
			break
		}

		length := ChunkSize

		//special last chunk size
		if offset + ChunkSize > c.torrent.Info.PieceLength {
			length = c.torrent.Info.PieceLength - offset
		}

		err := targetPeer.OutboundRequest(p, uint32(offset), uint32(length))
		if err != nil {
			logger.Error(err)
			continue
		}

		select {
		case chunk := <-targetPeer.Chunks:
			if chunk.Offset != offset {
				logger.Error(errors.New("chunk does not match expected chunk"))
				continue
			}

			for i, b := range chunk.Data {
				res[int64(i)+offset] = b
			}

			offset = offset + ChunkSize

			go c.GotChunk(result{
				piece:                p,
				chunkPositionInPiece: offset/ChunkSize - 1,
				bytes:                chunk.Data,
			})
			continue
		case <-time.After(chunkWaitTimeout * time.Second):
			if c.Peers[peerId] != nil {
				logger.Log(fmt.Sprintf("peer %s does not answer", c.Peers[peerId].Address))
			} else {
				logger.Log(fmt.Sprintf("peer has disconnected"))
			}
			continue
		}
	}

	if !p.Validate(res) {
		logger.Error(errors.New("piece does not match sha sum"))
	}

	return nil, res
}

func (c *Client) addPeers(response tracker.Response) error {
	var wg sync.WaitGroup

	wg.Add(len(response.Peers))

	for _, p := range response.Peers {
		go func(p tracker.Peer) {
			defer wg.Done()

			err, peer := peer.New(p.Ip, p.Port, c.torrent, func() { c.DeletePeer(fmt.Sprintf("%s:%d", p.Ip, p.Port)) })
			if err != nil {
				logger.Log(err.Error())
			} else {
				c.AddPeer(peer)
			}
		}(p)
	}

	wg.Wait()

	return nil
}

func (c *Client) AddPeer(p *peer.Peer) {
	if c.Peers[p.Address] != nil {
		return
	}

	logger.Log("adding peer : " + p.Address)

	if !p.Handshaked {
		err := p.OutboundHandshake(c.torrent.InfoHash, c.Id)
		if err != nil {
			logger.Log("could not handshake peer " + err.Error())
			c.DeletePeer(p.Address)
			return
		}
	}

	c.mapLock.Lock()
	c.Peers[p.Address] = p
	c.mapLock.Unlock()

	go c.ListenToPeerChan(p)

	return
}

func (c *Client) GotChunk(res result) {
	fmt.Printf("queueing chunk %d with current pointer %d\n", res.chunkPositionInPiece, c.latestChunk)
	send := func() {
		c.Results <- res
		c.latestChunk++
	}

	if res.chunkPositionInPiece > 0 && c.latestChunk < (res.chunkPositionInPiece-1) {
		for {
			if c.latestChunk == (res.chunkPositionInPiece - 1) {
				send()
				return
			}
		}
	} else {
		send()
	}

	return
}

func (c *Client) GotPiece(p *torrent.Piece, b []byte) {
	logger.Log(fmt.Sprintf("got full piece %d ", p.Index))

	c.finished[p.Index-c.torrent.Meta.SelectedPieces[0].Index] = true

	s := c.Tracker.State

	s.Downloaded += c.torrent.Info.PieceLength
	s.Left -= c.torrent.Info.PieceLength

	c.Tracker.Announce(&s)
	c.Have(p)

	return
}

func (c *Client) ListenToPeerChan(p *peer.Peer) {
	runtime.Gosched()

	check := func() {
		if p.CurrentJobs < 1 && p.Handshaken && p.Handshaked && !p.PeerChoking {
			c.AvailablePeers <- p.Address
		}
	}

	go func() {
		for {
			<-time.NewTicker(peerCheckInterval * time.Second).C
			check()
		}
	}()

	for msg := range p.Messages {
		switch msg.T {
		case "un_choke", "piece":
			check()
		}
	}
}

func (c *Client) DeletePeer(id string) {
	if c.Peers[id] == nil {
		return
	}

	logger.Log("deleting peer with id : " + id)

	c.mapLock.Lock()
	c.Peers[id].Conn.Close()
	c.Peers[id].Ticker.Stop()
	close(c.Peers[id].Messages)
	close(c.Peers[id].Chunks)
	delete(c.Peers, id)
	c.mapLock.Unlock()
	if len(c.Peers) == 0 {
		logger.Log("no moer peers")
	}
}

func (c *Client) Destroy() {
	c.listener.Close()

	for _, p := range c.Peers {
		c.DeletePeer(p.Address)
	}

	c.Tracker.Destroy()
}

func (c *Client) fatal(err error) {
	logger.Log(err.Error())
	c.Errors <- err
	close(c.Errors)
	c.listener.Close()
}
