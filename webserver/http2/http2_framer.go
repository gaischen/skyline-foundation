package http2

import (
	"bufio"
	"golang.org/x/net/http2"
	"io"
)

type framer struct {
	numWriters int32
	reader     io.Reader
	writer     *bufio.Writer
	fr         *http2.Framer
}
