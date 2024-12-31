package systemHelper

type SystemData struct {
	CPU    CpuData    `json:"cpu"`
	Memory MemoryData `json:"memory"`
	Disk   []DiskData `json:"disk"`
}

type CpuData struct {
	CoreNum  int     `json:"core_num"`
	TotalPer float64 `json:"total_per"` // 总体使用率
}

type MemoryData struct {
	Total       float64 `json:"total"`
	Available   float64 `json:"available"`
	Used        float64 `json:"used"`
	UsedPercent float64 `json:"used_percent"`
}

type DiskData struct {
	Device     string    `json:"device"`      // 分区
	MountPoint string    `json:"mount_point"` // 挂载点
	FsType     string    `json:"fs_type"`     // 文件系统类型
	Usage      DiskUsage `json:"usage"`
	Opts       []string  `json:"opts"`
}

type DiskUsage struct {
	Path              string  `json:"path"`
	FsType            string  `json:"fs_type"`
	Total             float64 `json:"total"`
	Free              float64 `json:"free"`
	Used              float64 `json:"used"`
	UsedPercent       float64 `json:"usedPercent"`
	InodesTotal       uint64  `json:"inodesTotal"`
	InodesUsed        uint64  `json:"inodesUsed"`
	InodesFree        uint64  `json:"inodesFree"`
	InodesUsedPercent float64 `json:"inodesUsedPercent"`
}

type ProgressData struct {
	PID           int32   `json:"pid"`
	Name          string  `json:"name"`
	CpuPercent    float64 `json:"cpu_percent"`     // cpu占用
	MemoryUsageMB float64 `json:"memory_usage_mb"` // 内存使用MB（常驻集）
	MemoryPer     float64 `json:"memory_per"`      // 内存占用
}
