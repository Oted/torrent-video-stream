package input

import (
	"fmt"
	"github.com/Oted/torrent-video-stream/lib/client"
	"github.com/Oted/torrent-video-stream/lib/logger"
	"net/http"
	"time"
)

type input struct {
	peerIp   string
	peerPort int
	ioHost   string
	ioPort   int
}

func NewIo(io_port int, io_host string, peer_port int) (error, *input) {
	err, peerIp := client.GetOutboundIP()
	if err != nil {
		return err, nil
	}

	s := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", io_host, io_port),
		WriteTimeout: time.Second * 60,
		ReadTimeout:  time.Second * 60,
	}

	i := input{
		ioHost:   "0.0.0.0",
		ioPort:   io_port,
		peerIp:   peerIp,
		peerPort: peer_port,
	}

	http.HandleFunc("/", i.handler)

	if err = s.ListenAndServe(); err != nil {
		return err, nil
	}

	return nil, &i
}

func (i *input) handler(w http.ResponseWriter, r *http.Request) {
	input := r.URL.Path[1:]

	if input[:6] == "magnet" {

	}

	logger.Log("received message " + string(input))

	err, torrent := torrentFromPath(string(input))
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("\n" + err.Error()))
		return
	}

	err, client := client.New(i.peerIp, i.peerPort, torrent)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("\n" + err.Error()))
		return
	}

	err = client.StartDownload()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("\n" + err.Error()))
		client.Destroy()
		return
	}

	//client.read and client.seek will be called until EOF
	http.ServeContent(w, r, torrent.Meta.TargetFileName, time.Time{}, client)

	//if there is an error then end
	for err := range client.Errors {
		w.Write([]byte("\n" + err.Error()))
		w.WriteHeader(500)
		client.Destroy()
	}
}
