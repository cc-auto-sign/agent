// Package task 提供任务执行相关功能
package task

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/mattn/go-shellwords"
)

// ExecuteCurlCommand 执行curl命令，安全地解析和执行curl请求
func ExecuteCurlCommand(cmdStr string) (string, error) {
	// 安全检查：确保命令以curl开头
	cmdStr = strings.TrimSpace(cmdStr)
	if !strings.HasPrefix(cmdStr, "curl") {
		return "", fmt.Errorf("命令必须以curl开头")
	}

	// 检查分号，防止多条命令执行
	if strings.Contains(cmdStr, ";") {
		return "", fmt.Errorf("不允许使用分号执行多条命令")
	}

	// 解析CURL命令并转换为HTTP请求
	return executeHTTPRequest(cmdStr)
}

// executeHTTPRequest 执行HTTP请求，处理复杂的curl命令解析
func executeHTTPRequest(curlCmd string) (string, error) {
	// 初始化HTTP请求参数
	url := ""
	method := "GET"
	headers := make(map[string]string)
	data := ""
	insecure := false

	// 使用更复杂的解析逻辑提取curl参数
	// 正确处理引号和转义
	parts, err := shellwords.Parse(curlCmd)
	if err != nil {
		return "", fmt.Errorf("解析curl命令失败: %v", err)
	}

	if len(parts) < 2 {
		return "", fmt.Errorf("无效的curl命令")
	}

	// 跳过第一个元素(curl命令本身)
	for i := 1; i < len(parts); i++ {
		arg := parts[i]

		// 处理URL (非选项参数)
		if !strings.HasPrefix(arg, "-") && url == "" {
			url = arg
			continue
		}

		// 处理选项
		switch arg {
		case "-X", "--request":
			if i+1 < len(parts) {
				method = parts[i+1]
				i++
			}
		case "-H", "--header":
			if i+1 < len(parts) {
				headerLine := parts[i+1]
				if colonIdx := strings.Index(headerLine, ":"); colonIdx != -1 {
					name := strings.TrimSpace(headerLine[:colonIdx])
					value := strings.TrimSpace(headerLine[colonIdx+1:])
					headers[name] = value
				}
				i++
			}
		case "-d", "--data", "--data-ascii", "--data-binary", "--data-raw":
			if i+1 < len(parts) {
				data = parts[i+1]
				// 如果没有明确指定方法，数据请求默认为POST
				if method == "GET" {
					method = "POST"
				}
				i++
			}
		case "-k", "--insecure":
			insecure = true
		}
	}

	if url == "" {
		return "", fmt.Errorf("未指定URL")
	}

	// 创建HTTP客户端
	client := &http.Client{}
	if insecure {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	// 创建请求
	req, err := http.NewRequest(method, url, strings.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 添加头信息
	for name, value := range headers {
		req.Header.Add(name, value)
	}

	// 如果是POST请求且未指定Content-Type，添加默认值
	if (method == "POST" || method == "PUT") && data != "" && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	// 执行请求
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("执行HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	return string(body), nil
}
