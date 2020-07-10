package main

import (
	"fmt"
	"github.com/kukayyou/commonlib/myconfig"
	"github.com/kukayyou/commonlib/myhttp"
	"github.com/kukayyou/commonlib/mylog"
	"github.com/kukayyou/commonlib/mysql"
	"github.com/kukayyou/commonlib/token"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/etcd"
	"github.com/micro/go-micro/web"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"userserver/config"
	//"github.com/micro/go-plugins/registry/consul"
	"userserver/routers"
	"context"
)

func main() {

	defer mylog.SugarLogger.Sync()

	if err:= initMySQL();err != nil{
		return
	}
	//初始化，mongdb
	if err:= InitMongodb();err != nil{
		return
	}
	//初始化路由
	ginRouter := routers.InitRouters()
	//新建一个consul注册的地址，也就是我们consul服务启动的机器ip+端口
	/*consulReg := consul.NewRegistry(
		registry.Addrs("192.168.109.131:8500"),
	)*/
	etcdReg := etcd.NewRegistry(
		registry.Addrs(config.EtcdAddress))
	//初始化go-micro熔断地址
	myhttp.EtcdAddr = config.EtcdAddress
	myhttp.ConsulAddr = config.ConsulAddress
	//注册服务
	microService:= web.NewService(
		web.Name("api.tutor.com.userserver"),
		//web.RegisterTTL(time.Second*30),//设置注册服务的过期时间
		//web.RegisterInterval(time.Second*20),//设置间隔多久再次注册服务
		web.Address(":18001"),
		web.Handler(ginRouter),
		web.Registry(etcdReg),
		)

	if err := microService.Run();err != nil{
		mylog.Error("server run error:%s",err.Error())
	}
}

func init() {
	//加载配置文件
	initConfig()
	//初始化日志
	mylog.InitLog(config.LogPath,"userserver", config.LogMaxAge, config.LogMaxSize, config.LogMaxBackups, int8(config.LogLevel))
	//初始化token
	token.Init(config.EtcdAddress)
}

func initConfig(){
	myconfig.LoadConfig("./conf/config.conf")
	config.ConsulAddress = myconfig.Config.GetString("consul_address")
	config.EtcdAddress = myconfig.Config.GetString("etcd_address")
	config.LogPath =  myconfig.Config.GetString("log_path")
	config.LogLevel,_ =  myconfig.Config.GetInt("log_level")
	config.LogMaxAge,_ =  myconfig.Config.GetInt("log_max_age")
	config.LogMaxSize,_ =  myconfig.Config.GetInt("log_max_size")
	config.LogMaxBackups,_ =  myconfig.Config.GetInt("log_max_backups")
	config.MongoReplicas = myconfig.Config.GetString("mongo_replicas")
	config.MongoMaxPoolSize,_ = myconfig.Config.GetInt64("mongo_maxPoolSize")
	config.MongoMinPoolSize,_ = myconfig.Config.GetInt64("mongo_minPoolSize")
	config.AppID = myconfig.Config.GetString("appid")
	config.AppSecret = myconfig.Config.GetString("appsecret")
	config.WxDomain = myconfig.Config.GetString("wx_domain")
}

//初始化mysql
func initMySQL() error{
	cons, err := myconfig.Config.GetInt("mysql_cons")
	if err != nil {
		mylog.Error("load mysql connections config error")
		return fmt.Errorf("load mysql connections config error")
	}
	datasrc := myconfig.Config.GetString("mysql_datasrc")
	if !mysql.InitConnectionPool(cons, datasrc) {
		mylog.Error("initialize connection pool fail")
		time.Sleep(time.Millisecond * 100)
		return fmt.Errorf("initialize connection pool fail")
	}
	sqlDebug,_ := myconfig.Config.GetBool("mysql_debug")
	if sqlDebug {
		mysql.OpenDebug()
	} else {
		mysql.CloseDebug()
	}
	mylog.Info("initialize connection pool success")

	return nil
}

func InitMongodb() (err error) {
	//链接单节点mongdb集群
	/*ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	Pool, err = mongo.Connect(ctx,
		options.Client().
			ApplyURI(config.MongoReplicas).
			SetMaxPoolSize(uint64(config.MongoMaxPoolSize)).
			SetMinPoolSize(uint64(config.MongoMinPoolSize)))
	if err != nil {
		return err
	}
	err = Pool.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}*/
	//链接单节点mongdb
	// Set client options
	clientOptions := options.Client().ApplyURI(config.MongoReplicas)

	// Connect to MongoDB
	config.Client, err = mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		mylog.Error(err.Error())
	}

	// Check the connection
	err = config.Client.Ping(context.TODO(), nil)

	if err != nil {
		mylog.Error(err.Error())
	}
	mylog.Info("mongodb init successfully")
	return nil
}