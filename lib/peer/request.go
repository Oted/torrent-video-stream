package peer

import (
	"encoding/binary"
	"github.com/Oted/torrent-video-stream/lib/torrent"
)

//request: <len=0013><id=6><index><begin><length>

type request struct {
	pre    uint32
	index  uint32
	begin  uint32
	length uint32
}

func (p *Peer) InboundRequest(d []byte) (error, *request) {
	//13 + 4
	var b [17]byte

	pre := binary.BigEndian.Uint32(d[0:4])
	index := binary.BigEndian.Uint32(b[5:9])
	begin := binary.BigEndian.Uint32(b[9:13])
	length := binary.BigEndian.Uint32(b[13:17])

	return nil, &request{
		pre: pre,
		index: index,
		begin: begin,
		length: length,
	}
}

func (p *Peer) OutboundRequest(piece *torrent.Piece, o uint32, chunkSize uint32) (error) {
	//13 + 4
	var b [17]byte

	var prefix [4]byte
	binary.BigEndian.PutUint32(prefix[:], uint32(13))
	copy(b[0:4], prefix[:])

	b[4] = 6

	var index [4]byte
	binary.BigEndian.PutUint32(index[:], uint32(piece.Index))
	copy(b[5:9], index[:])

	var begin [4]byte
	binary.BigEndian.PutUint32(begin[:], o)
	copy(b[9:13], begin[:])

	var length [4]byte
	binary.BigEndian.PutUint32(length[:], chunkSize)
	copy(b[13:17], length[:])

	err := p.Send(Message{
		T:    "request",
		Data: b[:],
	})

	if err != nil {
		return err
	}

	p.CurrentJobs++

	return nil
}
