package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sign_agent/cmd"
	"sign_agent/service"
)

const version = "1.0.0"

func main() {
	// 设置命令行参数
	var (
		configPath       string
		showVersion      bool
		regenerateKey    bool
		installService   bool
		uninstallService bool
		serviceStatus    bool
	)

	// 创建主命令
	mainCmd := flag.NewFlagSet("checkin-agent", flag.ExitOnError)
	mainCmd.StringVar(&configPath, "config", "", "配置文件路径 (默认: ./agent_config.json)")
	mainCmd.BoolVar(&showVersion, "version", false, "显示版本信息")
	mainCmd.BoolVar(&regenerateKey, "regenerate-key", false, "重新生成安全密钥")
	mainCmd.BoolVar(&installService, "install-service", false, "安装系统服务 (仅Linux)")
	mainCmd.BoolVar(&uninstallService, "uninstall-service", false, "卸载系统服务 (仅Linux)")
	mainCmd.BoolVar(&serviceStatus, "service-status", false, "检查服务状态 (仅Linux)")

	// 解析命令行参数
	if len(os.Args) > 1 && os.Args[1] == "serve" {
		// 处理serve子命令
		serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)
		serveCmd.StringVar(&configPath, "config", "", "配置文件路径 (默认: ./agent_config.json)")

		if err := serveCmd.Parse(os.Args[2:]); err != nil {
			log.Fatalf("解析参数失败: %v", err)
		}

		// 启动服务器
		if err := cmd.StartServer(configPath); err != nil {
			log.Fatalf("启动服务器失败: %v", err)
		}
		return
	}

	if err := mainCmd.Parse(os.Args[1:]); err != nil {
		log.Fatalf("解析参数失败: %v", err)
	}

	// 处理主命令选项
	if showVersion {
		fmt.Printf("自动签到系统 Agent 版本 %s\n", version)
		return
	}

	if regenerateKey {
		if err := cmd.RegenerateKey(configPath); err != nil {
			log.Fatalf("重新生成密钥失败: %v", err)
		}
		return
	}

	if installService {
		if err := service.InstallService(); err != nil {
			log.Fatalf("安装服务失败: %v", err)
		}
		return
	}

	if uninstallService {
		if err := service.UninstallService(); err != nil {
			log.Fatalf("卸载服务失败: %v", err)
		}
		return
	}

	if serviceStatus {
		status, err := service.ServiceStatus()
		if err != nil {
			log.Fatalf("获取服务状态失败: %v", err)
		}
		fmt.Println(status)
		return
	}

	// 如果没有指定特定选项，启动服务器
	if err := cmd.StartServer(configPath); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
