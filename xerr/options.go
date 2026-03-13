package xerr

import "github.com/pkg/errors"

type Option func(*Error)

func WithMessage(msg string) Option {
	return func(err *Error) {
		err.Message = msg
	}
}

func WithType(typ Type) Option {
	return func(err *Error) {
		err.Type = typ
		// also reset grpc code to match error type
		code, found := TypeToGRPCCode[typ]
		if found {
			err.Code = code
		}
	}
}

func WithMessageType() Option {
	return func(err *Error) {
		if err.Message == "" {
			err.Message = err.Type.String()
		}
	}
}

func WithStack() Option {
	return func(gerr *Error) {
		baseErr := gerr.err
		// no base error, no stacktrace
		if baseErr == nil {
			return
		}

		// don't override if already has stacktrace
		_, ok := baseErr.(stackTracer)
		if !ok {
			gerr.Wrap(errors.WithStack(baseErr))
		}
	}
}
