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

**ADK** (Alhamsya Development Kit) is a modular Go SDK yang menyediakan utilitas penting untuk membangun backend service yang production-ready. Setiap package didesain sebagai Go module independen sehingga dapat diimpor secara terpisah sesuai kebutuhan.

| Package | Deskripsi | Import Path |
|---------|-----------|-------------|
| **[`xerr`](#-xerr--structured-error-handling)** | Structured error handling dengan mapping otomatis ke gRPC & HTTP status codes | `github.com/alhamsya/adk/xerr` |
| **[`xhttp`](#-xhttp--http-client-dengan-retry)** | HTTP client dengan built-in retry, config-driven upstream management | `github.com/alhamsya/adk/xhttp` |
| **[`zlog`](#-zlog--structured-logging)** | Structured logging berbasis zerolog dengan diode buffering & context annotations | `github.com/alhamsya/adk/zlog` |

---

## 🛠 Instalasi

Setiap package adalah Go module independen. Install hanya yang dibutuhkan:

```bash
# Structured Error Handling
go get github.com/alhamsya/adk/xerr

# HTTP Client dengan Retry
go get github.com/alhamsya/adk/xhttp

# Structured Logging
go get github.com/alhamsya/adk/zlog
```

**Persyaratan:** Go 1.24+

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

    // Juga bisa menggunakan slog standard
    slog.Info("hello from slog", "user", "alhamsya")

    time.Sleep(100 * time.Millisecond) // tunggu async log flush
}
```

---

## 🔴 `xerr` — Structured Error Handling

Package `xerr` menyediakan error handling terstruktur yang secara otomatis memetakan error ke **gRPC status codes** dan **HTTP status codes**. Didesain untuk service yang melayani multiple protocol (REST + gRPC) secara bersamaan.

### Instalasi

```bash
go get github.com/alhamsya/adk/xerr
```

### Fitur Utama

- ✅ **Dual-protocol mapping** — satu error type, otomatis ter-mapping ke gRPC & HTTP codes
- ✅ **17 pre-defined error types** — mencakup skenario umum backend service
- ✅ **Pre-built error sentinels** — langsung pakai tanpa inisialisasi ulang
- ✅ **Functional options pattern** — kustomisasi error secara fleksibel
- ✅ **Error wrapping/unwrapping** — kompatibel dengan `errors.Is()` dan `errors.As()`
- ✅ **Stack trace support** — integrasi dengan `pkg/errors`
- ✅ **Pointer-safe** — `FromError` dengan options tidak memodifikasi global preset

### Membuat Error

```go
// Membuat error baru
err := xerr.New(xerr.TypeNotFound, "user not found")

// Membuat error dengan wrapping cause
err := xerr.NewWithWrap(xerr.TypeSystemError, originalErr, "database connection failed")

// Menggunakan pre-built sentinel errors
if errors.Is(err, xerr.ErrUnauthorized) {
    // handle unauthorized
}
```

### Menggunakan Pre-built Error Sentinels

ADK menyediakan error sentinel yang siap pakai:

```go
// Langsung digunakan
return xerr.ErrNotFound           // 404 / NotFound
return xerr.ErrUnauthorized       // 401 / Unauthenticated
return xerr.ErrInvalidParameter   // 400 / InvalidArgument
return xerr.ErrSystemError        // 500 / Internal
return xerr.ErrServiceBusy        // 429 / ResourceExhausted

// Kustomisasi dari sentinel — AMAN, tidak mengubah global
err := xerr.FromError(xerr.ErrNotFound, xerr.WithMessage("product ID 123 not found"))
```

### Konversi ke gRPC & HTTP

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
// Wrap external error
dbErr := sql.ErrNoRows
err := xerr.NewWithWrap(xerr.TypeNotFound, dbErr, "user not found")

// Unwrap
cause := errors.Unwrap(err) // sql.ErrNoRows

// Comparison
errors.Is(err, anotherXerr) // compares Code, Message, and Type
```

### Functional Options

```go
// Konversi dari generic error dengan options
err := xerr.FromError(someError,
    xerr.WithMessage("custom message"),      // override message
    xerr.WithType(xerr.TypeNotFound),        // override type (juga update gRPC code)
    xerr.WithStack(),                        // attach stack trace
    xerr.WithMessageType(),                  // gunakan nama type sebagai message jika kosong
)
```

### Error Type Mappings

Tabel lengkap mapping antara error `Type`, HTTP Status Code, dan gRPC Code:

| Type | Nilai | HTTP | gRPC Code |
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

`xerr` menyediakan interfaces untuk integrasi kustom:

```go
// Implementasikan interface ini pada error kustom Anda
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

### Contoh: Middleware gRPC

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

### Contoh: Middleware HTTP

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

## 🌐 `xhttp` — HTTP Client dengan Retry

Package `xhttp` menyediakan HTTP client dengan fitur retry built-in dan manajemen upstream berbasis konfigurasi. Cocok untuk service-to-service communication yang memerlukan resiliensi.

### Instalasi

```bash
go get github.com/alhamsya/adk/xhttp
```

### Fitur Utama

- ✅ **Built-in retry** — retry otomatis dengan `retry-go` (exponential & fixed backoff)
- ✅ **Config-driven upstreams** — kelola multiple upstream service via konfigurasi JSON/YAML
- ✅ **JSON auto-serialization** — request body otomatis di-marshal ke JSON
- ✅ **Response unmarshaling** — unmarshal response body ke struct dengan mudah
- ✅ **Per-upstream timeout** — konfigurasi timeout per upstream service
- ✅ **Extras metadata** — simpan metadata tambahan per upstream (API keys, partner IDs, dll.)

### Basic Usage

```go
client := xhttp.NewClient()
client.SetTimeout(10 * time.Second)

// GET request
res, err := client.Get(ctx, "https://api.example.com/users/123")
if err != nil {
    log.Fatal(err)
}

fmt.Println(res.StatusCode)  // 200
fmt.Println(string(res.Body)) // response body

// Unmarshal response
var user User
err = res.Unmarshal(&user)
```

### POST / PUT / DELETE

```go
// POST dengan body — otomatis di-serialize ke JSON
body := map[string]string{"name": "Alhamsya", "role": "engineer"}
res, err := client.Post(ctx, "https://api.example.com/users", body)

// PUT
res, err := client.Put(ctx, "https://api.example.com/users/123", updatedUser)

// DELETE
res, err := client.Delete(ctx, "https://api.example.com/users/123")
```

### Config-Driven Upstream Management

Cara yang lebih powerful untuk mengelola multiple upstream services:

#### Konfigurasi (JSON)

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

#### Konfigurasi (YAML)

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

#### Penggunaan

```go
// Inisialisasi client dengan config
client := xhttp.NewClient()
client.SetConfig(cfg) // cfg diload dari JSON/YAML

// Request menggunakan nama upstream & endpoint
res, err := client.GetEndpoint(ctx, "payment-service", "get-status")
res, err := client.PostEndpoint(ctx, "payment-service", "create-payment", paymentReq)
res, err := client.PutEndpoint(ctx, "notification-service", "send-email", emailReq)
res, err := client.DeleteEndpoint(ctx, "payment-service", "cancel-payment")
```

### Retry Strategy

| Strategi | Konstanta | Deskripsi |
|----------|-----------|-----------|
| **Exponential Backoff** | `xhttp.RetryTypeExponential` | Delay meningkat secara eksponensial antar retry |
| **Fixed Delay** | `xhttp.RetryTypeFixed` | Delay tetap antar retry |

### Struct Reference

```go
// Konfigurasi utama
type Config struct {
    Upstream map[string]Upstream
}

type Upstream struct {
    Host        string                 // Base URL upstream
    Timeout     string                 // Timeout per request (e.g., "5s")
    Type        string                 // Tipe service (e.g., "rest", "grpc")
    GlobalRetry RetryConfig            // Konfigurasi global retry
    Endpoints   map[string]Endpoint    // Daftar endpoint
    Extras      map[string]interface{} // Metadata tambahan
}

type Endpoint struct {
    Path string // Path endpoint (e.g., "/api/v1/users")
}

type RetryConfig struct {
    Attempts int    // Jumlah percobaan
    Delay    string // Delay antar retry (e.g., "100ms")
    Timeout  string // Timeout per retry attempt
    Type     string // "exponential" atau "fixed"
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

Package `zlog` menyediakan structured logging berkinerja tinggi yang dibangun di atas `zerolog` dengan kompatibilitas `log/slog`. Dilengkapi dengan diode buffering untuk non-blocking I/O dan context-scoped annotations.

### Instalasi

```bash
go get github.com/alhamsya/adk/zlog
```

### Fitur Utama

- ✅ **High-performance** — diode-buffered writer untuk async, non-blocking logging
- ✅ **Dual API** — gunakan `zerolog` langsung ATAU `log/slog` standard library
- ✅ **Context annotations** — inject metadata ke context, otomatis muncul di setiap log
- ✅ **RFC3339 timestamps** — format timestamp standar (field: `timestamp`)
- ✅ **Stack trace support** — otomatis log stack trace untuk errors dari `pkg/errors`
- ✅ **Environment-based level** — set level via env var `ZLOG_LEVEL`
- ✅ **Thread-safe annotations** — menggunakan `sync.Map` untuk concurrent access

### Basic Usage (Side-Effect Import)

Cukup import `zlog` untuk mengkonfigurasi `log/slog` secara otomatis:

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

Inject metadata ke context yang otomatis muncul di setiap log entry:

```go
// 1. Inisialisasi context dengan annotation
ctx := context.Background()
ctx = zlog.CtxWithAnnotation(ctx, zlog.DefaultAnnotation())

// 2. Tambahkan annotation
zlog.AddAnnotation(ctx, map[string]any{
    "user_id":    "12345",
    "request_id": "req-abc-789",
    "service":    "payment",
})

// 3. Log — annotation otomatis ter-attach
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

Error dari `pkg/errors` otomatis menyertakan stack trace:

```go
import "github.com/pkg/errors"

err := errors.New("database connection timeout")
logger.Error().Err(err).Msg("Failed to query database")
// Output: includes "stack" field with filtered stack frames
```

### Context Logger dengan Hooks

```go
func handler(ctx context.Context) {
    // Buat logger baru dengan hooks
    ctx, logger := zlog.NewContext(ctx, myCustomHook)

    // Update context fields
    logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
        return c.Str("handler", "payment")
    })

    // Ambil logger dari context
    log := zlog.FromContext(ctx)
    log.Info().Msg("Handler executed")
}
```

### Konfigurasi Log Level

Set log level via environment variable:

```bash
# Options: debug, info, warn, error
export ZLOG_LEVEL=info
```

| Level | Env Value | Deskripsi |
|-------|-----------|-----------|
| Debug | `debug` | Semua log ditampilkan (default) |
| Info | `info` | Info, warn, dan error |
| Warn | `warn` | Warn dan error saja |
| Error | `error` | Error saja |

### Annotation API

```go
// Inisialisasi context
ctx = zlog.CtxWithAnnotation(ctx, zlog.DefaultAnnotation())

// Tambah annotation (thread-safe, bisa dipanggil dari goroutine berbeda)
zlog.AddAnnotation(ctx, map[string]any{"key": "value"})

// Ambil annotation dari context
annotation := zlog.AnnotationFromCtx(ctx)

// Marshal ke JSON
data, _ := json.Marshal(annotation)
```

### Contoh Lengkap

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

    // Error log dengan stack trace
    err := errors.New("something went wrong")
    zlog.FromContext(ctx).Error().Err(err).Msg("failed process")

    // Tunggu async log flush
    time.Sleep(100 * time.Millisecond)
}
```

---

## 🏗 Arsitektur Project

```
adk/
├── go.mod                  # Root module
├── go.work                 # Go workspace definition
├── README.md               # Dokumentasi ini
│
├── xerr/                   # 🔴 Structured Error Handling
│   ├── go.mod              # Module: github.com/alhamsya/adk/xerr
│   ├── error.go            # Fungsi utama: New, FromError, Wrap, dll.
│   ├── enum.go             # Error type constants & mappings
│   ├── type.go             # Type definitions (Error struct)
│   ├── iface.go            # Interfaces (GRPCCode, GRPCErrType, HTTPStatusCode)
│   ├── options.go          # Functional options (WithMessage, WithType, dll.)
│   ├── stack.go            # Stack trace support
│   └── error_test.go       # Unit tests
│
├── xhttp/                  # 🌐 HTTP Client
│   ├── go.mod              # Module: github.com/alhamsya/adk/xhttp
│   ├── client.go           # HTTP client dengan retry
│   ├── config.go           # Konfigurasi upstream & endpoint
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

| Package | Dependency | Versi | Kegunaan |
|---------|------------|-------|----------|
| `xerr` | `github.com/pkg/errors` | v0.9.1 | Error wrapping & stack traces |
| `xerr` | `google.golang.org/grpc` | v1.79.2 | gRPC status codes |
| `xhttp` | `github.com/avast/retry-go/v5` | v5.0.0 | Retry mechanism |
| `xhttp` | `github.com/pkg/errors` | v0.9.1 | Error wrapping |
| `zlog` | `github.com/rs/zerolog` | v1.34.0 | High-performance logging |

---

## 🧪 Testing

Jalankan unit test untuk masing-masing package:

```bash
# Test semua package
go test ./xerr/... ./xhttp/... ./zlog/...

# Test dengan verbose output
go test -v ./xerr/...
go test -v ./xhttp/...

# Test dengan coverage
go test -cover ./xerr/... ./xhttp/...
```

---

## 📄 Lisensi

MIT License — lihat file [LICENSE](LICENSE) untuk detail.

---

<p align="center">
  Built with ❤️ by <a href="https://github.com/alhamsya">Alhamsya</a>
</p>
