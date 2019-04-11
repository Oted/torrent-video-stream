package client

import (
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/Oted/torrent-video-stream/lib/logger"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

func GetOutboundIP() (error, string) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return err, ""
	}
	defer conn.Close()

	remoteAddr := conn.LocalAddr().(*net.UDPAddr).String()

	for _, locals := range []string{"10.","192.168.","127.","169.254.","172.16.","224.",} {
		if strings.Index(remoteAddr, locals) == 0 {
			logger.Log("you seem to be behind NAT, cant listen to peers")
		}
	}

	return nil, remoteAddr
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