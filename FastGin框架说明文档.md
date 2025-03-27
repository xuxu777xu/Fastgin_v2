# FastGin框架说明文档

## 简介

FastGin是一个基于Gin框架的轻量级Web服务框架，专为高性能API开发设计，提供了完整的中间件支持、日志记录、错误处理等功能。框架采用模块化设计，便于开发者快速构建标准化的Web服务。

## 核心特性

### 1. 中间件支持
- **请求追踪**
  - 自动生成请求ID和追踪ID
  - 支持通过请求头自定义ID：`X-Request-ID` 和 `X-Trace-ID`
  - 响应头自动包含这些ID便于追踪

- **日志记录**
  - 自动记录所有请求的详细信息
  - 分级日志（INFO、WARN、ERROR）
  - 错误日志文件独立存储（logs/error.log）
  - 压缩格式的JSON日志，包含完整上下文信息
  - 支持请求体和响应体的记录

- **限流控制**
  - 支持基于令牌桶的限流
  - 可配置每秒请求数和突发流量
  - 默认配置：300次/秒，突发最大500次

- **CORS 支持**
  - 内置跨域请求支持
  - 可配置允许的源、方法和头部

### 2. 错误处理
- 统一的错误响应格式
- 详细的错误日志记录
- 支持错误堆栈追踪
- 自定义错误码和错误信息

### 3. 接口规范
- 统一的请求响应格式
- 支持参数验证
- 支持自定义验证规则

## 项目结构

```
Fastgin_v2/
├── api/                  # API处理层
│   ├── common/           # 公共组件
│   └── test/             # 测试模块 
├── apiServer/            # API服务实现
├── core/                 # 核心功能
├── middleware/           # 中间件组件
│   ├── cors.go           # 跨域处理
│   ├── logger.go         # 日志中间件
│   ├── ratelimit.go      # 限流中间件
│   └── request_id.go     # 请求ID生成
├── pkg/                  # 通用工具包
│   ├── config/           # 配置管理
│   └── logg/             # 日志工具
├── router/               # 路由管理
├── main.go               # 应用入口
└── settings.yaml         # 配置文件
```

## 快速开始

### 1. 安装依赖
```bash
# 初始化项目
go mod init your-project-name

# 安装核心依赖
go get -u github.com/gin-gonic/gin
go get -u github.com/google/uuid
go get -u github.com/go-resty/resty/v2
```

### 2. 创建基础路由
```go
package main

import (
    "your-project-name/router"
)

func main() {
    r := router.InitRouter()
    r.Run(":8080")
}
```

### 3. 添加新的处理器
```go
// handler/your_handler.go
type Handler struct {
    service ServiceInterface
}

func NewHandler() *Handler {
    return &Handler{
        service: NewService(),
    }
}

func (h *Handler) HandleRequest(c *gin.Context) {
    // 处理请求
}
```

### 4. 注册路由
```go
// router/router.go
func InitRouter() *gin.Engine {
    r := gin.New()
    
    // 使用中间件
    r.Use(gin.Recovery())
    r.Use(middleware.RequestID())
    r.Use(middleware.LoggerMiddleware())
    r.Use(middleware.Cors())
    r.Use(middleware.RateLimit(300, 500))

    // 注册路由
    apiGroup := r.Group("/api")
    {
        apiGroup.POST("/your-endpoint", yourHandler.HandleRequest)
    }

    return r
}
```

## API开发指南

### 1. 创建请求处理器

创建处理请求的Handler结构体：

```go
// api/module/handler.go
package module

import (
    "FastGin/api/common"
    "github.com/gin-gonic/gin"
)

// Handler 处理器结构体
type Handler struct {
    service ServiceInterface
}

// ServiceInterface 服务接口
type ServiceInterface interface {
    DoSomething(params map[string]string) (*common.Response, error)
}

// NewHandler 创建处理器实例
func NewHandler() *Handler {
    return &Handler{
        service: NewService(),
    }
}

// HandleRequest 处理请求方法
func (h *Handler) HandleRequest(c *gin.Context) {
    // 请求参数验证
    var req RequestStruct
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, common.Response{
            Code:    400,
            Message: "参数错误",
            Data:    err.Error(),
        })
        return
    }
    
    // 调用服务层处理业务逻辑
    result, err := h.service.DoSomething(params)
    if err != nil {
        c.JSON(500, common.Response{
            Code:    500,
            Message: "服务错误",
            Data:    nil,
        })
        return
    }
    
    // 返回成功响应
    c.JSON(200, result)
}
```

### 2. 创建Service实现

```go
// apiServer/service_module.go
package apiServer

import (
    "FastGin/api/common"
)

// ModuleService 模块服务结构体
type ModuleService struct {
    // 依赖项
}

// NewModuleService 创建服务实例
func NewModuleService() *ModuleService {
    return &ModuleService{}
}

// DoSomething 实现业务逻辑
func (s *ModuleService) DoSomething(params map[string]string) (*common.Response, error) {
    // 业务逻辑处理
    
    return &common.Response{
        Code:    200,
        Message: "操作成功",
        Data:    result,
    }, nil
}
```

## 中间件配置

### 1. 限流配置
```go
// 配置每秒处理300个请求，突发最大500个
r.Use(middleware.RateLimit(300, 500))
```

### 2. 跨域配置
默认允许所有源，可以通过修改 middleware/cors.go 自定义配置：
```go
func Cors() gin.HandlerFunc {
    config := cors.DefaultConfig()
    config.AllowOrigins = []string{"your-domain.com"}
    // ... 其他配置
}
```

### 3. 日志配置
日志文件位置和格式可在 middleware/logger.go 中配置：
```go
// 错误日志路径
errorLogPath := filepath.Join("logs", "error.log")
```

## 请求/响应格式

### 请求格式
```json
{
    "field1": "value1",
    "field2": "value2"
}
```

### 响应格式
```json
{
    "code": 200,
    "message": "success",
    "data": {
        // 响应数据
    }
}
```

## 配置管理

### 配置文件
项目使用YAML格式的配置文件(settings.yaml)：

```yaml
db:
  host: 127.0.0.1
  port: 3306
  user: root
  password: "password"

server:
  port: 8080
  mode: debug  # debug, release, test
```

### 读取配置
```go
// 配置文件的读取
config.RunSettingFile()
// 读取配置
cfg := core.ReadConfig(config.Options.File)
```

## 测试

### 1. 运行测试脚本
```bash
# 安装依赖
pip install requests colorama

# 运行完整测试
python test_all_apis.py

# 运行单个测试
python test_api.py
```

### 2. 测试覆盖功能
- 接口功能测试
- 参数验证测试
- 错误处理测试
- 限流功能测试
- 并发测试

## 日志示例

### 1. 正常请求日志
```json
{
    "level": "info",
    "message": "请求处理完成",
    "request_id": "xxx",
    "method": "POST",
    "path": "/api/login"
}
```

### 2. 错误日志（logs/error.log）
```json
{
    "t": "2024-01-21T10:30:45.123Z",
    "id": "xxx",
    "tr": "xxx",
    "st": 400,
    "m": "POST",
    "p": "/api/login",
    "e": "参数验证失败"
}
```

## 最佳实践

### 1. 错误处理
- 使用预定义的错误码
- 在处理器中使用 `c.Error()` 记录错误
- 保持错误信息的一致性

### 2. 日志记录
- 使用结构化日志
- 包含足够的上下文信息
- 区分不同级别的日志

### 3. 接口设计
- 使用统一的响应格式
- 实现参数验证
- 提供清晰的错误信息

### 4. 性能优化
- 合理配置限流参数
- 使用适当的缓存策略
- 监控响应时间

## 扩展开发

### 1. 添加新中间件
创建自定义中间件：

```go
// middleware/your_middleware.go
package middleware

import "github.com/gin-gonic/gin"

func YourMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 前置处理
        
        c.Next()
        
        // 后置处理
    }
}
```

注册中间件：

```go
// router/router.go
r.Use(middleware.YourMiddleware())
```

### 2. 扩展配置项
在settings.yaml中添加新的配置项：

```yaml
your_module:
  setting1: value1
  setting2: value2
```

读取配置项：

```go
cfg := core.ReadConfig(config.Options.File)
setting1 := cfg.YourModule.Setting1
```

## 许可证

MIT License 