// Package api 提供API服务相关功能
package api

// 类型定义部分，这些类型是从原始server.go文件移动过来的

// SystemInfo 系统信息结构体
type SystemInfo struct {
	TotalMemoryMB   int     `json:"total_memory_mb"`
	UsedMemoryMB    int     `json:"used_memory_mb"`
	MemoryUsagePerc float64 `json:"memory_usage_perc"`
	CPUUsagePerc    float64 `json:"cpu_usage_perc"`
}

// TaskRequest 任务执行请求结构体
type TaskRequest struct {
	Type      string `json:"type"`
	Command   string `json:"command"`
	SecureKey string `json:"secure_key"`
}

// Response API响应结构体
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
