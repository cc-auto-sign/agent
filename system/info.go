package system

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// GetMemoryInfo 获取系统内存信息
func GetMemoryInfo() (uint64, uint64, error) {
	// 尝试使用gopsutil库获取内存信息
	v, err := mem.VirtualMemory()
	if err == nil {
		return v.Total, v.Used, nil
	}

	// 如果gopsutil失败，尝试使用备用方法
	return getFallbackMemoryInfo()
}

// 备用方法：使用命令行工具获取内存信息
func getFallbackMemoryInfo() (uint64, uint64, error) {
	// 对于Windows系统
	if runtime.GOOS == "windows" {
		cmd := exec.Command("wmic", "OS", "get", "TotalVisibleMemorySize,FreePhysicalMemory", "/Value")
		output, err := cmd.Output()
		if err != nil {
			return 0, 0, err
		}

		outputStr := string(output)
		totalStr := strings.TrimSpace(strings.Replace(strings.Split(strings.Split(outputStr, "TotalVisibleMemorySize=")[1], "\n")[0], "\r", "", -1))
		freeStr := strings.TrimSpace(strings.Replace(strings.Split(strings.Split(outputStr, "FreePhysicalMemory=")[1], "\n")[0], "\r", "", -1))

		totalKB, _ := strconv.ParseUint(totalStr, 10, 64)
		freeKB, _ := strconv.ParseUint(freeStr, 10, 64)

		total := totalKB * 1024
		used := (totalKB - freeKB) * 1024

		return total, used, nil
	}

	// 对于Linux系统
	if runtime.GOOS == "linux" {
		cmd := exec.Command("free", "-b")
		output, err := cmd.Output()
		if err != nil {
			return 0, 0, err
		}

		lines := strings.Split(string(output), "\n")
		if len(lines) < 2 {
			return 0, 0, fmt.Errorf("意外的free命令输出格式")
		}

		memLine := strings.Fields(lines[1])
		if len(memLine) < 3 {
			return 0, 0, fmt.Errorf("无法解析内存信息")
		}

		total, _ := strconv.ParseUint(memLine[1], 10, 64)
		used, _ := strconv.ParseUint(memLine[2], 10, 64)

		return total, used, nil
	}

	// 对于macOS系统
	if runtime.GOOS == "darwin" {
		// 获取总内存
		cmd := exec.Command("sysctl", "-n", "hw.memsize")
		output, err := cmd.Output()
		if err != nil {
			return 0, 0, err
		}
		total, _ := strconv.ParseUint(strings.TrimSpace(string(output)), 10, 64)

		// 获取已使用内存
		cmd = exec.Command("vm_stat")
		output, err = cmd.Output()
		if err != nil {
			return 0, 0, err
		}

		lines := strings.Split(string(output), "\n")
		pageSize := uint64(4096) // 默认页大小
		free := uint64(0)

		for _, line := range lines {
			if strings.Contains(line, "Pages free:") {
				parts := strings.Split(line, ":")
				if len(parts) >= 2 {
					val := strings.TrimSpace(strings.Replace(parts[1], ".", "", -1))
					freePages, _ := strconv.ParseUint(val, 10, 64)
					free = freePages * pageSize
				}
				break
			}
		}

		return total, total - free, nil
	}

	return 0, 0, fmt.Errorf("不支持的操作系统")
}

// GetCPUUsage 获取CPU使用率
func GetCPUUsage() (float64, error) {
	// 尝试通过gopsutil获取CPU使用率，采样周期为200ms
	percentages, err := cpu.Percent(200*time.Millisecond, false)
	if err == nil && len(percentages) > 0 {
		return percentages[0], nil
	}

	// 如果gopsutil失败，尝试使用备用方法
	return getFallbackCPUUsage()
}

// 备用方法：使用命令行工具获取CPU使用率
func getFallbackCPUUsage() (float64, error) {
	// 对于Windows系统
	if runtime.GOOS == "windows" {
		cmd := exec.Command("wmic", "cpu", "get", "LoadPercentage")
		output, err := cmd.Output()
		if err != nil {
			return 0, err
		}

		lines := strings.Split(string(output), "\n")
		if len(lines) < 2 {
			return 0, fmt.Errorf("意外的wmic命令输出格式")
		}

		loadStr := strings.TrimSpace(lines[1])
		load, err := strconv.ParseFloat(loadStr, 64)
		if err != nil {
			return 0, err
		}

		return load, nil
	}

	// 对于Linux系统
	if runtime.GOOS == "linux" {
		// 获取两个时间点的CPU统计，计算使用率
		stat1, err := readCPUStat()
		if err != nil {
			return 0, err
		}

		time.Sleep(200 * time.Millisecond)

		stat2, err := readCPUStat()
		if err != nil {
			return 0, err
		}

		idle1 := stat1["idle"]
		idle2 := stat2["idle"]
		total1 := stat1["total"]
		total2 := stat2["total"]

		idleDelta := idle2 - idle1
		totalDelta := total2 - total1

		if totalDelta == 0 {
			return 0, nil
		}

		usage := 100.0 * (1.0 - float64(idleDelta)/float64(totalDelta))
		return usage, nil
	}

	// 对于macOS系统
	if runtime.GOOS == "darwin" {
		cmd := exec.Command("top", "-l", "1", "-n", "0", "-S")
		output, err := cmd.Output()
		if err != nil {
			return 0, err
		}

		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "CPU usage") {
				parts := strings.Split(line, ": ")
				if len(parts) < 2 {
					continue
				}
				usageStr := strings.Split(parts[1], "%")[0]
				usage, err := strconv.ParseFloat(usageStr, 64)
				if err != nil {
					return 0, err
				}
				return usage, nil
			}
		}
	}

	return 0, fmt.Errorf("不支持的操作系统")
}

// 读取Linux /proc/stat 获取CPU统计
func readCPUStat() (map[string]uint64, error) {
	if runtime.GOOS != "linux" {
		return nil, fmt.Errorf("该函数仅支持Linux系统")
	}

	cmd := exec.Command("cat", "/proc/stat")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("无法读取CPU统计")
	}

	cpuLine := strings.Fields(lines[0])
	if len(cpuLine) < 5 || cpuLine[0] != "cpu" {
		return nil, fmt.Errorf("无效的CPU统计格式")
	}

	stat := make(map[string]uint64)
	total := uint64(0)

	// CPU统计格式: cpu user nice system idle iowait irq softirq steal guest guest_nice
	for i := 1; i < len(cpuLine); i++ {
		val, _ := strconv.ParseUint(cpuLine[i], 10, 64)
		total += val
		if i == 4 { // idle位置
			stat["idle"] = val
		}
	}

	stat["total"] = total
	return stat, nil
}

// GetRuntimeInfo 获取Go运行时信息
func GetRuntimeInfo() map[string]interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return map[string]interface{}{
		"num_goroutines":  runtime.NumGoroutine(),
		"allocated_bytes": memStats.Alloc,
		"sys_bytes":       memStats.Sys,
	}
}
