package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/kukayyou/commonlib/mylog"
	"strings"
)

type CheckMailCodeController struct {
	BaseController
}

type CheckMailCodeRequest struct {
	UserID  int64  `json:"userId"`
	Captcha string `json:"captcha"`
	Token   string `json:"token"` //token
}

func (this CheckMailCodeController) CheckMailCodeApi(c *gin.Context) {
	defer this.FinishResponse(c)
	this.Prepare(c)

	var (
		params CheckMailCodeRequest
		err    error
	)
	if err = json.Unmarshal(this.ReqParams, &params); err != nil {
		mylog.Error("requestID:%s, Unmarshal request failed!", this.GetRequestId())
		this.Resp.Code = PARAMS_PARSE_ERROR
		this.Resp.Msg = "Unmarshal request failed!"
		return
	}
	if err := this.CheckToken(params.UserID, params.Token); err != nil {
		mylog.Error("requestID:%s, UserCheck failed!", this.GetRequestId())
		return
	}
	if err := checkCaptcha(this.GetRequestId(), params.Captcha, params.UserID); err != nil {
		if strings.Contains(err.Error(), "time out") {
			this.Resp.Code = CHECK_MAIL_CAPTCHA_TIMEOUT_ERROR
		} else if strings.Contains(err.Error(), "not equal") {
			this.Resp.Code = CHECK_MAIL_CAPTCHA_ERROR
		}
		this.Resp.Msg = "check captcha failed!"
		return
	}
}
