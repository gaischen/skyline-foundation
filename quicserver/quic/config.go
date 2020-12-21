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
	maxReceiveStreamFlowControlWindow := config.MaxReceiveStreamFlowControlWindow
	if maxReceiveStreamFlowControlWindow == 0 {
		maxReceiveStreamFlowControlWindow = protocol.DefaultMaxReceiveStreamFlowControlWindow
	}
	maxReceiveConnectionFlowControlWindow := config.MaxReceiveConnectionFlowControlWindow
	if maxReceiveConnectionFlowControlWindow == 0 {
		maxReceiveConnectionFlowControlWindow = protocol.DefaultMaxReceiveConnectionFlowControlWindow
	}
	maxIncomingStreams := config.MaxIncomingStreams
	if maxIncomingStreams == 0 {
		maxIncomingStreams = protocol.DefaultMaxIncomingStreams
	} else if maxIncomingStreams < 0 {
		maxIncomingStreams = 0
	}
	maxIncomingUniStreams := config.MaxIncomingUniStreams
	if maxIncomingUniStreams == 0 {
		maxIncomingUniStreams = protocol.DefaultMaxIncomingUniStreams
	} else if maxIncomingUniStreams < 0 {
		maxIncomingUniStreams = 0
	}
	return &Config{
		Versions:                              versions,
		HandshakeTimeout:                      handshakeTimeout,
		MaxIdleTimeout:                        idleTimeout,
		AcceptToken:                           config.AcceptToken,
		KeepAlive:                             config.KeepAlive,
		MaxReceiveStreamFlowControlWindow:     maxReceiveStreamFlowControlWindow,
		MaxReceiveConnectionFlowControlWindow: maxReceiveConnectionFlowControlWindow,
		MaxIncomingStreams:                    maxIncomingStreams,
		MaxIncomingUniStreams:                 maxIncomingUniStreams,
		ConnectionIDLength:                    config.ConnectionIDLength,
		StatelessResetKey:                     config.StatelessResetKey,
		TokenStore:                            config.TokenStore,
		Tracer:                                config.Tracer,
	}
}
