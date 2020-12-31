package main

import (
	"context"
	"fmt"
	"github.com/vanga-top/skyline-foundation/quicserver/quic"
	"io"
)

func main() {
	listener, err := quic.ListenAddr(saddr, generateTLSConfig(), nil)
	if err != nil {
		fmt.Println(err)
	}
	ctx, _ := context.WithCancel(context.Background())
	for {
		sess, err := listener.Accept(ctx)
		if err != nil {
			fmt.Println(err)
		} else {
			go dealSession01(sess)
		}
	}
}


func dealSession01(sess quic.Session) {
	ctx, _ := context.WithCancel(context.Background())
	stream, err := sess.AcceptStream(ctx)
	if err != nil {
		panic(err)
	} else {
		for {
			_, err = io.Copy(loggingWriter{stream}, stream)
		}
	}
}