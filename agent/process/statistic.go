package process

type ProcessInfo struct {
	ServiceName   string  `json:"ServiceName"`
	HostName      string  `json:"HostName"`
	Status        string  `json:"Status"`
	IP            string  `json:"IP"`
	ProcessName   string  `json:"ProcessName"`
	ProcessID     int32   `json:"ProcessID"`
	CreateTime    int64   `json:"CreateTime"`
	CPUPercent    float64 `json:"CPUPercent"`
	MemoryPercent float32 `json:"MemoryPercent"`
	MemoryRSS     uint64  `json:"MemoryRSS"`
	MemoryVMS     uint64  `json:"MemoryVMS"`
	DiskReadRate  float64 `json:"DiskReadRate"`
	DiskWriteRate float64 `json:"DiskWriteRate"`
	NetReadRate   float64 `json:"NetReadRate"`
	NetWriteRate  float64 `json:"NetWriteRate"`
}
