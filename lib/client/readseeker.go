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
	}

	c.seek = t

	return
}

func (c *Client) Read(p []byte) (n int, err error) {
	if c.done {
		return 0, io.EOF
	}

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
		}

		return
	}
}

