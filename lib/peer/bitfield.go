package peer

import (
	"encoding/binary"
	"errors"
	"strconv"
)

type bitfield struct {
}

func (p *Peer) InboundBitfield(d []byte) (error, *bitfield) {
	l := binary.BigEndian.Uint32(d[0:4])

	bf := make([]byte, l)
	copy(bf, d[5 : 5 + l - 1])

	index := 0
	for _, b := range bf {
		str := strconv.FormatInt(int64(b), 2)
		for _, c := range str {
			if index >= len(p.torrent.Info.Pieces) {
				index++
				continue
			}

			if c == 49 {
				p.Has[uint32(index)] = p.torrent.Info.Pieces[index]
			} else if c != 48 {
				return errors.New("invalid bitfield char recieved " + string(c)), nil
			}

			index++
		}
	}

	return nil, &bitfield{}
}

func (p *Peer) OutboundBitfield(infoHash [20]byte, peer_id [20]byte) (error) {
	return nil
}

