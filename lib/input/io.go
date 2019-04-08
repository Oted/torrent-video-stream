package input

import (
	"fmt"
	"github.com/Oted/torrent-video-stream/lib/client"
	"github.com/Oted/torrent-video-stream/lib/logger"
	"net"
	"net/http"
	"time"
)

type input struct {
	peerIp   net.IP
	peerPort int
	ioHost   string
	ioPort   int
	now      time.Time
}

func NewIo(io_port int, io_host string, peer_port int) (error, *input) {
	err, peerIp := client.GetOutboundIP()
	if err != nil {
		return err, nil
	}

	s := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", io_host, io_port),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	i := input{
		ioHost:   "0.0.0.0",
		ioPort:   io_port,
		peerIp:   peerIp,
		peerPort: peer_port,
		now:      time.Now(),
	}

	http.HandleFunc("/", i.handler)

	if err = s.ListenAndServe(); err != nil {
		return err, nil
	}

	return nil, &i
}

//TODO there will be request specific chunks maybe
func (i *input) handler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path

	logger.Log("recieved message " + string(filePath))

	err, torrent := torrentFromPath(string(filePath))
	if err != nil {
		w.Write([]byte("\n" + err.Error()))
		return
	}

	err, client := client.New(i.peerIp, i.peerPort, torrent)
	if err != nil {
		w.Write([]byte("\n" + err.Error()))
		return
	}

	err = client.StartDownload()
	if err != nil {
		w.Write([]byte("\n" + err.Error()))
		return
	}

	//client.read and client.seek will be called until EOF
	http.ServeContent(w, r, client.Torrent.Meta.TargetFile, i.now, client)

	//if there is an error then end
	for err := range client.Errors {
		w.Write([]byte("\n" + err.Error()))
		client.Destroy()
	}
}
