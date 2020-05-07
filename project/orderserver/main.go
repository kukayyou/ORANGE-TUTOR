package main

import (
	"github.com/kukayyou/commonlib/myconfig"
	"github.com/kukayyou/commonlib/myhttp"
	"github.com/kukayyou/commonlib/mylog"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/etcd"
	"github.com/micro/go-micro/web"
	//"github.com/micro/go-plugins/registry/consul"
	"go.uber.org/zap"
	"orderserver/config"
	"orderserver/routers"
)

var sugarLogger *zap.SugaredLogger

func main() {
	//初始化配置
	//initConfig()
	//初始化日志
	//initLog()
	//初始化路由
	ginRouter := routers.InitRouters()
	//新建一个consul注册的地址，也就是我们consul服务启动的机器ip+端口
	/*consulReg := consul.NewRegistry(
		registry.Addrs(config.ConsulAddress),
	)*/
	etcdReg := etcd.NewRegistry(
		registry.Addrs("192.168.109.131:12379"))
	myhttp.EtcdAddr = "192.168.109.131:12379"
	//注册服务
	microService:= web.NewService(
		web.Name("api.tutor.com.orderserver"),
		//web.RegisterTTL(time.Second*30),//设置注册服务的过期时间
		//web.RegisterInterval(time.Second*20),//设置间隔多久再次注册服务
		web.Address(":18002"),
		web.Handler(ginRouter),
		web.Registry(etcdReg),
		)

	microService.Run()
}

func initConfig() {
	myconfig.LoadConfig("./conf/config.conf")
	config.ConsulAddress = myconfig.Config.GetString("consul_address")
	config.LogPath =  myconfig.Config.GetString("log_path")
	config.LogLevel =  int8(myconfig.Config.GetInt64("log_level"))
	config.LogMaxAge =  int(myconfig.Config.GetInt64("log_max_age"))
	config.LogMaxSize =  int(myconfig.Config.GetInt64("log_max_size"))
	config.LogMaxBackups =  int(myconfig.Config.GetInt64("log_max_backups"))
}

func init() {
	initConfig()
	mylog.InitLog(config.LogPath,"orderserver", config.LogMaxAge, config.LogMaxSize, config.LogMaxBackups, config.LogLevel)
}
