package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/kukayyou/commonlib/myhttp"
	"github.com/kukayyou/commonlib/mylog"
)

type GetOrderController struct {
	BaseController
}

type RequestData struct {
	UserID int64 `json:"userId"`//userid
	Token string `json:"token"`//token
}

func (this GetOrderController)GetOrderInfosApi(c *gin.Context) {
	defer this.FinishResponse(c)
	this.Prepare(c)

	var params RequestData
	if err :=json.Unmarshal(this.ReqParams, &params);err!= nil{
		mylog.Error("requestID:%s, Unmarshal request error:%s", this.GetRequestId(), err.Error())
		this.Resp.Code = PARAMS_PARSE_ERROR
		this.Resp.Msg = "Unmarshal request failed!"
		return
	}

	if err:= this.userCheck(params.UserID, params.Token);err!=nil{
		mylog.Error("requestID:%s, UserCheck error:%s", this.GetRequestId(), err.Error())
		return
	}

	return
}

/*func (oi *OrderInfo) GetOrderInfos(requestID string) []OrderInfo {
	var orderInfo OrderInfo
	if userInfos := GetUserInfoByIDs(requestID);userInfos!=nil{
		orderInfo.UserInfos = userInfos
		o,_:= json.Marshal(orderInfo)
		mylog.Info("requestID:%s orderInfo is :%v", requestID, string(o))
		return []OrderInfo{orderInfo}
	}
	return nil
}*/

func GetUserInfoByIDs(requestID string) map[string]interface{} {
	resp := myhttp.RequestWithHytrix("api.tutor.com.userserver", "/userserver/user/infos", map[string]string{})
	if resp != nil{
		r,_:=json.Marshal(resp)
		mylog.Info("requestID:%s, resp is :%s", requestID, string(r))
		return resp
	}
	return nil
}
