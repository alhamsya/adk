<p align="center">
  <h1 align="center">🚀 ADK — Alhamsya Development Kit</h1>
  <p align="center">
    A collection of production-ready Go packages for building robust backend services.
    <br />
    <strong>Structured Errors · HTTP Client · Structured Logging</strong>
  </p>
</p>

<p align="center">
  <a href="https://pkg.go.dev/github.com/alhamsya/adk"><img src="https://pkg.go.dev/badge/github.com/alhamsya/adk.svg" alt="Go Reference"></a>
  <img src="https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License">
</p>

---

## 📦 Overview

**ADK** (Alhamsya Development Kit) is a modular Go SDK that provides essential utilities for building production-ready backend services. Each package is designed as an independent Go module, allowing you to import only what you need.

| Package | Description | Import Path |
|---------|-------------|-------------|
| **[`xerr`](#-xerr--structured-error-handling)** | Structured error handling with automatic gRPC & HTTP status code mapping | `github.com/alhamsya/adk/xerr` |
| **[`xhttp`](#-xhttp--http-client-with-retry)** | HTTP client with built-in retry and config-driven upstream management | `github.com/alhamsya/adk/xhttp` |
| **[`zlog`](#-zlog--structured-logging)** | Structured logging powered by zerolog with diode buffering & context annotations | `github.com/alhamsya/adk/zlog` |

---

## 🛠 Installation

Each package is an independent Go module. Install only what you need:

```bash
# Structured Error Handling
go get github.com/alhamsya/adk/xerr

# HTTP Client with Retry
go get github.com/alhamsya/adk/xhttp

# Structured Logging
go get github.com/alhamsya/adk/zlog
```

**Requirements:** Go 1.24+

---

## ⚡ Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log/slog"
    "time"

    "github.com/alhamsya/adk/xerr"
    "github.com/alhamsya/adk/xhttp"
    "github.com/alhamsya/adk/zlog"
)

func main() {
    // 🔴 Structured Errors
    err := xerr.New(xerr.TypeNotFound, "user not found")
    fmt.Println(err.HTTPStatus()) // 404
    fmt.Println(err.GRPCStatus()) // code: NotFound

    // 🌐 HTTP Client
    client := xhttp.NewClient()
    client.SetTimeout(10 * time.Second)
    res, _ := client.Get(context.Background(), "https://api.example.com/health")
    fmt.Println(res.StatusCode)

    // 📝 Structured Logging
    ctx := zlog.CtxWithAnnotation(context.Background(), zlog.DefaultAnnotation())
    zlog.AddAnnotation(ctx, map[string]any{"request_id": "req-123"})
    zlog.FromContext(ctx).Info().Msg("service started")

    // Also works with standard slog
    slog.Info("hello from slog", "user", "alhamsya")

    time.Sleep(100 * time.Millisecond) // wait for async log flush
}
```

---

## 🔴 `xerr` — Structured Error Handling

Package `xerr` provides structured error handling that automatically maps errors to both **gRPC status codes** and **HTTP status codes**. Designed for services that serve multiple protocols (REST + gRPC) simultaneously.

### Installation

```bash
go get github.com/alhamsya/adk/xerr
```

### Key Features

- ✅ **Dual-protocol mapping** — one error type automatically maps to both gRPC & HTTP codes
- ✅ **17 pre-defined error types** — covers common backend service scenarios
- ✅ **Pre-built error sentinels** — ready to use without re-initialization
- ✅ **Functional options pattern** — flexible error customization
- ✅ **Error wrapping/unwrapping** — fully compatible with `errors.Is()` and `errors.As()`
- ✅ **Stack trace support** — integration with `pkg/errors`
- ✅ **Pointer-safe** — `FromError` with options never modifies global presets

### Creating Errors

```go
// Create a new error
err := xerr.New(xerr.TypeNotFound, "user not found")

// Create an error with a wrapped cause
err := xerr.NewWithWrap(xerr.TypeSystemError, originalErr, "database connection failed")

// Use pre-built sentinel errors
if errors.Is(err, xerr.ErrUnauthorized) {
    // handle unauthorized
}
```

### Using Pre-built Error Sentinels

ADK provides ready-to-use error sentinels:

```go
// Use directly
return xerr.ErrNotFound           // 404 / NotFound
return xerr.ErrUnauthorized       // 401 / Unauthenticated
return xerr.ErrInvalidParameter   // 400 / InvalidArgument
return xerr.ErrSystemError        // 500 / Internal
return xerr.ErrServiceBusy        // 429 / ResourceExhausted

// Customize from sentinel — SAFE, does not modify the global preset
err := xerr.FromError(xerr.ErrNotFound, xerr.WithMessage("product ID 123 not found"))
```

### Converting to gRPC & HTTP

```go
err := xerr.New(xerr.TypeUnauthorized, "invalid token")

// gRPC
grpcStatus := err.GRPCStatus()
// grpcStatus.Code() == codes.Unauthenticated

// HTTP
httpCode := err.HTTPStatus()
// httpCode == 401
```

### Error Wrapping & Unwrapping

```go
// Wrap an external error
dbErr := sql.ErrNoRows
err := xerr.NewWithWrap(xerr.TypeNotFound, dbErr, "user not found")

// Unwrap
cause := errors.Unwrap(err) // sql.ErrNoRows

// Comparison
errors.Is(err, anotherXerr) // compares Code, Message, and Type
```

### Functional Options

```go
// Convert from a generic error with options
err := xerr.FromError(someError,
    xerr.WithMessage("custom message"),      // override message
    xerr.WithType(xerr.TypeNotFound),        // override type (also updates gRPC code)
    xerr.WithStack(),                        // attach stack trace
    xerr.WithMessageType(),                  // use type name as message if empty
)
```

### Error Type Mappings

Complete mapping table between error `Type`, HTTP Status Code, and gRPC Code:

| Type | Value | HTTP | gRPC Code |
|------|-------|------|-----------|
| `TypeOK` | `OK` | `200 OK` | `OK (0)` |
| `TypeInvalidParameter` | `INVALID_PARAMETER` | `400 Bad Request` | `InvalidArgument (3)` |
| `TypeUnauthorized` | `UNAUTHORIZED` | `401 Unauthorized` | `Unauthenticated (16)` |
| `TypeForbidden` | `FORBIDDEN` | `403 Forbidden` | `PermissionDenied (7)` |
| `TypeNoTradingPermission` | `NO_TRADING_PERMISSION` | `403 Forbidden` | `PermissionDenied (7)` |
| `TypeNoSubscription` | `NO_SUBSCRIPTION` | `403 Forbidden` | `PermissionDenied (7)` |
| `TypeNotFound` | `NOT_FOUND` | `404 Not Found` | `NotFound (5)` |
| `TypeDuplicateCall` | `DUPLICATE_CALL` | `409 Conflict` | `AlreadyExists (6)` |
| `TypeAlreadyExists` | `ALREADY_EXISTS` | `409 Conflict` | `AlreadyExists (6)` |
| `TypeAborted` | `ABORTED` | `409 Conflict` | `Aborted (10)` |
| `TypeServiceBusy` | `SERVICE_BUSY` | `429 Too Many Requests` | `ResourceExhausted (8)` |
| `TypeRequestCanceled` | `REQUEST_CANCELED` | `499 Client Closed Request` | `Canceled (1)` |
| `TypeSystemError` | `SYSTEM_ERROR` | `500 Internal Server Error` | `Internal (13)` |
| `TypeVendorError` | `VENDOR_ERROR` | `500 Internal Server Error` | `Internal (13)` |
| `TypeUnimplemented` | `UNIMPLEMENTED` | `501 Not Implemented` | `Unimplemented (12)` |
| `TypeBadGateway` | `BAD_GATEWAY` | `502 Bad Gateway` | `Unavailable (14)` |
| `TypeMaintenance` | `MAINTENANCE` | `503 Service Unavailable` | `Unavailable (14)` |
| `TypeGatewayTimeout` | `GATEWAY_TIMEOUT` | `504 Gateway Timeout` | `DeadlineExceeded (4)` |

### Interfaces

`xerr` provides interfaces for custom integration:

```go
// Implement these interfaces on your custom errors
type GRPCCode interface {
    Code() codes.Code
}

type GRPCErrType interface {
    Type() Type
}

type HTTPStatusCode interface {
    HTTPStatus() int
}
```

### Example: gRPC Middleware

```go
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
        resp, err := handler(ctx, req)
        if err != nil {
            xerrErr := xerr.FromError(err)
            return nil, xerrErr.GRPCStatus().Err()
        }
        return resp, nil
    }
}
```

### Example: HTTP Error Handler

```go
func ErrorHandler(w http.ResponseWriter, err error) {
    xerrErr := xerr.FromError(err)
    w.WriteHeader(xerrErr.HTTPStatus())
    json.NewEncoder(w).Encode(map[string]string{
        "error":   xerrErr.Type.String(),
        "message": xerrErr.Message,
    })
}
```

---

## 🌐 `xhttp` — HTTP Client with Retry

Package `xhttp` provides an HTTP client with built-in retry support and config-driven upstream management. Ideal for service-to-service communication that requires resilience.

### Installation

```bash
go get github.com/alhamsya/adk/xhttp
```

### Key Features

- ✅ **Built-in retry** — automatic retries with `retry-go` (exponential & fixed backoff)
- ✅ **Config-driven upstreams** — manage multiple upstream services via JSON/YAML configuration
- ✅ **JSON auto-serialization** — request body is automatically marshaled to JSON
- ✅ **Response unmarshaling** — easily unmarshal response body into structs
- ✅ **Per-upstream timeout** — configure timeout per upstream service
- ✅ **Extras metadata** — store additional metadata per upstream (API keys, partner IDs, etc.)

### Basic Usage

```go
client := xhttp.NewClient()
client.SetTimeout(10 * time.Second)

// GET request
res, err := client.Get(ctx, "https://api.example.com/users/123")
if err != nil {
    log.Fatal(err)
}

fmt.Println(res.StatusCode)    // 200
fmt.Println(string(res.Body))  // response body

// Unmarshal response
var user User
err = res.Unmarshal(&user)
```

### POST / PUT / DELETE

```go
// POST with body — automatically serialized to JSON
body := map[string]string{"name": "Alhamsya", "role": "engineer"}
res, err := client.Post(ctx, "https://api.example.com/users", body)

// PUT
res, err := client.Put(ctx, "https://api.example.com/users/123", updatedUser)

// DELETE
res, err := client.Delete(ctx, "https://api.example.com/users/123")
```

### Config-Driven Upstream Management

A more powerful way to manage multiple upstream services:

#### Configuration (JSON)

```json
{
  "upstream": {
    "payment-service": {
      "host": "https://payment.internal.svc",
      "timeout": "5s",
      "type": "rest",
      "global-retry": {
        "attempts": 3,
        "delay": "100ms",
        "type": "exponential"
      },
      "endpoints": {
        "create-payment": { "path": "/api/v1/payments" },
        "get-status":     { "path": "/api/v1/payments/status" }
      },
      "extras": {
        "api-key": "sk_live_xxx",
        "partner-id": "PARTNER001"
      }
    },
    "notification-service": {
      "host": "https://notif.internal.svc",
      "timeout": "3s",
      "endpoints": {
        "send-email": { "path": "/api/v1/email/send" },
        "send-push":  { "path": "/api/v1/push/send" }
      }
    }
  }
}
```

#### Configuration (YAML)

```yaml
upstream:
  payment-service:
    host: https://payment.internal.svc
    timeout: 5s
    type: rest
    global-retry:
      attempts: 3
      delay: 100ms
      type: exponential
    endpoints:
      create-payment:
        path: /api/v1/payments
      get-status:
        path: /api/v1/payments/status
    extras:
      api-key: sk_live_xxx
      partner-id: PARTNER001
```

#### Usage

```go
// Initialize client with config
client := xhttp.NewClient()
client.SetConfig(cfg) // cfg loaded from JSON/YAML

// Make requests using upstream & endpoint names
res, err := client.GetEndpoint(ctx, "payment-service", "get-status")
res, err := client.PostEndpoint(ctx, "payment-service", "create-payment", paymentReq)
res, err := client.PutEndpoint(ctx, "notification-service", "send-email", emailReq)
res, err := client.DeleteEndpoint(ctx, "payment-service", "cancel-payment")
```

### Retry Strategies

| Strategy | Constant | Description |
|----------|----------|-------------|
| **Exponential Backoff** | `xhttp.RetryTypeExponential` | Delay increases exponentially between retries |
| **Fixed Delay** | `xhttp.RetryTypeFixed` | Constant delay between retries |

### Struct Reference

```go
// Main configuration
type Config struct {
    Upstream map[string]Upstream
}

type Upstream struct {
    Host        string                 // Base URL of the upstream service
    Timeout     string                 // Per-request timeout (e.g., "5s")
    Type        string                 // Service type (e.g., "rest", "grpc")
    GlobalRetry RetryConfig            // Global retry configuration
    Endpoints   map[string]Endpoint    // Endpoint definitions
    Extras      map[string]interface{} // Additional metadata
}

type Endpoint struct {
    Path string // Endpoint path (e.g., "/api/v1/users")
}

type RetryConfig struct {
    Attempts int    // Number of retry attempts
    Delay    string // Delay between retries (e.g., "100ms")
    Timeout  string // Timeout per retry attempt
    Type     string // "exponential" or "fixed"
}

// Response
type Res struct {
    Header     http.Header
    StatusCode int
    Body       []byte
}
```

---

## 📝 `zlog` — Structured Logging

Package `zlog` provides high-performance structured logging built on top of `zerolog` with `log/slog` compatibility. It features diode buffering for non-blocking I/O and context-scoped annotations.

### Installation

```bash
go get github.com/alhamsya/adk/zlog
```

### Key Features

- ✅ **High-performance** — diode-buffered writer for async, non-blocking logging
- ✅ **Dual API** — use `zerolog` directly OR the standard `log/slog` library
- ✅ **Context annotations** — inject metadata into context, automatically included in every log entry
- ✅ **RFC3339 timestamps** — standard timestamp format (field: `timestamp`)
- ✅ **Stack trace support** — automatically logs stack traces for `pkg/errors` errors
- ✅ **Environment-based level** — set log level via the `ZLOG_LEVEL` env variable
- ✅ **Thread-safe annotations** — uses `sync.Map` for concurrent access

### Basic Usage (Side-Effect Import)

Simply import `zlog` to automatically configure `log/slog`:

```go
package main

import (
    "log/slog"
    _ "github.com/alhamsya/adk/zlog" // import for side-effects
)

func main() {
    slog.Info("Hello World", "user", "alhamsya")
    // Output: {"level":"info","timestamp":"2026-03-14T19:00:00+07:00","source":{"function":"..."},"msg":"Hello World","user":"alhamsya"}
}
```

### Context-Scoped Annotations

Inject metadata into the context that is automatically included in every log entry:

```go
// 1. Initialize context with annotation support
ctx := context.Background()
ctx = zlog.CtxWithAnnotation(ctx, zlog.DefaultAnnotation())

// 2. Add annotations
zlog.AddAnnotation(ctx, map[string]any{
    "user_id":    "12345",
    "request_id": "req-abc-789",
    "service":    "payment",
})

// 3. Log — annotations are automatically attached
logger := zlog.FromContext(ctx)
logger.Info().Msg("Payment processed")
// Output:
// {
//   "level": "info",
//   "annotation": {"user_id": "12345", "request_id": "req-abc-789", "service": "payment"},
//   "timestamp": "2026-03-14T19:00:00+07:00",
//   "message": "Payment processed"
// }
```

### Stack Traces

Errors from `pkg/errors` automatically include stack traces:

```go
import "github.com/pkg/errors"

err := errors.New("database connection timeout")
logger.Error().Err(err).Msg("Failed to query database")
// Output includes a "stack" field with filtered stack frames
```

### Context Logger with Hooks

```go
func handler(ctx context.Context) {
    // Create a new logger with hooks
    ctx, logger := zlog.NewContext(ctx, myCustomHook)

    // Update context fields
    logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
        return c.Str("handler", "payment")
    })

    // Retrieve logger from context
    log := zlog.FromContext(ctx)
    log.Info().Msg("Handler executed")
}
```

### Configuring Log Level

Set the log level via environment variable:

```bash
# Options: debug, info, warn, error
export ZLOG_LEVEL=info
```

| Level | Env Value | Description |
|-------|-----------|-------------|
| Debug | `debug` | All logs are shown (default) |
| Info | `info` | Info, warn, and error only |
| Warn | `warn` | Warn and error only |
| Error | `error` | Error only |

### Annotation API

```go
// Initialize context
ctx = zlog.CtxWithAnnotation(ctx, zlog.DefaultAnnotation())

// Add annotations (thread-safe, can be called from different goroutines)
zlog.AddAnnotation(ctx, map[string]any{"key": "value"})

// Retrieve annotation from context
annotation := zlog.AnnotationFromCtx(ctx)

// Marshal to JSON
data, _ := json.Marshal(annotation)
```

### Full Example

```go
package main

import (
    "context"
    "time"

    "github.com/alhamsya/adk/zlog"
    "github.com/pkg/errors"
)

func main() {
    ctx := context.Background()
    ctx = zlog.CtxWithAnnotation(ctx, zlog.DefaultAnnotation())

    // Inject metadata
    zlog.AddAnnotation(ctx, map[string]any{
        "data": "test",
    })

    zlog.AddAnnotation(ctx, map[string]any{
        "data1": "test1",
    })

    // Info log
    zlog.FromContext(ctx).Info().Msg("success send message")

    // Error log with stack trace
    err := errors.New("something went wrong")
    zlog.FromContext(ctx).Error().Err(err).Msg("failed process")

    // Wait for async log flush
    time.Sleep(100 * time.Millisecond)
}
```

---

## 🏗 Project Architecture

```
adk/
├── go.mod                  # Root module
├── go.work                 # Go workspace definition
├── README.md               # This documentation
│
├── xerr/                   # 🔴 Structured Error Handling
│   ├── go.mod              # Module: github.com/alhamsya/adk/xerr
│   ├── error.go            # Core functions: New, FromError, Wrap, etc.
│   ├── enum.go             # Error type constants & mappings
│   ├── type.go             # Type definitions (Error struct)
│   ├── iface.go            # Interfaces (GRPCCode, GRPCErrType, HTTPStatusCode)
│   ├── options.go          # Functional options (WithMessage, WithType, etc.)
│   ├── stack.go            # Stack trace support
│   └── error_test.go       # Unit tests
│
├── xhttp/                  # 🌐 HTTP Client
│   ├── go.mod              # Module: github.com/alhamsya/adk/xhttp
│   ├── client.go           # HTTP client with retry
│   ├── config.go           # Upstream & endpoint configuration
│   ├── client_test.go      # Unit tests
│   ├── client_config_test.go  # Config integration tests
│   └── config_test.go      # Config parsing tests
│
├── zlog/                   # 📝 Structured Logging
│   ├── go.mod              # Module: github.com/alhamsya/adk/zlog
│   ├── README.md           # Package-level documentation
│   ├── logger.go           # Logger setup, init, constructors
│   ├── annotation.go       # Context annotation system
│   └── types.go            # Type definitions & constants
│
└── cmd/                    # 📂 Example Programs
    ├── xerr/main.go        # xerr usage example
    └── zlog/main.go        # zlog usage example
```

---

## 🔗 Dependencies

| Package | Dependency | Version | Purpose |
|---------|------------|---------|---------|
| `xerr` | `github.com/pkg/errors` | v0.9.1 | Error wrapping & stack traces |
| `xerr` | `google.golang.org/grpc` | v1.79.2 | gRPC status codes |
| `xhttp` | `github.com/avast/retry-go/v5` | v5.0.0 | Retry mechanism |
| `xhttp` | `github.com/pkg/errors` | v0.9.1 | Error wrapping |
| `zlog` | `github.com/rs/zerolog` | v1.34.0 | High-performance logging |

---

## 🧪 Testing

Run unit tests for each package:

```bash
# Test all packages
go test ./xerr/... ./xhttp/... ./zlog/...

# Test with verbose output
go test -v ./xerr/...
go test -v ./xhttp/...

# Test with coverage
go test -cover ./xerr/... ./xhttp/...
```

---

## 📄 License

MIT License — see the [LICENSE](LICENSE) file for details.

---

<p align="center">
  Built with ❤️ by <a href="https://github.com/alhamsya">Alhamsya</a>
</p>
