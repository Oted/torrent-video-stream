package tracker

import (
	"errors"
	"github.com/Oted/torrent-video-stream/lib/torrent"
	"math/rand"
	"net"
	"strings"
)

/*
	req
		info_hash
		peer_id
		ip
		port
		uploaded
		downloaded
		left
		event //started, stopped, completed
	res
		failure reason:
		warning message:
		interval:
		min interval:
		tracker id:
		complete:
		incomplete:
		peers:
		peer id:
		ip:
		port:
		peers:

 */

type Tracker struct {
	torrent *torrent.Torrent
	udpCli  *net.UDPConn

	State    State
	Replies  []*Response
	Protocol string
}

type Response struct {
	action        uint32
	transactionId uint32
	interval      uint32
	leechers      uint32
	seeers        uint32
	Peers         []Peer
}

type Peer struct {
	Ip   string
	Port uint16
}

type State struct {
	InfoHash     [20]byte
	PeerId       [20]byte
	Ip           string
	Port         int16
	Uploaded     int64
	Downloaded   int64
	Left         int64
	Event        string
	Key          int32
	NumWant      *int32
	TrackerId    *string
	ConnectionId *uint64
}

func Create(t *torrent.Torrent, ip string, port int16, peerId [20]byte) (error, Tracker) {
	protocol := strings.Split(t.Announce, ":")[0]

	return nil, Tracker{
		torrent:  t,
		Protocol: protocol,
		State: State{ //default state
			InfoHash:     t.InfoHash,
			PeerId:       peerId,
			Ip:           ip,
			Port:         port,
			Uploaded:     0,
			Downloaded:   0,
			Left:         t.Info.Files[t.Meta.TargetIndex].Length,
			Event:        "started",
			NumWant:      nil,
			Key:          0,
			TrackerId:    nil,
			ConnectionId: nil,
		},
	}
}

func (t *Tracker) Announce(s *State) (error, *Response) {
	if s != nil {
		t.State = *s
	}

	switch t.Protocol {
	case "http":
		return t.announceHttp(t.torrent.Announce)
	case "https":
		return t.announceHttp(t.torrent.Announce)
	case "udp":
		url := strings.Replace(strings.Replace(t.torrent.Announce, "udp://", "", 1), "/announce", "", 1)
		return t.announceUDP(url)
	}

	return errors.New("unsupported protocol " + t.Protocol), nil
}

func (t *Tracker) Destroy() {
	if t.udpCli != nil {
		t.udpCli.Close()
	}
}

func newTransactionId() int32 {
	return rand.Int31()
}
