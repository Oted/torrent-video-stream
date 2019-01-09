package io

import (
	"fmt"
	"net"
)

//a connection listening for input on given port
func Listen(port int, cb func([]byte, net.Conn)) (error) {
	ioListener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	for {
		conn, err := ioListener.Accept()

		if err != nil {
			conn.Write([]byte(err.Error()))
			conn.Close()
		} else {
			b := make([]byte, 1024)
			length, err := conn.Read(b)

			if err != nil {
				conn.Write([]byte(err.Error()))
				conn.Close()
			} else {
				cb(b[0:length - 1], conn)
			}
		}
	}
}
