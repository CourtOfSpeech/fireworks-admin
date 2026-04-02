package db

import (
	"github.com/speech/fireworks-admin/internal/ent"
	"github.com/speech/fireworks-admin/internal/pkg/config"
)

// NewEntClient 创建 Ent 数据库客户端。
func NewEntClient(cfg *config.Config) (*ent.Client, func(), error) {
	client, err := ent.Open("postgres", cfg.Database.DSN())
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		_ = client.Close()
	}
	return client, cleanup, nil
}
