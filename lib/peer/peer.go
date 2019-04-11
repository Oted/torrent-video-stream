package peer

import (
	"fmt"
	"github.com/Oted/torrent-video-stream/lib/logger"
	"net"
	"time"
)

/*
0 - choke
1 - unchoke
2 - interested
3 - not interested
4 - have
5 - bitfield
6 - request
7 - piece
8 - cancel

19 - handshake?
*/

type Peer struct {
	Address    string //ip:port
	Port       uint16
	Ip         string
	Id         *string
	KeepAlives int
	conn       net.Conn
	listener   net.Listener
	Out        []*Message
	In         []*Message
	Handshaked bool //u sheked
	Handshaken bool //he sheked
}

func New(ip string, port uint16) (error, *Peer) {
	address := fmt.Sprintf("%s:%d", ip, port)

	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return err, nil
	}

	conn, err := net.DialTimeout("tcp", tcpAddr.String(), time.Second * 5)
	if err != nil {
		return err, nil
	}

	//try listening on the same port
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", tcpAddr.Port))
	if err != nil {
		return err, nil
	}

	p := Peer{
		Address:    address,
		Port:       port,
		Ip:         ip,
		Id:         nil,
		KeepAlives: 0,
		conn:       conn,
		listener: 	listener,
		Handshaked: false,
		Handshaken: false,
	}

	p.listen()

	return nil, &p
}

func (p *Peer) listen() {
	for {
		conn, err := p.listener.Accept()
		if err != nil {
			logger.Fatal(fmt.Sprintf("error listening on remote port : %s\n", err.Error()))
			return
		}

		data := make([]byte, 131072) //2^17?

		len, err := conn.Read(data)
		if err != nil {
			logger.Fatal(fmt.Sprintf("error reading data : %s\n", err.Error()))
			return
		}

		err = p.Recive(data[0 : len-1])
		if err != nil {
			logger.Fatal(fmt.Sprintf("error recieving : %s\n", err.Error()))
			return
		}
	}
}

func (p *Peer) Recive(b []byte) error {
	fmt.Printf("recieve message with length %d\n", len(b))
	err, m := NewMessage(b, len(p.In) > 0)
	if err != nil {
		return err
	}

	p.In = append(p.In, &m)
	return nil
}

func (p *Peer) Send(m Message) error {
	fmt.Printf("sending message %s with length %d\n", m.T, len(m.Data))
	_, err := p.conn.Write(m.Data)
	if err != nil {
		return err
	}

	if m.T == "handshake" {
		p.Handshaked = true
	}

	return nil
}

func (c *Peer) Destroy() {
	c.conn.Close()
}
