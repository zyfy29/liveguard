package cm

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_mergeMP4Files(t *testing.T) {
	type args struct {
		files []string
		dir   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MergeMP4Files(tt.args.files)
			if !tt.wantErr(t, err, fmt.Sprintf("MergeMP4Files(%v)", tt.args.files)) {
				return
			}
			t.Logf("Output file: %s", got)
		})
	}
}

func TestDownloadLiveAndConvertToAudio(t *testing.T) {
	type args struct {
		playListUrl string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Test DownloadPlaylistAudio",
			args: args{
				playListUrl: "https://cychengyuan-vod.48.cn/89653525/20240812/cy/1029882491506069504-ki5ixiimun5hkq47dw77.m3u8",
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DownloadPlaylistAudio(tt.args.playListUrl)
			if !tt.wantErr(t, err, fmt.Sprintf("DownloadPlaylistAudio(%v)", tt.args.playListUrl)) {
				return
			}
			t.Log(got)
		})
	}
}

func TestGetPlaylistDuration(t *testing.T) {
	type args struct {
		playListUrl string
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Test GetPlaylistDuration",
			args: args{
				playListUrl: "https://cychengyuan-vod.48.cn/6744/20240818/cy/1031720783931314176-w3vcjr1bmbd9d6uvdf1b.m3u8",
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPlaylistDuration(tt.args.playListUrl)
			if !tt.wantErr(t, err, fmt.Sprintf("GetPlaylistDuration(%v)", tt.args.playListUrl)) {
				return
			}
			t.Log(got)
		})
	}
}

func Test_parsePlayList(t *testing.T) {
	type args struct {
		playListUrl string
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "normal",
			args:    args{playListUrl: "https://cychengyuan-vod.48.cn/6744/20240818/cy/1031720783931314176-w3vcjr1bmbd9d6uvdf1b.m3u8"},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePlaylist(tt.args.playListUrl)
			if !tt.wantErr(t, err, fmt.Sprintf("parsePlaylist(%v)", tt.args.playListUrl)) {
				return
			}
			t.Log(got)
		})
	}
}
