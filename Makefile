# Makefile 汇总项目常用开发命令。
#
# run 用于本地启动服务，test 用于运行全部 Go 测试，
# build 用于编译服务二进制文件，fmt 用于格式化 Go 代码。
.PHONY: run test build fmt

run:
	go run ./cmd/server

test:
	go test ./...

build:
	go build -o bin/server ./cmd/server

fmt:
	go fmt ./...
