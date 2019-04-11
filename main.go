package main

import (
	"github.com/Oted/torrent-video-stream/lib/input"
	"os"
	"strconv"
)


func main() {
	//_ := os.Getenv("PEER_HOST")
	iH := os.Getenv("IO_HOST")
	pP := os.Getenv("PEER_PORT")
	iP := os.Getenv("IO_PORT")


	peerPort, err := strconv.Atoi(pP)
	if err != nil {
		panic(err)
	}

	ioPort, err := strconv.Atoi(iP)
	if err != nil {
		panic(err)
	}


	err, _ = input.NewIo(ioPort, iH, peerPort)

	if err != nil {
		panic(err)
	}
}
