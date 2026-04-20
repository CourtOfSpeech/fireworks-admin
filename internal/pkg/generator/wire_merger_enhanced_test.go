package generator

import (
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

// TestWireProviderMergeWithComments 测试带注释的情况。
func TestWireProviderMergeWithComments(t *testing.T) {
	existingContent := `// Package tenant 提供租户管理功能。
package tenant

import "github.com/google/wire"

// ProviderSet 租户模块依赖提供者集合。
// 包含租户模块所有需要注入的组件。
var ProviderSet = wire.NewSet(
	NewTenantRepo,
	NewTenantService,
	NewTenantHandler,
)
`

	newContent := `// Package tenant 提供租户管理功能。
package tenant

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewUserRepo,
	NewUserService,
	NewUserHandler,
)
`

	schema := &SchemaInfo{
		Name:        "User",
		PackageName: "tenant",
	}

	merger := NewWireProviderMerger(schema)
	result, err := merger.Merge(existingContent, newContent)
	if err != nil {
		t.Fatalf("合并失败: %v", err)
	}

	expectedProviders := []string{
		"NewTenantRepo",
		"NewTenantService",
		"NewTenantHandler",
		"NewUserRepo",
		"NewUserService",
		"NewUserHandler",
	}

	for _, provider := range expectedProviders {
		if !strings.Contains(result, provider) {
			t.Errorf("期望包含 Provider: %s", provider)
		}
	}

	t.Logf("合并结果:\n%s", result)
}

// TestWireProviderMergeWithDifferentIndentation 测试不同缩进风格。
func TestWireProviderMergeWithDifferentIndentation(t *testing.T) {
	existingContent := `package tenant

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
    NewTenantRepo,
    NewTenantService,
    NewTenantHandler,
)
`

	newContent := `package tenant

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
    NewUserRepo,
    NewUserService,
    NewUserHandler,
)
`

	schema := &SchemaInfo{
		Name:        "User",
		PackageName: "tenant",
	}

	merger := NewWireProviderMerger(schema)
	result, err := merger.Merge(existingContent, newContent)
	if err != nil {
		t.Fatalf("合并失败: %v", err)
	}

	expectedProviders := []string{
		"NewTenantRepo",
		"NewTenantService",
		"NewTenantHandler",
		"NewUserRepo",
		"NewUserService",
		"NewUserHandler",
	}

	for _, provider := range expectedProviders {
		if !strings.Contains(result, provider) {
			t.Errorf("期望包含 Provider: %s", provider)
		}
	}

	t.Logf("合并结果:\n%s", result)
}

// TestWireProviderMergeWithErrorHandling 测试错误处理。
func TestWireProviderMergeWithErrorHandling(t *testing.T) {
	invalidContent := `package tenant

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewTenantRepo,
	// 缺少闭合括号
`

	schema := &SchemaInfo{
		Name:        "User",
		PackageName: "tenant",
	}

	merger := NewWireProviderMerger(schema)
	_, err := merger.Merge(invalidContent, invalidContent)
	if err == nil {
		t.Error("期望返回错误，但没有返回")
	}

	t.Logf("错误信息: %v", err)
}

// TestDetectIndentation 测试缩进检测。
func TestDetectIndentation(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		wantIndent string
	}{
		{
			name: "制表符缩进",
			content: `var ProviderSet = wire.NewSet(
	NewTenantRepo,
	NewTenantService,
)`,
			wantIndent: "\t",
		},
		{
			name: "空格缩进",
			content: `var ProviderSet = wire.NewSet(
    NewTenantRepo,
    NewTenantService,
)`,
			wantIndent: "    ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines := strings.Split(tt.content, "\n")
			merger := NewWireProviderMerger(&SchemaInfo{})

			indent := merger.detectIndentation(lines, 0, len(lines)-1)
			if indent != tt.wantIndent {
				t.Errorf("期望缩进 %q，但得到 %q", tt.wantIndent, indent)
			}
		})
	}
}

// TestExtractProvidersWithComments 测试从带注释的代码中提取 Provider。
func TestExtractProvidersWithComments(t *testing.T) {
	content := `package tenant

import "github.com/google/wire"

// ProviderSet 依赖提供者集合
var ProviderSet = wire.NewSet(
	NewTenantRepo,    // Repository
	NewTenantService, // Service
	NewTenantHandler, // Handler
)
`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", content, parser.ParseComments)
	if err != nil {
		t.Fatalf("解析文件失败: %v", err)
	}

	schema := &SchemaInfo{
		Name:        "Tenant",
		PackageName: "tenant",
	}

	merger := NewWireProviderMerger(schema)
	providers := merger.extractProviders(file)

	expectedProviders := []string{
		"NewTenantRepo",
		"NewTenantService",
		"NewTenantHandler",
	}

	if len(providers) != len(expectedProviders) {
		t.Errorf("期望 %d 个 Provider，但得到 %d 个", len(expectedProviders), len(providers))
	}

	for i, expected := range expectedProviders {
		if providers[i] != expected {
			t.Errorf("期望 Provider[%d] = %s，但得到 %s", i, expected, providers[i])
		}
	}
}
