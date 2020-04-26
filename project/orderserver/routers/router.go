package routers

import "github.com/gin-gonic/gin"

func InitRouters() *gin.Engine {
	ginRouter := gin.Default()
	ginRouter.POST("/orders/", func(context *gin.Context) {
		context.String(200, "get orderinfos")
	})

	return ginRouter
}
