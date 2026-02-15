package zlog

import (
	"context"
	"encoding/json"
	"sync"
)

func CtxWithMetadata(ctx context.Context, md *Metadata) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if md == nil {
		md = &Metadata{&sync.Map{}}
	}

	return context.WithValue(ctx, loggingMetadata{}, md)
}

func MetadataFromCtx(ctx context.Context) *Metadata {
	val, ok := ctx.Value(loggingMetadata{}).(*Metadata)
	if ok {
		return val
	}

	return &Metadata{&sync.Map{}}
}

func InjectMetadata(ctx context.Context, values map[string]any) {
	md := MetadataFromCtx(ctx)
	for k, v := range values {
		md.Store(k, v)
	}
}

func (m *Metadata) toMap() map[string]any {
	res := map[string]any{}
	m.Range(func(key, value any) bool {
		res[key.(string)] = value
		return true
	})

	return res
}

func (m *Metadata) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.toMap())
}
