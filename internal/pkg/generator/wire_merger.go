package generator

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

// WireProviderMerger 负责 wire_provider.go 文件的智能合并。
// 支持在现有文件中添加新的 Provider，而不是简单覆盖。
type WireProviderMerger struct {
	schema *SchemaInfo
}

// NewWireProviderMerger 创建 wire_provider 合并器实例。
// 参数 schema 为 schema 信息。
// 返回初始化后的合并器实例。
func NewWireProviderMerger(schema *SchemaInfo) *WireProviderMerger {
	return &WireProviderMerger{
		schema: schema,
	}
}

// Merge 合并新的 Provider 到现有文件中。
// 参数 existingContent 为现有文件内容，newContent 为新生成的文件内容。
// 返回合并后的文件内容。
func (m *WireProviderMerger) Merge(existingContent, newContent string) (string, error) {
	fset := token.NewFileSet()

	existingFile, err := parser.ParseFile(fset, "", existingContent, parser.ParseComments)
	if err != nil {
		return "", fmt.Errorf("解析现有文件失败: %w", err)
	}

	newFile, err := parser.ParseFile(fset, "", newContent, parser.ParseComments)
	if err != nil {
		return "", fmt.Errorf("解析新文件失败: %w", err)
	}

	existingProviders := m.extractProviders(existingFile)
	newProviders := m.extractProviders(newFile)

	var providersToAdd []string
	for _, provider := range newProviders {
		if !contains(existingProviders, provider) {
			providersToAdd = append(providersToAdd, provider)
		}
	}

	if len(providersToAdd) == 0 {
		return existingContent, nil
	}

	modifiedContent, err := m.addProvidersToContentUsingAST(fset, existingFile, existingContent, providersToAdd)
	if err != nil {
		return "", fmt.Errorf("添加 Provider 失败: %w", err)
	}

	return modifiedContent, nil
}

// extractProviders 从 AST 文件中提取 ProviderSet 中的所有 Provider。
// 参数 file 为 Go 源文件的 AST 节点。
// 返回 Provider 名称列表。
func (m *WireProviderMerger) extractProviders(file *ast.File) []string {
	var providers []string

	ast.Inspect(file, func(n ast.Node) bool {
		genDecl, ok := n.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.VAR {
			return true
		}

		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok || len(valueSpec.Names) == 0 {
				continue
			}

			if valueSpec.Names[0].Name != "ProviderSet" {
				continue
			}

			if len(valueSpec.Values) == 0 {
				continue
			}

			callExpr, ok := valueSpec.Values[0].(*ast.CallExpr)
			if !ok {
				continue
			}

			selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
			if !ok || selectorExpr.Sel.Name != "NewSet" {
				continue
			}

			for _, arg := range callExpr.Args {
				if ident, ok := arg.(*ast.Ident); ok {
					providers = append(providers, ident.Name)
				}
			}
		}

		return true
	})

	return providers
}

// addProvidersToContentUsingAST 使用 AST 精确添加 Provider 到文件中。
// 参数 fset 为文件集，file 为 AST 文件节点，content 为原始内容，providers 为要添加的 Provider 列表。
// 返回修改后的文件内容。
func (m *WireProviderMerger) addProvidersToContentUsingAST(fset *token.FileSet, file *ast.File, content string, providers []string) (string, error) {
	var callExpr *ast.CallExpr

	ast.Inspect(file, func(n ast.Node) bool {
		if gd, ok := n.(*ast.GenDecl); ok && gd.Tok == token.VAR {
			for _, spec := range gd.Specs {
				if vs, ok := spec.(*ast.ValueSpec); ok && len(vs.Names) > 0 {
					if vs.Names[0].Name == "ProviderSet" && len(vs.Values) > 0 {
						if ce, ok := vs.Values[0].(*ast.CallExpr); ok {
							callExpr = ce
						}
					}
				}
			}
		}
		return true
	})

	if callExpr == nil {
		return "", fmt.Errorf("未找到 ProviderSet 定义")
	}

	lines := strings.Split(content, "\n")

	startPos := fset.Position(callExpr.Lparen)
	endPos := fset.Position(callExpr.Rparen)

	startLine := startPos.Line - 1
	endLine := endPos.Line - 1

	if startLine < 0 || endLine < 0 || startLine >= len(lines) || endLine >= len(lines) {
		return "", fmt.Errorf("无效的行号范围")
	}

	indent := m.detectIndentation(lines, startLine, endLine)

	var result []string
	result = append(result, lines[:endLine]...)

	for _, provider := range providers {
		result = append(result, indent+provider+",")
	}

	result = append(result, lines[endLine:]...)

	merged := strings.Join(result, "\n")

	formatted, err := formatGoCode(merged)
	if err != nil {
		return "", fmt.Errorf("格式化代码失败: %w", err)
	}

	return formatted, nil
}

// detectIndentation 检测现有代码的缩进风格。
// 参数 lines 为文件行列表，startLine 为 ProviderSet 起始行，endLine 为结束行。
// 返回检测到的缩进字符串。
func (m *WireProviderMerger) detectIndentation(lines []string, startLine, endLine int) string {
	for i := startLine; i <= endLine && i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && trimmed != "wire.NewSet(" && trimmed != ")" {
			indent := ""
			for _, ch := range line {
				if ch == ' ' || ch == '\t' {
					indent += string(ch)
				} else {
					break
				}
			}
			if indent != "" {
				return indent
			}
		}
	}
	return "\t"
}

// NeedsMerge 检查是否需要合并。
// 参数 filePath 为文件路径。
// 返回 true 表示文件存在且需要合并。
func (m *WireProviderMerger) NeedsMerge(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

// contains 检查字符串切片中是否包含指定字符串。
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
