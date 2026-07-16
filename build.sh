#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR"

if [ -z "${1:-}" ]; then
    echo "错误：请提供版本号，例如: ./build.sh 3.0.0"
    exit 1
fi
VERSION="$1"

mkdir -p build
rm -rf build/web || true

echo "复制 config.json..."
cp "$SCRIPT_DIR/backend/config.json" "$SCRIPT_DIR/build/config.json" || { echo "复制 config.json 失败"; exit 1; }
echo "OK config.json"

LDFLAGS="-X forum/controllers.FrontendVersion=${VERSION} -X forum/controllers.BackendVersion=${VERSION} -X forum/controllers.SwaggerVersion=${VERSION}"

compile() {
    local target="$1"
    local os="$2"
    local arch="$3"
    local output="$4"
    echo "编译: $target -> $output"
    (
        export CGO_ENABLED=0
        export GOOS="$os"
        export GOARCH="$arch"
        go build -ldflags "$LDFLAGS" -o "$output" .
    ) || { echo "$target 编译失败"; exit 1; }
    echo "编译成功"
}

echo "========================================"
echo "  编译论坛后端"
echo "  版本: ${VERSION}"
echo "========================================"
cd "$SCRIPT_DIR/backend"

compile "Linux-AMD64" linux amd64 ../build/server-linux-amd64
compile "Linux-ARM64" linux arm64 ../build/server-linux-arm64
compile "Android-ARM64" android arm64 ../build/server-android-arm64

echo ""
echo "========================================"
echo "  编译前端"
echo "========================================"
cd "$SCRIPT_DIR/frontend"
npm ci
npm run build

mv dist "$SCRIPT_DIR/build/web"

echo ""
echo "========================================"
echo "  创建压缩包"
echo "========================================"
cd "$SCRIPT_DIR/build"
zip -qr "../forum_v${VERSION}.zip" .

echo ""
echo "========================================"
echo "  构建完成: forum_v${VERSION}.zip"
echo "========================================"
echo ""
echo "编译产物:"
echo "  - server-linux-amd64      (Linux x64 后端)"
echo "  - server-linux-arm64      (Linux ARM64 后端)"
echo "  - server-android-arm64    (Android ARM64 后端)"
echo "  - web/                    (前端)"
echo ""
echo "架构: Repository-Service-Controller"
echo ""
