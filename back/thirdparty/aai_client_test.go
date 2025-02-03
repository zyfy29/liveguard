package thirdparty

import (
	"fmt"
	aai "github.com/AssemblyAI/assemblyai-go-sdk"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"os"
	"testing"
)

func TestTranscriptAndSummarize(t *testing.T) {
	type args struct {
		inputFile string
	}
	tests := []struct {
		name      string
		args      args
		tSavePath string
		sSavePath string
		wantErr   bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTranscription, gotSummary, err := TranscriptAndSummarize(tt.args.inputFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("TranscriptAndSummarize() error = %v, wantErr %v", err, tt.wantErr)
			}
			_ = os.WriteFile(tt.tSavePath, []byte(gotTranscription), os.ModePerm)
			_ = os.WriteFile(tt.sSavePath, []byte(gotSummary), os.ModePerm)
		})
	}
}

func TestTranscriptFromUrl(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test TranscriptFromUrl",
			args: args{
				url: "https://nim-nosdn.netease.im/NDA5MzEwOA==/bmltYV83MjUxODc3MzgzXzE3MjQxNzM0ODM4MjZfNmU4ZTg3ZjgtOTI4NC00OTY1LWFiZDgtODE5ZDM0ZjM0ZGQ2",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTranscription, err := TranscriptFromUrl(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("TranscriptFromUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("Transcription: %s", gotTranscription)
		})
	}
}

func TestGetTranscriptFromFile(t *testing.T) {
	type args struct {
		inputFile string
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTranscriptFromFile(tt.args.inputFile)
			if !tt.wantErr(t, err, fmt.Sprintf("GetTranscriptFromFile(%v)", tt.args.inputFile)) {
				return
			}
			t.Logf("Transcription: %s", *got.Text)
		})
	}
}

func TestGetSummaryFromTranscript(t *testing.T) {
	type args struct {
		transcriptID string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Test GetSummaryFromTranscript",
			args: args{
				transcriptID: "1dc3f2ee-3cd3-45cd-ac15-c282aa1666bc",
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSummaryFromTranscript(tt.args.transcriptID)
			if !tt.wantErr(t, err, fmt.Sprintf("GetSummaryFromTranscript(%v)", tt.args.transcriptID)) {
				return
			}
			t.Logf("Summary: %s", got)
		})
	}
}

func TestNihao2(t *testing.T) {
	client := aai.NewClient(viper.GetString("aai.token"))
	tr, err := client.Transcripts.Get(context.TODO(), "b401906a-dd9b-4864-8fb7-2207823cc395")
	assert.Nil(t, err)
	t.Log(*tr.Text)
}

func TestGetTranscriptFromID(t *testing.T) {
	type args struct {
		transcriptID string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Test GetTranscriptFromID",
			args: args{
				transcriptID: "b401906a-dd9b-4864-8fb7-2207823cc395",
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTranscriptFromID(tt.args.transcriptID)
			if !tt.wantErr(t, err, fmt.Sprintf("GetTranscriptFromID(%v)", tt.args.transcriptID)) {
				return
			}
			t.Logf("Transcription: %s", got)
		})
	}
}
