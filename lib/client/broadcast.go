package client

import (
	"fmt"
	"github.com/Oted/torrent-video-stream/lib/logger"
	"github.com/Oted/torrent-video-stream/lib/peer"
	"github.com/Oted/torrent-video-stream/lib/torrent"
	"sync"
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

func (c *Client) Have(piece *torrent.Piece) (error, *BRes) {
	for _, p := range c.Peers {
		p.OutboundHave(piece)
	}

	return nil, nil
}

func (c *Client) Choked(data []byte) (error, *BRes) {
	for _, p := range c.Peers {
		fmt.Println(p)
	}

	return nil, nil
}

func (c *Client) Interested() error {
	var wg sync.WaitGroup

	wg.Add(len(c.Peers))

	for _, p := range c.Peers {
		go func(p *peer.Peer) {
			defer wg.Done()

			err := p.OutboundInterested()
			if err != nil {
				logger.Error(err)
			}
		}(p)
	}

	wg.Wait()

	return nil
}

