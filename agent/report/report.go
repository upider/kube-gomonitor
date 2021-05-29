package report

import "gomonitor/agent/process"

type Reporter interface {
	Close()
	Report(processInfo *process.ProcessInfo)
}
