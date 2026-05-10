# go-task-api

一个用于练习 Go 后端开发的任务管理 API 项目。

## 当前技术栈

- Go
- Gin
- 内存存储
- Docker

## 后续计划

- 使用 GORM 接入 MySQL
- 使用 Redis 做缓存和登录状态管理
- Docker Compose

## 项目结构

```text
go-task-api/
  cmd/server/
    main.go            程序入口和 HTTP Server 启动
    app.go             应用初始化和内存状态
    routes.go          路由注册
    models.go          请求、响应和业务模型
    health_handler.go  健康检查接口
    auth_handler.go    注册和登录接口
    task_handler.go    任务 CRUD 接口
    response.go        JSON 响应、请求解析和通用工具
    middleware.go      请求日志中间件
  docs/                API 文档
```

## 快速启动

```bash
cp .env.example .env
make run
```

服务启动后访问：

```text
http://localhost:8080
```

## 当前功能

- 健康检查
- 用户注册和登录
- 返回演示用 token
- 基于内存的任务 CRUD
- 统一 JSON 响应格式
- 请求日志

## 接口测试

```bash
curl http://localhost:8080/health
```

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"username":"alice","password":"123456"}'
```

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"alice","password":"123456"}'
```

```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer demo-token-1' \
  -d '{"title":"learn Go API","content":"build the first simple version"}'
```

```bash
curl http://localhost:8080/api/v1/tasks
```

## Docker

使用 Docker Compose 启动本地开发环境：

```bash
cp .env.example .env
docker compose up --build
```
