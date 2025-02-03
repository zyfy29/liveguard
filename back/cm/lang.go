package cm

import (
	"github.com/teamlint/opencc"
	"log"
)

// Convert2Zhcn 繁体中文转简体中文
func Convert2Zhcn(text string) string {
	cvt, err := opencc.New("t2s")
	if err != nil {
		log.Printf("Failed to create opencc: %v", err)
		return text
	}
	res, err := cvt.Convert(text)
	if err != nil {
		log.Printf("Failed to convert text: %v", err)
		return text
	}
	return res
}
