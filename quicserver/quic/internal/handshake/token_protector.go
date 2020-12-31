package handshake

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"fmt"
	"golang.org/x/crypto/hkdf"
	"io"
)

//is used to create and verify a token
type tokenProtector interface {
	NewToken([]byte) ([]byte, error)
	DecodeToken([]byte) ([]byte, error)
}

const (
	tokenSecretSize = 32
	tokenNonceSize  = 32
)

type tokenProtectorImpl struct {
	rand   io.Reader
	secret []byte
}

func (s *tokenProtectorImpl) NewToken(data []byte) ([]byte, error) {
	nonce := make([]byte, tokenNonceSize)
	if _, err := s.rand.Read(nonce); err != nil {
		return nil, err
	}
	aead, aeadNonce, err := s.createAEAD(nonce)
	if err != nil {
		return nil, err
	}
	return append(nonce, aead.Seal(nil, aeadNonce, data, nil)...), nil
}

func (s *tokenProtectorImpl) DecodeToken(p []byte) ([]byte, error) {
	if len(p) < tokenNonceSize {
		return nil, fmt.Errorf("token too short: %d", len(p))
	}
	nonce := p[:tokenNonceSize]
	aead, aeadNonce, err := s.createAEAD(nonce)
	if err != nil {
		return nil, err
	}
	return aead.Open(nil, aeadNonce, p[tokenNonceSize:], nil)
}

func newTokenProtector(rand io.Reader) (tokenProtector, error) {
	secret := make([]byte, tokenSecretSize)
	if _, err := rand.Read(secret); err != nil {
		return nil, err
	}
	return &tokenProtectorImpl{
		rand:   rand,
		secret: secret,
	}, nil
}

func (s *tokenProtectorImpl) createAEAD(nonce []byte) (cipher.AEAD, []byte, error) {
	h := hkdf.New(sha256.New, s.secret, nonce, []byte("quic-go token source"))
	key := make([]byte, 32) // use a 32 byte key, in order to select AES-256
	if _, err := io.ReadFull(h, key); err != nil {
		return nil, nil, err
	}
	aeadNonce := make([]byte, 12)
	if _, err := io.ReadFull(h, aeadNonce); err != nil {
		return nil, nil, err
	}
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}
	aead, err := cipher.NewGCM(c)
	if err != nil {
		return nil, nil, err
	}
	return aead, aeadNonce, nil
}
