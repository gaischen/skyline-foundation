package quic

import (
	"errors"
	"github.com/vanga-top/skyline-foundation/quicserver/quic/internal/protocol"
)

func validateConfig(config *Config) error {
	if config == nil {
		return nil
	}
	if config.MaxIncomingStreams > 1<<60 {
		return errors.New("error in value for MaxIncomingStreams")
	}
	if config.MaxIncomingUniStreams > 1<<60 {
		return errors.New("error in value for MaxIncomingUniStreams")
	}
	return nil
}

func populateServerConfig(config *Config) *Config {
	config = populateConfig(config)
}

func populateConfig(config *Config) *Config {
	if config == nil {
		config = &Config{}
	}
	versions := config.Versions
	if len(versions) == 0 {
		versions = protocol.SupportedVersions
	}
	handshakeTimeout := protocol.DefaultHandshakeTimeout
	if config.HandshakeTimeout != 0 {
		handshakeTimeout = config.HandshakeTimeout
	}
	idleTimeout := protocol.DefaultIdleTimeout
	if config.MaxIdleTimeout != 0 {
		idleTimeout = config.MaxIdleTimeout
	}

}
