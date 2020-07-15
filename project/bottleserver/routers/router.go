package routers

import (
	"bottleserver/controllers"
	"github.com/gin-gonic/gin"
)

func InitRouters() *gin.Engine {
	ginRouter := gin.Default()
	root := ginRouter.Group("/bottleserver")
	//andriod、pc、mac用户
	bottleGroup := root.Group("/bottle")
	bottleGroup.POST("/list", controllers.BottleListController{}.BottleListApi)    //获取历史瓶子列表
	bottleGroup.POST("/throw", controllers.BottleThrowController{}.BottleThrowApi) //扔瓶子
	bottleGroup.POST("/pick", controllers.BottlePickController{}.BottlePickApi)    //捡瓶子
	bottleGroup.POST("/drop", controllers.BottleDropController{}.BottleDropApi)    //丢弃捡到的瓶子
	return ginRouter
}
