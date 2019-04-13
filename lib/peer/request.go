package peer

import (
	"encoding/binary"
	"github.com/Oted/torrent-video-stream/lib/torrent"
)

//request: <len=0013><id=6><index><begin><length>

type request struct {}

func (p *Peer) InboundRequest(d []byte) (error, *request) {
	return nil, &request{}
}

func (p *Peer) OutboundRequest(piece *torrent.Piece, o uint32) (error) {
	//13 + 4
	var b [17]byte

	var prefix [4]byte
	binary.BigEndian.PutUint32(prefix[:], uint32(13))
	copy(b[0:4], prefix[:])

	b[4] = 6

	var index [4]byte
	binary.BigEndian.PutUint32(index[:], uint32(piece.Index))
	copy(b[5:8], index[:])

	var begin [4]byte
	binary.BigEndian.PutUint32(begin[:], o)
	copy(b[8:11], begin[:])

	var length [4]byte
	binary.BigEndian.PutUint32(length[:], p.torrent.Meta.ChunkSize)
	copy(b[11:], length[:])

	err := p.Send(Message{
		T: "request",
		Data: b[:],
	})

	if err != nil {
		return err
	}

	p.CurrentJobs++

	return nil
}

