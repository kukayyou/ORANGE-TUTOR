package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/kukayyou/commonlib/mylog"
	"userserver/dao"
)

type UserInfosController struct {
	BaseController
}

type GetUserInfosRequest struct {
	UserID int64  `json:"UserId"`
	Token  string `json:"token"`
}

func (this UserInfosController) GetUserInfoApi(c *gin.Context) {
	defer this.FinishResponse(c)
	this.Prepare(c)
	var (
		params GetUserInfosRequest
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

	userInfo := &dao.UserInfo{UserID: params.UserID}
	if userInfo, err = userInfo.GetUserInfo(this.GetRequestId()); err != nil {
		this.Resp.Code = USER_GET_INFOS_ERROR
		this.Resp.Msg = "GetUserInfo failed!"
		return
	} else {
		userInfo.Passwd = "" //隐去密码
		this.Resp.Data = userInfo
	}

	return
}
