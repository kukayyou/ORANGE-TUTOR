package dao

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/kukayyou/commonlib/mylog"
	"github.com/kukayyou/commonlib/mytypeconv"
)

type UserInfo struct {
	UserID     int64  `json:"userId"`
	UserType   int64  `json:"userType"`   //0：大学生家教，1：家长
	UserName   string `json:"userName"`   //不支持修改
	NickName   string `json:"nickName"`   //昵称
	Passwd     string `json:"passwd"`     //密码
	Mobile     string `json:"mobile"`     //手机号
	Email      string `json:"email"`      //邮箱
	Sex        string `json:"sex"`        //male or female,不支持修改
	Age        uint64 `json:"age"`        //年龄
	University string `json:"university"` //学校，不支持修改
	Location   string `json:"location"`   //定位
	ProfilePic string `json:"profilePic"` //头像
	Token      string `json:"token"`
}

const USERSELECTPARAMS = "user_id,user_type,user_name,passwd,mobile,sex,age,email,university,nick_name,location,profile_pic"

func (uc UserInfo) GetUserInfo(requestID string) (userInfo UserInfo, err error) {
	sql := `SELECT %s FROM user WHERE user_id = %d`
	sql = fmt.Sprintf(sql, USERSELECTPARAMS, uc.UserID)

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
	userInfos := parseResult(result)

	return userInfos[0], nil
}

func (uc UserInfo) GetUserInfoByUserIDs(requestID string, userIDs []int64) (userInfos []UserInfo, err error) {
	var sql string
	if len(userIDs) > 0 {
		sql = `SELECT %s FROM user WHERE user_id IN ('%s')`
		sql = fmt.Sprintf(sql, USERSELECTPARAMS, mytypeconv.JoinInt64Array(userIDs, "','"))
	} else {
		return nil, fmt.Errorf("UserIDs is nil")
	}

	mylog.Info("requestID %s, sql:%s", requestID, sql)

	var (
		result []orm.Params
	)

	o := orm.NewOrm()
	o.Using("default")
	if num, err := o.Raw(sql).Values(&result); err != nil {
		mylog.Error("requestID:%s, GetUserInfo return error:%s", requestID, err.Error())
		return nil, err
	} else if num <= 0 {
		mylog.Error("requestID:%s, GetUserInfo result is null", requestID)
		return nil, fmt.Errorf(("GetUserInfo result is null"))
	}
	mylog.Info("requestID:%s, GetUserInfo result:%v", requestID, result)
	userInfos = parseResult(result)

	return userInfos, nil
}

func (uc UserInfo) CreateOrUpdateUserInfo(requestID string) error {
	sql := `INSERT INTO user (%s) VALUES (%d,%d,'%s','%s','%s','%s',%d,'%s','%s','%s','%s','%s') ON DUPLICATE KEY UPDATE user_type=%d,user_name='%s',passwd='%s',mobile='%s',sex='%s',age=%d,email='%s',university='%s',nick_name='%s',location='%s',profile_pic='%s'`
	sql = fmt.Sprintf(sql,
		USERSELECTPARAMS,
		uc.UserID,
		uc.UserType,
		mytypeconv.MysqlEscapeString(uc.UserName),
		uc.Passwd,
		uc.Mobile,
		uc.Sex,
		uc.Age,
		mytypeconv.MysqlEscapeString(uc.Email),
		mytypeconv.MysqlEscapeString(uc.University),
		uc.NickName,
		uc.Location,
		uc.ProfilePic,
		uc.UserType,
		mytypeconv.MysqlEscapeString(uc.UserName),
		uc.Passwd,
		uc.Mobile,
		uc.Sex,
		uc.Age,
		mytypeconv.MysqlEscapeString(uc.Email),
		mytypeconv.MysqlEscapeString(uc.University),
		uc.NickName,
		uc.Location,
		uc.ProfilePic)
	mylog.Info("requestID %s, sql:%s", requestID, sql)

	o := orm.NewOrm()
	o.Using("default")
	if _, err := o.Raw(sql).Exec(); err != nil {
		mylog.Error("requestID:%s, CreateOrUpdateUserInfo return error:%s", requestID, err.Error())
		return err
	}
	mylog.Info("requestID:%s, CreateOrUpdateUserInfo result:%v", requestID, uc)

	return nil
}

func (uc UserInfo) RegisterUserInfo(requestID string) (userInfo UserInfo, err error) {
	if uc.UserID, err = getLastUserID(requestID); err != nil {
		mylog.Error("requestID:%s, create userId failed, error:%s", requestID, err.Error())
		return uc, err
	}
	if err = uc.CreateOrUpdateUserInfo(requestID); err != nil {
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
		userID int64
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

	if mytypeconv.ToInt64(result[0], 0) == 0 {
		userID = mytypeconv.ToInt64(result[0], 0) + 1000001
	} else {
		userID = mytypeconv.ToInt64(result[0], 0) + 1
	}
	return userID, nil
}

func parseResult(data []orm.Params) (userInfos []UserInfo) {
	userInfos = make([]UserInfo, 0)
	for _, value := range data {
		userInfo := UserInfo{
			UserID:     mytypeconv.ToInt64(value["user_id"], 0),
			UserType:   mytypeconv.ToInt64(value["user_type"], 0),
			UserName:   mytypeconv.ToString(value["user_name"]),
			Passwd:     mytypeconv.ToString(value["passwd"]),
			Email:      mytypeconv.ToString(value["email"]),
			Mobile:     mytypeconv.ToString(value["mobile"]),
			Sex:        mytypeconv.ToString(value["sex"]),
			Age:        mytypeconv.ToUint64(value["age"]),
			University: mytypeconv.ToString(value["university"]),
			NickName:   mytypeconv.ToString(value["nick_name"]),
			Location:   mytypeconv.ToString(value["location"]),
			ProfilePic: mytypeconv.ToString(value["profile_pic"]),
		}
		userInfos = append(userInfos, userInfo)
	}
	return
}
