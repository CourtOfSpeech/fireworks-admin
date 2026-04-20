// Package generator 提供 CRUD 代码生成功能。
// 该文件实现了 Ent Schema 解析器，通过 AST 解析提取字段、索引和 Mixin 等元数据。
package generator

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// EntSchemaParser 是 Ent Schema 解析器的实现。
// 使用 Go AST 解析 ent schema 文件，提取字段、索引和 mixin 等元数据信息。
type EntSchemaParser struct {
	schemaDir string // schemaDir ent schema 文件所在目录
}

// NewEntSchemaParser 创建 Ent Schema 解析器实例。
// 参数 schemaDir 为 ent schema 文件所在的目录路径。
// 返回初始化后的解析器实例。
func NewEntSchemaParser(schemaDir string) *EntSchemaParser {
	return &EntSchemaParser{
		schemaDir: schemaDir,
	}
}

// Parse 解析指定名称的 Ent Schema。
// 参数 schemaName 为 schema 名称，如 "Tenant"。
// 返回解析后的 Schema 信息和可能的错误。
func (p *EntSchemaParser) Parse(schemaName string) (*SchemaInfo, error) {
	filePath := filepath.Join(p.schemaDir, strings.ToLower(schemaName)+".go")
	return p.ParseFile(filePath)
}

// ParseFile 解析指定路径的 Ent Schema 文件。
// 参数 filePath 为 schema 文件的绝对路径。
// 返回解析后的 Schema 信息和可能的错误。
func (p *EntSchemaParser) ParseFile(filePath string) (*SchemaInfo, error) {
	if filePath == "" {
		return nil, fmt.Errorf("文件路径不能为空")
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("schema 文件不存在: %s", filePath)
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("解析文件失败: %w", err)
	}

	schemaInfo := CreateSchemaInfo("")

	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			if p.isEntSchema(structType) {
				schemaInfo.Name = typeSpec.Name.Name
				schemaInfo.PackageName = ToSnakeCase(typeSpec.Name.Name)
				break
			}
		}
	}

	if schemaInfo.Name == "" {
		return nil, fmt.Errorf("未找到有效的 Ent Schema 定义")
	}

	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.FuncDecl:
			switch node.Name.Name {
			case "Fields":
				fields := p.parseFieldsMethod(node)
				for _, f := range fields {
					schemaInfo.AddField(f)
				}
			case "Indexes":
				indexes := p.parseIndexesMethod(node)
				for _, idx := range indexes {
					schemaInfo.AddIndex(idx)
				}
			case "Mixin":
				mixins := p.parseMixinMethod(node)
				for _, m := range mixins {
					schemaInfo.AddMixin(m)
					if m.Name == "SoftDelete" {
						schemaInfo.HasSoftDelete = true
					}
					if m.Name == "Status" {
						schemaInfo.HasStatus = true
					}
				}
			}
		}
		return true
	})

	p.addMixinFields(&schemaInfo)

	return &schemaInfo, nil
}

// ListSchemas 列出所有可用的 schema 名称。
// 返回 schema 名称列表和可能的错误。
func (p *EntSchemaParser) ListSchemas() ([]string, error) {
	entries, err := os.ReadDir(p.schemaDir)
	if err != nil {
		return nil, fmt.Errorf("读取 schema 目录失败: %w", err)
	}

	var schemas []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(name, ".go") {
			continue
		}

		if name == "mixin" || strings.HasPrefix(name, "_") {
			continue
		}

		baseName := strings.TrimSuffix(name, ".go")
		schemas = append(schemas, ToPascalCase(baseName))
	}

	return schemas, nil
}

// isEntSchema 检查结构体是否继承自 ent.Schema。
// 参数 structType 为 AST 结构体节点。
// 返回 true 表示该结构体是 ent.Schema 的子类。
func (p *EntSchemaParser) isEntSchema(structType *ast.StructType) bool {
	for _, field := range structType.Fields.List {
		if len(field.Names) == 0 {
			if ident, ok := field.Type.(*ast.Ident); ok {
				if ident.Name == "Schema" {
					return true
				}
			}
			if selectorExpr, ok := field.Type.(*ast.SelectorExpr); ok {
				if ident, ok := selectorExpr.X.(*ast.Ident); ok && ident.Name == "ent" {
					if selectorExpr.Sel.Name == "Schema" {
						return true
					}
				}
			}
		}
	}
	return false
}

// parseFieldsMethod 解析 Fields() 方法，提取字段信息。
// 参数 funcDecl 为 Fields 方法的 AST 节点。
// 返回解析出的字段信息列表。
func (p *EntSchemaParser) parseFieldsMethod(funcDecl *ast.FuncDecl) []FieldInfo {
	var fields []FieldInfo
	seenFields := make(map[string]bool)

	var allChainCalls []*ast.CallExpr

	ast.Inspect(funcDecl, func(n ast.Node) bool {
		callExpr, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if p.isFieldChainCall(callExpr) {
			allChainCalls = append(allChainCalls, callExpr)
		}

		return true
	})

	for _, callExpr := range allChainCalls {
		if p.isOutermostChainCall(callExpr, allChainCalls) {
			fieldInfo := p.parseFieldChain(callExpr)
			if fieldInfo.Name != "" && !seenFields[fieldInfo.Name] {
				fields = append(fields, fieldInfo)
				seenFields[fieldInfo.Name] = true
			}
		}
	}

	return fields
}

// isFieldChainCall 检查调用表达式是否是字段定义链的一部分。
// 参数 callExpr 为调用表达式。
// 返回 true 表示这是字段定义链的一部分。
func (p *EntSchemaParser) isFieldChainCall(callExpr *ast.CallExpr) bool {
	fun, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	if ident, ok := fun.X.(*ast.Ident); ok && ident.Name == "field" {
		return true
	}

	if innerCall, ok := fun.X.(*ast.CallExpr); ok {
		return p.isFieldChainCall(innerCall)
	}

	return false
}

// isOutermostChainCall 检查调用表达式是否是链式调用的最外层调用。
// 参数 callExpr 为调用表达式。
// 参数 allChainCalls 为所有链式调用列表。
// 返回 true 表示这是最外层调用。
func (p *EntSchemaParser) isOutermostChainCall(callExpr *ast.CallExpr, allChainCalls []*ast.CallExpr) bool {
	for _, other := range allChainCalls {
		if other == callExpr {
			continue
		}
		otherFun, ok := other.Fun.(*ast.SelectorExpr)
		if !ok {
			continue
		}
		if innerCall, ok := otherFun.X.(*ast.CallExpr); ok {
			if innerCall == callExpr {
				return false
			}
		}
	}

	return true
}

// parseFieldChain 从链式调用的末端开始解析整个字段定义链。
// 参数 endCall 为链式调用的末端调用表达式。
// 返回解析出的字段信息。
func (p *EntSchemaParser) parseFieldChain(endCall *ast.CallExpr) FieldInfo {
	fieldInfo := FieldInfo{
		IsRequired: true,
	}

	p.traverseFieldChainDown(endCall, &fieldInfo)

	return fieldInfo
}

// traverseFieldChainDown 从外层向内层遍历字段定义链。
// 参数 callExpr 为当前调用表达式。
// 参数 fieldInfo 为要填充的字段信息指针。
func (p *EntSchemaParser) traverseFieldChainDown(callExpr *ast.CallExpr, fieldInfo *FieldInfo) {
	fun, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}

	methodName := fun.Sel.Name
	p.processFieldConstraint(methodName, callExpr, fieldInfo)

	if innerCall, ok := fun.X.(*ast.CallExpr); ok {
		p.traverseFieldChainDown(innerCall, fieldInfo)
	}
}

// processFieldConstraint 处理单个字段约束。
// 参数 methodName 为约束方法名称。
// 参数 callExpr 为约束方法的调用表达式。
// 参数 fieldInfo 为要填充的字段信息指针。
func (p *EntSchemaParser) processFieldConstraint(methodName string, callExpr *ast.CallExpr, fieldInfo *FieldInfo) {
	switch methodName {
	case "NotEmpty":
		fieldInfo.IsRequired = true
	case "Optional":
		fieldInfo.IsOptional = true
		fieldInfo.IsRequired = false
	case "Unique":
		fieldInfo.IsUnique = true
	case "Immutable":
		fieldInfo.IsImmutable = true
	case "Default":
		fieldInfo.HasDefault = true
		if len(callExpr.Args) > 0 {
			fieldInfo.DefaultValue = p.exprToString(callExpr.Args[0])
		}
	case "Min":
		if len(callExpr.Args) > 0 {
			if basicLit, ok := callExpr.Args[0].(*ast.BasicLit); ok {
				val, _ := strconv.Atoi(basicLit.Value)
				fieldInfo.MinVal = val
			}
		}
	case "Max":
		if len(callExpr.Args) > 0 {
			if basicLit, ok := callExpr.Args[0].(*ast.BasicLit); ok {
				val, _ := strconv.Atoi(basicLit.Value)
				fieldInfo.MaxVal = val
			}
		}
	case "MaxLen":
		if len(callExpr.Args) > 0 {
			if basicLit, ok := callExpr.Args[0].(*ast.BasicLit); ok {
				val, _ := strconv.Atoi(basicLit.Value)
				fieldInfo.MaxLen = val
			}
		}
	case "Comment":
		if len(callExpr.Args) > 0 {
			if basicLit, ok := callExpr.Args[0].(*ast.BasicLit); ok {
				fieldInfo.Comment, _ = strconv.Unquote(basicLit.Value)
			}
		}
	case "Enum":
		fieldInfo.GoType = "string"
		fieldInfo.EnumValues = p.parseEnumValues(callExpr)
	case "String", "Int", "Int8", "Int16", "Int32", "Int64",
		"Uint", "Uint8", "Uint16", "Uint32", "Uint64",
		"Float32", "Float64", "Bool", "Time", "JSON", "Bytes", "UUID", "Text":
		fieldInfo.GoType = p.mapFieldTypeToGo(methodName)
		if len(callExpr.Args) > 0 {
			if basicLit, ok := callExpr.Args[0].(*ast.BasicLit); ok {
				name, _ := strconv.Unquote(basicLit.Value)
				fieldInfo.ColumnName = name
				fieldInfo.Name = ToPascalCase(name)
			}
		}
	}
}

// parseEnumValues 解析枚举值列表。
// 参数 callExpr 为 Enum 方法的调用表达式。
// 返回解析出的枚举值列表。
func (p *EntSchemaParser) parseEnumValues(callExpr *ast.CallExpr) []EnumValue {
	var values []EnumValue

	for _, arg := range callExpr.Args {
		if basicLit, ok := arg.(*ast.BasicLit); ok {
			val, _ := strconv.Unquote(basicLit.Value)
			values = append(values, EnumValue{
				Name:  ToPascalCase(val),
				Value: val,
			})
		}
	}

	return values
}

// parseIndexesMethod 解析 Indexes() 方法，提取索引信息。
// 参数 funcDecl 为 Indexes 方法的 AST 节点。
// 返回解析出的索引信息列表。
func (p *EntSchemaParser) parseIndexesMethod(funcDecl *ast.FuncDecl) []IndexInfo {
	var indexes []IndexInfo
	var allChainCalls []*ast.CallExpr

	ast.Inspect(funcDecl, func(n ast.Node) bool {
		callExpr, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if p.isIndexChainCall(callExpr) {
			allChainCalls = append(allChainCalls, callExpr)
		}

		return true
	})

	for _, callExpr := range allChainCalls {
		if p.isOutermostIndexChainCall(callExpr, allChainCalls) {
			indexInfo := p.parseIndexChain(callExpr)
			if len(indexInfo.Fields) > 0 {
				indexes = append(indexes, indexInfo)
			}
		}
	}

	return indexes
}

// isIndexChainCall 检查调用表达式是否是索引定义链的一部分。
// 参数 callExpr 为调用表达式。
// 返回 true 表示这是索引定义链的一部分。
func (p *EntSchemaParser) isIndexChainCall(callExpr *ast.CallExpr) bool {
	fun, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	if ident, ok := fun.X.(*ast.Ident); ok && ident.Name == "index" {
		return true
	}

	if innerCall, ok := fun.X.(*ast.CallExpr); ok {
		return p.isIndexChainCall(innerCall)
	}

	return false
}

// isOutermostIndexChainCall 检查调用表达式是否是索引链式调用的最外层调用。
// 参数 callExpr 为调用表达式。
// 参数 allChainCalls 为所有链式调用列表。
// 返回 true 表示这是最外层调用。
func (p *EntSchemaParser) isOutermostIndexChainCall(callExpr *ast.CallExpr, allChainCalls []*ast.CallExpr) bool {
	for _, other := range allChainCalls {
		if other == callExpr {
			continue
		}
		otherFun, ok := other.Fun.(*ast.SelectorExpr)
		if !ok {
			continue
		}
		if innerCall, ok := otherFun.X.(*ast.CallExpr); ok {
			if innerCall == callExpr {
				return false
			}
		}
	}

	return true
}

// parseIndexChain 从链式调用的末端开始解析整个索引定义链。
// 参数 endCall 为链式调用的末端调用表达式。
// 返回解析出的索引信息。
func (p *EntSchemaParser) parseIndexChain(endCall *ast.CallExpr) IndexInfo {
	indexInfo := IndexInfo{}

	p.traverseIndexChainDown(endCall, &indexInfo)

	return indexInfo
}

// traverseIndexChainDown 从外层向内层遍历索引定义链。
// 参数 callExpr 为当前调用表达式。
// 参数 indexInfo 为要填充的索引信息指针。
func (p *EntSchemaParser) traverseIndexChainDown(callExpr *ast.CallExpr, indexInfo *IndexInfo) {
	fun, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}

	methodName := fun.Sel.Name
	p.processIndexConstraint(methodName, callExpr, indexInfo)

	if innerCall, ok := fun.X.(*ast.CallExpr); ok {
		p.traverseIndexChainDown(innerCall, indexInfo)
	}
}

// processIndexConstraint 处理单个索引约束。
// 参数 methodName 为约束方法名称。
// 参数 callExpr 为约束方法的调用表达式。
// 参数 indexInfo 为要填充的索引信息指针。
func (p *EntSchemaParser) processIndexConstraint(methodName string, callExpr *ast.CallExpr, indexInfo *IndexInfo) {
	switch methodName {
	case "Fields":
		indexInfo.Fields = p.parseIndexFields(callExpr)
	case "Unique":
		indexInfo.IsUnique = true
	case "StorageKey":
		if len(callExpr.Args) > 0 {
			if basicLit, ok := callExpr.Args[0].(*ast.BasicLit); ok {
				indexInfo.StorageKey, _ = strconv.Unquote(basicLit.Value)
			}
		}
	case "Annotations":
		indexInfo.WhereClause = p.parseIndexAnnotations(callExpr)
	}
}

// parseIndexFields 解析索引字段列表。
// 参数 callExpr 为 Fields 方法的调用表达式。
// 返回解析出的字段名列表。
func (p *EntSchemaParser) parseIndexFields(callExpr *ast.CallExpr) []string {
	var fields []string

	for _, arg := range callExpr.Args {
		if basicLit, ok := arg.(*ast.BasicLit); ok {
			name, _ := strconv.Unquote(basicLit.Value)
			fields = append(fields, name)
		}
	}

	return fields
}

// parseIndexAnnotations 解析索引注解，提取 WHERE 条件。
// 参数 callExpr 为 Annotations 方法的调用表达式。
// 返回解析出的 WHERE 条件字符串。
func (p *EntSchemaParser) parseIndexAnnotations(callExpr *ast.CallExpr) string {
	var whereClause string

	for _, arg := range callExpr.Args {
		p.extractWhereClause(arg, &whereClause)
	}

	return whereClause
}

// extractWhereClause 从注解参数中提取 WHERE 条件。
// 参数 expr 为注解参数表达式。
// 参数 whereClause 为存储 WHERE 条件的字符串指针。
func (p *EntSchemaParser) extractWhereClause(expr ast.Expr, whereClause *string) {
	switch e := expr.(type) {
	case *ast.CallExpr:
		fun, ok := e.Fun.(*ast.SelectorExpr)
		if !ok {
			return
		}
		if fun.Sel.Name == "IndexWhere" && len(e.Args) > 0 {
			if basicLit, ok := e.Args[0].(*ast.BasicLit); ok {
				*whereClause, _ = strconv.Unquote(basicLit.Value)
			}
		}
	case *ast.CompositeLit:
		for _, elt := range e.Elts {
			p.extractWhereClause(elt, whereClause)
		}
	}
}

// parseMixinMethod 解析 Mixin() 方法，提取 mixin 信息。
// 参数 funcDecl 为 Mixin 方法的 AST 节点。
// 返回解析出的 mixin 信息列表。
func (p *EntSchemaParser) parseMixinMethod(funcDecl *ast.FuncDecl) []MixinInfo {
	var mixins []MixinInfo

	ast.Inspect(funcDecl, func(n ast.Node) bool {
		compositeLit, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}

		for _, elt := range compositeLit.Elts {
			if typeAssert, ok := elt.(*ast.TypeAssertExpr); ok {
				if ident, ok := typeAssert.X.(*ast.CompositeLit); ok {
					mixinInfo := p.parseMixinComposite(ident)
					if mixinInfo.Name != "" {
						mixins = append(mixins, mixinInfo)
					}
				}
			} else if ident, ok := elt.(*ast.Ident); ok {
				mixins = append(mixins, MixinInfo{
					Name:   ident.Name,
					Fields: p.getKnownMixinFields(ident.Name),
				})
			} else if composite, ok := elt.(*ast.CompositeLit); ok {
				mixinInfo := p.parseMixinComposite(composite)
				if mixinInfo.Name != "" {
					mixins = append(mixins, mixinInfo)
				}
			}
		}

		return true
	})

	return mixins
}

// parseMixinComposite 解析 mixin 复合字面量。
// 参数 compositeLit 为 mixin 的复合字面量 AST 节点。
// 返回解析出的 mixin 信息。
func (p *EntSchemaParser) parseMixinComposite(compositeLit *ast.CompositeLit) MixinInfo {
	mixinInfo := MixinInfo{}

	if selectorExpr, ok := compositeLit.Type.(*ast.SelectorExpr); ok {
		mixinInfo.Name = selectorExpr.Sel.Name
	} else if ident, ok := compositeLit.Type.(*ast.Ident); ok {
		mixinInfo.Name = ident.Name
	}

	mixinInfo.Fields = p.getKnownMixinFields(mixinInfo.Name)

	return mixinInfo
}

// getKnownMixinFields 获取已知 mixin 的字段列表。
// 参数 name 为 mixin 名称。
// 返回该 mixin 包含的字段名列表。
func (p *EntSchemaParser) getKnownMixinFields(name string) []string {
	mixinFields := map[string][]string{
		"Id":          {"id"},
		"TenantId":    {"tenant_id"},
		"Status":      {"status"},
		"CreateTime":  {"created_at"},
		"UpdateTime":  {"updated_at"},
		"SoftDelete":  {"deleted_at"},
		"CommonMixin": {"id", "tenant_id", "status", "created_at", "updated_at", "deleted_at"},
	}

	if fields, ok := mixinFields[name]; ok {
		return fields
	}
	return nil
}

// addMixinFields 向 schema 信息添加 mixin 字段。
// 参数 schemaInfo 为要添加字段的 schema 信息。
func (p *EntSchemaParser) addMixinFields(schemaInfo *SchemaInfo) {
	mixinFieldDefs := map[string][]FieldInfo{
		"Id": {
			{Name: "ID", ColumnName: "id", GoType: "string", IsPrimaryKey: true, HasDefault: true, IsImmutable: true, Comment: "主键"},
		},
		"TenantId": {
			{Name: "TenantID", ColumnName: "tenant_id", GoType: "string", HasDefault: true, IsImmutable: true, Comment: "租户ID"},
		},
		"Status": {
			{Name: "Status", ColumnName: "status", GoType: "int8", IsRequired: true, Comment: "状态"},
		},
		"CreateTime": {
			{Name: "CreatedAt", ColumnName: "created_at", GoType: "time.Time", HasDefault: true, IsImmutable: true, Comment: "创建时间"},
		},
		"UpdateTime": {
			{Name: "UpdatedAt", ColumnName: "updated_at", GoType: "time.Time", HasDefault: true, Comment: "更新时间"},
		},
		"SoftDelete": {
			{Name: "DeletedAt", ColumnName: "deleted_at", GoType: "time.Time", IsOptional: true, Comment: "删除时间"},
		},
	}

	for _, m := range schemaInfo.Mixins {
		if fields, ok := mixinFieldDefs[m.Name]; ok {
			for _, f := range fields {
				schemaInfo.AddField(f)
			}
		}
	}
}

// mapFieldTypeToGo 将 ent 字段类型映射到 Go 类型。
// 参数 entType 为 ent 字段类型名称。
// 返回对应的 Go 类型字符串。
func (p *EntSchemaParser) mapFieldTypeToGo(entType string) string {
	typeMap := map[string]string{
		"String":   "string",
		"Int":      "int",
		"Int8":     "int8",
		"Int16":    "int16",
		"Int32":    "int32",
		"Int64":    "int64",
		"Uint":     "uint",
		"Uint8":    "uint8",
		"Uint16":   "uint16",
		"Uint32":   "uint32",
		"Uint64":   "uint64",
		"Float32":  "float32",
		"Float64":  "float64",
		"Bool":     "bool",
		"Time":     "time.Time",
		"JSON":     "json.RawMessage",
		"Bytes":    "[]byte",
		"Enum":     "string",
		"UUID":     "uuid.UUID",
		"Text":     "string",
	}

	if goType, ok := typeMap[entType]; ok {
		return goType
	}
	return "interface{}"
}

// exprToString 将 AST 表达式转换为字符串。
// 参数 expr 为 AST 表达式节点。
// 返回表达式的字符串表示。
func (p *EntSchemaParser) exprToString(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.BasicLit:
		return e.Value
	case *ast.Ident:
		return e.Name
	case *ast.SelectorExpr:
		return p.exprToString(e.X) + "." + e.Sel.Name
	case *ast.CallExpr:
		return p.exprToString(e.Fun) + "()"
	default:
		return ""
	}
}

// GetFieldTypeJSONTag 获取字段类型的 JSON 标签。
// 参数 goType 为 Go 类型字符串。
// 返回对应的 JSON 类型标签。
func GetFieldTypeJSONTag(goType string) string {
	switch goType {
	case "string":
		return "string"
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return "integer"
	case "float32", "float64":
		return "number"
	case "bool":
		return "boolean"
	case "time.Time":
		return "string"
	case "uuid.UUID":
		return "string"
	default:
		return "object"
	}
}

// IsNumericType 检查类型是否为数值类型。
// 参数 goType 为 Go 类型字符串。
// 返回 true 表示是数值类型。
func IsNumericType(goType string) bool {
	numericTypes := map[string]bool{
		"int": true, "int8": true, "int16": true, "int32": true, "int64": true,
		"uint": true, "uint8": true, "uint16": true, "uint32": true, "uint64": true,
		"float32": true, "float64": true,
	}
	return numericTypes[goType]
}

// IsStringType 检查类型是否为字符串类型。
// 参数 goType 为 Go 类型字符串。
// 返回 true 表示是字符串类型。
func IsStringType(goType string) bool {
	return goType == "string" || goType == "[]byte"
}

// IsTimeType 检查类型是否为时间类型。
// 参数 goType 为 Go 类型字符串。
// 返回 true 表示是时间类型。
func IsTimeType(goType string) bool {
	return goType == "time.Time"
}

// IsUUIDType 检查类型是否为 UUID 类型。
// 参数 goType 为 Go 类型字符串。
// 返回 true 表示是 UUID 类型。
func IsUUIDType(goType string) bool {
	return goType == "uuid.UUID"
}
