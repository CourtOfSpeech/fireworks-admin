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
│   └── wire-best-practices.md    # Wire 依赖注入最佳实践
│
├── internal/                     # 私有应用代码目录
│   ├── app/                      # 应用层：应用初始化和组装
│   │   ├── app.go                # App 结构体定义，包含全局依赖
│   │   ├── server.go             # HTTP 服务器配置和启动
│   │   ├── router.go             # 路由注册，连接各模块 Handler
│   │   ├── wire.go               # Wire 依赖注入定义
│   │   └── wire_gen.go           # Wire 自动生成的依赖注入代码
│   │
│   ├── ent/                      # Ent ORM 生成代码目录
│   │   ├── enttest/              # Ent 测试工具
│   │   │   └── enttest.go        # 测试用的 Ent 客户端
│   │   ├── hook/                 # Ent 钩子
│   │   │   └── hook.go           # 数据库操作钩子定义
│   │   ├── migrate/              # 数据库迁移
│   │   │   ├── migrate.go        # 迁移工具
│   │   │   └── schema.go         # Schema 迁移定义
│   │   ├── predicate/            # 查询谓词
│   │   │   └── predicate.go      # 通用查询条件构建器
│   │   ├── runtime/              # 运行时配置
│   │   │   └── runtime.go        # Ent 运行时设置
│   │   ├── schema/               # 数据库 Schema 定义（手动编写）
│   │   │   ├── mixin/            # 公共 Mixin
│   │   │   │   └── common_mixin.go # 通用字段 Mixin（ID、时间戳、软删除）
│   │   │   └── teltent.go        # 租户表 Schema 定义
│   │   ├── teltent/              # 租户实体生成代码
│   │   │   ├── teltent.go        # 租户实体定义
│   │   │   └── where.go          # 租户查询条件
│   │   ├── client.go             # Ent 客户端
│   │   ├── ent.go                # Ent 核心定义
│   │   ├── generate.go           # 代码生成指令
│   │   ├── mutation.go           # 变更操作定义
│   │   ├── runtime.go            # 运行时配置
│   │   ├── teltent.go            # 租户实体
│   │   ├── teltent_create.go     # 租户创建操作
│   │   ├── teltent_delete.go     # 租户删除操作
│   │   ├── teltent_query.go      # 租户查询操作
│   │   ├── teltent_update.go     # 租户更新操作
│   │   └── tx.go                 # 事务处理
│   │
│   ├── features/                 # 功能模块目录（Feature-First 架构）
│   │   └── teltent/              # 租户模块
│   │       ├── entity.go         # 实体定义、请求/响应结构体
│   │       ├── handler.go        # HTTP 处理器，处理 API 请求
│   │       ├── provider.go       # Wire ProviderSet，提供依赖
│   │       ├── repository.go     # 数据访问层，数据库操作
│   │       └── service.go        # 业务逻辑层
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
│       ├── api/                  # API 响应工具
│       │   ├── page.go           # 分页查询和结果结构体
│       │   └── response.go       # 统一 API 响应封装
│       │
│       ├── config/               # 配置管理
│       │   ├── config.go         # 配置结构体定义
│       │   ├── loader.go         # 配置加载器（支持 TOML）
│       │   └── provider.go       # Wire ProviderSet
│       │
│       ├── db/                   # 数据库连接
│       │   ├── db.go             # Ent 客户端初始化
│       │   └── provider.go       # Wire ProviderSet
│       │
│       ├── idgen/                # ID 生成工具
│       │   └── uuid.go           # UUID v4/v7 生成工具
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
│       ├── teltent.http          # 租户 API 测试用例（HTTP 文件）
│       └── test_api.sh           # API 测试脚本
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
    ├── entity.go      # 数据结构定义
    ├── handler.go     # HTTP 处理器
    ├── service.go     # 业务逻辑
    ├── repository.go  # 数据访问
    └── provider.go    # 依赖注入
```

### 依赖注入（Wire）

使用 Google Wire 进行编译时依赖注入：

- 每个模块/包通过 `provider.go` 定义 `ProviderSet`
- `internal/app/wire.go` 组装所有依赖
- 运行 `wire` 命令生成 `wire_gen.go`

### 数据库（Ent ORM）

- Schema 定义在 `internal/ent/schema/`
- 生成代码在 `internal/ent/` 其他目录
- 使用 Atlas 进行数据库迁移

### HTTP 框架（Echo v5）

- 服务器配置在 `internal/app/server.go`
- 中间件在 `internal/middleware/`
- 路由注册在 `internal/app/router.go`

## 包说明

| 包名 | 路径 | 说明 |
|------|------|------|
| api | `internal/pkg/api` | API 响应封装，分页结构 |
| config | `internal/pkg/config` | 配置加载和管理 |
| db | `internal/pkg/db` | 数据库连接初始化 |
| idgen | `internal/pkg/idgen` | UUID v4/v7 生成工具 |
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
