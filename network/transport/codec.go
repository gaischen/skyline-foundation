package transport

type codec interface {
	Read(data []byte, obj interface{}) error
	Write(msg interface{}) ([]byte, error)
}
