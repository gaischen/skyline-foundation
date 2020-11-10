package examples

import (
	"encoding/binary"
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
		wel := &Message{id: atomic.AddUint32(&s.pkgId, 1), value: "say hello...", flag: uint32(1)}
		byt := Msg2Bytes(wel)
		_, err = conn.Write(byt)
		if err != nil {
			return
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	clientAddr := conn.RemoteAddr().String()
	fmt.Println("Connection success. Client address: ", clientAddr)
	defer conn.Close()

	for {
		buffer := make([]byte, 4)
		_, err := conn.Read(buffer)
		length := binary.BigEndian.Uint32(buffer[:])

		buf2 := make([]byte, length-4)
		_, err = conn.Read(buf2)

		buf3 := make([]byte, length)
		buf3 = append(buffer, buf2...)

		if err != nil {
			fmt.Println("Read error: ", err, clientAddr)
			return
		}

		fmt.Println("Client message: ", Bytes2Msg(buf3))
	}
}
