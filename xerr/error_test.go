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
				return status.New(codes.Unauthenticated, "sukab isn't allowed: io: read/write on closed pipe")
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

			if !reflect.DeepEqual(got.Message(), want.Message()) {
				t.Errorf("Error.GRPCStatus().Message() = %v, want %v", got.Message(), want.Message())
				t.Fail()
			}
		})
	}
}

func TestError_HTTPStatus(t *testing.T) {
	tests := []struct {
		name string
		e    *xerr.Error
		want int
	}{
		{
			name: "Unauthorized",
			e:    xerr.ErrUnauthorized,
			want: 401,
		},
		{
			name: "NotFound",
			e:    xerr.ErrNotFound,
			want: 404,
		},
		{
			name: "SystemError",
			e:    xerr.ErrSystemError,
			want: 500,
		},
		{
			name: "OK",
			e: &xerr.Error{
				Type: xerr.TypeOK,
			},
			want: 200,
		},
		{
			name: "Unknown type defaults to 500",
			e: &xerr.Error{
				Type: "UNKNOWN_TYPE",
			},
			want: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.HTTPStatus(); got != tt.want {
				t.Errorf("Error.HTTPStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_RefinedMappings(t *testing.T) {
	tests := []struct {
		name     string
		e        *xerr.Error
		wantHTTP int
		wantGRPC codes.Code
	}{
		{
			name:     "DuplicateCall",
			e:        xerr.ErrDuplicateCall,
			wantHTTP: 409,
			wantGRPC: codes.AlreadyExists,
		},
		{
			name:     "AlreadyExists",
			e:        xerr.ErrAlreadyExists,
			wantHTTP: 409,
			wantGRPC: codes.AlreadyExists,
		},
		{
			name:     "Unimplemented",
			e:        xerr.ErrUnimplemented,
			wantHTTP: 501,
			wantGRPC: codes.Unimplemented,
		},
		{
			name:     "Aborted",
			e:        xerr.ErrAborted,
			wantHTTP: 409,
			wantGRPC: codes.Aborted,
		},
		{
			name:     "RequestCanceled",
			e:        xerr.ErrRequestCanceled,
			wantHTTP: 499,
			wantGRPC: codes.Canceled,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.HTTPStatus(); got != tt.wantHTTP {
				t.Errorf("%s: HTTPStatus() = %v, want %v", tt.name, got, tt.wantHTTP)
			}
			if got := tt.e.Code; got != tt.wantGRPC {
				t.Errorf("%s: Code = %v, want %v", tt.name, got, tt.wantGRPC)
			}
		})
	}
}

func TestFromError_PointerSafety(t *testing.T) {
	// Original preset
	original := xerr.ErrUnauthorized
	originalMsg := original.Message

	// Create new error from preset with a custom message
	customMsg := "custom message"
	newErr := xerr.FromError(original, xerr.WithMessage(customMsg))

	if newErr.Message != customMsg {
		t.Errorf("newErr.Message = %v, want %v", newErr.Message, customMsg)
	}

	// Verify the global original preset wasn't modified
	if original.Message != originalMsg {
		t.Errorf("global preset was modified! original.Message = %v, want %v", original.Message, originalMsg)
	}
}

func TestError_Comparison(t *testing.T) {
	e1 := xerr.New(xerr.TypeUnauthorized, "msg1")
	e2 := xerr.New(xerr.TypeUnauthorized, "msg2")
	e3 := xerr.New(xerr.TypeUnauthorized, "msg1")

	if e1.Is(e2) {
		t.Errorf("e1 should not match e2 (different message)")
	}

	if !e1.Is(e3) {
		t.Errorf("e1 should match e3 (same content)")
	}
}
