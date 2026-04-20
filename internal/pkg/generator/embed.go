// Package generator 提供 CRUD 代码生成功能。
// 该文件使用 embed 嵌入模板文件，提供模板文件系统访问。
package generator

import "embed"

//go:embed all:templates
var templateFS embed.FS

// GetTemplateFS 获取嵌入的模板文件系统。
// 返回模板文件系统接口，用于代码生成器加载模板。
func GetTemplateFS() embed.FS {
	return templateFS
}
