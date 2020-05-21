package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/kukayyou/commonlib/mylog"
	"github.com/kukayyou/commonlib/mytypeconv"
	"userserver/dao"
)

type UpdateUserInfosRequest struct {
	UserID     int64       `json:"UserId"`
	Token      string      `json:"token"`
	UpdateData interface{} `json:"updateData"`
}

func (this UserInfosController) UpdateUserInfoApi(c *gin.Context) {
	defer this.FinishResponse(c)
	this.Prepare(c)
	var params UpdateUserInfosRequest
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
	newUserInfo := dao.UserInfo{UserID: params.UserID}
	if oldUserInfo, err := newUserInfo.GetUserInfo(this.GetRequestId()); err != nil {
		this.Resp.Code = USER_UPDATE_INFOS_ERROR
		this.Resp.Msg = "GetUserInfo failed!"
		return
	} else {
		updateData := mytypeconv.InterfaceToMap(params.UpdateData)
		for k, v := range updateData {
			switch k {
			case "passwd":
				oldUserInfo.Passwd = mytypeconv.ToString(v)
			case "mobile":
				oldUserInfo.Mobile = mytypeconv.ToString(v)
			case "email":
				oldUserInfo.Email = mytypeconv.ToString(v)
			case "age":
				oldUserInfo.Age = mytypeconv.ToUint64(v)
			case "location":
				oldUserInfo.Location = mytypeconv.ToString(v)
			case "profilePic":
				oldUserInfo.ProfilePic = mytypeconv.ToString(v)
			case "nickName":
				oldUserInfo.NickName = mytypeconv.ToString(v)
			}
		}
		newUserInfo = oldUserInfo
	}
	if err := newUserInfo.CreateOrUpdateUserInfo(this.GetRequestId()); err != nil {
		this.Resp.Code = USER_UPDATE_INFOS_ERROR
		this.Resp.Msg = "UpdateUserInfo failed!"
		return
	}
	return
}
