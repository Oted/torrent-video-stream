package client

import (
	"fmt"
	"github.com/Oted/torrent-video-stream/lib/peer"
	"github.com/Oted/torrent-video-stream/lib/torrent"
)

type BRes struct {
	Success int
	Fail int
}

func (c *Client) Want(data []byte) (error, *BRes) {
	for _, p := range c.Peers {
		fmt.Println(p)
	}


	return nil, nil
}

func (c *Client) Have(data []byte) (error, *BRes) {
	for _, p := range c.Peers {
		fmt.Println(p)
	}

	return nil, nil
}

func (c *Client) Choked(data []byte) (error, *BRes) {
	for _, p := range c.Peers {
		fmt.Println(p)
	}

	return nil, nil
}

func (c *Client) Interested(piece *torrent.Piece) (error, *BRes) {
	for _, p := range c.Peers {
		go func(p *peer.Peer) {
			//p.Interested(piece)
		}(p)
	}

	return nil, nil
}

