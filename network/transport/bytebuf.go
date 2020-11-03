package transport

import "io"

type ByteBuf struct {
	data []byte
}

func (b *ByteBuf) Bytes() []byte {
	return b.data
}

func (b *ByteBuf) Len() int {
	return len(b.data)
}

func (b *ByteBuf) Resize(n int) {
	if n > cap(b.data) {
		b.reserve(n)
	}
	b.data = b.data[0:n]
}

func (b *ByteBuf) reserve(n int) {
	if cap(b.data) >= n {
		return
	}
	m := cap(b.data)
	if m == 0 {
		m = 1024
	}
	for m < n {
		m *= 2
	}
	data := make([]byte, len(b.data), m)
	copy(data, b.data)
	b.data = data
}

//read full
func (b *ByteBuf) ReadFull(r io.Reader, n int) error {
	l := len(b.data)
	b.reserve(l + n)
	for {
		m, err := r.Read(b.data[l : l+n])
		b.data = b.data[0 : l+m]
		if b.Len() >= l+n {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *ByteBuf) Reset() {
	b.data = b.data[:0]
}
