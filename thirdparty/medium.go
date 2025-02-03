package thirdparty

import (
	"fmt"
	"github.com/medium/medium-sdk-go"
	"github.com/spf13/viper"
)

type SummarySubmitReq struct {
	Title      string   `json:"title"`
	Transcript string   `json:"transcript"`
	Summary    string   `json:"summary"`
	Tags       []string `json:"tags"`
}

func (req *SummarySubmitReq) Translate() string {
	return fmt.Sprintf("%s\n\n## 原始转录文本\n\n%s", req.Summary, req.Transcript)
}

func SubmitSummary(req SummarySubmitReq) (string, error) {
	client := medium.NewClientWithAccessToken(viper.GetString("medium.token"))
	u, err := client.GetUser("")
	if err != nil {
		return "", err
	}

	p, err := client.CreatePost(medium.CreatePostOptions{
		UserID:  u.ID,
		Title:   req.Title,
		Content: req.Translate(),
		Tags:    req.Tags,

		ContentFormat: medium.ContentFormatMarkdown,
		PublishStatus: medium.PublishStatusUnlisted,
	})
	if err != nil {
		return "", err
	}
	return p.URL, nil
}
