/*
app.go 定义应用的核心运行对象。

本文件负责创建 Gin 路由器、初始化 GORM 数据库连接，
并集中挂载全局中间件和业务路由。用户和任务数据由数据库持久化管理，
应用实例只保存路由器和数据库连接。
*/
package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type app struct {
	router *gin.Engine
	db     *gorm.DB
}

/*
newApp 创建并初始化一个应用实例。

它会创建 Gin Engine、初始化 GORM 数据库连接、
执行数据库表结构自动迁移，关闭可信代理默认配置，
并注册请求日志、异常恢复中间件和所有路由。
*/
func newApp() (*app, error) {
	db, err := openDatabase()
	if err != nil {
		return nil, err
	}

	a := &app{
		router: gin.New(),
		db:     db,
	}

	if err := a.router.SetTrustedProxies(nil); err != nil {
		return nil, err
	}

	a.router.Use(logRequests(), gin.Recovery())
	a.routes()
	return a, nil
}
