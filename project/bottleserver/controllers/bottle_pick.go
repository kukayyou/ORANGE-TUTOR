package controllers

import (
	"bottleserver/dao"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/kukayyou/commonlib/mylog"
)

type BottlePickController struct {
	BaseController
}

type BottlePickRequest struct {
	OpenID string `json:"openId"`
	Token  string `json:"token"`
}

func (this BottlePickController) BottlePickApi(c *gin.Context) {
	defer this.FinishResponse(c)
	this.Prepare(c)
	var (
		params BottlePickRequest
		err    error
	)
	if err = json.Unmarshal(this.ReqParams, &params); err != nil {
		mylog.Error("requestID:%s, Unmarshal request failed!", this.GetRequestId())
		this.Resp.Code = PARAMS_PARSE_ERROR
		this.Resp.Msg = "Unmarshal request failed!"
		return
	}
	//校验token
	if err := this.CheckToken(params.OpenID, params.Token); err != nil {
		mylog.Error("requestID:%s, UserCheck failed!", this.GetRequestId())
		this.Resp.Code = TOKEN_CHECK_ERROR
		this.Resp.Msg = "token check failed!"
		return
	}
	//查询用户信息
	miniUserInfoResp, err := dao.GetMiniUserInfo(this.GetRequestId(), params.OpenID)
	if err != nil {
		this.Resp.Code = BOTTLE_PICK_ERROR
		this.Resp.Msg = err.Error()
		return
	}
	bottle := &dao.BottleInfo{
		DestUserID:     params.OpenID,
		DestUserGender: miniUserInfoResp.Data.Gender,
		City:           miniUserInfoResp.Data.City,
	}
	//查询瓶子信息
	bottle, err = bottle.PickBottle(this.GetRequestId())
	if err != nil {
		this.Resp.Code = BOTTLE_PICK_ERROR
		this.Resp.Msg = err.Error()
		return
	}
	//更新瓶子接收人信息
	bottle.UpdatePickedBottle(this.GetRequestId())

	this.Resp.Data = bottle
	return
}
