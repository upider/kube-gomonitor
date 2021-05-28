package report

import (
	"fmt"
	"gomonitor/agent/process"
)

//Report send process info to db
func Report(processInfo *process.ProcessInfo) {
	fmt.Println(processInfo)
}
