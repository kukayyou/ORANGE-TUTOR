package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/kukayyou/commonlib/mylog"
	"orderserver/dao"
)

type DeleteDemandInfoRequest struct {
	UserID    int64   `json:"userId"`    //创建人用户id
	Token     string  `json:"token"`     //token
	DemandIDs []int64 `json:"demandIds"` //需求id
}

func (this DemandInfoController) DeleteDemandInfosApi(c *gin.Context) {
	defer this.FinishResponse(c)
	this.Prepare(c)

	var (
		params DeleteDemandInfoRequest
	)
	if err := json.Unmarshal(this.ReqParams, &params); err != nil {
		mylog.Error("requestID:%s, Unmarshal request failed!", this.GetRequestId())
		this.Resp.Code = PARAMS_PARSE_ERROR
		this.Resp.Msg = "Unmarshal request failed!"
		return
	}
	if err := this.userCheck(params.UserID, params.Token); err != nil {
		mylog.Error("requestID:%s, UserCheck error:%s", this.GetRequestId(), err.Error())
		return
	}
	demandInfo := dao.DemandInfo{}
	if err := demandInfo.DeleteDemandInfo(this.GetRequestId(), &params.DemandIDs); err != nil {
		this.Resp.Code = DEMAND_DELETE_ERROR
		this.Resp.Msg = "delete skill Info failed!"
		return
	}

	return
}
