package rest

import (
	"net/http"
	"testing"
)

func TestHTTPDeamon(t *testing.T) {
	Startup()
}

func TestPing(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	testRequest(t, req)
}
