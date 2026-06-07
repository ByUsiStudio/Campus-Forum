import customtkinter as ctk
from tkinter import messagebox
import re


class UserPageView(ctk.CTkScrollableFrame):
    """用户主页"""

    def __init__(self, parent, api, user, main_window, target_user_id: int, **kwargs):
        super().__init__(parent, fg_color="transparent")
        self.api = api
        self.user = user
        self.main_window = main_window
        self.target_user_id = target_user_id

        self.grid_columnconfigure(0, weight=1)
        self._load()

    def _load(self):
        result = self.api.get_user_info(self.target_user_id)
        if not result or "id" not in result:
            ctk.CTkLabel(self, text="😅 用户不存在",
                          font=("Microsoft YaHei", 16)).pack(pady=100)
            return

        target = result
        self._build_user_header(target)
        self._load_articles()

    def _build_user_header(self, target):
        """用户信息头部"""
        header = ctk.CTkFrame(self, corner_radius=12, fg_color=("#ffffff", "#252540"))
        header.pack(fill="x", padx=20, pady=(15, 10))

        inner = ctk.CTkFrame(header, fg_color="transparent")
        inner.pack(fill="x", padx=25, pady=20)
        inner.grid_columnconfigure(1, weight=1)

        # 头像
        avatar_text = target.get("display_name", "?")[0]
        ctk.CTkLabel(inner, text=avatar_text,
                      font=("Microsoft YaHei", 26, "bold"),
                      width=64, height=64, corner_radius=32,
                      fg_color=("#1a73e8", "#1565c0"),
                      text_color="white").grid(row=0, column=0, rowspan=3,
                                                padx=(0, 18))

        # 名称
        name_frame = ctk.CTkFrame(inner, fg_color="transparent")
        name_frame.grid(row=0, column=1, sticky="w")
        ctk.CTkLabel(name_frame, text=target.get("display_name", "未知"),
                      font=("Microsoft YaHei", 18, "bold")).pack(side="left")

        role = target.get("role", "")
        if role == "admin":
            ctk.CTkLabel(name_frame, text="管理员",
                          font=("Microsoft YaHei", 10),
                          fg_color="#ff6f00", text_color="white",
                          padx=6, pady=2, corner_radius=4).pack(side="left", padx=(8, 0))

        ctk.CTkLabel(inner, text=f"@{target.get('username', '')}",
                      font=("Microsoft YaHei", 12),
                      text_color=("#888888", "#aaaaaa")).grid(row=1, column=1, sticky="w")

        # 签名
        signature = target.get("signature", "")
        if signature:
            ctk.CTkLabel(inner, text=signature,
                          font=("Microsoft YaHei", 12),
                          text_color=("#555555", "#bbbbbb"),
                          anchor="w", wraplength=400).grid(row=2, column=1, sticky="w", pady=(4, 0))

        # 关注/粉丝统计
        stats_frame = ctk.CTkFrame(header, fg_color="transparent")
        stats_frame.pack(fill="x", padx=25, pady=(0, 15))

        for label, value in [("文章", target.get("article_count", 0)),
                               ("粉丝", target.get("follower_count", 0)),
                               ("关注", target.get("following_count", 0))]:
            item = ctk.CTkFrame(stats_frame, fg_color="transparent")
            item.pack(side="left", padx=(0, 25))
            ctk.CTkLabel(item, text=str(value),
                          font=("Microsoft YaHei", 14, "bold"),
                          text_color=("#1a73e8", "#4fc3f7")).pack()
            ctk.CTkLabel(item, text=label, font=("Microsoft YaHei", 10),
                          text_color=("#888888", "#aaaaaa")).pack()

        # 关注/取消关注按钮
        if self.target_user_id != self.user.get("id"):
            # 检查关注状态
            status_result = self.api.get_follow_status(self.target_user_id)
            is_following = status_result and status_result.get("is_following", False) if status_result else False

            self.follow_btn = ctk.CTkButton(stats_frame, text="已关注" if is_following else "+ 关注",
                                             font=("Microsoft YaHei", 12),
                                             height=34, corner_radius=17,
                                             fg_color=("#1a73e8", "#1565c0") if not is_following else ("#e0e0e0", "#444444"),
                                             text_color="white" if not is_following else ("#333333", "#cccccc"),
                                             command=lambda: self._toggle_follow(is_following))
            self.follow_btn.pack(side="right")

    def _toggle_follow(self, currently_following):
        if currently_following:
            result = self.api.unfollow_user(self.target_user_id)
            if result:
                self.follow_btn.configure(text="+ 关注",
                                           fg_color=("#1a73e8", "#1565c0"),
                                           text_color="white")
                # 刷新
                for w in self.winfo_children():
                    w.destroy()
                self._load()
        else:
            result = self.api.follow_user(self.target_user_id)
            if result:
                self.follow_btn.configure(text="已关注",
                                           fg_color=("#e0e0e0", "#444444"),
                                           text_color=("#333333", "#cccccc"))
                for w in self.winfo_children():
                    w.destroy()
                self._load()

    def _load_articles(self):
        """加载用户的文章"""
        ctk.CTkLabel(self, text="📄 他的文章",
                      font=("Microsoft YaHei", 16, "bold"),
                      text_color=("#1a73e8", "#4fc3f7")).pack(
            anchor="w", padx=20, pady=(10, 10))

        result = self.api.get_user_articles(self.target_user_id)
        if not result or "articles" not in result:
            ctk.CTkLabel(self, text="暂无文章", font=("Microsoft YaHei", 12),
                          text_color=("#888888", "#aaaaaa")).pack(pady=20)
            return

        articles = result["articles"]
        if not articles:
            ctk.CTkLabel(self, text="暂无文章", font=("Microsoft YaHei", 12),
                          text_color=("#888888", "#aaaaaa")).pack(pady=20)
            return

        for article in articles:
            card = ctk.CTkFrame(self, corner_radius=10,
                                 fg_color=("#ffffff", "#252540"))
            card.pack(fill="x", pady=4, padx=20)

            inner = ctk.CTkFrame(card, fg_color="transparent")
            inner.pack(fill="x", padx=16, pady=12)

            title = article.get("title", "无标题")
            if article.get("is_anonymous"):
                title = "🔒 " + title

            ctk.CTkLabel(inner, text=title,
                          font=("Microsoft YaHei", 14, "bold"),
                          anchor="w").pack(fill="x")

            content = article.get("content", "")
            preview = re.sub(r'[#*`>\-\[\]()!]', '', content)[:100]
            if len(content) > 100:
                preview += "..."

            ctk.CTkLabel(inner, text=preview,
                          font=("Microsoft YaHei", 11),
                          text_color=("#666666", "#aaaaaa"),
                          anchor="w", wraplength=600).pack(anchor="w", pady=(4, 6))

            btn_row = ctk.CTkFrame(inner, fg_color="transparent")
            btn_row.pack(fill="x")
            ctk.CTkButton(btn_row, text="查看", font=("Microsoft YaHei", 11),
                           height=30, fg_color="#1a73e8", corner_radius=6,
                           command=lambda a=article: self.main_window.show_article_detail(
                               a["id"], a)).pack(side="left")
