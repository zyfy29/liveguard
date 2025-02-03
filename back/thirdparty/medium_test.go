package thirdparty

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/medium/medium-sdk-go"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"os"
	"testing"
)

func mustReadTextFile(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	return string(bytes)
}

func TestMedium(t *testing.T) {
	m2 := medium.NewClientWithAccessToken(viper.GetString("medium.token"))
	u, err := m2.GetUser("")
	if err != nil {
		log.Fatal(err)
	}

	content := "# Hello, World!\n\nThis is a test post."
	p, err := m2.CreatePost(medium.CreatePostOptions{
		UserID:        u.ID,
		Title:         uuid.NewString(),
		Content:       content,
		ContentFormat: medium.ContentFormatMarkdown,
		PublishStatus: medium.PublishStatusPublic,
	})
	assert.Nil(t, err)

	t.Logf("Post created: %s", p.URL)
}

func TestTranslate(t *testing.T) {
	req := SummarySubmitReq{
		Title:      "Test Title",
		Transcript: "This is a transcript.",
		Summary:    "This is a summary.",
		Tags:       []string{"tag1", "tag2"},
	}

	result := req.Translate()
	t.Logf("Result: \n%s", result)
}

func TestSubmitSummary(t *testing.T) {
	type args struct {
		req SummarySubmitReq
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SubmitSummary(tt.args.req)
			if !tt.wantErr(t, err, fmt.Sprintf("SubmitSummary(%v)", tt.args.req)) {
				return
			}
			t.Log(got)
		})
	}
}
