package peer

import (
	"fmt"
	"github.com/Oted/torrent-video-stream/lib/torrent"
	"net"
	"time"
)

type Request struct {
	T    string
	Data []byte
}

type Peer struct {
	Address    string //ip:port
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

func (c *Peer) Receive(data []byte) error {
	//logger.Log("got data " + string(data) + "from cli " + c.Id)
	return nil
}

func (p *Peer) Choke() error {
	return nil
}

func (p *Peer) Unchoke() error {
	return nil
}

func (p *Peer) Interested(piece *torrent.Piece) error {
	//piece.
	return nil
}

func (p *Peer) NotInterested(piece *torrent.Piece) error {
	return nil
}

func (p *Peer) Have(piece *torrent.Piece) error {
	return nil
}

func (p *Peer) Bitfield() error {
	return nil
}

func (p *Peer) Request() error {
	return nil
}

func (p *Peer) Piece() error {
	return nil
}

func (p *Peer) Cancel() error {
	return nil
}

func (p *Peer) send(data []byte) error {
	//logger.Log("sent data " + string(data) + "to cli " + c.Id)

	r := Request{}

	p.Sent = append(p.Sent, &r)

	return nil
}

func (p *Peer) keepAlive() {

}

func (c *Peer) Destroy() {
	c.conn.Close()
}
