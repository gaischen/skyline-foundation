package http2

import (
	"fmt"
	"net"
	"testing"
)

func TestRunServer(t *testing.T) {
	ln, _ := net.Listen("tcp", ":7001")
	for {
		rawConn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		newHttp2Server(rawConn, &ServerConfig{})
	}
}
