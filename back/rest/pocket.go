package rest

import (
	"bearguard/cm"
	"bearguard/pocket"
	"bearguard/repo"
	"github.com/gin-gonic/gin"
	"strconv"
)

func setPocketRoutes(r *gin.Engine) {
	g := r.Group("/pocket")
	g.GET("/member", getMembers)
	g.GET("/live", getLivesByMember)
	g.GET("/live/duration", getLiveDuration)
}

// name and member_id only
func getMembers(c *gin.Context) {
	members, err := repo.GetDBMembers()
	if err != nil {
		ResponseServerError(c, err)
		return
	}
	ResponseOk(c, members)
}

func getLivesByMember(c *gin.Context) {
	ownerID := c.Query("owner_id")
	nextTime := c.Query("next_time")
	if nextTime == "" {
		nextTime = "0"
	}

	nextTimeInt, _ := strconv.ParseInt(nextTime, 10, 64)
	lives, nextTimeInt, err := pocket.GetClient().GetHistoryLives(ownerID, nextTimeInt)
	if err != nil {
		ResponseClientError(c, err)
		return
	}
	ResponseOk(c, gin.H{
		"lives":     lives,
		"next_time": nextTimeInt,
	})
}

func getLiveDuration(c *gin.Context) {
	playlistUrl := c.Query("playlist_url")
	seconds, err := cm.GetPlaylistDuration(playlistUrl)
	if err != nil {
		ResponseServerError(c, err)
		return
	}
	ResponseOk(c, gin.H{
		"duration": seconds,
	})
}
