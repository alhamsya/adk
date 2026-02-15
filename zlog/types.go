package zlog

import (
	"log/slog"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
)

// --- Types & Constants ---

type Logger struct {
	zerolog.Logger
	Config LoggerConfig
}

type LoggerConfig struct {
	CallerSkipFrameCount int
}

type LoggerOption func(cfg *LoggerConfig)

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
	DefaultDiode  diode.Writer
	DefaultLogger Logger
	defaultLevel  = slog.LevelDebug
)

var (
	DiodeBufferSize   = 1000
	DiodePollInterval = time.Millisecond
	DiodeDroppedLogFn = DefaultDiodeDroppedLogFn
)

const DateTimeZone = time.RFC3339

type Metadata struct {
	*sync.Map
}

type loggingMetadata struct{}
