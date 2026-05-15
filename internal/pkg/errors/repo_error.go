package errors

import (
	"strings"

	entgo "github.com/speech/fireworks-admin/internal/ent"
)

// ConstraintMapping 定义数据库约束冲突到业务错误的映射规则。
// Constraint 是约束名称的关键字（如 "uk_certificate_no"），
// 用于在错误信息中匹配具体的约束冲突类型。
// ErrFactory 是对应的业务错误构造函数。
type ConstraintMapping struct {
	Constraint string
	ErrFactory func(error) error
}

// RepoErrorParser Repository 层错误解析器。
// 将数据库返回的错误（未找到、约束冲突等）转换为对应的业务错误。
// 各业务模块通过 NewRepoErrorParser 创建实例，只需定义约束映射表即可复用解析逻辑。
type RepoErrorParser struct {
	notFoundFactory func(error) error
	constraints     []ConstraintMapping
}

// NewRepoErrorParser 创建 Repository 错误解析器实例。
// notFoundFactory 是未找到记录时的错误构造函数，
// constraints 是约束冲突到业务错误的映射列表。
func NewRepoErrorParser(notFoundFactory func(error) error, constraints []ConstraintMapping) *RepoErrorParser {
	return &RepoErrorParser{
		notFoundFactory: notFoundFactory,
		constraints:     constraints,
	}
}

// Parse 解析 Repository 层返回的错误并转换为业务错误。
// 支持处理：未找到错误、约束冲突错误（唯一键冲突）。
// 如果错误无法识别则返回原始错误。
func (p *RepoErrorParser) Parse(err error) error {
	if err == nil {
		return nil
	}

	if entgo.IsNotFound(err) {
		return p.notFoundFactory(err)
	}

	if entgo.IsConstraintError(err) {
		return p.parseConstraintError(err)
	}

	return err
}

// parseConstraintError 解析数据库约束冲突错误。
// 根据错误信息中的约束名称匹配映射表，返回对应的业务错误。
func (p *RepoErrorParser) parseConstraintError(err error) error {
	errMsg := err.Error()
	for _, m := range p.constraints {
		if strings.Contains(errMsg, m.Constraint) {
			return m.ErrFactory(err)
		}
	}
	return err
}
