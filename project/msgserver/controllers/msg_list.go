package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/kukayyou/commonlib/mylog"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"msgserver/dao"
)

type MsgListController struct {
	BaseController
}

type MsgListRequest struct {
	OpenID   string `json:"openId"`   //微信小程序用户唯一标识
	Token    string `json:"token"`    //token
	BottleID string `json:"bottleId"` //瓶子id
}

func (this MsgListController) MsgListApi(c *gin.Context) {
	//返回响应
	defer this.FinishResponse(c)
	var (
		params MsgListRequest
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
	//转换瓶子id为primitive.ObjectID
	id, err := primitive.ObjectIDFromHex(params.BottleID)
	if err != nil {
		this.Resp.Code = PARAMS_PARSE_ERROR
		this.Resp.Msg = err.Error()
		return
	}
	//token校验
	if err := this.CheckToken(params.OpenID, params.Token); err != nil {
		mylog.Error("requestID:%s, UserCheck failed!", this.GetRequestId())
		this.Resp.Code = TOKEN_CHECK_ERROR
		this.Resp.Msg = "token check failed!"
		return
	}
	bottle := &dao.BottleInfo{BottleID: id}
	bottleInfo, err := bottle.GetBottleInfoByID(this.GetRequestId())
	if err != nil {
		this.Resp.Code = MSG_LIST_ERROR
		this.Resp.Msg = err.Error()
	}
	this.Resp.Data = bottleInfo
	return
}
