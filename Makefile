# ビルド用の変数
BINARY_NAME=fmj
BUILD_DIR=bin
CMD_DIR=./cmd

all: build

# ビルドタスク
build:
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


