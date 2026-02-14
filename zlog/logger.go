package zlog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"github.com/rs/zerolog/pkgerrors"
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

const DateTimeZone = "2025-12-15T07:15:00+0700"

// --- Initialization ---

func init() {
	setupZeroLogGlobals()
	setupDefaultDiode()
	setupSlogDefault()
	setupDefaultLogger()
}

func setupZeroLogGlobals() {
	setDefaultLevel()
	zerolog.TimestampFieldName = "timestamp"
	zerolog.TimeFieldFormat = DateTimeZone
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.ErrorHandler = func(err error) {
		slog.Error(err.Error())
	}
}

func setupDefaultDiode() {
	DefaultDiode = diode.NewWriter(os.Stdout, DiodeBufferSize, DiodePollInterval, DiodeDroppedLogFn)
}

func setupSlogDefault() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(DefaultDiode, &slog.HandlerOptions{
		AddSource: true,
		Level:     defaultLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.TimeKey:
				return slog.String(zerolog.TimestampFieldName, a.Value.Time().Format(time.RFC3339))
			case slog.LevelKey:
				return slog.String(slog.LevelKey, strings.ToLower(a.Value.String()))
			}
			return a
		},
	})))
}

func setupDefaultLogger() {
	DefaultLogger = New(DefaultDiode, WithLoggerCallerSkipFrameCount(zerolog.CallerSkipFrameCount+2))
	zerolog.DefaultContextLogger = &DefaultLogger.Logger
}

// --- Helper Functions ---

func DefaultDiodeDroppedLogFn(dropCount int) {
	slog.Warn(
		fmt.Sprintf(
			"zLog: dropped %d logs due to buffer overflow",
			dropCount,
		),
	)
}

func setDefaultLevel() {
	envStr := os.Getenv(LogLevelEnvKey)
	defaultLevel = strToLevel(envStr)
	zerolog.SetGlobalLevel(slogLevelToZerologLevel(defaultLevel))
}

func strToLevel(str string) slog.Level {
	switch strings.ToLower(str) {
	case LevelDebug:
		return slog.LevelDebug
	case LevelInfo:
		return slog.LevelInfo
	case LevelWarn:
		return slog.LevelWarn
	case LevelError:
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}

func slogLevelToZerologLevel(sl slog.Level) zerolog.Level {
	switch sl {
	case slog.LevelDebug:
		return zerolog.DebugLevel
	case slog.LevelInfo:
		return zerolog.InfoLevel
	case slog.LevelWarn:
		return zerolog.WarnLevel
	case slog.LevelError:
		return zerolog.ErrorLevel
	default:
		return zerolog.DebugLevel
	}
}

// --- Constructors ---

// New creates a new Logger instance with the given output writer and options.
func New(output io.Writer, opts ...LoggerOption) Logger {
	config := LoggerConfig{
		CallerSkipFrameCount: zerolog.CallerSkipFrameCount + 1,
	}

	for _, o := range opts {
		o(&config)
	}

	logger := zerolog.
		New(output).
		With().
		Timestamp().
		CallerWithSkipFrameCount(config.CallerSkipFrameCount).
		Logger()

	return Logger{
		Logger: logger,
		Config: config,
	}
}

func WithLoggerCallerSkipFrameCount(skipCount int) LoggerOption {
	return func(cfg *LoggerConfig) {
		cfg.CallerSkipFrameCount = skipCount
	}
}

// FromContext returns the Logger associated with the ctx. If no logger
// is associated, it returns the global default logger.
func FromContext(ctx context.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx)
}

// NewContext returns a new context with the provided zerolog hooks attached to the logger.
func NewContext(ctx context.Context, hooks ...zerolog.Hook) (context.Context, *zerolog.Logger) {
	logger := zerolog.Ctx(ctx)
	for _, h := range hooks {
		l := logger.Hook(h).With().Logger()
		logger = &l
	}
	return logger.WithContext(ctx), logger
}
