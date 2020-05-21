package dao

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/kukayyou/commonlib/mylog"
	"github.com/kukayyou/commonlib/mytypeconv"
)

//用于管理教师的技能
type SkillInfo struct {
	UserID    int64  `json:"userId"`
	SkillID   int64  `json:"skillId"`
	SkillName string `json:"skillName"`
	Desc      string `json:"desc"`
}

const SKILLSELECTPARAMS = "id, user_id, skill_name, description"

func (si SkillInfo) GetSikllInfosByUserID(requestID string) ([]SkillInfo, error) {
	sql := `SELECT %s FROM user_skill WHERE user_id = %d ORDER BY id DESC`
	sql = fmt.Sprintf(sql, SKILLSELECTPARAMS, si.UserID)
	mylog.Info("requestID %s, sql:%s", requestID, sql)

	var (
		result []orm.Params
	)

	o := orm.NewOrm()
	o.Using("default")
	if num, err := o.Raw(sql).Values(&result); err != nil {
		mylog.Error("requestID:%s, GetSikllInfosByUserID return error:%s", requestID, err.Error())
		return nil, err
	} else if num <= 0 {
		mylog.Error("requestID:%s, GetDemandInfo result is null", requestID)
		return nil, fmt.Errorf(("GetSikllInfosByUserID result is null"))
	}
	mylog.Info("requestID:%s, GetSikllInfosByUserID result:%v", requestID, result)
	skillInfos := parseSkillResult(result)

	return skillInfos, nil
}

func (si SkillInfo) GetSkillInfoBySkillID(requestID string) (SkillInfo, error) {
	sql := `SELECT %s FROM user_skill WHERE user_id=%d AND id = %d`
	sql = fmt.Sprintf(sql, SKILLSELECTPARAMS, si.UserID, si.SkillID)
	mylog.Info("requestID %s, sql:%s", requestID, sql)

	var (
		result []orm.Params
	)

	o := orm.NewOrm()
	o.Using("default")
	if num, err := o.Raw(sql).Values(&result); err != nil {
		mylog.Error("requestID:%s, GetSkillInfoBySkillID return error:%s", requestID, err.Error())
		return si, err
	} else if num <= 0 {
		mylog.Error("requestID:%s, GetSkillInfoBySkillID result is null", requestID)
		return si, fmt.Errorf(("GetSkillInfoBySkillID result is null"))
	}
	mylog.Info("requestID:%s, GetSkillInfoBySkillID result:%v", requestID, result)
	skillInfos := parseSkillResult(result)

	return skillInfos[0], nil
}

func (si SkillInfo) CreateOrUpdateSkillInfo(requestID string) (SkillInfo, error) {
	if si.SkillID == 0 {
		si.SkillID, _ = getLastSkillID(requestID)
	}
	sql := `INSERT INTO user_skill (%s) VALUES (%d,%d,'%s','%s') ON DUPLICATE KEY UPDATE skill_name='%s',description='%s'`
	sql = fmt.Sprintf(sql,
		SKILLSELECTPARAMS,
		si.SkillID,
		si.UserID,
		mytypeconv.MysqlEscapeString(si.SkillName),
		mytypeconv.MysqlEscapeString(si.Desc),
		mytypeconv.MysqlEscapeString(si.SkillName),
		mytypeconv.MysqlEscapeString(si.Desc),
	)
	mylog.Info("requestID:%s, sql:%s", requestID, sql)

	o := orm.NewOrm()
	o.Using("default")
	if re, err := o.Raw(sql).Exec(); err != nil {
		mylog.Error("requestID:%s, CreateOrUpdateSkillInfo return error:%s", requestID, err.Error())
		return si, err
	} else {
		si.SkillID, _ = re.LastInsertId()
	}

	mylog.Info("requestID:%s, CreateOrUpdateDemandInfo result:%v", requestID, si)

	return si, nil
}

func (si SkillInfo) DeleteSkillInfo(requestID string, skillIDs []int64) error {
	sql := `DELETE FROM user_skill where id  IN ('%s')`
	sql = fmt.Sprintf(sql, mytypeconv.JoinInt64Array(skillIDs, "','"))
	mylog.Info("requestID %s, sql:%s", requestID, sql)

	o := orm.NewOrm()
	o.Using("default")
	if _, err := o.Raw(sql).Exec(); err != nil {
		mylog.Error("requestID:%s, DeleteSkillInfo return error:%s", requestID, err.Error())
		return err
	}
	mylog.Info("requestID:%s, DeleteSkillInfo success", requestID)

	return nil
}

func getLastSkillID(requestID string) (int64, error) {
	sql := `SELECT id FROM user_skill order by id DESC limit 1`
	mylog.Info("requestID %s, sql:%s", requestID, sql)

	var (
		result orm.ParamsList
	)

	o := orm.NewOrm()
	if num, err := o.Raw(sql).ValuesFlat(&result); err != nil {
		mylog.Error("requestID:%s, getLastDemandID return error:%s", requestID, err.Error())
		return 0, err
	} else if num <= 0 {
		mylog.Error("requestID:%s, getLastSkillID result is null", requestID)
		return 0, fmt.Errorf(("getLastSkillID result is null"))
	}
	mylog.Info("requestID:%s, getLastSkillID result:%v", requestID, result)

	return mytypeconv.ToInt64(result[0], 0) + 1, nil
}

func parseSkillResult(data []orm.Params) (skillInfos []SkillInfo) {
	skillInfos = make([]SkillInfo, 0)
	for _, value := range data {
		skill := SkillInfo{
			SkillID:   mytypeconv.ToInt64(value["id"], 0),
			UserID:    mytypeconv.ToInt64(value["user_id"], 0),
			SkillName: mytypeconv.ToString(value["skill_name"]),
			Desc:      mytypeconv.ToString(value["description"]),
		}
		skillInfos = append(skillInfos, skill)
	}

	return
}
