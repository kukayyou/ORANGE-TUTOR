package config

import (
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ConsulAddress    string //consul地址
	EtcdAddress      string //etcd地址
	LogPath          string //日志地址
	LogLevel         int    //日志级别
	LogMaxAge        int    //日志保留时长
	LogMaxSize       int    //日志最大size
	LogMaxBackups    int    //日志最多留备份数
	MongoReplicas    string //mongdb地址
	MongoMaxPoolSize int64  //最大连接池数
	MongoMinPoolSize int64  //最小连接池数
	WxDomain         string //微信域名
	AppID            string //微信小程序appid
	AppSecret        string //微信小程序appsecret
)

var (
	// Pool MongoDB连接池
	Pool   *mongo.Client
	Client *mongo.Client
)
