package systemHelper

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
	"sort"
	"strings"
	"time"
)

func GetSystemInfo(getLogicalDisk bool) (systemData SystemData, err error) {

	systemData = SystemData{
		CPU:    CpuData{},
		Memory: MemoryData{},
		Disk:   []DiskData{},
	}

	if systemData.CPU, err = GetCpuData(); err != nil {
		return
	}
	if systemData.Memory, err = GetMemoryData(); err != nil {
		return
	}
	if systemData.Disk, err = GetDiskData(getLogicalDisk); err != nil {
		return
	}

	return
}

// GetDiskData 获取硬盘信息
func GetDiskData(getLogicalDisk bool) (diskData []DiskData, err error) {

	diskData = []DiskData{}

	partitions, err := disk.Partitions(true)
	if err != nil {
		return
	}

	for _, partition := range partitions {

		var temp = DiskData{
			Device:     partition.Device,
			MountPoint: partition.Mountpoint,
			FsType:     partition.Fstype,
			Usage:      DiskUsage{},
			Opts:       nil,
		}

		// 获取每个分区的使用情况
		var usage *disk.UsageStat
		usage, err = disk.Usage(partition.Mountpoint)
		if err != nil {
			return
		}

		// getLogicalDisk 如果不获取逻辑盘，过滤一下
		if !getLogicalDisk {
			if usage.Total == 0 || usage.Used == 0 {
				continue
			}
			if usage.Fstype == "/dev/shm" || usage.Fstype == "/run" {
				// /dev/shm: 这是一个常见的 tmpfs 挂载点，用于共享内存。它通常用于进程间通信或需要快速读写的小文件的场景。
				// /run: 这是另一个 tmpfs 挂载点，通常用于存储系统运行时数据，比如进程 ID 文件、套接字文件等。这些数据在系统重启后不需要保留，因此适合使用 tmpfs。
				// 这两个挂载点的 tmpfs 文件系统都利用了系统的内存来存储数据，因此它们的总容量和使用情况会随着系统内存的变化而变化。
				continue
			}
		}

		temp.Usage.Total = float64(usage.Total) / 1024 / 1024 / 1024
		temp.Usage.Used = float64(usage.Used) / 1024 / 1024 / 1024
		temp.Usage.Free = float64(usage.Free) / 1024 / 1024 / 1024
		temp.Usage.UsedPercent = usage.UsedPercent

		diskData = append(diskData, temp)
	}

	return
}

// GetCpuData 获取cpu信息
func GetCpuData() (cpuData CpuData, err error) {
	cpuData = CpuData{}

	cpuData.CoreNum, err = cpu.Counts(true) // true表示逻辑核心数，false表示物理核心数
	if err != nil {
		return
	}

	// 获取整体 CPU 的使用率
	var percentages []float64
	// 第一个参数是一个 time.Duration 类型，用于指定测量间隔。在这个例子中，我们使用 1*time.Second，表示测量间隔为 1 秒。
	// 第二个参数是一个布尔值，表示是否按每个 CPU 核心返回使用率。如果为 true，则返回每个核心的使用率；如果为 false，则返回整体 CPU 的使用率。
	percentages, err = cpu.Percent(1*time.Second, false)
	if err != nil || len(percentages) == 0 {
		return
	}
	cpuData.TotalPer = percentages[0]

	return
}

// GetMemoryData 获取内存信息
func GetMemoryData() (memoryData MemoryData, err error) {
	memoryData = MemoryData{}
	//【2】获取内存信息
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return
	}

	memoryData.Total = float64(vmStat.Total) / 1024 / 1024
	memoryData.Available = float64(vmStat.Available) / 1024 / 1024
	memoryData.Used = float64(vmStat.Used) / 1024 / 1024
	memoryData.UsedPercent = vmStat.UsedPercent

	return
}

type SequenceFlag int

const CpuUsage SequenceFlag = 1
const MemoryUsage SequenceFlag = 2

// GetProgressRank 获取cpu占用进程排行榜
func GetProgressRank(topNum int, sequenceFlag SequenceFlag, nameFilter string) (processInfos []ProgressData, err error) {

	// 获取所有进程
	processes, err := process.Processes()
	if err != nil {
		return
	}

	processInfos = []ProgressData{}

	// 遍历每个进程，获取其CPU和内存使用率
	for _, progress := range processes {

		var progressData = ProgressData{
			PID: progress.Pid,
		}
		progressData.Name, _ = progress.Name()

		// 如果名字筛选不为空，筛选一下检查进程名称是否包含指定关键字
		if nameFilter != "" && !strings.Contains(progressData.Name, nameFilter) {
			continue
		}

		// cpu占用
		progressData.CpuPercent, _ = progress.CPUPercent()

		// 内存占用
		memInfo, _ := progress.MemoryInfo()
		progressData.MemoryUsage = memInfo.RSS

		processInfos = append(processInfos, progressData)
	}

	// 确定排序规则
	if sequenceFlag == CpuUsage {
		// 按 CPU 使用率排序
		sort.Slice(processInfos, func(i, j int) bool {
			return processInfos[i].CpuPercent > processInfos[j].CpuPercent
		})
	} else if sequenceFlag == MemoryUsage {
		// 按内存使用量排序
		sort.Slice(processInfos, func(i, j int) bool {
			return processInfos[i].MemoryUsage > processInfos[j].MemoryUsage
		})
	}

	if len(processInfos) > topNum {
		return processInfos[0:topNum], nil
	}
	return processInfos, nil
}
