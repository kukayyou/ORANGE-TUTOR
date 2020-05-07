package routers

import (
	"github.com/gin-gonic/gin"
	"orderserver/controllers"
)

func InitRouters() *gin.Engine {
	ginRouter := gin.Default()
	root := ginRouter.Group("/orderserver")
	order := root.Group("/order")
	order.POST("/infos", controllers.GetOrderController{}.GetOrderInfosApi)

	return ginRouter
}
