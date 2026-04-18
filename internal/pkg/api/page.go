// Package api 提供了 API 相关的通用类型和工具函数。
// 包含分页查询、响应格式化等常用功能，用于构建统一的 API 接口。
package api

// PageQuery 公共分页查询条件结构体。
// 可嵌入到具体的查询条件结构体中，提供统一的分页参数处理能力。
type PageQuery struct {
	Page     int `query:"page"`      // Page 页码，从1开始
	PageSize int `query:"page_size"` // PageSize 每页数量
}

// GetOffset 获取分页偏移量。
// 如果页码或每页数量未设置或无效，会自动设置为默认值（页码为1，每页数量为10）。
// 返回计算后的偏移量，用于数据库查询。
func (q *PageQuery) GetOffset() int {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 10
	}
	return (q.Page - 1) * q.PageSize
}

// GetLimit 获取分页限制数量。
// 如果每页数量未设置或无效，会自动设置为默认值10。
// 返回每页数量，用于数据库查询。
func (q *PageQuery) GetLimit() int {
	if q.PageSize <= 0 {
		q.PageSize = 10
	}
	return q.PageSize
}

// Normalize 规范化分页参数。
// 将无效的页码和每页数量设置为默认值（页码为1，每页数量为10）。
func (q *PageQuery) Normalize() {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 10
	}
}

// PageResult 公共分页结果结构体。
// T 为列表数据类型，支持泛型，可适配任意类型的数据列表。
type PageResult[T any] struct {
	List     []T   `json:"list"`      // List 当前页的数据列表
	Total    int64 `json:"total"`     // Total 总记录数
	Page     int   `json:"page"`      // Page 当前页码
	PageSize int   `json:"page_size"` // PageSize 每页数量
}

// NewPageResult 创建分页结果实例。
// list 是当前页的数据列表，total 是总记录数，
// page 是当前页码，pageSize 是每页数量。
func NewPageResult[T any](list []T, total int64, page, pageSize int) *PageResult[T] {
	return &PageResult[T]{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
}
