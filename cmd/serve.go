package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sign_agent/api"
	"sign_agent/config"
	"syscall"
)

// StartServer 启动服务器
func StartServer(configPath string) error {
	// 加载配置
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}

	// 创建API服务器
	server := api.NewServer(cfg)

	// 设置信号处理，优雅退出
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// 在后台启动服务器
	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 等待信号
	sig := <-sigCh
	log.Printf("接收到信号 %v, 正在优雅退出...", sig)

	// 停止服务器
	if err := server.Stop(); err != nil {
		return fmt.Errorf("停止服务器失败: %v", err)
	}

	return nil
}

// RegenerateKey 重新生成安全密钥
func RegenerateKey(configPath string) error {
	// 加载配置
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}

	// 重新生成密钥
	newKey, err := cfg.RegenerateSecureKey()
	if err != nil {
		return fmt.Errorf("重新生成安全密钥失败: %v", err)
	}

	fmt.Printf("安全密钥已重新生成: %s\n", newKey)
	return nil
}
