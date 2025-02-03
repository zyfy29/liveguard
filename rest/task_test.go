package rest

import (
	"net/http"
	"testing"
)

func TestGetTasks(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/task/", nil)
	testRequest(t, req)
}

func TestGetTaskById(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/task/6", nil)
	testRequest(t, req)
}
