package repo

import (
	"bearguard/cm"
	"database/sql"
	"fmt"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCRUDTask(t *testing.T) {
	// Define a task for testing
	task := Task{
		Title:     "Test Task",
		Status:    "Test Pending",
		LiveID:    uuid.NewString(),
		OwnerName: "Test Owner",
		ErrorInfo: "No Error",
		Details:   "Test Details",
	}

	// Insert the task
	err := InsertDBTask(task)
	assert.NoError(t, err, "InsertDBTask failed")

	// Retrieve the task by live_id
	retrievedTask, err := GetDBTaskByLiveID(task.LiveID)
	assert.NoError(t, err, "GetDBTaskByLiveID failed")
	assert.NotNil(t, retrievedTask, "Retrieved task is nil")
	assert.Equal(t, task.Title, retrievedTask.Title, "Task title mismatch")

	// Update the task status
	newStatus := "Completed"
	err = UpdateDBTaskStatusAndDetails(retrievedTask.ID, newStatus, "")
	assert.NoError(t, err, "UpdateDBTaskStatusAndDetails failed")

	// Verify the status update
	updatedTask, err := GetDBTaskByLiveID(task.LiveID)
	assert.NoError(t, err, "GetDBTaskByLiveID after update failed")
	assert.Equal(t, newStatus, updatedTask.Status, "Task status mismatch after update")

	// Delete the task
	err = DeleteDBTask(updatedTask.ID)
	assert.NoError(t, err, "DeleteDBTask failed")

	// Verify the task deletion
	_, err = GetDBTaskByLiveID(task.LiveID)
	assert.Error(t, err, "Expected error after deleting task")
	assert.True(t, errors.Is(err, sql.ErrNoRows), "Expected sql.ErrNoRows after deleting task")
}

func TestGetDBTaskByStatus(t *testing.T) {
	type args struct {
		status string
		limit  int
	}
	tests := []struct {
		name    string
		args    args
		want    []Task
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Get tasks by status",
			args: args{
				status: "Test Pending",
				limit:  10,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = InsertDBTask(Task{
				Status: "Test Pending",
				LiveID: uuid.NewString(),
			})
			got, err := GetDBTaskByStatus(tt.args.status, tt.args.limit)
			t.Logf("Got %d tasks with status %s", len(got), tt.args.status)
			if !tt.wantErr(t, err, fmt.Sprintf("GetDBTaskByStatus(%v, %v)", tt.args.status, tt.args.limit)) {
				return
			}
		})
	}
}

func TestInsertDBTask(t *testing.T) {
	type args struct {
		task Task
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Insert task",
			args: args{
				task: Task{
					Title:  "Test Task",
					Status: "Test Pending",
					LiveID: "test_live_id",
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, InsertDBTask(tt.args.task), fmt.Sprintf("InsertDBTask(%v)", tt.args.task))
		})
	}
}

func TestUpdateDBTaskStatusAndDetails(t *testing.T) {
	type args struct {
		id      int64
		status  string
		details string
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Update task status",
			args: args{
				id:     20,
				status: "test_status",
			},
			wantErr: assert.NoError,
		},
		{
			name: "Update task status and details",
			args: args{
				id:      20,
				status:  "test_status",
				details: "test_details",
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, UpdateDBTaskStatusAndDetails(tt.args.id, tt.args.status, tt.args.details), fmt.Sprintf("UpdateDBTaskStatusAndDetails(%v, %v, %v)", tt.args.id, tt.args.status, tt.args.details))
		})
	}
}

func TestUpdateDBTaskStatusAndErrorInfo(t *testing.T) {
	type args struct {
		id        int64
		status    string
		errorInfo string
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Update task status and error info",
			args: args{
				id:        20,
				status:    "test_status",
				errorInfo: "test_error_info",
			},
		},
		{
			name: "Update task status",
			args: args{
				id:     20,
				status: "test_status",
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, UpdateDBTaskStatusAndErrorInfo(tt.args.id, tt.args.status, tt.args.errorInfo), fmt.Sprintf("UpdateDBTaskStatusAndErrorInfo(%v, %v, %v)", tt.args.id, tt.args.status, tt.args.errorInfo))
		})
	}
}

func TestUpdateDBTaskDetails(t *testing.T) {
	type args struct {
		id      int64
		details string
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Update task details",
			args: args{
				id:      20,
				details: `{"key": "value"}`,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, UpdateDBTaskDetails(tt.args.id, tt.args.details), fmt.Sprintf("UpdateDBTaskDetails(%v, %v)", tt.args.id, tt.args.details))
		})
	}
}

func TestUpdateDBTaskErrorInfo(t *testing.T) {
	type args struct {
		id        int64
		errorInfo string
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Update task error info",
			args: args{
				id:        20,
				errorInfo: "",
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, UpdateDBTaskErrorInfo(tt.args.id, tt.args.errorInfo), fmt.Sprintf("UpdateDBTaskErrorInfo(%v, %v)", tt.args.id, tt.args.errorInfo))
		})
	}
}

func TestGetDBTasks(t *testing.T) {
	tests := []struct {
		name    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "Get tasks",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetDBTasks()
			if !tt.wantErr(t, err, fmt.Sprintf("GetDBTasks()")) {
				return
			}
			t.Logf("Got %d tasks", len(got))
			for i, v := range got {
				t.Log(i, v)
			}
		})
	}
}

func TestCreateLiveTask(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "normal",
			args: args{
				id: "1042148386357972992",
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, CreateLiveTask(tt.args.id), fmt.Sprintf("CreateLiveTask(%v)", tt.args.id))
		})
	}
}

func Test_restoreFailedTasks(t *testing.T) {
	p := gomonkey.NewPatches()
	defer p.Reset()
	tasks := []Task{
		// 投稿失敗のタスク
		{
			ID: 1,
			Details: cm.JsonMarshal(
				TaskDetail{
					Transcript: "123",
					Summary:    "456",
				},
			),
		},
		{
			ID: 2,
			Details: cm.JsonMarshal(
				TaskDetail{
					TranscriptID: "789",
				},
			),
		},
		{
			ID: 3,
			Details: cm.JsonMarshal(
				TaskDetail{
					FilePath: "abc",
				},
			),
		},
		{
			ID:      101,
			Details: `{"owner_id":"538697","file_path":"/root/bearguard/data/live-c3ad8474-5be4-4939-ad16-24c8925ca653.mp3","transcript_id":"","transcript":"","summary":"","post_url":"https://medium.com/@wasuremono127/%E6%B2%88%E5%B0%8F%E7%88%B1-2024-08-31-%E6%9D%A5%E5%95%A6-75f259365975"}`,
		},
	}

	p.ApplyFunc(GetDBTaskByID, func(id int64) (Task, error) {
		for _, task := range tasks {
			if task.ID == id {
				return task, nil
			}
		}
		panic("no such task")
	})
	p.ApplyFunc(UpdateDBTaskStatus, func(id int64, status string) error {
		switch int(id) {
		case 1:
			assert.Equal(t, TaskStatusFailed, status)
		case 2:
			assert.Equal(t, TaskStatusAwaitTranscript, status)
		case 3:
			assert.Equal(t, TaskStatusAwaitTranscript, status)
		case 101:
			assert.Equal(t, TaskStatusAwaitTranscript, status)
		}
		return nil
	})

	for _, task := range tasks {
		err := RestoreTask(task.ID)
		assert.NoError(t, err)
	}
}

func TestGetDBTaskByID(t *testing.T) {
	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "normal",
			args:    args{id: 6},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetDBTaskByID(tt.args.id)
			if !tt.wantErr(t, err, fmt.Sprintf("GetDBTaskByID(%v)", tt.args.id)) {
				return
			}
			t.Log(cm.JsonMarshal(got))
		})
	}
}
