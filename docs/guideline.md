# ADK — Usage & Contribution Guidelines

> Conventions for using and extending the Alhamsya Development Kit. Written for humans and AI agents working in this codebase.

## Choosing a module

- Serving gRPC and/or HTTP and want consistent errors → **xerr**.
- Calling JSON REST upstreams with retry → **xhttp**.
- Want structured JSON logs through `slog` → **zlog**.

Import only what you need — each is an independent module.

## Error handling

- **Classify at creation.** Return `xerr.New(Type, msg)` (or a sentinel) instead of `errors.New`. The `Type` is the contract; codes follow.
- **Normalise at boundaries, not everywhere.** Call `xerr.FromError(err, ...)` once, at the transport edge (gRPC interceptor, HTTP handler), not on every internal hop.
- **Match on `Type`, not on value.** Use `xerr.GetType(err) == xerr.TypeNotFound`. Do **not** rely on `errors.Is(err, xerr.ErrNotFound)` unless the message also matches — `Is` compares `Code`+`Message`+`Type`.
- **Don't wrap shared sentinels in place.** `xerr.Wrap(xerr.ErrX, cause)` mutates the global. Use `xerr.NewWithWrap(Type, cause, msg)` or the opts path of `FromError` (returns a copy).
- **Adding a new error type:** add the constant in `enum.go` **and** an entry in both `TypeToGRPCCode` and `TypeToHTTPStatus`, plus a sentinel `Err*`. Missing map entries silently degrade to Internal/500.

## HTTP client

- **Always check `res.StatusCode`.** Non-2xx is not an error and does not trigger retry. `err != nil` means transport/config failure only.
- **Prefer config-driven calls** (`*Endpoint`) over hard-coded URLs. Register upstreams in `Config`; call `SetConfig` before any `*Endpoint` call.
- **Declare retry in config**, not code: `attempts`, `delay` (Go duration string), `type` (`exponential`|`fixed`). Note invalid duration strings are silently ignored.
- **Bodies are JSON only.** Pass a struct to `WithBody`/`Post`/`Put`; read with `res.Unmarshal(&v)`.
- **Set timeouts intentionally.** `SetTimeout` is client-wide; per-upstream `timeout` layers a context deadline on `*Endpoint` calls.

## Logging

- **Import `zlog` for side effects** (`import _ ".../zlog"`) to activate the global slog→zerolog→diode wiring, then log through standard `slog`.
- **Set `ZLOG_LEVEL` before start** (`debug`|`info`|`warn`|`error`). It's read once at init; there is no runtime setter.
- **Use context annotations for request-scoped fields.** `CtxWithAnnotation(ctx, DefaultAnnotation())` → `AddAnnotation(ctx, ...)` → `FromContext(ctx)`. Without `DefaultAnnotation()`, fields are stored but not auto-attached.
- **Accept dropped logs under overload.** The diode is non-blocking by design; don't route audit-critical events through it.

## Code conventions

- **Package naming:** `x`/`z` prefix for kit packages (`xerr`, `xhttp`, `zlog`).
- **Functional options** for anything optional/extensible: `func(*T)` (or an interface wrapping it when the target type must stay hidden, as in `zlog.Option`).
- **Keep implementation types unexported.** Expose interfaces + option constructors; hide concrete structs (`callOptions`, `logger`, annotation internals).
- **Context first.** Any call that does I/O or carries request scope takes `ctx context.Context` as its first argument.

## Testing

- Table-driven tests are the norm (`xerr/error_test.go`, `xhttp/*_test.go`).
- `xhttp` uses `stretchr/testify` for assertions and `httptest` servers; other modules stay dependency-light in tests.
- When adding an error `Type`, add a test asserting both `GRPCStatus()` and `HTTPStatus()` mappings.

## Versioning & release

- Each module versions **independently**. Tag per module path (e.g. `xerr/vX.Y.Z`), not a single repo-wide tag.
- Root module (`github.com/alhamsya/adk`) is the umbrella for docs/examples; the `go.work` file is for local dev and is not consumed by downstream users.
- Keep Go directives aligned (currently `go 1.24` per module).

## Adding a new module to the kit

1. New directory with its own `go.mod` (`github.com/alhamsya/adk/<name>`).
2. Add it to `go.work` `use (...)`.
3. Follow the shared idioms: functional options, context-first, lean public surface, config-driven where operational.
4. Add a `docs/adk-<name>.md` reference and link it from the root `README.md` and `docs/system-design.md`.
