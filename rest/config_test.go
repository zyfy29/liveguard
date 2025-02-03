package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func Test_getPocketConfig(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/config/pocket", nil)
	testRequest(t, req)
}

func Test_setPocketConfig(t *testing.T) {
	reqBody := map[string]any{"token": "nihao", "interval": 1000}
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPut, "/config/pocket", bytes.NewBuffer(jsonBody))
	testRequest(t, req)
}
