package config

var (
	ConsulAddress string //consul地址
	LogPath       string //日志地址
	LogLevel      int8   //日志级别
	LogMaxAge     int    //日志保留时长
	LogMaxSize    int    //日志最大size
	LogMaxBackups int    //日志最多留备份数
)
