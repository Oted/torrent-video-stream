package peer

import (
	"github.com/Oted/torrent-video-stream/lib/logger"
	"net"
)

type Connection struct {
	Id         string
	Received   []*Request
	Sent       []*Request
	KeepAlives int
	Conn       net.Conn
}

func New(id string, conn net.Conn) (error, *Connection) {
	c := Connection{
		Id:         id,
		Received:   make([]*Request, 10),
		Sent:       make([]*Request, 10),
		KeepAlives: 0,
		Conn:       conn,
	}

	return nil, &c
}

func (c *Connection) Receive(data []byte) error {
	logger.Log("got data " + string(data) + "from cli " + c.Id)
	return nil
}

func (c *Connection) Send(data []byte) error {
	logger.Log("sent data " + string(data) + "to cli " + c.Id)
	return nil
}

func (c *Connection) keepAlive() {

}
