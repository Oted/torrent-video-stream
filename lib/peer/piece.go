package peer

import (
	"encoding/binary"
)

type piece struct {
	index uint32
	begin uint32
	block []byte
}

func (p *Peer) InboundPiece(d []byte) (error, *piece) {
	l := binary.BigEndian.Uint32(d[0:4])
	bf := make([]byte, l - 9)

	index := binary.BigEndian.Uint32(d[5:9])

	begin := binary.BigEndian.Uint32(d[9:13])

	copy(bf, d[13:])

	return nil, &piece{
		index: index,
		begin: begin,
		block: bf[:],
	}
}

func (p *Peer) OutboundPiece(ownedPieceChunk *piece, isFirst bool) (err error) {
	b :=  make([]byte, p.chunkSize + 9 + 4)

	var prefix [4]byte
	binary.BigEndian.PutUint32(prefix[:], uint32(p.chunkSize + 9))
	copy(b[0:4], prefix[:])

	b[4] = 4

	var index [4]byte
	binary.BigEndian.PutUint32(index[:], uint32(ownedPieceChunk.index))
	copy(b[5:9], index[:])

	var begin [4]byte
	binary.BigEndian.PutUint32(index[:], uint32(ownedPieceChunk.begin))
	copy(b[5:9], begin[:])

	if isFirst {
		err = p.Send(Message{
			T: "piece",
			Data: b[:],
		})
	} else {
		err = p.Send(Message{
			T: "piece",
			Data: b[:],
		})
	}

	if err != nil {
		return err
	}

	return nil}
