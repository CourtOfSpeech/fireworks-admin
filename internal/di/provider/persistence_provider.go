package provider

import (
	"github.com/speech/fireworks-admin/internal/domain/repo"
	"github.com/speech/fireworks-admin/internal/infrastructure/persistence/ent"
	persistenceRepo "github.com/speech/fireworks-admin/internal/infrastructure/persistence/repo"
)

// ProvideEntClient 提供 Ent 数据库客户端
// 参数:
//   - dsn: 数据库连接字符串
//
// 返回:
//   - *ent.Client: Ent 客户端实例
//   - func(): 清理函数
//   - error: 连接错误
func ProvideEntClient(dsn string) (*ent.Client, func(), error) {
	client, err := ent.Open("postgres", dsn)
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		_ = client.Close()
	}
	return client, cleanup, nil
}

// ProvideTeltentRepo 提供租户仓库实例
// 参数:
//   - client: Ent 数据库客户端
//
// 返回:
//   - repo.TeltentRepo: 租户仓库接口实例
func ProvideTeltentRepo(client *ent.Client) repo.TeltentRepo {
	return persistenceRepo.NewTeltentEnt(client)
}
