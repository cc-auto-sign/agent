// Package api 提供API服务相关功能
package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sign_agent/config"
)

// Server 结构体是从原始server.go移动过来的

// Server API服务器结构体
type Server struct {
	config *config.Config
	server *http.Server
}

// NewServer 创建一个新的API服务器
func NewServer(cfg *config.Config) *Server {
	return &Server{
		config: cfg,
	}
}

// Start 启动API服务
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// 注册API路由
	mux.HandleFunc("/api/system/info", s.handleAuthMiddleware(s.handleSystemInfo))
	mux.HandleFunc("/api/task/execute", s.handleAuthMiddleware(s.handleExecuteTask))
	mux.HandleFunc("/api/health", s.handleHealth)

	addr := fmt.Sprintf(":%d", s.config.GetPort())
	s.server = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Printf("API服务启动，监听地址: %s\n", addr)
	return s.server.ListenAndServe()
}

// Stop 停止API服务
func (s *Server) Stop() error {
	if s.server != nil {
		return s.server.Close()
	}
	return nil
}

// 健康检查端点
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "done",
	})
}
