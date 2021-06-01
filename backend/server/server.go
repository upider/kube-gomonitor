package server

import "context"

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
}

type MonitorServer interface {
	Start(ctx context.Context)
}
