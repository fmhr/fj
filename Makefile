# ビルド用の変数
BINARY_NAME=fj
BUILD_DIR=bin
BINARIES=$(BUILD_DIR)/$(BINARY_NAME)
CMD_DIR=./cmd/fj
GO_FILES:= $(shell find . -name '*.go' -type f)

# server
SERVER_BINARY=fj-server
SERVER_DIR=./cmd/server

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


server: $(SERVER_BINARY)

$(SERVER_BINARY): $(GO_FILES)
	go build -o $(BUILD_DIR)/$(SERVER_BINARY) $(SERVER_DIR)
