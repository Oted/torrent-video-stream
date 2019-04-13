package client

import (
	"errors"
	"fmt"
	"github.com/Oted/torrent-video-stream/lib/logger"
	"github.com/Oted/torrent-video-stream/lib/peer"
	"github.com/Oted/torrent-video-stream/lib/torrent"
	"github.com/Oted/torrent-video-stream/lib/tracker"
	"io"
	"net"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

//if piece empty or null then done
type result struct {
	piece *torrent.Piece
	bytes []byte
}

type Client struct {
	torrent        *torrent.Torrent
	MaxPeers       int
	Ip             string
	Id             [20]byte
	listener       net.Listener
	Peers          map[string]*peer.Peer //[ip:port]Peer
	Errors         chan error
	Tracker        tracker.Tracker
	AvailablePeers chan string
	Results        chan result
	Jobs           chan *torrent.Piece
	done           bool
	seek           int64
	finished       []bool
	mapLock        *sync.Mutex
}

const workers = 1

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
		torrent:        t,
		MaxPeers:       10,
		Ip:             ip,
		Id:             id,
		Tracker:        tracker,
		listener:       listener,
		Peers:          make(map[string]*peer.Peer),
		Errors:         make(chan error, 1),
		Jobs:           make(chan *torrent.Piece, len(t.Meta.SelectedPieces)),
		Results:        make(chan result, len(t.Meta.SelectedPieces)),
		done:           false,
		seek:           0,
		finished:       make([]bool, len(t.Meta.SelectedPieces)),
		mapLock:        &sync.Mutex{},
		AvailablePeers: make(chan string),
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

	c.startJobs()

	return nil
}

func (c *Client) startJobs() {
	for w := 1; w <= workers; w++ {
		go func(wor int) {
			for piece := range c.Jobs {
				err, res := c.getPiece(piece)

				if err != nil {
					c.fatal(err)
					return
				}

				//if the previous piece is not finished, must wait
				if piece.Index > c.torrent.Meta.SelectedPieces[0].Index && !c.finished[piece.Index-1-c.torrent.Meta.SelectedPieces[0].Index] {
					for {
						runtime.Gosched()

						if c.finished[piece.Index-1-c.torrent.Meta.SelectedPieces[0].Index] {
							c.Results <- result{
								piece: piece,
								bytes: res,
							}

							break
						}
					}
				} else {
					c.Results <- result{
						piece: piece,
						bytes: res,
					}
				}
			}
		}(w)
	}

	targetPieces := c.torrent.Meta.SelectedPieces

	for _, p := range targetPieces {
		c.Jobs <- p
	}

	close(c.Jobs)
}

func (c *Client) Read(p []byte) (n int, err error) {
	if c.done {
		return 0, io.EOF
	}

	//wait for results here
	for {
		result := <-c.Results

		if result.piece == nil {
			return 0, io.EOF
		}

		//this has to convert each pieces local byte slice to the global indexed slice
		for i, b := range result.bytes {
			//TODO here we can only return the bytes after the offset
			p[i] = b
		}

		c.finished[result.piece.Index-c.torrent.Meta.SelectedPieces[0].Index] = true
		return len(result.bytes), nil
	}
}

func (c *Client) Seek(offset int64, whence int) (t int64, e error) {
	logger.Log(fmt.Sprintf("seek called for offset %d whence %d", offset, whence))

	switch whence {
	case 0:
		t = offset
	case 1:
		t = c.seek + offset
	case 2:
		t = c.torrent.SelectedFile().Length + offset
	}

	if c.seek != t {
		c.seek = t
		//TODO here we must flush the jobs and restart with the pieces containing and after the offset!
	}

	return
}

//have to return full byte slice of piece
//which should match the torrent.info.piece.length
//except in the case of the last piece
func (c *Client) getPiece(p *torrent.Piece) (error, []byte) {
	res := make([]byte, c.torrent.Info.PieceLength)
	offset := uint32(0)

	runtime.Gosched()
	for peerId := range c.AvailablePeers {
		targetPeer := c.Peers[peerId]

		if targetPeer == nil {
			logger.Error(errors.New("non existing peer"))
			continue
		}

		if targetPeer.Has[uint32(p.Index)] == nil {
			logger.Log("peer does not have this piece")
			continue
		}

		err := targetPeer.OutboundRequest(p, uint32(offset))
		if err != nil {
			return err, nil
		}

		for chunk := range targetPeer.Chunks {
			fmt.Printf("GOT OUR FIRST CHUNK +%v\n", chunk)

			for i, b := range chunk.Data {
				res[uint32(i) + offset] = b
			}

			offset = offset + c.torrent.Meta.ChunkSize
			continue
		}
	}

	return nil, nil
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

//Peers calling us
func (c *Client) listen() {
	for {
		var p *peer.Peer

		logger.Log(fmt.Sprintf("listening for incomming messages on %s ", c.listener.Addr().String()))
		conn, err := c.listener.Accept()
		if err != nil {
			c.fatal(err)
			return
		}

		data := make([]byte, 131072) //2^17?

		len, err := conn.Read(data)
		if err != nil {
			c.fatal(err)
			return
		}

		//always process message
		defer func() {
			err := p.Recive(data[:len])
			if err != nil {
				c.fatal(err)
				return
			}
		}()

		addrs := strings.Split(conn.RemoteAddr().String(), ":")
		port, err := strconv.Atoi(addrs[1])
		if err != nil {
			c.fatal(err)
			return
		}

		p = c.Peers[addrs[0]+":"+addrs[1]]

		//is this an existing connection?
		if p != nil {
			return
		}

		//establish the connection with a dial
		err, p = peer.New(addrs[0], uint16(port), c.torrent, func() { c.DeletePeer(fmt.Sprintf("%s:%d", p.Ip, p.Port)) })
		if err != nil {
			c.fatal(err)
			return
		}

		c.AddPeer(p)
	}
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

func (c *Client) ListenToPeerChan(p *peer.Peer) {
	runtime.Gosched()

	for msg := range p.Messages {
		switch msg.T {
		case "un_choke", "piece" :
			if p.CurrentJobs == 0 && p.Handshaken && p.Handshaked && !p.PeerChoking {
				c.AvailablePeers <- p.Address
			}
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
		c.fatal(errors.New("no more peers left"))
	}
}

func (c *Client) Destroy() {
	c.listener.Close()

	for _, p := range c.Peers {
		p.Conn.Close()
	}

	c.Tracker.Destroy()
}

func (c *Client) fatal(err error) {
	logger.Log(err.Error())
	c.Errors <- err
	close(c.Errors)
	c.listener.Close()
}
