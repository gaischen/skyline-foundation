package handshake

//is used to create and verify a token
type tokenProtector interface {
	NewToken([]byte) ([]byte, error)
	DecodeToken([]byte) ([]byte, error)
}



