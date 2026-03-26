package repo

import (
	"context"

	"github.com/speech/fireworks-admin/internal/domain/entity"
)

// TeltentRepo 定义租户数据持久化操作接口。
type TeltentRepo interface {
	// Create 创建新租户并返回创建的实体。
	Create(ctx context.Context, req *entity.CreateTeltentReq) (*entity.Teltent, error)

	// Update 根据ID更新租户信息。
	Update(ctx context.Context, id string, req *entity.UpdateTeltentReq) (*entity.Teltent, error)

	// Delete 根据ID删除租户。
	Delete(ctx context.Context, id string) error

	// GetByID 根据ID查询租户。
	GetByID(ctx context.Context, id string) (*entity.Teltent, error)

	// FindByPage 根据查询条件分页查询租户列表。
	FindByPage(ctx context.Context, query *entity.TeltentQuery) ([]*entity.Teltent, int64, error)
}
