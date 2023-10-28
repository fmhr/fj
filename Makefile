GOPATH:=$(shell go env GOPATH)
GOBIN:=$(GOPATH)/bin
# ビルド用の変数
BINARY_NAME=fj
BUILD_DIR=$(GOBIN)
BINARIES=$(BUILD_DIR)/$(BINARY_NAME)
CMD_DIR=./cmd/fj
GO_FILES:= $(shell find . -name '*.go' -type f)

COMPILER_BINARY=fj-compiler
COMPILER_DIR=./cmd/compiler

WORKER_BINARY=fj-worker
WORKER_DIR=./cmd/worker

# ビルドタスク
.PHONY: build
build: $(BINARIES)
$(BINARIES): $(GO_FILES)
	@echo "Building..."
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)

# 実行ファイルのクリーンアップ
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

# 依存関係のダウンロード
.PHONY: deps
deps:
	go mod download

# テストタスク
.PHONY: test
test:
	go test -v ./...

# worker
.PHONY: worker
worker: $(WORKER_BINARY)
$(WORKER_BINARY): $(GO_FILES)
	@echo "Building... worker bunary..."
	CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(WORKER_BINARY) $(WORKER_DIR)

# compiler
.PHONY: compiler
compiler: $(COMPILER_BINARY)
$(COMPILER_BINARY): $(GO_FILES)
	@echo "Building... compiler bunary..."
	CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(COMPILER_BINARY) $(COMPILER_DIR)
