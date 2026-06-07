import customtkinter as ctk
from tkinter import messagebox


class HomeContent(ctk.CTkFrame):
    """首页 - 文章列表"""

    def __init__(self, parent, api, user, main_window, **kwargs):
        super().__init__(parent, fg_color="transparent")
        self.api = api
        self.user = user
        self.main_window = main_window
        self.current_page = 1
        self.total_pages = 1
        self.current_category = None
        self.articles = []

        self._build_ui()
        self._load_articles()

    def _build_ui(self):
        self.grid_columnconfigure(0, weight=1)
        self.grid_rowconfigure(1, weight=1)

        # ── 顶部栏 ──
        top = ctk.CTkFrame(self, fg_color="transparent")
        top.grid(row=0, column=0, sticky="ew", padx=20, pady=(15, 5))
        top.grid_columnconfigure(1, weight=1)

        ctk.CTkLabel(top, text="📋 最新文章", font=("Microsoft YaHei", 20, "bold"),
                      text_color=("#1a73e8", "#4fc3f7")).grid(row=0, column=0, sticky="w")

        # 搜索框
        search_frame = ctk.CTkFrame(top, fg_color=("#e8e8e8", "#2a2a3e"), corner_radius=20, height=36)
        search_frame.grid(row=0, column=1, padx=(20, 5), sticky="e")
        search_frame.grid_propagate(False)
        self.search_entry = ctk.CTkEntry(search_frame, placeholder_text="搜索文章...",
                                          font=("Microsoft YaHei", 12),
                                          fg_color="transparent", border_width=0, height=36, width=220)
        self.search_entry.pack(side="left", padx=(12, 5), pady=2)
        self.search_entry.bind("<Return>", lambda e: self._do_search())
        ctk.CTkButton(search_frame, text="🔍", width=36, height=30, corner_radius=15,
                       fg_color="transparent", hover_color=("#d0d0d0", "#3a3a4e"),
                       font=("Microsoft YaHei", 14), command=self._do_search).pack(side="right", padx=3)

        # ── 分类标签 ──
        cat_frame = ctk.CTkFrame(self, fg_color="transparent")
        cat_frame.grid(row=1, column=0, sticky="ew", padx=20, pady=(5, 5))
        cat_frame.grid_columnconfigure(0, weight=1)

        self.cat_scroll = ctk.CTkScrollableFrame(cat_frame, orientation="horizontal",
                                                  height=38, fg_color="transparent")
        self.cat_scroll.pack(fill="x")

        # 全部按钮
        self.cat_btns = []
        all_btn = ctk.CTkButton(self.cat_scroll, text="全部", font=("Microsoft YaHei", 12),
                                 height=32, corner_radius=16,
                                 fg_color=("#1a73e8", "#1565c0"), text_color="white",
                                 hover_color=("#1557b0", "#0d47a1"),
                                 command=lambda: self._filter_by_category(None))
        all_btn.pack(side="left", padx=2)
        self.cat_btns.append(("全部按钮", all_btn, None))

        # 加载分类
        categories = getattr(self.main_window, "categories", [])
        for cat in categories:
            btn = ctk.CTkButton(self.cat_scroll, text=cat.get("name", "未知"),
                                font=("Microsoft YaHei", 12),
                                height=32, corner_radius=16,
                                fg_color=("#e0e0e0", "#333355"),
                                text_color=("#333333", "#cccccc"),
                                hover_color=("#c8c8c8", "#444466"),
                                command=lambda cid=cat.get("id"): self._filter_by_category(cid))
            btn.pack(side="left", padx=2)
            self.cat_btns.append((cat.get("name", ""), btn, cat.get("id")))

        # ── 文章列表（可滚动） ──
        self.scroll_frame = ctk.CTkScrollableFrame(self, fg_color="transparent")
        self.scroll_frame.grid(row=2, column=0, sticky="nsew", padx=20, pady=(5, 10))
        self.grid_rowconfigure(2, weight=1)

        # ── 底部分页 ──
        self.paging_frame = ctk.CTkFrame(self, fg_color="transparent")
        self.paging_frame.grid(row=3, column=0, sticky="ew", padx=20, pady=(5, 15))

        self.prev_btn = ctk.CTkButton(self.paging_frame, text="◀ 上一页", font=("Microsoft YaHei", 12),
                                       height=32, command=self._prev_page)
        self.prev_btn.pack(side="left", padx=5)

        self.page_label = ctk.CTkLabel(self.paging_frame, text="第 1 页 / 共 1 页",
                                        font=("Microsoft YaHei", 12))
        self.page_label.pack(side="left", padx=15)

        self.next_btn = ctk.CTkButton(self.paging_frame, text="下一页 ▶", font=("Microsoft YaHei", 12),
                                       height=32, command=self._next_page)
        self.next_btn.pack(side="left", padx=5)

        # 加载中
        self.loading_label = ctk.CTkLabel(self.scroll_frame, text="正在加载...",
                                           font=("Microsoft YaHei", 14))

    def _highlight_category(self, active_id):
        for name, btn, cid in self.cat_btns:
            if cid == active_id:
                btn.configure(fg_color=("#1a73e8", "#1565c0"), text_color="white")
            else:
                btn.configure(fg_color=("#e0e0e0", "#333355"), text_color=("#333333", "#cccccc"))

    def _filter_by_category(self, category_id):
        self.current_category = category_id
        self.current_page = 1
        self._highlight_category(category_id)
        self._load_articles()

    def _do_search(self):
        keyword = self.search_entry.get().strip()
        if not keyword:
            return
        # 跳转到搜索
        from .search_view import SearchView
        if hasattr(self.main_window, '_switch_content'):
            self.main_window._switch_content(SearchView, keyword=keyword)
        else:
            self.main_window.show_search(keyword)

    def _load_articles(self):
        # 清空列表
        self.loading_label.pack(pady=50)
        for w in self.scroll_frame.winfo_children():
            if w != self.loading_label:
                w.destroy()
        self.update()

        result = self.api.get_articles(page=self.current_page, page_size=12,
                                        category_id=self.current_category)
        self.loading_label.pack_forget()

        if result and "articles" in result:
            self.articles = result["articles"]
            self.total_pages = result.get("total_pages", 1)
            self._render_articles()
            self._update_pagination()
        else:
            error = result.get("error", "加载失败") if result else "网络错误"
            ctk.CTkLabel(self.scroll_frame, text=f"😅 {error}",
                          font=("Microsoft YaHei", 14)).pack(pady=50)

    def _render_articles(self):
        if not self.articles:
            ctk.CTkLabel(self.scroll_frame, text="😅 暂无文章",
                          font=("Microsoft YaHei", 14)).pack(pady=50)
            return

        for article in self.articles:
            self._create_article_card(article)

    def _create_article_card(self, article):
        card = ctk.CTkFrame(self.scroll_frame, corner_radius=12,
                             fg_color=("#ffffff", "#252540"),
                             border_width=0)
        card.pack(fill="x", pady=6, padx=5)
        card.bind("<Button-1>", lambda e, a=article: self._open_article(a))

        # 内边距容器
        inner = ctk.CTkFrame(card, fg_color="transparent")
        inner.pack(fill="x", padx=18, pady=14)
        inner.grid_columnconfigure(1, weight=1)

        # 标题行
        title_text = article.get("title", "无标题")
        if article.get("is_anonymous"):
            title_text = "🔒 " + title_text
        title_lbl = ctk.CTkLabel(inner, text=title_text,
                                  font=("Microsoft YaHei", 15, "bold"),
                                  anchor="w", justify="left",
                                  text_color=("#222222", "#eeeeee"))
        title_lbl.grid(row=0, column=0, columnspan=3, sticky="w")
        title_lbl.bind("<Button-1>", lambda e, a=article: self._open_article(a))

        # 分类标签和作者
        info_frame = ctk.CTkFrame(inner, fg_color="transparent")
        info_frame.grid(row=1, column=0, columnspan=3, sticky="ew", pady=(6, 0))
        info_frame.grid_columnconfigure(2, weight=1)

        # 分类标签
        category = article.get("category", {})
        if category:
            cat_badge = ctk.CTkLabel(info_frame, text=category.get("name", ""),
                                      font=("Microsoft YaHei", 10),
                                      fg_color=("#e3f2fd", "#1a237e"),
                                      text_color=("#1565c0", "#90caf9"),
                                      corner_radius=4, padx=8, pady=2)
            cat_badge.grid(row=0, column=0, padx=(0, 8))

        # 作者
        author = article.get("user", {})
        author_text = author.get("display_name", "匿名") if not article.get("is_anonymous") else "匿名用户"
        author_lbl = ctk.CTkLabel(info_frame,
                                   text=f"👤 {author_text}",
                                   font=("Microsoft YaHei", 11),
                                   text_color=("#666666", "#aaaaaa"))
        author_lbl.grid(row=0, column=1, padx=(0, 8))
        if not article.get("is_anonymous") and author.get("id"):
            author_lbl.configure(cursor="hand2")
            author_lbl.bind("<Button-1>", lambda e, uid=author["id"]: self._open_user(uid))

        # 统计信息
        time_str = article.get("created_at", "")[:10]
        stats = f"👁 {article.get('view_count', 0)}  👍 {article.get('like_count', 0)}  💬 {article.get('comment_count', 0)}  🕐 {time_str}"
        stats_lbl = ctk.CTkLabel(info_frame, text=stats,
                                  font=("Microsoft YaHei", 10),
                                  text_color=("#888888", "#999999"))
        stats_lbl.grid(row=0, column=2, sticky="e")

        # 内容预览（简略）
        content = article.get("content", "")
        if content:
            # 去除markdown标记
            import re
            preview = re.sub(r'[#*`>\-\[\]()!]', '', content)[:120]
            if len(content) > 120:
                preview += "..."
            content_lbl = ctk.CTkLabel(inner, text=preview,
                                        font=("Microsoft YaHei", 11),
                                        text_color=("#666666", "#999999"),
                                        anchor="w", justify="left", wraplength=700)
            content_lbl.grid(row=2, column=0, columnspan=3, sticky="w", pady=(6, 0))
            content_lbl.bind("<Button-1>", lambda e, a=article: self._open_article(a))

    def _open_article(self, article):
        self.main_window.show_article_detail(article["id"], article)

    def _open_user(self, user_id):
        self.main_window.show_user_page(user_id)

    def _update_pagination(self):
        self.page_label.configure(text=f"第 {self.current_page} 页 / 共 {self.total_pages} 页")
        self.prev_btn.configure(state="normal" if self.current_page > 1 else "disabled")
        self.next_btn.configure(state="normal" if self.current_page < self.total_pages else "disabled")

    def _prev_page(self):
        if self.current_page > 1:
            self.current_page -= 1
            self._load_articles()

    def _next_page(self):
        if self.current_page < self.total_pages:
            self.current_page += 1
            self._load_articles()
