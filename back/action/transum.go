package action

import (
	"bearguard/cm"
	"bearguard/repo"
	"bearguard/thirdparty"
	"encoding/json"
	"github.com/AssemblyAI/assemblyai-go-sdk"
	"github.com/pkg/errors"
	"log"
	"time"
)

func WatchTaskAndCallLLM() {
	for {
		tasks, err := repo.GetDBTasksByStatus(repo.TaskStatusAwaitTranscript, 10)
		if err != nil {
			log.Printf("Failed to get tasks: %v", err)
		}

		for _, task := range tasks {
			if err := doTranscriptAndSummarize(task); err != nil {
				log.Printf("Failed to do transcript and summarize: %v", err)
			}
		}
		time.Sleep(10 * time.Second)
	}
}

func doTranscriptAndSummarize(task repo.Task) (err error) {
	defer func() {
		if err != nil {
			_ = repo.UpdateDBTaskStatusAndErrorInfo(task.ID, repo.TaskStatusFailed, err.Error())
		}
	}()
	if err = repo.UpdateDBTaskStatus(task.ID, repo.TaskStatusTranscribing); err != nil {
		return errors.Wrap(err, "failed to update task status")
	}
	var detail repo.TaskDetail
	if err = json.Unmarshal([]byte(task.Details), &detail); err != nil {
		return errors.Wrap(err, "failed to unmarshal task details")
	}
	if detail.TranscriptID == "" {
		var trans assemblyai.Transcript
		trans, err = thirdparty.GetTranscriptFromFile(detail.FilePath)
		if err != nil {
			return errors.Wrap(err, "failed to get transcript from file")
		}
		detail.TranscriptID = *trans.ID
		detail.Transcript = *trans.Text
		if err = repo.UpdateDBTaskStatusAndDetails(task.ID, repo.TaskStatusSummarizing, cm.JsonMarshal(detail)); err != nil {
			return errors.Wrap(err, "failed to update task status and details")
		}
	} else if detail.Transcript == "" {
		detail.Transcript, err = thirdparty.GetTranscriptFromID(detail.TranscriptID)
		if err != nil {
			return errors.Wrap(err, "failed to get transcript from id")
		}
		if err = repo.UpdateDBTaskStatusAndDetails(task.ID, repo.TaskStatusSummarizing, cm.JsonMarshal(detail)); err != nil {
			return errors.Wrap(err, "failed to update task status and details")
		}
	}

	summary, err := thirdparty.GetSummaryFromTranscript(detail.TranscriptID)
	if err != nil {
		return errors.Wrap(err, "failed to get summary from transcript")
	}
	detail.Summary = summary

	// changed from repo.TaskStatusAwaitSubmit
	if err = repo.UpdateDBTaskStatusAndDetails(task.ID, repo.TaskStatusSucceed, cm.JsonMarshal(detail)); err != nil {
		return errors.Wrap(err, "failed to update task status and details")
	}
	return
}
