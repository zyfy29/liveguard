package rest

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
)

var setRouteFuncs = []func(r *gin.Engine){
	setTaskRoutes,
	setPocketRoutes,
	setConfigRoutes,
}

func setRoute(r *gin.Engine) {
	r.GET("/ping", ping)

	for _, setRouteFunc := range setRouteFuncs {
		setRouteFunc(r)
	}
}

func ping(c *gin.Context) {
	ResponseOk(c, gin.H{
		"message": "pong",
	})
}

func Startup() {
	r := gin.Default()
	setRoute(r)
	if err := r.Run(fmt.Sprintf(":%d", viper.GetInt("app.back_port"))); err != nil {
		log.Panicf("Failed to start server: %v", err)
	}
}
