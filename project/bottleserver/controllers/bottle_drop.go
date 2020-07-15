package controllers

import (
	"bottleserver/dao"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/kukayyou/commonlib/mylog"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BottleDropController struct {
	BaseController
}

type BottleDropRequest struct {
	OpenID   string `json:"openId"`  //微信小程序用户唯一标识
	Token    string `json:"token"`   //token
	BottleID string `json"bottleId"` //瓶子id
}

func (this BottleDropController) BottleDropApi(c *gin.Context) {
	defer this.FinishResponse(c)
	this.Prepare(c)
	var (
		params BottleDropRequest
		err    error
	)
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
	//校验token
	if err := this.CheckToken(params.OpenID, params.Token); err != nil {
		mylog.Error("requestID:%s, UserCheck failed!", this.GetRequestId())
		this.Resp.Code = TOKEN_CHECK_ERROR
		this.Resp.Msg = "token check failed!"
		return
	}

	bottle := &dao.BottleInfo{BottleID: id}
	//查询瓶子信息
	err = bottle.DropBottle(this.GetRequestId())
	if err != nil {
		this.Resp.Code = BOTTLE_DROP_ERROR
		this.Resp.Msg = err.Error()
		return
	}

	return
}
