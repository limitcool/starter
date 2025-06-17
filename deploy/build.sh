#!/bin/bash

# Docker 构建脚本
# 用法: ./deploy/build.sh [tag]

set -e

# 获取项目根目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

# 切换到项目根目录
cd "${PROJECT_ROOT}"

# 获取版本信息
VERSION=${1:-$(git describe --tags --always --dirty 2>/dev/null || echo "dev")}
GIT_COMMIT=$(git rev-parse HEAD 2>/dev/null || echo "unknown")
BUILD_DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ)

# Docker 镜像信息
IMAGE_NAME="starter"
REGISTRY=${REGISTRY:-""}
FULL_IMAGE_NAME="${REGISTRY}${IMAGE_NAME}"

echo "=========================================="
echo "Building Docker image"
echo "=========================================="
echo "Image: ${FULL_IMAGE_NAME}:${VERSION}"
echo "Version: ${VERSION}"
echo "Git Commit: ${GIT_COMMIT}"
echo "Build Date: ${BUILD_DATE}"
echo "=========================================="

# 构建 Docker 镜像
docker build \
    --file deploy/Dockerfile \
    --build-arg VERSION="${VERSION}" \
    --build-arg GIT_COMMIT="${GIT_COMMIT}" \
    --build-arg BUILD_DATE="${BUILD_DATE}" \
    --tag "${FULL_IMAGE_NAME}:${VERSION}" \
    --tag "${FULL_IMAGE_NAME}:latest" \
    .

echo "=========================================="
echo "Build completed successfully!"
echo "=========================================="
echo "Images created:"
echo "  ${FULL_IMAGE_NAME}:${VERSION}"
echo "  ${FULL_IMAGE_NAME}:latest"
echo ""
echo "To run the container:"
echo "  docker run -p 8080:8080 ${FULL_IMAGE_NAME}:${VERSION}"
echo ""
echo "To check version:"
echo "  docker run --rm ${FULL_IMAGE_NAME}:${VERSION} version"
echo "=========================================="
