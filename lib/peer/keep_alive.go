package peer

type keep_alive struct {}

func (p *Peer) InboundKeepAlive(d []byte) (error, *keep_alive) {

	p.InKeepAlives++

	return nil, &keep_alive{}
}

func (p *Peer) OutboundKeepAlive() (error) {
	err := p.Send(Message{
		T: "keep_alive",
		Data: []byte{0,0,0,0},
	})

	if err != nil {
		return err
	}

	p.OutKeepAlives++

	return nil
}

