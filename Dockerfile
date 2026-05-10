# Dockerfile 用于构建 go-task-api 服务镜像。
#
# 第一阶段使用 Go 镜像编译 cmd/server 下的服务入口；
# 第二阶段使用轻量 Alpine 镜像运行编译产物，减少最终镜像体积。
FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
COPY . .
RUN go build -o /app/bin/server ./cmd/server

FROM alpine:3.20

WORKDIR /app
COPY --from=builder /app/bin/server /app/server
EXPOSE 8080
CMD ["/app/server"]
