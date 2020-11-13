package http2

import (
	"context"
	"github.com/vanga-top/skyline-foundation/webserver/credentials"
	"net"
)

type http2Server struct {
	id           string
	ctx          context.Context
	conn         net.Conn
	remoteAddr   net.Addr
	localAddr    net.Addr
	maxStreamID  uint32
	authInfo     credentials.AuthInfo
	writableChan chan int //sync write
	shutdownChan chan struct{}

}
