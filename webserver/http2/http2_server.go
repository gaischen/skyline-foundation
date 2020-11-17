package http2

import (
	"bytes"
	"context"
	"errors"
	"github.com/vanga-top/skyline-foundation/webserver/credentials"
	"github.com/vanga-top/skyline-foundation/webserver/keepalive"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"
	"math"
	"net"
	"sync"
	"time"
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

	mu              sync.Mutex
	state           transportState
	activeStreams   map[uint32]*Stream
	streamSendQuota uint32

	activity uint32
	kp       keepalive.ServerParameters
	idle     time.Time
}

func newHttp2Server(conn net.Conn, config *ServerConfig) (_ ServerTransport, err error) {
	framer := newFramer(conn)
	var settings []http2.Setting
	maxStreams := config.MaxStream
	if maxStreams == 0 {
		maxStreams = math.MaxUint32
	} else {
		settings = append(settings, http2.Setting{
			ID:  http2.SettingMaxConcurrentStreams,
			Val: maxStreams,
		})
	}

	if initialWindowSize != defaultWindowSize {
		settings = append(settings, http2.Setting{
			ID:  http2.SettingInitialWindowSize,
			Val: uint32(initialWindowSize),
		})
	}

	//触发client.handleSettings事件
	if err := framer.writeSettings(true, settings...); err != nil {
		return nil, errors.New("error in netH2Server writeSettings")
	}
	//触发client.windowUpdate事件
	if delta := uint32(initialWindowSize - defaultWindowSize); delta > 0 {
		if err := framer.writeWindowUpdate(true, 0, delta); err != nil {
			return nil, errors.New("error in netH2Server writeWindowUpdate")
		}
	}

	kp := keepalive.ServerParameters{}
	kp.Time = defaultServerKeepaliveTime
	kp.Timeout = defaultServerKeepaliveTimeout

	var buf bytes.Buffer

	server := &http2Server{
		ctx:             context.Background(),
		conn:            conn,
		remoteAddr:      conn.RemoteAddr(),
		localAddr:       conn.LocalAddr(),
		authInfo:        config.AuthInfo,
		framer:          framer,
		hBuf:            &buf,
		hEnc:            hpack.NewEncoder(&buf),
		maxStreams:      maxStreams,
		controlBuf:      newControlBuffer(),
		fc:              &transportInFlow{limit: initialConnWindowSize},
		sendQuotaPool:   newQuotaPool(defaultWindowSize),
		state:           reachable,
		writableChan:    make(chan int, 1),
		shutdownChan:    make(chan struct{}),
		activeStreams:   make(map[uint32]*Stream),
		streamSendQuota: defaultWindowSize,
		kp:              kp,
	}

	go server.controller()
	go server.keepalive()
	server.writableChan <- 0
	return server, nil
}

func (s *http2Server) controller() {
	for {
		select {
		case i := <-s.controlBuf.get():
			s.controlBuf.load()
			select {
			case <-s.writableChan:
				switch i := i.(type) {
				case *windowUpdate:
					s.framer.writeWindowUpdate(true, i.streamId, i.increment)
				case *settings:
					if i.ack {
						s.framer.writeSettingsAck(true)
						s.applySetting(i.ss)
					} else {
						s.framer.writeSettings(true, i.ss...)
					}
				case *resetStream:
					s.framer.writeRSTStream(true, i.streamId, i.code)
				case *goAway:
					s.mu.Lock()
					if s.state == closing {
						s.mu.Unlock()
						return
					}
					sid := s.maxStreamID
					s.state = draining
					s.mu.Unlock()
					s.framer.writeGoAway(true, sid, http2.ErrCodeNo, nil)
				case *flushIO:
					s.framer.flushWrite()
				case *ping:
					s.framer.writePing(true, i.ack, i.data)
				default:
				}
				s.writableChan <- 0
				continue
			case <-s.shutdownChan:
				return
			}
		case <-s.shutdownChan:
			return
		}
	}
}

func (s *http2Server) keepalive() {

}

func (s *http2Server) applySetting(ss []http2.Setting) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, set := range ss {
		if set.ID == http2.SettingInitialWindowSize {
			for _, stream := range s.activeStreams {
				stream.sendQuotaPool.add(int(set.Val) - int(s.streamSendQuota))
			}
			s.streamSendQuota = set.Val
		}
		if set.ID == http2.SettingMaxConcurrentStreams {
			s.maxStreams = set.Val
		}
	}
}
