package peer

import (
	"errors"
)

/*
19 "handshake"
- "keep_alive"
0 "choke"
1 "un_choke"
2 "interested"
3 "not_interested"
4 "have"
5 "bitfield"
6 "request"
*/

type Message struct {
	T    string
	Data []byte
}

func (p *Peer) decideMessageType(b []byte) (error, string) {
	if b[0] == 19 {
		return nil, "handshake"
	}

	if len(b) < 5 {
		return nil, "keep_alive"
	}

	switch b[4] {
	case 0:
		return nil, "choke"
	case 1:
		return nil, "un_choke"
	case 2:
		return nil, "interested"
	case 3:
		return nil, "not_interested"
	case 4:
		return nil, "have"
	case 5:
		return nil, "bitfield"
	case 6:
		return nil, "request"
	case 7:
		return nil, "piece"
	case 8:
		return nil, "cancel"
	case 9:
		return nil, "port"
	}

	if p.Out == 1 && p.In == 1 && p.Handshaked {
		return nil, "headerless_bitfield"
	}

	return errors.New("invalid message from "), ""
}
