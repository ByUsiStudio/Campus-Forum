@echo off
cd /d "%~dp0"
echo ============================
echo   校园论坛 PC 客户端
echo ============================
echo.

REM 检查 Python
where python >nul 2>nul
if %ERRORLEVEL% neq 0 (
    echo [错误] 未找到 Python，请先安装 Python 3.8+
    pause
    exit /b 1
)

REM 安装依赖
echo [信息] 检查依赖...
python -m pip install -r requirements.txt -q

echo.
echo [信息] 启动客户端...
echo.
python main.py
pause
