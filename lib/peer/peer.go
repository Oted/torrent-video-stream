package peer

import (
	"github.com/Oted/torrent-video-stream/lib/logger"
	"net"
)

type Peer struct {
	Id         string
	Received   []*Request
	Sent       []*Request
	KeepAlives int
	conn       net.Conn
}

func New(id string, conn net.Conn) (error, *Peer) {
	c := Peer{
		Id:         id,
		Received:   make([]*Request, 10),
		Sent:       make([]*Request, 10),
		KeepAlives: 0,
		conn:       conn,
	}

	return nil, &c
}

func (c *Peer) Receive(data []byte) error {
	logger.Log("got data " + string(data) + "from cli " + c.Id)
	return nil
}

func (c *Peer) Send(data []byte) error {
	logger.Log("sent data " + string(data) + "to cli " + c.Id)
	return nil
}

func (c *Peer) keepAlive() {

}

func (c *Peer) Destroy() {
	c.conn.Close()
}
