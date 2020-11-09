package examples

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestClientConnect(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	c := &Client{
		name:    "c1",
		ctx:     ctx,
		cancel:  cancel,
		wg:      sync.WaitGroup{},
		msgChan: make(chan *Message, 10),
	}
	conn, err := c.connect()
	if err != nil {
		return
	}
	c.reader = io.Reader(conn)
	c.conn = conn
	go c.handlerReceive()
	go c.handleMsgChan()
	c.wg.Wait()
}

type Client struct {
	name        string
	ctx         context.Context
	cancel      context.CancelFunc
	messageChan chan []byte
	wg          sync.WaitGroup
	reader      io.Reader
	conn        net.Conn
	closeOnce   sync.Once
	msgChan     chan *Message
	pkgId       uint32
}

func (c *Client) connect() (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:8080", time.Second*3)
	if err != nil {
		fmt.Println("error in connection", err)
		return nil, err
	}
	fmt.Println("connecting server with client name:", c.name)
	c.wg.Add(1)
	return conn, err
}

func (c *Client) handlerReceive() {
LOOP:
	for {
		select {
		case <-c.ctx.Done():
			c.closeConnection()
			fmt.Println("close chanel...")
			break LOOP
		default:
			b := make([]byte, 1024)
			l, err := c.conn.Read(b)
			if err != nil {
				fmt.Println("read chan error", err)
			}
			if errors.Cause(err) == io.EOF {
				c.cancel()
			}
			msg := BytesToMessage(b[:l])
			c.msgChan <- msg
		}
	}
}

func (c *Client) closeConnection() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("close error...", r)
		}
	}()
	c.closeOnce.Do(func() {
		c.wg.Done()
		c.conn.Close()
	})
}

func (c *Client) handleMsgChan() {
	for {
		select {
		case msg := <-c.msgChan:
			fmt.Println(msg)
			if msg.flag == 0 {
				rtn := &Message{flag: 1, id: atomic.AddUint32(&c.pkgId, 1), value: "go back"}
				byt := MessageToBytes(rtn)
				rtn.length = len(byt)
				c.conn.Write(byt)
			}
		default:
			continue
		}
	}
}
