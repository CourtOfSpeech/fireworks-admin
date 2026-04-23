package generator

import (
	"testing"
)

// TestUUIDFieldMapping 测试 UUID 字段映射为 string 类型。
func TestUUIDFieldMapping(t *testing.T) {
	parser := NewEntSchemaParser("internal/ent/schema")

	result := parser.mapFieldTypeToGo("UUID")
	expected := "string"

	if result != expected {
		t.Errorf("UUID 类型映射错误: 期望 %s, 实际 %s", expected, result)
	}

	t.Logf("✅ UUID 类型正确映射为 string 类型")
}
