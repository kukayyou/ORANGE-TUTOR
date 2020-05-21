package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/kukayyou/commonlib/mylog"
	"github.com/kukayyou/commonlib/mytypeconv"
	"orderserver/dao"
)

type UpdateSkillInfoRequest struct {
	UserID     int64       `json:"userId"`  //创建人用户id
	Token      string      `json:"token"`   //token
	SkillID    int64       `json:"skillId"` //需求id
	UpdateData interface{} `json:"updateData"`
}

func (this SkillInfoController) UpdateSkillInfosApi(c *gin.Context) {
	defer this.FinishResponse(c)
	this.Prepare(c)

	var (
		params UpdateSkillInfoRequest
	)
	if err := json.Unmarshal(this.ReqParams, &params); err != nil {
		mylog.Error("requestID:%s, Unmarshal request failed!", this.GetRequestId())
		this.Resp.Code = PARAMS_PARSE_ERROR
		this.Resp.Msg = "Unmarshal request failed!"
		return
	}
	if err := this.UserCheck(params.UserID, params.Token); err != nil {
		mylog.Error("requestID:%s, UserCheck error:%s", this.GetRequestId(), err.Error())
		return
	}
	newSkillInfo := dao.SkillInfo{
		SkillID: params.SkillID,
		UserID:  params.UserID}
	if oldSkillInfo, err := newSkillInfo.GetSkillInfoBySkillID(this.GetRequestId()); err != nil {
		this.Resp.Code = SKILL_QUERY_ERROR
		this.Resp.Msg = "get skill Info by SkillID failed!"
		return
	} else {
		updateData := mytypeconv.InterfaceToMap(params.UpdateData)
		for k, v := range updateData {
			switch k {
			case "skillName":
				oldSkillInfo.SkillName = mytypeconv.ToString(v)
			case "desc":
				oldSkillInfo.Desc = mytypeconv.ToString(v)
			}
		}
		newSkillInfo = oldSkillInfo
	}
	if skillInfo, err := newSkillInfo.CreateOrUpdateSkillInfo(this.GetRequestId()); err != nil {
		this.Resp.Code = SKILL_UPDATE_ERROR
		this.Resp.Msg = "update skill info failed!"
		return
	} else {
		this.Resp.Data = skillInfo
	}
	return
}
