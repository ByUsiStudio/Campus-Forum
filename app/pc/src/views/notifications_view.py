import customtkinter as ctk
from tkinter import messagebox


class NotificationsView(ctk.CTkFrame):
    """通知列表"""

    def __init__(self, parent, api, user, main_window, **kwargs):
        super().__init__(parent, fg_color="transparent")
        self.api = api
        self.user = user
        self.main_window = main_window

        self.grid_columnconfigure(0, weight=1)
        self.grid_rowconfigure(1, weight=1)
        self._build_ui()

    def _build_ui(self):
        # 标题 + 操作
        header = ctk.CTkFrame(self, fg_color="transparent")
        header.grid(row=0, column=0, sticky="ew", padx=25, pady=(15, 10))
        header.grid_columnconfigure(0, weight=1)

        ctk.CTkLabel(header, text="🔔 通知", font=("Microsoft YaHei", 20, "bold"),
                      text_color=("#1a73e8", "#4fc3f7")).grid(row=0, column=0, sticky="w")

        ctk.CTkButton(header, text="全部已读", font=("Microsoft YaHei", 12),
                       height=32, fg_color="transparent",
                       text_color=("#1a73e8", "#4fc3f7"),
                       command=self._read_all).grid(row=0, column=1, sticky="e")

        self.scroll = ctk.CTkScrollableFrame(self, fg_color="transparent")
        self.scroll.grid(row=1, column=0, sticky="nsew", padx=25, pady=(0, 10))

        self._load()

    def _load(self):
        for w in self.scroll.winfo_children():
            w.destroy()

        result = self.api.get_notifications()
        if not result:
            ctk.CTkLabel(self.scroll, text="😅 加载失败",
                          font=("Microsoft YaHei", 14)).pack(pady=50)
            return

        # 尝试不同的返回结构
        notifications = result.get("notifications", result.get("data", []))

        if not notifications:
            ctk.CTkLabel(self.scroll, text="✅ 暂无通知",
                          font=("Microsoft YaHei", 14)).pack(pady=50)
            return

        for notif in notifications:
            self._create_notif_card(notif)

    def _create_notif_card(self, notif):
        card = ctk.CTkFrame(self.scroll, corner_radius=10,
                             fg_color=("#ffffff", "#252540"))
        card.pack(fill="x", pady=3)

        is_read = notif.get("is_read", True)
        inner = ctk.CTkFrame(card, fg_color="transparent")
        inner.pack(fill="x", padx=16, pady=12)

        # 标题行
        title_frame = ctk.CTkFrame(inner, fg_color="transparent")
        title_frame.pack(fill="x")
        title_frame.grid_columnconfigure(0, weight=1)

        ntype = notif.get("type", "system")
        icon = "🔔" if not is_read else "🔕"

        ctk.CTkLabel(title_frame, text=f"{icon} {notif.get('title', '通知')}",
                      font=("Microsoft YaHei", 13, "bold") if not is_read else ("Microsoft YaHei", 13),
                      anchor="w").grid(row=0, column=0, sticky="w")

        time_str = notif.get("created_at", "")[:10]
        ctk.CTkLabel(title_frame, text=time_str,
                      font=("Microsoft YaHei", 10),
                      text_color=("#888888", "#aaaaaa")).grid(row=0, column=1, sticky="e")

        ctk.CTkLabel(inner, text=notif.get("content", ""),
                      font=("Microsoft YaHei", 11),
                      text_color=("#666666", "#aaaaaa"),
                      anchor="w", wraplength=500).pack(anchor="w", pady=(4, 6))

        # 操作按钮
        if not is_read:
            btn_row = ctk.CTkFrame(inner, fg_color="transparent")
            btn_row.pack(fill="x")
            ctk.CTkButton(btn_row, text="标记已读", font=("Microsoft YaHei", 10),
                           height=26, fg_color="transparent",
                           text_color=("#1a73e8", "#4fc3f7"),
                           command=lambda nid=notif["id"]: self._mark_read(nid)
                           ).pack(side="left")

    def _mark_read(self, nid):
        result = self.api.read_notification(nid)
        if result:
            self._load()

    def _read_all(self):
        result = self.api.read_all_notifications()
        if result:
            messagebox.showinfo("成功", "所有通知已标记为已读")
            self._load()
        else:
            messagebox.showerror("错误", "操作失败")
