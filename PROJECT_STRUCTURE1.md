# Fireworks-Admin 项目目录结构

> 本文档描述项目的目录结构和各文件/目录的作用

```
fireworks-admin/
├── cmd/                          # 应用程序入口目录
│   └── server/                   # HTTP 服务器入口
│       └── main.go               # 程序主入口，初始化并启动服务
│
├── configs/                      # 配置文件目录
│   ├── config.dev.toml           # 开发环境配置
│   └── config.prod.toml          # 生产环境配置
│
├── docs/                         # 项目文档目录
│   ├── echo-config-guide.md      # Echo 框架配置指南
│   ├── echo-middleware-guide.md  # Echo 中间件使用指南
│   ├── uuid-v7-guide.md          # UUID v7 使用说明
│   ├── wire-best-practices.md    # Wire 依赖注入最佳实践
│   └── optimization-analysis.md  # 性能优化分析文档
│
├── internal/                     # 私有应用代码目录
│   ├── app/                      # 应用层：应用初始化和组装
│   │   ├── app.go                # App 结构体定义，包含全局依赖
│   │   ├── echo.go               # Echo 实例创建、中间件配置、路由注册
│   │   ├── health.go             # 健康检查路由 (/health, /ready)
│   │   ├── router.go             # 路由注册器接口定义
│   │   ├── runner.go             # 应用启动、信号监听、优雅关闭
│   │   ├── server.go             # HTTP 服务器配置和启动
│   │   ├── wire.go               # Wire 依赖注入定义
│   │   ├── wire_gen.go           # Wire 自动生成的依赖注入代码
│   │   └── wire_providers.go     # Wire Provider 集合
│   │
│   ├── ent/                      # Ent ORM 生成代码目录（手动编写部分）
│   │   ├── schema/               # 数据库 Schema 定义（手动编写）
│   │   │   ├── mixin/            # 公共 Mixin
│   │   │   │   └── common_mixin.go  # 通用字段 Mixin（ID、租户ID、状态、时间戳、软删除）
│   │   │   └── tenant.go         # 租户表 Schema 定义
│   │   ├── enttest/              # Ent 测试工具
│   │   ├── hook/                 # Ent 钩子
│   │   ├── migrate/              # 数据库迁移
│   │   ├── predicate/            # 查询谓词
│   │   ├── runtime/              # 运行时配置
│   │   ├── tenant/               # 租户实体生成代码
│   │   ├── client.go             # Ent 客户端
│   │   ├── ent.go                # Ent 核心定义
│   │   ├── generate.go           # 代码生成指令
│   │   ├── mutation.go           # 变更操作定义
│   │   ├── runtime.go            # 运行时配置
│   │   ├── tenant.go             # 租户实体
│   │   ├── tenant_create.go      # 租户创建操作
│   │   ├── tenant_delete.go      # 租户删除操作
│   │   ├── tenant_query.go       # 租户查询操作
│   │   ├── tenant_update.go      # 租户更新操作
│   │   └── tx.go                 # 事务处理
│   │
│   ├── features/                 # 功能模块目录（Feature-First 架构）
│   │   └── tenant/               # 租户模块
│   │       ├── dto.go            # 数据传输对象（请求/响应结构体）
│   │       ├── entity.go         # 实体定义、常量
│   │       ├── errors.go         # 租户模块业务错误定义
│   │       ├── handler.go        # HTTP 处理器，处理 API 请求
│   │       ├── repository.go     # 数据访问层，数据库操作
│   │       ├── service.go        # 业务逻辑层
│   │       └── wire_provider.go  # Wire ProviderSet
│   │
│   ├── middleware/               # HTTP 中间件目录
│   │   ├── cors.go               # CORS 跨域中间件
│   │   ├── gzip.go               # Gzip 压缩中间件
│   │   ├── jwt.go                # JWT 认证中间件
│   │   ├── logger.go             # 请求日志中间件
│   │   ├── recover.go            # Panic 恢复中间件
│   │   ├── request_id.go         # 请求 ID 中间件
│   │   └── timeout.go            # 请求超时中间件
│   │
│   └── pkg/                      # 内部共享工具包目录
│       ├── api/                   # API 响应工具
│       │   ├── page.go           # 分页查询和结果结构体
│       │   └── response.go       # 统一 API 响应封装
│       │
│       ├── config/               # 配置管理
│       │   ├── config.go         # 配置结构体定义
│       │   ├── loader.go         # 配置加载器（支持 TOML）
│       │   └── provider.go       # Wire ProviderSet
│       │
│       ├── db/                   # 数据库连接和事务
│       │   ├── db.go             # Ent 客户端初始化
│       │   ├── provider.go       # Wire ProviderSet
│       │   └── tx.go             # 事务管理器 (WithinTx, TxManager)
│       │
│       ├── errors/               # 业务错误处理
│       │   ├── biz_error.go     # 业务错误类型 (BizError)
│       │   └── codes.go          # 错误码定义
│       │
│       ├── idgen/                # ID 生成工具
│       │   └── uuid.go           # UUID v4/v7 生成工具
│       │
│       ├── lifecycle/            # 生命周期管理
│       │   └── lifecycle.go       # 组件启动/停止钩子 (Hook, Lifecycle)
│       │
│       ├── logger/               # 日志工具
│       │   ├── logger.go         # Slog 日志封装
│       │   └── provider.go       # Wire ProviderSet
│       │
│       └── validator/            # 参数验证
│           └── validator.go      # Echo 请求验证器
│
├── test/                         # 测试目录
│   └── api/                      # API 测试
│       └── tenant.http           # 租户 API 测试用例（HTTP 文件）
│
├── .gitignore                    # Git 忽略文件配置
├── LICENSE                       # 开源许可证
├── Makefile                      # 构建和开发命令
├── PROJECT_STRUCTURE.md          # 项目目录结构文档
├── atlas.hcl                     # Atlas 数据库迁移配置
├── go.mod                        # Go 模块定义
└── go.sum                        # Go 模块校验
```

## 架构说明

### Feature-First 架构

项目采用 Feature-First（功能优先）架构，按业务功能组织代码：

```
features/
└── {module}/
    ├── dto.go         # 数据传输对象（请求/响应）
    ├── entity.go      # 实体定义、常量
    ├── errors.go      # 业务错误定义
    ├── handler.go     # HTTP 处理器
    ├── repository.go  # 数据访问层
    ├── service.go    # 业务逻辑层
    └── wire_provider.go  # 依赖注入
```

### 依赖注入（Wire）

使用 Google Wire 进行编译时依赖注入：

- 每个模块/包通过 `wire_provider.go` 或 `ProviderSet` 定义依赖
- `internal/app/wire.go` 组装所有依赖
- 运行 `wire` 命令生成 `wire_gen.go`

### 数据库（Ent ORM）

- Schema 定义在 `internal/ent/schema/`
  - `tenant.go` - 租户表 Schema
  - `mixin/common_mixin.go` - 公共字段 Mixin（ID、租户ID、状态、创建/更新时间、软删除）
- 生成代码在 `internal/ent/` 其他目录
- 使用 Atlas 进行数据库迁移

### HTTP 框架（Echo v5）

- 服务器配置在 `internal/app/server.go`
- 中间件在 `internal/middleware/`
- 路由注册在 `internal/app/echo.go`

### 错误处理

业务错误使用 `BizError` 类型：

- 错误码定义在 `internal/pkg/errors/codes.go`
- 错误类型在 `internal/pkg/errors/biz_error.go`
- 各模块在 `errors.go` 中定义自己的业务错误

### 事务管理

使用 `TxManager` 进行事务管理：

- 定义在 `internal/pkg/db/tx.go`
- `WithinTx` 方法支持事务嵌套
- `TxManager.DB(ctx)` 自动检测事务上下文

## 包说明

| 包名 | 路径 | 说明 |
|------|------|------|
| app | `internal/app` | 应用初始化、HTTP 服务器、路由组装 |
| ent | `internal/ent` | ORM 生成代码、Schema 定义 |
| features/tenant | `internal/features/tenant` | 租户模块业务逻辑 |
| middleware | `internal/middleware` | HTTP 中间件 |
| api | `internal/pkg/api` | API 响应封装，分页结构 |
| config | `internal/pkg/config` | 配置加载和管理 |
| db | `internal/pkg/db` | 数据库连接、事务管理 |
| errors | `internal/pkg/errors` | 业务错误处理 |
| idgen | `internal/pkg/idgen` | UUID v4/v7 生成工具 |
| lifecycle | `internal/pkg/lifecycle` | 组件生命周期管理 |
| logger | `internal/pkg/logger` | 结构化日志封装 |
| validator | `internal/pkg/validator` | 请求参数验证 |

## 快速命令

```bash
# 生成 Ent 代码
make ent-gen

# 生成 Wire 代码
make wire-gen

# 运行数据库迁移
make migrate

# 启动开发服务器
make run
```
