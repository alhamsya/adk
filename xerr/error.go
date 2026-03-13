package xerr

import (
	"fmt"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func New(errType Type, msg string) *Error {
	c, cFound := TypeToGRPCCode[errType]
	if !cFound {
		c = codes.Internal
	}
	return &Error{
		Code:    c,
		Message: msg,
		Type:    errType,
	}
}

func NewWithWrap(errType Type, cause error, msg string) *Error {
	err := New(errType, msg)
	if cause != nil {
		err.Wrap(cause)
	}
	return err
}

func (e *Error) Wrap(err error) {
	e.err = err
}

func (e Error) Is(err error) bool {
	var asPtr *Error
	var asNonPtr Error
	if errors.As(err, &asPtr) {
		return cmp(&e, asPtr)
	}
	if errors.As(err, &asNonPtr) {
		return cmp(&e, &asNonPtr)
	}
	return false
}

func cmp(e1 *Error, e2 *Error) bool {
	return e1.Code == e2.Code &&
		e1.Message == e2.Message &&
		e1.Type == e2.Type
}

func (e Error) Unwrap() error {
	return e.err
}

func Wrap(grpcErr *Error, err error) error {
	grpcErr.Wrap(err)
	return grpcErr
}

func (e Error) GRPCStatus() *status.Status {
	grpcStatus := status.New(e.Code, e.Error())
	return grpcStatus
}

func (e Error) HTTPStatus() int {
	sts, found := TypeToHTTPStatus[e.Type]
	if !found {
		return 500
	}
	return sts
}

func FromError(err error, opts ...Option) *Error {
	if err == nil {
		return &Error{
			Code:    codes.OK,
			Message: "",
			Type:    TypeOK,
		}
	}

	var ownErr *Error
	if errors.As(err, &ownErr) {
		// if err already an *Error,
		// and we have options, return a copy to avoid modifying the original (e.g. global presets)
		if len(opts) > 0 {
			res := *ownErr
			for _, opt := range opts {
				opt(&res)
			}
			return &res
		}
		return ownErr
	}

	msg := err.Error()
	var customMsg fmt.Stringer
	if errors.As(err, &customMsg) {
		msg = customMsg.String()
	}

	cde := codes.Internal
	var customCode GRPCCode
	if errors.As(err, &customCode) {
		cde = customCode.Code()
	}

	errorType := TypeSystemError
	var customErrorType GRPCErrType
	if errors.As(err, &customErrorType) {
		errorType = customErrorType.Type()
	}

	ownErr = &Error{
		Code:    cde,
		Message: msg,
		Type:    errorType,
		err:     err,
	}

	for _, opt := range opts {
		opt(ownErr)
	}

	return ownErr
}

func GetType(err error) Type {
	if err == nil {
		return TypeOK
	}

	var e *Error
	if errors.As(err, &e) {
		return e.Type
	}
	return TypeUnknown
}
