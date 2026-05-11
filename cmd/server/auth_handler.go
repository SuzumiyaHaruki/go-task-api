/*
auth_handler.go 实现账号和认证相关 HTTP 接口。

本文件包含用户注册、登录以及修改当前用户资料的处理逻辑。
注册接口会校验用户名和密码并把用户保存到数据库；登录接口会校验
用户名和密码并返回演示用 Bearer Token；资料修改接口允许当前用户
更新自己的用户名和密码。
*/
package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

/*
register 处理用户注册请求。

它从 JSON 请求体读取用户名和密码，校验用户名非空、密码长度不少于 6 位，
并检查用户名是否已经存在。注册成功后由数据库分配自增用户 ID，
接口会返回新用户信息。
*/
func (a *app) register(c *gin.Context) {
	var req registerRequest
	if !readJSON(c, &req) {
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" || len(req.Password) < 6 {
		writeError(c, http.StatusBadRequest, "username is required and password must be at least 6 characters")
		return
	}

	var existing user
	err := a.db.Where("username = ?", req.Username).First(&existing).Error
	if err == nil {
		writeError(c, http.StatusConflict, "username already exists")
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		writeError(c, http.StatusInternalServerError, "query user failed")
		return
	}

	u := user{
		Username: req.Username,
		Password: req.Password,
	}
	if err := a.db.Create(&u).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "create user failed")
		return
	}

	writeCreated(c, u)
}

/*
login 处理用户登录请求。

它根据请求体中的用户名查询数据库用户记录，并校验密码是否一致。
校验通过后返回演示用 token，供需要鉴权的任务写接口使用。
*/
func (a *app) login(c *gin.Context) {
	var req loginRequest
	if !readJSON(c, &req) {
		return
	}

	var u user
	if err := a.db.Where("username = ?", req.Username).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(c, http.StatusUnauthorized, "invalid username or password")
			return
		}
		writeError(c, http.StatusInternalServerError, "query user failed")
		return
	}

	if u.Password != req.Password {
		writeError(c, http.StatusUnauthorized, "invalid username or password")
		return
	}

	writeOK(c, map[string]string{
		"token": "demo-token-" + strconv.FormatInt(u.ID, 10),
		"type":  "Bearer",
	})
}

/*
updateCurrentUser 修改当前登录用户的用户名或密码。

请求体中的 username 和 password 都是可选字段，但至少需要提供一个。
username 会检查是否和其他用户重复；password 传入时必须不少于 6 个字符。
*/
func (a *app) updateCurrentUser(c *gin.Context) {
	userID, ok := a.authenticatedUserID(c)
	if !ok {
		return
	}

	var req updateUserRequest
	if !readJSON(c, &req) {
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" && req.Password == "" {
		writeError(c, http.StatusBadRequest, "username or password is required")
		return
	}
	if req.Password != "" && len(req.Password) < 6 {
		writeError(c, http.StatusBadRequest, "password must be at least 6 characters")
		return
	}

	var u user
	if err := a.db.First(&u, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(c, http.StatusUnauthorized, "missing or invalid Authorization header")
			return
		}
		writeError(c, http.StatusInternalServerError, "query user failed")
		return
	}

	if req.Username != "" && req.Username != u.Username {
		var existing user
		err := a.db.Where("username = ? AND id <> ?", req.Username, userID).First(&existing).Error
		if err == nil {
			writeError(c, http.StatusConflict, "username already exists")
			return
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(c, http.StatusInternalServerError, "query user failed")
			return
		}
		u.Username = req.Username
	}
	if req.Password != "" {
		u.Password = req.Password
	}

	if err := a.db.Save(&u).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "update user failed")
		return
	}

	writeOK(c, u)
}
