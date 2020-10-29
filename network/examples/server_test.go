package examples

import (
	"fmt"
	"net"
	"testing"
)

func TestStartServer(t *testing.T) {
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
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
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
		strBuffer := string(buffer[:recvLen])
		fmt.Println("Client message: ", strBuffer)
	}
}
