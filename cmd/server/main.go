// Package main 是 fireworks-admin 服务的程序入口包。
// 该包负责初始化并启动整个应用程序。
package main

import (
	"github.com/speech/fireworks-admin/internal/app"
)

// main 是程序的入口点。
// 该函数调用 app.Run 启动应用程序的主逻辑。
func main() {
	app.Run()
}
