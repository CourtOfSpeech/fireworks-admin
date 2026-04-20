// Package main 是 CRUD 代码生成器的命令行工具入口。
// 该工具用于根据 Ent Schema 自动生成增删改查代码。
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/speech/fireworks-admin/internal/pkg/generator"
)

// 命令行参数定义。
var (
	schemaName     = flag.String("name", "", "schema 名称，如 User（必填）")
	outputPath     = flag.String("output", "", "输出路径，如 internal/features/user（必填）")
	force          = flag.Bool("force", false, "强制覆盖已存在的文件")
	fileNamePrefix = flag.Bool("prefix", false, "在文件名前添加 schema 名称前缀（如 user_service.go）")
	help           = flag.Bool("help", false, "显示帮助信息")
	list           = flag.Bool("list", false, "列出所有可用的 schema")
)

// 默认配置。
const (
	defaultSchemaDir = "internal/ent/schema"
)

// main 是程序的入口点。
// 解析命令行参数并执行代码生成。
func main() {
	flag.Usage = printUsage
	flag.Parse()

	if *help {
		printUsage()
		os.Exit(0)
	}

	if *list {
		listSchemas()
		os.Exit(0)
	}

	if err := validateFlags(); err != nil {
		fmt.Fprintf(os.Stderr, "参数错误: %v\n\n", err)
		printUsage()
		os.Exit(1)
	}

	if err := runGenerator(); err != nil {
		fmt.Fprintf(os.Stderr, "生成失败: %v\n", err)
		os.Exit(1)
	}
}

// validateFlags 验证命令行参数是否有效。
// 返回参数验证错误，如果参数无效。
func validateFlags() error {
	if *schemaName == "" {
		return fmt.Errorf("-name 参数不能为空")
	}
	if *outputPath == "" {
		return fmt.Errorf("-output 参数不能为空")
	}
	return nil
}

// runGenerator 执行代码生成逻辑。
// 创建生成器实例并生成代码文件。
// 返回生成过程中可能发生的错误。
func runGenerator() error {
	projectRoot, err := getProjectRoot()
	if err != nil {
		return fmt.Errorf("获取项目根目录失败: %w", err)
	}

	schemaDir := filepath.Join(projectRoot, defaultSchemaDir)
	absOutputPath := *outputPath
	if !filepath.IsAbs(absOutputPath) {
		absOutputPath = filepath.Join(projectRoot, *outputPath)
	}

	config := generator.Config{
		SchemaName:     *schemaName,
		OutputPath:     absOutputPath,
		Force:          *force,
		FileNamePrefix: *fileNamePrefix,
	}

	templateFS := generator.GetTemplateFS()
	gen, err := generator.NewGeneratorWithDeps(config, schemaDir, templateFS)
	if err != nil {
		return fmt.Errorf("初始化生成器失败: %w", err)
	}

	result, err := gen.GenerateWithResult()
	if err != nil {
		return err
	}

	printResult(result)
	return nil
}

// getProjectRoot 获取项目根目录。
// 通过查找 go.mod 文件确定项目根目录。
// 返回项目根目录路径和可能的错误。
func getProjectRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	dir := cwd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return cwd, nil
		}
		dir = parent
	}
}

// listSchemas 列出所有可用的 schema。
// 解析 schema 目录并输出可用的 schema 名称列表。
func listSchemas() {
	projectRoot, err := getProjectRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "获取项目根目录失败: %v\n", err)
		return
	}

	schemaDir := filepath.Join(projectRoot, defaultSchemaDir)
	parser := generator.NewEntSchemaParser(schemaDir)

	schemas, err := parser.ListSchemas()
	if err != nil {
		fmt.Fprintf(os.Stderr, "读取 schema 列表失败: %v\n", err)
		return
	}

	if len(schemas) == 0 {
		fmt.Println("未找到任何 schema")
		return
	}

	fmt.Println("可用的 schema:")
	for _, name := range schemas {
		fmt.Printf("  - %s\n", name)
	}
}

// printResult 打印代码生成结果。
// 参数 result 为代码生成结果，包含生成的文件列表和警告信息。
func printResult(result *generator.GenerateResult) {
	fmt.Println("\n========================================")
	fmt.Println("代码生成完成！")
	fmt.Println("========================================")
	fmt.Printf("\nSchema: %s\n", *schemaName)
	fmt.Printf("输出路径: %s\n", *outputPath)

	generatedCount := 0
	skippedCount := 0

	fmt.Println("\n生成的文件:")
	for _, file := range result.Files {
		if file.Skipped {
			fmt.Printf("  [跳过] %s (文件已存在)\n", file.Path)
			skippedCount++
		} else {
			fmt.Printf("  [创建] %s\n", file.Path)
			generatedCount++
		}
	}

	fmt.Printf("\n统计: 创建 %d 个文件", generatedCount)
	if skippedCount > 0 {
		fmt.Printf(", 跳过 %d 个文件", skippedCount)
	}
	fmt.Println()

	if len(result.Warnings) > 0 {
		fmt.Println("\n警告:")
		for _, warning := range result.Warnings {
			fmt.Printf("  - %s\n", warning)
		}
	}

	if len(result.Errors) > 0 {
		fmt.Println("\n错误:")
		for _, err := range result.Errors {
			fmt.Printf("  - %v\n", err)
		}
	}
}

// printUsage 打印命令行使用说明。
func printUsage() {
	fmt.Println("CRUD 代码生成器")
	fmt.Println()
	fmt.Println("该工具根据 Ent Schema 自动生成完整的 CRUD 代码，包括：")
	fmt.Println("  - model.go        领域模型和常量定义")
	fmt.Println("  - repository.go   数据持久化层")
	fmt.Println("  - service.go      业务逻辑层")
	fmt.Println("  - handler.go      HTTP 处理器")
	fmt.Println("  - dto.go          数据传输对象")
	fmt.Println("  - errors.go       错误码和错误处理")
	fmt.Println("  - wire_provider.go Wire 依赖注入")
	fmt.Println()
	fmt.Println("用法:")
	fmt.Println("  generator -name <schema名称> -output <输出路径> [选项]")
	fmt.Println()
	fmt.Println("参数:")
	fmt.Println("  -name string")
	fmt.Println("        schema 名称，如 User（必填）")
	fmt.Println("  -output string")
	fmt.Println("        输出路径，如 internal/features/user（必填）")
	fmt.Println()
	fmt.Println("选项:")
	fmt.Println("  -force")
	fmt.Println("        强制覆盖已存在的文件")
	fmt.Println("  -prefix")
	fmt.Println("        在文件名前添加 schema 名称前缀（如 user_service.go）")
	fmt.Println("  -list")
	fmt.Println("        列出所有可用的 schema")
	fmt.Println("  -help")
	fmt.Println("        显示帮助信息")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  # 生成 User 模块的 CRUD 代码")
	fmt.Println("  generator -name User -output internal/features/user")
	fmt.Println()
	fmt.Println("  # 生成带前缀的文件名")
	fmt.Println("  generator -name User -output internal/features/user -prefix")
	fmt.Println()
	fmt.Println("  # 强制覆盖已存在的文件")
	fmt.Println("  generator -name User -output internal/features/user -force")
	fmt.Println()
	fmt.Println("  # 列出所有可用的 schema")
	fmt.Println("  generator -list")
}
