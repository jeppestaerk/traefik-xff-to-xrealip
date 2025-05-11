package traefik_xff_to_xrealip

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestXRealIpRewrite(t *testing.T) {
	tests := []struct {
		name            string
		xffHeader       string
		configDepth     int
		expectedRealIp  string
		expectHeaderSet bool // True if X-Real-Ip should be set, false otherwise
	}{
		{
			name:            "default depth (0), multiple IPs in X-Forwarded-For",
			xffHeader:       "203.0.113.5, 10.0.0.1, 192.168.1.1",
			configDepth:     0,
			expectedRealIp:  "203.0.113.5",
			expectHeaderSet: true,
		},
		{
			name:            "depth 1, multiple IPs in X-Forwarded-For",
			xffHeader:       "203.0.113.5, 10.0.0.1, 192.168.1.1",
			configDepth:     1,
			expectedRealIp:  "10.0.0.1",
			expectHeaderSet: true,
		},
		{
			name:            "depth 2, multiple IPs in X-Forwarded-For",
			xffHeader:       "203.0.113.5, 10.0.0.1, 192.168.1.1",
			configDepth:     2,
			expectedRealIp:  "192.168.1.1",
			expectHeaderSet: true,
		},
		{
			name:            "default depth (0), single IP in X-Forwarded-For",
			xffHeader:       "203.0.113.5",
			configDepth:     0,
			expectedRealIp:  "203.0.113.5",
			expectHeaderSet: true,
		},
		{
			name:            "default depth (0), X-Forwarded-For with spaces",
			xffHeader:       "  203.0.113.5  , 10.0.0.1",
			configDepth:     0,
			expectedRealIp:  "203.0.113.5",
			expectHeaderSet: true,
		},
		{
			name:            "no X-Forwarded-For header",
			xffHeader:       "", // Simulate not setting the header
			configDepth:     0,
			expectedRealIp:  "",
			expectHeaderSet: false,
		},
		{
			name:            "empty X-Forwarded-For header (just spaces)",
			xffHeader:       " ",
			configDepth:     0,
			expectedRealIp:  "",
			expectHeaderSet: false, // After filtering, ips list is empty
		},
		{
			name:            "X-Forwarded-For starts with comma, default depth (0)",
			xffHeader:       ",203.0.113.5",
			configDepth:     0,
			expectedRealIp:  "203.0.113.5", // Empty string before comma is filtered out
			expectHeaderSet: true,
		},
		{
			name:            "X-Forwarded-For has multiple commas, default depth (0)",
			xffHeader:       "1.1.1.1,,2.2.2.2",
			configDepth:     0,
			expectedRealIp:  "1.1.1.1",
			expectHeaderSet: true,
		},
		{
			name:            "X-Forwarded-For has multiple commas, depth 1",
			xffHeader:       "1.1.1.1,,2.2.2.2",
			configDepth:     1,
			expectedRealIp:  "2.2.2.2",
			expectHeaderSet: true,
		},
		{
			name:            "depth out of bounds (too high)",
			xffHeader:       "203.0.113.5, 10.0.0.1",
			configDepth:     2,
			expectedRealIp:  "",
			expectHeaderSet: false,
		},
		{
			name:            "depth out of bounds (negative)",
			xffHeader:       "203.0.113.5, 10.0.0.1",
			configDepth:     -1,
			expectedRealIp:  "",
			expectHeaderSet: false,
		},
		{
			name:            "depth 0, XFF contains only empty strings after split and trim",
			xffHeader:       ", ,, ,",
			configDepth:     0,
			expectedRealIp:  "",
			expectHeaderSet: false,
		},
		{
			name:            "depth 0, XFF is completely empty string",
			xffHeader:       "", // This is covered by "no X-Forwarded-For header" if header not set
			configDepth:     0,  // or by setting X-Forwarded-For: "" explicitly
			expectedRealIp:  "",
			expectHeaderSet: false,
		},
		{
			name:            "X-Forwarded-For set to empty string explicitly",
			xffHeader:       "XFF_EMPTY_STRING_SENTINEL", // Special value to indicate setting X-Forwarded-For: ""
			configDepth:     0,
			expectedRealIp:  "",
			expectHeaderSet: false, // xff != "" is true, but ips list will be empty
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.xffHeader == "XFF_EMPTY_STRING_SENTINEL" {
				req.Header.Set("X-Forwarded-For", "")
			} else if tt.xffHeader != "" || tt.name == "empty X-Forwarded-For header (just spaces)" {
				// Set header if xffHeader is not empty, or for the specific test case "empty X-Forwarded-For header (just spaces)"
				req.Header.Set("X-Forwarded-For", tt.xffHeader)
			}
			// For "no X-Forwarded-For header", the header is intentionally not set.

			rec := httptest.NewRecorder()
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ip := r.Header.Get("X-Real-Ip")
				if tt.expectHeaderSet {
					if ip != tt.expectedRealIp {
						t.Errorf("expected X-Real-Ip = %q, got %q", tt.expectedRealIp, ip)
					}
				} else {
					// Check if header exists at all
					if _, ok := r.Header["X-Real-Ip"]; ok {
						t.Errorf("X-Real-Ip header should not be set, but was found with value %q", ip)
					}
				}
			})

			cfg := CreateConfig()
			cfg.Depth = tt.configDepth // Set the depth from the test case
			mw, err := New(context.Background(), next, cfg, "test-xff-to-realip")
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}

			mw.ServeHTTP(rec, req)
		})
	}
}
