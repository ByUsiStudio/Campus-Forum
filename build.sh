set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR/backend"

if [ -z "${1:-}" ]; then
    echo "错误：请提供版本号，例如: ./build.sh 1.0.0"
    exit 1
fi
VERSION="$1"

mkdir -p ../build
rm -rf ../build/web || true

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
        go build -o "$output" . 
    ) || { echo "❌ $target 编译失败"; exit 1; }
    echo "✅ 编译成功"
}

compile "Windows-AMD64" windows amd64 ../build/server-windows-amd64.exe
compile "Windows-ARM64" windows arm64 ../build/server-windows-arm64.exe
compile "Linux-AMD64" linux amd64 ../build/server-linux-amd64
compile "Linux-ARM64" linux arm64 ../build/server-linux-arm64
compile "Android-ARM64" android arm64 ../build/server-android-arm64

cd ../frontend
npm ci
npm run build

mv dist ../build/web

cd ../build
zip -qr "../forum_v${VERSION}.zip" .

echo ""
echo "✅ 构建完成: forum_v${VERSION}.zip"
echo "⚠️  PHP测试服务器已启动于 0.0.0.0:8000 (按 Ctrl+C 退出)"
php -S 0.0.0.0:8000 -t ./web