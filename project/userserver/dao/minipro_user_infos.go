package dao

import (
	"context"
	"github.com/kukayyou/commonlib/mylog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"userserver/config"
)

type MiniUserInfo struct {
	ID         primitive.ObjectID `bson:"_id"`
	OpenID     string             `json:"openId"`     //微信登录返回的openid
	SessionKey string             `json:"sessionKey"` //微信登录返回的session_key
	Gender     int64              `json:"gender"`     //性别，0：未知， 1：男，2：女
	City       string             `json:"city"`       //用户所在城市
	Token      string             `json:"token"`
}

func (mu *MiniUserInfo) Create(requestID string) (*MiniUserInfo, error) {
	collection := config.Client.Database("follow_bottle").Collection("user")
	_, err := collection.InsertOne(context.TODO(), mu)
	if err != nil {
		mylog.Error("requestID:%s, Create mini user info failed. error:%s", requestID, err.Error())
		return nil, err
	}
	return mu, nil
}

func (mu *MiniUserInfo) Update(requestID string) (*MiniUserInfo, error) {
	collection := config.Client.Database("follow").Collection("user")
	updateOpts := options.Update().SetUpsert(true)
	update := bson.M{"$set": mu}
	_, err := collection.UpdateOne(context.TODO(), bson.M{"_id": mu.ID}, update, updateOpts)
	if err != nil {
		mylog.Error("requestID:%s, Update mini user info failed. error:%s", requestID, err.Error())
		return nil, err
	}
	return mu, nil
}
