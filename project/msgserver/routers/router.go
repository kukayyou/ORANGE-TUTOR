package routers

import (
	"github.com/gin-gonic/gin"
	"msgserver/controllers"
)

func InitRouters() *gin.Engine {
	ginRouter := gin.Default()
	root := ginRouter.Group("/msgserver")
	//andriod、pc、mac用户
	bottleGroup := root.Group("/msg")
	bottleGroup.POST("/list", controllers.MsgListController{}.MsgListApi)       //获取瓶子消息
	bottleGroup.POST("/replay", controllers.MsgReplayController{}.MsgReplayApi) //回复瓶子
	bottleGroup.POST("/delete", controllers.MsgDeleteController{}.MsgDeleteApi) //删除瓶子
	return ginRouter
}
