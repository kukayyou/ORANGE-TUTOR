package dao

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/kukayyou/commonlib/mylog"
	"github.com/kukayyou/commonlib/mytypeconv"
)

//用于管理家长发布的需求
type DemandInfo struct {
	UserID       int64  `json:"userId"`       //创建人用户id
	DemandID     int64  `json:"demandId"`     //需求id
	Title        string `json:"title"`        //需求名称
	SubjectType  string `json:"subjectType"`  //学科类别
	ClassLevel   int64  `json:"classLevel"`   //年级
	DemandStatus int64  `json:"demandStatus"` //需求状态，0：关闭，1：启用
	Desc         string `json:"desc"`         //需求描述
}

const DEMANDSELECTPARAMS = "id,user_id,title,subject_type,class_level,description,demand_status"

func (demand *DemandInfo) GetDemandInfosByUserID(requestID string) (*[]DemandInfo, error) {
	sql := `SELECT %s FROM user_demand WHERE user_id = %d ORDER BY id DESC`
	sql = fmt.Sprintf(sql, DEMANDSELECTPARAMS, demand.UserID)
	mylog.Info("requestID %s, sql:%s", requestID, sql)

	var (
		result []orm.Params
	)

	o := orm.NewOrm()
	o.Using("default")
	if num, err := o.Raw(sql).Values(&result); err != nil {
		mylog.Error("requestID:%s, GetDemandInfo return error:%s", requestID, err.Error())
		return nil, err
	} else if num <= 0 {
		mylog.Error("requestID:%s, GetDemandInfo result is null", requestID)
		return nil, fmt.Errorf(("GetDemandInfo result is null"))
	}
	mylog.Info("requestID:%s, GetDemandInfo result:%v", requestID, result)
	demandInfos := parseDemandResult(&result)

	return &demandInfos, nil
}

func (demand *DemandInfo) GetDemandInfoByDemandID(requestID string) (*DemandInfo, error) {
	sql := `SELECT %s FROM user_demand WHERE user_id=%d AND id = %d`
	sql = fmt.Sprintf(sql, DEMANDSELECTPARAMS, demand.UserID, demand.DemandID)
	mylog.Info("requestID %s, sql:%s", requestID, sql)

	var (
		result []orm.Params
	)

	o := orm.NewOrm()
	o.Using("default")
	if num, err := o.Raw(sql).Values(&result); err != nil {
		mylog.Error("requestID:%s, GetDemandInfo return error:%s", requestID, err.Error())
		return demand, err
	} else if num <= 0 {
		mylog.Error("requestID:%s, GetDemandInfo result is null", requestID)
		return demand, fmt.Errorf(("GetDemandInfo result is null"))
	}
	mylog.Info("requestID:%s, GetDemandInfo result:%v", requestID, result)
	demandInfos := parseDemandResult(&result)

	if len(demandInfos) > 0 {
		return &demandInfos[0], nil
	} else {
		return nil, fmt.Errorf("demandInfos is nil")
	}
}

func (demand *DemandInfo) CreateOrUpdateDemandInfo(requestID string) (*DemandInfo, error) {
	if demand.DemandID == 0 {
		demand.DemandID, _ = getLastDemandID(requestID)
	}
	sql := `INSERT INTO user_demand (%s) VALUES (%d,%d,'%s','%s',%d,'%s',%d) ON DUPLICATE KEY UPDATE title='%s',subject_type='%s',class_level=%d,description='%s',demand_status=%d`
	sql = fmt.Sprintf(sql,
		DEMANDSELECTPARAMS,
		demand.DemandID,
		demand.UserID,
		mytypeconv.MysqlEscapeString(demand.Title),
		mytypeconv.MysqlEscapeString(demand.SubjectType),
		demand.ClassLevel,
		mytypeconv.MysqlEscapeString(demand.Desc),
		1,
		mytypeconv.MysqlEscapeString(demand.Title),
		mytypeconv.MysqlEscapeString(demand.SubjectType),
		demand.ClassLevel,
		mytypeconv.MysqlEscapeString(demand.Desc),
		1,
	)
	mylog.Info("requestID:%s, sql:%s", requestID, sql)

	o := orm.NewOrm()
	o.Using("default")
	if re, err := o.Raw(sql).Exec(); err != nil {
		mylog.Error("requestID:%s, CreateOrUpdateDemandInfo return error:%s", requestID, err.Error())
		return demand, err
	} else {
		demand.DemandID, _ = re.LastInsertId()
	}

	mylog.Info("requestID:%s, CreateOrUpdateDemandInfo result:%v", requestID, demand)

	return demand, nil
}

func (demand *DemandInfo) DeleteDemandInfo(requestID string, demandIDs *[]int64) error {
	sql := `DELETE FROM user_demand where id  IN ('%s')`
	sql = fmt.Sprintf(sql, mytypeconv.JoinInt64Array(*demandIDs, "','"))
	mylog.Info("requestID %s, sql:%s", requestID, sql)

	o := orm.NewOrm()
	o.Using("default")
	if _, err := o.Raw(sql).Exec(); err != nil {
		mylog.Error("requestID:%s, DeleteDemandInfo return error:%s", requestID, err.Error())
		return err
	}
	mylog.Info("requestID:%s, DeleteDemandInfo success", requestID)

	return nil
}

func parseDemandResult(data *[]orm.Params) []DemandInfo {
	demandInfos := make([]DemandInfo, 0)
	for _, value := range *data {
		demand := DemandInfo{
			DemandID:     mytypeconv.ToInt64(value["id"], 0),
			UserID:       mytypeconv.ToInt64(value["user_id"], 0),
			Title:        mytypeconv.ToString(value["title"]),
			SubjectType:  mytypeconv.ToString(value["subject_type"]),
			ClassLevel:   mytypeconv.ToInt64(value["class_level"], 0),
			Desc:         mytypeconv.ToString(value["description"]),
			DemandStatus: mytypeconv.ToInt64(value["demand_status"], 0),
		}
		demandInfos = append(demandInfos, demand)
	}

	return demandInfos
}

func getLastDemandID(requestID string) (int64, error) {
	sql := `SELECT id FROM user_demand order by id DESC limit 1`
	mylog.Info("requestID %s, sql:%s", requestID, sql)

	var (
		result orm.ParamsList
	)

	o := orm.NewOrm()
	if num, err := o.Raw(sql).ValuesFlat(&result); err != nil {
		mylog.Error("requestID:%s, getLastDemandID return error:%s", requestID, err.Error())
		return 0, err
	} else if num <= 0 {
		mylog.Error("requestID:%s, getLastDemandID result is null", requestID)
		return 0, fmt.Errorf(("getLastDemandID result is null"))
	}
	mylog.Info("requestID:%s, getLastDemandID result:%v", requestID, result)

	return mytypeconv.ToInt64(result[0], 0) + 1, nil
}
