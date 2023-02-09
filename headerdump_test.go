package headerdump_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jaybubs/headerdump"
)

// There's honestly nothing to test here honestly, it's just a merciless header dump, if you can think of anything do a PR or something

func TestDemo(t *testing.T) {
	cfg := headerdump.CreateConfig()
	cfg.Prefix = "TESTPREFIX"

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := headerdump.New(ctx, next, cfg, "headerdump")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)
}
