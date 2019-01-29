package torrent

import (
	"crypto/sha1"
	"hash"
	"reflect"
)

type Piece struct {
	Index        int
	Raw          []byte
	Hash         hash.Hash
	Sum          [20]byte
	Announced    bool
	AmChoked     bool //seed?
	AmInterested bool //leech?
	Have         bool
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
			Index:        index,
			Raw:          slice,
			Hash:         sha,
			Sum:          sha1.Sum(slice),
			Announced:    false,
			AmChoked:     true,
			AmInterested: false,
			Have:         false,
		}

		pi = append(pi, &p)
		index++
	}

	return
}
