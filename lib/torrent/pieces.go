package torrent

import (
	"bytes"
	"crypto/sha1"
	"hash"
	"reflect"
)

type Piece struct {
	Index      int64
	RawSha     []byte
	Hash       hash.Hash
	Sum        [20]byte
	ByteOffset int64
	Data       []byte
}

func (p *Piece) Validate(b []byte) bool {
	sum := sha1.Sum(b)

	return bytes.Compare(sum[:], p.Sum[:]) == 0
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
			Index:  int64(index),
			RawSha: slice,
			Hash:   sha,
			Sum:    sha1.Sum(slice),
		}

		index++

		pi = append(pi, &p)
	}

	return
}
