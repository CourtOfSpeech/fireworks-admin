# UUID 与 UUID v7 完整使用指南

## 目录

1. [UUID 概述](#uuid-概述)
2. [UUID 版本对比](#uuid-版本对比)
3. [UUID v7 详解](#uuid-v7-详解)
4. [错误处理](#错误处理)
5. [Ent 框架集成](#ent-框架集成)
6. [完整代码示例](#完整代码示例)
7. [最佳实践](#最佳实践)

---

## UUID 概述

UUID（Universally Unique Identifier，通用唯一识别码）是一个 128 位的标识符，标准格式为 36 个字符（32 个十六进制数字 + 4 个连字符）：

```
550e8400-e29b-41d4-a716-446655440000
```

### UUID 结构

```
| 时间戳部分 | 版本 | 变体 | 随机/节点部分 |
| 123e4567-e89b-12d3-a456-426614174000
           ^    ^
        版本位  变体位
```

---

## UUID 版本对比

### 版本对比表

| 特性 | v1 | v4 | v7 |
|------|----|----|-----|
| **生成方式** | MAC地址 + 时间戳 | 随机数 | 时间戳 + 随机数 |
| **时间排序** | ✅ 是 | ❌ 否 | ✅ 是 |
| **碰撞概率** | 低 | 极低 | 极低 |
| **隐私问题** | ⚠️ 暴露MAC地址 | ✅ 无 | ✅ 无 |
| **数据库索引友好** | ⚠️ 一般 | ❌ 差 | ✅ 优秀 |
| **分布式系统** | ⚠️ 需协调 | ✅ 无需协调 | ✅ 无需协调 |
| **Go 函数** | `uuid.NewUUID()` | `uuid.New()` | `uuid.NewV7()` |
| **返回错误** | ✅ 是 | ❌ 否 | ✅ 是 |

### 各版本适用场景

```go
// v1: 基于时间和 MAC 地址
// 适用场景：需要时间排序，且不关心隐私泄露
// 注意：暴露机器 MAC 地址，有隐私风险
id, err := uuid.NewUUID() // v1

// v4: 纯随机
// 适用场景：不需要排序，只需要唯一性
// 优点：简单、无隐私问题
// 缺点：数据库索引性能差（随机分布导致页分裂）
id := uuid.New() // v4

// v7: 基于时间戳 + 随机数
// 适用场景：需要时间排序、数据库友好、分布式系统
// 优点：时间排序、索引友好、无隐私问题
id, err := uuid.NewV7() // v7
```

---

## UUID v7 详解

### UUID v7 结构

UUID v7 于 2024 年正式成为 RFC 9562 标准，结构如下：

```
|<-------- 48 位时间戳 ------->|<- 4位 ->|<- 12位 ->|<------ 62 位随机 ------->|
|  0  1  2  3  4  5  6  7  8  9  10 11 |  12 13  | 14 15    | 16-31 字节          |
|  unix_ms                          |  ver=7  | rand_a   | rand_b              |
```

### UUID v7 优势

1. **时间排序**：前 48 位是 Unix 毫秒时间戳，天然有序
2. **数据库友好**：顺序插入，减少 B+ 树页分裂
3. **分布式友好**：无需协调即可生成，后 62 位随机数保证唯一性
4. **隐私安全**：不包含 MAC 地址等敏感信息

### 时间排序示例

```go
package main

import (
    "fmt"
    "time"
    "github.com/google/uuid"
)

func main() {
    // 生成多个 UUID v7
    var ids []uuid.UUID
    for i := 0; i < 5; i++ {
        id, _ := uuid.NewV7()
        ids = append(ids, id)
        time.Sleep(10 * time.Millisecond)
    }
    
    // UUID v7 按时间顺序排列
    fmt.Println("UUID v7 列表（按生成顺序）：")
    for i, id := range ids {
        fmt.Printf("%d: %s\n", i+1, id)
    }
    
    // 输出示例：
    // 1: 019449a8-7c3b-7d2e-8f1a-5b3c2d1e0f0a
    // 2: 019449a8-7c45-7d2e-8f1a-5b3c2d1e0f0a
    // 3: 019449a8-7c4f-7d2e-8f1a-5b3c2d1e0f0a
    // 注意：前缀 019449a8-7c3b/7c45/7c4f 是时间戳，递增
}
```

### 从 UUID v7 提取时间戳

```go
package main

import (
    "fmt"
    "time"
    "github.com/google/uuid"
)

// ExtractTimestamp 从 UUID v7 提取时间戳
func ExtractTimestamp(id uuid.UUID) (time.Time, error) {
    // UUID v7 前 48 位是 Unix 毫秒时间戳
    // 提取前 6 个字节
    ts := int64(id[0])<<40 | int64(id[1])<<32 | int64(id[2])<<24 |
          int64(id[3])<<16 | int64(id[4])<<8 | int64(id[5])
    
    return time.UnixMilli(ts), nil
}

func main() {
    id, _ := uuid.NewV7()
    fmt.Printf("UUID: %s\n", id)
    
    t, _ := ExtractTimestamp(id)
    fmt.Printf("生成时间: %s\n", t.Format(time.RFC3339Nano))
}
```

---

## 错误处理

### 为什么 `uuid.NewV7()` 返回错误？

`uuid.NewV7()` 返回错误的原因：

1. **时间戳获取失败**：极少数情况下，系统时间获取可能失败
2. **随机数生成失败**：加密随机数生成器可能失败
3. **API 设计一致性**：Go 标准库风格，可能失败的操作都返回错误

```go
// uuid.NewV7 源码简化版
func NewV7() (UUID, error) {
    // 1. 获取当前时间
    now := timeNow()
    
    // 2. 生成随机数
    var buf [16]byte
    if _, err := rand.Read(buf[:]); err != nil {
        return Nil, err // 随机数生成失败
    }
    
    // 3. 组装 UUID v7
    // ...
    return uuid, nil
}
```

### 错误处理方式

#### 方式一：直接处理错误

```go
package main

import (
    "fmt"
    "log"
    "github.com/google/uuid"
)

func main() {
    id, err := uuid.NewV7()
    if err != nil {
        // 处理错误：记录日志并使用 v4 作为降级方案
        log.Printf("UUID v7 生成失败: %v, 降级使用 v4", err)
        id = uuid.New() // v4 不返回错误
    }
    fmt.Printf("UUID: %s\n", id)
}
```

#### 方式二：封装工具函数

```go
package uuidutil

import (
    "log/slog"
    "github.com/google/uuid"
)

// MustNewV7 生成 UUID v7，失败时 panic
// 仅用于初始化阶段或确定不会失败的场景
func MustNewV7() uuid.UUID {
    id, err := uuid.NewV7()
    if err != nil {
        panic(fmt.Sprintf("uuid.NewV7() failed: %v", err))
    }
    return id
}

// NewV7OrV4 生成 UUID v7，失败时降级为 v4
func NewV7OrV4() uuid.UUID {
    id, err := uuid.NewV7()
    if err != nil {
        slog.Warn("UUID v7 生成失败，降级使用 v4", "error", err)
        return uuid.New()
    }
    return id
}

// NewV7String 生成 UUID v7 字符串
func NewV7String() string {
    return NewV7OrV4().String()
}
```

#### 方式三：全局生成器（推荐）

```go
package uuidutil

import (
    "crypto/rand"
    "sync"
    "time"
    "github.com/google/uuid"
)

// Generator UUID 生成器
type Generator struct {
    mu sync.Mutex
}

var defaultGenerator = &Generator{}

// NewV7 生成 UUID v7，带重试机制
func (g *Generator) NewV7() (uuid.UUID, error) {
    g.mu.Lock()
    defer g.mu.Unlock()
    
    // 最多重试 3 次
    for i := 0; i < 3; i++ {
        id, err := uuid.NewV7()
        if err == nil {
            return id, nil
        }
        // 短暂等待后重试
        time.Sleep(time.Microsecond * 10)
    }
    
    // 重试失败，返回错误
    return uuid.Nil, fmt.Errorf("failed to generate UUID v7 after 3 retries")
}

// NewV7Safe 生成 UUID v7，失败时降级为 v4
func (g *Generator) NewV7Safe() uuid.UUID {
    id, err := g.NewV7()
    if err != nil {
        return uuid.New() // 降级为 v4
    }
    return id
}

// 包级函数
func NewV7() (uuid.UUID, error) {
    return defaultGenerator.NewV7()
}

func NewV7Safe() uuid.UUID {
    return defaultGenerator.NewV7Safe()
}
```

---

## Ent 框架集成

### 问题：Default 函数签名限制

Ent 的 `Default` 函数签名要求返回单个值，不支持返回错误：

```go
// Ent Default 函数签名
field.UUID("id", uuid.UUID{}).
    Default(func() uuid.UUID {
        // 这里只能返回 uuid.UUID，不能返回 error
        return uuid.New() // v4 可以，因为不返回错误
    })

// ❌ 错误：无法直接使用 uuid.NewV7()
field.UUID("id", uuid.UUID{}).
    Default(func() uuid.UUID {
        id, err := uuid.NewV7() // err 无法处理！
        if err != nil {
            // 这里怎么办？不能返回 error
        }
        return id
    })
```

### 解决方案

#### 方案一：封装不返回错误的函数（推荐）

```go
package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/schema/field"
    "github.com/google/uuid"
)

// Teltent 租户实体
type Teltent struct {
    ent.Schema
}

// newUUIDv7 生成 UUID v7，失败时降级为 v4
// 该函数签名符合 Ent Default 要求
func newUUIDv7() uuid.UUID {
    id, err := uuid.NewV7()
    if err != nil {
        // 降级为 v4
        return uuid.New()
    }
    return id
}

func (Teltent) Fields() []ent.Field {
    return []ent.Field{
        field.UUID("id", uuid.UUID{}).
            Default(newUUIDv7). // 使用封装函数
            Unique().
            Immutable(),
        field.String("name").NotEmpty(),
        field.String("email").NotEmpty().Unique(),
        field.String("phone").NotEmpty().Unique(),
        field.Time("created_at").Default(time.Now).Immutable(),
        field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
    }
}
```

#### 方案二：使用包级工具函数

```go
// pkg/uuid/uuid.go
package uuidutil

import "github.com/google/uuid"

// NewV7Safe 生成 UUID v7，失败时降级为 v4
func NewV7Safe() uuid.UUID {
    id, err := uuid.NewV7()
    if err != nil {
        return uuid.New()
    }
    return id
}

// NewV7String 生成 UUID v7 字符串
func NewV7String() string {
    return NewV7Safe().String()
}
```

```go
// internal/infrastructure/persistence/ent/schema/teltent.go
package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/schema/field"
    "github.com/google/uuid"
    uuidutil "github.com/speech/fireworks-admin/pkg/uuid"
)

type Teltent struct {
    ent.Schema
}

func (Teltent) Fields() []ent.Field {
    return []ent.Field{
        field.UUID("id", uuid.UUID{}).
            Default(uuidutil.NewV7Safe). // 使用工具函数
            Unique().
            Immutable(),
        // ... 其他字段
    }
}
```

#### 方案三：使用 Hook 在创建前生成

```go
package schema

import (
    "context"
    "entgo.io/ent"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/mixin"
    "github.com/google/uuid"
)

// UUIDMixin UUID mixin，自动生成 UUID v7
type UUIDMixin struct {
    mixin.Schema
}

func (UUIDMixin) Fields() []ent.Field {
    return []ent.Field{
        field.UUID("id", uuid.UUID{}).
            Unique().
            Immutable(),
    }
}

// Hooks 在创建时自动生成 UUID v7
func (UUIDMixin) Hooks() []ent.Hook {
    return []ent.Hook{
        hook.On(
            func(next ent.Mutator) ent.Mutator {
                return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
                    // 仅在创建时生成
                    if m.Op().Is(ent.OpCreate) {
                        if _, exists := m.Field("id"); !exists {
                            id, err := uuid.NewV7()
                            if err != nil {
                                id = uuid.New() // 降级
                            }
                            m.SetField("id", id)
                        }
                    }
                    return next.Mutate(ctx, m)
                })
            },
            ent.OpCreate,
        ),
    }
}

// 使用 mixin
type Teltent struct {
    ent.Schema
}

func (Teltent) Mixin() []ent.Mixin {
    return []ent.Mixin{
        UUIDMixin{},
    }
}
```

---

## 完整代码示例

### 项目 UUID 工具包

```go
// pkg/uuid/uuid.go
package uuidutil

import (
    "fmt"
    "log/slog"
    "time"

    "github.com/google/uuid"
)

// MustNewV7 生成 UUID v7，失败时 panic
// 仅用于初始化阶段
func MustNewV7() uuid.UUID {
    id, err := uuid.NewV7()
    if err != nil {
        panic(fmt.Sprintf("uuid.NewV7() failed: %v", err))
    }
    return id
}

// NewV7Safe 生成 UUID v7，失败时降级为 v4
// 推荐用于业务代码
func NewV7Safe() uuid.UUID {
    id, err := uuid.NewV7()
    if err != nil {
        slog.Warn("UUID v7 生成失败，降级使用 v4", "error", err)
        return uuid.New()
    }
    return id
}

// NewV7String 生成 UUID v7 字符串
func NewV7String() string {
    return NewV7Safe().String()
}

// ExtractTimestamp 从 UUID v7 提取时间戳
func ExtractTimestamp(id uuid.UUID) (time.Time, error) {
    if id.Version() != 7 {
        return time.Time{}, fmt.Errorf("not a UUID v7")
    }
    
    // 提取前 48 位时间戳
    ts := int64(id[0])<<40 | int64(id[1])<<32 | int64(id[2])<<24 |
          int64(id[3])<<16 | int64(id[4])<<8 | int64(id[5])
    
    return time.UnixMilli(ts), nil
}

// Parse 解析 UUID 字符串
func Parse(s string) (uuid.UUID, error) {
    return uuid.Parse(s)
}

// MustParse 解析 UUID 字符串，失败时 panic
func MustParse(s string) uuid.UUID {
    return uuid.MustParse(s)
}
```

### Ent Schema 完整示例

```go
// internal/infrastructure/persistence/ent/schema/teltent.go
package schema

import (
    "time"

    "entgo.io/ent"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
    "github.com/google/uuid"
    uuidutil "github.com/speech/fireworks-admin/pkg/uuid"
)

// Teltent 租户实体
type Teltent struct {
    ent.Schema
}

// Fields 定义字段
func (Teltent) Fields() []ent.Field {
    return []ent.Field{
        // 主键：UUID v7
        field.UUID("id", uuid.UUID{}).
            Default(uuidutil.NewV7Safe).
            Unique().
            Immutable().
            Comment("主键，UUID v7"),
        
        // 租户名称
        field.String("name").
            NotEmpty().
            MaxLen(100).
            Comment("租户名称"),
        
        // 邮箱
        field.String("email").
            NotEmpty().
            MaxLen(255).
            Unique().
            Comment("邮箱地址"),
        
        // 手机号
        field.String("phone").
            NotEmpty().
            MaxLen(20).
            Unique().
            Comment("手机号码"),
        
        // 创建时间
        field.Time("created_at").
            Default(time.Now).
            Immutable().
            Comment("创建时间"),
        
        // 更新时间
        field.Time("updated_at").
            Default(time.Now).
            UpdateDefault(time.Now).
            Comment("更新时间"),
    }
}

// Indexes 定义索引
func (Teltent) Indexes() []ent.Index {
    return []ent.Index{
        // UUID v7 主键天然有序，索引性能优秀
        field.Index("id"),
        field.Index("email").Unique(),
        field.Index("phone").Unique(),
        // 时间范围查询索引
        field.Index("created_at"),
    }
}
```

### 业务代码使用示例

```go
// internal/usecase/teltent_usecase.go
package usecase

import (
    "context"
    
    "github.com/google/uuid"
    uuidutil "github.com/speech/fireworks-admin/pkg/uuid"
)

// Create 创建租户（手动指定 ID）
func (u *teltentUsecase) Create(ctx context.Context, req *entity.CreateTeltentReq) (*entity.Teltent, error) {
    // 方式一：使用工具函数生成 UUID v7
    id := uuidutil.NewV7Safe()
    
    // 方式二：直接使用（需要处理错误）
    // id, err := uuid.NewV7()
    // if err != nil {
    //     return nil, fmt.Errorf("generate uuid: %w", err)
    // }
    
    return u.repo.Create(ctx, id, req)
}

// GetByID 根据 ID 查询，同时提取创建时间
func (u *teltentUsecase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Teltent, error) {
    teltent, err := u.repo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // 从 UUID v7 提取创建时间（可选）
    if createdTime, err := uuidutil.ExtractTimestamp(id); err == nil {
        slog.Debug("UUID 创建时间", "id", id, "extracted_time", createdTime)
    }
    
    return teltent, nil
}
```

---

## 最佳实践

### 1. 选择正确的 UUID 版本

| 场景 | 推荐版本 | 原因 |
|------|----------|------|
| 主键（数据库） | v7 | 时间排序，索引友好 |
| 分布式系统 ID | v7 | 无需协调，有序 |
| 会话 Token | v4 | 随机性高，不可预测 |
| 临时 ID | v4 | 简单快速 |
| 需要追溯创建时间 | v7 | 可提取时间戳 |

### 2. 错误处理策略

```go
// ✅ 推荐：降级策略
func generateID() uuid.UUID {
    id, err := uuid.NewV7()
    if err != nil {
        log.Printf("UUID v7 failed: %v, using v4", err)
        return uuid.New()
    }
    return id
}

// ❌ 不推荐：直接 panic（除非初始化阶段）
func generateID() uuid.UUID {
    id, err := uuid.NewV7()
    if err != nil {
        panic(err) // 业务代码不应 panic
    }
    return id
}

// ✅ 推荐：返回错误让调用者处理
func createEntity() (*Entity, error) {
    id, err := uuid.NewV7()
    if err != nil {
        return nil, fmt.Errorf("generate id: %w", err)
    }
    // ...
}
```

### 3. Ent Schema 最佳实践

```go
// ✅ 推荐：使用工具函数
field.UUID("id", uuid.UUID{}).
    Default(uuidutil.NewV7Safe)

// ❌ 不推荐：直接调用（无法处理错误）
field.UUID("id", uuid.UUID{}).
    Default(func() uuid.UUID {
        id, _ := uuid.NewV7() // 忽略错误
        return id
    })
```

### 4. 数据库索引优化

```sql
-- UUID v7 作为主键，索引性能优秀
-- 因为时间有序，B+ 树插入是追加操作

-- 对比 UUID v4：
-- UUID v4 随机分布，导致大量页分裂
-- 索引碎片化，查询性能下降

-- 推荐索引策略
CREATE INDEX idx_created_at ON teltent(created_at);
-- UUID v7 的 id 字段本身就有时间顺序，可以用于范围查询
```

### 5. 迁移现有 v4 到 v7

```go
// 渐进式迁移：新记录使用 v7，旧记录保持 v4
func (r *repo) Create(ctx context.Context, req *CreateReq) (*Entity, error) {
    entity := &Entity{
        ID: uuidutil.NewV7Safe(), // 新记录使用 v7
        // ...
    }
    return r.save(ctx, entity)
}

// 查询时兼容两种版本
func (r *repo) GetByID(ctx context.Context, id uuid.UUID) (*Entity, error) {
    // UUID v4 和 v7 都可以正常查询
    return r.find(ctx, id)
}
```

---

## 总结

| 问题 | 解决方案 |
|------|----------|
| `uuid.NewV7()` 返回错误 | 使用 `NewV7Safe()` 降级为 v4 |
| Ent Default 不支持错误 | 封装不返回错误的函数 |
| 需要时间排序 | 使用 UUID v7 |
| 需要高随机性 | 使用 UUID v4 |

**核心建议**：在项目中统一使用 `pkg/uuid` 工具包，封装好错误处理逻辑，业务代码无需关心底层实现。
