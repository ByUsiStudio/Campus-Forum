import customtkinter as ctk
from tkinter import messagebox
import re


class ArticleDetailView(ctk.CTkScrollableFrame):
    """文章详情页"""

    def __init__(self, parent, api, user, main_window, article_id: int,
                 article_data: dict = None, **kwargs):
        super().__init__(parent, fg_color="transparent")
        self.api = api
        self.user = user
        self.main_window = main_window
        self.article_id = article_id
        self.article_data = article_data or {}
        self.liked = False
        self.comment_liked = {}
        self.comments_page = 1
        self.total_comment_pages = 1

        self.grid_columnconfigure(0, weight=1)
        self._build_ui()

    def _build_ui(self):
        # 加载文章详情
        result = self.api.get_article_detail(self.article_id, page=self.comments_page)
        if not result or "article" not in result:
            ctk.CTkLabel(self, text="😅 加载文章失败",
                          font=("Microsoft YaHei", 16)).pack(pady=100)
            return

        article = result["article"]
        self.liked = result.get("liked", False)
        self.comment_liked = result.get("comment_liked", {})

        # ── 返回按钮 ──
        back_btn = ctk.CTkButton(self, text="← 返回", font=("Microsoft YaHei", 13),
                                  height=32, fg_color="transparent",
                                  text_color=("#1a73e8", "#4fc3f7"),
                                  hover_color=("#e8eaf6", "#2a2a4a"),
                                  command=lambda: self.main_window._show_home(),
                                  anchor="w", width=80)
        back_btn.pack(anchor="w", padx=20, pady=(10, 5))

        # ── 文章标题 ──
        title = article.get("title", "无标题")
        ctk.CTkLabel(self, text=title, font=("Microsoft YaHei", 22, "bold"),
                      text_color=("#1a1a1a", "#eeeeee"),
                      anchor="w", justify="left", wraplength=700).pack(
            anchor="w", padx=20, pady=(5, 5))

        # ── 元信息 ──
        meta_frame = ctk.CTkFrame(self, fg_color="transparent")
        meta_frame.pack(fill="x", padx=20, pady=(0, 15))

        author = article.get("user", {})
        anon = article.get("is_anonymous", False)
        author_name = author.get("display_name", "匿名") if not anon else "匿名用户"
        ctk.CTkLabel(meta_frame, text=f"👤 {author_name}",
                      font=("Microsoft YaHei", 12),
                      text_color=("#666666", "#aaaaaa")).pack(side="left", padx=(0, 15))

        category = article.get("category", {})
        if category:
            cat_lbl = ctk.CTkLabel(meta_frame, text=category.get("name", ""),
                                    font=("Microsoft YaHei", 11),
                                    fg_color=("#e3f2fd", "#1a237e"),
                                    text_color=("#1565c0", "#90caf9"),
                                    corner_radius=4, padx=8, pady=2)
            cat_lbl.pack(side="left", padx=(0, 15))

        created = article.get("created_at", "")[:19].replace("T", " ")
        ctk.CTkLabel(meta_frame, text=f"🕐 {created}",
                      font=("Microsoft YaHei", 11),
                      text_color=("#888888", "#999999")).pack(side="left")

        # 操作按钮区域
        action_frame = ctk.CTkFrame(self, fg_color="transparent")
        action_frame.pack(fill="x", padx=20, pady=(0, 15))

        self.like_btn = ctk.CTkButton(action_frame,
                                       text=f"👍 {article.get('like_count', 0)}",
                                       font=("Microsoft YaHei", 12),
                                       height=34,
                                       fg_color=("#e8eaf6", "#333366") if not self.liked else ("#1a73e8", "#1565c0"),
                                       text_color=("#1a73e8", "#4fc3f7") if not self.liked else "white",
                                       corner_radius=8, command=self._toggle_like)
        self.like_btn.pack(side="left", padx=(0, 8))

        fav_result = self.api.check_favorite(self.article_id)
        self.is_favorited = False
        if fav_result:
            self.is_favorited = fav_result.get("favorited", fav_result.get("is_favorited", False))

        fav_text = "⭐ 已收藏" if self.is_favorited else "☆ 收藏"
        self.fav_btn = ctk.CTkButton(action_frame, text=fav_text,
                                      font=("Microsoft YaHei", 12), height=34,
                                      fg_color=("#fff8e1", "#3a3020") if self.is_favorited else ("#e8e8e8", "#333355"),
                                      text_color=("#f57f17", "#ffd54f") if self.is_favorited else ("#666666", "#aaaaaa"),
                                      corner_radius=8, command=self._toggle_favorite)
        self.fav_btn.pack(side="left", padx=(0, 8))

        # 统计
        stats = f"👁 {article.get('view_count', 0)}   💬 {article.get('comment_count', 0)}"
        ctk.CTkLabel(action_frame, text=stats, font=("Microsoft YaHei", 11),
                      text_color=("#888888", "#aaaaaa")).pack(side="left", padx=8)

        # 编辑/删除按钮（仅作者）
        if not anon and author.get("id") == self.user.get("id"):
            ctk.CTkButton(action_frame, text="✏️ 编辑", font=("Microsoft YaHei", 12),
                           height=34, fg_color="transparent",
                           text_color=("#1a73e8", "#4fc3f7"),
                           hover_color=("#e8eaf6", "#2a2a4a"),
                           command=lambda: self._edit_article(article)).pack(side="right", padx=2)
            ctk.CTkButton(action_frame, text="🗑 删除", font=("Microsoft YaHei", 12),
                           height=34, fg_color="transparent",
                           text_color=("#e53935", "#ef5350"),
                           hover_color=("#ffebee", "#4a2020"),
                           command=lambda: self._delete_article(article["id"])).pack(side="right", padx=2)

        # ── 分隔线 ──
        ctk.CTkFrame(self, height=1, fg_color=("#d0d0d0", "#444444")).pack(
            fill="x", padx=20, pady=(0, 15))

        # ── 文章内容（简易渲染） ──
        content_html = article.get("content_html", "")
        content_md = article.get("content", "")

        content_frame = ctk.CTkFrame(self, fg_color=("#ffffff", "#1e1e38"),
                                      corner_radius=8)
        content_frame.pack(fill="x", padx=20, pady=(0, 20))
        content_frame.grid_columnconfigure(0, weight=1)

        display_text = content_html if content_html else content_md
        self._render_markdown(content_frame, display_text)

        # ── 评论区 ──
        ctk.CTkLabel(self, text=f"💬 评论 ({article.get('comment_count', 0)})",
                      font=("Microsoft YaHei", 16, "bold"),
                      text_color=("#1a73e8", "#4fc3f7")).pack(
            anchor="w", padx=20, pady=(0, 10))

        # 评论输入
        input_frame = ctk.CTkFrame(self, fg_color=("#ffffff", "#1e1e38"), corner_radius=8)
        input_frame.pack(fill="x", padx=20, pady=(0, 15))
        input_frame.grid_columnconfigure(0, weight=1)

        self.comment_entry = ctk.CTkTextbox(input_frame, height=70,
                                              font=("Microsoft YaHei", 12),
                                              fg_color=("#f8f9fa", "#252540"),
                                              corner_radius=6)
        self.comment_entry.grid(row=0, column=0, padx=12, pady=(12, 8), sticky="ew")

        self.anon_comment_var = ctk.BooleanVar(value=False)
        anon_check = ctk.CTkCheckBox(input_frame, text="匿名评论",
                                      variable=self.anon_comment_var,
                                      font=("Microsoft YaHei", 11))
        anon_check.grid(row=1, column=0, sticky="w", padx=12, pady=(0, 8))

        ctk.CTkButton(input_frame, text="发表评论", font=("Microsoft YaHei", 13, "bold"),
                       height=36, corner_radius=6,
                       command=self._do_comment).grid(
            row=1, column=0, sticky="e", padx=12, pady=(0, 8))

        # ── 评论列表 ──
        self.comments_frame = ctk.CTkFrame(self, fg_color="transparent")
        self.comments_frame.pack(fill="x", padx=20, pady=(0, 20))

        comments = result.get("comments", [])
        if comments:
            for comment in comments:
                self._render_comment(comment)
        else:
            ctk.CTkLabel(self.comments_frame, text="暂无评论，来发表第一条评论吧",
                          font=("Microsoft YaHei", 12),
                          text_color=("#888888", "#aaaaaa")).pack(pady=20)

    def _render_markdown(self, parent, text):
        """简易Markdown渲染"""
        lines = text.split("\n")
        in_code_block = False
        code_lines = []

        for line in lines:
            # 代码块
            if line.strip().startswith("```"):
                if in_code_block:
                    # 结束代码块
                    code_text = "\n".join(code_lines)
                    code_widget = ctk.CTkTextbox(parent, height=min(200, len(code_lines) * 22 + 20),
                                                  font=("Consolas", 12),
                                                  fg_color=("#1e1e2e", "#0a0a1a"),
                                                  text_color=("#e0e0e0", "#e0e0e0"),
                                                  corner_radius=6)
                    code_widget.insert("1.0", code_text)
                    code_widget.configure(state="disabled")
                    code_widget.pack(fill="x", padx=10, pady=5)
                    code_lines = []
                    in_code_block = False
                else:
                    in_code_block = True
                continue
            if in_code_block:
                code_lines.append(line)
                continue

            # 空行
            if not line.strip():
                continue

            # 标题
            if line.startswith("### "):
                ctk.CTkLabel(parent, text=line[4:],
                              font=("Microsoft YaHei", 14, "bold"),
                              anchor="w", justify="left", wraplength=680).pack(
                    fill="x", padx=10, pady=4)
            elif line.startswith("## "):
                ctk.CTkLabel(parent, text=line[3:],
                              font=("Microsoft YaHei", 16, "bold"),
                              anchor="w", justify="left", wraplength=680).pack(
                    fill="x", padx=10, pady=4)
            elif line.startswith("# "):
                ctk.CTkLabel(parent, text=line[2:],
                              font=("Microsoft YaHei", 18, "bold"),
                              anchor="w", justify="left", wraplength=680).pack(
                    fill="x", padx=10, pady=4)
            elif line.startswith("- ") or line.startswith("* "):
                ctk.CTkLabel(parent, text=f"  • {line[2:]}",
                              font=("Microsoft YaHei", 12),
                              anchor="w", justify="left", wraplength=680).pack(
                    fill="x", padx=10, pady=2)
            # 数字列表
            elif line.strip() and line[0].isdigit() and line.strip()[1:2] == ".":
                ctk.CTkLabel(parent, text=line,
                              font=("Microsoft YaHei", 12),
                              anchor="w", justify="left", wraplength=680).pack(
                    fill="x", padx=10, pady=2)
            else:
                # 普通文本 - 去除链接标记、加粗等
                text_clean = re.sub(r'\[([^\]]+)\]\([^)]+\)', r'\1', line)
                text_clean = re.sub(r'\*\*(.+?)\*\*', r'\1', text_clean)
                text_clean = re.sub(r'\*(.+?)\*', r'\1', text_clean)
                ctk.CTkLabel(parent, text=text_clean,
                              font=("Microsoft YaHei", 12),
                              anchor="w", justify="left", wraplength=680).pack(
                    fill="x", padx=10, pady=3)

    def _render_comment(self, comment, parent=None):
        if parent is None:
            parent = self.comments_frame

        is_reply = comment.get("parent_id") is not None

        card = ctk.CTkFrame(parent, corner_radius=8,
                             fg_color=("#ffffff", "#1e1e38"),
                             border_width=0)
        card.pack(fill="x", pady=4)
        inner = ctk.CTkFrame(card, fg_color="transparent")
        inner.pack(fill="x", padx=14, pady=10)

        user_info = comment.get("user", {})
        anon = comment.get("is_anonymous", False)
        name = user_info.get("display_name", "匿名") if not anon else "匿名用户"

        # 用户信息行
        name_lbl = ctk.CTkLabel(inner, text=f"👤 {name}",
                                 font=("Microsoft YaHei", 11, "bold"),
                                 text_color=("#333333", "#cccccc"))
        name_lbl.pack(anchor="w")

        if not anon and user_info.get("id"):
            name_lbl.configure(cursor="hand2")
            name_lbl.bind("<Button-1>", lambda e, uid=user_info["id"]: self._open_user(uid))

        time_str = comment.get("created_at", "")[:19].replace("T", " ")
        ctk.CTkLabel(inner, text=time_str,
                      font=("Microsoft YaHei", 10),
                      text_color=("#999999", "#888888")).pack(anchor="w")

        # 评论内容
        ctk.CTkLabel(inner, text=comment.get("content", ""),
                      font=("Microsoft YaHei", 12),
                      anchor="w", justify="left", wraplength=600).pack(
            anchor="w", pady=(4, 4))

        # 操作行
        action_row = ctk.CTkFrame(inner, fg_color="transparent")
        action_row.pack(fill="x")

        cid = comment["id"]
        is_liked = self.comment_liked.get(str(cid), False)
        like_txt = f"👍 {comment.get('like_count', 0)}"
        like_btn = ctk.CTkButton(action_row, text=like_txt,
                                  font=("Microsoft YaHei", 10), height=26,
                                  fg_color="transparent",
                                  text_color=("#1a73e8", "#4fc3f7") if not is_liked else "#ff6f00",
                                  hover_color=("#e8eaf6", "#2a2a4a"),
                                  command=lambda cid=cid, btn_ref=[]: self._toggle_comment_like(cid, btn_ref))
        like_btn.pack(side="left", padx=(0, 8))

        # 回复按钮
        ctk.CTkButton(action_row, text="💬 回复",
                       font=("Microsoft YaHei", 10), height=26,
                       fg_color="transparent",
                       text_color=("#666666", "#aaaaaa"),
                       hover_color=("#e8eaf6", "#2a2a4a"),
                       command=lambda cid=cid: self._show_reply_input(cid)).pack(side="left")

        # 删除（仅作者）
        if not anon and user_info.get("id") == self.user.get("id"):
            ctk.CTkButton(action_row, text="删除",
                           font=("Microsoft YaHei", 10), height=26,
                           fg_color="transparent",
                           text_color=("#e53935", "#ef5350"),
                           hover_color=("#ffebee", "#4a2020"),
                           command=lambda cid=cid: self._delete_comment(cid)).pack(side="right")

        # 回复列表
        replies = comment.get("replies", [])
        if replies:
            reply_frame = ctk.CTkFrame(card, fg_color="transparent")
            reply_frame.pack(fill="x", padx=(30, 14), pady=(0, 8))
            for reply in replies:
                self._render_comment(reply, parent=reply_frame)

    def _show_reply_input(self, parent_id):
        """显示回复输入框"""
        dialog = ctk.CTkToplevel(self)
        dialog.title("回复评论")
        dialog.geometry("450x250")
        dialog.resizable(False, False)
        dialog.transient(self)

        ctk.CTkLabel(dialog, text="回复评论",
                      font=("Microsoft YaHei", 16, "bold")).pack(pady=(15, 10))

        text_box = ctk.CTkTextbox(dialog, height=80, font=("Microsoft YaHei", 12))
        text_box.pack(fill="x", padx=20, pady=(0, 10))

        anon_var = ctk.BooleanVar(value=False)
        ctk.CTkCheckBox(dialog, text="匿名回复",
                         variable=anon_var).pack(padx=20, anchor="w")

        def submit():
            content = text_box.get("1.0", "end-1c").strip()
            if not content:
                messagebox.showerror("错误", "请输入回复内容")
                return
            result = self.api.add_comment(self.article_id, content,
                                           parent_id=parent_id,
                                           is_anonymous=anon_var.get())
            if result and "comment" in result:
                dialog.destroy()
                self._rebuild()
            else:
                error = result.get("error", "回复失败") if result else "网络错误"
                messagebox.showerror("错误", error)

        btn_frame = ctk.CTkFrame(dialog, fg_color="transparent")
        btn_frame.pack(fill="x", padx=20, pady=(10, 15))
        ctk.CTkButton(btn_frame, text="取消", fg_color="transparent",
                       command=dialog.destroy).pack(side="left", padx=5)
        ctk.CTkButton(btn_frame, text="提交", fg_color=("#1a73e8", "#1565c0"),
                       command=submit).pack(side="right", padx=5)

    def _toggle_like(self):
        if self.liked:
            self.api.unlike_article(self.article_id)
            self.liked = False
            self.like_btn.configure(
                fg_color=("#e8eaf6", "#333366"),
                text_color=("#1a73e8", "#4fc3f7"))
        else:
            self.api.like_article(self.article_id)
            self.liked = True
            self.like_btn.configure(
                fg_color=("#1a73e8", "#1565c0"),
                text_color="white")

    def _toggle_favorite(self):
        if self.is_favorited:
            self.api.unfavorite_article(self.article_id)
            self.is_favorited = False
            self.fav_btn.configure(text="☆ 收藏",
                                    fg_color=("#e8e8e8", "#333355"),
                                    text_color=("#666666", "#aaaaaa"))
        else:
            self.api.favorite_article(self.article_id)
            self.is_favorited = True
            self.fav_btn.configure(text="⭐ 已收藏",
                                    fg_color=("#fff8e1", "#3a3020"),
                                    text_color=("#f57f17", "#ffd54f"))

    def _toggle_comment_like(self, comment_id, btn_ref):
        is_liked = self.comment_liked.get(str(comment_id), False)
        if is_liked:
            self.api.unlike_comment(comment_id)
            self.comment_liked[str(comment_id)] = False
        else:
            self.api.like_comment(comment_id)
            self.comment_liked[str(comment_id)] = True
        self._rebuild()

    def _do_comment(self):
        content = self.comment_entry.get("1.0", "end-1c").strip()
        if not content:
            messagebox.showerror("错误", "请输入评论内容")
            return
        result = self.api.add_comment(self.article_id, content,
                                       is_anonymous=self.anon_comment_var.get())
        if result and "comment" in result:
            self.comment_entry.delete("1.0", "end")
            self._rebuild()
        else:
            error = result.get("error", "评论失败") if result else "网络错误"
            messagebox.showerror("错误", error)

    def _delete_comment(self, comment_id):
        if messagebox.askyesno("确认", "确定要删除这条评论吗？"):
            result = self.api.delete_comment(comment_id)
            if result and "message" in result:
                self._rebuild()
            else:
                error = result.get("error", "删除失败") if result else "网络错误"
                messagebox.showerror("错误", error)

    def _delete_article(self, article_id):
        if messagebox.askyesno("确认删除", "确定要删除这篇文章吗？\n此操作不可恢复。"):
            result = self.api.delete_article(article_id)
            if result and "message" in result:
                messagebox.showinfo("成功", "文章已删除")
                self.main_window._show_home()
            else:
                error = result.get("error", "删除失败") if result else "网络错误"
                messagebox.showerror("错误", error)

    def _edit_article(self, article):
        self.main_window.show_create_with_data(article)

    def _rebuild(self):
        """刷新页面"""
        for w in self.winfo_children():
            w.destroy()
        self._build_ui()

    def _open_user(self, user_id):
        self.main_window.show_user_page(user_id)
