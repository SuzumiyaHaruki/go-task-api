/*
task_handler.go 实现任务管理相关 HTTP 接口。

本文件包含任务的列表、创建、详情查询、更新和删除逻辑。
任务数据通过 GORM 保存到 MySQL，每条任务都归属于登录用户，
所有任务接口都需要携带登录后获得的演示 token。
*/
package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

/*
listTasks 返回当前用户的全部任务。

它会从演示 token 中解析用户 ID，并只查询 user_id 等于当前用户 ID 的任务。
当前按 ID 升序查询，并以数组形式返回。
*/
func (a *app) listTasks(c *gin.Context) {
	userID, ok := a.authenticatedUserID(c)
	if !ok {
		return
	}

	var tasks []task
	if err := a.db.Where("user_id = ?", userID).Order("id asc").Find(&tasks).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "list tasks failed")
		return
	}

	writeOK(c, tasks)
}

/*
createTask 创建一条新任务。

该接口会从演示 token 中解析当前用户 ID，并校验任务标题不能为空。
如果请求未提供状态，默认使用 todo；创建成功后 GORM 会写入创建和更新时间。
*/
func (a *app) createTask(c *gin.Context) {
	userID, ok := a.authenticatedUserID(c)
	if !ok {
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

	item := task{
		UserID:  userID,
		Title:   req.Title,
		Content: req.Content,
		Status:  req.Status,
	}

	if err := a.db.Create(&item).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "create task failed")
		return
	}

	writeCreated(c, item)
}

/*
getTask 根据路径 ID 查询单个任务。

如果 ID 不合法会返回 400；如果任务不存在，或任务不属于当前用户，会返回 404。
*/
func (a *app) getTask(c *gin.Context) {
	userID, ok := a.authenticatedUserID(c)
	if !ok {
		return
	}

	id, ok := parseID(c)
	if !ok {
		return
	}

	var item task
	if err := a.db.Where("id = ? AND user_id = ?", id, userID).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(c, http.StatusNotFound, "task not found")
			return
		}
		writeError(c, http.StatusInternalServerError, "query task failed")
		return
	}

	writeOK(c, item)
}

/*
updateTask 根据路径 ID 更新任务。

该接口会从演示 token 中解析当前用户 ID。请求体中非空的 title、content、status
会覆盖原任务字段，并由 GORM 自动刷新 UpdatedAt。
*/
func (a *app) updateTask(c *gin.Context) {
	userID, ok := a.authenticatedUserID(c)
	if !ok {
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

	var item task
	if err := a.db.Where("id = ? AND user_id = ?", id, userID).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(c, http.StatusNotFound, "task not found")
			return
		}
		writeError(c, http.StatusInternalServerError, "query task failed")
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
	if err := a.db.Save(&item).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "update task failed")
		return
	}

	writeOK(c, item)
}

/*
deleteTask 根据路径 ID 删除任务。

该接口会从演示 token 中解析当前用户 ID，只允许删除当前用户自己的任务。
删除成功后返回被删除任务的 ID。
*/
func (a *app) deleteTask(c *gin.Context) {
	userID, ok := a.authenticatedUserID(c)
	if !ok {
		return
	}

	id, ok := parseID(c)
	if !ok {
		return
	}

	result := a.db.Where("user_id = ?", userID).Delete(&task{}, id)
	if result.Error != nil {
		writeError(c, http.StatusInternalServerError, "delete task failed")
		return
	}
	if result.RowsAffected == 0 {
		writeError(c, http.StatusNotFound, "task not found")
		return
	}

	writeOK(c, map[string]int64{"deleted_id": id})
}
