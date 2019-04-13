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
	return nil, Tracker{
		torrent:  t,
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
	target := t.torrent.Announce

	if strings.Contains(target, "http") {
		for _, url := range t.torrent.AnnounceList {
			if strings.Split(url, ":")[0] == "udp" {
				target = url
				break
			}
		}
	}

	if s != nil {
		t.State = *s
	}

	protocol := strings.Split(target, ":")[0]


	switch protocol {
	case "http":
		return t.announceHttp(target)
	case "https":
		return t.announceHttp(target)
	case "udp":
		url := strings.Replace(strings.Replace(target, "udp://", "", 1), "/announce", "", 1)
		return t.announceUDP(url)
	}

	return errors.New("unsupported protocol " + protocol), nil
}

func (t *Tracker) Destroy() {
	if t.udpCli != nil {
		t.udpCli.Close()
	}
}

func newTransactionId() int32 {
	return rand.Int31()
}
