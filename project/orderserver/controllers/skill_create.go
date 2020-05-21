package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/kukayyou/commonlib/mylog"
	"orderserver/dao"
)

type SkillInfoController struct {
	BaseController
}

type CreateSkillRequest struct {
	Token string `json:"token"` //token
	dao.SkillInfo
}

func (this SkillInfoController) CreateSkillApi(c *gin.Context) {
	defer this.FinishResponse(c)
	this.Prepare(c)

	var (
		params CreateSkillRequest
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
	skillInfo := dao.SkillInfo{
		UserID:    params.UserID,
		SkillName: params.SkillName,
		Desc:      params.Desc,
	}
	if skillInfo, err = skillInfo.CreateOrUpdateSkillInfo(this.GetRequestId()); err != nil {
		this.Resp.Code = SKILL_CREATE_ERROR
		this.Resp.Msg = "create skill info failed!"
		return
	} else {
		this.Resp.Data = skillInfo
	}
	return
}
