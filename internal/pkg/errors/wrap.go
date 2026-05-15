package errors

// WrapRepoError 解析 Repository 层错误并包装为业务错误。
// 如果错误已经是 BizError 则直接返回，否则包装为内部错误。
// 此函数替代各业务模块中重复的 wrapError 函数。
func WrapRepoError(err error, parser *RepoErrorParser) error {
	if err == nil {
		return nil
	}
	parsed := parser.Parse(err)
	if _, ok := parsed.(*BizError); ok {
		return parsed
	}
	return Internal(parsed)
}
