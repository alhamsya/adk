package zlog

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"github.com/rs/zerolog/pkgerrors"
)

func init() {
	setDefaultLevel()
	zerolog.TimestampFieldName = "timestamp"
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	DefaultDiode = diode.NewWriter(os.Stdout, DiodeBufferSize, DiodePollInterval, DiodeDroppedLogFn)

	slog.SetDefault(slog.New(slog.NewJSONHandler(DefaultDiode, &slog.HandlerOptions{
		AddSource: true,
		Level:     defaultLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.TimeKey:
				return slog.Int64(zerolog.TimestampFieldName, a.Value.Time().UnixMilli())
			case slog.LevelKey:
				return slog.String(slog.LevelKey, strings.ToLower(a.Value.String()))
			}
			return a
		},
	})))

	DefaultLogger = New(NewLevelWriter(DefaultDiode), WithLoggerCallerSkipFrameCount(zerolog.CallerSkipFrameCount+2))

	zerolog.DefaultContextLogger = &DefaultLogger.Logger
	zerolog.ErrorHandler = func(err error) {
		slog.Error(err.Error())
	}
}

func DefaultDiodeDroppedLogFn(dropCount int) {
	slog.Warn(
		fmt.Sprintf(
			"zLog: you might logging things a bit too fast. We just dropped %d logs!!!",
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
	levelMap := map[string]slog.Level{
		LevelDebug: slog.LevelDebug,
		LevelInfo:  slog.LevelInfo,
		LevelWarn:  slog.LevelWarn,
		LevelError: slog.LevelError,
	}
	for k, v := range levelMap {
		if strings.EqualFold(k, str) {
			return v
		}
	}

	// default:
	return slog.LevelDebug
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

func New(output io.Writer, opts ...LoggerOption) zLog {
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

	return zLog{
		Logger: logger,
		Config: config,
	}
}

type levelWriter struct {
	w io.Writer
}

func (lw *levelWriter) Write(p []byte) (n int, err error) {
	return lw.w.Write(p)
}

func NewLevelWriter(w io.Writer) *levelWriter {
	return &levelWriter{w}
}

func WithLoggerCallerSkipFrameCount(skipCount int) LoggerOption {
	return func(cfg *LoggerConfig) {
		cfg.CallerSkipFrameCount = skipCount
	}
}
