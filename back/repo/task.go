package repo

import (
	"bearguard/cm"
	"bearguard/pocket"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

const (
	TaskStatusPending         = "pending"
	TaskStatusDownloading     = "downloading"
	TaskStatusAwaitTranscript = "await_transcript"
	TaskStatusTranscribing    = "transcribing"
	TaskStatusSummarizing     = "summarizing"
	TaskStatusAwaitSubmit     = "await_submit"
	TaskStatusSucceed         = "succeed"
	TaskStatusFailed          = "failed"
)

type Task struct {
	ID        int64        `db:"id" json:"id"`
	Title     string       `db:"title" json:"title"`
	Status    string       `db:"status" json:"status"`
	LiveID    string       `db:"live_id" json:"live_id"`
	LiveTime  sql.NullTime `db:"live_time" json:"live_time"`
	OwnerName string       `db:"owner_name" json:"owner_name"`
	ErrorInfo string       `db:"error_info" json:"error_info"`
	Details   string       `db:"details" json:"-"`

	Created time.Time    `db:"created" json:"created"`
	Updated sql.NullTime `db:"updated" json:"updated"`
}

type TaskDetail struct {
	OwnerID      string `json:"owner_id"`
	FilePath     string `json:"file_path"`
	TranscriptID string `json:"transcript_id"`
	Transcript   string `json:"transcript"`
	Summary      string `json:"summary"`
	PostURL      string `json:"post_url"`
	Duration     int    `json:"duration"` // seconds
}

func GetDBTasks() ([]Task, error) {
	tasks := []Task{}
	query := `SELECT * FROM task ORDER BY id desc`
	conn := getDbMust()
	defer conn.Close()
	err := conn.Select(&tasks, query)
	return tasks, err
}

func GetDBTaskByID(id int64) (Task, error) {
	var task Task
	query := `SELECT * FROM task WHERE id = $1`
	conn := getDbMust()
	defer conn.Close()
	err := conn.Get(&task, query, id)
	return task, err
}

func GetDBTaskByLiveID(liveID string) (Task, error) {
	var task Task
	query := `SELECT * FROM task WHERE live_id = $1`
	conn := getDbMust()
	defer conn.Close()
	err := conn.Get(&task, query, liveID)
	return task, err
}

func GetDBTasksByStatus(status string, limit int) ([]Task, error) {
	tasks := []Task{}
	query := `SELECT * FROM task WHERE status = $1 ORDER BY id desc`
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}
	conn := getDbMust()
	defer conn.Close()
	err := conn.Select(&tasks, query, status)
	return tasks, err
}

func InsertDBTask(task Task) error {
	query := `INSERT INTO task (title, status, live_id, live_time, owner_name, error_info, details) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT DO NOTHING`
	conn := getDbMust()
	defer conn.Close()
	if _, err := conn.Exec(query, task.Title, task.Status, task.LiveID, task.LiveTime, task.OwnerName, task.ErrorInfo, task.Details); err != nil {
		return fmt.Errorf("failed to insert task: %w", err)
	}
	return nil
}

func DeleteDBTask(id int64) error {
	deleteTaskSQL := `DELETE FROM task WHERE id = $1`
	conn := getDbMust()
	defer conn.Close()
	if _, err := conn.Exec(deleteTaskSQL, id); err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	return nil
}

func UpdateDBTaskStatus(id int64, status string) error {
	updateTaskSQL := `UPDATE task SET status = $1 WHERE id = $2`
	conn := getDbMust()
	defer conn.Close()
	if _, err := conn.Exec(updateTaskSQL, status, id); err != nil {
		return fmt.Errorf("failed to update task status: %w", err)
	}
	return nil
}

func UpdateDBTaskErrorInfo(id int64, errorInfo string) error {
	updateTaskSQL := `UPDATE task SET error_info = $1 WHERE id = $2`
	conn := getDbMust()
	defer conn.Close()
	if _, err := conn.Exec(updateTaskSQL, errorInfo, id); err != nil {
		return fmt.Errorf("failed to update task error info: %w", err)
	}
	return nil
}

func UpdateDBTaskDetails(id int64, details string) error {
	updateTaskSQL := `UPDATE task SET details = $1 WHERE id = $2`
	conn := getDbMust()
	defer conn.Close()
	if _, err := conn.Exec(updateTaskSQL, details, id); err != nil {
		return fmt.Errorf("failed to update task details: %w", err)
	}
	return nil
}

func UpdateDBTaskStatusAndDetails(id int64, status string, details string) error {
	updateTaskSQL := `UPDATE task SET status = $1, details = $2 WHERE id = $3`
	conn := getDbMust()
	defer conn.Close()
	if _, err := conn.Exec(updateTaskSQL, status, details, id); err != nil {
		return fmt.Errorf("failed to update task status and details: %w", err)
	}
	return nil
}

func UpdateDBTaskStatusAndErrorInfo(id int64, status string, errorInfo string) error {
	updateTaskSQL := `UPDATE task SET status = $1, error_info = $2 WHERE id = $3`
	conn := getDbMust()
	defer conn.Close()
	if _, err := conn.Exec(updateTaskSQL, status, errorInfo, id); err != nil {
		return fmt.Errorf("failed to update task status and error info: %w", err)
	}
	return nil
}

func CreateLiveTask(id string) (err error) {
	live, err := pocket.GetClient().GetLiveInfo(id)
	if err != nil {
		return err
	}
	duration, err := cm.GetPlaylistDuration(live.PlayStreamPath)
	if err != nil {
		return err
	}

	ctime, _ := strconv.ParseInt(live.Ctime, 10, 64)
	task := Task{
		Title:     live.Title,
		Status:    TaskStatusPending,
		LiveID:    id,
		LiveTime:  sql.NullTime{Time: time.UnixMilli(ctime), Valid: true},
		OwnerName: live.User.UserName,
		Details:   cm.JsonMarshal(TaskDetail{OwnerID: live.User.UserId, Duration: duration}),
	}
	if err = InsertDBTask(task); err != nil {
		return errors.Wrap(err, "failed to insert task")
	}
	return
}

func RestoreTask(id int64) error {
	task, err := GetDBTaskByID(id)
	if err != nil {
		return errors.Wrapf(err, "Failed to get task by id: %d", id)
	}

	// downloading -> pending
	if task.Status == TaskStatusDownloading {
		if err := UpdateDBTaskStatus(task.ID, TaskStatusPending); err != nil {
			return errors.Wrap(err, "Failed to update live status")
		}
	}

	// failed, check detail
	var detail TaskDetail
	if err = json.Unmarshal([]byte(task.Details), &detail); err != nil {
		return errors.Wrap(err, "Failed to unmarshal task details")
	}
	if detail.Transcript != "" && detail.Summary != "" {
		// 投稿失敗の場合は自動的に再試行するので、スキップ
		return nil
	} else if detail.TranscriptID != "" || detail.FilePath != "" {
		// ダウンロード成功の場合、transcriptを試行
		if err := UpdateDBTaskStatus(task.ID, TaskStatusAwaitTranscript); err != nil {
			return errors.Wrap(err, "Failed to update live status")
		}
	}
	return nil
}

func RestoreFailedTasks() error {
	tasks, err := GetDBTasksByStatus(TaskStatusFailed, 0)
	if err != nil {
		return errors.Wrap(err, "Failed to get tasks")
	}

	for _, task := range tasks {
		var detail TaskDetail
		if err = json.Unmarshal([]byte(task.Details), &detail); err != nil {
			return errors.Wrap(err, "Failed to unmarshal task details")
		}

		// 投稿失敗の場合は自動的に再試行するので、スキップ
		if detail.Transcript != "" && detail.Summary != "" {
			continue
		} else if detail.TranscriptID != "" || detail.FilePath != "" {
			// ダウンロード成功の場合
			if err := UpdateDBTaskStatus(task.ID, TaskStatusAwaitTranscript); err != nil {
				return errors.Wrap(err, "Failed to update live status")
			}
		}
		//else if err := repo.UpdateDBTaskStatus(task.ID, repo.TaskStatusPending); err != nil {
		//	// ダウンロードさえ失敗の場合
		//	log.Printf("Failed to do download live: %v", err)
		//}
	}
	return nil
}
