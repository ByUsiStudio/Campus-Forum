@echo off
setlocal EnableDelayedExpansion

if "%~1"=="" (
    echo 错误：请提供版本号，例如: build.bat 1.0.0
    exit /b 1
)
set VERSION=%~1

echo ========================================
echo   Campus Forum 构建脚本 (Windows)
echo ========================================
echo.

:: 创建构建目录
if not exist "build" mkdir build
if exist "build\web" rmdir /s /q "build\web"

:: 编译后端
echo [1/3] 编译后端...
cd backend

echo    - 编译 Windows-AMD64...
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=0
go build -o "..\build\server-windows-amd64.exe" .
if errorlevel 1 (
    echo    [失败] Windows-AMD64
    exit /b 1
)
echo    [成功] Windows-AMD64

echo    - 编译 Linux-AMD64...
set GOOS=linux
set GOARCH=amd64
go build -o "..\build\server-linux-amd64" .
if errorlevel 1 (
    echo    [失败] Linux-AMD64
    exit /b 1
)
echo    [成功] Linux-AMD64

echo    - 编译 Linux-ARM64...
set GOARCH=arm64
go build -o "..\build\server-linux-arm64" .
if errorlevel 1 (
    echo    [失败] Linux-ARM64
    exit /b 1
)
echo    [成功] Linux-ARM64

cd ..

:: 构建前端
echo.
echo [2/3] 构建前端...
cd frontend
call npm install
if errorlevel 1 (
    echo    [失败] npm install
    exit /b 1
)
echo    [成功] npm install

call npm run build
if errorlevel 1 (
    echo    [失败] npm run build
    exit /b 1
)
echo    [成功] npm run build

move /y dist "..\build\web" >nul 2>&1

cd ..

:: 创建压缩包
echo.
echo [3/3] 创建压缩包...
powershell -Command "Compress-Archive -Path 'build\*' -DestinationPath 'forum_v%VERSION%.zip' -Force"
if errorlevel 1 (
    echo    [失败] 创建压缩包
    exit /b 1
)
echo    [成功] forum_v%VERSION%.zip

echo.
echo ========================================
echo   构建完成: forum_v%VERSION%.zip
echo ========================================
echo.
echo 产物位置: %cd%\forum_v%VERSION%.zip
echo.
