package http2

import (
	"context"
	"github.com/vanga-top/skyline-foundation/h2server/codes"
	"github.com/vanga-top/skyline-foundation/h2server/metadata"
	"golang.org/x/net/http2"
	"io"
	url2 "net/url"
	"sync"
	"time"
)

type Stream struct {
	StreamExtra
	id               uint32
	ct               clientTransport
	ctx              context.Context
	cancel           context.CancelFunc
	done             chan struct{}
	goAway           chan struct{}
	t                time.Time
	recvCompress     string //接收端压缩算法
	sendCompress     string //
	buf              *recvBuffer
	recvBufferReader io.Reader
	flowControl      *inFlow
	recvQuota        uint32 // 累计的请求配额
	windowHandler    func(int)
	sendQuotaPool    *quotaPool
	headerChan       chan struct{}
	header           metadata.MD
	trailer          metadata.MD
	mu               sync.RWMutex
	headerOK         bool
	state            streamState
	headerDone       bool //headerchan close
	statusCode       codes.Code
	statusDesc       string
	rstStream        bool //是否发送RST_STREAM FRAME到server
	rstError         http2.ErrCode
}

type StreamExtra struct {
	method     string
	vhost      string
	schema     string
	httpMethod string
	httpPath   string
	uri        string
	userAgent  string
	did        string
	refer      string
	url        *url2.URL
}

type streamState uint8

const (
	streamActive    streamState = iota
	streamWriteDone             // 已经发送了EndStream
	streamReadDone              // 已经接收了EndStream
	streamDone                  // 整个stream结束
)
