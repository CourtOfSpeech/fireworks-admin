package response

// PageQuery 公共分页查询条件
// 可嵌入到具体的查询条件结构体中
type PageQuery struct {
	Page     int `query:"page"`      // 页码，从1开始
	PageSize int `query:"page_size"` // 每页数量
}

// GetOffset 获取分页偏移量
func (q *PageQuery) GetOffset() int {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 10
	}
	return (q.Page - 1) * q.PageSize
}

// GetLimit 获取分页限制
func (q *PageQuery) GetLimit() int {
	if q.PageSize <= 0 {
		q.PageSize = 10
	}
	return q.PageSize
}

// Normalize 规范化分页参数
func (q *PageQuery) Normalize() {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 10
	}
}

// PageResult 公共分页结果
// T 为列表数据类型
type PageResult[T any] struct {
	List     []T   `json:"list"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
}

// NewPageResult 创建分页结果
func NewPageResult[T any](list []T, total int64, page, pageSize int) *PageResult[T] {
	return &PageResult[T]{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
}
