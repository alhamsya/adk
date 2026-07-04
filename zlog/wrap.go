package zlog

import "github.com/rs/zerolog"

// Logger wraps zerolog.Logger so the timestamp is written immediately after the
// level field on every event (level, timestamp, ...), instead of zerolog's
// default where the timestamp hook runs at Msg time and lands after any event
// fields. Returned by FromContext and NewContext.
//
// All other zerolog.Logger methods are promoted unchanged. Note that methods
// which return a bare zerolog.Logger (e.g. Level, Sample, Output, or
// With().Logger()) yield a raw logger without this timestamp placement.
type Logger struct {
	zerolog.Logger
	annotation *Annotation
}

// with places the timestamp right after the level, then the default annotation
// (if any), before any caller-supplied fields.
func (l *Logger) with(e *zerolog.Event) *zerolog.Event {
	e = e.Timestamp()
	if l.annotation != nil {
		e = e.Interface("annotation", l.annotation)
	}
	return e
}

func (l *Logger) Trace() *zerolog.Event { return l.with(l.Logger.Trace()) }
func (l *Logger) Debug() *zerolog.Event { return l.with(l.Logger.Debug()) }
func (l *Logger) Info() *zerolog.Event  { return l.with(l.Logger.Info()) }
func (l *Logger) Warn() *zerolog.Event  { return l.with(l.Logger.Warn()) }
func (l *Logger) Error() *zerolog.Event { return l.with(l.Logger.Error()) }

func (l *Logger) Err(err error) *zerolog.Event { return l.with(l.Logger.Err(err)) }
func (l *Logger) Fatal() *zerolog.Event        { return l.with(l.Logger.Fatal()) }
func (l *Logger) Panic() *zerolog.Event        { return l.with(l.Logger.Panic()) }
func (l *Logger) Log() *zerolog.Event          { return l.with(l.Logger.Log()) }

func (l *Logger) WithLevel(level zerolog.Level) *zerolog.Event {
	return l.with(l.Logger.WithLevel(level))
}
