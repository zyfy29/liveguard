package pocket

import (
	"bearguard/cm"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Client struct {
	mu           sync.Mutex
	lastCallTime time.Time
	rateLimit    time.Duration

	token   string
	appInfo string
}

func newClient(token string, appInfo map[string]string, interval int) *Client {
	appInfo["deviceId"] = strings.ToUpper(uuid.NewString())
	c := Client{
		mu:           sync.Mutex{},
		lastCallTime: time.Time{},
		rateLimit:    time.Duration(interval) * time.Millisecond,
		token:        token,
		appInfo:      cm.JsonMarshal(appInfo),
	}
	return &c
}

func (c *Client) doPocketRequest(req *http.Request) (*http.Response, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	elapsed := now.Sub(c.lastCallTime)
	if elapsed < c.rateLimit {
		waitTime := c.rateLimit - elapsed
		time.Sleep(waitTime)
	}
	c.lastCallTime = time.Now()

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("token", c.token)
	req.Header.Set("appInfo", c.appInfo)
	client := &http.Client{}
	resp, err := client.Do(req)
	return resp, err
}

type formatter[T any] struct {
	c *Client
}

// req should be able to marshaled to json
func (f formatter[T]) doRestWithResult(method, url string, req any) (t T, err error) {
	if method != "GET" {
		method = "POST"
	}
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return
	}

	r, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return
	}

	httpResp, err := f.c.doPocketRequest(r)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()

	// レスポンスボディを読み込む
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return
	}

	resp := &Resp[T]{}
	if err = json.Unmarshal(body, resp); err != nil {
		return
	}
	if !resp.Success {
		err = fmt.Errorf("request failed with code %d, message: %s", resp.Status, resp.Message)
		return
	}
	return resp.Content, nil
}
