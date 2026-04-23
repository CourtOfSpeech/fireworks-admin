// Package generator 提供 CRUD 代码生成功能。
// 该包定义了代码生成器所需的数据结构和接口，用于根据 Ent Schema 自动生成增删改查代码。
package generator

// SchemaInfo 表示从 Ent Schema 解析出的完整信息。
// 包含 schema 名称、字段列表、索引列表、mixin 信息等元数据，
// 用于后续的代码模板渲染。
type SchemaInfo struct {
	Name         string       // Name schema 名称，如 "Tenant"
	PackageName  string       // PackageName 包名，通常为小写形式，如 "tenant"
	Fields       []FieldInfo  // Fields 字段信息列表
	Indexes      []IndexInfo  // Indexes 索引信息列表
	Mixins       []MixinInfo  // Mixins mixin 信息列表
	HasSoftDelete bool        // HasSoftDelete 是否包含软删除字段
	HasStatus    bool         // HasStatus 是否包含状态字段
	Comment      string       // Comment schema 注释描述
}

// FieldInfo 表示 Schema 中单个字段的详细信息。
// 包含字段名、类型、约束条件、注释等元数据，
// 用于生成 model、dto、repository 等代码。
type FieldInfo struct {
	Name         string // Name 字段名，如 "CertificateNo"
	ColumnName   string // ColumnName 数据库列名，如 "certificate_no"
	GoType       string // GoType Go 语言类型，如 "string"、"int8"、"time.Time"
	IsPrimaryKey bool   // IsPrimaryKey 是否为主键
	IsRequired   bool   // IsRequired 是否必填（非 Optional）
	IsUnique     bool   // IsUnique 是否唯一
	IsOptional   bool   // IsOptional 是否可选
	IsImmutable  bool   // IsImmutable 是否不可变
	IsUUID       bool   // IsUUID 是否为 UUID 类型
	HasDefault   bool   // HasDefault 是否有默认值
	DefaultValue string // DefaultValue 默认值表达式
	MaxLen       int    // MaxLen 最大长度（字符串类型）
	MinVal       int    // MinVal 最小值（数值类型）
	MaxVal       int    // MaxVal 最大值（数值类型）
	Comment      string // Comment 字段注释
	EnumValues   []EnumValue // EnumValues 枚举值列表（枚举类型）
}

// EnumValue 表示枚举字段的单个枚举值。
// 用于生成枚举常量和相关验证逻辑。
type EnumValue struct {
	Name  string // Name 枚举值名称，如 "TenantTypeCompany"
	Value string // Value 枚举值，如 "1"
	Comment string // Comment 枚举值注释
}

// IndexInfo 表示 Schema 中索引的详细信息。
// 包含索引字段、是否唯一、约束名称等元数据，
// 用于生成错误处理和唯一性验证相关代码。
type IndexInfo struct {
	Fields      []string // Fields 索引包含的字段列表
	IsUnique    bool     // IsUnique 是否为唯一索引
	StorageKey  string   // StorageKey 索引约束名称，如 "uk_email"
	WhereClause string   // WhereClause 部分索引的 WHERE 条件
}

// MixinInfo 表示 Schema 使用的 Mixin 信息。
// 用于识别公共字段（如 ID、状态、时间戳等），
// 避免在生成代码时重复处理这些字段。
type MixinInfo struct {
	Name   string   // Name mixin 名称，如 "Id"、"Status"、"SoftDelete"
	Fields []string // Fields mixin 包含的字段名列表
}

// CreateFieldInfo 创建字段信息实例。
// 参数 name 为字段名，goType 为 Go 类型。
// 返回初始化后的字段信息结构体。
func CreateFieldInfo(name, goType string) FieldInfo {
	return FieldInfo{
		Name:       name,
		ColumnName: ToSnakeCase(name),
		GoType:     goType,
	}
}

// CreateIndexInfo 创建索引信息实例。
// 参数 fields 为索引字段列表，isUnique 表示是否唯一索引。
// 返回初始化后的索引信息结构体。
func CreateIndexInfo(fields []string, isUnique bool) IndexInfo {
	return IndexInfo{
		Fields:   fields,
		IsUnique: isUnique,
	}
}

// CreateSchemaInfo 创建 schema 信息实例。
// 参数 name 为 schema 名称。
// 返回初始化后的 schema 信息结构体。
func CreateSchemaInfo(name string) SchemaInfo {
	return SchemaInfo{
		Name:        name,
		PackageName: ToSnakeCase(name),
		Fields:      make([]FieldInfo, 0),
		Indexes:     make([]IndexInfo, 0),
		Mixins:      make([]MixinInfo, 0),
	}
}

// AddField 向 SchemaInfo 添加字段信息。
// 参数 field 为要添加的字段信息。
func (s *SchemaInfo) AddField(field FieldInfo) {
	s.Fields = append(s.Fields, field)
}

// AddIndex 向 SchemaInfo 添加索引信息。
// 参数 index 为要添加的索引信息。
func (s *SchemaInfo) AddIndex(index IndexInfo) {
	s.Indexes = append(s.Indexes, index)
}

// AddMixin 向 SchemaInfo 添加 mixin 信息。
// 参数 mixin 为要添加的 mixin 信息。
func (s *SchemaInfo) AddMixin(mixin MixinInfo) {
	s.Mixins = append(s.Mixins, mixin)
}

// GetUniqueIndexes 获取所有唯一索引。
// 返回唯一索引列表，用于生成唯一性约束错误处理代码。
func (s *SchemaInfo) GetUniqueIndexes() []IndexInfo {
	result := make([]IndexInfo, 0)
	for _, idx := range s.Indexes {
		if idx.IsUnique {
			result = append(result, idx)
		}
	}
	return result
}

// GetRequiredFields 获取所有必填字段。
// 返回必填字段列表，用于生成创建请求 DTO 的验证规则。
func (s *SchemaInfo) GetRequiredFields() []FieldInfo {
	result := make([]FieldInfo, 0)
	for _, f := range s.Fields {
		if f.IsRequired && !f.IsPrimaryKey {
			result = append(result, f)
		}
	}
	return result
}

// GetOptionalFields 获取所有可选字段。
// 返回可选字段列表，用于生成更新请求 DTO。
func (s *SchemaInfo) GetOptionalFields() []FieldInfo {
	result := make([]FieldInfo, 0)
	for _, f := range s.Fields {
		if !f.IsPrimaryKey && !f.IsImmutable {
			result = append(result, f)
		}
	}
	return result
}

// GetQueryFields 获取可用于查询的字段。
// 返回查询字段列表，排除主键和特殊字段。
func (s *SchemaInfo) GetQueryFields() []FieldInfo {
	result := make([]FieldInfo, 0)
	skipFields := map[string]bool{
		"ID":         true,
		"CreatedAt":  true,
		"UpdatedAt":  true,
		"DeletedAt":  true,
	}
	for _, f := range s.Fields {
		if !skipFields[f.Name] {
			result = append(result, f)
		}
	}
	return result
}

// HasMixin 检查是否包含指定的 mixin。
// 参数 name 为 mixin 名称。
// 返回 true 表示包含该 mixin。
func (s *SchemaInfo) HasMixin(name string) bool {
	for _, m := range s.Mixins {
		if m.Name == name {
			return true
		}
	}
	return false
}

// HasField 检查是否包含指定的字段。
// 参数 name 为字段名称。
// 返回 true 表示包含该字段。
func (s *SchemaInfo) HasField(name string) bool {
	for _, f := range s.Fields {
		if f.Name == name {
			return true
		}
	}
	return false
}

// ToSnakeCase 将驼峰命名转换为蛇形命名。
// 参数 str 为驼峰命名的字符串。
// 返回转换后的蛇形命名字符串。
func ToSnakeCase(str string) string {
	if str == "" {
		return ""
	}
	
	result := make([]byte, 0, len(str)+5)
	for i := 0; i < len(str); i++ {
		c := str[i]
		if c >= 'A' && c <= 'Z' {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, c+32)
		} else {
			result = append(result, c)
		}
	}
	return string(result)
}

// ToPascalCase 将蛇形命名转换为帕斯卡命名（大驼峰）。
// 参数 str 为蛇形命名的字符串。
// 返回转换后的帕斯卡命名字符串。
func ToPascalCase(str string) string {
	if str == "" {
		return ""
	}
	
	result := make([]byte, 0, len(str))
	upperNext := true
	for i := 0; i < len(str); i++ {
		c := str[i]
		if c == '_' {
			upperNext = true
			continue
		}
		if upperNext {
			if c >= 'a' && c <= 'z' {
				result = append(result, c-32)
			} else {
				result = append(result, c)
			}
			upperNext = false
		} else {
			result = append(result, c)
		}
	}
	return string(result)
}

// ToCamelCase 将蛇形命名转换为小驼峰命名。
// 参数 str 为蛇形命名的字符串。
// 返回转换后的小驼峰命名字符串。
func ToCamelCase(str string) string {
	pascal := ToPascalCase(str)
	if len(pascal) == 0 {
		return ""
	}
	result := []byte(pascal)
	if result[0] >= 'A' && result[0] <= 'Z' {
		result[0] += 32
	}
	return string(result)
}
