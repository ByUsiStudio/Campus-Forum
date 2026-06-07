import customtkinter as ctk
from tkinter import messagebox


class CreateEditView(ctk.CTkScrollableFrame):
    """创建/编辑文章页面"""

    def __init__(self, parent, api, user, main_window, edit_data: dict = None, **kwargs):
        super().__init__(parent, fg_color="transparent")
        self.api = api
        self.user = user
        self.main_window = main_window
        self.edit_data = edit_data
        self.is_edit = edit_data is not None

        self.grid_columnconfigure(0, weight=1)
        self._build_ui()

    def _build_ui(self):
        # 标题
        title_text = "✏️ 编辑文章" if self.is_edit else "✍️ 发布文章"
        ctk.CTkLabel(self, text=title_text, font=("Microsoft YaHei", 22, "bold"),
                      text_color=("#1a73e8", "#4fc3f7")).pack(
            anchor="w", padx=30, pady=(20, 20))

        # ── 标题 ──
        ctk.CTkLabel(self, text="文章标题", font=("Microsoft YaHei", 13, "bold"),
                      anchor="w").pack(fill="x", padx=30, pady=(0, 5))
        self.title_entry = ctk.CTkEntry(self, placeholder_text="请输入文章标题",
                                         height=40, font=("Microsoft YaHei", 14))
        self.title_entry.pack(fill="x", padx=30, pady=(0, 15))
        if self.is_edit:
            self.title_entry.insert(0, self.edit_data.get("title", ""))

        # ── 分类 ──
        ctk.CTkLabel(self, text="选择分类", font=("Microsoft YaHei", 13, "bold"),
                      anchor="w").pack(fill="x", padx=30, pady=(0, 5))

        categories = getattr(self.main_window, "categories", [])
        cat_names = ["请选择分类"] + [c.get("name", f"分类{c['id']}") for c in categories]
        self.category_var = ctk.StringVar(value=cat_names[0])

        if self.is_edit and self.edit_data.get("category"):
            edit_cat_name = self.edit_data["category"].get("name", "")
            if edit_cat_name in cat_names:
                self.category_var.set(edit_cat_name)

        self.cat_combo = ctk.CTkComboBox(self, variable=self.category_var,
                                          values=cat_names,
                                          font=("Microsoft YaHei", 13),
                                          height=38, state="readonly")
        self.cat_combo.pack(fill="x", padx=30, pady=(0, 15))

        # ── 内容 ──
        ctk.CTkLabel(self, text="文章内容（支持 Markdown）",
                      font=("Microsoft YaHei", 13, "bold"),
                      anchor="w").pack(fill="x", padx=30, pady=(0, 5))

        self.content_text = ctk.CTkTextbox(self, height=350,
                                            font=("Microsoft YaHei", 12),
                                            corner_radius=8,
                                            fg_color=("#ffffff", "#1a1a35"),
                                            border_width=1,
                                            border_color=("#d0d0d0", "#444466"))
        self.content_text.pack(fill="x", padx=30, pady=(0, 5))
        if self.is_edit:
            self.content_text.insert("1.0", self.edit_data.get("content", ""))

        # 内容提示
        ctk.CTkLabel(self, text="💡 提示：支持 Markdown 语法，可用 **加粗**、*斜体*、# 标题、```代码``` 等",
                      font=("Microsoft YaHei", 10),
                      text_color=("#888888", "#999999")).pack(
            anchor="w", padx=30, pady=(0, 15))

        # ── 选项 ──
        options_frame = ctk.CTkFrame(self, fg_color="transparent")
        options_frame.pack(fill="x", padx=30, pady=(0, 15))

        self.anon_var = ctk.BooleanVar(
            value=self.edit_data.get("is_anonymous", False) if self.is_edit else False)
        ctk.CTkCheckBox(options_frame, text="匿名发布",
                         variable=self.anon_var,
                         font=("Microsoft YaHei", 12)).pack(side="left", padx=(0, 20))

        if not self.is_edit:
            self.draft_var = ctk.BooleanVar(value=False)
            ctk.CTkCheckBox(options_frame, text="保存为草稿",
                             variable=self.draft_var,
                             font=("Microsoft YaHei", 12)).pack(side="left")

        # ── 按钮 ──
        btn_frame = ctk.CTkFrame(self, fg_color="transparent")
        btn_frame.pack(fill="x", padx=30, pady=(5, 20))
        btn_frame.grid_columnconfigure(0, weight=1)
        btn_frame.grid_columnconfigure(1, weight=1)

        ctk.CTkButton(btn_frame, text="取消", font=("Microsoft YaHei", 14),
                       height=40, fg_color=("#e0e0e0", "#333355"),
                       text_color=("#333333", "#cccccc"),
                       corner_radius=8,
                       command=lambda: self.main_window._show_home()).grid(
            row=0, column=0, padx=(0, 5), sticky="ew")

        submit_text = "更新文章" if self.is_edit else "发布文章"
        self.submit_btn = ctk.CTkButton(btn_frame, text=submit_text,
                                         font=("Microsoft YaHei", 14, "bold"),
                                         height=40, corner_radius=8,
                                         command=self._do_submit)
        self.submit_btn.grid(row=0, column=1, padx=(5, 0), sticky="ew")

    def _get_category_id(self) -> int:
        selected = self.category_var.get()
        categories = getattr(self.main_window, "categories", [])
        for c in categories:
            if c.get("name") == selected:
                return c["id"]
        return 0

    def _do_submit(self):
        title = self.title_entry.get().strip()
        content = self.content_text.get("1.0", "end-1c").strip()
        category_id = self._get_category_id()

        if not title:
            messagebox.showerror("错误", "请输入文章标题")
            return
        if not content:
            messagebox.showerror("错误", "请输入文章内容")
            return
        if category_id == 0:
            messagebox.showerror("错误", "请选择文章分类")
            return

        self.submit_btn.configure(state="disabled", text="提交中...")
        self.update()

        if self.is_edit:
            result = self.api.update_article(
                self.edit_data["id"], title, content, category_id,
                is_anonymous=self.anon_var.get())
            success_msg = "文章已更新"
        else:
            status = "draft" if self.draft_var.get() else "published"
            result = self.api.create_article(
                title, content, category_id,
                is_anonymous=self.anon_var.get(), status=status)
            success_msg = "草稿已保存" if status == "draft" else "文章已发布"

        if result and "message" in result:
            messagebox.showinfo("成功", success_msg)
            self.main_window._show_home()
        elif result and "article" in result:
            messagebox.showinfo("成功", success_msg)
            self.main_window._show_home()
        else:
            error = result.get("error", "提交失败") if result else "网络错误"
            messagebox.showerror("错误", error)
            self.submit_btn.configure(state="normal", text="发布文章")
