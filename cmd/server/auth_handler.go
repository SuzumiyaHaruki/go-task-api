/*
auth_handler.go 实现认证相关 HTTP 接口。

本文件包含用户注册和登录处理逻辑。注册接口会校验用户名和密码，
并把用户保存到内存 map；登录接口会校验用户名和密码，并返回演示用
Bearer Token。
*/
package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

/*
register 处理用户注册请求。

它从 JSON 请求体读取用户名和密码，校验用户名非空、密码长度不少于 6 位，
并检查用户名是否已经存在。注册成功后会分配自增用户 ID，并返回新用户信息。
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

	a.mu.Lock()
	defer a.mu.Unlock()

	if _, exists := a.usersByName[req.Username]; exists {
		writeError(c, http.StatusConflict, "username already exists")
		return
	}

	u := user{
		ID:       a.nextUserID,
		Username: req.Username,
		Password: req.Password,
	}
	a.nextUserID++
	a.usersByName[u.Username] = u

	writeCreated(c, u)
}

/*
login 处理用户登录请求。

它根据请求体中的用户名查找内存中的用户记录，并校验密码是否一致。
校验通过后返回演示用 token，供需要鉴权的任务写接口使用。
*/
func (a *app) login(c *gin.Context) {
	var req loginRequest
	if !readJSON(c, &req) {
		return
	}

	a.mu.Lock()
	u, exists := a.usersByName[req.Username]
	a.mu.Unlock()

	if !exists || u.Password != req.Password {
		writeError(c, http.StatusUnauthorized, "invalid username or password")
		return
	}

	writeOK(c, map[string]string{
		"token": "demo-token-" + strconv.FormatInt(u.ID, 10),
		"type":  "Bearer",
	})
}
