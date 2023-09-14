# ビルド用の変数
BINARY_NAME=fj
BUILD_DIR=bin
BINARIES=$(BUILD_DIR)/$(BINARY_NAME)
CMD_DIR=./cmd
GO_FILES:= $(shell find . -name '*.go' -type f)

all: build

# ビルドタスク
build: $(BINARIES)

$(BINARIES): $(GO_FILES)
	@echo "Building..."
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)

# 実行ファイルのクリーンアップ
clean:
	@echo "Cleaning up..."
	rm -rf $(BUILD_DIR)

# 依存関係のダウンロード
deps:
	@echo "Downloading dependencies..."
	go mod download

# テストタスク
test:
	@echo "Running tests..."
	go test -v ./...


