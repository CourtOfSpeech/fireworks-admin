package generator

import (
	"testing"
)

// TestParseFileEmptyPath 测试空文件路径的错误处理。
func TestParseFileEmptyPath(t *testing.T) {
	parser := NewEntSchemaParser("internal/ent/schema")
	_, err := parser.ParseFile("")
	if err == nil {
		t.Error("期望返回错误，但没有返回")
	}
	if err.Error() != "文件路径不能为空" {
		t.Errorf("期望错误信息 '文件路径不能为空'，但得到 '%s'", err.Error())
	}
}

// TestParseFileNotExists 测试文件不存在的错误处理。
func TestParseFileNotExists(t *testing.T) {
	parser := NewEntSchemaParser("internal/ent/schema")
	_, err := parser.ParseFile("internal/ent/schema/not_exists.go")
	if err == nil {
		t.Error("期望返回错误，但没有返回")
	}
	expectedMsg := "schema 文件不存在: internal/ent/schema/not_exists.go"
	if err.Error() != expectedMsg {
		t.Errorf("期望错误信息 '%s'，但得到 '%s'", expectedMsg, err.Error())
	}
}

// TestGenerateWithEmptyGenerators 测试生成器列表为空的错误处理。
func TestGenerateWithEmptyGenerators(t *testing.T) {
	g := NewGenerator(Config{
		SchemaName: "Tenant",
		OutputPath: "internal/features/tenant",
	})

	g.SetParser(NewEntSchemaParser("internal/ent/schema"))
	g.SetTemplateEngine(NewTemplateEngine())
	g.SetGenerators([]FileGenerator{})

	err := g.Generate()
	if err == nil {
		t.Error("期望返回错误，但没有返回")
	}
	if err.Error() != "文件生成器列表为空" {
		t.Errorf("期望错误信息 '文件生成器列表为空'，但得到 '%s'", err.Error())
	}
}

// TestGenerateWithResultWithEmptyGenerators 测试 GenerateWithResult 方法生成器列表为空的错误处理。
func TestGenerateWithResultWithEmptyGenerators(t *testing.T) {
	g := NewGenerator(Config{
		SchemaName: "Tenant",
		OutputPath: "internal/features/tenant",
	})

	g.SetParser(NewEntSchemaParser("internal/ent/schema"))
	g.SetTemplateEngine(NewTemplateEngine())
	g.SetGenerators([]FileGenerator{})

	_, err := g.GenerateWithResult()
	if err == nil {
		t.Error("期望返回错误，但没有返回")
	}
	if err.Error() != "文件生成器列表为空" {
		t.Errorf("期望错误信息 '文件生成器列表为空'，但得到 '%s'", err.Error())
	}
}
