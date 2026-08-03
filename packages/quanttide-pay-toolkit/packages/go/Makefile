.PHONY: test vet fmt lint

# 运行全部单元测试
test:
	go test ./...

# 静态检查
vet:
	go vet ./...

# 格式化代码
fmt:
	gofmt -w .

# 运行 golangci-lint（需先安装：go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest）
lint:
	golangci-lint run ./...
