package protocol

import "time"

// DefaultHandshakeTimeout is the default timeout for a connection until the crypto handshake succeeds.
const DefaultHandshakeTimeout = 10 * time.Second

// DefaultIdleTimeout is the default idle timeout
const DefaultIdleTimeout = 30 * time.Second
