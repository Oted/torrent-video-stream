package peer

import (
	"github.com/Oted/torrent-video-stream/lib/logger"
	"net"
)

type Request struct {
	T    string
	Data []byte
}

type Peer struct {
	Port       int16
	Ip         string
	Id         *string
	Received   []*Request
	Sent       []*Request
	KeepAlives int
	Handshaked bool
	conn       net.Conn
}

func New(ip string, port int16) (error, *Peer) {
	c := Peer{
		Port:       port,
		Ip:         ip,
		Id:         nil,
		Received:   make([]*Request, 10),
		Sent:       make([]*Request, 10),
		KeepAlives: 0,
		conn:       nil,
		Handshaked: false,
	}

	return nil, &c
}

func (c *Peer) Handshake() error {
	return nil
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

func (c *Peer) keepAlive() {

}

func (c *Peer) Destroy() {
	c.conn.Close()
}
