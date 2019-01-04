package tracker

import (
	"github.com/Oted/torrent-video-stream/lib/client"
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

type instance struct {
	torrent *torrent.Torrent
	client  *client.Client
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

func Create(t *torrent.Torrent, c *client.Client) (error, instance) {
	cli := http.DefaultClient

	return nil, instance{
		torrent: t,
		client:  c,
		httpCli: cli,
	}
}

//func (t *instance) Announce(event string, ) {
//	http.NewRequest("GET")
//}
