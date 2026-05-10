/*
app.go 定义应用的核心运行对象。

本文件负责创建 Gin 路由器、初始化内存中的用户和任务存储，
并集中挂载全局中间件和业务路由。当前项目使用内存 map 保存数据，
因此这里也维护自增 ID 和互斥锁，保证并发请求下读写状态更稳定。
*/
package main

import (
	"sync"

	"github.com/gin-gonic/gin"
)

type app struct {
	router      *gin.Engine
	mu          sync.Mutex
	nextUserID  int64
	nextTaskID  int64
	usersByName map[string]user
	tasks       map[int64]task
}

/*
newApp 创建并初始化一个应用实例。

它会创建 Gin Engine、初始化用户和任务的内存存储、
关闭可信代理默认配置，并注册请求日志、异常恢复中间件和所有路由。
*/
func newApp() *app {
	a := &app{
		router:      gin.New(),
		nextUserID:  1,
		nextTaskID:  1,
		usersByName: make(map[string]user),
		tasks:       make(map[int64]task),
	}

	if err := a.router.SetTrustedProxies(nil); err != nil {
		panic(err)
	}

	a.router.Use(logRequests(), gin.Recovery())
	a.routes()
	return a
}

/*
nextID 生成新的任务 ID。

该方法通过互斥锁保护自增计数器，避免并发创建任务时产生重复 ID。
*/
func (a *app) nextID() int64 {
	a.mu.Lock()
	defer a.mu.Unlock()

	id := a.nextTaskID
	a.nextTaskID++
	return id
}
