package dao

import (
	"bottleserver/config"
	"context"
	"github.com/kukayyou/commonlib/mylog"
	"github.com/kukayyou/commonlib/mytypeconv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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

//查询用户历史瓶子列表，以降序按id排列
func (bc *BottleInfo) GetUserHistoryBottleList(requestID string) ([]BottleInfo, error) {
	bottles := make([]BottleInfo, 0)
	collection := config.Client.Database("follow_bottle").Collection("bottle")
	filter := bson.M{"$or": []bson.M{bson.M{"souruserid": bc.SourUserID}, bson.M{"destuserid": bc.DestUserID}}}
	findOpt := options.Find().SetSort(bson.D{{"_id", 1}}) //sort 1:升序, -1:降序
	cursor, err := collection.Find(context.TODO(), filter, findOpt)

	if err != nil {
		mylog.Error("requestID：%s, GetHistoryBottleList error:%s", requestID, err.Error())
		return nil, err
	}
	for cursor.Next(context.TODO()) {
		var bottleInfo BottleInfo
		if err := cursor.Decode(&bottleInfo); err == nil {
			if len(bottleInfo.DestUserID) > 0 && //只获取有回复的瓶子且用户未删除的
				!mytypeconv.Contains(bottleInfo.DeleteUser, bc.SourUserID) {
				bottles = append(bottles, bottleInfo)
			}
		} else {
			mylog.Error("requestID：%s, Unmarshal data error:%s", requestID, err.Error())
		}
	}
	return bottles, nil
}

//捡瓶子
func (bc *BottleInfo) PickBottle(requestID string) (*BottleInfo, error) {
	collection := config.Client.Database("follow_bottle").Collection("bottle")
	filter := bson.M{
		"destuserid":     "",                               //查询未被捞走
		"sourusergender": bson.M{"$ne": bc.DestUserGender}, //查询性别不同的
		"city":           bc.City,                          //查询同城的
		"souruserid":     bson.M{"$ne": bc.DestUserID},     //查询不是自己扔出的
	}
	singResult := collection.FindOne(context.TODO(), filter)

	var bottleInfo BottleInfo
	if err := singResult.Decode(&bottleInfo); err != nil {
		mylog.Error("requestID：%s, Unmarshal data error:%s", requestID, err.Error())
		return nil, err
	}
	bottleInfo.DestUserID = bc.DestUserID
	bottleInfo.DestUserGender = bc.DestUserGender
	return &bottleInfo, nil
}

//更新瓶子接收人及其性别
func (bc *BottleInfo) UpdatePickedBottle(requestID string) error {
	collection := config.Client.Database("follow_bottle").Collection("bottle")
	filter := bson.M{"_id": bc.BottleID}
	update := bson.M{"$set": bson.M{"destuserid": bc.DestUserID, "destusergender": bc.DestUserGender}}
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		mylog.Error("requestID：%s, Update Picked Bottle error:%s", requestID, err.Error())
	}
	return err
}

//扔瓶子
func (bc *BottleInfo) ThrowBottle(requestID string) error {
	collection := config.Client.Database("follow_bottle").Collection("bottle")
	if _, err := collection.InsertOne(context.TODO(), bc); err != nil {
		mylog.Error("requestID：%s, InsertOne data error:%s", requestID, err.Error())
		return err
	}
	return nil
}

//丢弃捡到的瓶子
func (bc *BottleInfo) DropBottle(requestID string) error {
	collection := config.Client.Database("follow_bottle").Collection("bottle")
	filter := bson.M{"_id": bc.BottleID}
	update := bson.M{"$set": bson.M{"destuserid": "", "destusergender": 0}}
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		mylog.Error("requestID：%s, Drop Bottle error:%s", requestID, err.Error())
	}
	return err
}
