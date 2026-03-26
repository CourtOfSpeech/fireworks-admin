# Echo v5 Config 配置详解

## 概述

Echo v5 的 `Config` 结构体用于配置 Echo 实例的各种行为。通过 `echo.NewWithConfig()` 方法可以创建自定义配置的 Echo 实例。

## Config 结构体定义

```go
type Config struct {
    Logger              *slog.Logger
    HTTPErrorHandler    HTTPErrorHandler
    Router              *Router
    Filesystem          fs.FS
    Binder              Binder
    Validator           Validator
    Renderer            Renderer
    JSONSerializer      JSONSerializer
    IPExtractor         IPExtractor
    FormParseMaxMemory  int64
}
```

## 配置属性详解

### 1. Logger

**类型**: `*slog.Logger`

**作用**: 应用程序的结构化日志记录器，用于记录应用程序运行时的各种日志信息。

**默认值**: 如果未设置，会创建一个默认的 TextHandler 写入 stdout。

**使用场景**:
- 需要自定义日志格式（JSON、Text）
- 需要控制日志级别
- 需要添加自定义字段（如请求ID、用户ID等）
- 需要将日志输出到特定位置（文件、网络等）

**配置示例**:

```go
package main

import (
    "log/slog"
    "os"
    "time"
    "github.com/labstack/echo/v5"
)

func main() {
    // 示例1: 基本配置
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    
    e := echo.NewWithConfig(echo.Config{
        Logger: logger,
    })
    
    // 示例2: 带自定义配置的 logger
    opts := &slog.HandlerOptions{
        Level:     slog.LevelDebug,
        AddSource: true,
        ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
            if a.Key == slog.TimeKey {
                t := a.Value.Time()
                a.Value = slog.StringValue(t.Format("2006-01-02 15:04:05"))
            }
            return a
        },
    }
    
    logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
    
    e := echo.NewWithConfig(echo.Config{
        Logger: logger,
    })
    
    // 示例3: 使用自定义 logger 包
    logger := logger.NewLogger("debug", "json", true)
    
    e := echo.NewWithConfig(echo.Config{
        Logger: logger,
    })
}
```

**最佳实践**:
- ✅ 在生产环境使用 JSON 格式，便于日志分析
- ✅ 设置合适的日志级别（生产环境用 Info/Error）
- ✅ 添加源码位置信息便于调试
- ✅ 统一使用项目中的 logger 包

---

### 2. HTTPErrorHandler

**类型**: `HTTPErrorHandler`

**作用**: 集中式错误处理器，处理从处理器和中间件返回的错误，将其转换为适当的 HTTP 响应。

**默认值**: 如果未设置，使用 `DefaultHTTPErrorHandler(false)`。

**函数签名**:
```go
type HTTPErrorHandler func(c *Context, err error)
```

**使用场景**:
- 需要自定义错误响应格式
- 需要根据错误类型返回不同的状态码
- 需要记录错误日志
- 需要统一错误处理逻辑

**配置示例**:

```go
package main

import (
    "errors"
    "fmt"
    "log/slog"
    "net/http"
    "github.com/labstack/echo/v5"
)

// 自定义错误类型
type AppError struct {
    Code    int
    Message string
    Err     error
}

func (e *AppError) Error() string {
    return fmt.Sprintf("code=%d, message=%s, error=%v", e.Code, e.Message, e.Err)
}

// 自定义错误处理器
func customHTTPErrorHandler(c *echo.Context, err error) {
    // 检查响应是否已提交
    if resp, uErr := echo.UnwrapResponse(c.Response()); uErr == nil {
        if resp.Committed {
            return
        }
    }
    
    code := http.StatusInternalServerError
    message := "Internal Server Error"
    
    // 处理不同类型的错误
    var appErr *AppError
    if errors.As(err, &appErr) {
        code = appErr.Code
        message = appErr.Message
    }
    
    // 处理 Echo 的 HTTPError
    var he *echo.HTTPError
    if errors.As(err, &he) {
        code = he.Code
        message = he.Message
    }
    
    // 记录错误日志
    slog.Error("HTTP error",
        slog.Int("status", code),
        slog.String("method", c.Request().Method),
        slog.String("path", c.Request().URL.Path),
        slog.Any("error", err),
    )
    
    // 返回 JSON 错误响应
    if c.Request().Method == http.MethodHead {
        _ = c.NoContent(code)
        return
    }
    
    _ = c.JSON(code, map[string]interface{}{
        "success": false,
        "error": map[string]interface{}{
            "code":    code,
            "message": message,
        },
    })
}

func main() {
    e := echo.NewWithConfig(echo.Config{
        HTTPErrorHandler: customHTTPErrorHandler,
    })
    
    // 使用示例
    e.GET("/error", func(c *echo.Context) error {
        return &AppError{
            Code:    400,
            Message: "参数错误",
            Err:     errors.New("invalid parameter"),
        }
    })
    
    e.Start(":8080")
}
```

**最佳实践**:
- ✅ 检查响应是否已提交
- ✅ 根据错误类型返回不同的状态码
- ✅ 记录详细的错误日志
- ✅ 返回统一格式的错误响应
- ✅ 处理 HEAD 请求的特殊情况

---

### 3. Router

**类型**: `*Router`

**作用**: HTTP 请求路由器，负责将 URL 路径映射到对应的处理器。

**默认值**: 如果未设置，使用 `NewRouter()`。

**使用场景**:
- 需要自定义路由匹配规则
- 需要添加路由前缀
- 需要实现动态路由

**配置示例**:

```go
package main

import (
    "github.com/labstack/echo/v5"
)

func main() {
    // 创建自定义路由器
    router := echo.NewRouter()
    
    // 配置 Echo
    e := echo.NewWithConfig(echo.Config{
        Router: router,
    })
    
    // 添加路由
    e.GET("/", func(c *echo.Context) error {
        return c.String(200, "Hello, World!")
    })
    
    e.Start(":8080")
}
```

**最佳实践**:
- ✅ 大多数情况下使用默认路由器即可
- ✅ 如需自定义路由逻辑，可以实现自己的 Router

---

### 4. Filesystem

**类型**: `fs.FS`

**作用**: 文件系统接口，用于提供静态文件服务。支持 `os.DirFS`、`embed.FS` 和自定义实现。

**默认值**: 如果未设置，默认使用当前工作目录。

**使用场景**:
- 提供静态文件服务（HTML、CSS、JS、图片等）
- 嵌入静态文件到二进制文件中
- 从特定目录提供文件

**配置示例**:

```go
package main

import (
    "embed"
    "os"
    "github.com/labstack/echo/v5"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
    // 示例1: 使用 os.DirFS 从目录提供文件
    e := echo.NewWithConfig(echo.Config{
        Filesystem: os.DirFS("./public"),
    })
    
    e.Static("/", "static")
    
    // 示例2: 使用 embed.FS 嵌入文件
    e2 := echo.NewWithConfig(echo.Config{
        Filesystem: staticFiles,
    })
    
    e2.Static("/", "static")
    
    // 示例3: 自定义文件系统
    e3 := echo.NewWithConfig(echo.Config{
        Filesystem: &CustomFS{},
    })
    
    e.Start(":8080")
}

// 自定义文件系统
type CustomFS struct{}

func (c *CustomFS) Open(name string) (fs.File, error) {
    // 自定义文件打开逻辑
    return os.Open(name)
}
```

**最佳实践**:
- ✅ 开发环境使用 `os.DirFS` 便于调试
- ✅ 生产环境使用 `embed.FS` 将文件嵌入二进制
- ✅ 静态文件单独存放，便于管理

---

### 5. Binder

**类型**: `Binder`

**作用**: 处理 HTTP 请求数据到 Go 结构体的自动绑定。支持 JSON、XML、表单数据、查询参数和路径参数。

**默认值**: 如果未设置，使用 `DefaultBinder`。

**使用场景**:
- 需要自定义数据绑定逻辑
- 需要支持额外的数据格式
- 需要修改默认绑定行为

**配置示例**:

```go
package main

import (
    "encoding/json"
    "github.com/labstack/echo/v5"
)

// 自定义 Binder
type CustomBinder struct{}

func (cb *CustomBinder) Bind(i interface{}, c *echo.Context) error {
    // 先尝试 JSON
    if err := json.NewDecoder(c.Request().Body).Decode(i); err != nil {
        // 如果失败，尝试表单数据
        return echo.NewHTTPError(400, "invalid request data")
    }
    return nil
}

func main() {
    // 示例1: 使用默认 Binder
    e := echo.New()
    
    // 示例2: 使用自定义 Binder
    e2 := echo.NewWithConfig(echo.Config{
        Binder: &CustomBinder{},
    })
    
    type UserRequest struct {
        Name  string `json:"name"`
        Email string `json:"email"`
    }
    
    e2.POST("/users", func(c *echo.Context) error {
        var req UserRequest
        if err := c.Bind(&req); err != nil {
            return err
        }
        return c.JSON(200, req)
    })
    
    e.Start(":8080")
}
```

**最佳实践**:
- ✅ 大多数情况下使用默认 Binder 即可
- ✅ 使用结构体标签指定字段名和验证规则
- ✅ 配合 Validator 使用进行数据验证

---

### 6. Validator

**类型**: `Validator`

**作用**: 数据验证器，用于验证绑定到结构体的数据。

**默认值**: 如果未设置，不进行验证。

**使用场景**:
- 需要验证请求数据的格式和内容
- 需要自定义验证规则
- 需要返回详细的验证错误信息

**配置示例**:

```go
package main

import (
    "github.com/go-playground/validator/v10"
    "github.com/labstack/echo/v5"
)

// 自定义验证器
type CustomValidator struct {
    validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
    return cv.validator.Struct(i)
}

func main() {
    e := echo.NewWithConfig(echo.Config{
        Validator: &CustomValidator{
            validator: validator.New(),
        },
    })
    
    type UserRequest struct {
        Name  string `validate:"required,min=2,max=50"`
        Email string `validate:"required,email"`
        Age   int    `validate:"required,gte=0,lte=130"`
    }
    
    e.POST("/users", func(c *echo.Context) error {
        var req UserRequest
        if err := c.Bind(&req); err != nil {
            return err
        }
        
        // 自动验证
        if err := c.Validate(req); err != nil {
            return echo.NewHTTPError(400, err.Error())
        }
        
        return c.JSON(200, req)
    })
    
    e.Start(":8080")
}
```

**最佳实践**:
- ✅ 使用 `go-playground/validator` 进行验证
- ✅ 定义清晰的验证规则
- ✅ 返回友好的错误信息
- ✅ 在项目启动时初始化验证器

---

### 7. Renderer

**类型**: `Renderer`

**作用**: 模板渲染器，用于渲染 HTML 模板。

**默认值**: 如果未设置，不支持模板渲染。

**使用场景**:
- 需要渲染 HTML 页面
- 需要支持多种模板引擎
- 需要自定义模板函数

**配置示例**:

```go
package main

import (
    "html/template"
    "io"
    "github.com/labstack/echo/v5"
)

// 自定义模板渲染器
type TemplateRenderer struct {
    templates *template.Template
}

func (r *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c *echo.Context) error {
    return r.templates.ExecuteTemplate(w, name, data)
}

func main() {
    // 示例1: 使用 html/template
    e := echo.NewWithConfig(echo.Config{
        Renderer: &TemplateRenderer{
            templates: template.Must(template.ParseGlob("views/*.html")),
        },
    })
    
    e.GET("/", func(c *echo.Context) error {
        return c.Render(200, "index.html", map[string]interface{}{
            "title": "首页",
            "user":  "张三",
        })
    })
    
    // 示例2: 使用自定义模板函数
    funcs := template.FuncMap{
        "upper": strings.ToUpper,
        "lower": strings.ToLower,
    }
    
    e2 := echo.NewWithConfig(echo.Config{
        Renderer: &TemplateRenderer{
            templates: template.Must(template.New("").Funcs(funcs).ParseGlob("views/*.html")),
        },
    })
    
    e.Start(":8080")
}
```

**最佳实践**:
- ✅ 在项目启动时预编译所有模板
- ✅ 使用模板继承和布局
- ✅ 添加常用的模板函数
- ✅ 处理模板渲染错误

---

### 8. JSONSerializer

**类型**: `JSONSerializer`

**作用**: JSON 序列化器，处理 JSON 编码和解码。可以替换默认的 `encoding/json` 实现。

**默认值**: 如果未设置，使用 `DefaultJSONSerializer`（基于 `encoding/json`）。

**使用场景**:
- 需要使用更快的 JSON 库（如 `jsoniter`）
- 需要自定义 JSON 序列化行为
- 需要处理特殊的数据类型

**配置示例**:

```go
package main

import (
    "encoding/json"
    "github.com/labstack/echo/v5"
    jsoniter "github.com/json-iterator/go"
)

// 自定义 JSON 序列化器（使用 jsoniter）
type JSONiterSerializer struct{}

func (s *JSONiterSerializer) Serialize(v interface{}) ([]byte, error) {
    return jsoniter.Marshal(v)
}

func (s *JSONiterSerializer) Deserialize(data []byte, v interface{}) error {
    return jsoniter.Unmarshal(data, v)
}

func main() {
    // 示例1: 使用默认 JSON 序列化器
    e := echo.New()
    
    // 示例2: 使用 jsoniter（性能更好）
    e2 := echo.NewWithConfig(echo.Config{
        JSONSerializer: &JSONiterSerializer{},
    })
    
    e2.GET("/json", func(c *echo.Context) error {
        return c.JSON(200, map[string]interface{}{
            "message": "Hello, World!",
            "time":    time.Now(),
        })
    })
    
    e.Start(":8080")
}
```

**最佳实践**:
- ✅ 默认的 `encoding/json` 已经足够好
- ✅ 如需更高性能，可以使用 `jsoniter`
- ✅ 注意不同 JSON 库的兼容性

---

### 9. IPExtractor

**类型**: `IPExtractor`

**作用**: 从请求中提取真实客户端 IP 地址的策略，特别是在代理或负载均衡器后面时很重要。用于限流、访问控制和日志记录。

**默认值**: 如果未设置，会检查 `X-Forwarded-For` 和 `X-Real-IP` 头。

**使用场景**:
- 应用部署在反向代理后面
- 需要获取真实客户端 IP
- 需要实现 IP 限流
- 需要记录访问日志

**配置示例**:

```go
package main

import (
    "net"
    "strings"
    "github.com/labstack/echo/v5"
)

// 自定义 IP 提取器
type CustomIPExtractor struct{}

func (e *CustomIPExtractor) ExtractIP(r *http.Request) (string, error) {
    // 1. 尝试从 X-Real-IP 获取
    if ip := r.Header.Get("X-Real-IP"); ip != "" {
        return ip, nil
    }
    
    // 2. 尝试从 X-Forwarded-For 获取
    if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
        // 取第一个 IP
        ips := strings.Split(xff, ",")
        if len(ips) > 0 {
            return strings.TrimSpace(ips[0]), nil
        }
    }
    
    // 3. 从 RemoteAddr 获取
    ip, _, err := net.SplitHostPort(r.RemoteAddr)
    if err != nil {
        return "", err
    }
    
    return ip, nil
}

func main() {
    // 示例1: 使用默认 IP 提取器
    e := echo.New()
    
    // 示例2: 使用自定义 IP 提取器
    e2 := echo.NewWithConfig(echo.Config{
        IPExtractor: &CustomIPExtractor{},
    })
    
    // IP 限流中间件
    e2.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c *echo.Context) error {
            ip := c.RealIP()
            
            // 检查 IP 是否在黑名单中
            if isBlocked(ip) {
                return echo.NewHTTPError(403, "IP blocked")
            }
            
            return next(c)
        }
    })
    
    e2.GET("/", func(c *echo.Context) error {
        return c.String(200, "Your IP: "+c.RealIP())
    })
    
    e.Start(":8080")
}

func isBlocked(ip string) bool {
    // 实现 IP 黑名单检查
    return false
}
```

**最佳实践**:
- ✅ 在反向代理后面时必须配置
- ✅ 信任特定的代理服务器
- ✅ 验证 IP 地址格式
- ✅ 处理 IPv6 地址

---

### 10. FormParseMaxMemory

**类型**: `int64`

**作用**: 解析 multipart 表单时的默认内存限制。参见 `(*http.Request).ParseMultipartForm`。

**默认值**: 0（使用 Go 的默认值 32MB）。

**使用场景**:
- 处理文件上传
- 需要限制表单数据大小
- 防止内存耗尽攻击

**配置示例**:

```go
package main

import (
    "github.com/labstack/echo/v5"
)

func main() {
    // 示例1: 默认配置（32MB）
    e := echo.New()
    
    // 示例2: 自定义内存限制（10MB）
    e2 := echo.NewWithConfig(echo.Config{
        FormParseMaxMemory: 10 << 20, // 10MB
    })
    
    // 示例3: 更大的限制（100MB）
    e3 := echo.NewWithConfig(echo.Config{
        FormParseMaxMemory: 100 << 20, // 100MB
    })
    
    e3.POST("/upload", func(c *echo.Context) error {
        // 解析 multipart 表单
        form, err := c.MultipartForm()
        if err != nil {
            return err
        }
        
        // 处理文件
        files := form.File["files"]
        for _, file := range files {
            // 保存文件
            src, err := file.Open()
            if err != nil {
                return err
            }
            defer src.Close()
            
            // ... 保存文件逻辑
        }
        
        return c.String(200, "上传成功")
    })
    
    e.Start(":8080")
}
```

**最佳实践**:
- ✅ 根据实际需求设置合理的大小
- ✅ 文件上传使用流式处理，避免内存占用过大
- ✅ 添加文件大小验证
- ✅ 限制上传文件数量

---

## 完整配置示例

### 示例1: 基础配置

```go
package main

import (
    "log/slog"
    "os"
    "github.com/labstack/echo/v5"
    "github.com/go-playground/validator/v10"
)

func main() {
    // 创建 logger
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    
    // 创建 Echo 实例
    e := echo.NewWithConfig(echo.Config{
        Logger:             logger,
        FormParseMaxMemory: 10 << 20, // 10MB
    })
    
    // 添加路由
    e.GET("/", func(c *echo.Context) error {
        return c.String(200, "Hello, World!")
    })
    
    // 启动服务器
    e.Start(":8080")
}
```

### 示例2: 生产环境配置

```go
package main

import (
    "log/slog"
    "os"
    "time"
    "github.com/labstack/echo/v5"
    "github.com/go-playground/validator/v10"
    jsoniter "github.com/json-iterator/go"
)

// 自定义验证器
type CustomValidator struct {
    validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
    return cv.validator.Struct(i)
}

// 自定义 JSON 序列化器
type JSONiterSerializer struct{}

func (s *JSONiterSerializer) Serialize(v interface{}) ([]byte, error) {
    return jsoniter.Marshal(v)
}

func (s *JSONiterSerializer) Deserialize(data []byte, v interface{}) error {
    return jsoniter.Unmarshal(data, v)
}

// 自定义错误处理器
func customHTTPErrorHandler(c *echo.Context, err error) {
    // 错误处理逻辑
    c.JSON(500, map[string]interface{}{
        "success": false,
        "error":   err.Error(),
    })
}

func main() {
    // 创建 logger
    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level:     slog.LevelInfo,
        AddSource: true,
    }))
    
    // 创建 Echo 实例
    e := echo.NewWithConfig(echo.Config{
        Logger:             logger,
        HTTPErrorHandler:   customHTTPErrorHandler,
        Validator:          &CustomValidator{validator: validator.New()},
        JSONSerializer:     &JSONiterSerializer{},
        FormParseMaxMemory: 10 << 20, // 10MB
    })
    
    // 添加中间件
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())
    e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
    
    // 添加路由
    e.GET("/", func(c *echo.Context) error {
        return c.String(200, "Hello, World!")
    })
    
    // 启动服务器
    e.Start(":8080")
}
```

### 示例3: 项目中的实际使用

```go
package http

import (
    "log/slog"
    "os"
    "github.com/labstack/echo/v5"
    "github.com/labstack/echo/v5/middleware"
    "github.com/speech/fireworks-admin/internal/infrastructure/config"
    "github.com/speech/fireworks-admin/pkg/logger"
    "github.com/speech/fireworks-admin/pkg/response"
    "github.com/speech/fireworks-admin/pkg/validate"
)

type Server struct {
    echo   *echo.Echo
    config *config.Config
    log    *slog.Logger
}

func NewServer(cfg *config.Config, log *slog.Logger) (*Server, error) {
    // 创建 Echo 实例
    e := echo.NewWithConfig(echo.Config{
        Logger:           log,
        HTTPErrorHandler: customHTTPErrorHandler,
        Validator:        validate.NewValidator(),
    })
    
    // 添加中间件
    e.Use(middleware.RequestID())
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())
    e.Use(middleware.CORS())
    
    return &Server{
        echo:   e,
        config: cfg,
        log:    log,
    }, nil
}

func (s *Server) Start() error {
    logger.Info("server starting", slog.Int("port", s.config.Server.Port))
    return s.echo.Start(fmt.Sprintf(":%d", s.config.Server.Port))
}

// 自定义错误处理器
func customHTTPErrorHandler(c *echo.Context, err error) {
    // 错误处理逻辑
    logger.Error("HTTP error",
        slog.Int("status", 500),
        slog.String("method", c.Request().Method),
        slog.String("path", c.Request().URL.Path),
        slog.Any("error", err),
    )
    
    _ = response.Error(c, 500, "Internal Server Error")
}
```

---

## 配置最佳实践

### 1. Logger 配置

```go
// ✅ 推荐：使用项目统一的 logger
logger := logger.NewLogger("debug", "json", true)

e := echo.NewWithConfig(echo.Config{
    Logger: logger,
})

// ❌ 不推荐：每个地方都创建新的 logger
e := echo.New()
```

### 2. 错误处理

```go
// ✅ 推荐：统一的错误处理器
func customHTTPErrorHandler(c *echo.Context, err error) {
    // 检查响应是否已提交
    if resp, uErr := echo.UnwrapResponse(c.Response()); uErr == nil {
        if resp.Committed {
            return
        }
    }
    
    // 处理错误
    code := 500
    message := "Internal Server Error"
    
    // 返回统一格式的响应
    _ = c.JSON(code, map[string]interface{}{
        "success": false,
        "error": map[string]interface{}{
            "code":    code,
            "message": message,
        },
    })
}

// ❌ 不推荐：每个处理器都自己处理错误
e.GET("/users", func(c *echo.Context) error {
    if err != nil {
        return c.JSON(500, map[string]interface{}{
            "error": err.Error(),
        })
    }
    // ...
})
```

### 3. 数据验证

```go
// ✅ 推荐：使用验证器
type UserRequest struct {
    Name  string `validate:"required,min=2,max=50"`
    Email string `validate:"required,email"`
    Age   int    `validate:"required,gte=0,lte=130"`
}

e.Validator = &CustomValidator{validator: validator.New()}

e.POST("/users", func(c *echo.Context) error {
    var req UserRequest
    if err := c.Bind(&req); err != nil {
        return err
    }
    
    if err := c.Validate(req); err != nil {
        return err
    }
    
    // 处理请求
    return c.JSON(200, req)
})

// ❌ 不推荐：手动验证
e.POST("/users", func(c *echo.Context) error {
    name := c.FormValue("name")
    if name == "" {
        return errors.New("name is required")
    }
    if len(name) < 2 {
        return errors.New("name too short")
    }
    // ...
})
```

### 4. 文件上传

```go
// ✅ 推荐：设置合理的内存限制
e := echo.NewWithConfig(echo.Config{
    FormParseMaxMemory: 10 << 20, // 10MB
})

e.POST("/upload", func(c *echo.Context) error {
    file, err := c.FormFile("file")
    if err != nil {
        return err
    }
    
    // 检查文件大小
    if file.Size > 10<<20 {
        return echo.NewHTTPError(400, "file too large")
    }
    
    // 处理文件
    // ...
})

// ❌ 不推荐：不限制文件大小
e.POST("/upload", func(c *echo.Context) error {
    file, err := c.FormFile("file")
    // 直接处理，没有大小检查
})
```

---

## 总结

Echo v5 的 Config 提供了丰富的配置选项，可以根据项目需求灵活配置：

| 配置项 | 重要性 | 推荐配置 |
|--------|--------|----------|
| Logger | ⭐⭐⭐⭐⭐ | 使用项目统一的 logger |
| HTTPErrorHandler | ⭐⭐⭐⭐⭐ | 实现统一的错误处理 |
| Validator | ⭐⭐⭐⭐ | 使用 go-playground/validator |
| FormParseMaxMemory | ⭐⭐⭐ | 根据实际需求设置 |
| JSONSerializer | ⭐⭐ | 默认即可，需要性能时用 jsoniter |
| IPExtractor | ⭐⭐ | 反向代理后必须配置 |
| Binder | ⭐ | 默认即可 |
| Renderer | ⭐ | 需要 HTML 时配置 |
| Router | ⭐ | 默认即可 |
| Filesystem | ⭐ | 静态文件服务时配置 |

**核心建议**:
1. ✅ 必须配置 Logger 和 HTTPErrorHandler
2. ✅ 推荐配置 Validator 进行数据验证
3. ✅ 根据实际需求配置其他选项
4. ✅ 保持配置简洁，不要过度配置
5. ✅ 在项目启动时统一初始化所有配置
