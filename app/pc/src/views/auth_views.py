import customtkinter as ctk
from tkinter import messagebox
from ..api_client import ApiClient
from .home_view import MainWindow


class AuthWindow(ctk.CTk):
    """登录/注册/密码重置窗口"""

    def __init__(self, api: ApiClient):
        super().__init__()
        self.api = api
        self.title("校园论坛 - 登录")
        self.geometry("400x560")
        self.resizable(False, False)
        self.update_idletasks()
        w, h = 400, 560
        sw = self.winfo_screenwidth()
        sh = self.winfo_screenheight()
        self.geometry(f"{w}x{h}+{(sw-w)//2}+{(sh-h)//2}")

        self._build_login_view()

    def _clear(self):
        for w in self.winfo_children():
            w.destroy()

    def _show_error(self, msg: str):
        messagebox.showerror("错误", msg)

    def _show_info(self, msg: str):
        messagebox.showinfo("提示", msg)

    # ── 登录界面 ──────────────────────────────────────────────

    def _build_login_view(self):
        self._clear()
        self.title("校园论坛 - 登录")
        self.geometry("400x560")

        main = ctk.CTkFrame(self, fg_color="transparent")
        main.pack(fill="both", expand=True, padx=40, pady=40)

        # 标题
        ctk.CTkLabel(main, text="校园论坛", font=("Microsoft YaHei", 28, "bold"),
                      text_color=("#1a73e8", "#4fc3f7")).pack(pady=(30, 5))
        ctk.CTkLabel(main, text="欢迎回来", font=("Microsoft YaHei", 13),
                      text_color=("#666666", "#aaaaaa")).pack(pady=(0, 25))

        ctk.CTkLabel(main, text="用户名", font=("Microsoft YaHei", 12),
                      anchor="w").pack(fill="x", padx=5)
        self.login_user = ctk.CTkEntry(main, placeholder_text="请输入用户名", height=38,
                                        font=("Microsoft YaHei", 13))
        self.login_user.pack(fill="x", pady=(3, 12))

        ctk.CTkLabel(main, text="密码", font=("Microsoft YaHei", 12),
                      anchor="w").pack(fill="x", padx=5)
        self.login_pwd = ctk.CTkEntry(main, placeholder_text="请输入密码", height=38,
                                       show="●", font=("Microsoft YaHei", 13))
        self.login_pwd.pack(fill="x", pady=(3, 5))

        forgot_frame = ctk.CTkFrame(main, fg_color="transparent")
        forgot_frame.pack(fill="x", pady=(0, 20))
        ctk.CTkButton(forgot_frame, text="忘记密码?", font=("Microsoft YaHei", 11),
                       fg_color="transparent", hover_color=("#e0e0e0", "#333333"),
                       text_color=("#1a73e8", "#4fc3f7"), width=0, height=20,
                       command=self._build_reset_view).pack(side="right")

        self.login_btn = ctk.CTkButton(main, text="登 录", height=42,
                                        font=("Microsoft YaHei", 15, "bold"),
                                        corner_radius=8, command=self._do_login)
        self.login_btn.pack(fill="x", pady=(0, 12))

        ctk.CTkButton(main, text="还没有账号？立即注册", font=("Microsoft YaHei", 12),
                       fg_color="transparent", hover_color=("#e0e0e0", "#333333"),
                       text_color=("#1a73e8", "#4fc3f7"), width=0, height=20,
                       command=self._build_register_view).pack()

    def _do_login(self):
        username = self.login_user.get().strip()
        password = self.login_pwd.get()
        if not username or not password:
            self._show_error("请输入用户名和密码")
            return
        self.login_btn.configure(state="disabled", text="登录中...")
        self.update()
        result = self.api.login(username, password)
        if result and "token" in result:
            profile = result.get("user", self.api.get_profile())
            if profile and "id" in profile:
                self.destroy()
                app = MainWindow(self.api, profile)
                app.mainloop()
                return
            self._show_error("登录成功但获取用户信息失败")
        else:
            error = result.get("error", "登录失败") if result else "网络错误"
            self._show_error(error)
        self.login_btn.configure(state="normal", text="登 录")

    # ── 注册界面 ──────────────────────────────────────────────

    def _build_register_view(self):
        self._clear()
        self.title("校园论坛 - 注册")
        self.geometry("400x640")

        main = ctk.CTkFrame(self, fg_color="transparent")
        main.pack(fill="both", expand=True, padx=40, pady=30)

        ctk.CTkLabel(main, text="校园论坛", font=("Microsoft YaHei", 26, "bold"),
                      text_color=("#1a73e8", "#4fc3f7")).pack(pady=(20, 5))
        ctk.CTkLabel(main, text="创建新账号", font=("Microsoft YaHei", 14),
                      text_color=("#666666", "#aaaaaa")).pack(pady=(0, 20))

        ctk.CTkLabel(main, text="用户名", font=("Microsoft YaHei", 12),
                      anchor="w").pack(fill="x", padx=5)
        self.reg_user = ctk.CTkEntry(main, placeholder_text="4-20个字符", height=36)
        self.reg_user.pack(fill="x", pady=(3, 10))

        ctk.CTkLabel(main, text="昵称", font=("Microsoft YaHei", 12),
                      anchor="w").pack(fill="x", padx=5)
        self.reg_name = ctk.CTkEntry(main, placeholder_text="你的昵称", height=36)
        self.reg_name.pack(fill="x", pady=(3, 10))

        ctk.CTkLabel(main, text="QQ号", font=("Microsoft YaHei", 12),
                      anchor="w").pack(fill="x", padx=5)
        self.reg_qq = ctk.CTkEntry(main, placeholder_text="用于找回密码", height=36)
        self.reg_qq.pack(fill="x", pady=(3, 10))

        ctk.CTkLabel(main, text="密码", font=("Microsoft YaHei", 12),
                      anchor="w").pack(fill="x", padx=5)
        self.reg_pwd = ctk.CTkEntry(main, placeholder_text="至少6位", height=36, show="●")
        self.reg_pwd.pack(fill="x", pady=(3, 5))

        self.reg_btn = ctk.CTkButton(main, text="注 册", height=42,
                                      font=("Microsoft YaHei", 15, "bold"),
                                      corner_radius=8, command=self._do_register)
        self.reg_btn.pack(fill="x", pady=(15, 10))

        ctk.CTkButton(main, text="已有账号？返回登录", font=("Microsoft YaHei", 12),
                       fg_color="transparent", hover_color=("#e0e0e0", "#333333"),
                       text_color=("#1a73e8", "#4fc3f7"), width=0, height=20,
                       command=self._build_login_view).pack()

    def _do_register(self):
        username = self.reg_user.get().strip()
        display_name = self.reg_name.get().strip()
        qq = self.reg_qq.get().strip()
        password = self.reg_pwd.get()
        if not all([username, display_name, qq, password]):
            self._show_error("请填写所有字段")
            return
        if len(password) < 6:
            self._show_error("密码至少6位")
            return
        self.reg_btn.configure(state="disabled", text="注册中...")
        self.update()
        result = self.api.register(username, qq, display_name, password)
        if result and "message" in result:
            self._show_info("注册成功！请登录")
            self._build_login_view()
        else:
            error = result.get("error", "注册失败") if result else "网络错误"
            self._show_error(error)
        self.reg_btn.configure(state="normal", text="注 册")

    # ── 密码重置界面 ──────────────────────────────────────────

    def _build_reset_view(self):
        self._clear()
        self.title("校园论坛 - 重置密码")
        self.geometry("400x560")

        self.reset_identifier = None

        main = ctk.CTkFrame(self, fg_color="transparent")
        main.pack(fill="both", expand=True, padx=40, pady=30)

        ctk.CTkLabel(main, text="校园论坛", font=("Microsoft YaHei", 26, "bold"),
                      text_color=("#1a73e8", "#4fc3f7")).pack(pady=(20, 5))
        ctk.CTkLabel(main, text="重置密码", font=("Microsoft YaHei", 14),
                      text_color=("#666666", "#aaaaaa")).pack(pady=(0, 20))

        ctk.CTkLabel(main, text="QQ号", font=("Microsoft YaHei", 12),
                      anchor="w").pack(fill="x", padx=5)
        self.reset_qq = ctk.CTkEntry(main, placeholder_text="注册时绑定的QQ号", height=36)
        self.reset_qq.pack(fill="x", pady=(3, 10))

        self.send_code_btn = ctk.CTkButton(main, text="发送验证码", height=36,
                                            font=("Microsoft YaHei", 13),
                                            command=self._do_send_code)
        self.send_code_btn.pack(fill="x", pady=(0, 15))

        ctk.CTkLabel(main, text="验证码", font=("Microsoft YaHei", 12),
                      anchor="w").pack(fill="x", padx=5)
        self.reset_code = ctk.CTkEntry(main, placeholder_text="邮箱中收到的验证码", height=36)
        self.reset_code.pack(fill="x", pady=(3, 10))

        ctk.CTkLabel(main, text="新密码", font=("Microsoft YaHei", 12),
                      anchor="w").pack(fill="x", padx=5)
        self.reset_pwd = ctk.CTkEntry(main, placeholder_text="至少6位", height=36, show="●")
        self.reset_pwd.pack(fill="x", pady=(3, 5))

        self.reset_btn = ctk.CTkButton(main, text="重置密码", height=42,
                                        font=("Microsoft YaHei", 15, "bold"),
                                        corner_radius=8, command=self._do_reset)
        self.reset_btn.pack(fill="x", pady=(15, 10))

        ctk.CTkButton(main, text="返回登录", font=("Microsoft YaHei", 12),
                       fg_color="transparent", hover_color=("#e0e0e0", "#333333"),
                       text_color=("#1a73e8", "#4fc3f7"), width=0, height=20,
                       command=self._build_login_view).pack()

    def _do_send_code(self):
        qq = self.reset_qq.get().strip()
        if not qq:
            self._show_error("请输入QQ号")
            return
        self.send_code_btn.configure(state="disabled", text="发送中...")
        self.update()
        result = self.api.send_reset_code(qq)
        if result and "identifier" in result:
            self.reset_identifier = result["identifier"]
            self._show_info("验证码已发送到您的QQ邮箱")
            self._countdown(60)
        else:
            error = result.get("error", "发送失败") if result else "网络错误"
            self._show_error(error)
            self.send_code_btn.configure(state="normal", text="发送验证码")

    def _countdown(self, remaining):
        if remaining <= 0:
            self.send_code_btn.configure(state="normal", text="发送验证码")
            return
        self.send_code_btn.configure(text=f"重新发送({remaining}s)")
        self.after(1000, lambda: self._countdown(remaining - 1))

    def _do_reset(self):
        qq = self.reset_qq.get().strip()
        code = self.reset_code.get().strip()
        password = self.reset_pwd.get()
        if not all([qq, code, password]):
            self._show_error("请填写所有字段")
            return
        if len(password) < 6:
            self._show_error("密码至少6位")
            return
        if not self.reset_identifier:
            self._show_error("请先发送验证码")
            return
        self.reset_btn.configure(state="disabled", text="重置中...")
        self.update()
        result = self.api.reset_password(qq, code, self.reset_identifier, password)
        if result and "message" in result:
            self._show_info("密码重置成功！请登录")
            self._build_login_view()
        else:
            error = result.get("error", "重置失败") if result else "网络错误"
            self._show_error(error)
        self.reset_btn.configure(state="normal", text="重置密码")
