package client

import (
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/Oted/torrent-video-stream/lib/peer"
	"github.com/Oted/torrent-video-stream/lib/torrent"
	"github.com/Oted/torrent-video-stream/lib/tracker"
	"net"
	"os"
	"time"
)

type Client struct {
	Ip           net.IP
	Id           [20]byte
	listener     net.Listener
	peers        map[string]*peer.Peer //[ip:port]Peer
	Errors       chan error
	Tracker      tracker.Tracker
	DataChannel  chan []byte
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
		Ip:          ip,
		Id:          id,
		Tracker:     tracker,
		listener:    listener,
		peers:       make(map[string]*peer.Peer),
		Errors:      make(chan error, 1),
		DataChannel: make(chan []byte, torrent.Info.Files[torrent.Meta.VidIndex].Length),
	}

	return nil, &c
}

/*
	order from here:
	 - establish connection with tracker
	 - gather peer data
	 - listen for connections
	 - handshake with peers
	 - request first chunk of mov
	 - seed every bit after dl
	 - open

 */
func (c *Client) Download() {
	err, res := c.Tracker.Announce(nil)
	if err != nil {
		c.fatal(err)
		return
	}

	err = c.establishPeers(*res)
	if err != nil {
		c.fatal(err)
		return
	}

	err = c.listen()
	if err != nil {
		c.fatal(err)
		return
	}

	//we now have our peers and are listening, we can start handshake

}

func (c *Client) establishPeers(response tracker.Response) error {
	for _, p := range response.Peers {
		err, peer := peer.New(p.Ip, p.Port)
		if err != nil {
			return err
		}

		key := p.Ip + string(p.Port)

		if c.peers[key] != nil {
			continue
		}

		c.peers[key] = peer
	}

	return nil
}

func (c *Client) listen() error {
	for {
		//TODO needs routing
		conn, err := c.listener.Accept()
		if err != nil {
			return err
		}

		data := make([]byte, 131072) //2^17?

		len, err := conn.Read(data)
		if err != nil {
			return err
		}

		addr := conn.RemoteAddr().String()

		err, peer := peer.New(addr, conn)
		if err != nil {
			return err
		}

		c.AddPeer(addr, peer)

		err = peer.Receive(data[0 : len-1])
		if err != nil {
			return err
		}
	}
}

func (c *Client) AddPeer(id string, p *peer.Peer) {
	c.peers[id] = p
}

func (c *Client) DeletePeer(id string) {
	delete(c.peers, id)
}

func (c *Client) Destroy() {
	c.listener.Close()

	for _, p := range c.peers {
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
