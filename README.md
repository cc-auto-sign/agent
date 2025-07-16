# 自动签到系统 Agent

这是一个服务器节点代理程序，可以接收指令并执行任务，同时提供系统资源信息。

## 功能特性

- 通过HTTP API对外提供服务器系统信息（内存、CPU使用率）
- 接收并执行客户端下发的任务
- 安全认证机制，使用安全密钥保护API访问
- 支持作为系统服务运行
- 提供命令行工具管理服务

## 安装与使用

### 编译

```bash
go build -o checkin-agent
```

### 初次运行

```bash
./checkin-agent
```

初次运行会自动创建配置文件和生成安全密钥。

### 命令行选项

- `--config`：指定配置文件路径（默认：./agent_config.json）
- `--version`：显示版本信息
- `--regenerate-key`：重新生成安全密钥
- `--install-service`：安装为系统服务（仅Linux）
- `--uninstall-service`：卸载系统服务（仅Linux）
- `--service-status`：查看服务状态（仅Linux）

### 作为服务运行

在Linux系统上安装为服务（需要root权限）：

```bash
sudo ./checkin-agent --install-service
```

然后可以使用systemctl管理服务：

```bash
# 启动服务
systemctl start checkin-agent

# 设置开机自启
systemctl enable checkin-agent

# 查看服务状态
systemctl status checkin-agent

# 停止服务
systemctl stop checkin-agent
```

## API接口

### 系统信息

```
GET /api/system/info
```

需要在请求中提供安全密钥（HTTP头`X-Secure-Key`或JSON/表单参数`secure_key`）。

返回系统信息，包括：
- 总内存（MB，整数）
- 已使用内存（MB，整数）
- 内存使用率（百分比，保留两位小数）
- CPU使用率（百分比，保留两位小数）

### 执行任务

```
POST /api/task/execute
```

请求体：

```json
{
  "type": "任务类型",
  "command": "任务命令",
  "secure_key": "安全密钥"
}
```

支持的任务类型：
- `1`: 执行curl命令，安全解析并执行HTTP请求（支持忽略SSL验证）
- `2`: Node.js命令执行（尚未实现）
- `3`: Python命令执行（尚未实现）
**这里我的想法是用类似dify的docker沙盒去执行代码，防止有问题的代码**

注意：
- curl命令会被智能解析而不是直接执行，可以安全地处理URL中的特殊字符(&、|、$等)、JSON数据等
- 不允许使用分号(;)来链接多个命令
- 支持标准curl选项如`-H`(设置头信息)、`-d`(发送数据)、`-X`(设置请求方法)、`-k`(忽略SSL验证)等
- 可以传递复杂的JSON或包含特殊字符的参数，因为系统会正确解析引号内的内容

### 健康检查

```
GET /api/health
```

检查服务是否正常运行。

## 配置文件

配置文件为JSON格式，默认路径为`./agent_config.json`：

```json
{
  "secure_key": "生成的安全密钥",
  "port": 8080
}
```

## 安全性

- 所有API请求都需要提供有效的安全密钥
- 安全密钥在初次运行时自动生成，也可以使用命令重新生成
- 建议将配置文件权限设置为仅管理员可读
