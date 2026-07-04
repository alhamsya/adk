# adk-zlog — Structured Logging

> Reference for the `github.com/alhamsya/adk/zlog` module. Written for humans and AI agents working in this codebase.

## What it is

`zlog` configures the process's global logging on import. It points `log/slog` at a stdlib JSON handler and configures `zerolog`'s globals separately — **both write through one shared `diode`** (lock-free, non-blocking) writer to stdout, with aligned field conventions so their output looks identical. You get JSON logs with RFC3339 timestamps, filtered stack traces (zerolog path), and per-request context annotations — while still logging through the standard `slog` API.

## When to use it

- You want standard `slog` calls to emit high-throughput structured JSON without touching zerolog directly.
- You need request-scoped fields (user id, request id) attached automatically.
- You accept **dropping logs under extreme load** in exchange for never blocking the caller.

## Install

```bash
go get github.com/alhamsya/adk/zlog
```

Dep: `github.com/rs/zerolog` (+ `diode`, `pkgerrors`).

## Import for side effects

Most of the value is in the package `init()`. Import it (blank if you only use `slog`):

```go
import _ "github.com/alhamsya/adk/zlog"
```

`init()` performs, once, at process start:
- Reads log level from `ZLOG_LEVEL` env (`debug`|`info`|`warn`|`error`, default `debug`).
- Sets zerolog globals: `timestamp` field, RFC3339 time format, stack marshaler, error handler → `slog.Error`.
- Creates the default diode writer to `os.Stdout` (buffer **1000**, poll **1ms**).
- Sets `slog.SetDefault` to a JSON handler writing through the diode (with source, level-lowercased, RFC3339 timestamp).
- Installs a default context logger for zerolog.

## Basic usage — via slog

```go
import (
    "log/slog"
    _ "github.com/alhamsya/adk/zlog"
)

slog.Info("hello", "user", "alhamsya")
slog.Warn("careful", "n", 3)
```

## Context loggers

```go
func FromContext(ctx context.Context) *zerolog.Logger
func NewContext(ctx context.Context, hooks ...zerolog.Hook) (context.Context, *zerolog.Logger)
```

- `FromContext` returns the ctx-scoped zerolog logger (or the global default). If the context carries a **default** annotation, the returned logger auto-includes it.
- `NewContext` attaches zerolog hooks and stores the logger back on the context.

```go
ctx, logger := zlog.NewContext(ctx)
logger.Info().Msg("action")
```

## Context annotations

Attach request-scoped key-values that ride along in every log line.

```go
// 1. Enable auto-injection on the context
ctx = zlog.CtxWithAnnotation(ctx, zlog.DefaultAnnotation())

// 2. Add fields (thread-safe, in-place)
zlog.AddAnnotation(ctx, map[string]any{"user_id": "12345", "request_id": "req-abc"})

// 3. Log — annotation is attached automatically
zlog.FromContext(ctx).Info().Msg("action performed")
// {"level":"info","annotation":{"user_id":"12345","request_id":"req-abc"},"timestamp":"...","message":"action performed"}
```

API:

```go
DefaultAnnotation() Option                              // mark annotation as auto-injected
CtxWithAnnotation(ctx, opts ...Option) context.Context  // seed an annotation map on ctx
AddAnnotation(ctx, values map[string]any)               // merge fields (sync.Map under the hood)
AnnotationFromCtx(ctx) *Annotation                      // inspect/pass around manually
```

`Annotation` wraps a `sync.Map` and implements `json.Marshaler`.

## Stack traces

Errors implementing `pkg/errors` `StackTracer` are logged with a `stack` field. Runtime frames (`runtime/`, `*.s`, `proc.go`) are filtered out for readability.

```go
import "github.com/pkg/errors"
err := errors.New("boom")
zlog.FromContext(ctx).Error().Err(err).Msg("operation failed") // includes filtered stack
```

## Gotchas

- **Logs can be dropped by design.** The diode is non-blocking; on buffer overflow it drops entries and emits `slog.Warn("zLog: dropped N logs...")`. This is the throughput/latency trade-off — don't rely on zlog for audit-grade delivery.
- **Annotations are only auto-injected when `DefaultAnnotation()` was passed** to `CtxWithAnnotation`. Without it, `AddAnnotation` still stores values, but `FromContext` won't attach them — you'd read them manually via `AnnotationFromCtx`.
- **Log level is read once at `init` from `ZLOG_LEVEL`.** Set the env before the process starts; there is no exported runtime setter. Default level is `debug`.
- **It configures process-global state.** Importing `zlog` changes `slog`'s default handler and zerolog globals for the whole process — intended, but be aware in libraries/tests.
- **Output goes to `os.Stdout`** via the default diode; there is no exported constructor to redirect it (the logger constructor and diode are unexported).
