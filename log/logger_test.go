package log

import (
	"github.com/skyline/skyline-foundation/log/level"
	"testing"
)

func TestLoggerInit(t *testing.T) {
	logger := NewLogger("test", level.ERROR)
	logger.Debug("test", "   ha")
}
