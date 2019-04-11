package client

import (
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
	"time"
)

//if piece empty or null then done
type result struct {
	piece *torrent.Piece
	bytes []byte
}

type Client struct {
	Torrent  *torrent.Torrent
	MaxPeers int
	Ip       string
	Id       [20]byte
	listener net.Listener
	Peers    map[string]*peer.Peer //[ip:port]Peer
	Errors   chan error
	Tracker  tracker.Tracker
	Results  chan result
	Jobs     chan *torrent.Piece
	done     bool
	seek     int64
	finished []bool
	mapLock  *sync.Mutex
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
		Torrent:  t,
		MaxPeers: 10,
		Ip:       ip,
		Id:       id,
		Tracker:  tracker,
		listener: listener,
		Peers:    make(map[string]*peer.Peer),
		Errors:   make(chan error, 1),
		Jobs:     make(chan *torrent.Piece, len(t.Meta.SelectedPieces)),
		Results:  make(chan result, len(t.Meta.SelectedPieces)),
		done:     false,
		seek:     0,
		finished: make([]bool, len(t.Meta.SelectedPieces)),
		mapLock:  &sync.Mutex{},
	}

	return nil, &c
}

func (c *Client) StartDownload() error {
	go c.listen()

	err, res := c.Tracker.Announce(nil)
	if err != nil {
		return err
	}

	err = c.addPeers(*res)
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
				if piece.Index > c.Torrent.Meta.SelectedPieces[0].Index && !c.finished[piece.Index-1-c.Torrent.Meta.SelectedPieces[0].Index] {
					for {
						runtime.Gosched()

						if c.finished[piece.Index-1-c.Torrent.Meta.SelectedPieces[0].Index] {
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

	targetPieces := c.Torrent.Meta.SelectedPieces

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

		c.finished[result.piece.Index-c.Torrent.Meta.SelectedPieces[0].Index] = true
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
		t = c.Torrent.SelectedFile().Length + offset
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
	time.Sleep(2 * time.Second)
	return nil, nil
}

func (c *Client) addPeers(response tracker.Response) error {
	var wg sync.WaitGroup

	wg.Add(len(response.Peers))

	for _, p := range response.Peers {
		go func(p tracker.Peer) {
			defer wg.Done()
			err, peer := peer.New(p.Ip, p.Port)
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

		logger.Log(fmt.Sprintf("listening for incomming messages on %s ",c.listener.Addr().String()))
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
			err := p.Recive(data[0 : len-1])
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
		err, p = peer.New(addrs[0], uint16(port))
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

	if !p.Handshaked {
		message := peer.CreateHandshakeMessage(c.Torrent.InfoHash, c.Id)
		err := p.Send(message)

		if err != nil {
			logger.Log("could not handshake peer " + err.Error())
			c.DeletePeer(p.Address)
			return
		}
	}

	c.mapLock.Lock()
	c.Peers[p.Address] = p
	c.mapLock.Unlock()
	return
}

func (c *Client) DeletePeer(id string) {
	c.mapLock.Lock()
	c.Peers[id].Destroy()
	c.mapLock.Unlock()
	delete(c.Peers, id)
}

func (c *Client) Destroy() {
	c.listener.Close()

	for _, p := range c.Peers {
		p.Destroy()
	}

	c.Tracker.Destroy()
}

func (c *Client) fatal(err error) {
	logger.Log(err.Error())
	c.Errors <- err
	close(c.Errors)
	c.listener.Close()
}
