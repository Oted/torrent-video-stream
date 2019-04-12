package peer

import (
	"errors"
	"fmt"
	"github.com/Oted/torrent-video-stream/lib/logger"
	"github.com/Oted/torrent-video-stream/lib/torrent"
	"net"
	"time"
)

type Peer struct {
	torrent        *torrent.Torrent
	Address        string //ip:port
	Port           uint16
	Ip             string
	Id             [20]byte
	InKeepAlives   int
	OutKeepAlives  int
	Conn           net.Conn
	Out            int
	In             int
	Handshaked     bool
	Handshaken     bool
	Has            map[uint32]*torrent.Piece
	delete         func()
	AmChoking      bool
	AmInterested   bool
	PeerChoking    bool
	PeerInterested bool
	Ticker         *time.Ticker
	Messages       chan Message
}

const DialTimeout = 7
const KeepAliveInterval = 120

func New(ip string, port uint16, t *torrent.Torrent, delete func()) (error, *Peer) {
	address := fmt.Sprintf("%s:%d", ip, port)

	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return err, nil
	}

	conn, err := net.DialTimeout("tcp", tcpAddr.String(), time.Second*DialTimeout)
	if err != nil {
		return err, nil
	}

	ticker := time.NewTicker(time.Second * KeepAliveInterval)

	p := Peer{
		torrent:        t,
		Address:        address,
		Port:           port,
		Ip:             ip,
		InKeepAlives:   0,
		OutKeepAlives:  0,
		Conn:           conn,
		Handshaked:     false,
		Handshaken:     false,
		Out:            0,
		In:             0,
		delete:         delete,
		AmChoking:      true,
		AmInterested:   false,
		PeerChoking:    true,
		PeerInterested: false,
		Ticker:         ticker,
		Has:            make(map[uint32]*torrent.Piece),
		Messages:       make(chan Message),
	}

	go p.listen()
	go func() {
		for {
			<-ticker.C

			p.OutboundKeepAlive()
		}
	}()

	return nil, &p
}

func (p *Peer) listen() {
	data := make([]byte, 16384) //2^14?

	for {
		len, err := p.Conn.Read(data[:])
		if err != nil {
			logger.Error(err)
			p.delete()
			break
		} else {
			err = p.Recive(data[:len])
			if err != nil {
				logger.Error(err)
				return
			}
		}

	}
}

func (p *Peer) Recive(b []byte) error {
	if len(b) < 1 {
		return errors.New("invalid message")
	}

	err, t := decideMessageType(b)
	if err != nil {
		return err
	}

	logger.Log(fmt.Sprintf("recived message %s from %s ", t, p.Address))

	switch t {
	case "handshake":
		err, _ := p.InboundHandshake(b)
		if err != nil {
			p.delete()
			return err
		}
	case "choke":
		err, _ := p.InboundChoke(b)
		if err != nil {
			return err
		}
	case "un_choke":
		err, _ := p.InboundUnChoke(b)
		if err != nil {
			return err
		}
	case "interested":
		err, _ := p.InboundInterested(b)
		if err != nil {
			return err
		}
	case "not_interested":
		err, _ := p.InboundNotInterested(b)
		if err != nil {
			return err
		}
	case "have":
	case "bitfield":
		err, _ := p.InboundBitfield(b)
		if err != nil {
			p.delete()
			return err
		}

	case "request":
	case "piece":
	case "cancel":
	case "port":

	}

	p.Messages <- t

	p.In++
	return nil
}

func (p *Peer) Send(m Message) error {
	logger.Log(fmt.Sprintf("sending message %s to %s ", m.T, p.Address))
	_, err := p.Conn.Write(m.Data)
	if err != nil {
		return err
	}

	p.Out++

	return nil
}
