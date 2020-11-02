package transport

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
		b.Resize(n)
	}
	b.data = b.data[0:n]
}
