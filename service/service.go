package service

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
)

// 安装系统服务模板
const systemdServiceTemplate = `[Unit]
Description=Auto Checkin Agent Service
After=network.target

[Service]
Type=simple
User=root
ExecStart={{.ExecPath}} serve
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
`

// InstallService 安装系统服务
func InstallService() error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("仅支持在Linux系统上安装服务")
	}

	// 检查是否有root权限
	if os.Geteuid() != 0 {
		return fmt.Errorf("安装系统服务需要root权限，请使用sudo运行")
	}

	// 获取当前可执行文件路径
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("获取可执行文件路径失败: %v", err)
	}

	absExecPath, err := filepath.Abs(execPath)
	if err != nil {
		return fmt.Errorf("获取绝对路径失败: %v", err)
	}

	// 准备服务文件数据
	data := struct {
		ExecPath string
	}{
		ExecPath: absExecPath,
	}

	// 创建服务文件
	serviceName := "checkin-agent.service"
	servicePath := "/etc/systemd/system/" + serviceName

	// 检查服务文件是否已存在
	if _, err := os.Stat(servicePath); err == nil {
		log.Printf("服务文件已存在: %s，将覆盖...", servicePath)
	}

	// 生成服务文件内容
	tmpl, err := template.New("service").Parse(systemdServiceTemplate)
	if err != nil {
		return fmt.Errorf("解析服务模板失败: %v", err)
	}

	file, err := os.Create(servicePath)
	if err != nil {
		return fmt.Errorf("创建服务文件失败: %v", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("写入服务文件失败: %v", err)
	}

	// 重新加载systemd配置
	cmd := exec.Command("systemctl", "daemon-reload")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("重新加载systemd配置失败: %v", err)
	}

	log.Printf("服务文件已安装: %s", servicePath)
	log.Printf("使用以下命令启动服务: systemctl start %s", serviceName)
	log.Printf("使用以下命令设置开机自启: systemctl enable %s", serviceName)

	return nil
}

// UninstallService 卸载系统服务
func UninstallService() error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("仅支持在Linux系统上卸载服务")
	}

	// 检查是否有root权限
	if os.Geteuid() != 0 {
		return fmt.Errorf("卸载系统服务需要root权限，请使用sudo运行")
	}

	serviceName := "checkin-agent.service"

	// 停止服务
	cmd := exec.Command("systemctl", "stop", serviceName)
	_ = cmd.Run() // 忽略错误，可能服务未运行

	// 禁用服务
	cmd = exec.Command("systemctl", "disable", serviceName)
	_ = cmd.Run() // 忽略错误，可能服务未启用

	// 删除服务文件
	servicePath := "/etc/systemd/system/" + serviceName
	if err := os.Remove(servicePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("删除服务文件失败: %v", err)
	}

	// 重新加载systemd配置
	cmd = exec.Command("systemctl", "daemon-reload")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("重新加载systemd配置失败: %v", err)
	}

	log.Printf("服务已成功卸载: %s", serviceName)
	return nil
}

// IsServiceActive 检查服务是否正在运行
func IsServiceActive() (bool, error) {
	if runtime.GOOS != "linux" {
		return false, fmt.Errorf("仅支持在Linux系统上检查服务状态")
	}

	cmd := exec.Command("systemctl", "is-active", "checkin-agent.service")
	output, err := cmd.Output()
	if err != nil {
		// 如果命令执行失败，服务可能未安装或未运行
		return false, nil
	}

	return strings.TrimSpace(string(output)) == "active", nil
}

// ServiceStatus 获取服务状态信息
func ServiceStatus() (string, error) {
	if runtime.GOOS != "linux" {
		return "", fmt.Errorf("仅支持在Linux系统上获取服务状态")
	}

	cmd := exec.Command("systemctl", "status", "checkin-agent.service")
	output, err := cmd.CombinedOutput()
	if err != nil {
		// 命令可能因服务未运行而返回非零状态码，但我们仍然想返回输出
		return string(output), nil
	}

	return string(output), nil
}
