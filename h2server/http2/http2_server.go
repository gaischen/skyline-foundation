package http2

import (
	"bytes"
	"context"
	"errors"
	"github.com/vanga-top/skyline-foundation/webserver/codes"
	"github.com/vanga-top/skyline-foundation/webserver/credentials"
	"github.com/vanga-top/skyline-foundation/webserver/keepalive"
	"github.com/vanga-top/skyline-foundation/webserver/metadata"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"
	"io"
	"math"
	"net"
	"sync"
	"sync/atomic"
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

func (s *http2Server) SetId(id string) {
	panic("implement me")
}

func (s *http2Server) HandleStream(handle func(stream *Stream)) {
	frame, err := s.framer.readFrame()
	if err == io.EOF || err == io.ErrUnexpectedEOF {
		s.Close()
		return
	}
	if err != nil {
		s.Close()
		return
	}

	//标记为活跃流
	atomic.StoreUint32(&s.activity, 1)
	sf, ok := frame.(*http2.SettingsFrame)
	//第一个必须是setting
	if !ok {
		s.Close()
		return
	}
	s.handleSettings(sf)

	for {
		frame, err := s.framer.readFrame()
		atomic.StoreUint32(&s.activity, 1)
		if err != nil {
			if se, ok := err.(http2.StreamError); ok {
				s.mu.Lock()
				stream := s.activeStreams[se.StreamID]
				s.mu.Unlock()
				if stream != nil {
					s.closeStream(stream)
				}
				s.controlBuf.put(&resetStream{se.StreamID, se.Code})
				continue
			}
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				s.Close()
				return
			}
			s.Close()
			return
		}
		switch frame := frame.(type) {
		case *http2.MetaHeadersFrame:

		case *http2.DataFrame:
			s.handleData(frame)
		}
	}

}

func (s *http2Server) handleSettings(frame *http2.SettingsFrame) {
	if frame.IsAck() {
		return
	}

	var ss []http2.Setting
	frame.ForeachSetting(func(setting http2.Setting) error {
		ss = append(ss, setting)
		return nil
	})
	s.controlBuf.put(&settings{ack: true, ss: ss})
}

func (s *http2Server) WriteHeader(stream *Stream, md metadata.MD) error {
	panic("implement me")
}

func (s *http2Server) Write(stream *Stream, data []byte, opts *Options) error {
	panic("implement me")
}

func (s *http2Server) WriteStatus(stream *Stream, statusCode codes.Code, statusDesc string) error {
	panic("implement me")
}

func (s *http2Server) Push(stream *Stream, data []byte, flags http2.Flags) error {
	panic("implement me")
}

func (s *http2Server) RemoteAddr() net.Addr {
	panic("implement me")
}

func (s *http2Server) Drain() {
	panic("implement me")
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
	p := &ping{}
	var pingSent bool
	keepalive := time.NewTimer(s.kp.Time)
	defer func() {
		if !keepalive.Stop() {
			<-keepalive.C
		}
	}()
	for {
		select {
		case <-keepalive.C:
			if atomic.CompareAndSwapUint32(&s.activity, 1, 0) {
				pingSent = false
				keepalive.Reset(s.kp.Time)
				continue
			}
			if pingSent {
				s.Close()
				keepalive.Reset(time.Duration(math.MaxInt64))
				return
			}
			pingSent = true
			s.controlBuf.put(p)
			//等待ping的buffer
			keepalive.Reset(s.kp.Timeout)
		case <-s.shutdownChan:
			return
		}
	}
}

func (s *http2Server) Close() (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.state == closing {
		return errors.New("transport close() has bean called... ")
	}
	s.state = closing
	streams := s.activeStreams
	s.activeStreams = nil
	close(s.shutdownChan)
	err = s.conn.Close()
	for _, stream := range streams {
		if stream.cancel != nil {
			stream.cancel()
		}
	}
	return
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

func (s *http2Server) closeStream(stream *Stream) {
	s.mu.Lock()
	delete(s.activeStreams, stream.id)

	if s.state == draining && len(s.activeStreams) == 0 {
		defer s.Close()
	}
	s.mu.Unlock()
	stream.cancel()
	stream.mu.Lock()
	if stream.state == streamDone {
		stream.mu.Unlock()
		return
	}
	stream.state = streamDone
	stream.mu.Unlock()
}

func (s *http2Server) handleData(frame *http2.DataFrame) {
	size := frame.Header().Length

}
