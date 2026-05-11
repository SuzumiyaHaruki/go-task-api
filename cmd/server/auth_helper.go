/*
auth_helper.go 提供需要数据库参与的认证辅助逻辑。

本文件基于演示 token 解析当前用户 ID，并确认该用户仍然存在于数据库中。
任务接口和用户资料接口都会复用这里的函数。
*/
package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

/*
authenticatedUserID 解析并校验当前请求的用户 ID。

它先从演示 token 中解析用户 ID，再确认该用户仍然存在于数据库中。
校验失败时会直接写入 401 响应。
*/
func (a *app) authenticatedUserID(c *gin.Context) (int64, bool) {
	userID, ok := parseDemoUserID(c)
	if !ok {
		return 0, false
	}

	var u user
	if err := a.db.Select("id").First(&u, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(c, http.StatusUnauthorized, "missing or invalid Authorization header")
			return 0, false
		}
		writeError(c, http.StatusInternalServerError, "query user failed")
		return 0, false
	}

	return userID, true
}
