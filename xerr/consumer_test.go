package xerr

import (
	"errors"
	"testing"
)

func TestRespMsg(t *testing.T) {
	tests := []struct {
		name          string
		setup         func() *RespMsg
		wantIsError   bool
		wantRequeue   bool
		wantAck       bool
	}{
		{
			name: "success scenario",
			setup: func() *RespMsg {
				return New()
			},
			wantIsError: false,
			wantRequeue: false,
			wantAck:     true,
		},
		{
			name: "error with requeue",
			setup: func() *RespMsg {
				return New().Err(errors.New("db error"))
			},
			wantIsError: true,
			wantRequeue: true,
			wantAck:     false,
		},
		{
			name: "error with ignore",
			setup: func() *RespMsg {
				return New().Err(errors.New("invalid payload")).Ignored()
			},
			wantIsError: true,
			wantRequeue: false,
			wantAck:     true,
		},
		{
			name: "ignore without error",
			setup: func() *RespMsg {
				return New().Ignored()
			},
			wantIsError: false,
			wantRequeue: false,
			wantAck:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.setup()
			if got := r.IsError(); got != tt.wantIsError {
				t.Errorf("IsError() = %v, want %v", got, tt.wantIsError)
			}
			if got := r.ShouldRequeue(); got != tt.wantRequeue {
				t.Errorf("ShouldRequeue() = %v, want %v", got, tt.wantRequeue)
			}
			if got := r.ShouldAck(); got != tt.wantAck {
				t.Errorf("ShouldAck() = %v, want %v", got, tt.wantAck)
			}
		})
	}
}

func TestRespMsg_Reset(t *testing.T) {
	r := New().Err(errors.New("error")).Ignored()
	r.Reset()
	if r.Error != nil || r.Ignore != false {
		t.Errorf("Reset() failed, state: %+v", r)
	}
}
