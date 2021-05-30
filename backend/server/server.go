package server

import "context"

type MonitorServer interface {
	Start(ctx context.Context)
}
