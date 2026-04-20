// Package generator 提供 CRUD 代码生成功能。
// 该包定义了代码生成器的核心结构和接口，用于根据 Ent Schema 自动生成增删改查代码。
package generator

import (
	"fmt"
	"go/format"
	"io/fs"
	"os"
	"path/filepath"
)

// ErrFileExists 文件已存在错误。
// 当尝试写入已存在的文件且未启用强制覆盖时返回此错误。
var ErrFileExists = fmt.Errorf("文件已存在")

// Config 代码生成器配置。
// 包含生成代码所需的所有配置参数。
type Config struct {
	SchemaName     string // SchemaName schema 名称，如 "User"
	OutputPath     string // OutputPath 输出路径，如 "internal/features/user"
	Force          bool   // Force 是否强制覆盖已存在的文件
	FileNamePrefix bool   // FileNamePrefix 是否在文件名前添加 schema 名称前缀
}

// Generator 代码生成器核心结构。
// 负责解析 Ent Schema 并生成对应的 CRUD 代码文件。
type Generator struct {
	config     Config          // config 生成器配置
	parser     SchemaParser    // parser schema 解析器
	templates  TemplateLoader  // templates 模板加载器
	engine     *TemplateEngine // engine 模板渲染引擎
	writer     FileWriter      // writer 文件写入器
	generators []FileGenerator // generators 文件生成器列表
}

// NewGenerator 创建代码生成器实例。
// 参数 config 为生成器配置。
// 返回初始化后的生成器实例。
func NewGenerator(config Config) *Generator {
	return &Generator{
		config:     config,
		engine:     NewTemplateEngine(),
		generators: DefaultGenerators(),
	}
}

// SetParser 设置 schema 解析器。
// 参数 parser 为实现了 SchemaParser 接口的解析器实例。
func (g *Generator) SetParser(parser SchemaParser) {
	g.parser = parser
}

// SetTemplateLoader 设置模板加载器。
// 参数 loader 为实现了 TemplateLoader 接口的加载器实例。
func (g *Generator) SetTemplateLoader(loader TemplateLoader) {
	g.templates = loader
}

// SetFileWriter 设置文件写入器。
// 参数 writer 为实现了 FileWriter 接口的写入器实例。
func (g *Generator) SetFileWriter(writer FileWriter) {
	g.writer = writer
}

// SetTemplateEngine 设置模板渲染引擎。
// 参数 engine 为模板渲染引擎实例。
func (g *Generator) SetTemplateEngine(engine *TemplateEngine) {
	g.engine = engine
}

// SetGenerators 设置文件生成器列表。
// 参数 generators 为文件生成器列表。
func (g *Generator) SetGenerators(generators []FileGenerator) {
	g.generators = generators
}

// LoadTemplatesFromFS 从文件系统加载模板。
// 参数 fsys 为文件系统接口，root 为模板根目录。
// 返回可能的错误。
func (g *Generator) LoadTemplatesFromFS(fsys fs.FS, root string) error {
	return g.engine.LoadFromFS(fsys, root)
}

// Generate 执行代码生成。
// 该方法依次执行：解析 schema -> 加载模板 -> 渲染模板 -> 写入文件。
// 返回生成过程中可能发生的错误。
func (g *Generator) Generate() error {
	if g.parser == nil {
		return fmt.Errorf("schema 解析器未设置")
	}
	if g.engine == nil {
		return fmt.Errorf("模板引擎未设置")
	}
	if len(g.generators) == 0 {
		return fmt.Errorf("文件生成器列表为空")
	}
	if g.writer == nil {
		g.writer = NewOSFileWriter(g.config.Force)
	}

	schema, err := g.parser.Parse(g.config.SchemaName)
	if err != nil {
		return fmt.Errorf("解析 schema 失败: %w", err)
	}

	schema.Comment = g.config.SchemaName

	if err := g.writer.MkdirAll(g.config.OutputPath); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	for _, gen := range g.generators {
		filename, content, err := gen.Generate(schema)
		if err != nil {
			return fmt.Errorf("生成 %s 失败: %w", gen.Name(), err)
		}

		if content == "" {
			continue
		}

		outputPath := filepath.Join(g.config.OutputPath, filename)

		formatted, err := formatGoCode(content)
		if err != nil {
			fmt.Printf("警告: 格式化 %s 失败: %v，将使用原始内容\n", filename, err)
			formatted = content
		}

		if g.config.Force {
			if err := g.writer.Write(outputPath, formatted); err != nil {
				return fmt.Errorf("写入文件 %s 失败: %w", outputPath, err)
			}
		} else {
			if err := g.writer.WriteIfNotExists(outputPath, formatted); err != nil {
				if err == ErrFileExists {
					fmt.Printf("跳过已存在的文件: %s\n", outputPath)
					continue
				}
				return fmt.Errorf("写入文件 %s 失败: %w", outputPath, err)
			}
		}

		fmt.Printf("生成文件: %s\n", outputPath)
	}

	return nil
}

// GenerateWithResult 执行代码生成并返回详细结果。
// 该方法依次执行：解析 schema -> 加载模板 -> 渲染模板 -> 写入文件。
// 返回生成结果，包含生成的文件列表和可能的错误。
func (g *Generator) GenerateWithResult() (*GenerateResult, error) {
	result := &GenerateResult{
		Files:    make([]GeneratedFile, 0),
		Warnings: make([]string, 0),
		Errors:   make([]error, 0),
	}

	if g.parser == nil {
		return nil, fmt.Errorf("schema 解析器未设置")
	}
	if g.engine == nil {
		return nil, fmt.Errorf("模板引擎未设置")
	}
	if len(g.generators) == 0 {
		return nil, fmt.Errorf("文件生成器列表为空")
	}
	if g.writer == nil {
		g.writer = NewOSFileWriter(g.config.Force)
	}

	schema, err := g.parser.Parse(g.config.SchemaName)
	if err != nil {
		return nil, fmt.Errorf("解析 schema 失败: %w", err)
	}

	schema.Comment = g.config.SchemaName

	if err := g.writer.MkdirAll(g.config.OutputPath); err != nil {
		return nil, fmt.Errorf("创建输出目录失败: %w", err)
	}

	for _, gen := range g.generators {
		if wireGen, ok := gen.(*WireProviderGenerator); ok {
			outputPath := filepath.Join(g.config.OutputPath, "wire_provider.go")

			filename, content, skipped, err := wireGen.GenerateWithMerge(schema, outputPath)
			if err != nil {
				result.Errors = append(result.Errors, fmt.Errorf("生成 %s 失败: %w", gen.Name(), err))
				continue
			}

			if skipped {
				generatedFile := GeneratedFile{
					Path:    outputPath,
					Name:    filename,
					Skipped: true,
				}
				result.Files = append(result.Files, generatedFile)
				continue
			}

			if content == "" {
				continue
			}

			formatted, err := formatGoCode(content)
			if err != nil {
				result.Warnings = append(result.Warnings, fmt.Sprintf("格式化 %s 失败: %v", filename, err))
				formatted = content
			}

			generatedFile := GeneratedFile{
				Path: outputPath,
				Name: filename,
			}

			if g.config.Force {
				if err := g.writer.Write(outputPath, formatted); err != nil {
					result.Errors = append(result.Errors, fmt.Errorf("写入文件 %s 失败: %w", outputPath, err))
					continue
				}
			} else {
				if err := g.writer.Write(outputPath, formatted); err != nil {
					result.Errors = append(result.Errors, fmt.Errorf("写入文件 %s 失败: %w", outputPath, err))
					continue
				}
			}

			generatedFile.Size = int64(len(formatted))
			result.Files = append(result.Files, generatedFile)
			continue
		}

		filename, content, err := gen.Generate(schema)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("生成 %s 失败: %w", gen.Name(), err))
			continue
		}

		if content == "" {
			continue
		}

		outputPath := filepath.Join(g.config.OutputPath, filename)

		formatted, err := formatGoCode(content)
		if err != nil {
			result.Warnings = append(result.Warnings, fmt.Sprintf("格式化 %s 失败: %v", filename, err))
			formatted = content
		}

		generatedFile := GeneratedFile{
			Path: outputPath,
			Name: filename,
		}

		if g.config.Force {
			if err := g.writer.Write(outputPath, formatted); err != nil {
				result.Errors = append(result.Errors, fmt.Errorf("写入文件 %s 失败: %w", outputPath, err))
				continue
			}
		} else {
			if err := g.writer.WriteIfNotExists(outputPath, formatted); err != nil {
				if err == ErrFileExists {
					generatedFile.Skipped = true
					result.Files = append(result.Files, generatedFile)
					continue
				}
				result.Errors = append(result.Errors, fmt.Errorf("写入文件 %s 失败: %w", outputPath, err))
				continue
			}
		}

		generatedFile.Size = int64(len(formatted))
		result.Files = append(result.Files, generatedFile)
	}

	return result, nil
}

// formatGoCode 使用 gofmt 格式化 Go 代码。
// 参数 content 为原始 Go 代码字符串。
// 返回格式化后的代码和可能的错误。
func formatGoCode(content string) (string, error) {
	formatted, err := format.Source([]byte(content))
	if err != nil {
		return "", err
	}
	return string(formatted), nil
}

// SchemaParser 定义 Schema 解析器接口。
// 负责从 Ent Schema 文件中解析出结构化的 Schema 信息。
type SchemaParser interface {
	// Parse 解析指定名称的 Ent Schema。
	// 参数 schemaName 为 schema 名称，如 "Tenant"。
	// 返回解析后的 Schema 信息和可能的错误。
	Parse(schemaName string) (*SchemaInfo, error)

	// ParseFile 解析指定路径的 Ent Schema 文件。
	// 参数 filePath 为 schema 文件的绝对路径。
	// 返回解析后的 Schema 信息和可能的错误。
	ParseFile(filePath string) (*SchemaInfo, error)

	// ListSchemas 列出所有可用的 schema 名称。
	// 返回 schema 名称列表和可能的错误。
	ListSchemas() ([]string, error)
}

// TemplateLoader 定义模板加载器接口。
// 负责加载和管理代码生成模板。
type TemplateLoader interface {
	// Load 加载指定名称的模板。
	// 参数 name 为模板名称，如 "model"、"repository"。
	// 返回模板内容和可能的错误。
	Load(name string) (string, error)

	// LoadAll 加载所有模板。
	// 返回模板名称到内容的映射和可能的错误。
	LoadAll() (map[string]string, error)

	// LoadFromFS 从文件系统加载模板。
	// 参数 fsys 为文件系统接口，root 为模板根目录。
	// 返回可能的错误。
	LoadFromFS(fsys fs.FS, root string) error
}

// FileWriter 定义文件写入器接口。
// 负责将生成的代码写入文件系统。
type FileWriter interface {
	// Write 将内容写入指定路径的文件。
	// 参数 path 为文件绝对路径，content 为文件内容。
	// 返回可能的错误。
	Write(path, content string) error

	// WriteIfNotExists 仅在文件不存在时写入。
	// 参数 path 为文件绝对路径，content 为文件内容。
	// 返回可能的错误，如果文件已存在返回 ErrFileExists。
	WriteIfNotExists(path, content string) error

	// Exists 检查文件是否存在。
	// 参数 path 为文件绝对路径。
	// 返回 true 表示文件存在。
	Exists(path string) bool

	// MkdirAll 创建目录及其所有父目录。
	// 参数 path 为目录绝对路径。
	// 返回可能的错误。
	MkdirAll(path string) error
}

// OSFileWriter 基于操作系统的文件写入器实现。
// 支持文件存在检查、目录创建和强制覆盖功能。
type OSFileWriter struct {
	force bool // force 是否强制覆盖已存在的文件
}

// NewOSFileWriter 创建操作系统文件写入器实例。
// 参数 force 表示是否强制覆盖已存在的文件。
// 返回初始化后的文件写入器实例。
func NewOSFileWriter(force bool) *OSFileWriter {
	return &OSFileWriter{
		force: force,
	}
}

// Write 将内容写入指定路径的文件。
// 如果文件已存在，将覆盖原文件。
// 参数 path 为文件绝对路径，content 为文件内容。
// 返回可能的错误。
func (w *OSFileWriter) Write(path, content string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
}

// WriteIfNotExists 仅在文件不存在时写入。
// 参数 path 为文件绝对路径，content 为文件内容。
// 返回可能的错误，如果文件已存在返回 ErrFileExists。
func (w *OSFileWriter) WriteIfNotExists(path, content string) error {
	if w.Exists(path) {
		return ErrFileExists
	}
	return w.Write(path, content)
}

// Exists 检查文件是否存在。
// 参数 path 为文件绝对路径。
// 返回 true 表示文件存在。
func (w *OSFileWriter) Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// MkdirAll 创建目录及其所有父目录。
// 参数 path 为目录绝对路径。
// 返回可能的错误。
func (w *OSFileWriter) MkdirAll(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}
	return nil
}

// FileGenerator 定义单个文件生成器接口。
// 负责生成特定类型的代码文件。
type FileGenerator interface {
	// Generate 生成代码文件。
	// 参数 schema 为 schema 信息。
	// 返回生成的文件名、内容和可能的错误。
	Generate(schema *SchemaInfo) (filename, content string, err error)

	// Name 返回生成器名称。
	// 如 "model"、"repository"、"service" 等。
	Name() string
}

// baseGenerator 基础生成器，包含模板引擎引用。
type baseGenerator struct {
	engine         *TemplateEngine
	fileNamePrefix bool
	force          bool
}

// getFilename 根据配置生成文件名。
// 参数 baseName 为基础文件名，如 "model"、"repository"。
// 参数 schema 为 schema 信息。
// 返回完整的文件名。
func (g *baseGenerator) getFilename(baseName string, schema *SchemaInfo) string {
	if g.fileNamePrefix {
		return schema.PackageName + "_" + baseName + ".go"
	}
	return baseName + ".go"
}

// ModelGenerator model.go 文件生成器。
// 负责生成领域模型和常量定义。
type ModelGenerator struct {
	baseGenerator
}

// NewModelGenerator 创建 model 生成器实例。
// 参数 engine 为模板渲染引擎，fileNamePrefix 表示是否在文件名前添加 schema 名称前缀。
// 返回初始化后的生成器实例。
func NewModelGenerator(engine *TemplateEngine, fileNamePrefix bool) *ModelGenerator {
	return &ModelGenerator{
		baseGenerator: baseGenerator{engine: engine, fileNamePrefix: fileNamePrefix},
	}
}

// Name 返回生成器名称。
func (g *ModelGenerator) Name() string {
	return "model"
}

// Generate 生成 model.go 文件内容。
// 参数 schema 为 schema 信息。
// 返回生成的文件名、内容和可能的错误。
func (g *ModelGenerator) Generate(schema *SchemaInfo) (string, string, error) {
	content, err := g.engine.RenderWithSchema("model", schema)
	if err != nil {
		return "", "", fmt.Errorf("渲染 model 模板失败: %w", err)
	}
	return g.getFilename("model", schema), content, nil
}

// RepositoryGenerator repository.go 文件生成器。
// 负责生成数据持久化层代码。
type RepositoryGenerator struct {
	baseGenerator
}

// NewRepositoryGenerator 创建 repository 生成器实例。
// 参数 engine 为模板渲染引擎，fileNamePrefix 表示是否在文件名前添加 schema 名称前缀。
// 返回初始化后的生成器实例。
func NewRepositoryGenerator(engine *TemplateEngine, fileNamePrefix bool) *RepositoryGenerator {
	return &RepositoryGenerator{
		baseGenerator: baseGenerator{engine: engine, fileNamePrefix: fileNamePrefix},
	}
}

// Name 返回生成器名称。
func (g *RepositoryGenerator) Name() string {
	return "repository"
}

// Generate 生成 repository.go 文件内容。
// 参数 schema 为 schema 信息。
// 返回生成的文件名、内容和可能的错误。
func (g *RepositoryGenerator) Generate(schema *SchemaInfo) (string, string, error) {
	content, err := g.engine.RenderWithSchema("repository", schema)
	if err != nil {
		return "", "", fmt.Errorf("渲染 repository 模板失败: %w", err)
	}
	return g.getFilename("repository", schema), content, nil
}

// ServiceGenerator service.go 文件生成器。
// 负责生成业务逻辑层代码。
type ServiceGenerator struct {
	baseGenerator
}

// NewServiceGenerator 创建 service 生成器实例。
// 参数 engine 为模板渲染引擎，fileNamePrefix 表示是否在文件名前添加 schema 名称前缀。
// 返回初始化后的生成器实例。
func NewServiceGenerator(engine *TemplateEngine, fileNamePrefix bool) *ServiceGenerator {
	return &ServiceGenerator{
		baseGenerator: baseGenerator{engine: engine, fileNamePrefix: fileNamePrefix},
	}
}

// Name 返回生成器名称。
func (g *ServiceGenerator) Name() string {
	return "service"
}

// Generate 生成 service.go 文件内容。
// 参数 schema 为 schema 信息。
// 返回生成的文件名、内容和可能的错误。
func (g *ServiceGenerator) Generate(schema *SchemaInfo) (string, string, error) {
	content, err := g.engine.RenderWithSchema("service", schema)
	if err != nil {
		return "", "", fmt.Errorf("渲染 service 模板失败: %w", err)
	}
	return g.getFilename("service", schema), content, nil
}

// HandlerGenerator handler.go 文件生成器。
// 负责生成 HTTP 处理器代码。
type HandlerGenerator struct {
	baseGenerator
}

// NewHandlerGenerator 创建 handler 生成器实例。
// 参数 engine 为模板渲染引擎，fileNamePrefix 表示是否在文件名前添加 schema 名称前缀。
// 返回初始化后的生成器实例。
func NewHandlerGenerator(engine *TemplateEngine, fileNamePrefix bool) *HandlerGenerator {
	return &HandlerGenerator{
		baseGenerator: baseGenerator{engine: engine, fileNamePrefix: fileNamePrefix},
	}
}

// Name 返回生成器名称。
func (g *HandlerGenerator) Name() string {
	return "handler"
}

// Generate 生成 handler.go 文件内容。
// 参数 schema 为 schema 信息。
// 返回生成的文件名、内容和可能的错误。
func (g *HandlerGenerator) Generate(schema *SchemaInfo) (string, string, error) {
	content, err := g.engine.RenderWithSchema("handler", schema)
	if err != nil {
		return "", "", fmt.Errorf("渲染 handler 模板失败: %w", err)
	}
	return g.getFilename("handler", schema), content, nil
}

// DtoGenerator dto.go 文件生成器。
// 负责生成数据传输对象代码。
type DtoGenerator struct {
	baseGenerator
}

// NewDtoGenerator 创建 dto 生成器实例。
// 参数 engine 为模板渲染引擎，fileNamePrefix 表示是否在文件名前添加 schema 名称前缀。
// 返回初始化后的生成器实例。
func NewDtoGenerator(engine *TemplateEngine, fileNamePrefix bool) *DtoGenerator {
	return &DtoGenerator{
		baseGenerator: baseGenerator{engine: engine, fileNamePrefix: fileNamePrefix},
	}
}

// Name 返回生成器名称。
func (g *DtoGenerator) Name() string {
	return "dto"
}

// Generate 生成 dto.go 文件内容。
// 参数 schema 为 schema 信息。
// 返回生成的文件名、内容和可能的错误。
func (g *DtoGenerator) Generate(schema *SchemaInfo) (string, string, error) {
	content, err := g.engine.RenderWithSchema("dto", schema)
	if err != nil {
		return "", "", fmt.Errorf("渲染 dto 模板失败: %w", err)
	}
	return g.getFilename("dto", schema), content, nil
}

// ErrorsGenerator errors.go 文件生成器。
// 负责生成错误码和错误处理函数代码。
type ErrorsGenerator struct {
	baseGenerator
}

// NewErrorsGenerator 创建 errors 生成器实例。
// 参数 engine 为模板渲染引擎，fileNamePrefix 表示是否在文件名前添加 schema 名称前缀。
// 返回初始化后的生成器实例。
func NewErrorsGenerator(engine *TemplateEngine, fileNamePrefix bool) *ErrorsGenerator {
	return &ErrorsGenerator{
		baseGenerator: baseGenerator{engine: engine, fileNamePrefix: fileNamePrefix},
	}
}

// Name 返回生成器名称。
func (g *ErrorsGenerator) Name() string {
	return "errors"
}

// Generate 生成 errors.go 文件内容。
// 参数 schema 为 schema 信息。
// 返回生成的文件名、内容和可能的错误。
func (g *ErrorsGenerator) Generate(schema *SchemaInfo) (string, string, error) {
	content, err := g.engine.RenderWithSchema("errors", schema)
	if err != nil {
		return "", "", fmt.Errorf("渲染 errors 模板失败: %w", err)
	}
	return g.getFilename("errors", schema), content, nil
}

// WireProviderGenerator wire_provider.go 文件生成器。
// 负责生成 Wire 依赖注入提供者代码。
type WireProviderGenerator struct {
	baseGenerator
	force bool
}

// NewWireProviderGenerator 创建 wire_provider 生成器实例。
// 参数 engine 为模板渲染引擎，fileNamePrefix 表示是否在文件名前添加 schema 名称前缀，force 表示是否强制覆盖。
// 返回初始化后的生成器实例。
func NewWireProviderGenerator(engine *TemplateEngine, fileNamePrefix bool, force bool) *WireProviderGenerator {
	return &WireProviderGenerator{
		baseGenerator: baseGenerator{engine: engine, fileNamePrefix: fileNamePrefix},
		force:         force,
	}
}

// Name 返回生成器名称。
func (g *WireProviderGenerator) Name() string {
	return "wire_provider"
}

// Generate 生成 wire_provider.go 文件内容。
// 参数 schema 为 schema 信息。
// 返回生成的文件名、内容和可能的错误。
// 在非强制覆盖模式下，如果文件已存在，会智能合并新的 Provider。
// 注意：wire_provider.go 文件名不受 -prefix 参数影响，始终为 wire_provider.go。
func (g *WireProviderGenerator) Generate(schema *SchemaInfo) (string, string, error) {
	newContent, err := g.engine.RenderWithSchema("wire_provider", schema)
	if err != nil {
		return "", "", fmt.Errorf("渲染 wire_provider 模板失败: %w", err)
	}

	return "wire_provider.go", newContent, nil
}

// GenerateWithMerge 生成或合并 wire_provider.go 文件内容。
// 参数 schema 为 schema 信息，existingPath 为现有文件路径（如果存在）。
// 返回生成的文件名、内容、是否跳过和可能的错误。
// 注意：wire_provider.go 文件名不受 -prefix 参数影响，始终为 wire_provider.go。
func (g *WireProviderGenerator) GenerateWithMerge(schema *SchemaInfo, existingPath string) (string, string, bool, error) {
	newContent, err := g.engine.RenderWithSchema("wire_provider", schema)
	if err != nil {
		return "", "", false, fmt.Errorf("渲染 wire_provider 模板失败: %w", err)
	}

	filename := "wire_provider.go"

	if g.force {
		return filename, newContent, false, nil
	}

	merger := NewWireProviderMerger(schema)

	if !merger.NeedsMerge(existingPath) {
		return filename, newContent, false, nil
	}

	existingContent, err := os.ReadFile(existingPath)
	if err != nil {
		return "", "", false, fmt.Errorf("读取现有文件失败: %w", err)
	}

	mergedContent, err := merger.Merge(string(existingContent), newContent)
	if err != nil {
		return "", "", false, fmt.Errorf("合并文件失败: %w", err)
	}

	if mergedContent == string(existingContent) {
		return filename, "", true, nil
	}

	return filename, mergedContent, false, nil
}

// GenerateResult 代码生成结果。
// 包含生成的文件信息和可能的错误。
type GenerateResult struct {
	Files    []GeneratedFile // Files 生成的文件列表
	Warnings []string        // Warnings 警告信息列表
	Errors   []error         // Errors 错误列表
}

// GeneratedFile 生成的文件信息。
type GeneratedFile struct {
	Path    string // Path 文件绝对路径
	Name    string // Name 文件名
	Size    int64  // Size 文件大小（字节）
	Skipped bool   // Skipped 是否跳过（文件已存在且未强制覆盖）
}

// DefaultGenerators 返回默认的文件生成器列表。
// 包含 model、repository、service、handler、dto、errors、wire_provider 生成器。
func DefaultGenerators() []FileGenerator {
	return []FileGenerator{}
}

// DefaultGeneratorsWithEngine 返回带有模板引擎的默认文件生成器列表。
// 参数 engine 为模板渲染引擎，fileNamePrefix 表示是否在文件名前添加 schema 名称前缀，force 表示是否强制覆盖。
// 包含 model、repository、service、handler、dto、errors、wire_provider 生成器。
func DefaultGeneratorsWithEngine(engine *TemplateEngine, fileNamePrefix bool, force bool) []FileGenerator {
	return []FileGenerator{
		NewModelGenerator(engine, fileNamePrefix),
		NewRepositoryGenerator(engine, fileNamePrefix),
		NewServiceGenerator(engine, fileNamePrefix),
		NewHandlerGenerator(engine, fileNamePrefix),
		NewDtoGenerator(engine, fileNamePrefix),
		NewErrorsGenerator(engine, fileNamePrefix),
		NewWireProviderGenerator(engine, fileNamePrefix, force),
	}
}

// NewGeneratorWithDeps 创建带有完整依赖的代码生成器实例。
// 参数 config 为生成器配置，schemaDir 为 ent schema 目录路径，templateFS 为模板文件系统。
// 返回初始化后的生成器实例和可能的错误。
func NewGeneratorWithDeps(config Config, schemaDir string, templateFS fs.FS) (*Generator, error) {
	g := NewGenerator(config)

	parser := NewEntSchemaParser(schemaDir)
	g.SetParser(parser)

	engine := NewTemplateEngine()
	if err := engine.LoadFromFS(templateFS, "templates"); err != nil {
		return nil, fmt.Errorf("加载模板失败: %w", err)
	}
	g.SetTemplateEngine(engine)

	g.SetGenerators(DefaultGeneratorsWithEngine(engine, config.FileNamePrefix, config.Force))

	g.SetFileWriter(NewOSFileWriter(config.Force))

	return g, nil
}
