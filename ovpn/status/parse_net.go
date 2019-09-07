package status

import (
	"bytes"
	"io"
	"log"
	"net"
	"time"
)

const (
	bufferSize = 512
)

// expect "for more info"
// send "status ${NUM}\n"
// expect "END"
// send "quit\n"
// expect eof

// ParseAddr func parse OpenVPN `status 2` from connect write and return status struct、error
func ParseAddr(management string) (*Status, error) {
	conn, err := net.DialTimeout("tcp", management, 2*time.Second)
	if err != nil {
		return &Status{Result: "connect to open server false"}, err
	}

	defer conn.Close()

	// read first sentence
	buf := make([]byte, 75)
	if _, err := io.ReadFull(conn, buf); err != nil {
		return &Status{Result: "conn read false"}, err
	}

	//write `status 2`
	_, err = conn.Write([]byte("status 2\n"))
	if err != nil {
		return &Status{Result: "connect write false"}, err
	}
	// log.Printf("wrote: %d bytes", size)

	var data bytes.Buffer
	// read `status 2` result
	buf = make([]byte, bufferSize)
	for {
		var n int
		n, err = conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Printf("read error: %s", err)
			}
			break
		}
		data.Write(buf[:n])
		if n < bufferSize {
			if pos := bytes.Index(buf[:n], []byte("END")); pos > -1 {
				// log.Printf("found END at %d", pos)
				break
			}
		}
	}

	log.Printf("ovpn parsed %d", data.Len())
	// log.Printf("parsed %s", data)
	return Parse(bytes.NewReader(data.Bytes()))
}