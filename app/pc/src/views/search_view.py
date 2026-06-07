import customtkinter as ctk
import re


class SearchView(ctk.CTkFrame):
    """搜索结果页面"""

    def __init__(self, parent, api, user, main_window, keyword: str = "", **kwargs):
        super().__init__(parent, fg_color="transparent")
        self.api = api
        self.user = user
        self.main_window = main_window
        self.keyword = keyword
        self.current_page = 1
        self.total_pages = 1

        self.grid_columnconfigure(0, weight=1)
        self.grid_rowconfigure(2, weight=1)
        self._build_ui()

    def _build_ui(self):
        # 标题
        ctk.CTkLabel(self, text=f"🔍 搜索结果: \"{self.keyword}\"",
                      font=("Microsoft YaHei", 18, "bold"),
                      text_color=("#1a73e8", "#4fc3f7")).grid(
            row=0, column=0, sticky="w", padx=25, pady=(15, 5))

        self.result_frame = ctk.CTkScrollableFrame(self, fg_color="transparent")
        self.result_frame.grid(row=2, column=0, sticky="nsew", padx=25, pady=(5, 10))

        self._search()

    def _search(self):
        for w in self.result_frame.winfo_children():
            w.destroy()

        result = self.api.search_articles(self.keyword, page=self.current_page)
        if not result or "articles" not in result:
            ctk.CTkLabel(self.result_frame, text="😅 搜索失败",
                          font=("Microsoft YaHei", 14)).pack(pady=50)
            return

        articles = result["articles"]
        self.total_pages = result.get("total_pages", 1)

        # 结果统计
        total = result.get("total", len(articles))
        ctk.CTkLabel(self.result_frame, text=f"共找到 {total} 条结果",
                      font=("Microsoft YaHei", 12),
                      text_color=("#888888", "#aaaaaa")).pack(anchor="w", pady=(0, 10))

        if not articles:
            ctk.CTkLabel(self.result_frame, text="😅 未找到相关文章",
                          font=("Microsoft YaHei", 14)).pack(pady=50)
            return

        for article in articles:
            self._create_card(article)

        # 分页
        if self.total_pages > 1:
            page_frame = ctk.CTkFrame(self.result_frame, fg_color="transparent")
            page_frame.pack(pady=10)
            ctk.CTkButton(page_frame, text="上一页", font=("Microsoft YaHei", 12),
                           height=32, command=self._prev,
                           state="normal" if self.current_page > 1 else "disabled"
                           ).pack(side="left", padx=5)
            ctk.CTkLabel(page_frame, text=f"{self.current_page}/{self.total_pages}",
                          font=("Microsoft YaHei", 12)).pack(side="left", padx=10)
            ctk.CTkButton(page_frame, text="下一页", font=("Microsoft YaHei", 12),
                           height=32, command=self._next,
                           state="normal" if self.current_page < self.total_pages else "disabled"
                           ).pack(side="left", padx=5)

    def _create_card(self, article):
        card = ctk.CTkFrame(self.result_frame, corner_radius=10,
                             fg_color=("#ffffff", "#252540"))
        card.pack(fill="x", pady=4)

        inner = ctk.CTkFrame(card, fg_color="transparent")
        inner.pack(fill="x", padx=16, pady=12)

        # 高亮关键词
        title = article.get("title", "无标题")
        ctk.CTkLabel(inner, text=title,
                      font=("Microsoft YaHei", 14, "bold"),
                      anchor="w").pack(fill="x")

        # 内容片段
        content = article.get("content", "")
        preview = re.sub(r'[#*`>\-\[\]()!]', '', content)[:150]
        if len(content) > 150:
            preview += "..."
        ctk.CTkLabel(inner, text=preview,
                      font=("Microsoft YaHei", 11),
                      text_color=("#666666", "#aaaaaa"),
                      anchor="w", wraplength=600).pack(anchor="w", pady=(4, 6))

        stats = f"👁 {article.get('view_count', 0)}  👍 {article.get('like_count', 0)}  💬 {article.get('comment_count', 0)}"
        ctk.CTkLabel(inner, text=stats, font=("Microsoft YaHei", 11),
                      text_color=("#888888", "#aaaaaa")).pack(anchor="w")

        # 查看按钮
        ctk.CTkButton(inner, text="查看详情", font=("Microsoft YaHei", 11),
                       height=30, fg_color="#1a73e8", corner_radius=6,
                       command=lambda a=article: self.main_window.show_article_detail(
                           a["id"], a)).pack(anchor="e", pady=(6, 0))

    def _prev(self):
        if self.current_page > 1:
            self.current_page -= 1
            self._search()

    def _next(self):
        if self.current_page < self.total_pages:
            self.current_page += 1
            self._search()
