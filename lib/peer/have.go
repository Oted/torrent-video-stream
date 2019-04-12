package peer

import "encoding/binary"

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

func (p *Peer) OutboundHave(index int) (error) {
	err := p.Send(Message{
		T: "have",
		Data: []byte{0, 0, 0, 1, 2},
	})

	if err != nil {
		return err
	}

	//TODO maybe something with have

	return nil
}
