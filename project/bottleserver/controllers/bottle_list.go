package controllers

import (
	"bottleserver/dao"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/kukayyou/commonlib/mylog"
)

type BottleListController struct {
	BaseController
}

type BottleListRequest struct {
	OpenID string `json:"openId"` //微信小程序用户唯一标识
	Token  string `json:"token"`  //token
}

func (this BottleListController) BottleListApi(c *gin.Context) {
	//返回响应
	defer this.FinishResponse(c)
	var (
		params BottleListRequest
		err    error
	)

	//获取请求参数
	this.Prepare(c)
	//解析请求参数
	if err = json.Unmarshal(this.ReqParams, &params); err != nil {
		mylog.Error("requestID:%s, Unmarshal request failed!", this.GetRequestId())
		this.Resp.Code = PARAMS_PARSE_ERROR
		this.Resp.Msg = "Unmarshal request failed!"
		return
	}

	//token校验
	if err := this.CheckToken(params.OpenID, params.Token); err != nil {
		mylog.Error("requestID:%s, UserCheck failed!", this.GetRequestId())
		this.Resp.Code = TOKEN_CHECK_ERROR
		this.Resp.Msg = "token check failed!"
		return
	}
	bottle := &dao.BottleInfo{SourUserID: params.OpenID, DestUserID: params.OpenID}
	bottleList, err := bottle.GetUserHistoryBottleList(this.GetRequestId())
	if err != nil {
		this.Resp.Code = BOTTLE_GET_INFOS_ERROR
		this.Resp.Msg = err.Error()
	}
	this.Resp.Data = bottleList
	return
}
