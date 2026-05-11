# go-task-api

一个用于练习 Go 后端开发的任务管理 API 项目。

## 当前技术栈

- Go
- Gin
- GORM
- MySQL
- Docker

## 后续计划

安全与认证：

- 使用 bcrypt 对用户密码进行哈希存储
- 使用正式 JWT 替换演示 token
- 使用 Redis 管理登录状态、token 黑名单或缓存数据
- 将鉴权逻辑改造成 Gin 中间件，减少 handler 中的重复校验

接口能力：

- 为任务列表增加分页查询
- 为任务列表增加 status 状态筛选
- 限制任务状态只能使用 todo、doing、done、cancelled 等固定值
- 为任务增加 priority、due_date、completed_at 等字段
- 支持删除当前用户账号，并验证用户和任务的级联删除

代码结构：

- 拆分 handler、service、repository 层，减少 handler 直接操作数据库
- 将环境变量和数据库配置集中到 config 模块
- 设计统一业务错误码，避免只使用 HTTP 状态码作为 code

工程质量：

- 增加注册、登录、任务隔离等核心流程测试
- 增强 Makefile，加入 docker-up、docker-down、tidy、lint 等命令
- 引入数据库迁移工具管理表结构版本，逐步替代 AutoMigrate

## 项目结构

```text
go-task-api/
  cmd/server/
    main.go            程序入口和 HTTP Server 启动
    app.go             应用初始化和依赖组装
    database.go        GORM 和 MySQL 连接初始化
    routes.go          路由注册
    models.go          请求、响应和数据库模型
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
docker compose up --build
```

服务启动后访问：

```text
http://localhost:8080
```

## 当前功能

- 健康检查
- 用户注册和登录
- 修改当前用户用户名和密码
- 返回演示用 token
- 基于 GORM 和 MySQL 的任务 CRUD
- 每个用户只能管理自己的任务
- 启动时自动迁移 users 和 tasks 表
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
curl -X PUT http://localhost:8080/api/v1/users/me \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer demo-token-1' \
  -d '{"username":"alice_new","password":"abcdef"}'
```

```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer demo-token-1' \
  -d '{"title":"learn Go API","content":"build the first simple version"}'
```

```bash
curl http://localhost:8080/api/v1/tasks \
  -H 'Authorization: Bearer demo-token-1'
```

## 本地运行

如果已经在本机启动 MySQL，并创建了 `.env.example` 中对应的数据库和账号，
也可以直接运行：

```bash
make run
```
