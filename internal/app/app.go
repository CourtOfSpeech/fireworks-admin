// Package app 提供应用程序的核心初始化和运行机制。
// 该包实现了依赖注入容器、HTTP 服务器、路由注册、健康检查等核心功能，
// 并通过 Wire 框架自动管理依赖关系。
package app

import (
	"log/slog"

	"github.com/speech/fireworks-admin/internal/ent"
	"github.com/speech/fireworks-admin/internal/pkg/config"
	"github.com/speech/fireworks-admin/internal/pkg/lifecycle"
)

// App 应用依赖容器，聚合了应用程序运行所需的所有核心依赖。
// 该结构体由 Wire 框架自动注入，包含配置、日志、数据库客户端、
// 路由注册器、生命周期管理器和 HTTP 服务器等组件。
type App struct {
	Config     *config.Config      // 应用配置
	Logger     *slog.Logger        // 全局日志器
	EntClient  *ent.Client         // Ent ORM 数据库客户端
	Registrars []RouterRegistrar   // 路由注册器列表
	Lifecycle  *lifecycle.Lifecycle // 生命周期管理器
	Server     *Server             // HTTP 服务器
}
