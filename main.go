package traefik_xff_to_xrealip

import (
	"context"
	"net/http"
	"strings"
)

type Config struct{}

func CreateConfig() *Config {
	return &Config{}
}

type Middleware struct {
	next http.Handler
}

func New(_ context.Context, _ *Config, next http.Handler, _ string) (http.Handler, error) {
	return &Middleware{next: next}, nil
}

func (m *Middleware) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	xff := req.Header.Get("X-Forwarded-For")
	if xff != "" {
		ip := strings.TrimSpace(strings.Split(xff, ",")[0])
		req.Header.Set("X-Real-Ip", ip)
	}
	m.next.ServeHTTP(rw, req)
}
