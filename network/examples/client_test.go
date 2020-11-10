package examples

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/vanga-top/skyline-foundation/log"
	"github.com/vanga-top/skyline-foundation/log/level"
	"io"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var logger log.Logger = log.NewLogger("net_example", level.ERROR)

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
	c.isAvailable = true
	go c.handlerReceive()
	go c.handleMsgChan()
	go c.heartbeat()
	c.wg.Wait()
}

func (c *Client) heartbeat() {
	t := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-t.C:
			str := "go tick " + strconv.Itoa(int(c.pkgId))
			c.send(build(c, str))
		case <-c.ctx.Done():
			return
		}
	}
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
	isAvailable bool
}

func (c *Client) send(byt []byte) {
	if c.isAvailable {
		c.conn.Write(byt)
	}
}

func (c *Client) connect() (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:8080", time.Second*3)
	if err != nil {
		logger.Debug("error in connection", err)
		return nil, err
	}
	logger.Debug("connecting server with client name:", c.name)
	c.wg.Add(1)
	return conn, err
}

func (c *Client) handlerReceive() {
LOOP:
	for {
		select {
		case <-c.ctx.Done():
			c.closeConnection()
			logger.Debug("close chanel...")
			break LOOP
		default:
			b := make([]byte, 1024)
			l, err := c.conn.Read(b)
			if err != nil {
				logger.Debug("read chan error", err)
			}
			if errors.Cause(err) == io.EOF {
				c.cancel()
				logger.Debug("read IO EOF of chan..")
				return
			}
			msg := Bytes2Msg(b[:l])
			c.msgChan <- msg
		}
	}
}

func (c *Client) closeConnection() {
	defer func() {
		if r := recover(); r != nil {
			logger.Debug("close error...", r)
		}
	}()
	c.closeOnce.Do(func() {
		c.isAvailable = false
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
				c.conn.Write(build(c, "go back"))
			}
		case <-c.ctx.Done():
			logger.Debug("close chanel")
			return
		default:
			continue
		}
	}
}

func build(c *Client, value string) []byte {
	rtn := &Message{flag: uint32(1), id: atomic.AddUint32(&c.pkgId, 1), value: value}
	byt := Msg2Bytes(rtn)
	return byt
}
