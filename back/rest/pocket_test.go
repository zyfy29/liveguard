package rest

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

func Test_getMyName(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/pocket/me", nil)
	testRequest(t, req)
}

func Test_getMembers(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/pocket/member", nil)
	testRequest(t, req)
}

func Test_getLivesByMember(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/pocket/live?owner_id=63566&next_time=1020093092891267072", nil)
	testRequest(t, req)
}

func Test_getLiveDuration(t *testing.T) {
	playlistUrl := "https://cychengyuan-vod.48.cn/6744/20240818/cy/1031720783931314176-w3vcjr1bmbd9d6uvdf1b.m3u8"
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/pocket/live/duration?playlist_url=%s", url.QueryEscape(playlistUrl)), nil)
	testRequest(t, req)
}
