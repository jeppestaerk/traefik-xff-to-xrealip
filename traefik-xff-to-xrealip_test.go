package traefik_xff_to_xrealip

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestXRealIpRewrite(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Forwarded-For", "203.0.113.5, 10.0.0.1")

	rec := httptest.NewRecorder()
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.Header.Get("X-Real-Ip")
		if ip != "203.0.113.5" {
			t.Errorf("expected X-Real-Ip = 203.0.113.5, got %s", ip)
		}
	})
	mw, _ := New(context.Background(), CreateConfig(), next)
	mw.ServeHTTP(rec, req)
}
