package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/kukayyou/commonlib/mylog"
	"io/ioutil"
)

type BaseController struct {
	mylog.LogInfo
	ReqParams []byte
}

func (bc *BaseController) Prepare(c *gin.Context) {
	bc.SetRequestId()

	bc.ReqParams, _ = ioutil.ReadAll(c.Request.Body)

	mylog.Info("requestId:%s, params : %s", bc.GetRequestId(), string(bc.ReqParams))
}
