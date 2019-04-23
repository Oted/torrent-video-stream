package peer

import (
	"encoding/binary"
	"github.com/Oted/torrent-video-stream/lib/torrent"
)

type have struct {
	index uint32
}

func (p *Peer) InboundHave(d []byte) (error, *have) {

	index := binary.BigEndian.Uint32(d[1:4])

	p.Has[index] = p.torrent.Info.Pieces[index]

	return nil, &have{
		index:index,
	}
}

func (p *Peer) OutboundHave(piece *torrent.Piece) (error) {
	var b [9]byte

	var prefix [4]byte
	binary.BigEndian.PutUint32(prefix[:], uint32(5))
	copy(b[0:4], prefix[:])

	b[4] = 4

	var index [4]byte
	binary.BigEndian.PutUint32(index[:], uint32(piece.Index))
	copy(b[5:9], index[:])

	err := p.Send(Message{
		T: "have",
		Data: b[:],
	})

	if err != nil {
		return err
	}

	p.Have[uint32(piece.Index)] = piece

	return nil
}
