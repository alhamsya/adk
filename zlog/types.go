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

// Annotation holds the key-value pairs for logging metadata.
type Annotation struct {
	*sync.Map
	isDefault bool
}

// loggingAnnotation is a context key for storing the Annotation.
type loggingAnnotation struct{}

// Option conforms to the option pattern for configuring context behavior.
// It is an interface to allow the implementation to be hidden.
type Option interface {
	apply(*Annotation)
}

type optionFunc func(*Annotation)
