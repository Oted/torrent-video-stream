package peer

import (
	"bytes"
	"errors"
)

type handshake struct {
	ident    string
	infoHash [20]byte
	peerId   [20]byte
}

func (p *Peer) InboundHandshake(d []byte) (error, *handshake) {
	l := d[0]
	ident := string(d[1 : l+1])

	var infoHash [20]byte
	copy(infoHash[:], d[l+1+8 : l+1+20+8])

	if bytes.Compare(infoHash[:], p.torrent.InfoHash[:]) != 0 {
		return errors.New("infohash does not match"), nil
	}

	var peerId [20]byte
	copy(peerId[:], d[l+1+20+8 : l+1+40+8])


	p.Id = peerId
	p.Handshaken = true

	return nil, &handshake{
		ident:ident,
		infoHash: infoHash,
		peerId:peerId,
	}
}

func (p *Peer) OutboundHandshake(infoHash [20]byte, peer_id [20]byte) (error) {
	var b [68]byte

	b[0] = 19

	copy(b[1:20], "BitTorrent protocol")
	copy(b[20:28], []byte{0, 0, 0, 0, 0, 0, 0, 0})

	copy(b[28:48], infoHash[:])
	copy(b[48:68], peer_id[:])

	err := p.Send(Message{
		T: "handshake",
		Data: b[:],
	})

	if err != nil {
		return err
	}

	p.Handshaked = true

	return nil
}

