package client

import (
	"fmt"
	"github.com/Oted/torrent-video-stream/lib/logger"
	"io"
)

func (c *Client) Seek(offset int64, whence int) (t int64, e error) {
	logger.Log(fmt.Sprintf("seek called for offset %d whence %d", offset, whence))

	switch whence {
	case 0:
		t = offset
	case 1:
		t = c.seek + offset
	case 2:
		t = c.torrent.SelectedFile().Length + offset
	}

	if c.seek != t {
		c.seek = t
		//TODO here we must flush the jobs and restart with the pieces containing and after the offset!
	}

	return
}

func (c *Client) Read(p []byte) (n int, err error) {
	if c.done {
		return 0, io.EOF
	}

	fmt.Printf("read called!\n")
	for {
		result := <-c.Results

		if result.piece == nil {
			return 0, io.EOF
		}

		//byteOffset :=
		//	((result.piece.Index - c.torrent.Meta.SelectedPieces[0].Index) * c.torrent.Info.PieceLength) +
		//	((result.chunkPositionInPiece - c.torrent.SelectedFile().Start/ChunkSize) * ChunkSize)

		for i, b := range result.bytes {
			index := int64(i)
			if index < 0 || index > c.torrent.SelectedFile().Length {
				continue
			}

			n++
			p[index] = b
			//fmt.Printf("writing byte %d to pos %d\n", b, index + byteOffset)
		}

		//logger.Log(fmt.Sprintf("sending chunk %d of piece %d with length %d and start index %d",
		//	result.chunkPositionInPiece,
		//	result.piece.Index,
		//	n,
		//	byteOffset))
		return
	}
}

