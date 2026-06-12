#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR"

if [ -z "${1:-}" ]; then
    echo "错误：请提供版本号，例如: ./build.sh 1.0.0"
    exit 1
fi
VERSION="$1"

mkdir -p build
rm -rf build/web || true

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
    ) || { echo "❌ $target 编译失败"; exit 1; }
    echo "✅ 编译成功"
}

# 编译论坛后端
echo "========================================"
echo "  编译论坛后端"
echo "========================================"
cd "$SCRIPT_DIR/backend"

compile "Linux-AMD64" linux amd64 ../build/server-linux-amd64
compile "Linux-ARM64" linux arm64 ../build/server-linux-arm64
compile "Android-ARM64" android arm64 ../build/server-android-arm64

# 编译IM服务
echo ""
echo "========================================"
echo "  编译IM服务"
echo "========================================"
cd "$SCRIPT_DIR/backend/sdk/im-server/launcher"

compile_im() {
    local target="$1"
    local os="$2"
    local arch="$3"
    local output="$4"
    echo "编译: IM $target -> $output"
    (
        export CGO_ENABLED=0
        export GOOS="$os"
        export GOARCH="$arch"
        go build -o "$output" .
    ) || { echo "❌ IM $target 编译失败"; exit 1; }
    echo "✅ IM编译成功"
}

compile_im "Windows-AMD64" windows amd64 "$SCRIPT_DIR/build/im-server-windows-amd64.exe"
compile_im "Windows-ARM64" windows arm64 "$SCRIPT_DIR/build/im-server-windows-arm64.exe"
compile_im "Linux-AMD64" linux amd64 "$SCRIPT_DIR/build/im-server-linux-amd64"
compile_im "Linux-ARM64" linux arm64 "$SCRIPT_DIR/build/im-server-linux-arm64"

# 编译前端
echo ""
echo "========================================"
echo "  编译前端"
echo "========================================"
cd "$SCRIPT_DIR/frontend"
npm ci
npm run build

mv dist "$SCRIPT_DIR/build/web"

# 创建压缩包
echo ""
echo "========================================"
echo "  创建压缩包"
echo "========================================"
cd "$SCRIPT_DIR/build"
zip -qr "../forum_v${VERSION}.zip" .

echo ""
echo "========================================"
echo "  ✅ 构建完成: forum_v${VERSION}.zip"
echo "========================================"
echo ""
echo "编译产物:"
echo "  - server-*.exe/.sh    (论坛后端)"
echo "  - im-server-*.exe/.sh (IM服务)"
echo "  - web/                (前端)"
echo ""