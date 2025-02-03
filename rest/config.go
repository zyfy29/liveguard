package rest

import (
	"bearguard/cm"
	"github.com/gin-gonic/gin"
)

func setConfigRoutes(r *gin.Engine) {
	g := r.Group("/config")
	g.GET("/pocket", getPocketConfig)
	g.GET("/aai", getAaiConfig)
	g.GET("/medium", getMediumConfig)
	g.PUT("/pocket", setPocketConfig)
	g.PUT("/aai", setAaiConfig)
	g.PUT("/medium", setMediumConfig)
}

func getPocketConfig(c *gin.Context) {
	config := cm.GetConfig()
	ResponseOk(c, config.Pocket)
}

func getAaiConfig(c *gin.Context) {
	config := cm.GetConfig()
	ResponseOk(c, config.AAI)
}

func getMediumConfig(c *gin.Context) {
	config := cm.GetConfig()
	ResponseOk(c, config.Medium)
}

func setPocketConfig(c *gin.Context) {
	var req cm.PocketConfig
	if err := c.BindJSON(&req); err != nil {
		ResponseClientError(c, err)
		return
	}
	count := cm.SetConfig(&cm.Config{Pocket: req})
	ResponseOk(c, gin.H{
		"updated_count": count,
	})
}

func setAaiConfig(c *gin.Context) {
	var req cm.AAIConfig
	if err := c.BindJSON(&req); err != nil {
		ResponseClientError(c, err)
		return
	}
	count := cm.SetConfig(&cm.Config{AAI: req})
	ResponseOk(c, gin.H{
		"updated_count": count,
	})
}

func setMediumConfig(c *gin.Context) {
	var req cm.MediumConfig
	if err := c.BindJSON(&req); err != nil {
		ResponseClientError(c, err)
		return
	}
	count := cm.SetConfig(&cm.Config{Medium: req})
	ResponseOk(c, gin.H{
		"updated_count": count,
	})
}
