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
	"time"

	"userserver/config"
	//"github.com/micro/go-plugins/registry/consul"
	"userserver/routers"
)

func main() {

	defer mylog.SugarLogger.Sync()
	if err:= initMySQL();err != nil{
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
	myconfig.LoadConfig("./conf/config.conf")
	config.ConsulAddress = myconfig.Config.GetString("consul_address")
	config.EtcdAddress = myconfig.Config.GetString("etcd_address")
	config.LogPath =  myconfig.Config.GetString("log_path")
	config.LogLevel,_ =  myconfig.Config.GetInt("log_level")
	config.LogMaxAge,_ =  myconfig.Config.GetInt("log_max_age")
	config.LogMaxSize,_ =  myconfig.Config.GetInt("log_max_size")
	config.LogMaxBackups,_ =  myconfig.Config.GetInt("log_max_backups")
	mylog.InitLog(config.LogPath,"userserver", config.LogMaxAge, config.LogMaxSize, config.LogMaxBackups, int8(config.LogLevel))
	token.Init(config.EtcdAddress)
}

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
