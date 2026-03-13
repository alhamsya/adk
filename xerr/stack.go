package xerr

import "github.com/pkg/errors"

func (e Error) StackTrace() errors.StackTrace {
	if e.err == nil {
		return nil
	}
	st, ok := e.err.(stackTracer)
	if !ok {
		return nil
	}
	return st.StackTrace()
}
