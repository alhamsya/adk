package xhttp

import (
	"encoding/json"
	"testing"
)

func TestConfigUnmarshal(t *testing.T) {
	jsonConfig := `
{
  "upstream": {
    "mandiri": {
      "endpoints": {
        "status-rdn": {
          "path": "/openapi/onboarding/v1.0/inquiryStatusRDN"
        }
      },
      "extras": {
        "branch-code": "99999",
        "participant-id": "CC001",
        "participant-name": "PT. Mandiri Sekuritas",
        "partner-id": "SANDBOX",
        "product-type": "MTBINV-OL"
      },
      "global-retry": {
        "attempts": 3,
        "delay": "1s",
        "timeout": "10s",
        "type": "exponential"
      },
      "host": "https://sandbox.bankmandiri.co.id",
      "timeout": "5s",
      "type": "rest"
    }
  }
}`

	var cfg Config
	err := json.Unmarshal([]byte(jsonConfig), &cfg)
	if err != nil {
		t.Fatalf("failed to unmarshal config: %v", err)
	}

	upstream, ok := cfg.Upstream["mandiri"]
	if !ok {
		t.Fatal("expected upstream 'mandiri' to exist")
	}

	if upstream.Host != "https://sandbox.bankmandiri.co.id" {
		t.Errorf("expected host 'https://sandbox.bankmandiri.co.id', got '%s'", upstream.Host)
	}

	if upstream.Timeout != "5s" {
		t.Errorf("expected timeout '5s', got '%s'", upstream.Timeout)
	}

	endpoint, ok := upstream.Endpoints["status-rdn"]
	if !ok {
		t.Fatal("expected endpoint 'status-rdn' to exist")
	}

	if endpoint.Path != "/openapi/onboarding/v1.0/inquiryStatusRDN" {
		t.Errorf("expected path '/openapi/onboarding/v1.0/inquiryStatusRDN', got '%s'", endpoint.Path)
	}

	if upstream.GlobalRetry.Attempts != 3 {
		t.Errorf("expected retry attempts 3, got %d", upstream.GlobalRetry.Attempts)
	}

	if val, ok := upstream.Extras["branch-code"]; !ok || val != "99999" {
		t.Errorf("expected extras branch-code '99999', got '%v'", val)
	}
}
