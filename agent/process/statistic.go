package process

import (
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/process"
)

type ProcessInfo struct {
	NativateProcess *process.Process
	Tags            *ProcessTags
	Fields          *ProcessFields
}

type ProcessTags struct {
	ServiceName string `json:"ServiceName"`
	HostName    string `json:"HostName"`
	IP          string `json:"IP"`
	ProcessName string `json:"ProcessName"`
}

type ProcessFields struct {
	CreateTime    int64   `json:"CreateTime"`
	ProcessID     int32   `json:"ProcessID"`
	Status        string  `json:"Status"`
	CPUPercent    float64 `json:"CPUPercent"`
	MemoryPercent float32 `json:"MemoryPercent"`
	MemoryRSS     uint64  `json:"MemoryRSS"`
	MemoryVMS     uint64  `json:"MemoryVMS"`
	DiskReadRate  float64 `json:"DiskReadRate"`
	DiskWriteRate float64 `json:"DiskWriteRate"`
	NetReadRate   float64 `json:"NetReadRate"`
	NetWriteRate  float64 `json:"NetWriteRate"`
}

type ProcessInfoUpdate interface {
	Update()
}

//Update 更新进程所有信息
func (processInfo *ProcessInfo) Update() {
	proce := processInfo.NativateProcess

	isRunning, _ := proce.IsRunning()
	//is not running, return
	if !isRunning {
		return
	}

	processInfo.Fields.CPUPercent, _ = proce.CPUPercent()
	processInfo.Fields.MemoryPercent, _ = proce.MemoryPercent()
	memoryInfo, _ := proce.MemoryInfo()
	processInfo.Fields.MemoryRSS = memoryInfo.RSS
	processInfo.Fields.MemoryVMS = memoryInfo.VMS

	Netstat("tcp", processInfo)
}

//NewProcessInfo 构造新的ProcessInfo
func NewProcessInfo(pid int32, monitorService string, monitorIP string) (*ProcessInfo, error) {
	var processInfo ProcessInfo
	hostInfo, _ := host.Info()
	hostName := hostInfo.Hostname
	proce, err := process.NewProcess(pid)
	if err != nil {
		return nil, err
	}
	processInfo.NativateProcess = proce
	processName, _ := proce.Name()
	processInfo.Tags = &ProcessTags{
		monitorService,
		hostName,
		monitorIP,
		processName,
	}

	createTime, _ := proce.CreateTime()
	CPUPercent, _ := proce.CPUPercent()
	MemoryPercent, _ := proce.MemoryPercent()
	memoryInfo, _ := proce.MemoryInfo()
	MemoryRSS := memoryInfo.RSS
	MemoryVMS := memoryInfo.VMS
	status, _ := proce.Status()
	processInfo.Fields = &ProcessFields{
		createTime,
		proce.Pid,
		status,
		CPUPercent,
		MemoryPercent,
		MemoryRSS,
		MemoryVMS,
		0, 0, 0, 0,
	}
	return &processInfo, nil
}
