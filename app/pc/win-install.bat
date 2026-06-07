@echo off
cd /d "%~dp0"
echo ============================
echo   校园论坛 PC 客户端 - 安装
echo ============================
echo.

where python >nul 2>nul
if %ERRORLEVEL% neq 0 (
    echo [错误] 未找到 Python
    pause
    exit /b 1
)

echo [信息] 安装依赖...
python -m pip install -r requirements.txt -q
if %ERRORLEVEL% neq 0 (
    echo [错误] 安装依赖失败
    pause
    exit /b 1
)

echo.
echo [信息] 安装完成！
echo [信息] 请运行 run.bat 启动客户端
echo.
pause
