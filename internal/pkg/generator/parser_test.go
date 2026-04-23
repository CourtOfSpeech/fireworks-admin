// Package generator 提供 CRUD 代码生成功能。
// 该文件包含 Ent Schema 解析器的单元测试。
package generator

import (
	"os"
	"path/filepath"
	"testing"
)

// TestNewEntSchemaParser 测试创建解析器实例。
func TestNewEntSchemaParser(t *testing.T) {
	schemaDir := "/path/to/schema"
	parser := NewEntSchemaParser(schemaDir)

	if parser == nil {
		t.Fatal("解析器实例不应为 nil")
	}
	if parser.schemaDir != schemaDir {
		t.Errorf("schemaDir 期望 %s, 实际 %s", schemaDir, parser.schemaDir)
	}
}

// TestMapFieldTypeToGo 测试 ent 字段类型到 Go 类型的映射。
func TestMapFieldTypeToGo(t *testing.T) {
	parser := NewEntSchemaParser("")

	tests := []struct {
		entType  string
		expected string
	}{
		{"String", "string"},
		{"Int", "int"},
		{"Int8", "int8"},
		{"Int16", "int16"},
		{"Int32", "int32"},
		{"Int64", "int64"},
		{"Uint", "uint"},
		{"Uint8", "uint8"},
		{"Uint16", "uint16"},
		{"Uint32", "uint32"},
		{"Uint64", "uint64"},
		{"Float32", "float32"},
		{"Float64", "float64"},
		{"Bool", "bool"},
		{"Time", "time.Time"},
		{"JSON", "json.RawMessage"},
		{"Bytes", "[]byte"},
		{"UUID", "string"},
		{"Text", "string"},
		{"Unknown", "interface{}"},
	}

	for _, tt := range tests {
		t.Run(tt.entType, func(t *testing.T) {
			result := parser.mapFieldTypeToGo(tt.entType)
			if result != tt.expected {
				t.Errorf("mapFieldTypeToGo(%s) = %s, 期望 %s", tt.entType, result, tt.expected)
			}
		})
	}
}

// TestParseTenantSchema 测试解析 Tenant schema。
func TestParseTenantSchema(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("获取当前目录失败: %v", err)
	}

	schemaDir := filepath.Join(cwd, "../../ent/schema")
	parser := NewEntSchemaParser(schemaDir)

	schemaInfo, err := parser.Parse("Tenant")
	if err != nil {
		t.Fatalf("解析 Tenant schema 失败: %v", err)
	}

	if schemaInfo.Name != "Tenant" {
		t.Errorf("schema 名称期望 Tenant, 实际 %s", schemaInfo.Name)
	}

	if schemaInfo.PackageName != "tenant" {
		t.Errorf("包名期望 tenant, 实际 %s", schemaInfo.PackageName)
	}

	if len(schemaInfo.Fields) == 0 {
		t.Error("字段列表不应为空")
	}

	fieldNames := make(map[string]bool)
	for _, f := range schemaInfo.Fields {
		fieldNames[f.Name] = true
	}

	expectedFields := []string{"certificate_no", "name", "type", "contact_name", "email", "phone", "expired_at"}
	for _, name := range expectedFields {
		if !fieldNames[name] {
			t.Errorf("缺少字段: %s", name)
		}
	}

	if len(schemaInfo.Indexes) == 0 {
		t.Error("索引列表不应为空")
	}

	if !schemaInfo.HasSoftDelete {
		t.Error("HasSoftDelete 应为 true")
	}

	if !schemaInfo.HasStatus {
		t.Error("HasStatus 应为 true")
	}
}

// TestParseFile 测试解析指定路径的 schema 文件。
func TestParseFile(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("获取当前目录失败: %v", err)
	}

	filePath := filepath.Join(cwd, "../../ent/schema/tenant.go")
	parser := NewEntSchemaParser("")

	schemaInfo, err := parser.ParseFile(filePath)
	if err != nil {
		t.Fatalf("解析文件失败: %v", err)
	}

	if schemaInfo.Name != "Tenant" {
		t.Errorf("schema 名称期望 Tenant, 实际 %s", schemaInfo.Name)
	}
}

// TestListSchemas 测试列出所有 schema。
func TestListSchemas(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("获取当前目录失败: %v", err)
	}

	schemaDir := filepath.Join(cwd, "../../ent/schema")
	parser := NewEntSchemaParser(schemaDir)

	schemas, err := parser.ListSchemas()
	if err != nil {
		t.Fatalf("列出 schema 失败: %v", err)
	}

	if len(schemas) == 0 {
		t.Error("schema 列表不应为空")
	}

	found := false
	for _, s := range schemas {
		if s == "Tenant" {
			found = true
			break
		}
	}
	if !found {
		t.Error("应包含 Tenant schema")
	}
}

// TestGetKnownMixinFields 测试获取已知 mixin 字段。
func TestGetKnownMixinFields(t *testing.T) {
	parser := NewEntSchemaParser("")

	tests := []struct {
		mixinName       string
		expectedFields  []string
	}{
		{"Id", []string{"id"}},
		{"TenantId", []string{"tenant_id"}},
		{"Status", []string{"status"}},
		{"CreateTime", []string{"created_at"}},
		{"UpdateTime", []string{"updated_at"}},
		{"SoftDelete", []string{"deleted_at"}},
		{"CommonMixin", []string{"id", "tenant_id", "status", "created_at", "updated_at", "deleted_at"}},
		{"Unknown", nil},
	}

	for _, tt := range tests {
		t.Run(tt.mixinName, func(t *testing.T) {
			fields := parser.getKnownMixinFields(tt.mixinName)
			if len(fields) != len(tt.expectedFields) {
				t.Errorf("getKnownMixinFields(%s) 返回 %d 个字段, 期望 %d 个", tt.mixinName, len(fields), len(tt.expectedFields))
				return
			}
			for i, f := range fields {
				if f != tt.expectedFields[i] {
					t.Errorf("getKnownMixinFields(%s)[%d] = %s, 期望 %s", tt.mixinName, i, f, tt.expectedFields[i])
				}
			}
		})
	}
}

// TestGetFieldTypeJSONTag 测试获取 JSON 类型标签。
func TestGetFieldTypeJSONTag(t *testing.T) {
	tests := []struct {
		goType   string
		expected string
	}{
		{"string", "string"},
		{"int", "integer"},
		{"int8", "integer"},
		{"int64", "integer"},
		{"uint", "integer"},
		{"float32", "number"},
		{"float64", "number"},
		{"bool", "boolean"},
		{"time.Time", "string"},
		{"uuid.UUID", "string"},
		{"unknown", "object"},
	}

	for _, tt := range tests {
		t.Run(tt.goType, func(t *testing.T) {
			result := GetFieldTypeJSONTag(tt.goType)
			if result != tt.expected {
				t.Errorf("GetFieldTypeJSONTag(%s) = %s, 期望 %s", tt.goType, result, tt.expected)
			}
		})
	}
}

// TestIsNumericType 测试数值类型检查。
func TestIsNumericType(t *testing.T) {
	numericTypes := []string{"int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "float32", "float64"}
	nonNumericTypes := []string{"string", "bool", "time.Time", "uuid.UUID"}

	for _, tpe := range numericTypes {
		if !IsNumericType(tpe) {
			t.Errorf("IsNumericType(%s) 应返回 true", tpe)
		}
	}

	for _, tpe := range nonNumericTypes {
		if IsNumericType(tpe) {
			t.Errorf("IsNumericType(%s) 应返回 false", tpe)
		}
	}
}

// TestIsStringType 测试字符串类型检查。
func TestIsStringType(t *testing.T) {
	stringTypes := []string{"string", "[]byte"}
	nonStringTypes := []string{"int", "bool", "time.Time"}

	for _, tpe := range stringTypes {
		if !IsStringType(tpe) {
			t.Errorf("IsStringType(%s) 应返回 true", tpe)
		}
	}

	for _, tpe := range nonStringTypes {
		if IsStringType(tpe) {
			t.Errorf("IsStringType(%s) 应返回 false", tpe)
		}
	}
}

// TestIsTimeType 测试时间类型检查。
func TestIsTimeType(t *testing.T) {
	if !IsTimeType("time.Time") {
		t.Error("IsTimeType(time.Time) 应返回 true")
	}
	if IsTimeType("string") {
		t.Error("IsTimeType(string) 应返回 false")
	}
}

// TestIsUUIDType 测试 UUID 类型检查。
func TestIsUUIDType(t *testing.T) {
	if !IsUUIDType("uuid.UUID") {
		t.Error("IsUUIDType(uuid.UUID) 应返回 true")
	}
	if IsUUIDType("string") {
		t.Error("IsUUIDType(string) 应返回 false")
	}
}

// TestFieldConstraints 测试字段约束解析。
func TestFieldConstraints(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("获取当前目录失败: %v", err)
	}

	schemaDir := filepath.Join(cwd, "../../ent/schema")
	parser := NewEntSchemaParser(schemaDir)

	schemaInfo, err := parser.Parse("Tenant")
	if err != nil {
		t.Fatalf("解析 Tenant schema 失败: %v", err)
	}

	fieldMap := make(map[string]FieldInfo)
	for _, f := range schemaInfo.Fields {
		fieldMap[f.Name] = f
	}

	certificateNo, ok := fieldMap["certificate_no"]
	if !ok {
		t.Fatal("缺少 certificate_no 字段")
	}
	if certificateNo.MaxLen != 50 {
		t.Errorf("certificate_no MaxLen 期望 50, 实际 %d", certificateNo.MaxLen)
	}
	if !certificateNo.IsRequired {
		t.Error("certificate_no 应为必填字段")
	}

	expiredAt, ok := fieldMap["expired_at"]
	if !ok {
		t.Fatal("缺少 expired_at 字段")
	}
	if !expiredAt.IsOptional {
		t.Error("expired_at 应为可选字段")
	}

	typeField, ok := fieldMap["type"]
	if !ok {
		t.Fatal("缺少 type 字段")
	}
	if !typeField.HasDefault {
		t.Error("type 应有默认值")
	}
	if typeField.MinVal != 1 {
		t.Errorf("type MinVal 期望 1, 实际 %d", typeField.MinVal)
	}
	if typeField.MaxVal != 2 {
		t.Errorf("type MaxVal 期望 2, 实际 %d", typeField.MaxVal)
	}
}

// TestIndexInfo 测试索引信息解析。
func TestIndexInfo(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("获取当前目录失败: %v", err)
	}

	schemaDir := filepath.Join(cwd, "../../ent/schema")
	parser := NewEntSchemaParser(schemaDir)

	schemaInfo, err := parser.Parse("Tenant")
	if err != nil {
		t.Fatalf("解析 Tenant schema 失败: %v", err)
	}

	indexMap := make(map[string]IndexInfo)
	for _, idx := range schemaInfo.Indexes {
		if len(idx.Fields) > 0 {
			indexMap[idx.Fields[0]] = idx
		}
	}

	emailIdx, ok := indexMap["email"]
	if !ok {
		t.Fatal("缺少 email 索引")
	}
	if !emailIdx.IsUnique {
		t.Error("email 索引应为唯一索引")
	}
	if emailIdx.StorageKey != "uk_email" {
		t.Errorf("email 索引 StorageKey 期望 uk_email, 实际 %s", emailIdx.StorageKey)
	}
	if emailIdx.WhereClause != "deleted_at IS NULL" {
		t.Errorf("email 索引 WhereClause 期望 'deleted_at IS NULL', 实际 '%s'", emailIdx.WhereClause)
	}
}

// TestMixinInfo 测试 mixin 信息解析。
func TestMixinInfo(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("获取当前目录失败: %v", err)
	}

	schemaDir := filepath.Join(cwd, "../../ent/schema")
	parser := NewEntSchemaParser(schemaDir)

	schemaInfo, err := parser.Parse("Tenant")
	if err != nil {
		t.Fatalf("解析 Tenant schema 失败: %v", err)
	}

	expectedMixins := []string{"Id", "Status", "CreateTime", "UpdateTime", "SoftDelete"}
	mixinMap := make(map[string]bool)
	for _, m := range schemaInfo.Mixins {
		mixinMap[m.Name] = true
	}

	for _, name := range expectedMixins {
		if !mixinMap[name] {
			t.Errorf("缺少 mixin: %s", name)
		}
	}

	if !schemaInfo.HasMixin("SoftDelete") {
		t.Error("HasMixin(SoftDelete) 应返回 true")
	}
}
