package tracker

import (
	"github.com/Oted/torrent-video-stream/lib/torrent"
	"net/http"
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
	httpCli *http.Client
}

type Info struct {
	PeerId     string
	Ip         *string
	Port       int
	Uploaded   int64
	Downloaded int64
	Left       int64
	Event      string
	NumWant    *string
	Key        *string
	TrackerId  *string
}

func Create(t *torrent.Torrent) (error, Tracker) {
	cli := http.DefaultClient

	return nil, Tracker{
		torrent: t,
		httpCli: cli,
	}
}

//func (t *Tracker) Announce(event string, ) {

//}