package router

/*


func parse(data []byte) (error, *Request) {
	return nil,nil
}
	handshake: <pstrlen><pstr><reserved><info_hash><peer_id>
		pstrlen: string length of <pstr>, as a single raw byte
		pstr: string identifier of the protocol
		reserved: eight (8) reserved bytes. All current implementations use all zeroes. Each bit in these bytes can be used to change the behavior of the protocol. An email from Bram suggests that trailing bits should be used first, so that leading bits may be used to change the meaning of trailing bits.
		info_hash: 20-byte SHA1 hash of the info key in the metainfo file. This is the same info_hash that is transmitted in tracker requests.
		peer_id: 20-byte string used as a unique ID for the client. This is usually the same peer_id that is transmitted in tracker requests (but not always e.g. an anonymity option in Azureus).
	keep-alive: <len=0000>
	choke: <len=0001><id=0>
	unchoke: <len=0001><id=1>
	interested: <len=0001><id=2>
	have: <len=0005><id=4><piece index>
	bitfield: <len=0001+X><id=5><bitfield>
	request: <len=0013><id=6><index><begin><length>
		index: integer specifying the zero-based piece index
		begin: integer specifying the zero-based byte offset within the piece
		length: integer specifying the requested length.
	piece: <len=0009+X><id=7><index><begin><block>
		index: integer specifying the zero-based piece index
		begin: integer specifying the zero-based byte offset within the piece
		block: block of data, which is a subset of the piece specified by index.

	cancel: <len=0013><id=8><index><begin><length>
	port: <len=0003><id=9><listen-port>
 */
