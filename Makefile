GOPATH:=$(shell go env GOPATH)
GOBIN:=$(GOPATH)/bin
# ビルド用の変数
BINARY_NAME=fj
BUILD_DIR=$(GOBIN)
BINARIES=$(BUILD_DIR)/$(BINARY_NAME)
CMD_DIR=./
GO_FILES:= $(shell find . -name '*.go' -type f)

# バージョン情報の生成
VERSION_DATE:=$(shell date +%Y.%m.%d)
VERSION_FILE:=.version
VERSION_COUNT:=1

# .versionファイルから前回のビルド情報を読み込む
ifneq (,$(wildcard $(VERSION_FILE)))
    LAST_DATE:=$(shell cut -d'.' -f1-3 $(VERSION_FILE))
    LAST_COUNT:=$(shell cut -d'.' -f4 $(VERSION_FILE))
    ifeq ($(VERSION_DATE),$(LAST_DATE))
        VERSION_COUNT:=$(shell echo $$(($(LAST_COUNT) + 1)))
    endif
endif

VERSION:=$(VERSION_DATE).$(VERSION_COUNT)
LDFLAGS:=-X github.com/fmhr/fj/cmd.Version=$(VERSION)

COMPILER_BINARY=fj-compiler
COMPILER_DIR=./cmd/compiler

WORKER_BINARY=fj-worker
WORKER_DIR=./cmd/worker

# ビルドタスク
.PHONY: build
build: $(BINARIES)
$(BINARIES): $(GO_FILES)
	@echo "Building version $(VERSION)..."
	go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)
	@echo $(VERSION) > $(VERSION_FILE)
	@echo "Build complete: $(BINARY_NAME) version $(VERSION)"

# 実行ファイルのクリーンアップ
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)/$(BINARY_NAME)
	rm -f $(VERSION_FILE)

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

