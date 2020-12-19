package quic

import (
	"crypto/tls"
	"net"
	"sync"
)

//start listen quic addr
func ListenAddr(addr string, tlsConfig *tls.Config, config *Config) (Listener, error) {

	return nil, nil
}

type basicServer struct {
	mutex               sync.Mutex
	acceptEarlySessions bool
	tlsConfig           *tls.Config
	config              *Config

	conn net.PacketConn
	//if the server is started with listenAddr we create a packet conn.
	//if it is started with listen we take a packet conn as a parameter
	createdPacketConn bool
}
