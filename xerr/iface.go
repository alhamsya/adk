package xerr

import (
	"fmt"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
)

type GRPCErr interface {
	fmt.Stringer
	GRPCCode
	GRPCErrType
}

type GRPCCode interface {
	// Code should return a gRPC Status Code.
	Code() codes.Code
}

type GRPCErrType interface {
	// Type should return a Stockbit's Error Type.
	Type() Type
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}
