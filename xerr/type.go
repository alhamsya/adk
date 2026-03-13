package xerr

import (
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
	//TODO implement me
	panic("implement me")
}
