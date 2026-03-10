package xhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/avast/retry-go/v5"
	"github.com/pkg/errors"
)

const (
	// ApplicationJSON content type
	ApplicationJSON = "application/json"

	// RetryTypeExponential indicates exponential backoff retry strategy
	RetryTypeExponential = "exponential"
	// RetryTypeFixed indicates fixed delay retry strategy
	RetryTypeFixed = "fixed"
)

// CallOption configures the call options.
type CallOption func(*callOptions)

type retryConfig struct {
	retryOpts []retry.Option
}

type callOptions struct {
	retryCfg retryConfig
	reqBody  interface{}
}

// WithBody sets the request body.
func WithBody(body interface{}) CallOption {
	return func(c *callOptions) {
		c.reqBody = body
	}
}

// Res represents an HTTP response.
type Res struct {
	Header     http.Header
	StatusCode int
	Body       []byte
}

// Unmarshal parses the JSON-encoded data and stores the result in the value pointed to by v.
func (r *Res) Unmarshal(v interface{}) error {
	if err := json.Unmarshal(r.Body, v); err != nil {
		return errors.Wrap(err, "failed to unmarshal response body")
	}
	return nil
}

// Client wraps http.Client to provide retry mechanism.
type Client struct {
	cli    *http.Client
	config *Config
}

// NewClient creates a new Client with default configuration.
func NewClient() *Client {
	return &Client{
		cli: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SetTimeout sets the timeout for the underlying http client.
func (c *Client) SetTimeout(timeout time.Duration) {
	c.cli.Timeout = timeout
}

// SetConfig sets the configuration for the client.
func (c *Client) SetConfig(cfg Config) {
	c.config = &cfg
}

// Get performs a GET request.
func (c *Client) Get(ctx context.Context, url string, opts ...CallOption) (*Res, error) {
	return c.Do(ctx, http.MethodGet, url, opts...)
}

// Post performs a POST request.
func (c *Client) Post(ctx context.Context, url string, body interface{}, opts ...CallOption) (*Res, error) {
	opts = append(opts, WithBody(body))
	return c.Do(ctx, http.MethodPost, url, opts...)
}

// Put performs a PUT request.
func (c *Client) Put(ctx context.Context, url string, body interface{}, opts ...CallOption) (*Res, error) {
	opts = append(opts, WithBody(body))
	return c.Do(ctx, http.MethodPut, url, opts...)
}

// Delete performs a DELETE request.
func (c *Client) Delete(ctx context.Context, url string, opts ...CallOption) (*Res, error) {
	return c.Do(ctx, http.MethodDelete, url, opts...)
}

// Do performs an HTTP request with the given method, url, and options.
func (c *Client) Do(ctx context.Context, method, url string, opts ...CallOption) (*Res, error) {
	// Default call options
	cOpts := &callOptions{}
	for _, opt := range opts {
		opt(cOpts)
	}

	var reqBody io.Reader
	if cOpts.reqBody != nil {
		reqBytes, err := json.Marshal(cOpts.reqBody)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal request body")
		}
		reqBody = bytes.NewBuffer(reqBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new request")
	}

	if cOpts.reqBody != nil {
		req.Header.Set("Content-Type", ApplicationJSON)
	}

	var resp *http.Response
	var resBody []byte

	// Operation to be retried
	op := func() error {
		// Reset body if needed for retry
		if req.Body != nil && req.GetBody != nil {
			bodyReadCloser, err := req.GetBody()
			if err != nil {
				return errors.Wrap(err, "failed to get request body")
			}
			req.Body = bodyReadCloser
		}

		var err error
		resp, err = c.cli.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		resBody, err = io.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "failed to read response body")
		}

		return nil
	}

	// Configure retry options
	// Default to 1 attempt unless overridden
	finalRetryOpts := []retry.Option{
		retry.Attempts(1),
		retry.Context(ctx),
	}
	finalRetryOpts = append(finalRetryOpts, cOpts.retryCfg.retryOpts...)

	// Execute with retry
	err = retry.New(
		finalRetryOpts...,
	).Do(op)

	if err != nil {
		return nil, errors.Wrap(err, "http request failed")
	}

	return &Res{
		Header:     resp.Header,
		StatusCode: resp.StatusCode,
		Body:       resBody,
	}, nil
}

// GetEndpoint performs a GET request using configuration.
func (c *Client) GetEndpoint(ctx context.Context, upstream, endpoint string, opts ...CallOption) (*Res, error) {
	return c.DoEndpoint(ctx, http.MethodGet, upstream, endpoint, opts...)
}

// PostEndpoint performs a POST request using configuration.
func (c *Client) PostEndpoint(ctx context.Context, upstream, endpoint string, body interface{}, opts ...CallOption) (*Res, error) {
	opts = append(opts, WithBody(body))
	return c.DoEndpoint(ctx, http.MethodPost, upstream, endpoint, opts...)
}

// PutEndpoint performs a PUT request using configuration.
func (c *Client) PutEndpoint(ctx context.Context, upstream, endpoint string, body interface{}, opts ...CallOption) (*Res, error) {
	opts = append(opts, WithBody(body))
	return c.DoEndpoint(ctx, http.MethodPut, upstream, endpoint, opts...)
}

// DeleteEndpoint performs a DELETE request using configuration.
func (c *Client) DeleteEndpoint(ctx context.Context, upstream, endpoint string, opts ...CallOption) (*Res, error) {
	return c.DoEndpoint(ctx, http.MethodDelete, upstream, endpoint, opts...)
}

// DoEndpoint performs an HTTP request using configuration.
func (c *Client) DoEndpoint(ctx context.Context, method, upstreamName, endpointName string, opts ...CallOption) (*Res, error) {
	if c.config == nil {
		return nil, errors.New("client config is not set")
	}

	upstream, ok := c.config.Upstream[upstreamName]
	if !ok {
		return nil, errors.Errorf("upstream '%s' not found in config", upstreamName)
	}

	endpoint, ok := upstream.Endpoints[endpointName]
	if !ok {
		return nil, errors.Errorf("endpoint '%s' not found in upstream '%s'", endpointName, upstreamName)
	}

	url := upstream.Host + endpoint.Path

	// Configure retry options from config
	if upstream.GlobalRetry.Attempts > 0 {
		retryOpts := []retry.Option{
			retry.Attempts(uint(upstream.GlobalRetry.Attempts)),
		}

		if upstream.GlobalRetry.Delay != "" {
			delay, err := time.ParseDuration(upstream.GlobalRetry.Delay)
			if err == nil {
				retryOpts = append(retryOpts, retry.Delay(delay))
			}
		}

		switch upstream.GlobalRetry.Type {
		case RetryTypeExponential:
			retryOpts = append(retryOpts, retry.DelayType(retry.BackOffDelay))
		case RetryTypeFixed:
			retryOpts = append(retryOpts, retry.DelayType(retry.FixedDelay))
		}

		// Add retry options to call options
		opts = append(opts, func(c *callOptions) {
			c.retryCfg.retryOpts = append(c.retryCfg.retryOpts, retryOpts...)
		})
	}

	// Handle timeout if configured
	if upstream.Timeout != "" {
		timeout, err := time.ParseDuration(upstream.Timeout)
		if err == nil {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}
	}

	return c.Do(ctx, method, url, opts...)
}
