package examples

import (
	"fmt"
	"net"
	"sync/atomic"
	"testing"
)

func TestServerRunning(t *testing.T) {
	s := &Server{name: "test_s"}
	s.start()
}

type Server struct {
	name  string
	pkgId uint32
}

func (s *Server) start() {
	lner, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Listener create error: ", err)
		return
	}
	fmt.Println("Waiting for client...")
	for {
		conn, err := lner.Accept()
		if err != nil {
			fmt.Println("Accept error: ", err)
			return
		}
		wel := &Message{id: atomic.AddUint32(&s.pkgId, 1), value: "welcome hello..."}
		byt := MessageToBytes(wel)
		wel.length = len(byt)
		fmt.Println("send length:", wel.length)
		conn.Write(byt)
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	clientAddr := conn.RemoteAddr().String()
	fmt.Println("Connection success. Client address: ", clientAddr)
	defer conn.Close()

	for {
		buffer := make([]byte, 1024)
		recvLen, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Read error: ", err, clientAddr)
			return
		}
		msg := BytesToMessage(buffer[:recvLen])
		fmt.Println("Client message: ", msg)
	}
}
