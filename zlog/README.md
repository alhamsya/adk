# zlog

`zlog` is a structured logging package for ADK, built on top of `zerolog` and compatible with `log/slog`.

## Installation

To use `zlog` in your project, install it via `go get`:

```bash
go get github.com/alhamsya/adk/zlog
```

## Features

### Timestamp Format (Timestampz)

`zlog` automatically configures `zerolog` and `slog` to use **RFC3339** format (e.g., `2026-02-15T15:04:05+07:00`) for the `timestamp` field. This ensures compatibility with systems expecting `timestampz`.

### Annotation Injection

You can inject annotation into the context, which will be automatically included in all subsequent log entries created from that context.

#### 1. Initialize Context
First, initialize the context to hold annotation. You can pass options to configure behavior.
To enable automatic injection into the logger (from `zlog.FromContext`), use `zlog.DefaultAnnotation()`.

```go
ctx := context.Background()
ctx = zlog.CtxWithAnnotation(ctx, zlog.DefaultAnnotation())
```

#### 2. Inject Annotation
Inject annotation key-value pairs. This updates the annotation in-place (thread-safe).

```go
zlog.AddAnnotation(ctx, map[string]any{
    "user_id": "12345",
    "request_id": "req-abc-789",
})
```

#### 3. Log with Annotation
When you create a logger from this context, the annotation is automatically attached if `DefaultAnnotation()` was used.

```go
logger := zlog.FromContext(ctx)
logger.Info().Msg("Action performed")
// Output: {"level":"info","annotation":{"user_id":"12345","request_id":"req-abc-789"},"timestamp":"...","message":"Action performed"}
```

### Stack Traces

`zlog` is configured to automatically log stack traces for errors that implement the `pkg/errors.StackTracer` interface.

```go
import "github.com/pkg/errors"

err := errors.New("something went wrong")
logger.Error().Err(err).Msg("Operation failed")
// Output will include a "stack" field with the stack trace.
```

## Usage

### Basic Usage

`zlog` automatically configures `log/slog` to use `zerolog` backend with a diode writer for high performance.

```go
package main

import (
    "log/slog"
    "time"
    _ "github.com/alhamsya/adk/zlog" // Import for side-effects (init)
)

func main() {
    // Standard slog usage
    slog.Info("Hello World", "user", "alhamsya")
    
    // With timestamp
    slog.Warn("Warning message", "time", time.Now())
}
```

### Advanced Usage (Context)

You can use `zlog.NewContext` and `zlog.FromContext` to manage loggers with context.

```go
import "github.com/alhamsya/adk/zlog"

func handler(ctx context.Context) {
    // Add fields to logger in context
    ctx, logger := zlog.NewContext(ctx)
    logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
        return c.Str("request_id", "123")
    })
    
    // Retrieve logger
    log := zlog.FromContext(ctx)
    log.Info().Msg("Request processed")
}
```
