package traefik_xff_to_xrealip

import (
	"context"
	"net/http"
	"strings"
)

// Config holds the middleware configuration.
type Config struct {
	// Depth is the index of the IP address to select from the X-Forwarded-For header.
	// Defaults to 0 (the first IP).
	Depth int `json:"depth,omitempty" yaml:"depth,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Depth: 0,
	}
}

// Middleware is the XFF to XRealIP middleware.
type Middleware struct {
	next  http.Handler
	depth int
}

// New creates a new middleware instance.
func New(_ context.Context, next http.Handler, config *Config, _ string) (http.Handler, error) {
	return &Middleware{
		next:  next,
		depth: config.Depth,
	}, nil
}

// ServeHTTP handles the HTTP request.
func (m *Middleware) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	xff := req.Header.Get("X-Forwarded-For")
	if xff != "" {
		ipCandidates := strings.Split(xff, ",")
		var ips []string
		for _, ipStr := range ipCandidates {
			trimmedIP := strings.TrimSpace(ipStr)
			if trimmedIP != "" {
				ips = append(ips, trimmedIP)
			}
		}

		if len(ips) > 0 {
			// Ensure the configured depth is within the bounds of the available IPs.
			if m.depth >= 0 && m.depth < len(ips) {
				req.Header.Set("X-Real-Ip", ips[m.depth])
			}
		}
	}
	m.next.ServeHTTP(rw, req)
}
