package client

import (
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"
)

func GetOutboundIP() (error, string) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return err, ""
	}
	defer conn.Close()

	localAddr := conn.RemoteAddr().(*net.UDPAddr)
	return nil, localAddr.String()
}

func GetPublicIP() (error, string) {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString("\n")
		os.Exit(1)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err, ""
	}

	return nil, string(b)
}

func generateId() [20]byte {
	now := time.Now().Unix()
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(now))
	binary.LittleEndian.PutUint64(b, uint64(os.Getpid()))

	return sha1.Sum(b)
}

//over the first 100 ports
func findAvailPort(start int) (error, net.Listener, int16) {
	for i := start; i <= start+100; i++ {
		//One connection per torrent.
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", i))
		if err == nil {
			return nil, listener, int16(i)
		}
	}

	return errors.New("no port avail"), nil, int16(-1)
}