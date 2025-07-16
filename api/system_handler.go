// Package api 提供API服务相关功能
package api

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"sign_agent/system"
)

// handleSystemInfo 处理系统信息请求
func (s *Server) handleSystemInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 获取系统信息
	totalMem, usedMem, err := system.GetMemoryInfo()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: fmt.Sprintf("获取内存信息失败: %v", err),
		})
		return
	}

	cpuUsage, err := system.GetCPUUsage()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: fmt.Sprintf("获取CPU信息失败: %v", err),
		})
		return
	}

	// 计算内存使用率百分比
	memoryUsagePerc := 0.0
	if totalMem > 0 {
		memoryUsagePerc = float64(usedMem) / float64(totalMem) * 100.0
		// 保留两位小数
		memoryUsagePerc = math.Round(memoryUsagePerc*100) / 100
	}

	// 转换为MB，并精确到整数
	totalMemoryMB := math.Round(float64(totalMem) / (1024 * 1024))
	usedMemoryMB := math.Round(float64(usedMem) / (1024 * 1024))

	// CPU使用率保留两位小数
	cpuUsage = math.Round(cpuUsage*100) / 100

	sysInfo := SystemInfo{
		TotalMemoryMB:   int(totalMemoryMB),
		UsedMemoryMB:    int(usedMemoryMB),
		MemoryUsagePerc: memoryUsagePerc,
		CPUUsagePerc:    cpuUsage,
	}

	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data:    sysInfo,
	})
}
