package network

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/vanga-top/skyline-foundation/log"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

const (
	tcpConnectionTimeout = 3 * time.Second
)

var connId uint32
var recvBuf = new(sync.Pool)
var writeBuf = new(sync.Pool)

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

	buf    []byte
	logger log.Logger
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

	cnt := &connection{
		id:         atomic.AddUint32(&connId, 1),
		conn:       conn,
		reader:     io.Reader(conn),
		writer:     io.Writer(conn),
		localAddr:  conn.LocalAddr().String(),
		remoteAddr: conn.RemoteAddr().String(),
	}
	return cnt, nil
}

func (conn *connection) incReadPkgNum() {
	atomic.AddUint32(&conn.readPkgs, 1)
}

func (conn *connection) incWritePkgNum(delta uint32) {
	atomic.AddUint32(&conn.writePkgs, delta)
}

func acquireRecvBuf() *ByteBuf {
	v := recvBuf.Get()
	if v == nil {
		v = new(ByteBuf)
	}
	return v.(*ByteBuf)
}

func releaseWriteBuf(buf *bytes.Buffer) {
	buf.Reset()
	writeBuf.Put(buf)
}

func releaseRecvBuf(buf *ByteBuf) {
	buf.Reset()
	recvBuf.Put(buf)
}

func acquireWriteBuf() *bytes.Buffer {
	v := writeBuf.Get()
	if v == nil {
		v = new(bytes.Buffer)
	}
	return v.(*bytes.Buffer)
}

//return pkg length ; content length ; error
func (conn *connection) getPkgLength() (int, int, error) {
	b := make([]byte, 4) //todo 这里写死了
	if l, err := io.ReadFull(conn.reader, b); err != nil || l != 4 {
		if conn.isClosed {
			logger.Debug("[error] tcp connection is closed..")
		}
		logger.Debug("[error] recv tcp pkg failed:%v", err)
		return 0, 0, err
	} else {
		length := binary.BigEndian.Uint32(b[:])
		return int(length), int(length) - 4, err
	}
}

//receive
func (conn *connection) recv() (*ByteBuf, int, error) {
	var (
		length     int //require pkg length
		dataLength int //require content length
		err        error
	)
	if length, dataLength, err = conn.getPkgLength(); err != nil {
		return nil, dataLength, err
	}
	buf := acquireRecvBuf()
	err = buf.ReadFull(conn.reader, dataLength)
	if err != nil {
		if conn.isClosed {
			err = nil
		}
	}
	logger.Debug("[debug] receive data peer addr: %s, local addr:%s", conn.remoteAddr, conn.localAddr)

	atomic.AddUint32(&conn.readBytes, uint32(dataLength))
	conn.incReadPkgNum()
	return buf, length, nil
}

func (conn *connection) send(data interface{}) (int, error) {

	if buffers, ok := data.([]*bytes.Buffer); ok {
		defer func() {
			for _, buffer := range buffers {
				releaseWriteBuf(buffer)
			}
		}()
		bufBytes := make([][]byte, len(buffers))
		for _, buffer := range buffers {
			bufBytes = append(bufBytes, buffer.Bytes())
		}
		netBuf := net.Buffers(bufBytes)
		if length, err := netBuf.WriteTo(conn.conn); err == nil { //batch writer
			atomic.StoreInt64(&conn.writeLastTime, time.Now().Unix())
			atomic.AddUint32(&conn.writeBytes, (uint32)(length))
			conn.incWritePkgNum((uint32)(len(buffers)))
			length += length
			return int(length), err
		} else {
			logger.Error("[error] send data failed..", err)
			return int(length), err
		}
	}

	if data, ok := data.(*bytes.Buffer); ok {
		defer releaseWriteBuf(data)
		num, err := conn.writer.Write(data.Bytes())
		if err == nil {
			atomic.StoreInt64(&conn.writeLastTime, time.Now().Unix())
			atomic.AddUint32(&conn.writeBytes, (uint32)(num))
			conn.incWritePkgNum(1)
			return num, err
		} else {
			logger.Debug("[debug] send data failed...single model")
			return 0, err
		}
	}

	return 0, errors.New("illegal data type !!!")
}

func (conn *connection) close() {
	if conn == nil {
		return
	}

	err := conn.conn.Close()
	if err != nil {
		logger.Error("[error] error closing connection...")
	}
	conn.isClosed = true //todo
}
