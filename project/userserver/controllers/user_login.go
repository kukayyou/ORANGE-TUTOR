package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/kukayyou/commonlib/mylog"
	"github.com/kukayyou/commonlib/token"
	"userserver/dao"
)

type UserLoginController struct {
	BaseController
}

type LoginRequest struct {
	UserID   string  `json:"userId"`
	UserName string `json:"userName"`
	Passwd   string `json:"passwd"`
}

func (this UserLoginController) UserLoginApi(c *gin.Context) {
	defer this.FinishResponse(c)
	this.Prepare(c)

	var (
		params LoginRequest
		err    error
	)
	if err = json.Unmarshal(this.ReqParams, &params); err != nil {
		mylog.Error("requestID:%s, Unmarshal request failed!", this.GetRequestId())
		this.Resp.Code = PARAMS_PARSE_ERROR
		this.Resp.Msg = "Unmarshal request failed!"
		return
	}

	userInfo := &dao.UserInfo{UserName: params.UserName, UserID: params.UserID}
	if userInfo, err = userInfo.GetUserInfo(this.GetRequestId()); err == nil {
		if userInfo.Passwd != params.Passwd {
			this.Resp.Code = USER_LOGIN_ERROR
			this.Resp.Msg = "user or password is wrong!"
			return
		}
		userTokenInfo := token.UserInfo{
			UserID:   userInfo.UserID,
		}
		userInfo.Token, _ = token.CreateUserToken(userTokenInfo, int64(^uint(0)>>1))
	} else {
		this.Resp.Code = USER_LOGIN_ERROR
		this.Resp.Msg = "get user info failed!"
		return
	}
	userInfo.Passwd = "" //隐去密码
	this.Resp.Data = userInfo
	return
}
