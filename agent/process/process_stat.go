package process

import (
	// "auto-monitor/packet"

	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/process"
)

//ProcessStat 得到进程所有信息
func ProcessStat(pid int32, serviceName string) (*ProcessInfo, error) {
	proce, err := process.NewProcess(pid)
	if err != nil {
		return nil, err
	}

	var processInfo ProcessInfo
	processInfo.ServiceName = serviceName
	processInfo.ProcessID = pid
	processInfo.Status, _ = proce.Status()

	isRunning, _ := proce.IsRunning()
	//is not running, return
	if !isRunning {
		return &processInfo, nil
	}

	hostInfo, _ := host.Info()
	processInfo.HostName = hostInfo.Hostname
	processInfo.ProcessName, _ = proce.Name()
	processInfo.CreateTime, _ = proce.CreateTime()
	processInfo.CPUPercent, _ = proce.CPUPercent()
	processInfo.MemoryPercent, _ = proce.MemoryPercent()
	memoryInfo, _ := proce.MemoryInfo()
	processInfo.MemoryRSS = memoryInfo.RSS
	processInfo.MemoryVMS = memoryInfo.VMS

	Netstat("tcp", &processInfo)

	return &processInfo, nil
}
