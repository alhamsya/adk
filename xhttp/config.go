package xhttp

// Config holds the configuration for the xhttp package.
type Config struct {
	Upstream map[string]Upstream `json:"upstream" yaml:"upstream"`
}

// Upstream holds the configuration for a specific upstream service.
type Upstream struct {
	Endpoints   map[string]Endpoint    `json:"endpoints" yaml:"endpoints"`
	Extras      map[string]interface{} `json:"extras" yaml:"extras"`
	GlobalRetry RetryConfig            `json:"global-retry" yaml:"global-retry"`
	Host        string                 `json:"host" yaml:"host"`
	Timeout     string                 `json:"timeout" yaml:"timeout"`
	Type        string                 `json:"type" yaml:"type"`
}

// Endpoint holds the configuration for a specific endpoint.
type Endpoint struct {
	Path string `json:"path" yaml:"path"`
}

// RetryConfig holds the configuration for retry logic.
type RetryConfig struct {
	Attempts int    `json:"attempts" yaml:"attempts"`
	Delay    string `json:"delay" yaml:"delay"`
	Timeout  string `json:"timeout" yaml:"timeout"`
	Type     string `json:"type" yaml:"type"`
}
