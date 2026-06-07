import os
import sys
import json
import customtkinter as ctk
from src.api_client import ApiClient
from src.views.auth_views import AuthWindow
from src.views.home_view import MainWindow

# ── 路径设置 ──
BASE_DIR = os.path.dirname(os.path.abspath(__file__))
CONFIG_PATH = os.path.join(BASE_DIR, "config.json5")


def load_config():
    """加载配置文件"""
    try:
        with open(CONFIG_PATH, "r", encoding="utf-8") as f:
            text = f.read()
            # 移除注释 (JSON5风格)
            lines = []
            for line in text.split("\n"):
                if "//" in line:
                    idx = line.index("//")
                    # 判断是否在字符串中
                    in_str = False
                    clean_idx = idx
                    for i, ch in enumerate(line):
                        if ch == '"':
                            in_str = not in_str
                        if ch == '/' and not in_str and i > 0 and line[i-1] == '/':
                            clean_idx = i - 1
                            break
                    lines.append(line[:clean_idx])
                else:
                    lines.append(line)
            return json.loads("\n".join(lines))
    except Exception as e:
        print(f"加载配置文件失败: {e}")
        return {"baseApi": "http://localhost:3620"}


def main():
    # 配置主题
    ctk.set_appearance_mode("dark")
    ctk.set_default_color_theme("blue")

    config = load_config()
    base_url = config.get("baseApi")
    if not base_url:
        print("[错误] 配置文件 config.json5 中未设置 baseApi 地址")
        sys.exit(1)
    api = ApiClient(base_url)

    # 尝试自动登录
    api.load_token()
    if api.token:
        profile = api.get_profile()
        if profile and "id" in profile:
            app = MainWindow(api, profile)
            app.mainloop()
            return

    # 显示登录窗口
    auth = AuthWindow(api)
    auth.mainloop()


if __name__ == "__main__":
    main()
