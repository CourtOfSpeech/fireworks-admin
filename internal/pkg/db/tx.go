package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/speech/fireworks-admin/internal/ent"
)

// TxOption 现在操作的是标准库的 sql.TxOptions
type TxOption func(*sql.TxOptions)

// WithIsolationLevel 设置事务隔离级别（使用标准库类型）
func WithIsolationLevel(level sql.IsolationLevel) TxOption {
	return func(o *sql.TxOptions) {
		o.Isolation = level
	}
}

// WithReadOnly 设置为只读事务
func WithReadOnly() TxOption {
	return func(o *sql.TxOptions) {
		o.ReadOnly = true
	}
}

// TxManager 提供基础的 DB 访问能力
type TxManager struct {
	client *ent.Client
}

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
