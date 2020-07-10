package usercheck

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kukayyou/commonlib/mylog"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"userserver/config"
)

type WxLoginResp struct {
	OpenID     string `json:"openid"`      //用户唯一标识
	SessionKey string `json:"session_key"` //会话密钥
	UnionID    string `json:"unionid"`     //用户在开放平台的唯一标识符
	ErrCode    int64  `json:"errcode"`     //错误码
	ErrMsg     string `json:"errmsg"`      //错误信息
}

type WxLoginReq struct {
	AppID      string `json:"appid"`      //小程序 appId
	Secret     string `json:"secret"`     //小程序 appSecret
	Js_Code    string `json:"js_code"`    //登录时获取的 code
	Grant_Type string `json:"grant_type"` //授权类型，此处只需填写 authorization_code
}

func (wr *WxLoginReq) CheckMiniUser(requestID string) (*WxLoginResp, error) {
	var wxResp WxLoginResp
	reqUrl := config.WxDomain + "/sns/jscode2session"
	reqUrl = buildQueryUrl(reqUrl, *wr)
	resp, err := http.Get(reqUrl)
	defer resp.Body.Close()
	if err != nil {
		mylog.Error("requestID:%s, CheckMiniUser http request failed, error:%s", requestID, err.Error())
		return nil, err
	} else {
		body, _ := ioutil.ReadAll(resp.Body)

		if err = json.Unmarshal(body, &wxResp); err != nil {
			mylog.Error("requestID:%s, Unmarshal http response failed, error:%s", requestID, err.Error())
			return nil, err
		} else if wxResp.ErrCode != 0 {
			/*错误码说明：
			-1 ：系统繁忙，此时请开发者稍候再试
			0：请求成功
			40029：code 无效
			45011：频率限制，每个用户每分钟100次
			*/
			mylog.Error("requestID:%s, wx login errcode is:%d", requestID, wxResp.ErrCode)
			return &wxResp, fmt.Errorf("wx login errcode is:%d", wxResp.ErrCode)
		}
	}

	return &wxResp, nil
}

func buildQueryUrl(reqUrl string, params WxLoginReq) string {
	var paramBody string
	t := reflect.TypeOf(params)
	v := reflect.ValueOf(params)

	var buf bytes.Buffer
	for i := 0; i < t.NumField(); i++ {
		buf.WriteString(strings.ToLower(t.Field(i).Name))
		buf.WriteByte('=')
		buf.WriteString(strings.ToLower(v.Field(i).String()))
		buf.WriteByte('&')
	}
	paramBody = buf.String()
	paramBody = paramBody[0 : len(paramBody)-1]

	if len(paramBody) > 0 {
		if strings.Index(reqUrl, "?") != -1 {
			reqUrl += "&" + paramBody
		} else {
			reqUrl = reqUrl + "?" + paramBody
		}
	}

	return reqUrl
}
