package rest

import (
	"bearguard/repo"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strconv"
)

func setTaskRoutes(r *gin.Engine) {
	g := r.Group("/task")
	g.GET("/", getTasks)
	g.DELETE("/:id", deleteTask)
	g.GET("/:id", getTaskDetail)
	g.POST("/", createTask)
	g.POST("/retry", restoreTask)
}

func getTasks(c *gin.Context) {
	taskStatus := c.Query("status")
	var tasks []repo.Task
	var err error
	if len(taskStatus) == 0 {
		tasks, err = repo.GetDBTasks()
	} else {
		tasks, err = repo.GetDBTasksByStatus(taskStatus, 0)
	}
	if err != nil {
		ResponseServerError(c, err)
		return
	}
	ResponseOk(c, tasks)
}

func deleteTask(c *gin.Context) {
	taskPk, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := repo.DeleteDBTask(taskPk); err != nil {
		ResponseServerError(c, err)
		return
	}
	ResponseOk(c, nil)
}

type taskAndDetail struct {
	repo.Task
	repo.TaskDetail
}

func getTaskDetail(c *gin.Context) {
	taskPk, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	task, err := repo.GetDBTaskByID(taskPk)
	if err != nil {
		ResponseClientError(c, err)
	}
	detail := repo.TaskDetail{}
	if err := json.Unmarshal([]byte(task.Details), &detail); err != nil {
		ResponseServerError(c, err)
		return
	}
	ResponseOk(c, taskAndDetail{
		Task:       task,
		TaskDetail: detail,
	})
}

func createTask(c *gin.Context) {
	var req struct {
		LiveID string `json:"live_id"`
	}
	if err := c.BindJSON(&req); err != nil {
		ResponseClientError(c, err)
		return
	}
	if err := repo.CreateLiveTask(req.LiveID); err != nil {
		ResponseServerError(c, err)
		return
	}
	ResponseOk(c, nil)
}

func restoreTask(c *gin.Context) {
	var req struct {
		TaskID int64 `json:"task_id"`
	}
	if err := c.BindJSON(&req); err != nil {
		ResponseClientError(c, err)
		return
	}
	if err := repo.RestoreTask(req.TaskID); err != nil {
		ResponseServerError(c, err)
		return
	}
	ResponseOk(c, nil)
}
