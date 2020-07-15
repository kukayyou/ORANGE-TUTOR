package routers

import (
	"github.com/gin-gonic/gin"
	"userserver/controllers"
)

func InitRouters() *gin.Engine {
	ginRouter := gin.Default()
	root := ginRouter.Group("/userserver")
	//andriod、pc、mac用户
	userGroup := root.Group("/user")
	userGroup.POST("/registry", controllers.UserRegisterController{}.UserRegisterApi)
	userGroup.POST("/login", controllers.UserLoginController{}.UserLoginApi)
	userGroup.POST("/infos", controllers.UserInfosController{}.GetUserInfoApi)
	userGroup.POST("/update", controllers.UserInfosController{}.UpdateUserInfoApi)
	//微信小程序用户
	mini := root.Group("/miniuser")
	mini.POST("/login", controllers.MiniUserLoginController{}.MiniUserLoginApi)
	mini.POST("/infos", controllers.MiniUserInfosController{}.MiniUserInfosApi)
	return ginRouter
}
