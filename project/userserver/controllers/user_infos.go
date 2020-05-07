package controllers

import (
	"github.com/gin-gonic/gin"
)

type UserInfo struct {
	UserID   uint64 `json:"userId"`
	UserName string `json:"userName"`
	Mobile   string `json:"mobile"`
	Email    string `json:"email"`
	Sex      string `json:"sex"` //male or female
	Age      uint64 `json:"age"`
}

func GetUserInfosApi(c *gin.Context) {
	var userInfo UserInfo
	data := userInfo.GetUserInfo()

	c.JSON(200,
		gin.H{
			"status": "1",
			"data":   data,
		})
	return
}

func (uc *UserInfo) GetUserInfo() []UserInfo {
	var userInfo UserInfo
	return []UserInfo{userInfo}
}
