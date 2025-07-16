// Package api 提供API服务相关功能
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sign_agent/task"
)

// handleExecuteTask 处理任务执行请求
func (s *Server) handleExecuteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "仅支持POST请求",
		})
		return
	}

	var taskReq TaskRequest
	err := json.NewDecoder(r.Body).Decode(&taskReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: fmt.Sprintf("无法解析请求体: %v", err),
		})
		return
	}

	switch taskReq.Type {
	case "1": // curl命令执行
		result, err := task.ExecuteCurlCommand(taskReq.Command)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{
				Success: false,
				Message: fmt.Sprintf("执行curl命令失败: %v", err),
			})
			return
		}

		json.NewEncoder(w).Encode(Response{
			Success: true,
			Data:    result,
		})

	case "2": // 预留给Node.js执行
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Node.js命令执行功能尚未实现",
		})

	case "3": // 预留给Python执行
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Python命令执行功能尚未实现",
		})

	default:
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: fmt.Sprintf("不支持的任务类型: %s", taskReq.Type),
		})
	}
}
