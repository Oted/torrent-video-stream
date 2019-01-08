package client

import (
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/Oted/torrent-video-stream/lib/peer"
	"net"
	"os"
	"time"
)

//client is just another name for our peer instance

type Client struct {
	Id       [20]byte
	listener *net.Listener
	Errors   chan error
	Peers    map[string]*peer.Connection
}

func New(ip string, startPort int) (error, *Client) {
	err, port := findAvailPort(startPort)
	if err != nil {
		return err, nil
	}

	//One connection per torrent.
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		panic(err)
	}

	c := Client{
		Id:       generateId(),
		listener: &listener,
		Peers:    make(map[string]*peer.Connection),
		Errors:   make(chan error),
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err, nil
		}

		go c.handleRequest(conn)
	}

	return nil, &c
}

//loog over the first 100 ports
func findAvailPort(start int) (error, int) {
	for i := start; i <= start + 100; i++ {
		conn, _ := net.DialTimeout("tcp", string(start), 20)
		if conn != nil {
			return nil, start
		}
	}

	return errors.New("no port avail"), 0
}

func (c *Client) handleRequest(conn net.Conn) error {
	data := make([]byte, 1024)

	_, err := conn.Read(data)
	if err != nil {
		return err
	}

	addr := conn.RemoteAddr().String()

	err, peer := peer.New(addr, conn)
	if err != nil {
		return err
	}

	err = peer.Receive(data)
	if err != nil {
		return err
	}

	c.AddPeer(addr, peer)
	return nil
}

func (c *Client) AddPeer(id string, p *peer.Connection) {
	c.Peers[id] = p
}

func (c *Client) DeletePeer(id string) {
	delete(c.Peers, id)
}

func generateId() [20]byte {
	now := time.Now().Unix()
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(now))
	binary.LittleEndian.PutUint64(b, uint64(os.Getpid()))

	return sha1.Sum(b)
}
