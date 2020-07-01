package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/kukayyou/commonlib/mylog"
)

type SendMsgController struct {
	BaseController
}

type SendMsgRequest struct {
	Token string `json:"token"` //token
}

func (this SendMsgController) SendMsgApi(c *gin.Context) {
	defer this.FinishResponse(c)
	this.Prepare(c)

	var (
		params SendMsgRequest
		err    error
	)
	if err = json.Unmarshal(this.ReqParams, &params); err != nil {
		mylog.Error("requestID:%s, Unmarshal request failed!", this.GetRequestId())
		this.Resp.Code = PARAMS_PARSE_ERROR
		this.Resp.Msg = "Unmarshal request failed!"
		return
	}
}
