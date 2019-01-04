package client

import (
	"crypto/sha1"
	"encoding/binary"
	"net"
	"os"
	"time"
)

//client is just another name for our peer instance

type Client struct {
	Id     [20]byte
	client *net.TCPListener
}

func New() (error, *Client) {

	net.ListenTCP()

	return nil, &Client{
		Id: generateId(),
	}
}

func generateId() [20]byte {
	now := time.Now().Unix()
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(now))
	binary.LittleEndian.PutUint64(b, uint64(os.Getpid()))

	return sha1.Sum(b)
}
