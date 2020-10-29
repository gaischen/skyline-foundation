package transport

import (
	"github.com/skyline/skyline-foundation/log"
	"github.com/skyline/skyline-foundation/log/level"
	"io"
	"net"
	"time"
)

const (
	tcpConnectionTimeout = 3 * time.Second
)

var logger log.Logger = log.NewLogger("connection", level.WARN)

type connection struct {
	id uint32 // connection id, unique one process

	readBytes  uint32
	readPkgs   uint32
	writeBytes uint32
	writePkgs  uint32

	readLastTime  int64
	writeLastTime int64

	isClosed   bool
	conn       net.Conn
	reader     io.Reader
	writer     io.Writer
	localAddr  string
	remoteAddr string

	buf []byte
}

/**
new connection
*/
func newConnection(addr string) (*connection, error) {
	conn, err := net.DialTimeout("tcp", addr, tcpConnectionTimeout)
	if err != nil {
		logger.Error(err)
		return nil, err
	}


	return nil, nil
}
