/*
middleware.go 定义项目使用的 Gin 中间件。

当前文件提供请求日志中间件，用于记录每次 HTTP 请求的方法、路径和耗时，
帮助开发和排查接口调用情况。
*/
package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

/*
logRequests 创建请求日志中间件。

该中间件会在请求进入时记录开始时间，请求处理完成后输出 HTTP 方法、
请求路径和总耗时。
*/
func logRequests() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		log.Printf("%s %s %s", c.Request.Method, c.Request.URL.Path, time.Since(start))
	}
}
