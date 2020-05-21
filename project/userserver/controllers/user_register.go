package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/kukayyou/commonlib/mylog"
	"github.com/kukayyou/commonlib/token"
	"userserver/dao"
)

type UserRegisterController struct {
	BaseController
}

type RegisterRequest struct {
	UserName string `json:"userName"` //用户名
	UserType int    `json:"userType"` //0：大学生家教，1：家长
	Passwd   string `json:"passwd"`   //密码
}

func (this UserRegisterController) UserRegisterApi(c *gin.Context) {
	//返回响应
	defer this.FinishResponse(c)
	var (
		params RegisterRequest
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

	userInfo := dao.UserInfo{
		UserName: params.UserName,
		Passwd:   params.Passwd,
	}
	//注册用户
	if userInfo, err = userInfo.RegisterUserInfo(this.GetRequestId()); err != nil {
		this.Resp.Code = USER_REGISTER_ERROR
		this.Resp.Msg = "register user failed!"
		return
	} else {
		userTokenInfo := token.UserInfo{
			UserID:   userInfo.UserID,
			UserName: userInfo.UserName,
			Passwd:   userInfo.Passwd,
		}
		//创建token
		if userInfo.Token, err = token.CreateToken(userTokenInfo, int64(^uint(0)>>1)); err != nil {
			this.Resp.Code = USER_REGISTER_ERROR
			this.Resp.Msg = "create  token failed!"
			return
		}
	}
	userInfo.Passwd = "" //隐去密码
	this.Resp.Data = userInfo
	return
}
