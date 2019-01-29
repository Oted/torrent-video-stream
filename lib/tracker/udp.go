package tracker

/*
 http://www.bittorrent.org/beps/bep_0015.html
 */

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

func (t *Tracker) announceUDP(url string) (error, *Response) {
	if t.udpCli == nil {
		err := t.handshakeUDP(url)
		if err != nil {
			return err, nil
		}
	}

	tid := newTransactionId()

	req := new(bytes.Buffer)
	binary.Write(req, binary.BigEndian, t.State.ConnectionId)
	binary.Write(req, binary.BigEndian, uint32(1))
	binary.Write(req, binary.BigEndian, uint32(tid))
	binary.Write(req, binary.BigEndian, t.State.InfoHash)
	binary.Write(req, binary.BigEndian, t.State.PeerId)
	binary.Write(req, binary.BigEndian, t.State.Downloaded)
	binary.Write(req, binary.BigEndian, t.State.Left)
	binary.Write(req, binary.BigEndian, t.State.Uploaded)

	var ev int32
	switch t.State.Event {
	case "started":
		ev = 2
	case "completed":
		ev = 1
	case "stopped":
		ev = 3
	default:
		ev = 0
	}

	binary.Write(req, binary.BigEndian, ev)
	binary.Write(req, binary.BigEndian, int32(0))
	binary.Write(req, binary.BigEndian, t.State.Key)
	if t.State.NumWant == nil {
		binary.Write(req, binary.BigEndian, int32(-1))
	} else {
		binary.Write(req, binary.BigEndian, t.State.NumWant)
	}
	binary.Write(req, binary.BigEndian, t.State.Port)

	_, err := t.udpCli.Write(req.Bytes())
	if err != nil {
		return err, nil
	}

	res := make([]byte, 1024)
	_, _, err = t.udpCli.ReadFrom(res)
	if err != nil {
		return err, nil
	}

	action := binary.BigEndian.Uint32(res[:4])
	transactionId := binary.BigEndian.Uint32(res[4:8])
	interval := binary.BigEndian.Uint32(res[8:12])
	leechers := binary.BigEndian.Uint32(res[12:16])
	seeers := binary.BigEndian.Uint32(res[16:20])
	peerData := res[20:]

	var peers []Peer
	n := 0
	for {
		ip := binary.BigEndian.Uint32(peerData[n : n+4])
		port := binary.BigEndian.Uint16(peerData[n+4 : n+6])

		if port < 1 {
			break
		}

		peers = append(peers, Peer{
			Ip: fmt.Sprintf("%d.%d.%d.%d",byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip)),
			Port: port,
		})

		n += 6
	}

	return nil, &Response{
		action,
		transactionId,
		interval,
		leechers,
		seeers,
		peers,
	}
}

func (t *Tracker) handshakeUDP(url string) (err error) {
	raddr, err := net.ResolveUDPAddr("udp", url)
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return err
	}

	tid := newTransactionId()

	req := new(bytes.Buffer)
	binary.Write(req, binary.BigEndian, uint64(0x41727101980))
	binary.Write(req, binary.BigEndian, uint32(0))
	binary.Write(req, binary.BigEndian, uint32(tid))

	if req.Len() != 16 {
		return errors.New("invalid length of request")
	}

	_, err = conn.Write(req.Bytes())
	if err != nil {
		return err
	}

	res := make([]byte, 16)
	_, _, err = conn.ReadFrom(res)
	if err != nil {
		return err
	}

	t.udpCli = conn
	connectionId := binary.BigEndian.Uint64(res[8:16])
	t.State.ConnectionId = &connectionId
	return
}
