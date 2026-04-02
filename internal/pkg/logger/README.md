# Logger 包使用指南

## 概述

logger 包封装了高性能的日志函数，使用 `slog.LogAttrs` 方法避免反射开销，提供 8 个日志函数。

## 性能优势

- 使用 `LogAttrs` 方法，避免反射开销
- 当日志级别被禁用时，性能提升 **5 倍**
- 零内存分配（在禁用级别时）

## 可用函数

### INFO 级别
```go
logger.Info(msg string, attrs ...slog.Attr)
logger.InfoCtx(ctx context.Context, msg string, attrs ...slog.Attr)
```

### DEBUG 级别
```go
logger.Debug(msg string, attrs ...slog.Attr)
logger.DebugCtx(ctx context.Context, msg string, attrs ...slog.Attr)
```

### WARN 级别
```go
logger.Warn(msg string, attrs ...slog.Attr)
logger.WarnCtx(ctx context.Context, msg string, attrs ...slog.Attr)
```

### ERROR 级别
```go
logger.Error(msg string, attrs ...slog.Attr)
logger.ErrorCtx(ctx context.Context, msg string, attrs ...slog.Attr)
```

## 使用示例

### 1. 初始化日志记录器

```go
package main

import (
    "github.com/speech/fireworks-admin/pkg/logger"
)

func main() {
    // 初始化日志记录器
    logger.NewLogger("debug", "json", true)
    
    // 现在可以使用封装的函数了
    logger.Info("服务启动成功")
}
```

### 2. 基本使用

```go
// 使用默认 context
logger.Info("用户登录",
    slog.Int("user_id", 123),
    slog.String("username", "张三"),
    slog.String("ip", "192.168.1.1"),
)

// 使用指定的 context
ctx := context.Background()
logger.InfoCtx(ctx, "请求处理",
    slog.String("method", "GET"),
    slog.String("path", "/api/users"),
    slog.Int("status", 200),
)
```

### 3. 错误日志

```go
// 记录错误
logger.Error("数据库连接失败",
    slog.Any("error", err),
    slog.String("host", "localhost"),
    slog.Int("port", 5432),
)

// 使用 context
logger.ErrorCtx(ctx, "请求处理失败",
    slog.Any("error", err),
    slog.String("method", "POST"),
    slog.String("path", "/api/users"),
)
```

### 4. 调试日志

```go
logger.Debug("调试信息",
    slog.String("module", "auth"),
    slog.Any("config", map[string]interface{}{
        "timeout": 30,
        "retry":   3,
    }),
)
```

### 5. 警告日志

```go
logger.Warn("配置警告",
    slog.String("key", "jwt_secret"),
    slog.String("issue", "使用默认密钥"),
)
```

### 6. 使用分组

```go
logger.Info("请求处理完成",
    slog.Group("request",
        slog.String("method", "GET"),
        slog.String("path", "/api/users"),
        slog.Int("status", 200),
    ),
    slog.Group("response",
        slog.Int("size", 1024),
        slog.Duration("latency", 50*time.Millisecond),
    ),
)
```

## 常用 slog.Attr 方法

```go
// 基本类型
slog.String("key", "value")
slog.Int("count", 100)
slog.Int64("id", 123456)
slog.Float64("price", 99.99)
slog.Bool("enabled", true)

// 时间类型
slog.Time("created_at", time.Now())
slog.Duration("elapsed", time.Second)

// 通用类型（用于 error、struct 等）
slog.Any("error", err)
slog.Any("user", user)

// 分组
slog.Group("request",
    slog.String("method", "GET"),
    slog.String("path", "/api/users"),
)
```

## 性能对比

| 方式 | 正常输出 | 禁用级别时 |
|------|----------|------------|
| slog.Info + Attr | 1787 ns/op | 165.7 ns/op |
| slog.LogAttrs | 1669 ns/op | **31.66 ns/op** ⭐ |

**结论**：当生产环境禁用 DEBUG/INFO 级别时，性能提升 **5 倍**！

## 最佳实践

### 1. 在项目初始化时调用 NewLogger

```go
func main() {
    // 初始化日志
    logger.NewLogger("debug", "json", true)
    
    // 其他初始化代码...
}
```

### 2. 使用封装的函数而不是 slog 直接调用

```go
// ❌ 不推荐
slog.Info("message", "key", "value")

// ✅ 推荐
logger.Info("message", slog.String("key", "value"))
```

### 3. 在 HTTP 处理器中使用 context

```go
func HandleRequest(c echo.Context) error {
    ctx := c.Request().Context()
    
    logger.InfoCtx(ctx, "处理请求",
        slog.String("method", c.Request().Method),
        slog.String("path", c.Request().URL.Path),
    )
    
    return nil
}
```

### 4. 错误处理时记录错误

```go
if err != nil {
    logger.Error("操作失败",
        slog.Any("error", err),
        slog.String("operation", "database_query"),
    )
    return err
}
```

## 配置说明

```go
logger.NewLogger(
    level,      // 日志级别: "debug", "info", "warn", "error"
    format,     // 输出格式: "json", "text"
    addSource,  // 是否添加源码位置: true, false
)
```

## 特性

- ✅ 高性能：使用 LogAttrs 方法
- ✅ 类型安全：使用 slog.Attr
- ✅ 自动脱敏：password 字段自动替换为 `***REDACTED***`
- ✅ 时间格式化：时间自动格式化为 `2006-01-02 15:04:05`
- ✅ 源码位置：可选添加文件名和行号
- ✅ Context 支持：支持传递 context 用于追踪

## 迁移指南

### 从 slog 迁移

```go
// 之前
slog.Info("user login", "user_id", 123, "ip", "192.168.1.1")

// 现在
logger.Info("user login",
    slog.Int("user_id", 123),
    slog.String("ip", "192.168.1.1"),
)
```

### 从 log.Error 迁移

```go
// 之前
log.Error("failed to connect", "error", err)

// 现在
logger.Error("failed to connect",
    slog.Any("error", err),
)
```
