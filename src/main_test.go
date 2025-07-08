package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func runTestServer() *httptest.Server {
	conf := getConfig("./config.dev.toml")
	handler := setupHandler(conf)
	return httptest.NewServer(handler)
}

func TestSetupHandler(t *testing.T) {
	ts := runTestServer()
	defer ts.Close()

	t.Run("it should return 200 when home is ok", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("%s/home", ts.URL))

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		assert.Equal(t, 200, resp.StatusCode)
	})
}
