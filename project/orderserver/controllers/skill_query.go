package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/kukayyou/commonlib/mylog"
	"orderserver/dao"
)

type QuerySkillInfosRequest struct {
	UserID int64  `json:"userId"` //创建人用户id
	Token  string `json:"token"`  //token
}

func (this SkillInfoController) GetSkillInfosApi(c *gin.Context) {
	defer this.FinishResponse(c)
	this.Prepare(c)

	var (
		params QuerySkillInfosRequest
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
	skillInfo := dao.SkillInfo{UserID: params.UserID}
	if skillInfos, err := skillInfo.GetSikllInfosByUserID(this.GetRequestId()); err != nil {
		this.Resp.Code = SKILL_QUERY_ERROR
		this.Resp.Msg = "get skill info failed!"
		return
	} else {
		this.Resp.Data = skillInfos
	}
	return
}
