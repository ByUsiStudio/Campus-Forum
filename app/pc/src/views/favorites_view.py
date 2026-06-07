import customtkinter as ctk
from tkinter import messagebox


class FavoritesView(ctk.CTkFrame):
    """收藏夹"""

    def __init__(self, parent, api, user, main_window, **kwargs):
        super().__init__(parent, fg_color="transparent")
        self.api = api
        self.user = user
        self.main_window = main_window

        self.grid_columnconfigure(0, weight=1)
        self.grid_rowconfigure(1, weight=1)
        self._build_ui()

    def _build_ui(self):
        ctk.CTkLabel(self, text="⭐ 我的收藏", font=("Microsoft YaHei", 20, "bold"),
                      text_color=("#1a73e8", "#4fc3f7")).grid(
            row=0, column=0, sticky="w", padx=25, pady=(15, 10))

        self.scroll = ctk.CTkScrollableFrame(self, fg_color="transparent")
        self.scroll.grid(row=1, column=0, sticky="nsew", padx=25, pady=(0, 10))

        self._load()

    def _load(self):
        for w in self.scroll.winfo_children():
            w.destroy()

        result = self.api.get_favorites()
        if not result:
            ctk.CTkLabel(self.scroll, text="😅 加载失败",
                          font=("Microsoft YaHei", 14)).pack(pady=50)
            return

        # 支持不同字段名
        articles = result.get("articles", result.get("favorites", []))

        if not articles:
            ctk.CTkLabel(self.scroll, text="😅 还没有收藏任何文章",
                          font=("Microsoft YaHei", 14)).pack(pady=50)
            return

        for article in articles:
            card = ctk.CTkFrame(self.scroll, corner_radius=10,
                                 fg_color=("#ffffff", "#252540"))
            card.pack(fill="x", pady=4)

            inner = ctk.CTkFrame(card, fg_color="transparent")
            inner.pack(fill="x", padx=16, pady=12)

            ctk.CTkLabel(inner, text=article.get("title", "无标题"),
                          font=("Microsoft YaHei", 14, "bold"),
                          anchor="w").pack(fill="x")

            stats = f"👁 {article.get('view_count', 0)}  👍 {article.get('like_count', 0)}"
            ctk.CTkLabel(inner, text=stats, font=("Microsoft YaHei", 11),
                          text_color=("#888888", "#aaaaaa")).pack(anchor="w", pady=(4, 6))

            btn_row = ctk.CTkFrame(inner, fg_color="transparent")
            btn_row.pack(fill="x")

            ctk.CTkButton(btn_row, text="查看", font=("Microsoft YaHei", 11),
                           height=30, fg_color="#1a73e8", corner_radius=6,
                           command=lambda a=article: self.main_window.show_article_detail(
                               a["id"], a)).pack(side="left")
