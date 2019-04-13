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

	//wait for results here
	for {
		result := <-c.Results

		if result.piece == nil {
			return 0, io.EOF
		}

		startPos := result.piece.Index * c.torrent.Info.PieceLength + (ChunkSize * result.chunkPositionInPiece) - c.selectedFileOffset

		for i, b := range result.bytes {
			if startPos+int64(i) < 0 {
				continue
			}

			p[startPos+int64(i)] = b
		}

		return len(result.bytes), nil
	}
}

