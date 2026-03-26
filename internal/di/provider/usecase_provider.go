package provider

import (
	"github.com/speech/fireworks-admin/internal/domain/repo"
	"github.com/speech/fireworks-admin/internal/usecase"
)

// ProvideTeltentUsecase 提供租户用例实例
// 参数:
//   - teltentRepo: 租户仓库接口实例
// 返回:
//   - *usecase.TeltentUsecase: 租户用例实例
func ProvideTeltentUsecase(teltentRepo repo.TeltentRepo) *usecase.TeltentUsecase {
	return usecase.NewTeltentUsecase(teltentRepo)
}
