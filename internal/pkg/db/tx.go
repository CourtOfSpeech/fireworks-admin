package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/speech/fireworks-admin/internal/ent"
)

// TxOption 定义事务选项的函数类型，用于配置 sql.TxOptions。
type TxOption func(*sql.TxOptions)

// WithIsolationLevel 设置事务隔离级别（使用标准库类型）。
func WithIsolationLevel(level sql.IsolationLevel) TxOption {
	return func(o *sql.TxOptions) {
		o.Isolation = level
	}
}

// WithReadOnly 设置为只读事务。
func WithReadOnly() TxOption {
	return func(o *sql.TxOptions) {
		o.ReadOnly = true
	}
}

// TxManager 事务管理器，提供数据库访问和事务管理能力。
// 支持从 Context 中提取事务，实现事务传播。
type TxManager struct {
	client *ent.Client // Ent 客户端实例
}

// NewTxManager 创建新的事务管理器实例。
func NewTxManager(client *ent.Client) *TxManager {
	return &TxManager{client: client}
}

// DB 根据 Context 返回对应的 Client。
// 如果 Context 中存在事务，则返回绑定了事务的 Client；否则返回原始 Client。
func (r *TxManager) DB(ctx context.Context) *ent.Client {
	// ent 框架原生支持从 Context 中提取事务
	tx := ent.TxFromContext(ctx)
	if tx != nil {
		return tx.Client()
	}
	return r.client
}

// WithinTx 在事务中执行给定的函数。
// 如果 Context 中已存在事务，则复用该事务（事务传播）；
// 否则开启新事务，并根据函数执行结果提交或回滚。
// 支持通过 TxOption 配置事务隔离级别和只读模式。
func (r *TxManager) WithinTx(ctx context.Context, fn func(ctx context.Context) error, opts ...TxOption) error {
	// 1. 检查是否已经在事务中
	if tx := ent.TxFromContext(ctx); tx != nil {
		// 已经在事务中了，直接执行逻辑，不要开启新事务
		return fn(ctx)
	}

	// 2. 初始化标准库的事务配置
	txOpts := &sql.TxOptions{}
	for _, opt := range opts {
		opt(txOpts)
	}

	// 3. 开启事务 (ent 内部会透传 sql.TxOptions 给底层的驱动)
	tx, err := r.client.BeginTx(ctx, txOpts)
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}

	txCtx := ent.NewTxContext(ctx, tx)

	defer func() {
		if v := recover(); v != nil {
			_ = tx.Rollback()
			panic(v)
		}
	}()

	if err := fn(txCtx); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			return errors.Join(fmt.Errorf("rollback failed"), rerr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}
