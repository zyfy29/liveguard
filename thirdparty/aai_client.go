package thirdparty

import (
	"bearguard/cm"
	"context"
	aai "github.com/AssemblyAI/assemblyai-go-sdk"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

var prompt = `
- 你是一位专业的内容总结专家。
- 你需要对以下偶像直播中的发言内容进行详细总结。
- 请按照以下要求进行总结：
  1. **详细性**: 请详细总结直播发言中提到的每一个话题，尽可能多地呈现细节，不要遗漏任何内容。总结时不需要缩短文字，尽量丰富内容。
  2. **结构化**: 使用Markdown格式组织总结内容。使用二级标题"##"标识每个话题，按话题出现的先后顺序分点作答，确保每个时间段的内容都有清晰的总结。
  3. **列表和重点**: 对于每个话题的关键点，使用有序列表"1", "2"标注，并使用**加粗**标记重要内容。必要时，可以使用引用块">"突出显示特别的发言或情感表达。
  4. **情感分析**: 对发言者在每个时间段的情感进行分析，注重情感的变化和表达，并用专门的小节进行说明。
  5. **准确性**: 总结时仅基于提供的转录内容，不添加任何额外的信息。使用准确的中文表达，确保语句通顺，并保持专业性。
  6. **格式示例**:
    
## 话题一：XXXX
1. **关键点1**: 详细描述……
2. **关键点2**: 详细描述……

## 话题二：XXXX
1. **关键点1**: 详细描述……
2. **关键点2**: 详细描述……

> 情感分析：在此部分，发言者显得XXXX……

- 请开始总结，并生成尽可能长的文本。
`

func GetTranscriptFromFile(inputFile string) (aai.Transcript, error) {
	client := aai.NewClient(viper.GetString("aai.token"))
	start := time.Now()
	f, _ := os.Open(inputFile)
	res, err := client.Transcripts.TranscribeFromReader(context.TODO(), f, &aai.TranscriptOptionalParams{LanguageCode: "zh"})
	if err != nil {
		log.Printf("TranscribeFromReader Error: %v\n", err)
		return aai.Transcript{}, err
	}
	*res.Text = cm.Convert2Zhcn(*res.Text)
	end := time.Now()
	log.Printf("TranscribeFromReader took %s\n", end.Sub(start))
	return res, err
}

func GetTranscriptFromID(transcriptID string) (string, error) {
	client := aai.NewClient(viper.GetString("aai.token"))
	res, err := client.Transcripts.Get(context.TODO(), transcriptID)
	if err != nil {
		log.Printf("Transcripts.Get Error: %v\n", err)
		return "", err
	}
	return *res.Text, err
}

func GetSummaryFromTranscript(transcriptID string) (string, error) {
	client := aai.NewClient(viper.GetString("aai.token"))
	var params aai.LeMURTaskParams
	params.Prompt = aai.String(prompt)
	params.TranscriptIDs = []string{transcriptID}
	params.FinalModel = "anthropic/claude-3-5-sonnet"
	params.MaxOutputSize = aai.Int64(4000)

	start := time.Now()
	result, err := client.LeMUR.Task(context.TODO(), params)
	end := time.Now()
	log.Printf("LeMUR Task took %s\n", end.Sub(start))

	if err != nil {
		log.Printf("LeMUR Error: %v\n", err)
		return "", err
	}
	log.Printf("LeMUR Response: %s\n", *result.Response)
	return *result.Response, err
}

// TranscriptAndSummarize experimental
func TranscriptAndSummarize(inputFile string) (transcription, summary string, err error) {
	ctx := context.Background()

	client := aai.NewClient(viper.GetString("aai.token"))
	opts := aai.TranscriptOptionalParams{
		LanguageCode: "zh",
		//SpeakerLabels: aai.Bool(true),
	}

	start := time.Now()
	f, _ := os.Open(inputFile)
	transcript, err := client.Transcripts.TranscribeFromReader(ctx, f, &opts)
	if err != nil {
		log.Printf("TranscribeFromReader Error: %v\n", err)
		return "", "", err
	}
	log.Printf("Transcript ID: %s\n", *transcript.ID)
	log.Printf("Transcript Text: %s\n", *transcript.Text)
	end := time.Now()
	log.Printf("TranscribeFromURL took %s\n", end.Sub(start))

	var params aai.LeMURTaskParams
	params.Prompt = aai.String(prompt)
	params.TranscriptIDs = []string{aai.ToString(transcript.ID)}
	params.FinalModel = "anthropic/claude-3-5-sonnet"
	params.MaxOutputSize = aai.Int64(4000)

	start = time.Now()
	result, err := client.LeMUR.Task(ctx, params)
	end = time.Now()
	log.Printf("LeMUR Task took %s\n", end.Sub(start))

	if err != nil {
		log.Printf("LeMUR Error: %v\n", err)
	} else {
		log.Printf("LeMUR Response: %s\n", *result.Response)
	}
	return *transcript.Text, *result.Response, nil
}

// TranscriptFromUrl experimental
func TranscriptFromUrl(url string) (transcription string, err error) {
	//ctx := context.Background()

	client := aai.NewClient(viper.GetString("aai.token"))
	opts := aai.TranscriptOptionalParams{
		LanguageCode: "zh",
	}

	start := time.Now()
	transcript, err := client.Transcripts.TranscribeFromURL(context.TODO(), url, &opts)
	if err != nil {
		log.Printf("TranscribeFromURL Error: %v\n", err)
		return "", err
	}
	log.Printf("Transcript ID: %s\n", *transcript.ID)
	log.Printf("Transcript Text: %s\n", *transcript.Text)
	end := time.Now()
	log.Printf("TranscribeFromURL took %s\n", end.Sub(start))
	return *transcript.Text, nil
}
