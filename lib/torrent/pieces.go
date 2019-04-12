package torrent

import (
	"crypto/sha1"
	"hash"
	"reflect"
)

type Piece struct {
	Index      int
	RawSha     []byte
	Hash       hash.Hash
	Sum        [20]byte
	ByteOffset int64
	Chunks     []Chunk
}

type Chunk struct {
	Offset int64
	Done   bool
}

type Pieces []*Piece

func parsePieces(i interface{}) (pi Pieces) {
	s := reflect.ValueOf(i).String()
	index := 0

	for sliceFloor := 0; sliceFloor < len(s); sliceFloor += 20 {
		slice := []byte(s[sliceFloor:sliceFloor+20])
		sha := sha1.New()
		sha.Write(slice)

		p := Piece{
			Index:  index,
			RawSha: slice,
			Hash:   sha,
			Sum:    sha1.Sum(slice),
		}

		index++

		pi = append(pi, &p)
	}

	return
}
