package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kukayyou/commonlib/mylog"
	"github.com/kukayyou/commonlib/token"
	"io/ioutil"
)

//错误码
const (
	PARAMS_PARSE_ERROR      = 1001 + iota //请求参数解析错误
	TOKEN_CHECK_ERROR       = 1001 + iota //token验证错误
	USER_CHECK_ERROR        = 1001 + iota //用户验证错误，非本人操作
	USER_REGISTER_ERROR     = 1001 + iota //注册错误
	USER_LOGIN_ERROR        = 1001 + iota //登录错误
	USER_GET_INFOS_ERROR    = 1001 + iota //获取用户信息错误
	USER_UPDATE_INFOS_ERROR = 1001 + iota //更新用户信息错误
	DEMAND_CREATE_ERROR     = 1001 + iota //创建需求错误
	DEMAND_UPDATE_ERROR     = 1001 + iota //更新需求错误
	DEMAND_QUERY_ERROR      = 1001 + iota //查询需求错误
	DEMAND_DELETE_ERROR     = 1001 + iota //删除需求错误
	SKILL_CREATE_ERROR      = 1001 + iota //创建需求错误
	SKILL_UPDATE_ERROR      = 1001 + iota //更新需求错误
	SKILL_QUERY_ERROR       = 1001 + iota //查询需求错误
	SKILL_DELETE_ERROR      = 1001 + iota //删除需求错误
	BOTTLE_GET_INFOS_ERROR  = 1001 + iota //获取瓶子列表错误
	BOTTLE_THROW_ERROR      = 1001 + iota //扔瓶子错误
	BOTTLE_PICK_ERROR       = 1001 + iota //捡瓶子错误
	BOTTLE_DROP_ERROR       = 1001 + iota //丢弃瓶子错误
	MSG_REPLAY_ERROR       = 1001 + iota //回复瓶子消息错误
	MSG_DELETE_ERROR       = 1001 + iota //删除瓶子消息错误
	MSG_LIST_ERROR       = 1001 + iota //拉取瓶子消息错误
)

type BaseController struct {
	mylog.LogInfo
	ReqParams   []byte
	ServerToken string
	Resp        Response
}

type Response struct {
	Code      int         `json:"code"`      //错误码
	Msg       string      `json:"msg"`       //错误信息
	RequestID string      `json:"requestId"` //请求id
	Data      interface{} `json:"data"`      //返回数据
}

func (bc *BaseController) Prepare(c *gin.Context) {
	//设置requestid
	bc.SetRequestId()
	//设置请求url
	bc.SetRequestUrl(c.Request.RequestURI)
	//设置返回requestid
	bc.Resp.RequestID = bc.GetRequestId()
	//获取请求参数
	bc.ReqParams, _ = ioutil.ReadAll(c.Request.Body)

	mylog.Info("requestId:%s, requestUrl:%s, params : %s", bc.GetRequestId(), bc.GetRequestUrl(), string(bc.ReqParams))
}

func (bc *BaseController) FinishResponse(c *gin.Context) {
	if len(bc.Resp.Msg) <= 0 {
		bc.Resp.Msg = "success"
	}
	c.JSON(200,
		gin.H{
			"errcode":      bc.Resp.Code,
			"errmsg":       bc.Resp.Msg,
			"requestId": bc.Resp.RequestID,
			"data":      bc.Resp.Data,
		})
	r, _ := json.Marshal(bc.Resp)
	mylog.Info("requestUrl:%s, response data:%s", bc.GetRequestUrl(), string(r))
}

func (bc *BaseController) CheckToken(userID, tokenData string) (err error) {
	if len(bc.ServerToken) == 0 {
		err = bc.userCheck(userID, tokenData)
	} else {
		err = bc.serverCheck()
	}
	return
}

func (bc *BaseController) userCheck(userID, tokenData string) error {
	if claim, err := token.CheckUserToken(tokenData); err != nil {
		bc.Resp.Code = TOKEN_CHECK_ERROR
		bc.Resp.Msg = "token check failed!"
		return fmt.Errorf("token check failed!")
	} else if claim.UserData.UserID != userID {
		bc.Resp.Code = USER_CHECK_ERROR
		bc.Resp.Msg = "user is invilid!"
		return fmt.Errorf("user is invalid!")
	}
	return nil
}

func (bc *BaseController) serverCheck() error {
	if _, err := token.CheckServerToken(bc.ServerToken); err != nil {
		bc.Resp.Code = TOKEN_CHECK_ERROR
		bc.Resp.Msg = "token check failed!"
		return fmt.Errorf("token check failed!")
	}
	return nil
}
