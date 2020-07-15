package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/kukayyou/commonlib/mylog"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"msgserver/dao"
)

type MsgDeleteController struct {
	BaseController
}

type MsgDeleteRequest struct {
	OpenID   string `json:"openId"`
	Token    string `json:"token"`
	BottleID string `json:"bottleId"` //瓶子id
}

func (this MsgDeleteController) MsgDeleteApi(c *gin.Context) {
	defer this.FinishResponse(c)
	this.Prepare(c)
	var (
		params MsgDeleteRequest
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
	delUserID := make([]string, 0)
	delUserID = append(delUserID, params.OpenID)
	bottleInfo := &dao.BottleInfo{
		BottleID:   id,
		DeleteUser: delUserID,
	}
	if err := bottleInfo.DeleteMsg(this.GetRequestId()); err != nil {
		this.Resp.Code = MSG_DELETE_ERROR
		this.Resp.Msg = err.Error()
		return
	}

	return
}
