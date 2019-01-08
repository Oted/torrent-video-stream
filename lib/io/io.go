package io

import (
	"fmt"
	"net"
)

//a connection listening for input on given port
func Listen(port int, cb func([]byte)) (error) {
	ioListener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	for {
		conn, err := ioListener.Accept()
		if err != nil {
			conn.Write([]byte(err.Error()))
		} else {
			b := make([]byte, 128)
			_, err := conn.Read(b)

			if err != nil {
				conn.Write([]byte(err.Error()))
			} else {
				cb(b)
			}
		}
	}
}
