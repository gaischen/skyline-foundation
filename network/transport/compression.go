package transport

type compression interface {
	DeCompress(in []byte) ([]byte, error)
	Compress(body []byte) ([]byte, error)
}

type lz4Compression struct {
}

func (l *lz4Compression) DeCompress(in []byte) ([]byte, error) {
	panic("implement me")
}

func (l *lz4Compression) Compress(body []byte) ([]byte, error) {
	panic("implement me")
}
