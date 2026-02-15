package zlog

import (
	"log/slog"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
)

// --- Types & Constants ---

type logger struct {
	zerolog.Logger
	Config loggerConfig
}

type loggerConfig struct {
	CallerSkipFrameCount int
}

type loggerOption func(cfg *loggerConfig)

const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelError = "error"
	LevelWarn  = "warn"

	// LogLevelEnvKey is the environment variable used to detect default level
	// for the global logger.
	LogLevelEnvKey = "ZLOG_LEVEL"
)

var (
	defaultDiode  diode.Writer
	defaultLogger logger
	defaultLevel  = slog.LevelDebug
)

var (
	diodeBufferSize   = 1000
	diodePollInterval = time.Millisecond
	diodeDroppedLogFn = defaultDiodeDroppedLogFn
)

type Metadata struct {
	*sync.Map
}

type loggingMetadata struct{}
