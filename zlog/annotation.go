package zlog

import (
	"context"
	"encoding/json"
	"sync"
)

func CtxWithAnnotation(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	md := &Annotation{&sync.Map{}}

	return context.WithValue(ctx, loggingAnnotation{}, md)
}

func AnnotationFromCtx(ctx context.Context) *Annotation {
	val, ok := ctx.Value(loggingAnnotation{}).(*Annotation)
	if ok {
		return val
	}

	return &Annotation{&sync.Map{}}
}

func AddAnnotation(ctx context.Context, values map[string]any) {
	md := AnnotationFromCtx(ctx)
	for k, v := range values {
		md.Store(k, v)
	}
}

func (m *Annotation) toMap() map[string]any {
	res := map[string]any{}
	m.Range(func(key, value any) bool {
		res[key.(string)] = value
		return true
	})

	return res
}

func (m *Annotation) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.toMap())
}
