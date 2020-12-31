package handshake

import "io"

type TokenGenerator struct {
	tokenProtector tokenProtector
}

func NewTokenGenerator(rand io.Reader) (*TokenGenerator, error) {
	
}
