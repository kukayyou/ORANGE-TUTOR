package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/kukayyou/commonlib/mylog"
	"userserver/dao"
)

type MiniUserInfosController struct {
	BaseController
}

type MiniUserInfosRequest struct {
	OpenID      string `json:"openId"`      //用户唯一标识
	Token       string `json:"token"`       //客户端调用使用的token
	ServerToken string `json:"serverToken"` //server调用使用的token
}

type MiniUserInfosResp struct {
	OpenID string `json:"openId"` //用户唯一标识
	Gender int64  `json:"gender"` //token
	City   string `json:"city"`   //token
}

func (this MiniUserInfosController) MiniUserInfosApi(c *gin.Context) {
	//构建返回参数
	defer this.FinishResponse(c)
	this.Prepare(c)

	var (
		params MiniUserInfosRequest
		err    error
	)
	//解析参数
	if err = json.Unmarshal(this.ReqParams, &params); err != nil {
		mylog.Error("requestID:%s, Unmarshal request failed!", this.GetRequestId())
		this.Resp.Code = PARAMS_PARSE_ERROR
		this.Resp.Msg = "Unmarshal request failed!"
		return
	}
	this.ServerToken = params.ServerToken
	//校验token
	err = this.CheckToken(params.OpenID, params.Token)
	if err != nil {
		this.Resp.Code = USER_CHECK_ERROR
		this.Resp.Msg = err.Error()
		return
	}
	//查询数据库，用户信息
	miniUserInfo := &dao.MiniUserInfo{OpenID: params.OpenID}
	miniUserInfo, err = miniUserInfo.Query(this.GetRequestId())
	if err != nil {
		this.Resp.Code = USER_GET_INFOS_ERROR
		this.Resp.Msg = err.Error()
		return
	}
	this.Resp.Data = miniUserInfo
	return
}
