package peer

import "github.com/Oted/torrent-video-stream/lib/torrent"

//request: <len=0013><id=6><index><begin><length>

type request struct {}

func (p *Peer) InboundRequest(d []byte) (error, *request) {
	return nil, &request{}
}

func (p *Peer) OutboundRequest(piece *torrent.Piece, offset int64) (error) {
	var b [67]byte

	b[0] = 19

	copy(b[1:20], "BitTorrent protocol")
	copy(b[20:28], []byte{0, 0, 0, 6, 0, 0, 0, 0})


	err := p.Send(Message{
		T: "request",
		Data: b[:],
	})

	if err != nil {
		return err
	}

	return nil
}

