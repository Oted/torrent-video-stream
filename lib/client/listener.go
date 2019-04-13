package client

import (
	"fmt"
	"github.com/Oted/torrent-video-stream/lib/logger"
	"github.com/Oted/torrent-video-stream/lib/peer"
	"strconv"
	"strings"
)

//Peers calling us
func (c *Client) listen() {
	for {
		var p *peer.Peer

		logger.Log(fmt.Sprintf("listening for incomming messages on %s ", c.listener.Addr().String()))
		conn, err := c.listener.Accept()
		if err != nil {
			c.fatal(err)
			return
		}

		data := make([]byte, 131072) //2^17?

		len, err := conn.Read(data)
		if err != nil {
			c.fatal(err)
			return
		}

		//always process message
		defer func() {
			err := p.Recive(data[:len])
			if err != nil {
				c.fatal(err)
				return
			}
		}()

		addrs := strings.Split(conn.RemoteAddr().String(), ":")
		port, err := strconv.Atoi(addrs[1])
		if err != nil {
			c.fatal(err)
			return
		}

		p = c.Peers[addrs[0]+":"+addrs[1]]

		//is this an existing connection?
		if p != nil {
			return
		}

		//establish the connection with a dial
		err, p = peer.New(addrs[0], uint16(port), c.torrent, func() { c.DeletePeer(fmt.Sprintf("%s:%d", p.Ip, p.Port)) })
		if err != nil {
			c.fatal(err)
			return
		}

		c.AddPeer(p)
	}
}

