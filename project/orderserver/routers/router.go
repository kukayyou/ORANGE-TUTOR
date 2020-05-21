package routers

import (
	"github.com/gin-gonic/gin"
	"orderserver/controllers"
)

func InitRouters() *gin.Engine {
	ginRouter := gin.Default()
	root := ginRouter.Group("/orderserver")
	//签约管理
	order := root.Group("/order")
	order.POST("/infos", controllers.GetOrderController{}.GetOrderInfosApi)
	//家长需求管理
	demand := root.Group("/demand")
	demand.POST("/create", controllers.DemandInfoController{}.CreateDemandApi)
	demand.POST("/update", controllers.DemandInfoController{}.UpdateDemandInfosApi)
	demand.POST("/querybyuserid", controllers.DemandInfoController{}.GetDemandInfosApi)
	demand.POST("/delete", controllers.DemandInfoController{}.DeleteDemandInfosApi)
	//教师技能管理
	skill := root.Group("/skill")
	skill.POST("/create", controllers.SkillInfoController{}.CreateSkillApi)
	skill.POST("/update", controllers.SkillInfoController{}.UpdateSkillInfosApi)
	skill.POST("/querybyuserid", controllers.SkillInfoController{}.GetSkillInfosApi)
	skill.POST("/delete", controllers.SkillInfoController{}.DeleteSkillInfosApi)

	return ginRouter
}
