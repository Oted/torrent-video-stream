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
	Handshaked bool
	conn       net.Conn
}

func New(id string, conn net.Conn) (error, *Peer) {
	c := Peer{
		Id:         id,
		Received:   make([]*Request, 10),
		Sent:       make([]*Request, 10),
		KeepAlives: 0,
		conn:       conn,
		Handshaked: false,
	}

	return nil, &c
}

func (c *Peer) Receive(data []byte) error {
	logger.Log("got data " + string(data) + "from cli " + c.Id)
	return nil
}

func (c *Peer) Send(data []byte) error {
	logger.Log("sent data " + string(data) + "to cli " + c.Id)

	r := Request{}

	c.Sent = append(c.Sent, &r)

	return nil
}

//infohash,
func (c *Peer) Handshake() error {
	//c.Send()

	c.Handshaked = true
	return nil
}

func (c *Peer) keepAlive() {

}

func (c *Peer) Destroy() {
	c.conn.Close()
}
