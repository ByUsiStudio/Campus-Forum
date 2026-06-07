import customtkinter as ctk
from tkinter import messagebox


class MyArticlesView(ctk.CTkFrame):
    """我的文章 / 草稿箱"""

    def __init__(self, parent, api, user, main_window, mode="published", **kwargs):
        super().__init__(parent, fg_color="transparent")
        self.api = api
        self.user = user
        self.main_window = main_window
        self.mode = mode  # "published" or "drafts"
        self.current_page = 1
        self.total_pages = 1

        self.grid_columnconfigure(0, weight=1)
        self.grid_rowconfigure(1, weight=1)
        self._build_ui()

    def _build_ui(self):
        title = "📄 我的文章" if self.mode == "published" else "📝 草稿箱"
        ctk.CTkLabel(self, text=title, font=("Microsoft YaHei", 20, "bold"),
                      text_color=("#1a73e8", "#4fc3f7")).grid(
            row=0, column=0, sticky="w", padx=25, pady=(15, 10))

        self.scroll = ctk.CTkScrollableFrame(self, fg_color="transparent")
        self.scroll.grid(row=1, column=0, sticky="nsew", padx=25, pady=(0, 10))

        self._load()

    def _load(self):
        for w in self.scroll.winfo_children():
            w.destroy()

        if self.mode == "published":
            result = self.api.get_my_articles(page=self.current_page, page_size=12)
        else:
            result = self.api.get_my_drafts(page=self.current_page, page_size=12)

        if not result or "articles" not in result:
            ctk.CTkLabel(self.scroll, text="😅 暂无内容",
                          font=("Microsoft YaHei", 14)).pack(pady=50)
            return

        articles = result["articles"]
        self.total_pages = result.get("total_pages", 1)

        if not articles:
            ctk.CTkLabel(self.scroll, text="😅 暂无内容",
                          font=("Microsoft YaHei", 14)).pack(pady=50)
            return

        for article in articles:
            self._create_card(article)

    def _create_card(self, article):
        card = ctk.CTkFrame(self.scroll, corner_radius=10,
                             fg_color=("#ffffff", "#252540"))
        card.pack(fill="x", pady=4)

        inner = ctk.CTkFrame(card, fg_color="transparent")
        inner.pack(fill="x", padx=16, pady=12)

        ctk.CTkLabel(inner, text=article.get("title", "无标题"),
                      font=("Microsoft YaHei", 14, "bold"),
                      anchor="w").pack(fill="x")

        created = article.get("created_at", "")[:10]
        stats = f"👁 {article.get('view_count', 0)}  👍 {article.get('like_count', 0)}  💬 {article.get('comment_count', 0)}  🕐 {created}"
        ctk.CTkLabel(inner, text=stats, font=("Microsoft YaHei", 11),
                      text_color=("#888888", "#aaaaaa")).pack(anchor="w", pady=(4, 6))

        # 操作按钮
        btn_row = ctk.CTkFrame(inner, fg_color="transparent")
        btn_row.pack(fill="x")

        ctk.CTkButton(btn_row, text="查看", font=("Microsoft YaHei", 11),
                       height=30, fg_color="#1a73e8", corner_radius=6,
                       command=lambda a=article: self.main_window.show_article_detail(a["id"], a)
                       ).pack(side="left", padx=(0, 5))

        if self.mode == "drafts":
            ctk.CTkButton(btn_row, text="编辑", font=("Microsoft YaHei", 11),
                           height=30, fg_color=("#e8e8e8", "#333355"),
                           text_color=("#333333", "#cccccc"), corner_radius=6,
                           command=lambda a=article: self.main_window.show_create_with_data(a)
                           ).pack(side="left", padx=(0, 5))

            ctk.CTkButton(btn_row, text="发布", font=("Microsoft YaHei", 11),
                           height=30, fg_color=("#4caf50", "#2e7d32"),
                           corner_radius=6,
                           command=lambda aid=article["id"]: self._publish_draft(aid)
                           ).pack(side="left", padx=(0, 5))

        ctk.CTkButton(btn_row, text="删除", font=("Microsoft YaHei", 11),
                       height=30, fg_color=("#e53935", "#c62828"),
                       corner_radius=6,
                       command=lambda aid=article["id"]: self._delete(aid)
                       ).pack(side="right")

    def _publish_draft(self, article_id):
        result = self.api.publish_draft(article_id)
        if result and "message" in result:
            messagebox.showinfo("成功", "文章已发布")
            self._load()
        else:
            error = result.get("error", "发布失败") if result else "网络错误"
            messagebox.showerror("错误", error)

    def _delete(self, article_id):
        if messagebox.askyesno("确认删除", "确定要删除吗？"):
            result = self.api.delete_article(article_id)
            if result and "message" in result:
                messagebox.showinfo("成功", "已删除")
                self._load()
            else:
                error = result.get("error", "删除失败") if result else "网络错误"
                messagebox.showerror("错误", error)
