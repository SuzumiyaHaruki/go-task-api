/*
main.go 是应用程序入口文件。

本文件负责读取服务端口配置、初始化应用实例，并启动 Gin HTTP 服务。
如果启动失败，会记录错误并结束进程。
*/
package main

import (
	"log"
)

/*
main 启动 go-task-api 服务。

它会从 APP_PORT 环境变量读取监听端口，未配置时默认使用 8080，
然后创建应用并在对应端口启动 HTTP 服务器。
*/
func main() {
	port := getenv("APP_PORT", "8080")

	a := newApp()

	log.Printf("go-task-api listening on http://localhost:%s", port)
	if err := a.router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
