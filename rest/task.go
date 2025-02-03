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
	g.GET("/:id", getTaskDetail)
	g.POST("/", createTask)
	g.POST("/retry", restoreTasks)
}

func getTasks(c *gin.Context) {
	tasks, err := repo.GetDBTasks()
	if err != nil {
		ResponseServerError(c, err)
		return
	}
	ResponseOk(c, tasks)
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

func restoreTasks(c *gin.Context) {
	if err := repo.RestoreTasks(); err != nil {
		ResponseServerError(c, err)
		return
	}
	ResponseOk(c, nil)
}
