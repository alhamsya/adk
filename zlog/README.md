# zlog

`zlog` is a structured logging package for ADK, built on top of `zerolog` and compatible with `log/slog`.

## Installation

To use `zlog` in your project, install it via `go get`:

```bash
go get github.com/alhamsya/adk/zlog
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
