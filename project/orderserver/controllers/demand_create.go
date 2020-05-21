package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/kukayyou/commonlib/mylog"
	"orderserver/dao"
	"time"
)

type DemandInfoController struct {
	BaseController
}

type CreateDemandInfoRequest struct {
	Token string `json:"token"` //token
	dao.DemandInfo
}

func (this DemandInfoController) CreateDemandApi(c *gin.Context) {
	defer this.FinishResponse(c)
	this.Prepare(c)

	var (
		params CreateDemandInfoRequest
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
	demandInfo := dao.DemandInfo{
		UserID:       params.UserID,
		Title:        params.Title,
		SubjectType:  params.SubjectType,
		ClassLevel:   params.ClassLevel,
		DemandStatus: params.DemandStatus,
		Desc:         params.Desc,
	}
	for i := 0; i < 100000; i++ {
		demandInfo.CreateOrUpdateDemandInfo(this.GetRequestId())
		time.Sleep(time.Microsecond * 1)
	}
	if demandInfo, err = demandInfo.CreateOrUpdateDemandInfo(this.GetRequestId()); err != nil {
		this.Resp.Code = DEMAND_CREATE_ERROR
		this.Resp.Msg = "create demand info failed!"
		return
	} else {
		this.Resp.Data = demandInfo
	}
	return
}
