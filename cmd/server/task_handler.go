/*
task_handler.go 实现任务管理相关 HTTP 接口。

本文件包含任务的列表、创建、详情查询、更新和删除逻辑。
任务数据暂存在应用内存 map 中，写操作需要携带登录后获得的演示 token。
*/
package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

/*
listTasks 返回当前内存中的全部任务。

它会在锁保护下复制任务 map 中的数据，并以数组形式返回。
*/
func (a *app) listTasks(c *gin.Context) {
	a.mu.Lock()
	defer a.mu.Unlock()

	tasks := make([]task, 0, len(a.tasks))
	for _, item := range a.tasks {
		tasks = append(tasks, item)
	}

	writeOK(c, tasks)
}

/*
createTask 创建一条新任务。

该接口要求请求携带演示 token，并校验任务标题不能为空。
如果请求未提供状态，默认使用 todo；创建成功后会写入创建和更新时间。
*/
func (a *app) createTask(c *gin.Context) {
	if !hasDemoToken(c) {
		writeError(c, http.StatusUnauthorized, "missing or invalid Authorization header")
		return
	}

	var req taskRequest
	if !readJSON(c, &req) {
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		writeError(c, http.StatusBadRequest, "title is required")
		return
	}
	if req.Status == "" {
		req.Status = "todo"
	}

	now := time.Now()
	item := task{
		ID:        a.nextID(),
		Title:     req.Title,
		Content:   req.Content,
		Status:    req.Status,
		CreatedAt: now,
		UpdatedAt: now,
	}

	a.mu.Lock()
	a.tasks[item.ID] = item
	a.mu.Unlock()

	writeCreated(c, item)
}

/*
getTask 根据路径 ID 查询单个任务。

如果 ID 不合法会返回 400；如果任务不存在会返回 404。
*/
func (a *app) getTask(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	a.mu.Lock()
	item, exists := a.tasks[id]
	a.mu.Unlock()

	if !exists {
		writeError(c, http.StatusNotFound, "task not found")
		return
	}

	writeOK(c, item)
}

/*
updateTask 根据路径 ID 更新任务。

该接口要求请求携带演示 token。请求体中非空的 title、content、status
会覆盖原任务字段，并刷新 UpdatedAt。
*/
func (a *app) updateTask(c *gin.Context) {
	if !hasDemoToken(c) {
		writeError(c, http.StatusUnauthorized, "missing or invalid Authorization header")
		return
	}

	id, ok := parseID(c)
	if !ok {
		return
	}

	var req taskRequest
	if !readJSON(c, &req) {
		return
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	item, exists := a.tasks[id]
	if !exists {
		writeError(c, http.StatusNotFound, "task not found")
		return
	}

	if strings.TrimSpace(req.Title) != "" {
		item.Title = strings.TrimSpace(req.Title)
	}
	if req.Content != "" {
		item.Content = req.Content
	}
	if req.Status != "" {
		item.Status = req.Status
	}
	item.UpdatedAt = time.Now()
	a.tasks[id] = item

	writeOK(c, item)
}

/*
deleteTask 根据路径 ID 删除任务。

该接口要求请求携带演示 token。删除成功后返回被删除任务的 ID。
*/
func (a *app) deleteTask(c *gin.Context) {
	if !hasDemoToken(c) {
		writeError(c, http.StatusUnauthorized, "missing or invalid Authorization header")
		return
	}

	id, ok := parseID(c)
	if !ok {
		return
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	if _, exists := a.tasks[id]; !exists {
		writeError(c, http.StatusNotFound, "task not found")
		return
	}
	delete(a.tasks, id)

	writeOK(c, map[string]int64{"deleted_id": id})
}
