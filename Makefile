# ビルド用の変数
BINARY_NAME=fj
BUILD_DIR=bin
BINARIES=$(BUILD_DIR)/$(BINARY_NAME)
CMD_DIR=./cmd/fj
GO_FILES:= $(shell find . -name '*.go' -type f)

all: build

# ビルドタスク
build: $(BINARIES)

$(BINARIES): $(GO_FILES)
	@echo "Building..."
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)

# 実行ファイルのクリーンアップ
clean:
	rm -rf $(BUILD_DIR)

# 依存関係のダウンロード
deps:
	go mod download

# テストタスク
test:
	go test -v ./...


# old server
SERVER_BINARY=fj-server
SERVER_DIR=./cmd/server
server: $(SERVER_BINARY)
$(SERVER_BINARY): $(GO_FILES)
	CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(SERVER_BINARY) $(SERVER_DIR)

# worker
WORKER_BINARY=fj-worker
WORKER_DIR=./cmd/worker
worker: $(WORKER_BINARY)
$(WORKER_BINARY): $(GO_FILES)
	CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(WORKER_BINARY) $(WORKER_DIR)

# compiler
COMPILER_BINARY=fj-compiler
COMPILER_DIR=./cmd/compiler
compiler: $(COMPILER_BINARY)
$(COMPILER_BINARY): $(GO_FILES)
	CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(COMPILER_BINARY) $(COMPILER_DIR)