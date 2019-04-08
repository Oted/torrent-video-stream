package router

import (
	"github.com/Oted/torrent-video-stream/lib/peer"
	"net"
)

func Out2(conn net.Conn, data []byte) error {
	return nil
}

func Out1(peer peer.Peer, t int, msg []byte) error {
	return nil
}