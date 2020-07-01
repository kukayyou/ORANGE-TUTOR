package routers

import (
	"github.com/gin-gonic/gin"
	"notifyserver/controllers"
)

func InitRouters() *gin.Engine {
	ginRouter := gin.Default()
	root := ginRouter.Group("/notifyserver")
	//验证码发送
	captcha := root.Group("/captcha")
	captcha.POST("/mailsender", controllers.SendMailController{}.SendMailApi)
	captcha.POST("/msgsender", controllers.SendMsgController{}.SendMsgApi)
	//验证码合法性校验
	check := root.Group("/check")
	check.POST("/mailcode", controllers.CheckMailCodeController{}.CheckMailCodeApi)
	check.POST("/msgcode", controllers.CheckMsgCodeController{}.CheckMsgCodeApi)

	return ginRouter
}
