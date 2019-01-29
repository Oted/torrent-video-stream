package peer

import (
	"fmt"
	"net"
	"time"
)

type HandshakeRequest struct {
	Request
}

type Request struct {
	T    string
	Data []byte
}

type Peer struct {
	Address    string
	Port       uint16
	Ip         string
	Id         *string
	Received   []*Request
	Sent       []*Request
	KeepAlives int
	Handshaked bool
	conn       net.Conn
}

func New(ip string, port uint16) (error, *Peer) {
	address := fmt.Sprintf("%s:%d", ip, port)

	conn, err := net.DialTimeout("tcp", address, time.Second * 2)
	if err != nil {
		return err, nil
	}

	c := Peer{
		Address:    address,
		Port:       port,
		Ip:         ip,
		Id:         nil,
		Received:   make([]*Request, 10),
		Sent:       make([]*Request, 10),
		KeepAlives: 0,
		conn:       conn,
		Handshaked: false,
	}

	return nil, &c
}

func (c *Peer) Handshake(h HandshakeRequest) error {
	return nil
}

func (c *Peer) Receive(data []byte) error {
	//logger.Log("got data " + string(data) + "from cli " + c.Id)
	return nil
}

func (c *Peer) Send(data []byte) error {
	//logger.Log("sent data " + string(data) + "to cli " + c.Id)

	r := Request{}

	c.Sent = append(c.Sent, &r)

	return nil
}

func (c *Peer) keepAlive() {

}

func (c *Peer) Destroy() {
	c.conn.Close()
}
