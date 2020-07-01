package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/kukayyou/commonlib/mylog"
	"github.com/kukayyou/commonlib/mytypeconv"
	"orderserver/dao"
)

type UpdateDemandInfoRequest struct {
	UserID     int64       `json:"userId"`   //创建人用户id
	Token      string      `json:"token"`    //token
	DemandID   int64       `json:"demandId"` //需求id
	UpdateData interface{} `json:"updateData"`
}

func (this DemandInfoController) UpdateDemandInfosApi(c *gin.Context) {
	defer this.FinishResponse(c)
	this.Prepare(c)

	var (
		params UpdateDemandInfoRequest
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
	newDemandInfo := &dao.DemandInfo{DemandID: params.DemandID}
	if oldDemandInfo, err := newDemandInfo.GetDemandInfoByDemandID(this.GetRequestId()); err != nil {
		this.Resp.Code = DEMAND_QUERY_ERROR
		this.Resp.Msg = "get demand Info by demandID failed!"
		return
	} else {
		updateData := mytypeconv.InterfaceToMap(params.UpdateData)
		for k, v := range updateData {
			switch k {
			case "title":
				oldDemandInfo.Title = mytypeconv.ToString(v)
			case "subjectType":
				oldDemandInfo.SubjectType = mytypeconv.ToString(v)
			case "classLevel":
				oldDemandInfo.ClassLevel = mytypeconv.ToInt64(v, 0)
			case "demandStatus":
				oldDemandInfo.DemandStatus = mytypeconv.ToInt64(v, 0)
			case "desc":
				oldDemandInfo.Desc = mytypeconv.ToString(v)
			}
		}
		newDemandInfo = oldDemandInfo
	}
	if demandInfo, err := newDemandInfo.CreateOrUpdateDemandInfo(this.GetRequestId()); err != nil {
		this.Resp.Code = DEMAND_UPDATE_ERROR
		this.Resp.Msg = "update demand info failed!"
		return
	} else {
		this.Resp.Data = demandInfo
	}
	return
}
