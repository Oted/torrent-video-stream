package client

import (
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/Oted/torrent-video-stream/lib/download"
	"github.com/Oted/torrent-video-stream/lib/logger"
	"github.com/Oted/torrent-video-stream/lib/peer"
	"github.com/Oted/torrent-video-stream/lib/router"
	"github.com/Oted/torrent-video-stream/lib/torrent"
	"github.com/Oted/torrent-video-stream/lib/tracker"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Client struct {
	Torrent     *torrent.Torrent
	MaxPeers    int
	Ip          net.IP
	Id          [20]byte
	listener    net.Listener
	Peers       map[string]*peer.Peer //[ip:port]Peer
	Errors      chan error
	Tracker     tracker.Tracker
	DataChannel chan []byte
	Jobs        chan job
	sync.WaitGroup
}

type job struct {
}

func New(ip net.IP, startPort int, torrent *torrent.Torrent) (error, *Client) {
	err, listener, port := findAvailPort(ip.String(), startPort)
	if err != nil {
		return err, nil
	}

	id := generateId()

	err, tracker := tracker.Create(torrent, ip.String(), port, id)
	if err != nil {
		return err, nil
	}

	c := Client{
		Torrent:     torrent,
		MaxPeers:    10,
		Ip:          ip,
		Id:          id,
		Tracker:     tracker,
		listener:    listener,
		Peers:       make(map[string]*peer.Peer),
		Errors:      make(chan error, 1),
		Jobs:        make(chan job, 100),
		DataChannel: make(chan []byte, torrent.SelectVideoFile().Length),
	}

	return nil, &c
}

func (c *Client) Start() {
	go c.listen()

	err, res := c.Tracker.Announce(nil)
	if err != nil {
		c.fatal(err)
		return
	}

	err = c.addPeers(*res)
	if err != nil {
		c.fatal(err)
		return
	}

	err, targetPieces := c.Torrent.SelectVideoPieces()
	if err != nil {
		c.fatal(err)
		return
	}

	for _, p := range targetPieces {
		go download.Download(p)
	}
}

func (c *Client) download(p *torrent.Piece) {


}

func (c *Client) addPeers(response tracker.Response) error {
	var wg sync.WaitGroup

	wg.Add(len(response.Peers))

	for _, p := range response.Peers {
		go func(p tracker.Peer) {
			defer wg.Done()

			if len(c.Peers) > c.MaxPeers {
				return
			}

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

func (c *Client) listen() {
	//this cant return really...
	for {
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

		addrs := strings.Split(conn.RemoteAddr().String(), ":")
		port, err := strconv.Atoi(addrs[1])
		if err != nil {
			c.fatal(err)
			return
		}

		err, peer := peer.New(addrs[0], uint16(port))
		if err != nil {
			c.fatal(err)
			return
		}

		err = c.AddPeer(peer)
		if err != nil {
			c.fatal(err)
			return
		}

		router.In(*peer, data[0:len-1])
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

//message goes out to all peers
func (c *Client) Broadcast(data []byte) {
	for _, p := range c.Peers {
		fmt.Println(p)
	}
}

func (c *Client) Destroy() {
	c.listener.Close()

	for _, p := range c.Peers {
		p.Destroy()
	}

	c.Tracker.Destroy()
}

func (c *Client) fatal(err error) {
	c.Errors <- err
	close(c.Errors)
}

//over the first 100 ports
func findAvailPort(ip string, start int) (error, net.Listener, int16) {
	for i := start; i <= start+100; i++ {
		//One connection per torrent.
		listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, i))
		if err == nil {
			return nil, listener, int16(i)
		}
	}

	return errors.New("no port avail"), nil, int16(-1)
}

func generateId() [20]byte {
	now := time.Now().Unix()
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(now))
	binary.LittleEndian.PutUint64(b, uint64(os.Getpid()))

	return sha1.Sum(b)
}
