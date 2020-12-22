package protocol

import "time"

// DefaultHandshakeTimeout is the default timeout for a connection until the crypto handshake succeeds.
const DefaultHandshakeTimeout = 10 * time.Second

// DefaultIdleTimeout is the default idle timeout
const DefaultIdleTimeout = 30 * time.Second

// DefaultMaxReceiveStreamFlowControlWindow is the default maximum stream-level flow control window for receiving data, for the server
const DefaultMaxReceiveStreamFlowControlWindow = 6 * (1 << 20) // 6 MB

// DefaultMaxReceiveConnectionFlowControlWindow is the default connection-level flow control window for receiving data, for the server
const DefaultMaxReceiveConnectionFlowControlWindow = 15 * (1 << 20) // 12 MB

// DefaultMaxIncomingStreams is the maximum number of streams that a peer may open
const DefaultMaxIncomingStreams = 100

// DefaultMaxIncomingUniStreams is the maximum number of unidirectional streams that a peer may open
const DefaultMaxIncomingUniStreams = 100


// DefaultConnectionIDLength is the connection ID length that is used for multiplexed connections
// if no other value is configured.
const DefaultConnectionIDLength = 4

// TokenValidity is the duration that a (non-retry) token is considered valid
const TokenValidity = 24 * time.Hour

// RetryTokenValidity is the duration that a retry token is considered valid
const RetryTokenValidity = 10 * time.Second