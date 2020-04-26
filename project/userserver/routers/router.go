package routers

import "github.com/gin-gonic/gin"

func InitRouters() *gin.Engine {
	ginRouter := gin.Default()
	ginRouter.POST("/users/", func(context *gin.Context) {
		context.String(200, "get userinfos")
	})

	return ginRouter
}
