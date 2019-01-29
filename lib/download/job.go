package download

import (
	"github.com/Oted/torrent-video-stream/lib/peer"
	"github.com/Oted/torrent-video-stream/lib/torrent"
)

type job struct {
	piece *torrent.Piece
	Data  chan []byte
	Error chan error
	peer  *peer.Peer
}

func Download(piece *torrent.Piece, peer *peer.Peer) job {
	j := job{
		piece: piece,
		peer:  peer,
		Data:  make(chan []byte),
		Error: make(chan error),
	}

	return j
}
