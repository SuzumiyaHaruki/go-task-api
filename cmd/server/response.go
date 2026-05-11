/*
response.go 提供 HTTP 请求和响应辅助函数。

本文件负责读取 JSON 请求体、解析路径参数、解析演示 token，
并用统一的 apiResponse 格式写出成功或错误响应。同时也提供读取环境变量
并设置默认值的工具函数。
*/
package main

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

/*
readJSON 将请求 JSON 绑定到目标结构体。

如果请求体不是合法 JSON，函数会直接写入 400 错误响应并返回 false；
调用方可根据返回值决定是否继续处理业务逻辑。
*/
func readJSON(c *gin.Context, dst interface{}) bool {
	if err := c.ShouldBindJSON(dst); err != nil {
		writeError(c, http.StatusBadRequest, "invalid json body")
		return false
	}

	return true
}

/*
parseID 从路由参数中解析正整数 ID。

当路径中的 id 缺失、不是整数或小于等于 0 时，会写入 400 错误响应，
并通过第二个返回值通知调用方解析失败。
*/
func parseID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		writeError(c, http.StatusBadRequest, "invalid id")
		return 0, false
	}

	return id, true
}

/*
parseDemoUserID 从演示用 Bearer Token 中解析当前用户 ID。

当前 token 格式为 "Bearer demo-token-{userID}"，例如登录接口返回的
"Bearer demo-token-1" 会解析出用户 ID 1。该实现只适合示例项目使用，
不等同于生产环境的 JWT 鉴权。
*/
func parseDemoUserID(c *gin.Context) (int64, bool) {
	value := c.GetHeader("Authorization")
	const prefix = "Bearer demo-token-"
	if !strings.HasPrefix(value, prefix) {
		writeError(c, http.StatusUnauthorized, "missing or invalid Authorization header")
		return 0, false
	}

	id, err := strconv.ParseInt(strings.TrimPrefix(value, prefix), 10, 64)
	if err != nil || id <= 0 {
		writeError(c, http.StatusUnauthorized, "missing or invalid Authorization header")
		return 0, false
	}

	return id, true
}

/*
writeOK 写入 200 成功响应。

data 会被放入统一响应结构的 Data 字段中返回给客户端。
*/
func writeOK(c *gin.Context, data interface{}) {
	writeJSON(c, http.StatusOK, apiResponse{Code: 0, Message: "ok", Data: data})
}

/*
writeCreated 写入 201 创建成功响应。

通常用于注册用户或创建任务等资源创建成功的场景。
*/
func writeCreated(c *gin.Context, data interface{}) {
	writeJSON(c, http.StatusCreated, apiResponse{Code: 0, Message: "created", Data: data})
}

/*
writeError 写入错误响应。

status 同时作为 HTTP 状态码和响应体中的 Code，message 用于描述错误原因。
*/
func writeError(c *gin.Context, status int, message string) {
	writeJSON(c, status, apiResponse{Code: status, Message: message})
}

/*
writeJSON 按统一结构写出 JSON 响应。

它封装 Gin 的 JSON 输出，并额外写入换行，方便命令行工具查看响应内容。
*/
func writeJSON(c *gin.Context, status int, resp apiResponse) {
	c.JSON(status, resp)
	c.Writer.Write([]byte("\n"))
}

/*
getenv 读取环境变量并提供默认值。

当环境变量不存在或内容为空白字符串时，会返回 fallback。
*/
func getenv(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}
