# Echo 中间件完整指南

## 概述

Echo v5 提供了丰富的内置中间件，用于处理常见的 Web 应用需求，如日志记录、错误恢复、CORS、限流等。本文档详细介绍所有内置中间件的作用、配置和使用方法。

## 当前项目已使用的中间件

```go
CORS
```

根据 `/Users/xinjiang/Desktop/fireworks-admin/internal/infrastructure/http/server.go`，项目当前使用了以下中间件：

```go
e.Use(middleware.RequestID())  // 请求ID
e.Use(middleware.Logger())     // 日志记录
e.Use(middleware.Recover())    // 错误恢复
e.Use(middleware.CORS())       // 跨域资源共享
```

## 中间件分类

### 1. 核心中间件（必须使用）

- RequestID - 请求ID生成
- Logger - 日志记录
- Recover - 错误恢复

### 2. 安全中间件（强烈推荐）

- CORS - 跨域资源共享
- CSRF - 跨站请求伪造保护
- JWT - JWT认证
- BasicAuth - 基础认证
- KeyAuth - API Key认证

### 3. 性能优化中间件（推荐使用）

- Gzip - 响应压缩
- RateLimiter - 请求限流
- RequestLogger - 请求日志

### 4. 功能增强中间件（按需使用）

- BodyLimit - 请求体大小限制
- Timeout - 请求超时
- Proxy - 反向代理
- Rewrite - URL重写
- Static - 静态文件服务

***

## 详细说明

### 1. RequestID - 请求ID生成器

**作用**: 为每个请求生成唯一的请求ID，便于日志追踪和问题排查。

**重要性**: ⭐⭐⭐⭐⭐

**使用场景**:

- 需要追踪请求链路
- 需要在日志中关联同一请求的所有记录
- 分布式系统中追踪请求

**配置示例**:

```go
package main

import (
    "github.com/labstack/echo/v5"
    "github.com/labstack/echo/v5/middleware"
)

func main() {
    e := echo.New()
    
    // 基本使用
    e.Use(middleware.RequestID())
    
    // 自定义配置
    e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
        Generator: func() string {
            return "custom-" + uuid.New().String()
        },
        TargetHeader: "X-Custom-Request-ID",
    }))
    
    e.GET("/", func(c *echo.Context) error {
        // 获取请求ID
        requestID := c.Response().Header().Get(echo.HeaderXRequestID)
        return c.String(200, "Request ID: "+requestID)
    })
    
    e.Start(":8080")
}
```

**配置选项**:

| 选项           | 类型            | 默认值          | 说明           |
| ------------ | ------------- | ------------ | ------------ |
| Generator    | func() string | 随机生成         | 自定义ID生成函数    |
| TargetHeader | string        | X-Request-ID | 响应头中的请求ID字段名 |

**最佳实践**:

- ✅ 必须在所有中间件之前添加
- ✅ 在日志中记录请求ID
- ✅ 将请求ID传递给下游服务

***

### 2. Logger - 日志记录器

**作用**: 记录HTTP请求的详细信息，包括请求方法、路径、状态码、响应时间等。

**重要性**: ⭐⭐⭐⭐⭐

**使用场景**:

- 记录访问日志
- 性能监控
- 问题排查

**配置示例**:

```go
package main

import (
    "github.com/labstack/echo/v5"
    "github.com/labstack/echo/v5/middleware"
)

func main() {
    e := echo.New()
    
    // 基本使用
    e.Use(middleware.Logger())
    
    // 自定义配置
    e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
        Format: `{"time":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}",` +
            `"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
            `"status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}",` +
            `"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n",
    }))
    
    // 使用自定义 logger
    e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
        Skipper: func(c *echo.Context) bool {
            // 跳过健康检查
            return c.Request().URL.Path == "/health"
        },
    }))
    
    e.Start(":8080")
}
```

**配置选项**:

| 选项      | 类型                   | 默认值  | 说明        |
| ------- | -------------------- | ---- | --------- |
| Skipper | func(\*Context) bool | nil  | 跳过日志记录的条件 |
| Format  | string               | 默认格式 | 日志格式模板    |

**日志变量**:

| 变量                     | 说明            |
| ---------------------- | ------------- |
| ${time\_rfc3339}       | RFC3339格式时间   |
| ${time\_rfc3339\_nano} | RFC3339纳秒格式时间 |
| ${id}                  | 请求ID          |
| ${remote\_ip}          | 客户端IP         |
| ${host}                | 主机名           |
| ${method}              | HTTP方法        |
| ${uri}                 | 请求URI         |
| ${path}                | 请求路径          |
| ${protocol}            | 协议版本          |
| ${referer}             | 来源页面          |
| ${user\_agent}         | 用户代理          |
| ${status}              | HTTP状态码       |
| ${error}               | 错误信息          |
| ${latency}             | 响应时间（纳秒）      |
| ${latency\_human}      | 响应时间（人类可读）    |
| ${bytes\_in}           | 请求体大小         |
| ${bytes\_out}          | 响应体大小         |

**最佳实践**:

- ✅ 使用JSON格式便于日志分析
- ✅ 记录请求ID便于追踪
- ✅ 跳过健康检查等频繁请求
- ✅ 记录响应时间监控性能

***

### 3. Recover - 错误恢复

**作用**: 捕获处理器中的 panic，防止应用崩溃，返回500错误。

**重要性**: ⭐⭐⭐⭐⭐

**使用场景**:

- 防止应用因 panic 崩溃
- 记录 panic 错误信息
- 返回友好的错误响应

**配置示例**:

```go
package main

import (
    "github.com/labstack/echo/v5"
    "github.com/labstack/echo/v5/middleware"
)

func main() {
    e := echo.New()
    
    // 基本使用
    e.Use(middleware.Recover())
    
    // 自定义配置
    e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
        Skipper: func(c *echo.Context) bool {
            return false
        },
        LogLevel: 0, // 0=DEBUG, 1=INFO, 2=WARN, 3=ERROR
        LogErrorFunc: func(c *echo.Context, err error, stack []byte) error {
            // 自定义错误处理
            logger.Error("Panic recovered",
                slog.Any("error", err),
                slog.String("stack", string(stack)),
                slog.String("path", c.Request().URL.Path),
            )
            return nil
        },
        DisablePrintStack: false,
        DisableStackAll:   false,
    }))
    
    e.GET("/panic", func(c *echo.Context) error {
        panic("something went wrong")
    })
    
    e.Start(":8080")
}
```

**配置选项**:

| 选项                | 类型                                    | 默认值   | 说明                |
| ----------------- | ------------------------------------- | ----- | ----------------- |
| Skipper           | func(\*Context) bool                  | nil   | 跳过恢复的条件           |
| LogLevel          | int                                   | 0     | 日志级别              |
| LogErrorFunc      | func(\*Context, error, \[]byte) error | nil   | 自定义错误处理函数         |
| DisablePrintStack | bool                                  | false | 禁用打印堆栈            |
| DisableStackAll   | bool                                  | false | 禁用打印所有goroutine堆栈 |

**最佳实践**:

- ✅ 必须在所有中间件之前添加
- ✅ 记录详细的错误信息和堆栈
- ✅ 发送错误通知（邮件、Slack等）
- ✅ 返回友好的错误响应

***

### 4. CORS - 跨域资源共享

**作用**: 处理跨域请求，允许前端应用从不同域名访问API。

**重要性**: ⭐⭐⭐⭐⭐

**使用场景**:

- 前后端分离架构
- API需要被多个前端应用访问
- 开发环境调试

**配置示例**:

```go
package main

import (
    "github.com/labstack/echo/v5"
    "github.com/labstack/echo/v5/middleware"
)

func main() {
    e := echo.New()
    
    // 基本使用（允许所有来源）
    e.Use(middleware.CORS())
    
    // 自定义配置
    e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
        AllowOrigins: []string{
            "https://example.com",
            "https://www.example.com",
            "http://localhost:3000",
        },
        AllowMethods: []string{
            http.MethodGet,
            http.MethodPost,
            http.MethodPut,
            http.MethodDelete,
            http.MethodOptions,
        },
        AllowHeaders: []string{
            echo.HeaderOrigin,
            echo.HeaderContentType,
            echo.HeaderAccept,
            echo.HeaderAuthorization,
            echo.HeaderXRequestID,
        },
        AllowCredentials: true,
        ExposeHeaders: []string{
            echo.HeaderXRequestID,
        },
        MaxAge: 86400, // 24小时
    }))
    
    e.Start(":8080")
}
```

**配置选项**:

| 选项               | 类型        | 默认值     | 说明          |
| ---------------- | --------- | ------- | ----------- |
| AllowOrigins     | \[]string | \["\*"] | 允许的来源       |
| AllowMethods     | \[]string | 所有方法    | 允许的HTTP方法   |
| AllowHeaders     | \[]string | \[]     | 允许的请求头      |
| AllowCredentials | bool      | false   | 是否允许携带凭证    |
| ExposeHeaders    | \[]string | \[]     | 暴露的响应头      |
| MaxAge           | int       | 0       | 预检请求缓存时间（秒） |

**最佳实践**:

- ✅ 生产环境明确指定允许的来源
- ✅ 只允许必要的HTTP方法
- ✅ 设置合理的缓存时间
- ✅ 开发环境可以使用 "\*"

***

### 5. Gzip - 响应压缩

**作用**: 压缩HTTP响应，减少传输数据量，提高性能。

**重要性**: ⭐⭐⭐⭐

**使用场景**:

- 减少网络传输时间
- 节省带宽
- 提高页面加载速度

**配置示例**:

```go
package main

import (
    "github.com/labstack/echo/v5"
    "github.com/labstack/echo/v5/middleware"
)

func main() {
    e := echo.New()
    
    // 基本使用
    e.Use(middleware.Gzip())
    
    // 自定义配置
    e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
        Level: 6, // 压缩级别（1-9）
        Skipper: func(c *echo.Context) bool {
            // 跳过小文件
            return c.Request().URL.Path == "/small"
        },
    }))
    
    e.GET("/large", func(c *echo.Context) error {
        return c.String(200, "large content...")
    })
    
    e.Start(":8080")
}
```

**配置选项**:

| 选项      | 类型                   | 默认值 | 说明              |
| ------- | -------------------- | --- | --------------- |
| Level   | int                  | -1  | 压缩级别（-1=默认，1-9） |
| Skipper | func(\*Context) bool | nil | 跳过压缩的条件         |

**最佳实践**:

- ✅ 压缩级别设置为 6（平衡性能和压缩率）
- ✅ 跳过已压缩的文件（图片、视频等）
- ✅ 跳过小文件（< 1KB）

***

### 6. RateLimiter - 请求限流

**作用**: 限制请求频率，防止滥用和DDoS攻击。

**重要性**: ⭐⭐⭐⭐

**使用场景**:

- 防止API滥用
- 保护后端服务
- 公平分配资源

**配置示例**:

```go
package main

import (
    "time"
    "github.com/labstack/echo/v5"
    "github.com/labstack/echo/v5/middleware"
)

func main() {
    e := echo.New()
    
    // 基本使用（内存存储）
    e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
    
    // 自定义配置
    e.Use(middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
        Skipper: func(c *echo.Context) bool {
            // 跳过健康检查
            return c.Request().URL.Path == "/health"
        },
        Store: middleware.NewRateLimiterMemoryStoreWithConfig(
            middleware.RateLimiterMemoryStoreConfig{
                Rate:      10,      // 每秒10个请求
                Burst:     30,      // 突发30个请求
                ExpiresIn: 3 * time.Minute, // 过期时间
            },
        ),
        IdentifierExtractor: func(c *echo.Context) (string, error) {
            // 基于IP限流
            return c.RealIP(), nil
        },
        ErrorHandler: func(c *echo.Context, err error) error {
            return c.JSON(429, map[string]interface{}{
                "error": "rate limit exceeded",
            })
        },
        DenyHandler: func(c *echo.Context, identifier string, err error) error {
            return c.JSON(429, map[string]interface{}{
                "error": "too many requests",
            })
        },
    }))
    
    e.Start(":8080")
}
```

**配置选项**:

| 选项                  | 类型                                   | 默认值 | 说明      |
| ------------------- | ------------------------------------ | --- | ------- |
| Store               | RateLimiterStore                     | -   | 限流存储    |
| IdentifierExtractor | func(\*Context) (string, error)      | IP  | 标识符提取函数 |
| ErrorHandler        | func(\*Context, error) error         | -   | 错误处理函数  |
| DenyHandler         | func(\*Context, string, error) error | -   | 拒绝处理函数  |

**最佳实践**:

- ✅ 根据业务需求设置合理的限制
- ✅ 区分不同用户/IP的限流策略
- ✅ 返回清晰的错误信息
- ✅ 生产环境使用 Redis 存储

***

### 7. BodyLimit - 请求体大小限制

**作用**: 限制请求体大小，防止大文件上传导致内存耗尽。

**重要性**: ⭐⭐⭐⭐

**使用场景**:

- 限制文件上传大小
- 防止大请求攻击
- 保护服务器资源

**配置示例**:

```go
package main

import (
    "github.com/labstack/echo/v5"
    "github.com/labstack/echo/v5/middleware"
)

func main() {
    e := echo.New()
    
    // 基本使用（全局限制）
    e.Use(middleware.BodyLimit("2M"))
    
    // 自义配置
    e.Use(middleware.BodyLimitWithConfig(middleware.BodyLimitConfig{
        Skipper: func(c *echo.Context) bool {
            // 跳过特定路由
            return c.Request().URL.Path == "/upload/large"
        },
        Limit: "10M", // 10MB
    }))
    
    // 路由级别限制
    e.POST("/upload", uploadHandler, middleware.BodyLimit("5M"))
    
    e.Start(":8080")
}
```

**配置选项**:

| 选项      | 类型                   | 默认值  | 说明            |
| ------- | -------------------- | ---- | ------------- |
| Limit   | string               | "4M" | 大小限制（支持K、M、G） |
| Skipper | func(\*Context) bool | nil  | 跳过限制的条件       |

**最佳实践**:

- ✅ 根据业务需求设置合理的大小
- ✅ 文件上传使用流式处理
- ✅ 返回清晰的错误信息

***

### 8. Timeout - 请求超时

**作用**: 设置请求处理超时时间，防止长时间运行的请求占用资源。

**重要性**: ⭐⭐⭐

**使用场景**:

- 防止请求长时间挂起
- 保护服务器资源
- 提高用户体验

**配置示例**:

```go
package main

import (
    "time"
    "github.com/labstack/echo/v5"
    "github.com/labstack/echo/v5/middleware"
)

func main() {
    e := echo.New()
    
    // 基本使用
    e.Use(middleware.Timeout(30 * time.Second))
    
    // 自定义配置
    e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
        Skipper: func(c *echo.Context) bool {
            // 跳过长时间任务
            return c.Request().URL.Path == "/long-task"
        },
        Timeout: 30 * time.Second,
        ErrorMessage: "Request timeout",
        OnTimeoutRouteErrorHandler: func(err error, c *echo.Context) {
            logger.Error("Request timeout",
                slog.String("path", c.Request().URL.Path),
                slog.Any("error", err),
            )
        },
    }))
    
    e.Start(":8080")
}
```

**配置选项**:

| 选项                         | 类型                     | 默认值 | 说明       |
| -------------------------- | ---------------------- | --- | -------- |
| Timeout                    | time.Duration          | 0   | 超时时间     |
| ErrorMessage               | string                 | ""  | 超时错误消息   |
| OnTimeoutRouteErrorHandler | func(error, \*Context) | nil | 超时错误处理函数 |

**最佳实践**:

- ✅ 设置合理的超时时间
- ✅ 记录超时错误
- ✅ 返回友好的错误信息
- ✅ 长时间任务使用异步处理

***

### 9. JWT - JWT认证

**作用**: 使用JWT令牌进行身份认证。

**重要性**: ⭐⭐⭐⭐

**使用场景**:

- API认证
- 单点登录
- 无状态认证

**配置示例**:

```go
package main

import (
    "github.com/labstack/echo/v5"
    "github.com/labstack/echo/v5/middleware"
    "github.com/golang-jwt/jwt/v5"
)

func main() {
    e := echo.New()
    
    // 公开路由
    e.POST("/login", loginHandler)
    
    // 受保护路由
    r := e.Group("/api")
    r.Use(middleware.JWT([]byte("secret")))
    
    r.GET("/profile", func(c *echo.Context) error {
        user := c.Get("user").(*jwt.Token)
        claims := user.Claims.(jwt.MapClaims)
        name := claims["name"].(string)
        return c.String(200, "Welcome "+name)
    })
    
    // 自定义配置
    r.Use(middleware.JWTWithConfig(middleware.JWTConfig{
        SigningKey: []byte("secret"),
        SigningMethod: "HS256",
        TokenLookup: "header:Authorization:Bearer ",
        AuthScheme: "Bearer",
        Claims: &jwt.MapClaims{},
        Skipper: func(c *echo.Context) bool {
            // 跳过特定路由
            return c.Request().URL.Path == "/api/public"
        },
        ErrorHandler: func(err error) error {
            return echo.NewHTTPError(401, "invalid token")
        },
        SuccessHandler: func(c *echo.Context) {
            // 认证成功处理
        },
    }))
    
    e.Start(":8080")
}
```

**配置选项**:

| 选项             | 类型                   | 默认值                            | 说明        |
| -------------- | -------------------- | ------------------------------ | --------- |
| SigningKey     | interface{}          | -                              | 签名密钥      |
| SigningMethod  | string               | "HS256"                        | 签名方法      |
| TokenLookup    | string               | "header:Authorization:Bearer " | 令牌查找位置    |
| AuthScheme     | string               | "Bearer"                       | 认证方案      |
| Claims         | jwt.Claims           | jwt.MapClaims                  | 自定义Claims |
| Skipper        | func(\*Context) bool | nil                            | 跳过认证的条件   |
| ErrorHandler   | func(error) error    | -                              | 错误处理函数    |
| SuccessHandler | func(\*Context)      | -                              | 成功处理函数    |

**最佳实践**:

- ✅ 使用强密钥
- ✅ 设置合理的过期时间
- ✅ 使用HTTPS传输
- ✅ 定期轮换密钥

***

### 10. BasicAuth - 基础认证

**作用**: 使用用户名和密码进行HTTP基础认证。

**重要性**: ⭐⭐⭐

**使用场景**:

- 简单的API认证
- 内部系统认证
- 临时访问控制

**配置示例**:

```go
package main

import (
    "github.com/labstack/echo/v5"
    "github.com/labstack/echo/v5/middleware"
)

func main() {
    e := echo.New()
    
    // 基本使用
    e.Use(middleware.BasicAuth(func(username, password string, c *echo.Context) (bool, error) {
        // 验证用户名和密码
        if username == "admin" && password == "secret" {
            return true, nil
        }
        return false, nil
    }))
    
    // 自定义配置
    e.Use(middleware.BasicAuthWithConfig(middleware.BasicAuthConfig{
        Skipper: func(c *echo.Context) bool {
            return c.Request().URL.Path == "/health"
        },
        Validator: func(username, password string, c *echo.Context) (bool, error) {
            // 从数据库验证
            user, err := authenticateUser(username, password)
            if err != nil {
                return false, err
            }
            if user != nil {
                c.Set("user", user)
                return true, nil
            }
            return false, nil
        },
        Realm: "Restricted",
    }))
    
    e.Start(":8080")
}
```

**配置选项**:

| 选项        | 类型                                            | 默认值          | 说明      |
| --------- | --------------------------------------------- | ------------ | ------- |
| Validator | func(string, string, \*Context) (bool, error) | -            | 验证函数    |
| Skipper   | func(\*Context) bool                          | nil          | 跳过认证的条件 |
| Realm     | string                                        | "Restricted" | 认证域     |

**最佳实践**:

- ✅ 使用HTTPS传输
- ✅ 密码加密存储
- ✅ 限制登录失败次数
- ✅ 生产环境推荐使用JWT

***

### 11. KeyAuth - API Key认证

**作用**: 使用API Key进行认证。

**重要性**: ⭐⭐⭐

**使用场景**:

- API服务认证
- 服务间调用
- 第三方集成

**配置示例**:

```go
package main

import (
    "github.com/labstack/echo/v5"
    "github.com/labstack/echo/v5/middleware"
)

func main() {
    e := echo.New()
    
    // 基本使用
    e.Use(middleware.KeyAuth(func(key string, c *echo.Context) (bool, error) {
        // 验证API Key
        return key == "valid-api-key", nil
    }))
    
    // 自定义配置
    e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
        Skipper: func(c *echo.Context) bool {
            return c.Request().URL.Path == "/health"
        },
        KeyLookup: "header: X-API-Key",
        AuthScheme: "",
        Validator: func(key string, c *echo.Context) (bool, error) {
            // 从数据库验证
            valid, err := validateAPIKey(key)
            if err != nil {
                return false, err
            }
            if valid {
                c.Set("api_key", key)
                return true, nil
            }
            return false, nil
        },
    }))
    
    e.Start(":8080")
}
```

**配置选项**:

| 选项         | 类型                                    | 默认值                    | 说明      |
| ---------- | ------------------------------------- | ---------------------- | ------- |
| KeyLookup  | string                                | "header:Authorization" | Key查找位置 |
| AuthScheme | string                                | "Bearer"               | 认证方案    |
| Validator  | func(string, \*Context) (bool, error) | -                      | 验证函数    |
| Skipper    | func(\*Context) bool                  | nil                    | 跳过认证的条件 |

**最佳实践**:

- ✅ 使用强随机Key
- ✅ 定期轮换Key
- ✅ 记录Key使用情况
- ✅ 限制Key权限

***

### 12. CSRF - 跨站请求伪造保护

**作用**: 防止CSRF攻击，保护表单提交。

**重要性**: ⭐⭐⭐⭐

**使用场景**:

- 表单提交保护
- 防止恶意网站伪造请求
- Web应用安全

**配置示例**:

```go
package main

import (
    "github.com/labstack/echo/v5"
    "github.com/labstack/echo/v5/middleware"
)

func main() {
    e := echo.New()
    
    // 基本使用
    e.Use(middleware.CSRF())
    
    // 自定义配置
    e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
        Skipper: func(c *echo.Context) bool {
            // 跳过API路由
            return c.Request().URL.Path[:4] == "/api"
        },
        TokenLookup: "form:_csrf",
        ContextKey: "csrf",
        Name: "csrf",
        Cookie: &http.Cookie{
            Name:     "_csrf",
            MaxAge:   86400,
            HttpOnly: true,
            Secure:   true,
            SameSite: http.SameSiteStrictMode,
        },
    }))
    
    e.GET("/form", func(c *echo.Context) error {
        csrf := c.Get("csrf").(string)
        return c.Render(200, "form.html", map[string]interface{}{
            "csrf": csrf,
        })
    })
    
    e.POST("/form", func(c *echo.Context) error {
        return c.String(200, "Form submitted")
    })
    
    e.Start(":8080")
}
```

**配置选项**:

| 选项          | 类型            | 默认值                   | 说明          |
| ----------- | ------------- | --------------------- | ----------- |
| TokenLookup | string        | "header:X-CSRF-Token" | Token查找位置   |
| ContextKey  | string        | "csrf"                | Context中的键名 |
| Name        | string        | "csrf"                | Cookie名称    |
| Cookie      | \*http.Cookie | 默认Cookie              | Cookie配置    |

**最佳实践**:

- ✅ 所有表单提交必须包含CSRF Token
- ✅ API路由可以跳过CSRF
- ✅ 使用HTTPS
- ✅ 设置合理的Cookie属性

***

### 13. Proxy - 反向代理

**作用**: 将请求代理到后端服务。

**重要性**: ⭐⭐⭐

**使用场景**:

- API网关
- 负载均衡
- 服务转发

**配置示例**:

```go
package main

import (
    "github.com/labstack/echo/v5"
    "github.com/labstack/echo/v5/middleware"
)

func main() {
    e := echo.New()
    
    // 基本使用
    e.Use(middleware.Proxy(middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{
        {
            URL: mustParse("http://localhost:8081"),
        },
        {
            URL: mustParse("http://localhost:8082"),
        },
    })))
    
    // 自定义配置
    e.Use(middleware.ProxyWithConfig(middleware.ProxyConfig{
        Skipper: func(c *echo.Context) bool {
            return c.Request().URL.Path == "/health"
        },
        Balancer: middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{
            {
                URL: mustParse("http://localhost:8081"),
                Name: "server1",
            },
            {
                URL: mustParse("http://localhost:8082"),
                Name: "server2",
            },
        }),
        Rewrite: map[string]string{
            "/api/*": "/$1",
        },
    }))
    
    e.Start(":8080")
}

func mustParse(rawURL string) *url.URL {
    u, err := url.Parse(rawURL)
    if err != nil {
        panic(err)
    }
    return u
}
```

**配置选项**:

| 选项       | 类型                   | 默认值 | 说明      |
| -------- | -------------------- | --- | ------- |
| Balancer | Balancer             | -   | 负载均衡器   |
| Rewrite  | map\[string]string   | -   | URL重写规则 |
| Skipper  | func(\*Context) bool | nil | 跳过代理的条件 |

**最佳实践**:

- ✅ 使用健康检查
- ✅ 配置合理的超时时间
- ✅ 记录代理日志
- ✅ 处理错误情况

***

### 14. Rewrite - URL重写

**作用**: 重写URL路径，便于路由管理。

**重要性**: ⭐⭐⭐

**使用场景**:

- URL美化
- 版本控制
- 向后兼容

**配置示例**:

```go
package main

import (
    "github.com/labstack/echo/v5"
    "github.com/labstack/echo/v5/middleware"
)

func main() {
    e := echo.New()
    
    // 基本使用
    e.Use(middleware.Rewrite(map[string]string{
        "/old":              "/new",
        "/api/*":            "/v1/$1",
        "/users/:id/*":      "/user/$1/$2",
    }))
    
    // 自定义配置
    e.Use(middleware.RewriteWithConfig(middleware.RewriteConfig{
        Skipper: func(c *echo.Context) bool {
            return false
        },
        Rules: map[string]string{
            "^/api/v1/(.*)": "/v1/$1",
            "^/api/v2/(.*)": "/v2/$1",
        },
    }))
    
    e.Start(":8080")
}
```

**配置选项**:

| 选项    | 类型                 | 默认值 | 说明   |
| ----- | ------------------ | --- | ---- |
| Rules | map\[string]string | -   | 重写规则 |

**最佳实践**:

- ✅ 使用正则表达式匹配
- ✅ 保持URL简洁
- ✅ 记录重写日志
- ✅ 测试重写规则

***

### 15. Static - 静态文件服务

**作用**: 提供静态文件服务。

**重要性**: ⭐⭐⭐

**使用场景**:

- 提供HTML、CSS、JS文件
- 提供图片、视频等资源
- SPA应用部署

**配置示例**:

```go
package main

import (
    "github.com/labstack/echo/v5"
    "github.com/labstack/echo/v5/middleware"
)

func main() {
    e := echo.New()
    
    // 基本使用
    e.Use(middleware.Static("static"))
    
    // 自定义配置
    e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
        Root: "static",
        Skipper: func(c *echo.Context) bool {
            return c.Request().URL.Path[:4] == "/api"
        },
        Index: "index.html",
        HTML5: true, // 支持HTML5 History模式
        Browse: false,
    }))
    
    // 路由级别
    e.Static("/static", "assets")
    e.File("/", "public/index.html")
    
    e.Start(":8080")
}
```

**配置选项**:

| 选项      | 类型                   | 默认值          | 说明                |
| ------- | -------------------- | ------------ | ----------------- |
| Root    | string               | ""           | 静态文件根目录           |
| Skipper | func(\*Context) bool | nil          | 跳过静态服务的条件         |
| Index   | string               | "index.html" | 默认文件              |
| HTML5   | bool                 | false        | 支持HTML5 History模式 |
| Browse  | bool                 | false        | 是否允许浏览目录          |

**最佳实践**:

- ✅ 设置缓存头
- ✅ 使用CDN加速
- ✅ 压缩静态文件
- ✅ 使用版本号管理

***

## 中间件组合最佳实践

### 生产环境推荐配置

```go
package main

import (
    "time"
    "github.com/labstack/echo/v5"
    "github.com/labstack/echo/v5/middleware"
)

func main() {
    e := echo.New()
    
    // 1. 核心中间件（必须）
    e.Use(middleware.RequestID())
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())
    
    // 2. 安全中间件（强烈推荐）
    e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
        AllowOrigins: []string{"https://example.com"},
        AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
        AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
        AllowCredentials: true,
        MaxAge: 86400,
    }))
    
    // 3. 性能优化中间件（推荐）
    e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
        Level: 6,
    }))
    
    e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
    
    // 4. 功能增强中间件（按需）
    e.Use(middleware.BodyLimit("10M"))
    e.Use(middleware.Timeout(30 * time.Second))
    
    // 5. 认证中间件（按路由）
    api := e.Group("/api")
    api.Use(middleware.JWT([]byte("secret")))
    
    e.Start(":8080")
}
```

### 中间件执行顺序

中间件按照添加顺序执行，建议顺序：

1. **RequestID** - 生成请求ID
2. **Logger** - 记录请求日志
3. **Recover** - 错误恢复
4. **CORS** - 跨域处理
5. **Gzip** - 响应压缩
6. **RateLimiter** - 限流
7. **BodyLimit** - 请求体限制
8. **Timeout** - 超时控制
9. **认证中间件** - JWT/BasicAuth/KeyAuth

***

## 中间件选择指南

### 按场景选择

| 场景    | 推荐中间件                                                |
| ----- | ---------------------------------------------------- |
| API服务 | RequestID, Logger, Recover, CORS, RateLimiter, JWT   |
| Web应用 | RequestID, Logger, Recover, CORS, CSRF, Gzip, Static |
| 内部服务  | RequestID, Logger, Recover, BasicAuth                |
| 微服务   | RequestID, Logger, Recover, RateLimiter, Proxy       |

### 按重要性选择

| 重要性  | 中间件                                                   |
| ---- | ----------------------------------------------------- |
| 必须使用 | RequestID, Logger, Recover                            |
| 强烈推荐 | CORS, RateLimiter                                     |
| 推荐使用 | Gzip, BodyLimit, Timeout                              |
| 按需使用 | JWT, BasicAuth, KeyAuth, CSRF, Proxy, Rewrite, Static |

***

## 总结

Echo v5 提供了丰富的内置中间件，覆盖了Web应用的常见需求。合理使用中间件可以：

- ✅ 提高应用安全性
- ✅ 改善性能
- ✅ 简化开发
- ✅ 统一处理逻辑

**核心建议**:

1. 必须使用 RequestID、Logger、Recover
2. 根据业务需求选择其他中间件
3. 注意中间件的执行顺序
4. 定期审查和优化中间件配置

