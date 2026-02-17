package zlog

import (
	"context"
	"encoding/json"
	"sync"
)

func (f optionFunc) apply(a *Annotation) {
	f(a)
}

// DefaultAnnotation returns an Option that marks the annotation as default.
// When this option is used in CtxWithAnnotation, loggers created from the context
// (via FromContext) will automatically include the current annotations in every log entry.
func DefaultAnnotation() Option {
	return optionFunc(func(a *Annotation) {
		a.isDefault = true
	})
}

// CtxWithAnnotation initializes a new context with an empty Annotation map.
// If valid Options are provided, they modify the annotation behavior (e.g., enabling auto-injection).
func CtxWithAnnotation(ctx context.Context, opts ...Option) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	md := &Annotation{
		Map: &sync.Map{},
	}

	for _, opt := range opts {
		opt.apply(md)
	}

	return context.WithValue(ctx, loggingAnnotation{}, md)
}

func annotationFromCtx(ctx context.Context) *Annotation {
	val, ok := ctx.Value(loggingAnnotation{}).(*Annotation)
	if ok {
		return val
	}

	return &Annotation{&sync.Map{}, false}
}

// AddAnnotation adds key-value pairs to the annotation map stored in the context.
// These values will be logged if the context was configured with DefaultAnnotation.
func AddAnnotation(ctx context.Context, values map[string]any) {
	md := annotationFromCtx(ctx)
	for k, v := range values {
		md.Store(k, v)
	}
}

// AnnotationFromCtx returns the current Annotation object from the context.
// This is useful if you need to manually inspect or pass around the annotations.
func AnnotationFromCtx(ctx context.Context) *Annotation {
	return annotationFromCtx(ctx)
}

func (m *Annotation) toMap() map[string]any {
	res := map[string]any{}
	m.Range(func(key, value any) bool {
		res[key.(string)] = value
		return true
	})

	return res
}

// MarshalJSON conforms to the json.Marshaler interface.
func (m *Annotation) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.toMap())
}
