# adk-xhttp — HTTP Client with Retry

> Reference for the `github.com/alhamsya/adk/xhttp` module. Written for humans and AI agents working in this codebase.

## What it is

A thin wrapper over `net/http.Client` that adds JSON marshalling, a retry mechanism (`avast/retry-go`), and an optional **config-driven upstream registry** so call sites reference `("payments", "createOrder")` instead of hard-coded URLs.

## When to use it

- You call JSON REST services and want body marshal/unmarshal handled.
- You want retry policy (attempts, delay, backoff) declared in config, not code.
- You manage several upstreams with per-upstream host, timeout, and retry.

## Install

```bash
go get github.com/alhamsya/adk/xhttp
```

Deps: `github.com/avast/retry-go/v5`, `github.com/pkg/errors`.

## Client

```go
c := xhttp.NewClient()          // default 30s timeout
c.SetTimeout(10 * time.Second)  // override transport timeout
c.SetConfig(cfg)                // required before *Endpoint calls
```

### Direct calls (explicit URL)

```go
func (c *Client) Get(ctx, url string, opts ...CallOption) (*Res, error)
func (c *Client) Post(ctx, url string, body any, opts ...CallOption) (*Res, error)
func (c *Client) Put(ctx, url string, body any, opts ...CallOption) (*Res, error)
func (c *Client) Delete(ctx, url string, opts ...CallOption) (*Res, error)
func (c *Client) Do(ctx, method, url string, opts ...CallOption) (*Res, error)
```

### Config-driven calls (upstream + endpoint names)

```go
func (c *Client) GetEndpoint(ctx, upstream, endpoint string, opts ...CallOption) (*Res, error)
func (c *Client) PostEndpoint(ctx, upstream, endpoint string, body any, opts ...CallOption) (*Res, error)
func (c *Client) PutEndpoint(ctx, upstream, endpoint string, body any, opts ...CallOption) (*Res, error)
func (c *Client) DeleteEndpoint(ctx, upstream, endpoint string, opts ...CallOption) (*Res, error)
func (c *Client) DoEndpoint(ctx, method, upstream, endpoint string, opts ...CallOption) (*Res, error)
```

URL = `upstream.Host + endpoint.Path`. Retry and timeout come from the upstream config unless overridden per call.

## Call options

```go
WithBody(body any) CallOption   // JSON-marshalled; auto-applied by Post/Put
```
`Post`/`Put`/`*Endpoint` variants inject `WithBody` for you.

## Response

```go
type Res struct {
    Header     http.Header
    StatusCode int
    Body       []byte
}
func (r *Res) Unmarshal(v any) error   // json.Unmarshal(Body, v)
```

## Config

```go
type Config struct {
    Upstream map[string]Upstream `json:"upstream" yaml:"upstream"`
}
type Upstream struct {
    Endpoints   map[string]Endpoint    `json:"endpoints"`
    Extras      map[string]any         `json:"extras"`
    GlobalRetry RetryConfig            `json:"global-retry"`
    Host        string                 `json:"host"`
    Timeout     string                 `json:"timeout"`   // Go duration, e.g. "5s"
    Type        string                 `json:"type"`
}
type Endpoint struct { Path string `json:"path"` }
type RetryConfig struct {
    Attempts int    `json:"attempts"`
    Delay    string `json:"delay"`    // Go duration, e.g. "200ms"
    Timeout  string `json:"timeout"`
    Type     string `json:"type"`     // "exponential" | "fixed"
}
```

### YAML example

```yaml
upstream:
  payments:
    host: https://api.payments.internal
    timeout: 5s
    global-retry:
      attempts: 3
      delay: 200ms
      type: exponential   # or: fixed
    endpoints:
      createOrder: { path: /v1/orders }
```

```go
c := xhttp.NewClient()
c.SetConfig(cfg)
res, err := c.PostEndpoint(ctx, "payments", "createOrder", body)
if err != nil { /* transport/config error */ }
if res.StatusCode >= 300 { /* handle HTTP error yourself */ }
var out Order
_ = res.Unmarshal(&out)
```

## Retry behaviour

- **Default is 1 attempt** (no retry) unless config or a call option raises it.
- Config drives it: `Attempts` > 0 enables retry; `Type` selects `exponential` (backoff) or `fixed` delay; `Delay` parsed as a Go duration.
- Retry uses `retry.Context(ctx)`, so a cancelled context stops retrying.

## Gotchas

- **Non-2xx responses are NOT errors and do NOT trigger retry.** The retried operation only returns an error on transport failure or body-read failure. An HTTP 500 comes back as a normal `*Res` with `StatusCode == 500` and `err == nil`. **You must check `res.StatusCode` yourself.**
- **`*Endpoint` calls fail if `SetConfig` was never called** (`"client config is not set"`), or if the upstream/endpoint name is missing.
- **Bodies are always JSON.** `WithBody` marshals with `encoding/json` and sets `Content-Type: application/json`. No form/multipart support.
- **`SetTimeout` sets the transport-level timeout** (applies to every call on the client). Per-upstream `timeout` adds a `context.WithTimeout` on top for `*Endpoint` calls.
- **Invalid duration strings are silently ignored** (`time.ParseDuration` error → that retry/timeout option is skipped, no error surfaced).
