package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/kukayyou/commonlib/mylog"
	"orderserver/dao"
)

type QueryDemandInfosRequest struct {
	UserID int64  `json:"userId"` //创建人用户id
	Token  string `json:"token"`  //token
}

func (this DemandInfoController) GetDemandInfosApi(c *gin.Context) {
	defer this.FinishResponse(c)
	this.Prepare(c)

	var (
		params QueryDemandInfosRequest
		err    error
	)
	if err = json.Unmarshal(this.ReqParams, &params); err != nil {
		mylog.Error("requestID:%s, Unmarshal request failed!", this.GetRequestId())
		this.Resp.Code = PARAMS_PARSE_ERROR
		this.Resp.Msg = "Unmarshal request failed!"
		return
	}
	if err := this.UserCheck(params.UserID, params.Token); err != nil {
		mylog.Error("requestID:%s, UserCheck error:%s", this.GetRequestId(), err.Error())
		return
	}
	demandInfo := dao.DemandInfo{UserID: params.UserID}
	if demandInfos, err := demandInfo.GetDemandInfosByUserID(this.GetRequestId()); err != nil {
		this.Resp.Code = DEMAND_QUERY_ERROR
		this.Resp.Msg = "get demand info failed!"
		return
	} else {
		this.Resp.Data = demandInfos
	}
	return
}
