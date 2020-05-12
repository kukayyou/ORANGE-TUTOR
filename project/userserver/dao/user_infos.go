package dao

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/kukayyou/commonlib/mylog"
	"github.com/kukayyou/commonlib/mytypeconv"
)

type UserInfo struct {
	UserID     int64  `json:"userId"`
	UserName   string `json:"userName"`
	Passwd     string `json:"passwd"`
	Mobile     string `json:"mobile"`
	Email      string `json:"email"`
	Sex        string `json:"sex"` //male or female
	Age        uint64 `json:"age"`
	University string `json:"university"`
	Token      string `json:"token"`
}

func (uc UserInfo) GetUserInfo(requestID string) (userInfo UserInfo, err error) {
	sql := `SELECT user_id,user_name,passwd,mobile,sex,age,email  FROM user WHERE user_id = %d`
	sql = fmt.Sprintf(sql, uc.UserID)
	mylog.Info("requestID %s, sql:%s", requestID, sql)

	var (
		result []orm.Params
	)

	o := orm.NewOrm()
	o.Using("default")
	if num, err := o.Raw(sql).Values(&result); err != nil {
		mylog.Error("requestID:%s, GetUserInfo return error:%s", requestID, err.Error())
		return userInfo, err
	} else if num <= 0 {
		mylog.Error("requestID:%s, GetUserInfo result is null", requestID)
		return userInfo, fmt.Errorf(("GetUserInfo result is null"))
	}
	mylog.Info("requestID:%s, GetUserInfo result:%v", requestID, result)
	userInfo = parseResult(result[0])

	return userInfo, nil
}

func (uc UserInfo) UpdateUserInfo(requestID string) error {
	sql := `INSERT INTO user (user_id,user_name,passwd,mobile,sex,age,email) VALUES (%d,'%s','%s','%s','%s',%d,'%s') ON DUPLICATE KEY UPDATE user_name='%s',passwd='%s',mobile='%s',sex='%s',age=%d,email='%s'`
	sql = fmt.Sprintf(sql, uc.UserID,
		mytypeconv.MysqlEscapeString(uc.UserName),
		uc.Passwd,
		uc.Mobile,
		uc.Sex,
		uc.Age,
		mytypeconv.MysqlEscapeString(uc.Email),
		mytypeconv.MysqlEscapeString(uc.UserName),
		uc.Passwd,
		uc.Mobile,
		uc.Sex,
		uc.Age,
		mytypeconv.MysqlEscapeString(uc.Email))
	mylog.Info("requestID %s, sql:%s", requestID, sql)

	o := orm.NewOrm()
	o.Using("default")
	if _, err := o.Raw(sql).Exec(); err != nil {
		mylog.Error("requestID:%s, getLastUserID return error:%s", requestID, err.Error())
		return err
	}
	mylog.Info("requestID:%s, getLastUserID result:%v", requestID)

	return nil
}

func (uc UserInfo) RegisterUserInfo(requestID string) (userInfo UserInfo, err error) {

	if uc.UserID, err = getLastUserID(requestID); err != nil {
		mylog.Error("requestID:%s, create userId failed, error:%s", requestID, err.Error())
		return uc, err
	}
	if err = uc.UpdateUserInfo(requestID); err != nil {
		mylog.Error("requestID:%s, create userInfo failed, error:%s", requestID, err.Error())
		return uc, err
	}

	return uc, nil
}

func getLastUserID(requestID string) (int64, error) {
	sql := `SELECT user_id FROM user order by user_id DESC limit 1`
	mylog.Info("requestID %s, sql:%s", requestID, sql)

	var (
		result orm.ParamsList
	)

	o := orm.NewOrm()
	if num, err := o.Raw(sql).ValuesFlat(&result); err != nil {
		mylog.Error("requestID:%s, getLastUserID return error:%s", requestID, err.Error())
		return 0, err
	} else if num <= 0 {
		mylog.Error("requestID:%s, getLastUserID result is null", requestID)
		return 0, fmt.Errorf(("getLastUserID result is null"))
	}
	mylog.Info("requestID:%s, getLastUserID result:%v", requestID, result)

	return mytypeconv.ToInt64(result[0], 0) + 1, nil
}

func parseResult(data orm.Params) (userInfo UserInfo) {
	userInfo.UserID = mytypeconv.ToInt64(data["user_id"], 0)
	userInfo.UserName = mytypeconv.ToString(data["user_name"])
	userInfo.Passwd = mytypeconv.ToString(data["passwd"])
	userInfo.Email = mytypeconv.ToString(data["email"])
	userInfo.Mobile = mytypeconv.ToString(data["mobile"])
	userInfo.Sex = mytypeconv.ToString(data["sex"])
	userInfo.Age, _ = mytypeconv.ToUint64(data["age"])
	return
}
