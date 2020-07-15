package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/kukayyou/commonlib/mylog"
	"github.com/kukayyou/commonlib/token"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"userserver/config"
	"userserver/dao"
	"userserver/usercheck"
)

type MiniUserLoginController struct {
	BaseController
}

type MiniUserLoginRequest struct {
	WxCode string `json:"wxcode"` //wx.login返回的code
	Gender int64  `json:"gender"` //性别，0：未知，1：男，2：女
	City   string `json:"city"`   //用户所在城市
}

type MiniUserLoginResp struct {
	OpenID string `json:"openId"` //用户唯一标识
	Token  string `json:"token"`  //token
}

func (this MiniUserLoginController) MiniUserLoginApi(c *gin.Context) {
	//构建返回参数
	defer this.FinishResponse(c)
	this.Prepare(c)

	var (
		params MiniUserLoginRequest
		err    error
	)
	//解析参数
	if err = json.Unmarshal(this.ReqParams, &params); err != nil {
		mylog.Error("requestID:%s, Unmarshal request failed!", this.GetRequestId())
		this.Resp.Code = PARAMS_PARSE_ERROR
		this.Resp.Msg = "Unmarshal request failed!"
		return
	}
	wlReq := &usercheck.WxLoginReq{
		AppID:     	config.AppID,
		Secret: 	config.AppSecret,
		Js_Code:    params.WxCode,
		Grant_Type: "authorization_code",
	}
	//请求微信验证auth.code2Session接口，验证code
	resp, err := wlReq.CheckMiniUser(this.GetRequestId())
	if err != nil {
		if resp != nil && resp.ErrCode != 0{
			this.Resp.Code = resp.ErrCode
		}else {
			this.Resp.Code = USER_LOGIN_ERROR
		}
		this.Resp.Msg = err.Error()
		return
	}
	//创建token
	userToken, err := token.CreateUserToken(token.UserInfo{UserID: "1000002"}, int64(^uint(0)>>1))
	if err != nil {
		this.Resp.Code = USER_LOGIN_ERROR
		this.Resp.Msg = err.Error()
		return
	}
	miniUserInfo := &dao.MiniUserInfo{
		ID:primitive.NewObjectID(),
		OpenID:     "1000002",//resp.OpenID,
		SessionKey: "",//resp.SessionKey,
		Gender:     params.Gender,
		City:       params.City,
		Token:userToken,
	}
	//用户信息入库
	if err := miniUserInfo.Create(this.GetRequestId()); err != nil{
		mylog.Error("requestID:%s, MiniUser insert into db failed! error:%s", this.GetRequestId(), err.Error())
		this.Resp.Code = USER_LOGIN_ERROR
		this.Resp.Msg = err.Error()
		return
	}
	//返回登录数据
	this.Resp.Data = MiniUserLoginResp{OpenID:"1000001", Token:userToken}
	return
}
