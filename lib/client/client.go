package client

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Oted/torrent-video-stream/lib/logger"
	"github.com/Oted/torrent-video-stream/lib/peer"
	"github.com/Oted/torrent-video-stream/lib/router"
	"github.com/Oted/torrent-video-stream/lib/torrent"
	"github.com/Oted/torrent-video-stream/lib/tracker"
	"io"
	"net"
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
	Torrent    *torrent.Torrent
	MaxPeers   int
	Ip         net.IP
	Id         [20]byte
	listener   net.Listener
	Peers      map[string]*peer.Peer //[ip:port]Peer
	Errors     chan error
	Tracker    tracker.Tracker
	Results    chan result
	Jobs       chan *torrent.Piece
	done       bool
	wg         sync.WaitGroup
	seek       int64
	finished   []bool
	waitBuffer *bytes.Buffer
}

const workers = 1

func New(ip net.IP, startPort int, t *torrent.Torrent) (error, *Client) {
	err, listener, port := findAvailPort(ip.String(), startPort)
	if err != nil {
		return err, nil
	}

	id := generateId()

	err, tracker := tracker.Create(t, ip.String(), port, id)
	if err != nil {
		return err, nil
	}

	c := Client{
		Torrent:    t,
		MaxPeers:   10,
		Ip:         ip,
		Id:         id,
		Tracker:    tracker,
		listener:   listener,
		Peers:      make(map[string]*peer.Peer),
		Errors:     make(chan error, 1),
		Jobs:       make(chan *torrent.Piece, len(t.SelectPieces())),
		Results:    make(chan result, len(t.SelectPieces())),
		done:       false,
		seek:       0,
		finished:   make([]bool, len(t.SelectPieces())),
		waitBuffer: bytes.NewBuffer(nil),
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

	//start the workers, depending on how many, that decides the amount of concurrent pieces
	for w := 1; w <= workers; w++ {
		go func() {
			for piece := range c.Jobs {
				err, res := c.getPiece(piece)

				if err != nil {
					c.fatal(err)
					return
				}

				//if the previous piece is not finished, must wait
				if piece.Index > 0 && !c.finished[piece.Index - 1] {
					for {
						if c.finished[piece.Index - 1] {
							c.Results <- result{
								piece: piece,
								bytes: res,
							}
						}
					}
				} else {
					c.Results <- result{
						piece: piece,
						bytes: res,
					}
				}

			}
		}()
	}

	targetPieces := c.Torrent.SelectPieces()

	for _, p := range targetPieces {
		c.Jobs <- p
	}

	close(c.Jobs)

	return nil
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
			p[i] = b
		}

		c.finished[result.piece.Index] = true
		return len(result.bytes), nil
	}
}

func (c *Client) Seek(offset int64, whence int) (int64, error) {

}

//have to return full byte slice of piece
//which should match the torrent.info.piece.length
//except in the case of the last piece
func (c *Client) getPiece(p *torrent.Piece) (error, []byte) {
	fmt.Println("getting piece")

	//c.Broadcast()

	fmt.Println(p)
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

	if len(c.Peers) < 1 {
		return errors.New("no peer connections")
	}

	return nil
}

//Peers calling us
func (c *Client) listen() {
	for {
		var p *peer.Peer

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
			router.In(*p, data[0:len-1])
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

		//add the peer
		err = c.AddPeer(p)
		if err != nil {
			c.fatal(err)
			return
		}
	}
}

func (c *Client) AddPeer(p *peer.Peer) error {
	if c.Peers[p.Address] != nil {
		return nil
	}

	c.Peers[p.Address] = p
	return nil
}

func (c *Client) DeletePeer(id string) {
	c.Peers[id].Destroy()
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
}
