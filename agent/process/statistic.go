package process

import (
	"gomonitor/utils"
	"os"

	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/process"
	"k8s.io/client-go/rest"
)

var (
	NODE_NAME = "NODE_NAME"
	NODE_IP   = "NODE_IP"
	POD_NAME  = "POD_NAME"
	POD_IP    = "POD_IP"
)

type ProcessInfo struct {
	NativateProcess *process.Process
	Tags            *ProcessTags
	Fields          *ProcessFields
}

type ProcessTags struct {
	ServiceName string `json:"serviceName"`
	ProcessName string `json:"processName"`
	ProcessIP   string `json:"processIP"`
	NodeName    string `json:"nodeName"`
	NodeIP      string `json:"nodeIP"`
	PodName     string `json:"podName"`
	PodIP       string `json:"podIP"`
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
	_, err := utils.CheckInK8s()
	var hostName string
	var hostIP string
	if err == rest.ErrNotInCluster {
		hostInfo, _ := host.Info()
		hostName = hostInfo.Hostname
		hostIP = monitorIP
	} else {
		hostName = os.Getenv(NODE_NAME)
		hostIP = os.Getenv(NODE_IP)
	}

	proce, err := process.NewProcess(pid)
	if err != nil {
		return nil, err
	}
	processInfo.NativateProcess = proce
	processName, _ := proce.Name()
	processInfo.Tags = &ProcessTags{
		ServiceName: monitorService,
		ProcessIP:   monitorIP,
		ProcessName: processName,
		NodeName:    hostName,
		NodeIP:      hostIP,
		PodName:     os.Getenv(POD_NAME),
		PodIP:       os.Getenv(POD_IP),
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
