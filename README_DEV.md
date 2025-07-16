# 开发者文档

## 项目结构

```
/
├── api/                # API服务相关代码
│   ├── middleware.go   # 中间件
│   ├── server_base.go  # 服务器基础结构
│   ├── system_handler.go # 系统信息处理器
│   ├── task_handler.go # 任务执行处理器
│   └── types.go        # API类型定义
├── cmd/                # 命令处理
│   └── serve.go        # 服务启动逻辑
├── config/             # 配置管理
│   └── config.go       # 配置操作
├── service/            # 系统服务相关
│   └── service.go      # 服务安装与管理
├── system/             # 系统信息相关
│   └── info.go         # 获取系统信息
├── task/               # 任务执行相关
│   ├── curl.go         # curl命令执行
│   └── task.go         # 任务定义
├── main.go             # 主程序
├── go.mod              # Go模块定义
└── README.md           # 使用说明
```

## 代码组织

项目按功能模块拆分为多个包，每个包负责特定的功能：

- **api**: 处理HTTP API相关的请求和响应
- **cmd**: 处理命令行指令
- **config**: 负责配置的加载、保存和验证
- **service**: 管理系统服务（安装、卸载等）
- **system**: 提供系统信息获取功能
- **task**: 处理各类任务的执行

## 开发指南

### 添加新的API端点

1. 在 `api` 包中创建处理函数
2. 在 `server_base.go` 的 `Start()` 方法中注册路由

### 添加新的任务类型

1. 在 `task` 包中创建新的任务执行函数
2. 在 `api/task_handler.go` 中的 `handleExecuteTask` 方法中添加新的任务类型处理

### 修改配置

配置管理在 `config` 包中，修改 `Config` 结构可添加新的配置项。

## 测试

每个包应当有自己的单元测试，确保功能正常工作。

## 安全性考虑

1. 所有API请求都需要验证安全密钥
2. 任务执行需要防止命令注入
3. 配置文件应当设置适当的权限（仅管理员可读）
