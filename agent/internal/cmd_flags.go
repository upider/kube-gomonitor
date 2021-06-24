package internal

type CmdFlags struct {
	NacosIPs        []string
	NacosPort       uint64
	Help            bool
	MonitorIP       string
	MonitorPid      int32
	MonitorService  string
	MonitorInterval int64
	DBUrl           string
	Organization    string
	Bucket          string
	Token           string
}
