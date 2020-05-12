package routers

import (
	"github.com/gin-gonic/gin"
	"userserver/controllers"
)

func InitRouters() *gin.Engine {
	ginRouter := gin.Default()
	root := ginRouter.Group("/userserver")
	userGroup := root.Group("/user")
	userGroup.POST("/registry", controllers.UserRegisterController{}.UserRegisterApi)
	userGroup.POST("/login", controllers.UserLoginController{}.UserLoginApi)
	userGroup.POST("/infos", controllers.UserInfosController{}.GetUserInfosApi)
	userGroup.POST("/update", controllers.UserInfosController{}.UpdateUserInfosApi)
	return ginRouter
}
