package http2

import (
	"bytes"
	"context"
	"github.com/vanga-top/skyline-foundation/webserver/credentials"
	"golang.org/x/net/http2/hpack"
	"net"
	"sync"
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

	framer *framer
	hBuf   *bytes.Buffer
	hEnc   *hpack.Encoder

	maxStreams    uint32           //最大并发
	controlBuf    *controlBuffer   //
	fc            *transportInFlow //针对输入流的控制
	sendQuotaPool *quotaPool       //输出流控制

	mu            sync.Mutex
	state         transportState
	activeStreams map[uint32]*Stream
	streamSenQuota uint32

	activity uint32

}
