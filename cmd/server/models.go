/*
models.go 定义服务使用的数据结构。

本文件集中保存 API 统一响应格式、认证请求体、任务请求体，
以及当前内存存储使用的用户和任务模型。
*/
package main

import "time"

/*
apiResponse 是统一的 HTTP JSON 响应结构。

Code 表示业务状态码或 HTTP 状态码，Message 表示可读提示信息，
Data 用于承载接口返回的数据。
*/
type apiResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

/*
registerRequest 表示注册接口的 JSON 请求体。

Username 是用户登录名，Password 是用户密码。
*/
type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

/*
loginRequest 表示登录接口的 JSON 请求体。

它复用用户名和密码字段，用于校验用户身份。
*/
type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

/*
taskRequest 表示创建或更新任务时提交的 JSON 请求体。

Title 是任务标题，Content 是任务内容，Status 是任务状态。
更新任务时允许只传部分字段。
*/
type taskRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Status  string `json:"status"`
}

/*
user 表示内存中的用户记录。

Password 字段不会序列化到 JSON 响应中，避免接口直接返回密码。
*/
type user struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
}

/*
task 表示内存中的任务记录。

它包含任务基础信息、当前状态，以及创建和最后更新时间。
*/
type task struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
