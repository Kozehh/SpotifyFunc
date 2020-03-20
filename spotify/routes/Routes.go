package routes

import (
	"../Controllers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/v1")
	{
		v1.GET("followers", Controllers.GetFollowers)
		v1.POST("followers", Controllers.CreateAFollower)
		v1.GET("followers/:name", Controllers.GetAFollower)
		v1.PUT("follower/:name", Controllers.UpdateAFollower)
		v1.DELETE("follower/:name", Controllers.DeleteAFollower)
	}
	return r
}
