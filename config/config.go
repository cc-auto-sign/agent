package config

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
)

// Config 配置结构
type Config struct {
	SecureKey string `json:"secure_key"`
	Port      int    `json:"port"`
	filePath  string // 配置文件路径
}

// LoadConfig 加载配置文件
func LoadConfig(configPath string) (*Config, error) {
	// 如果未指定配置路径，使用默认路径
	if configPath == "" {
		configPath = "./agent_config.json"
	}

	// 确保路径是绝对路径
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, fmt.Errorf("获取配置文件绝对路径失败: %v", err)
	}

	// 检查配置文件是否存在
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		// 配置文件不存在，创建默认配置
		return createDefaultConfig(absPath)
	} else if err != nil {
		return nil, fmt.Errorf("检查配置文件状态失败: %v", err)
	}

	// 读取配置文件
	file, err := os.Open(absPath)
	if err != nil {
		return nil, fmt.Errorf("打开配置文件失败: %v", err)
	}
	defer file.Close()

	// 解析配置
	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 设置文件路径
	config.filePath = absPath

	// 验证配置
	if err := config.validate(); err != nil {
		return nil, err
	}

	return &config, nil
}

// 创建默认配置
func createDefaultConfig(path string) (*Config, error) {
	// 生成安全密钥
	secureKey, err := generateSecureKey()
	if err != nil {
		return nil, fmt.Errorf("生成安全密钥失败: %v", err)
	}

	// 创建默认配置
	config := &Config{
		SecureKey: secureKey,
		Port:      8080,
		filePath:  path,
	}

	// 保存配置
	if err := config.Save(); err != nil {
		return nil, err
	}

	fmt.Printf("已创建默认配置文件: %s\n", path)
	fmt.Printf("安全密钥: %s\n", secureKey)

	return config, nil
}

// 生成安全密钥
func generateSecureKey() (string, error) {
	// 生成32字节随机数据
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	// 转换为十六进制字符串
	return hex.EncodeToString(bytes), nil
}

// 验证配置
func (c *Config) validate() error {
	if c.SecureKey == "" {
		return fmt.Errorf("配置缺少安全密钥")
	}
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("端口号无效: %d", c.Port)
	}
	return nil
}

// Save 保存配置到文件
func (c *Config) Save() error {
	// 序列化配置
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}

	// 确保目录存在
	dir := filepath.Dir(c.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	// 写入文件
	if err := ioutil.WriteFile(c.filePath, data, 0600); err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	return nil
}

// GetSecureKey 获取安全密钥
func (c *Config) GetSecureKey() string {
	return c.SecureKey
}

// GetPort 获取端口号
func (c *Config) GetPort() int {
	return c.Port
}

// RegenerateSecureKey 重新生成安全密钥
func (c *Config) RegenerateSecureKey() (string, error) {
	// 生成新密钥
	newKey, err := generateSecureKey()
	if err != nil {
		return "", err
	}

	// 更新配置
	c.SecureKey = newKey

	// 保存配置
	if err := c.Save(); err != nil {
		return "", err
	}

	return newKey, nil
}
