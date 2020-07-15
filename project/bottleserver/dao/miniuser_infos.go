package dao

import (
	"fmt"
	"github.com/kukayyou/commonlib/myhttp"
	"github.com/kukayyou/commonlib/mylog"
	"github.com/kukayyou/commonlib/token"
	"github.com/goinggo/mapstructure"
)

type MiniUserInfoResp struct {
	ErrCode int64        `json:"errcode"` //错误码
	ErrMsg  string       `json:"errmsg"`  //错误信息
	Data    MiniUserInfo `json:"data"`    //数据
}

type MiniUserInfo struct {
	OpenID     string `json:"openId"`     //微信登录返回的openid
	SessionKey string `json:"sessionKey"` //微信登录返回的session_key
	Gender     int64  `json:"gender"`     //性别，0：未知， 1：男，2：女
	City       string `json:"city"`       //用户所在城市
	Token      string `json:"token"`
}

//查询用户信息
func GetMiniUserInfo(requestID, openID string) (*MiniUserInfoResp, error) {
	//查询用户信息
	req := make(map[string]string)
	req["openId"] = openID
	req["serverToken"], _ = token.CreateServerToken("bottleserver", int64(^uint(0)>>1))

	resp := myhttp.RequestWithHytrix("api.tutor.com.userserver", "/userserver/miniuser/infos", req)
	if resp == nil {
		mylog.Error("requestID:%s, get mini user info is null", requestID)
		return nil, fmt.Errorf("get mini user info is null")
	}
	var miniUserInfoResp MiniUserInfoResp
	if err := mapstructure.Decode(resp, &miniUserInfoResp); err != nil {
		mylog.Error("requestID:%s, Decode user info error:%s", requestID, err.Error())
		return nil, err
	}
	return &miniUserInfoResp, nil
}
