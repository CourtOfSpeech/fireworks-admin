package generator

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

// TestWireProviderMerge 测试 wire_provider 合并功能。
func TestWireProviderMerge(t *testing.T) {
	existingContent := `// Package tenant 提供租户管理功能。
package tenant

import "github.com/google/wire"

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

// TestWireProviderMergeNoDuplicates 测试不添加重复的 Provider。
func TestWireProviderMergeNoDuplicates(t *testing.T) {
	existingContent := `// Package tenant 提供租户管理功能。
package tenant

import "github.com/google/wire"

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
	NewTenantRepo,
	NewTenantService,
	NewTenantHandler,
)
`

	schema := &SchemaInfo{
		Name:        "Tenant",
		PackageName: "tenant",
	}

	merger := NewWireProviderMerger(schema)
	result, err := merger.Merge(existingContent, newContent)
	if err != nil {
		t.Fatalf("合并失败: %v", err)
	}

	if result != existingContent {
		t.Errorf("期望内容不变，但得到了不同的结果")
	}

	t.Logf("合并结果:\n%s", result)
}

// TestExtractProviders 测试提取 Provider 功能。
func TestExtractProviders(t *testing.T) {
	content := `package tenant

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewTenantRepo,
	NewTenantService,
	NewTenantHandler,
)
`

	schema := &SchemaInfo{
		Name:        "Tenant",
		PackageName: "tenant",
	}

	merger := NewWireProviderMerger(schema)

	fset, err := parseFile(content)
	if err != nil {
		t.Fatalf("解析文件失败: %v", err)
	}

	providers := merger.extractProviders(fset)

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

// parseFile 解析 Go 源文件内容并返回 AST。
func parseFile(content string) (*ast.File, error) {
	fset := token.NewFileSet()
	return parser.ParseFile(fset, "", content, parser.ParseComments)
}
