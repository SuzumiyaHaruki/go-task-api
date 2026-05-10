/*
health_handler.go 实现服务健康检查接口。

健康检查接口用于让开发者、部署平台或监控系统快速确认 HTTP 服务
是否已经启动并能正常响应请求。
*/
package main

import "github.com/gin-gonic/gin"

/*
health 返回服务运行状态。

当前实现固定返回 status=ok，表示 API 进程可访问。
*/
func (a *app) health(c *gin.Context) {
	writeOK(c, map[string]string{"status": "ok"})
}
