# ADK — System Design

> Architecture and design rationale for the Alhamsya Development Kit. Written for humans and AI agents working in this codebase.

## Overview

ADK is a **multi-module Go workspace**: a collection of small, single-purpose backend packages that are meant to be adopted à la carte. Each package is its own Go module with its own `go.mod`, so a consumer can `go get` just `xerr` without pulling in `zlog` or `xhttp`.

| Module | Import path | Responsibility | External deps |
|--------|-------------|----------------|---------------|
| `xerr` | `github.com/alhamsya/adk/xerr` | Typed errors → gRPC + HTTP status | grpc, pkg/errors |
| `xhttp` | `github.com/alhamsya/adk/xhttp` | JSON HTTP client, retry, upstream registry | avast/retry-go, pkg/errors |
| `zlog` | `github.com/alhamsya/adk/zlog` | slog + zerolog over a shared diode, annotations | rs/zerolog |
| root | `github.com/alhamsya/adk` | Umbrella module for docs/examples | — |

Go version: **1.24** per module; the `go.work` workspace pins **1.25** for local development.

## Module independence

The three feature modules are **mutually independent** — none imports another. Verified from `go.mod`:

- `xerr` → grpc, pkg/errors
- `xhttp` → retry-go, pkg/errors (carries a `replace ../zlog` directive, but no `require` — dormant, not an actual dependency)
- `zlog` → zerolog only

```
        ┌────────────────────── go.work (dev, Go 1.25) ──────────────────────┐
        │                                                                     │
   ┌────┴────┐        ┌──────────┐        ┌──────────┐        ┌──────────┐
   │  root   │        │   xerr   │        │  xhttp   │        │   zlog   │
   │  (adk)  │        │ go.mod   │        │ go.mod   │        │ go.mod   │
   └─────────┘        └──────────┘        └──────────┘        └──────────┘
       │                   │                   │                   │
   docs/examples      grpc + errors      retry-go + errors      zerolog
```

> Note: an earlier automated graph pass inferred a `zlog → xerr` call edge. That is **false** — `zlog`'s `newLogger` calls `zerolog.New`, not `xerr.New`. There is no cross-module coupling in code.

## Design principles

1. **One classification, many transports.** `xerr` fixes the idea that an error's meaning is protocol-agnostic; the `Type` is the source of truth and gRPC/HTTP codes are derived from it via lookup tables. Add a protocol → add a table, not a new error taxonomy.

2. **Functional options everywhere.** All three modules use the `func(*T)` option pattern for extensibility without breaking constructors:
   - `xerr.Option` → `WithType`, `WithStack`, ...
   - `xhttp.CallOption` → `WithBody`, ...
   - `zlog.Option` (interface form, to hide the concrete `Annotation`) → `DefaultAnnotation`.

3. **Config-driven over code-driven.** `xhttp` pushes hosts, timeouts, and retry policy into a `Config` map so operational knobs live in YAML/JSON, not recompiled Go. Call sites name `(upstream, endpoint)`.

4. **Non-blocking by default.** `zlog` favours caller latency over guaranteed delivery: the diode writer drops logs under overload rather than blocking the hot path. A deliberate, documented trade-off.

5. **Context as the carrier.** Request-scoped data (log annotations, cancellation, timeouts) flows through `context.Context` rather than globals or thread-locals.

6. **Lean public surface.** Concrete implementation types are unexported (`callOptions`, `logger`, `Annotation` internals). Consumers touch interfaces and option constructors only.

## Cross-cutting themes (parallel, not shared)

- **Stack traces via `pkg/errors`.** Both `xerr` (`WithStack`, `StackTrace()`) and `zlog` (stack marshaler, runtime-frame filtering) integrate the same `pkg/errors` `StackTracer` interface — independently. They interoperate at runtime without depending on each other.
- **Adapter interfaces.** `xerr` accepts foreign errors that implement `GRPCCode` / `GRPCErrType` / `HTTPStatusCode`, so other packages integrate without importing `xerr`.

## How the modules compose in a service

```
 inbound request
      │
      ▼
 ctx = zlog.CtxWithAnnotation(ctx, zlog.DefaultAnnotation())   ── request-scoped fields
      │
      ▼
 res, err := xhttpClient.PostEndpoint(ctx, "payments", "createOrder", body)  ── outbound call
      │
      ▼
 e := xerr.FromError(err, xerr.WithStack())   ── normalise at the boundary
      │
      ├── gRPC handler:  return nil, e.GRPCStatus().Err()
      └── HTTP handler:  w.WriteHeader(e.HTTPStatus())
```

Each module is useful alone; together they form a request → call → error → response spine with consistent logging.
