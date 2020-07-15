package controllers

import (
	"bottleserver/dao"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/kukayyou/commonlib/mylog"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BottleThrowController struct {
	BaseController
}

type BottleThrowRequest struct {
	OpenID  string            `json:"openId"`        //微信小程序用户唯一标识
	Token   string            `json:"token"`         //token
	Content dao.BottleContent `json:"bottleContent"` //瓶子信息
}

func (this BottleThrowController) BottleThrowApi(c *gin.Context) {
	defer this.FinishResponse(c)
	this.Prepare(c)

	var (
		params BottleThrowRequest
		err    error
	)
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
	//查询用户信息
	miniUserInfoResp, err := dao.GetMiniUserInfo(this.GetRequestId(), params.OpenID)
	if err != nil {
		this.Resp.Code = BOTTLE_THROW_ERROR
		this.Resp.Msg = err.Error()
		return
	}

	bottle := dao.BottleInfo{
		BottleID:       primitive.NewObjectID(),
		SourUserID:     params.OpenID,
		SourUserGender: miniUserInfoResp.Data.Gender,
		City:           miniUserInfoResp.Data.City,
	}
	bottle.Content = append(bottle.Content, params.Content)

	//瓶子信息入库
	if err = bottle.ThrowBottle(this.GetRequestId()); err != nil {
		this.Resp.Code = BOTTLE_THROW_ERROR
		this.Resp.Msg = err.Error()
		return
	}
	return
}
