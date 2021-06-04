package server

type ServerFlags struct {
	NacosIP             string
	NacosPort           uint64
	NamespaceId         string
	ServiceName         string
	ServerIP            string
	MonitorServices     []string
	MonitorServiceGroup string
	DBUrl               string
	Organization        string
	Bucket              string
	Token               string
	Interval            uint64
}

type MonitorServer interface {
	Start()
}
