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
	TotalGB           float64 `json:"total_gb"`
	FreeGB            float64 `json:"free_gb"`
	UsedGB            float64 `json:"used_gb"`
	UsedPercent       float64 `json:"used_percent"`
	InodesTotal       uint64  `json:"inodes_total"`
	InodesUsed        uint64  `json:"inodes_used"`
	InodesFree        uint64  `json:"inodes_free"`
	InodesUsedPercent float64 `json:"inodes_used_percent"`
}

type ProgressData struct {
	PID           int32    `json:"pid"`
	Name          string   `json:"name"`
	CpuPercent    float64  `json:"cpu_percent"`     // cpu占用
	MemoryUsageMB float64  `json:"memory_usage_mb"` // 内存使用MB（常驻集）
	MemoryPer     float64  `json:"memory_per"`      // 内存占用
	CmdLines      []string `json:"cmd_lines"`       // 命令参数
}
