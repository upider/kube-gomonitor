package server

type ServerFlags struct {
	AgentImage           string
	NacosIPs             []string
	NacosPort            uint64
	Namespaces           []string
	ServiceName          string
	ServerIP             string
	MonitorServices      []string
	MonitorServiceGroups []string
	DBUrl                string
	Organization         string
	Bucket               string
	Token                string
	Interval             uint64
}

type MonitorServer interface {
	Start()
}
