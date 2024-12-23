# 变量定义
APP_NAME = starter
VERSION ?= $(shell git describe --tags --always)
BUILD_DIR = build

# 默认架构
ARCH ?= amd64

# 编译参数
LDFLAGS = -s -w \
	-X main.Version=$(VERSION) \
	-X main.BuildTime=$(shell date -u +%Y-%m-%d_%H:%M:%S)

# 默认目标
.PHONY: all
all: linux

# 创建构建目录
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

# 构建 Linux 版本
.PHONY: linux
linux: $(BUILD_DIR)
	@echo "Building for Linux ($(ARCH))..."
	GOOS=linux GOARCH=$(ARCH) go build -ldflags="$(LDFLAGS)" \
		-o $(BUILD_DIR)/$(APP_NAME)-linux-$(ARCH)

# 构建 Windows 版本
.PHONY: windows
windows: $(BUILD_DIR)
	@echo "Building for Windows ($(ARCH))..."
	GOOS=windows GOARCH=$(ARCH) go build -ldflags="$(LDFLAGS)" \
		-o $(BUILD_DIR)/$(APP_NAME)-windows-$(ARCH).exe

# 构建 MacOS 版本
.PHONY: darwin
darwin: $(BUILD_DIR)
	@echo "Building for MacOS ($(ARCH))..."
	GOOS=darwin GOARCH=$(ARCH) go build -ldflags="$(LDFLAGS)" \
		-o $(BUILD_DIR)/$(APP_NAME)-darwin-$(ARCH)

# 构建所有平台 amd64 版本
.PHONY: all-arch
all-arch: linux windows darwin

# 构建所有平台 arm64 版本
.PHONY: all-arm64
all-arm64: ARCH=arm64
all-arm64: linux windows darwin

# 清理构建产物
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

# 运行测试
.PHONY: test
test:
	go test -v ./...

# 运行程序
.PHONY: run
run:
	go run main.go

# 帮助信息
.PHONY: help
help:
	@echo "Make targets:"
	@echo "  linux    - Build for Linux (ARCH=amd64/arm64)"
	@echo "  windows  - Build for Windows (ARCH=amd64/arm64)"
	@echo "  darwin   - Build for MacOS (ARCH=amd64/arm64)"
	@echo "  all-arch - Build for all platforms (amd64)"
	@echo "  all-arm64 - Build for all platforms (arm64)"
	@echo "  clean    - Clean build directory"
	@echo "  test     - Run tests"
	@echo "  run      - Run the application"
	@echo ""
	@echo "Examples:"
	@echo "  make linux ARCH=arm64   - Build Linux arm64 version"
	@echo "  make darwin ARCH=arm64  - Build MacOS arm64 version"
	@echo "  make all-arch          - Build all platforms in amd64"
	@echo "  make all-arm64         - Build all platforms in arm64"