package main

import (
	"errors"
	"fmt"
	"github.com/Oted/torrent-video-stream/lib/client"
	"github.com/Oted/torrent-video-stream/lib/io"
	"github.com/Oted/torrent-video-stream/lib/torrent"
	"github.com/Oted/torrent-video-stream/lib/tracker"
	"github.com/zeebo/bencode"
	"io/ioutil"
	"os"
	"strconv"
)

type Input struct {
	P2P_IP   string
	P2P_PORT int
	IO_PORT  int
}

func main() {
	if len(os.Args) < 2 {
		panic(errors.New("no path"))
	}

	err, input := NewEnvs()
	if err != nil {
		panic(err)
	}

	err = io.Listen(input.IO_PORT, func(message []byte) {
		//defaults to path for now
		err, torrent := torrentFromPath()
		if err != nil {
			panic(err)
		}

		err, client := client.New(input.P2P_IP, input.IO_PORT)
		if err != nil {
			panic(err)
		}

		err, tracker := tracker.Create(torrent, client)
		if err != nil {
			panic(err)
		}

		fmt.Println(tracker)
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
