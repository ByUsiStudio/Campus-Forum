import customtkinter as ctk
from tkinter import messagebox


class ProfileView(ctk.CTkScrollableFrame):
    """个人资料页面"""

    def __init__(self, parent, api, user, main_window, **kwargs):
        super().__init__(parent, fg_color="transparent")
        self.api = api
        self.user = user
        self.main_window = main_window

        self.grid_columnconfigure(0, weight=1)
        self._load_and_build()

    def _load_and_build(self):
        # 获取最新资料
        profile = self.api.get_profile()
        if profile and "id" in profile:
            self.user.update(profile)
        self._build_ui()

    def _build_ui(self):
        # 标题
        ctk.CTkLabel(self, text="👤 个人资料", font=("Microsoft YaHei", 22, "bold"),
                      text_color=("#1a73e8", "#4fc3f7")).pack(
            anchor="w", padx=30, pady=(20, 25))

        # 卡片
        card = ctk.CTkFrame(self, corner_radius=12, fg_color=("#ffffff", "#252540"))
        card.pack(fill="x", padx=30, pady=(0, 15))
        inner = ctk.CTkFrame(card, fg_color="transparent")
        inner.pack(fill="x", padx=25, pady=20)
        inner.grid_columnconfigure(1, weight=1)

        # 头像
        avatar_text = self.user.get("display_name", "?")[0]
        ctk.CTkLabel(inner, text=avatar_text,
                      font=("Microsoft YaHei", 28, "bold"),
                      width=70, height=70, corner_radius=35,
                      fg_color=("#1a73e8", "#1565c0"),
                      text_color="white").grid(row=0, column=0, rowspan=4,
                                                padx=(0, 20), pady=(0, 0))

        # 信息字段
        fields = [
            ("用户名", self.user.get("username", "")),
            ("昵称", self.user.get("display_name", "")),
            ("QQ号", self.user.get("qq_number", "")),
            ("角色", "管理员" if self.user.get("role") == "admin" else "用户"),
            ("个性签名", self.user.get("signature", "未设置")),
            ("注册时间", str(self.user.get("created_at", ""))[:10]),
        ]
        for i, (label, value) in enumerate(fields):
            ctk.CTkLabel(inner, text=f"{label}:", font=("Microsoft YaHei", 12, "bold"),
                          text_color=("#666666", "#aaaaaa")).grid(
                row=i, column=0, sticky="w", pady=4)
            ctk.CTkLabel(inner, text=str(value) if value else "-",
                          font=("Microsoft YaHei", 12),
                          anchor="w", wraplength=300).grid(
                row=i, column=1, sticky="w", pady=4, padx=(10, 0))

        # ── 编辑资料 ──
        edit_card = ctk.CTkFrame(self, corner_radius=12, fg_color=("#ffffff", "#252540"))
        edit_card.pack(fill="x", padx=30, pady=(0, 20))
        edit_inner = ctk.CTkFrame(edit_card, fg_color="transparent")
        edit_inner.pack(fill="x", padx=25, pady=20)
        edit_inner.grid_columnconfigure(1, weight=1)

        ctk.CTkLabel(edit_inner, text="编辑资料", font=("Microsoft YaHei", 16, "bold")).grid(
            row=0, column=0, columnspan=2, sticky="w", pady=(0, 15))

        ctk.CTkLabel(edit_inner, text="昵称", font=("Microsoft YaHei", 12),
                      anchor="w").grid(row=1, column=0, sticky="w", pady=5)
        self.name_entry = ctk.CTkEntry(edit_inner, height=36,
                                        font=("Microsoft YaHei", 13))
        self.name_entry.insert(0, self.user.get("display_name", ""))
        self.name_entry.grid(row=1, column=1, sticky="ew", padx=(10, 0), pady=5)

        ctk.CTkLabel(edit_inner, text="签名", font=("Microsoft YaHei", 12),
                      anchor="w").grid(row=2, column=0, sticky="w", pady=5)
        self.sig_entry = ctk.CTkEntry(edit_inner, height=36,
                                       font=("Microsoft YaHei", 13),
                                       placeholder_text="设置个性签名")
        self.sig_entry.insert(0, self.user.get("signature", ""))
        self.sig_entry.grid(row=2, column=1, sticky="ew", padx=(10, 0), pady=5)

        ctk.CTkButton(edit_inner, text="保存修改", font=("Microsoft YaHei", 13, "bold"),
                       height=36, corner_radius=6,
                       command=self._do_update).grid(
            row=3, column=1, sticky="e", padx=(10, 0), pady=(10, 0))

        # ── 管理入口 ──
        if self.user.get("role") == "admin":
            admin_card = ctk.CTkFrame(self, corner_radius=12, fg_color=("#fff3e0", "#3a2a10"))
            admin_card.pack(fill="x", padx=30, pady=(0, 20))
            admin_inner = ctk.CTkFrame(admin_card, fg_color="transparent")
            admin_inner.pack(fill="x", padx=25, pady=15)

            ctk.CTkLabel(admin_inner, text="⚙️ 管理面板",
                          font=("Microsoft YaHei", 16, "bold"),
                          text_color=("#e65100", "#ffb74d")).pack(anchor="w")
            ctk.CTkLabel(admin_inner, text="你有管理员权限，可以管理论坛内容",
                          font=("Microsoft YaHei", 11),
                          text_color=("#888888", "#aaaaaa")).pack(anchor="w", pady=(4, 10))

            # 这里可以添加更多管理功能
            ctk.CTkButton(admin_inner, text="查看用户列表",
                           font=("Microsoft YaHei", 12),
                           height=34, corner_radius=6,
                           fg_color=("#e65100", "#bf360c"),
                           command=self._show_admin_panel).pack(side="left", padx=(0, 8))

    def _do_update(self):
        name = self.name_entry.get().strip()
        sig = self.sig_entry.get().strip()
        if not name:
            messagebox.showerror("错误", "昵称不能为空")
            return
        result = self.api.update_profile(display_name=name, signature=sig)
        if result and "message" in result:
            messagebox.showinfo("成功", "资料已更新")
            # 刷新
            for w in self.winfo_children():
                w.destroy()
            self._load_and_build()
        else:
            error = result.get("error", "更新失败") if result else "网络错误"
            messagebox.showerror("错误", error)

    def _show_admin_panel(self):
        """管理员面板入口"""
        messagebox.showinfo("管理面板", "管理功能正在开发中...\n请通过网页后台管理系统进行操作。")
