package action

import (
	"bearguard/cm"
	"bearguard/pocket"
	"bearguard/repo"
	"encoding/json"
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
		// liveまだ生成されていない
		return
	}

	if err = repo.UpdateDBTaskStatus(task.ID, repo.TaskStatusDownloading); err != nil {
		return errors.Wrap(err, "failed to update task status")
	}

	liveAudio, err := cm.DownloadPlaylistAudio(liveInfo.PlayStreamPath)
	if err != nil {
		return errors.Wrap(err, "failed to download live and convert to audio")
	}
	log.Printf("Live audio merged to %s", liveAudio)

	var detail repo.TaskDetail
	if err = json.Unmarshal([]byte(task.Details), &detail); err != nil {
		return errors.Wrap(err, "failed to unmarshal task details")
	}
	detail.FilePath = liveAudio

	if err = repo.UpdateDBTaskStatusAndDetails(task.ID, repo.TaskStatusAwaitTranscript, cm.JsonMarshal(detail)); err != nil {
		return errors.Wrap(err, "failed to update task status and details")
	}
	return
}
