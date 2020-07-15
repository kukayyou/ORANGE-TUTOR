package dao

import (
	"context"
	"github.com/kukayyou/commonlib/mylog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"msgserver/config"
)

type BottleInfo struct {
	BottleID       primitive.ObjectID `bson:"_id"`            //瓶子id
	SourUserID     string             `json:"sourUserId"`     //源用户id
	SourUserGender int64              `json:"sourUserGender"` //源用户性别
	DestUserID     string             `json:"destUserId"`     //目的用户id
	DestUserGender int64              `json:"destUserGender"` //目的用户性别
	Content        []BottleContent    `json:"bottleContent"`  //瓶子内容
	City           string             `json:"city"`           //位置信息，格式：【经度，维度】
	DeleteUser     []string           `json:"deleteUser"`     //删除消息的用户
	IsNewMsg       int64              `json:"isNewMsg"`       //是否有新消息, 0:无，1：有
}

type BottleContent struct {
	ContentType int64    `json:"contentType"` //瓶子内容类型，1:源用户内容，2：目的用户内容
	Text        string   `json:"text"`        //瓶子内容
	Pic         []string `json:"pic"`         //图片
	TimeStamp   int64    `json:"timeStamp"`   //时间戳
}

//查询瓶子信息
func (bc *BottleInfo) GetBottleInfoByID(requestID string) (*BottleInfo, error) {
	collection := config.Client.Database("follow_bottle").Collection("bottle")
	filter := bson.M{"_id": bc.BottleID}
	singReslut := collection.FindOne(context.TODO(), filter)

	var bottleInfo BottleInfo
	if err := singReslut.Decode(&bottleInfo); err != nil {
		mylog.Error("requestID：%s, Unmarshal data error:%s", requestID, err.Error())
		return nil, err
	}
	//更新瓶子消息已读状态
	bottleInfo.UpdateMsgReadStatus(requestID)
	return &bottleInfo, nil
}

//删除瓶子消息
func (bc *BottleInfo) DeleteMsg(requestID string) error {
	delUSerID := make([]string, 0)
	delUSerID = append(delUSerID, bc.DeleteUser...)

	//获取原瓶子信息
	bottleInfo, err := bc.GetBottleInfoByID(requestID)
	if err != nil {
		return err
	}

	//更新瓶子删除人信息
	collection := config.Client.Database("follow_bottle").Collection("bottle")
	filter := bson.M{"_id": bc.BottleID}
	bottleInfo.DeleteUser = append(bottleInfo.DeleteUser, delUSerID...)
	update := bson.M{"$set": bson.M{"deleteuser": bottleInfo.DeleteUser}}

	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		mylog.Error("requestID：%s, Delete Msg error:%s", requestID, err.Error())
	}
	return err
}

func (bc *BottleInfo) ReplayMsg(requestID string) error {
	bottleContent := make([]BottleContent, 0)
	bottleContent = append(bottleContent, bc.Content...)

	//获取原瓶子信息
	bottleInfo, err := bc.GetBottleInfoByID(requestID)
	if err != nil {
		return err
	}

	//更新瓶子消息内容
	collection := config.Client.Database("follow_bottle").Collection("bottle")
	filter := bson.M{"_id": bc.BottleID}
	bottleInfo.Content = append(bottleInfo.Content, bottleContent...)
	update := bson.M{"$set": bson.M{"content": bottleInfo.Content, "isnewmsg": 1}}

	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		mylog.Error("requestID：%s, Replay Msg error:%s", requestID, err.Error())
	}
	return err
}

//更新瓶子消息已读状态
func (bc *BottleInfo) UpdateMsgReadStatus(requestID string) error {
	collection := config.Client.Database("follow_bottle").Collection("bottle")
	filter := bson.M{"_id": bc.BottleID}
	update := bson.M{"$set": bson.M{"isnewmsg": 0}}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		mylog.Error("requestID：%s, Update Msg Read Status error:%s", requestID, err.Error())
	}
	return err
}
