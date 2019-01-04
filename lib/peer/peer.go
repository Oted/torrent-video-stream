package peer

import "github.com/Oted/torrentz/lib/torrent"
/*
	0 - choke
	1 - unchoke
	2 - interested
	3 - not interested
	4 - have
	5 - bitfield
	6 - request
	7 - piece
	8 - cancel
*/

//a peer is created per connection established with another peer and has


type Peer struct {
	Pieces []*torrent.Piece
}
