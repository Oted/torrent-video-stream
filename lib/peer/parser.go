package peer


type Request struct {
	T    string
	Data []byte
}

/*
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

func parse(data []byte) (error, *Request) {

}
