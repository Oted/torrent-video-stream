package tracker

import (
	"errors"
	"fmt"
	"github.com/Oted/torrent-video-stream/lib/torrent"
	"io/ioutil"
	"net/http"
	"net/url"
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
	State
}

type State struct {
	InfoHash   string
	PeerId     string
	Ip         string
	Port       string
	Uploaded   int64
	Downloaded int64
	Left       int64
	Event      string
	NumWant    *string
	Key        *string
	TrackerId  *string
}

func Create(t *torrent.Torrent, ip string, port int, peerId [20]byte) (error, Tracker) {
	cli := http.DefaultClient

	peerEnc := url.QueryEscape(string(peerId[:20]))
	infoEnc := url.QueryEscape(string(t.InfoHash[:20]))

	return nil, Tracker{
		torrent: t,
		httpCli: cli,
		State: State{
			InfoHash:   infoEnc,
			PeerId:     peerEnc,
			Ip:         ip,
			Port:       fmt.Sprintf("%d", port),
			Uploaded:   0,
			Downloaded: 0,
			Left:       t.Info.Files[t.Meta.VidIndex].Length,
			Event:      "started",
			NumWant:    nil,
			Key:        nil,
			TrackerId:  nil,
		},
	}
}

func (t *Tracker) Announce() error {
	req, err := http.NewRequest("GET", t.torrent.Announce, nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("info_hash", t.State.InfoHash)
	q.Add("peer_id", t.State.PeerId)
	q.Add("ip", t.State.Ip)
	q.Add("port", t.State.Port)
	q.Add("downloaded", fmt.Sprintf("%d", t.State.Downloaded))
	q.Add("uploaded", fmt.Sprintf("%d", t.State.Uploaded))
	q.Add("left", fmt.Sprintf("%d", t.State.Left))

	req.URL.RawQuery = q.Encode()

	fmt.Println(req.URL.RawQuery)
	res, err := t.httpCli.Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))

	if res.StatusCode != 200 {
		return errors.New(fmt.Sprintf("%s : %d", "invalid statuscode from tracker", res.StatusCode))
	}

	return nil
}
