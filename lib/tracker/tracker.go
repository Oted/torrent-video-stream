package tracker

import (
	"github.com/Oted/torrentz/lib/client"
	"github.com/Oted/torrentz/lib/torrent"
	"net/http"
)

/*
	info_hash
	peer_id
	ip
	port
	uploaded
	downloaded
	left
	event
*/

type instance struct {
	torrent *torrent.Torrent
	client  *client.Client
	httpCli *http.Client
}

type Info struct {
	PeerId     string
	Ip         string
	Port       int
	Uploaded   int64
	Downloaded int64
	Left       int64
	Event      string
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
