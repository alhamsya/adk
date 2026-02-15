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
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.ErrorStackMarshaler = marshalStack
	zerolog.ErrorHandler = func(err error) {
		slog.Error(err.Error())
	}
}

func setupDefaultDiode() {
	defaultDiode = diode.NewWriter(os.Stdout, diodeBufferSize, diodePollInterval, defaultDiodeDroppedLogFn)
}

func setupSlogDefault() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(defaultDiode, &slog.HandlerOptions{
		AddSource:   true,
		Level:       defaultLevel,
		ReplaceAttr: slogReplaceAttr,
	})))
}

func setupDefaultLogger() {
	defaultLogger = newLogger(defaultDiode, withLoggerCallerSkipFrameCount(zerolog.CallerSkipFrameCount+2))
	zerolog.DefaultContextLogger = &defaultLogger.Logger
}

// --- Helper Functions ---

func marshalStack(err error) interface{} {
	stack := pkgerrors.MarshalStack(err)
	st, ok := stack.([]map[string]string)
	if !ok {
		return stack
	}

	// Filter out runtime frames
	filtered := make([]map[string]string, 0, len(st))
	for _, frame := range st {
		src := frame["source"]
		// Filter out internal GO runtime frames
		if strings.Contains(src, "runtime/") || strings.HasSuffix(src, ".s") || src == "proc.go" {
			continue
		}
		filtered = append(filtered, frame)
	}
	return filtered
}

func slogReplaceAttr(groups []string, a slog.Attr) slog.Attr {
	switch a.Key {
	case slog.TimeKey:
		return slog.String(zerolog.TimestampFieldName, a.Value.Time().Format(time.RFC3339))
	case slog.LevelKey:
		return slog.String(slog.LevelKey, strings.ToLower(a.Value.String()))
	}
	return a
}

func defaultDiodeDroppedLogFn(dropCount int) {
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

// newLogger creates a new Logger instance with the given output writer and options.
func newLogger(output io.Writer, opts ...loggerOption) logger {
	config := loggerConfig{
		CallerSkipFrameCount: zerolog.CallerSkipFrameCount + 1,
	}

	for _, o := range opts {
		o(&config)
	}

	loggerInstance := zerolog.
		New(output).
		With().
		Timestamp().
		Stack().
		Logger()

	return logger{
		Logger: loggerInstance,
		Config: config,
	}
}

func withLoggerCallerSkipFrameCount(skipCount int) loggerOption {
	return func(cfg *loggerConfig) {
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
	l := zerolog.Ctx(ctx)
	for _, h := range hooks {
		loggerInstance := l.Hook(h).With().Logger()
		l = &loggerInstance
	}
	return l.WithContext(ctx), l
}
