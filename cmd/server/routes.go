/*
routes.go 集中注册 HTTP 路由。

本文件把健康检查、认证接口和任务管理接口挂载到 Gin 路由器上，
统一维护 API 版本前缀和各资源路径。
*/
package main

/*
routes 注册应用所有路由。

/health 用于健康检查；/api/v1/auth 下提供注册和登录；
/api/v1/tasks 下提供任务列表、创建、详情、更新和删除接口。
*/
func (a *app) routes() {
	a.router.GET("/health", a.health)

	api := a.router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", a.register)
			auth.POST("/login", a.login)
		}

		tasks := api.Group("/tasks")
		{
			tasks.GET("", a.listTasks)
			tasks.POST("", a.createTask)
			tasks.GET("/:id", a.getTask)
			tasks.PUT("/:id", a.updateTask)
			tasks.DELETE("/:id", a.deleteTask)
		}
	}
}
