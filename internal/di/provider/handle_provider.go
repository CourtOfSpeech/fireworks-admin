package provider

import (
	"github.com/speech/fireworks-admin/internal/infrastructure/http/handle"
	"github.com/speech/fireworks-admin/internal/usecase"
)

// ProvideTeltentHandle 提供租户处理器实例
// 参数:
//   - teltentUsecase: 租户用例实例
// 返回:
//   - *handle.TeltentHandle: 租户处理器实例
func ProvideTeltentHandle(teltentUsecase *usecase.TeltentUsecase) *handle.TeltentHandle {
	return handle.NewTeltentHandle(teltentUsecase)
}
