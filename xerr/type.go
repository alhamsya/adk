package xerr

import (
	"fmt"

	"google.golang.org/grpc/codes"
)

type Type string

func (t Type) String() string {
	return string(t)
}

type Error struct {
	Code    codes.Code
	Message string
	Type    Type
	err     error
}

func (e Error) Error() string {
	if e.err == nil {
		return e.Message
	}
	return fmt.Sprintf("%s: %s", e.Message, e.err.Error())
}
