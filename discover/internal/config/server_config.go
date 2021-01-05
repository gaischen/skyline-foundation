package config

type ServerConfig struct {
	Addr       string
	Port       int
	ServerType int //1 leader 0 slave

	logPath      string //log 输出地址
	logFormatter string

	serverList string //获取server list的地址，需要轮询去获取
}

func NewDefaultConfig() *ServerConfig {
	return &ServerConfig{}
}
