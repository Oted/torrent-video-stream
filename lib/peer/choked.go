package peer

type choke struct {}

func (p *Peer) InboundChoke(d []byte) (error, *choke) {
	p.PeerChoking = true

	return nil, &choke{}
}

func (p *Peer) InboundUnChoke(d []byte) (error, *choke) {
	p.PeerChoking = false

	return nil, &choke{}
}

func (p *Peer) OutboundChoke() (error) {
	err := p.Send(Message{
		T: "choke",
		Data: []byte{0, 0, 0, 1, 0},
	})

	if err != nil {
		return err
	}

	p.AmChoking = true

	return nil
}

func (p *Peer) OutboundUnChoke() (error) {
	err := p.Send(Message{
		T: "un_choke",
		Data: []byte{0, 0, 0, 1, 1},
	})

	if err != nil {
		return err
	}

	p.AmChoking = false

	return nil
}