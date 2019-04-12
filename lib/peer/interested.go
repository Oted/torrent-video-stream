package peer

type interested struct {}

func (p *Peer) InboundInterested(d []byte) (error, *interested) {
	p.PeerInterested = true

	return nil, &interested{}
}

func (p *Peer) InboundNotInterested(d []byte) (error, *interested) {
	p.AmInterested = true

	return nil, &interested{}
}

func (p *Peer) OutboundInterested() (error) {
	err := p.Send(Message{
		T: "interested",
		Data: []byte{0, 0, 0, 1, 2},
	})

	if err != nil {
		return err
	}

	p.AmInterested = true

	return nil
}

func (p *Peer) OutboundNotInterested() (error) {
	err := p.Send(Message{
		T: "not_interested",
		Data: []byte{0, 0, 0, 1, 3},
	})

	if err != nil {
		return err
	}

	p.AmInterested = false

	return nil
}

