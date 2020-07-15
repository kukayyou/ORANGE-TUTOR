package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/kukayyou/commonlib/mylog"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"msgserver/dao"
)

type MsgReplayController struct {
	BaseController
}

type MsgReplayRequest struct {
	OpenID   string            `json:"openId"`        //微信小程序用户唯一标识
	Token    string            `json:"token"`         //token
	BottleID string            `json:"bottleId"`      //瓶子id
	Content  dao.BottleContent `json:"bottleContent"` //瓶子回复内容
}

func (this MsgReplayController) MsgReplayApi(c *gin.Context) {
	defer this.FinishResponse(c)
	this.Prepare(c)

	var (
		params MsgReplayRequest
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
	//token校验
	if err := this.CheckToken(params.OpenID, params.Token); err != nil {
		mylog.Error("requestID:%s, UserCheck failed!", this.GetRequestId())
		this.Resp.Code = TOKEN_CHECK_ERROR
		this.Resp.Msg = "token check failed!"
		return
	}
	//查询用户信息
	miniUserInfoResp, err := dao.GetMiniUserInfo(this.GetRequestId(), params.OpenID)
	if err != nil {
		this.Resp.Code = MSG_REPLAY_ERROR
		this.Resp.Msg = err.Error()
		return
	}

	bottleContent := make([]dao.BottleContent, 0)
	bottleContent = append(bottleContent, params.Content)
	bottleInfo := &dao.BottleInfo{
		BottleID: id,
		Content:  bottleContent,
	}
	//更新数据库瓶子消息
	if err := bottleInfo.ReplayMsg(this.GetRequestId()); err != nil {
		this.Resp.Code = MSG_REPLAY_ERROR
		this.Resp.Msg = err.Error()
		return
	}
	this.Resp.Data = miniUserInfoResp
	return
}

func PushMsg(bottleContent dao.BottleContent) error {
	return nil
}
