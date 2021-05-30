package report

import (
	"context"
)

type Reporter interface {
	Start(ctx context.Context)
	Close()
	Report()
}
