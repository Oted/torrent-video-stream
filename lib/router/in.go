package router

import "github.com/Oted/torrent-video-stream/lib/peer"

/*
0 - choke
1 - unchoke
2 - interested
3 - not interested
4 - have
5 - bitfield
6 - request
7 - piece
8 - cancel


19 - handshake?

'choke', 'unchoke', 'interested', and 'not interested' have no payload.
*/

func In(peer peer.Peer, data []byte) error {
	//t := data[0]

	//data
	//message

	return nil
}
