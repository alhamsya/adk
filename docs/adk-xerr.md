# adk-xerr — Structured Error Handling

> Reference for the `github.com/alhamsya/adk/xerr` module. Written for humans and AI agents working in this codebase.

## What it is

`xerr` is a typed error package that maps one error to **both** a gRPC status code and an HTTP status code. You classify an error once (by `Type`), and the transport layer derives the right code for whichever protocol it speaks.

## When to use it

- You have a service that serves gRPC **and** HTTP (or might later) and want one error taxonomy.
- You want pre-built, reusable error values (sentinels) instead of `errors.New` scattered everywhere.
- You need to normalise arbitrary `error` values into a known type at a boundary (e.g. a gRPC interceptor).

## Install

```bash
go get github.com/alhamsya/adk/xerr
```

Deps: `google.golang.org/grpc` (codes/status), `github.com/pkg/errors` (stack traces).

## Core types

```go
type Type string          // classification, e.g. "NOT_FOUND"

type Error struct {
    Code    codes.Code    // gRPC code, derived from Type
    Message string
    Type    Type
    err     error         // wrapped cause (unexported)
}
```

`Error` implements `error`, `Unwrap`, `Is`, `GRPCStatus() *status.Status`, `HTTPStatus() int`, and `StackTrace() errors.StackTrace`.

## Constructors

| Func | Purpose |
|------|---------|
| `New(t Type, msg string) *Error` | New typed error. Code derived from `t` (falls back to `codes.Internal` for unknown types). |
| `NewWithWrap(t Type, cause error, msg string) *Error` | Same, wrapping a cause. |
| `Wrap(e *Error, cause error) error` | Attach a cause to an existing `*Error`, returns it. |
| `FromError(err error, opts ...Option) *Error` | Normalise any `error` into `*Error` (see below). |

## `FromError` — the boundary function

The one to use where foreign errors enter your typed world (interceptors, handler edges).

- `nil` → `TypeOK` / `codes.OK`.
- Already an `*Error`:
  - no opts → returns it **as-is** (same pointer).
  - with opts → returns a **copy** with opts applied (protects global sentinels from mutation).
- Any other error → builds `TypeSystemError` / `codes.Internal`, then upgrades using adapter interfaces if the cause implements them:
  - `fmt.Stringer` → overrides `Message`.
  - `GRPCCode` → overrides `Code`.
  - `GRPCErrType` → overrides `Type`.

## Options (functional)

```go
WithMessage(msg string)   // set Message
WithType(t Type)          // set Type AND resync Code from TypeToGRPCCode
WithMessageType()         // if Message empty, set it to Type.String()
WithStack()               // wrap cause with pkg/errors stack (skips if already stacked or no cause)
```

## Adapter interfaces

Let foreign error types plug into `FromError` without importing `xerr`:

```go
type GRPCCode     interface{ Code() codes.Code }
type GRPCErrType  interface{ Type() xerr.Type }
type HTTPStatusCode interface{ HTTPStatus() int }
type GRPCErr      interface{ fmt.Stringer; GRPCCode; GRPCErrType } // all three
```

## Transport mapping

```go
e := xerr.New(xerr.TypeNotFound, "user not found")
e.GRPCStatus()  // *status.Status with codes.NotFound
e.HTTPStatus()  // 404  (unknown Type → 500)
xerr.GetType(err) // Type of any error; nil→TypeOK, non-xerr→TypeUnknown
```

### Type → code reference

| Type | gRPC code | HTTP |
|------|-----------|------|
| `TypeOK` | OK | 200 |
| `TypeInvalidParameter` | InvalidArgument | 400 |
| `TypeUnauthorized` | Unauthenticated | 401 |
| `TypeForbidden` / `TypeNoTradingPermission` / `TypeNoSubscription` | PermissionDenied | 403 |
| `TypeNotFound` | NotFound | 404 |
| `TypeDuplicateCall` / `TypeAlreadyExists` | AlreadyExists | 409 |
| `TypeAborted` | Aborted | 409 |
| `TypeServiceBusy` | ResourceExhausted | 429 |
| `TypeRequestCanceled` | Canceled | 499 |
| `TypeSystemError` / `TypeVendorError` | Internal | 500 |
| `TypeUnimplemented` | Unimplemented | 501 |
| `TypeBadGateway` | Unavailable | 502 |
| `TypeMaintenance` | Unavailable | 503 |
| `TypeGatewayTimeout` | DeadlineExceeded | 504 |

Every `Type` has a matching sentinel `Err*` (e.g. `xerr.ErrNotFound`, `xerr.ErrServiceBusy`) whose `Message` equals the type string.

## Patterns

**Return a typed error:**
```go
if user == nil {
    return xerr.New(xerr.TypeNotFound, "user not found")
}
```

**Reuse a sentinel, add a cause:**
```go
return xerr.NewWithWrap(xerr.TypeServiceBusy, cause, "queue full")
// or
return xerr.Wrap(xerr.ErrServiceBusy, cause) // NOTE: mutates the sentinel, see gotcha
```

**Normalise at a boundary and answer the client:**
```go
e := xerr.FromError(err, xerr.WithStack())
http.Error(w, e.Error(), e.HTTPStatus())
```

## Gotchas

- **`Is` compares `Code` + `Message` + `Type`, not identity.** `xerr.New(TypeNotFound, "user not found").Is(xerr.ErrNotFound)` is **false** because the sentinel's message is `"NOT_FOUND"`. Match on `GetType(err) == xerr.TypeNotFound` when you only care about the class.
- **`Wrap(sentinel, cause)` mutates the global sentinel** (it sets the unexported `err` on the shared value). Prefer `NewWithWrap` or `FromError(sentinel, WithStack())` (opts path returns a copy) when wrapping a shared `Err*`.
- **Unknown `Type` → `codes.Internal` / HTTP 500.** `New` with a `Type` not in the map silently downgrades to Internal.
- **`WithType` resyncs `Code`; direct field assignment does not.** Setting `e.Type` by hand leaves `e.Code` stale — use `WithType`.
