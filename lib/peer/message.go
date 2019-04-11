package peer

import (
	"encoding/binary"
	"fmt"
)

type Message struct {
	T    string
	Data []byte
}

func NewMessage(data []byte, isFirst bool) (err error, m Message) {
	if isFirst {
		err, m = parseHandshake(data)
	} else {
		err, m = parseMessage(data)
	}

	return nil, m
}

func parseMessage(d []byte) (error, Message) {
	return nil, Message{}
}

func parseHandshake(d []byte) (error, Message) {
	length := binary.BigEndian.Uint32(d[:1])
	fmt.Printf("got handshake length : %d\n", length)

	ident := binary.BigEndian.Uint32(d[1:length])
	fmt.Printf("got handshake identifier : %d\n", ident)

	sha := d[9+length : 20+length]
	fmt.Printf("got handshake sha : %s\n", sha)

	peerId := d[9+length : 20+length]
	fmt.Printf("got handshake sha : %s\n", peerId)

	return nil, Message{
		T:    "handshake",
		Data: d,
	}
}

func CreateHandshakeMessage(infoHash [20]byte, peer_id [20]byte) Message {
	var b [68]byte

	b[0] = 19

	copy(b[1:20], "BitTorrent protocol")
	copy(b[20:28], []byte{0,0,0,0,0,0,0,0})

	copy(b[28:48], infoHash[:])
	copy(b[48:68], peer_id[:])

	return Message{
		T:    "handshake",
		Data: b[:],
	}
}
