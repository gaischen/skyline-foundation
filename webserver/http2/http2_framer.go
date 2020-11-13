package http2

import (
	"bufio"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"
	"io"
	"net"
)

const (
	http2IOWriteBufSize      = 8 * 1024
	http2InitHeaderTableSize = 4096
)

type framer struct {
	numWriters int32
	reader     io.Reader
	writer     *bufio.Writer
	fr         *http2.Framer
}

func newFramer(conn net.Conn) *framer {
	f := &framer{
		reader: conn,
		writer: bufio.NewWriterSize(conn, http2IOWriteBufSize),
	}
	f.fr = http2.NewFramer(f.writer, f.reader)
	f.fr.SetReuseFrames()
	f.fr.ReadMetaHeaders = hpack.NewDecoder(http2InitHeaderTableSize, nil)
	return f
}

func (f *framer) writeSettings(forceFlush bool, settings ...http2.Setting) error {
	if err := f.fr.WriteSettings(settings...); err != nil {
		return err
	}
	if forceFlush {
		return f.writer.Flush()
	}
	return nil
}

func (f *framer) writeWindowUpdate(forceFlush bool, streamID, incr uint32) error {
	if err := f.fr.WriteWindowUpdate(streamID, incr); err != nil {
		return err
	}
	if forceFlush {
		return f.writer.Flush()
	}
	return nil
}

func (f *framer) writeSettingsAck(forceFlush bool) error {
	if err := f.fr.WriteSettingsAck(); err != nil {
		return err
	}
	if forceFlush {
		f.writer.Flush()
	}
	return nil
}
