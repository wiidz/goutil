package systemHelper

import (
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
	"os"
	"sort"
	"strings"
	"syscall"
)

type SequenceFlag int

const CpuUsage SequenceFlag = 1
const MemoryUsage SequenceFlag = 2

// GetProgressRank 获取cpu占用进程排行榜
func GetProgressRank(topNum int, sequenceFlag SequenceFlag, nameFilter string) (processInfos []ProgressData, err error) {

	//【1】获取系统总内存信息，用于计算内存占用
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return
	}

	//【2】获取所有进程
	processes, err := process.Processes()
	if err != nil {
		return
	}

	processInfos = []ProgressData{}

	//【3】遍历每个进程，获取其CPU和内存使用率
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
		progressData.MemoryUsageMB = float64(memInfo.RSS) / 1024 / 1024
		progressData.MemoryPer = (float64(memInfo.RSS) / float64(vmStat.Total)) * 100

		// 获取命令行参数
		progressData.CmdLines, _ = progress.CmdlineSlice()

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
			return processInfos[i].MemoryUsageMB > processInfos[j].MemoryUsageMB
		})
	}

	if len(processInfos) > topNum {
		return processInfos[0:topNum], nil
	}
	return processInfos, nil
}

// KillProgress 杀死进程
func KillProgress(pid int, force bool) (err error) {

	// 查找进程并发送 SIGTERM 信号
	process, err := os.FindProcess(pid)
	if err != nil {
		return
	}

	// 发送 SIGTERM 信号
	signal := syscall.SIGTERM
	if force {
		signal = syscall.SIGKILL
	}
	err = process.Signal(signal)
	return
}
