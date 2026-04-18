package errors

// 业务错误码常量定义
// 错误码采用分段设计，不同范围表示不同类型的错误
const (
	// ErrInternal 系统内部错误码，表示服务器内部错误
	ErrInternal = 100000

	// ErrInvalidParam 请求参数无效错误码，表示客户端传入的参数不符合要求
	ErrInvalidParam = 140001
)
