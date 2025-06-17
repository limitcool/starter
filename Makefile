# 变量定义
APP_NAME = starter
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT ?= $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
BUILD_TIME ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
BUILD_DIR = build

# 错误码生成配置
ERROR_MD = tools/errorgen/error_codes.md
ERROR_CODE_FILE = internal/pkg/errorx/code_gen.go

# 默认架构
ARCH ?= amd64

# 编译参数
LDFLAGS = -s -w \
	-X github.com/limitcool/starter/internal/version.Version=$(VERSION) \
	-X github.com/limitcool/starter/internal/version.GitCommit=$(GIT_COMMIT) \
	-X github.com/limitcool/starter/internal/version.BuildDate=$(BUILD_TIME)

# 构建标志
BUILDFLAGS = -trimpath

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
	GOOS=linux GOARCH=$(ARCH) go build $(BUILDFLAGS) -ldflags="$(LDFLAGS)" \
		-o $(BUILD_DIR)/$(APP_NAME)-linux-$(ARCH)

# 构建 Windows 版本
.PHONY: windows
windows: $(BUILD_DIR)
	@echo "Building for Windows ($(ARCH))..."
	GOOS=windows GOARCH=$(ARCH) go build $(BUILDFLAGS) -ldflags="$(LDFLAGS)" \
		-o $(BUILD_DIR)/$(APP_NAME)-windows-$(ARCH).exe

# 构建 MacOS 版本
.PHONY: darwin
darwin: $(BUILD_DIR)
	@echo "Building for MacOS ($(ARCH))..."
	GOOS=darwin GOARCH=$(ARCH) go build $(BUILDFLAGS) -ldflags="$(LDFLAGS)" \
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

# 显示版本信息
.PHONY: version
version:
	@echo "Project: $(APP_NAME)"
	@echo "Version: $(VERSION)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Build Time: $(BUILD_TIME)"

# 快速构建（开发用）
.PHONY: build-dev
build-dev:
	@echo "Building $(APP_NAME) for development..."
	go build $(BUILDFLAGS) -ldflags="$(LDFLAGS)" -o $(APP_NAME)$(if $(filter windows,$(shell go env GOOS)),.exe,)

# Docker 相关目标
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	./deploy/build.sh $(VERSION)

.PHONY: docker-run
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 -p 6060:6060 --name $(APP_NAME)-container $(APP_NAME):$(VERSION)

.PHONY: docker-run-bg
docker-run-bg:
	@echo "Running Docker container in background..."
	docker run -d -p 8080:8080 -p 6060:6060 --name $(APP_NAME)-container $(APP_NAME):$(VERSION)

.PHONY: docker-stop
docker-stop:
	@echo "Stopping Docker container..."
	docker stop $(APP_NAME)-container || true
	docker rm $(APP_NAME)-container || true

.PHONY: docker-logs
docker-logs:
	@echo "Showing Docker container logs..."
	docker logs -f $(APP_NAME)-container

.PHONY: docker-version
docker-version:
	@echo "Checking Docker image version..."
	docker run --rm $(APP_NAME):$(VERSION) version

# 生成错误码
.PHONY: gen-errors
gen-errors:
	@echo "Generating error codes from $(ERROR_MD) to $(ERROR_CODE_FILE)..."
	go run tools/errorgen/main.go $(ERROR_MD) $(ERROR_CODE_FILE)

# 生成proto文件
.PHONY: proto
proto:
	@echo "Generating protobuf code..."
	# 生成proto代码
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative internal/proto/v1/system.proto
	# 移动生成的文件到gen目录
	mkdir -p internal/proto/gen/v1
	mv internal/proto/v1/*.pb.go internal/proto/gen/v1/

# 帮助信息
.PHONY: help
help:
	@echo "Make targets:"
	@echo "  linux         - Build for Linux (ARCH=amd64/arm64)"
	@echo "  windows       - Build for Windows (ARCH=amd64/arm64)"
	@echo "  darwin        - Build for MacOS (ARCH=amd64/arm64)"
	@echo "  all-arch      - Build for all platforms (amd64)"
	@echo "  all-arm64     - Build for all platforms (arm64)"
	@echo "  build-dev     - Quick build for development"
	@echo "  clean         - Clean build directory"
	@echo "  test          - Run tests"
	@echo "  run           - Run the application"
	@echo "  version       - Show version information"
	@echo "  gen-errors    - Generate error codes from Markdown definition"
	@echo "  proto         - Generate protobuf code from proto files"
	@echo ""
	@echo "Docker targets:"
	@echo "  docker-build  - Build Docker image with version info"
	@echo "  docker-run    - Run Docker container (foreground)"
	@echo "  docker-run-bg - Run Docker container (background)"
	@echo "  docker-stop   - Stop and remove Docker container"
	@echo "  docker-logs   - Show Docker container logs"
	@echo "  docker-version- Check Docker image version"
	@echo ""
	@echo "Examples:"
	@echo "  make linux ARCH=arm64     - Build Linux arm64 version"
	@echo "  make darwin ARCH=arm64    - Build MacOS arm64 version"
	@echo "  make all-arch             - Build all platforms in amd64"
	@echo "  make all-arm64            - Build all platforms in arm64"
	@echo "  make build-dev            - Quick development build"
	@echo "  make version              - Show build version information"
	@echo "  make docker-build         - Build Docker image"
	@echo "  make docker-run VERSION=v1.0.0 - Run specific version"
	@echo "  make gen-errors           - Generate error codes from $(ERROR_MD)"
	@echo "  make proto                - Generate protobuf code"