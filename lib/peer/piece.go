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

func (p *Peer) OutboundPiece() (error) {
	return nil
}
