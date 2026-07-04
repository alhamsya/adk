package zlog

import (
	"bytes"
	"context"
	"strings"
	"sync"
	"testing"

	"github.com/rs/zerolog"
)

// assertOrder fails unless keys appear in s in the given order.
func assertOrder(t *testing.T, s string, keys ...string) {
	t.Helper()
	prev := -1
	for _, k := range keys {
		i := strings.Index(s, k)
		if i < 0 {
			t.Fatalf("missing %s in: %s", k, s)
		}
		if i < prev {
			t.Fatalf("%s out of order in: %s", k, s)
		}
		prev = i
	}
}

func TestLogger_TimestampRightAfterLevel(t *testing.T) {
	// With a caller field: level, timestamp, user.
	var b1 bytes.Buffer
	l1 := &Logger{Logger: zerolog.New(&b1).With().Stack().Logger()}
	l1.Info().Str("user", "x").Msg("hello")
	assertOrder(t, b1.String(), `"level"`, `"timestamp"`, `"user"`)

	// With a default annotation: level, timestamp, annotation.
	var b2 bytes.Buffer
	ann := &Annotation{Map: &sync.Map{}, isDefault: true}
	ann.Store("request_id", "req-1")
	l2 := &Logger{Logger: zerolog.New(&b2).With().Stack().Logger(), annotation: ann}
	l2.Info().Str("user", "x").Msg("hi")
	assertOrder(t, b2.String(), `"level"`, `"timestamp"`, `"annotation"`, `"user"`)
}

// TestFromContext_Order exercises the real FromContext path (annotation
// extraction + wrapper) against a captured buffer.
func TestFromContext_Order(t *testing.T) {
	var buf bytes.Buffer
	base := zerolog.New(&buf).With().Stack().Logger()
	ctx := base.WithContext(context.Background())
	ctx = CtxWithAnnotation(ctx, DefaultAnnotation())
	AddAnnotation(ctx, map[string]any{"request_id": "req-1"})

	FromContext(ctx).Info().Str("user", "x").Msg("hi")
	assertOrder(t, buf.String(), `"level"`, `"timestamp"`, `"annotation"`, `"request_id"`, `"user"`)
}
