// Package generator 提供 CRUD 代码生成功能。
// 该包定义了代码生成器的核心结构和接口，用于根据 Ent Schema 自动生成增删改查代码。
package generator

import (
	"bytes"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
	"unicode"
)

// TemplateEngine 模板渲染引擎。
// 负责从文件系统加载模板文件，并提供模板渲染功能。
// 支持模板缓存，避免重复加载相同模板。
type TemplateEngine struct {
	mu         sync.RWMutex                   // mu 读写锁，保护并发访问
	templates  map[string]*template.Template  // templates 已解析的模板缓存
	rawContent map[string]string              // rawContent 原始模板内容缓存
	funcMap    template.FuncMap               // funcMap 模板辅助函数映射
	templateFS fs.FS                          // templateFS 模板文件系统
	rootDir    string                         // rootDir 模板根目录
}

// NewTemplateEngine 创建模板渲染引擎实例。
// 初始化模板缓存和辅助函数映射。
// 返回初始化后的模板引擎实例。
func NewTemplateEngine() *TemplateEngine {
	engine := &TemplateEngine{
		templates:  make(map[string]*template.Template),
		rawContent: make(map[string]string),
		funcMap:    make(template.FuncMap),
	}
	engine.registerFuncs()
	return engine
}

// registerFuncs 注册模板辅助函数。
// 将所有自定义辅助函数注册到模板引擎的函数映射中。
func (e *TemplateEngine) registerFuncs() {
	e.funcMap = template.FuncMap{
		"toLower":      toLower,
		"toUpper":      toUpper,
		"toSnakeCase":  toSnakeCase,
		"toPascalCase": toPascalCase,
		"toCamelCase":  toCamelCase,
		"firstLower":   firstLower,
		"firstUpper":   firstUpper,
		"plural":       plural,
		"singular":     singular,
		"add":          add,
	}
}

// Load 加载指定名称的模板。
// 参数 name 为模板名称，如 "model"、"repository"。
// 返回模板内容和可能的错误。
// 如果模板已缓存，直接返回缓存内容。
func (e *TemplateEngine) Load(name string) (string, error) {
	e.mu.RLock()
	if content, ok := e.rawContent[name]; ok {
		e.mu.RUnlock()
		return content, nil
	}
	e.mu.RUnlock()

	if e.templateFS == nil {
		return "", fmt.Errorf("模板文件系统未初始化，请先调用 LoadFromFS")
	}

	filename := name + ".go.tmpl"
	filepath := filepath.Join(e.rootDir, filename)

	content, err := fs.ReadFile(e.templateFS, filepath)
	if err != nil {
		return "", fmt.Errorf("读取模板文件 %s 失败: %w", filepath, err)
	}

	e.mu.Lock()
	e.rawContent[name] = string(content)
	e.mu.Unlock()

	return string(content), nil
}

// LoadAll 加载所有模板。
// 返回模板名称到内容的映射和可能的错误。
// 遍历模板目录，加载所有 .tmpl 文件。
func (e *TemplateEngine) LoadAll() (map[string]string, error) {
	e.mu.RLock()
	if len(e.rawContent) > 0 {
		result := make(map[string]string, len(e.rawContent))
		for k, v := range e.rawContent {
			result[k] = v
		}
		e.mu.RUnlock()
		return result, nil
	}
	e.mu.RUnlock()

	if e.templateFS == nil {
		return nil, fmt.Errorf("模板文件系统未初始化，请先调用 LoadFromFS")
	}

	result := make(map[string]string)

	err := fs.WalkDir(e.templateFS, e.rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".tmpl") {
			return nil
		}

		content, err := fs.ReadFile(e.templateFS, path)
		if err != nil {
			return fmt.Errorf("读取模板文件 %s 失败: %w", path, err)
		}

		name := e.extractTemplateName(path)
		result[name] = string(content)

		e.mu.Lock()
		e.rawContent[name] = string(content)
		e.mu.Unlock()

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// LoadFromFS 从文件系统加载模板。
// 参数 fsys 为文件系统接口，root 为模板根目录。
// 返回可能的错误。
// 加载所有 .tmpl 文件并解析为模板对象。
func (e *TemplateEngine) LoadFromFS(fsys fs.FS, root string) error {
	e.mu.Lock()
	e.templateFS = fsys
	e.rootDir = root
	e.mu.Unlock()

	return e.loadAndParseTemplates()
}

// loadAndParseTemplates 加载并解析所有模板文件。
// 遍历模板目录，解析每个 .tmpl 文件并缓存。
// 返回可能的错误。
func (e *TemplateEngine) loadAndParseTemplates() error {
	return fs.WalkDir(e.templateFS, e.rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".tmpl") {
			return nil
		}

		content, err := fs.ReadFile(e.templateFS, path)
		if err != nil {
			return fmt.Errorf("读取模板文件 %s 失败: %w", path, err)
		}

		name := e.extractTemplateName(path)

		tmpl, err := template.New(name).Funcs(e.funcMap).Parse(string(content))
		if err != nil {
			return fmt.Errorf("解析模板 %s 失败: %w", name, err)
		}

		e.mu.Lock()
		e.templates[name] = tmpl
		e.rawContent[name] = string(content)
		e.mu.Unlock()

		return nil
	})
}

// extractTemplateName 从文件路径提取模板名称。
// 参数 path 为模板文件的相对路径。
// 返回模板名称（不含扩展名）。
// 例如：templates/model.go.tmpl -> model
func (e *TemplateEngine) extractTemplateName(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	return strings.TrimSuffix(name, filepath.Ext(name))
}

// Render 使用模板和数据渲染生成代码。
// 参数 name 为模板名称，data 为模板数据。
// 返回渲染后的内容和可能的错误。
func (e *TemplateEngine) Render(name string, data interface{}) (string, error) {
	e.mu.RLock()
	tmpl, ok := e.templates[name]
	e.mu.RUnlock()

	if !ok {
		content, err := e.Load(name)
		if err != nil {
			return "", err
		}

		tmpl, err = template.New(name).Funcs(e.funcMap).Parse(content)
		if err != nil {
			return "", fmt.Errorf("解析模板 %s 失败: %w", name, err)
		}

		e.mu.Lock()
		e.templates[name] = tmpl
		e.mu.Unlock()
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("渲染模板 %s 失败: %w", name, err)
	}

	return buf.String(), nil
}

// RenderWithSchema 使用 SchemaInfo 渲染模板。
// 参数 name 为模板名称，schema 为 schema 信息。
// 返回渲染后的内容和可能的错误。
// 这是一个便捷方法，自动构造模板数据结构。
func (e *TemplateEngine) RenderWithSchema(name string, schema *SchemaInfo) (string, error) {
	data := &TemplateData{
		Schema: schema,
	}
	return e.Render(name, data)
}

// GetTemplateNames 获取所有已加载的模板名称。
// 返回模板名称列表。
func (e *TemplateEngine) GetTemplateNames() []string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	names := make([]string, 0, len(e.templates))
	for name := range e.templates {
		names = append(names, name)
	}
	return names
}

// ClearCache 清除模板缓存。
// 用于重新加载模板或释放内存。
func (e *TemplateEngine) ClearCache() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.templates = make(map[string]*template.Template)
	e.rawContent = make(map[string]string)
}

// TemplateData 模板渲染数据结构。
// 包含渲染模板所需的所有数据。
type TemplateData struct {
	Schema *SchemaInfo // Schema schema 信息
}

// toLower 将字符串转换为小写。
// 参数 str 为输入字符串。
// 返回转换后的小写字符串。
// 例如：toLower("Tenant") -> "tenant"
func toLower(str string) string {
	return strings.ToLower(str)
}

// toUpper 将字符串转换为大写。
// 参数 str 为输入字符串。
// 返回转换后的大写字符串。
// 例如：toUpper("tenant") -> "TENANT"
func toUpper(str string) string {
	return strings.ToUpper(str)
}

// toSnakeCase 将驼峰命名转换为下划线命名（蛇形命名）。
// 参数 str 为驼峰命名的字符串。
// 返回转换后的蛇形命名字符串。
// 例如：toSnakeCase("TenantName") -> "tenant_name"
func toSnakeCase(str string) string {
	if str == "" {
		return ""
	}

	var result strings.Builder
	result.Grow(len(str) + 5)

	for i, r := range str {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// toPascalCase 将下划线命名转换为帕斯卡命名（大驼峰）。
// 参数 str 为下划线命名的字符串。
// 返回转换后的帕斯卡命名字符串。
// 例如：toPascalCase("tenant_name") -> "TenantName"
func toPascalCase(str string) string {
	if str == "" {
		return ""
	}

	parts := strings.Split(str, "_")
	var result strings.Builder
	result.Grow(len(str))

	for _, part := range parts {
		if part == "" {
			continue
		}
		runes := []rune(part)
		runes[0] = unicode.ToUpper(runes[0])
		result.WriteString(string(runes))
	}

	return result.String()
}

// toCamelCase 将下划线命名或帕斯卡命名转换为驼峰命名（小驼峰）。
// 参数 str 为输入字符串。
// 返回转换后的驼峰命名字符串。
// 例如：toCamelCase("TenantName") -> "tenantName"
//       toCamelCase("tenant_name") -> "tenantName"
func toCamelCase(str string) string {
	if str == "" {
		return ""
	}

	pascal := toPascalCase(str)
	runes := []rune(pascal)
	if len(runes) > 0 {
		runes[0] = unicode.ToLower(runes[0])
	}
	return string(runes)
}

// firstLower 将字符串首字母转换为小写。
// 参数 str 为输入字符串。
// 返回首字母小写的字符串。
// 例如：firstLower("Tenant") -> "tenant"
func firstLower(str string) string {
	if str == "" {
		return ""
	}

	runes := []rune(str)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// firstUpper 将字符串首字母转换为大写。
// 参数 str 为输入字符串。
// 返回首字母大写的字符串。
// 例如：firstUpper("tenant") -> "Tenant"
func firstUpper(str string) string {
	if str == "" {
		return ""
	}

	runes := []rune(str)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// plural 将单词转换为复数形式。
// 参数 str 为单数形式的单词。
// 返回复数形式的单词。
// 例如：plural("tenant") -> "tenants"
func plural(str string) string {
	if str == "" {
		return ""
	}

	if strings.HasSuffix(str, "y") {
		return strings.TrimSuffix(str, "y") + "ies"
	}

	if strings.HasSuffix(str, "s") ||
		strings.HasSuffix(str, "x") ||
		strings.HasSuffix(str, "z") ||
		strings.HasSuffix(str, "ch") ||
		strings.HasSuffix(str, "sh") {
		return str + "es"
	}

	if strings.HasSuffix(str, "f") {
		return strings.TrimSuffix(str, "f") + "ves"
	}

	if strings.HasSuffix(str, "fe") {
		return strings.TrimSuffix(str, "fe") + "ves"
	}

	return str + "s"
}

// singular 将单词转换为单数形式。
// 参数 str 为复数形式的单词。
// 返回单数形式的单词。
// 例如：singular("tenants") -> "tenant"
func singular(str string) string {
	if str == "" {
		return ""
	}

	if strings.HasSuffix(str, "ies") {
		return strings.TrimSuffix(str, "ies") + "y"
	}

	if strings.HasSuffix(str, "ves") {
		base := strings.TrimSuffix(str, "ves")
		return base + "f"
	}

	if strings.HasSuffix(str, "es") {
		base := strings.TrimSuffix(str, "es")
		if strings.HasSuffix(base, "s") ||
			strings.HasSuffix(base, "x") ||
			strings.HasSuffix(base, "z") ||
			strings.HasSuffix(base, "ch") ||
			strings.HasSuffix(base, "sh") {
			return base
		}
		return str
	}

	if strings.HasSuffix(str, "s") && !strings.HasSuffix(str, "ss") {
		return strings.TrimSuffix(str, "s")
	}

	return str
}

// add 将两个整数相加。
// 参数 a 和 b 为要相加的整数。
// 返回两数之和。
// 例如：add(1, 2) -> 3
func add(a, b int) int {
	return a + b
}
