package client

import (
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"os"
	"time"
)

func GetOutboundIP() (error, net.IP) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return err, nil
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return nil, localAddr.IP
}

func generateId() [20]byte {
	now := time.Now().Unix()
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(now))
	binary.LittleEndian.PutUint64(b, uint64(os.Getpid()))

	return sha1.Sum(b)
}

//over the first 100 ports
func findAvailPort(ip string, start int) (error, net.Listener, int16) {
	for i := start; i <= start+100; i++ {
		//One connection per torrent.
		listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, i))
		if err == nil {
			return nil, listener, int16(i)
		}
	}

	return errors.New("no port avail"), nil, int16(-1)
}