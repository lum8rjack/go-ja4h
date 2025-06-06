package ja4h

import (
	"net/http"
	"strings"
	"testing"
)

func TestJA4H_Fingerprint(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		headers  map[string]string
		body     string
		expected string
	}{
		{
			name:     "GET no headers",
			method:   "GET",
			headers:  map[string]string{},
			body:     "",
			expected: "ge11nn000000_e3b0c44298fc_000000000000_000000000000",
		},
		{
			name:     "GET Accept-Language",
			method:   "GET",
			headers:  map[string]string{"Accept-Language": "en-us"},
			body:     "",
			expected: "ge11nn01enus_6ec18c3c2e22_000000000000_000000000000",
		},
		{
			name:     "POST Accept-Language Cookie",
			method:   "POST",
			headers:  map[string]string{"Accept-Language": "en-us", "Cookie": "admin=true"},
			body:     "",
			expected: "po11cn01enus_fab33050c907_8c6976e5b541_8c6976e5b541",
		},
		{
			name:     "POST Accept Accept-Language Cookie Referer",
			method:   "POST",
			headers:  map[string]string{"Accept": "application/json", "Accept-Language": "en", "Cookie": "admin=true", "Referer": "https://example.com"},
			body:     "",
			expected: "po11cr02en00_e5fdf4927470_8c6976e5b541_8c6976e5b541",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, "http://example.com", strings.NewReader(tt.body))
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}
			fp := JA4H(req)
			if fp != tt.expected {
				t.Errorf("ja4h mismatch: got %q, want %q", fp, tt.expected)
			}
		})
	}
}
