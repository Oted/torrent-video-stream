package main

import (
	"errors"
	"github.com/Oted/torrent-video-stream/lib/client"
	"github.com/Oted/torrent-video-stream/lib/io"
	"github.com/Oted/torrent-video-stream/lib/logger"
	"github.com/Oted/torrent-video-stream/lib/torrent"
	"github.com/zeebo/bencode"
	"io/ioutil"
	"net"
	"os"
	"strconv"
)

type Input struct {
	P2P_IP   string
	P2P_PORT int
	IO_PORT  int
}

func main() {
	err, input := NewEnvs()
	if err != nil {
		panic(err)
	}

	err, ip := getOutboundIP()
	if err != nil {
		panic(err)
	}

	err = io.Listen(input.IO_PORT, func(message []byte, conn net.Conn) {
		//TODO there will be request specific chunks maybe
		logger.Log("recieved message " + string(message))
		defer conn.Close()

		err, torrent := torrentFromPath(string(message))
		if err != nil {
			conn.Write([]byte("\n" + err.Error()))
			return
		}

		err, client := client.New(ip, input.P2P_PORT, torrent)
		if err != nil {
			conn.Write([]byte("\n" + err.Error()))
			return
		}

		go client.Download()

		//we are reading until the end of the stream
		go func() {
			for data := range client.DataChannel {
				//TODO what happens on done?
				conn.Write(data)
			}
		}()

		//if there is an error then end
		for err := range client.Errors {
			conn.Write([]byte("\n" + err.Error()))
			client.Destroy()
		}
	})

	if err != nil {
		panic(err)
	}
}

//implement parser
func NewEnvs() (error, *Input) {
	ip := os.Getenv("P2P_IP")
	port1 := os.Getenv("P2P_PORT")
	port2 := os.Getenv("IO_PORT")

	if ip == "" || port1 == "" || port2 == "" {
		return errors.New("missing input envs"), nil
	}

	p2p, err := strconv.Atoi(port1)
	if err != nil {
		return err, nil
	}

	io, err := strconv.Atoi(port2)
	if err != nil {
		return err, nil
	}

	return nil, &Input{
		P2P_IP:   ip,
		P2P_PORT: p2p,
		IO_PORT:  io,
	}
}

func torrentFromPath(path string) (error, *torrent.Torrent) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err, nil
	}

	var data map[string]interface{}

	err = bencode.DecodeBytes(file, &data)
	if err != nil {
		return err, nil
	}

	err, torrent := torrent.Create(data)
	if err != nil {
		return err, nil
	}

	return nil, torrent
}


func getOutboundIP() (error, net.IP) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return err, nil
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return nil, localAddr.IP
}
