package zlog

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"log/slog"
	"time"
)

type zLog struct {
	zerolog.Logger
	Config LoggerConfig
}

type LoggerConfig struct {
	CallerSkipFrameCount int
}

var (
	DefaultDiode  diode.Writer
	DefaultLogger zLog
	defaultLevel  = slog.LevelDebug
)

var (
	DiodeBufferSize   = 1000
	DiodePollInterval = time.Millisecond
	DiodeDroppedLogFn = DefaultDiodeDroppedLogFn
)

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
