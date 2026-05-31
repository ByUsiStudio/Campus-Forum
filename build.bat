@echo off
setlocal EnableDelayedExpansion

if "%~1"=="" (
    echo Error: Please provide version number, e.g.: build.bat 1.0.0
    exit /b 1
)
set VERSION=%~1

echo ========================================
echo   Campus Forum Build Script (Windows)
echo ========================================
echo.

:: Create build directory
if not exist "build" mkdir build
if exist "build\web" rmdir /s /q "build\web"

:: Build backend with version injection
echo [1/3] Building backend...
cd backend

:: Set ldflags for version injection
set LDFLAGS=-X forum/controllers.FrontendVersion=%VERSION% -X forum/controllers.BackendVersion=%VERSION% -X forum/controllers.SwaggerVersion=%VERSION%

echo    - Building Windows-AMD64...
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=0
go build -ldflags "%LDFLAGS%" -o "..\build\server-windows-amd64.exe" .
if errorlevel 1 (
    echo    [FAIL] Windows-AMD64
    exit /b 1
)
echo    [OK] Windows-AMD64

echo    - Building Linux-AMD64...
set GOOS=linux
set GOARCH=amd64
go build -ldflags "%LDFLAGS%" -o "..\build\server-linux-amd64" .
if errorlevel 1 (
    echo    [FAIL] Linux-AMD64
    exit /b 1
)
echo    [OK] Linux-AMD64

echo    - Building Linux-ARM64...
set GOARCH=arm64
go build -ldflags "%LDFLAGS%" -o "..\build\server-linux-arm64" .
if errorlevel 1 (
    echo    [FAIL] Linux-ARM64
    exit /b 1
)
echo    [OK] Linux-ARM64

cd ..

:: Build frontend
echo.
echo [2/3] Building frontend...
cd frontend
call npm install
if errorlevel 1 (
    echo    [FAIL] npm install
    exit /b 1
)
echo    [OK] npm install

call npm run build
if errorlevel 1 (
    echo    [FAIL] npm run build
    exit /b 1
)
echo    [OK] npm run build

move /y dist "..\build\web" >nul 2>&1

cd ..

:: Create zip archive
echo.
echo [3/3] Creating archive...
powershell -Command "Compress-Archive -Path 'build\*' -DestinationPath 'forum_v%VERSION%.zip' -Force"
if errorlevel 1 (
    echo    [FAIL] Create archive
    exit /b 1
)
echo    [OK] forum_v%VERSION%.zip

echo.
echo ========================================
echo   Build Complete: forum_v%VERSION%.zip
echo ========================================
echo.
echo Output: %cd%\forum_v%VERSION%.zip
echo.
