cd backend
go mod tidy
echo "正在编译Windows的AMD64版校园论坛服务端...."
export CGO_ENABLED=0
export GOOS=windows
export GOARCH=amd64
go build -o ../build/forum.exe .
echo "编译成功"

echo "正在编译Windows的ARM64版校园论坛服务端...."
export CGO_ENABLED=0
export GOOS=windows
export GOARCH=arm64
go build -o ../build/forum-arm64.exe .
echo "编译成功"

echo "正在编译Linux的AMD64版校园论坛服务端...."
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64
go build -o ../build/forum-amd64 .
echo "编译成功"

echo "正在编译MacOS的AMD64版校园论坛服务端...."
export CGO_ENABLED=0
export GOOS=darwin
export GOARCH=amd64
go build -o ../build/forum-apple-amd64 .
echo "编译成功"

echo "正在编译MacOS的ARM64版校园论坛服务端...."
export CGO_ENABLED=0
export GOOS=darwin
export GOARCH=arm64
go build -o ../build/forum-apple-arm64 .
echo "编译成功"

echo "正在编译Linux的ARM64版校园论坛服务端...."
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=arm64
go build -o ../build/forum-arm64 .
echo "编译成功"

echo "正在编译Linux的ARMv7版校园论坛服务端...."
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=arm
export GOARM=7
go build -o ../build/forum-armv7 .
echo "编译成功"

echo "正在编译Linux的Mipsle版校园论坛服务端...."
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=mipsle
go build -o ../build/forum-mips .
echo "编译成功"

echo "正在编译Android的ARM64版校园论坛服务端...."
export CGO_ENABLED=0
export GOOS=android
export GOARCH=arm64
go build -o ../build/forum-android .
echo "编译成功"
cd ../frontend
npm install
npm run build
echo " "
echo "正在转移编译产物中...."
rm -rvf ../build/web
mv dist ../build/web
echo "正在打包编译后产物中..."
cd ../build
zip -r ../forum_v$1.zip .
echo " "
echo "正在运行测试，如果想退出测试，那么请 Ctrl + C"
php -S 0.0.0.0:8000 -t ./web