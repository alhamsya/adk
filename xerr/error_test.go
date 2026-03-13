package xerr_test

import (
	"github.com/alhamsya/adk/xerr"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"reflect"
	"testing"
)

func TestError_GRPCStatus(t *testing.T) {
	tests := []struct {
		name string
		e    func() *xerr.Error
		want func() *status.Status
	}{
		{
			name: "correct Status given err preset",
			e:    func() *xerr.Error { return xerr.ErrUnauthorized },
			want: func() *status.Status {
				grpcStatus := status.New(codes.Unauthenticated, xerr.ErrUnauthorized.Error())
				return grpcStatus
			},
		},
		{
			name: "correct Status given OK",
			e: func() *xerr.Error {
				return &xerr.Error{
					Code:    codes.OK,
					Message: "OK",
					Type:    "OK",
				}
			},
			want: func() *status.Status {
				return status.New(codes.OK, "OK")
			},
		},
		{
			name: "correct Status given wrapped err",
			e: func() *xerr.Error {
				er := xerr.New(xerr.TypeUnauthorized, "sukab isn't allowed")
				er.Wrap(io.ErrClosedPipe)
				return er
			},
			want: func() *status.Status {
				grpcStatus := status.New(codes.Unauthenticated, "gRPC error 16: sukab isn't allowed: io: read/write on closed pipe")
				return grpcStatus
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.e().GRPCStatus()
			want := tt.want()

			if !reflect.DeepEqual(got.Details(), want.Details()) {
				t.Errorf("Error.GRPCStatus().Details() = %v, want %v", got.Details(), want.Details())
				t.Fail()
			}

			if !reflect.DeepEqual(got.Code(), want.Code()) {
				t.Errorf("Error.GRPCStatus().Code() = %v, want %v", got.Code(), want.Code())
				t.Fail()
			}

			if !reflect.DeepEqual(got.Message(), want.Message()) {
				t.Errorf("Error.GRPCStatus().Message() = %v, want %v", got.Message(), want.Message())
				t.Fail()
			}
		})
	}
}
