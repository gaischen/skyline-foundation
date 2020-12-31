package handshake

import "io"

type TokenGenerator struct {
	tokenProtector tokenProtector
}

func NewTokenGenerator(rand io.Reader) (*TokenGenerator, error) {
	tokenProtector, err := newTokenProtector(rand)
	if err != nil {
		return nil, err
	}
	return &TokenGenerator{
		tokenProtector: tokenProtector,
	}, nil
}
