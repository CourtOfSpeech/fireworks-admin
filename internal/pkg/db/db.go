// Package db 提供数据库连接和事务管理功能。
// 该包封装了 Ent ORM 客户端的创建、连接池配置以及事务管理器，
// 支持 Wire 依赖注入和生命周期管理。
package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/speech/fireworks-admin/internal/ent"
	"github.com/speech/fireworks-admin/internal/pkg/config"
	"github.com/speech/fireworks-admin/internal/pkg/lifecycle"
	"github.com/speech/fireworks-admin/internal/pkg/logger"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/lib/pq"
)

// NewEntClient 根据配置创建并初始化 Ent 数据库客户端。
// 该函数会打开数据库连接、配置连接池参数，并返回客户端实例及清理函数。
func NewEntClient(lc *lifecycle.Lifecycle, cfg *config.Config) (*ent.Client, error) {
	db, err := sql.Open("postgres", cfg.Database.DSN())
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(defaultMaxOpenConns(cfg.Database.MaxOpenConns))
	db.SetMaxIdleConns(defaultMaxIdleConns(cfg.Database.MaxIdleConns))
	db.SetConnMaxLifetime(defaultConnMaxLifetime(cfg.Database.ConnMaxLifetime))
	db.SetConnMaxIdleTime(defaultConnMaxIdleTime(cfg.Database.ConnMaxIdleTime))

	drv := entsql.OpenDB(dialect.Postgres, db)
	client := ent.NewClient(ent.Driver(drv))

	lc.Append(lifecycle.Hook{
		Name: "Database",
		OnStart: func(ctx context.Context) error {
			return db.PingContext(ctx)
		},
		OnStop: func(ctx context.Context) error {
			logger.Info(ctx, "正在关闭数据库连接池...")
			return client.Close()
		},
	})
	return client, nil
}

// defaultMaxOpenConns 返回最大打开连接数，当配置值为 0 时使用默认值 25。
func defaultMaxOpenConns(v int) int {
	if v <= 0 {
		return 25
	}
	return v
}

// defaultMaxIdleConns 返回最大空闲连接数，当配置值为 0 时使用默认值 5。
func defaultMaxIdleConns(v int) int {
	if v <= 0 {
		return 5
	}
	return v
}

// defaultConnMaxLifetime 返回连接最大存活时间（秒），当配置值为 0 时使用默认值 3600 秒（1小时）。
func defaultConnMaxLifetime(v int) time.Duration {
	if v <= 0 {
		return 3600 * time.Second
	}
	return time.Duration(v) * time.Second
}

// defaultConnMaxIdleTime 返回空闲连接最大存活时间（秒），当配置值为 0 时使用默认值 600 秒（10分钟）。
func defaultConnMaxIdleTime(v int) time.Duration {
	if v <= 0 {
		return 600 * time.Second
	}
	return time.Duration(v) * time.Second
}
