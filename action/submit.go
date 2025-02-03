package action

import (
	"bearguard/cm"
	"bearguard/repo"
	"bearguard/thirdparty"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"strings"
	"time"
)

func WatchAndSubmitToMedium() {
	for {
		tasks, err := repo.GetDBTaskByStatus(repo.TaskStatusAwaitSubmit, 10)
		if err != nil {
			log.Printf("Failed to get tasks: %v", err)
		}

		for _, task := range tasks {
			if err := doSubmit(task); err != nil {
				log.Printf("Failed to do submit: %v", err)
			}
		}
		time.Sleep(10 * time.Second)
	}

}

func doSubmit(task repo.Task) (err error) {
	defer func() {
		if err != nil {
			_ = repo.UpdateDBTaskStatusAndErrorInfo(task.ID, repo.TaskStatusFailed, err.Error())
		}
	}()
	var detail repo.TaskDetail
	if err = json.Unmarshal([]byte(task.Details), &detail); err != nil {
		return errors.Wrap(err, "failed to unmarshal task details")
	}
	req := thirdparty.SummarySubmitReq{
		Title:      fmt.Sprintf("%s-%s-\"%s\"", task.OwnerName, task.LiveTime.Time.Format("2006-01-02"), task.Title),
		Transcript: detail.Transcript,
		Summary:    detail.Summary,
		Tags:       []string{"live", "summary", task.OwnerName},
	}
	if task.LiveTime.Valid {
		req.Tags = append(req.Tags, task.LiveTime.Time.Format("200601"))
	}
	postUrl, err := thirdparty.SubmitSummary(req)
	if err != nil {
		return errors.Wrap(err, "failed to submit summary")
	}
	detail.PostURL = postUrl
	if err = repo.UpdateDBTaskStatusAndDetails(task.ID, repo.TaskStatusSucceed, cm.JsonMarshal(detail)); err != nil {
		return errors.Wrap(err, "failed to update task status and details")
	}
	if err = repo.UpdateDBTaskErrorInfo(task.ID, ""); err != nil {
		return errors.Wrap(err, "failed to update task post url")

	}
	return
}

func WatchAndRetrySubmit() {
	for {
		tasks, err := repo.GetDBTaskByStatus(repo.TaskStatusFailed, 10)
		if err != nil {
			log.Printf("Failed to get tasks: %v", err)
		}

		for _, task := range tasks {
			if err := retrySubmitFailedTask(task); err != nil {
				log.Printf("Failed to retry failed task: %v", err)
			}
		}
		time.Sleep(3600 * time.Second)
	}
}

func retrySubmitFailedTask(task repo.Task) (err error) {
	var detail repo.TaskDetail
	if err = json.Unmarshal([]byte(task.Details), &detail); err != nil {
		return errors.Wrap(err, "failed to unmarshal task details")
	}
	if detail.Summary != "" && strings.Contains(task.ErrorInfo, "failed to submit summary") {
		return doSubmit(task)
	}
	return
}
