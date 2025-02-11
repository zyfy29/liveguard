package action

import (
	"bearguard/cm"
	"bearguard/pocket"
	"bearguard/repo"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"strings"
	"time"
)

func WatchTaskAndDownloadLive() {
	for {
		tasks, err := repo.GetDBTasksByStatus(repo.TaskStatusPending, 10)
		if err != nil {
			log.Printf("Failed to get tasks: %v", err)
		}

		for _, task := range tasks {
			if err := doDownloadLive(task); err != nil {
				log.Printf("Failed to do download live: %v", err)
			}
		}
		time.Sleep(10 * time.Second)
	}
}

func doDownloadLive(task repo.Task) (err error) {
	defer func() {
		if err != nil {
			_ = repo.UpdateDBTaskStatusAndErrorInfo(task.ID, repo.TaskStatusFailed, err.Error())
		}
	}()
	liveInfo, err := pocket.GetClient().GetLiveInfo(task.LiveID)
	if err != nil {
		return errors.Wrap(err, "failed to get live info")
	}
	if liveInfo.LiveId == "" || strings.HasPrefix(liveInfo.PlayStreamPath, "rtmp") {
		return
	}

	if err = repo.UpdateDBTaskStatus(task.ID, repo.TaskStatusDownloading); err != nil {
		return errors.Wrap(err, "failed to update task status")
	}

	// 非同期でダウンロード開始
	progressChan, resultChan := cm.DownloadPlaylistAudio(liveInfo.PlayStreamPath)
	totalSegments := -1

	// 進捗監視
	for progressChan != nil || resultChan != nil {
		select {
		case progress := <-progressChan:
			if totalSegments == -1 {
				totalSegments = progress
			} else {
				_ = repo.UpdateDBTaskErrorInfo(task.ID, fmt.Sprintf("%d/%d", progress, totalSegments))
			}

		case result := <-resultChan:
			if result.Err != nil {
				return errors.Wrap(result.Err, "failed to download live and convert to audio")
			}

			log.Printf("Live audio merged to %s", result.FilePath)

			var detail repo.TaskDetail
			if err := json.Unmarshal([]byte(task.Details), &detail); err != nil {
				return errors.Wrap(err, "failed to unmarshal task details")
			}
			detail.FilePath = result.FilePath

			_ = repo.UpdateDBTaskErrorInfo(task.ID, "")
			if err := repo.UpdateDBTaskStatusAndDetails(task.ID, repo.TaskStatusAwaitTranscript, cm.JsonMarshal(detail)); err != nil {
				return errors.Wrap(err, "failed to update task status and details")
			}
			return nil
		}
	}

	return errors.New("unexpected termination of download process")
}
