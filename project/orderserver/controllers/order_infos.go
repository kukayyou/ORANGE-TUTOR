package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/kukayyou/commonlib/myhttp"
	"github.com/kukayyou/commonlib/mylog"
	"encoding/json"
	"time"
)

type GetOrderController struct {
	BaseController
}

type RequestData struct {
	Data int `json:"data"`
}

type OrderInfo struct {
	OrderID   string                 `json:"orderId"`
	UserInfos map[string]interface{} `json:"userInfos"`
}

type UserInfo struct {
	UserID   uint64 `json:"userId"`
	UserName string `json:"userName"`
	Mobile   string `json:"mobile"`
	Email    string `json:"email"`
	Sex      string `json:"sex"` //male or female
	Age      uint64 `json:"age"`
}

func (this GetOrderController)GetOrderInfosApi(c *gin.Context) {
	this.Prepare(c)
	var params RequestData
	json.Unmarshal(this.ReqParams, &params)

	var orderInfo OrderInfo
	mylog.Info("requestID:%s:, GetOrderInfosApi start ... ", this.GetRequestId())
	data := orderInfo.GetOrderInfos(this.GetRequestId())
	c.JSON(200,
		gin.H{
			"status": "1",
			"data":   data,
		})

	go func() {
		time.Sleep(time.Second*2)
		mylog.Info("requestID:%s:, 延时日志：%s", this.GetRequestId(), time.Now().Format("2006-01-02 15:04:05"))
	}()

	return
}

func (oi *OrderInfo) GetOrderInfos(requestID string) []OrderInfo {
	var orderInfo OrderInfo
	if userInfos := GetUserInfoByIDs(requestID);userInfos!=nil{
		orderInfo.UserInfos = userInfos
		o,_:= json.Marshal(orderInfo)
		mylog.Info("requestID:%s orderInfo is :%v", requestID, string(o))
		return []OrderInfo{orderInfo}
	}
	return nil
}

func GetUserInfoByIDs(requestID string) map[string]interface{} {
	resp := myhttp.RequestWithHytrix("api.tutor.com.userserver", "/userserver/user/infos", map[string]string{})
	if resp != nil{
		r,_:=json.Marshal(resp)
		mylog.Info("requestID:%s, resp is :%s", requestID, string(r))
		return resp
	}
	return nil
}
