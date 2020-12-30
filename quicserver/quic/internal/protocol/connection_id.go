package protocol

type ConnectionID []byte


// Bytes returns the byte representation
func (c ConnectionID) Bytes() []byte {
	return []byte(c)
}