package routers

import (
	"github.com/gin-gonic/gin"
	"userserver/controllers"
)

func InitRouters() *gin.Engine {
	ginRouter := gin.Default()
	root := ginRouter.Group("/userserver")
	userGroup := root.Group("/user")
	userGroup.POST("/infos", controllers.GetUserInfosApi)
	return ginRouter
}
