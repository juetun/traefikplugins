package traefikplugins_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/juetun/traefikplugins"
	"github.com/juetun/traefikplugins/logic"
)

// CreateConfig creates the default plugin configuration.
func CreateConfig() *logic.Config {
	return &logic.Config{}
}

func TestDemo(t *testing.T) {
	cfg := CreateConfig()

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := traefikplugins.New(ctx, next, cfg, "demo-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	assertHeader(t, req, "X-Host", "localhost")
	assertHeader(t, req, "X-URL", "http://localhost")
	assertHeader(t, req, "X-Method", "GET")
	assertHeader(t, req, "X-Demo", "test")
}

func assertHeader(t *testing.T, req *http.Request, key, expected string) {
	t.Helper()

	if req.Header.Get(key) != expected {
		t.Errorf("invalid header value: %s expected:%s", req.Header.Get(key), expected)
	}
}
