import os
import customtkinter as ctk
from ..api_client import ApiClient
from .article_detail_view import ArticleDetailView
from .create_edit_view import CreateEditView

# ── 图标映射（纯文本，后续可替换为真实图标） ──
NAV_ICONS = {
    "首页": "🏠",
    "发布文章": "✍️",
    "我的文章": "📄",
    "草稿箱": "📝",
    "收藏夹": "⭐",
    "个人资料": "👤",
    "关注管理": "👥",
    "通知": "🔔",
    "退出登录": "🚪",
}


class MainWindow(ctk.CTk):
    """主窗口 - 包含侧边栏导航和内容区域"""

    def __init__(self, api: ApiClient, user: dict):
        super().__init__()
        self.api = api
        self.user = user
        self.current_frame = None

        self.title(f"校园论坛 - {user.get('display_name', user.get('username', ''))}")
        self.geometry("1200x750")
        self.minsize(900, 600)
        self.update_idletasks()
        w, h = 1200, 750
        sw = self.winfo_screenwidth()
        sh = self.winfo_screenheight()
        self.geometry(f"{w}x{h}+{(sw-w)//2}+{(sh-h)//2}")

        # 加载分类数据
        self.categories = []
        self._load_categories()

        # 构建 UI
        self._build_ui()

        # 默认显示首页
        self._show_home()
        self._select_nav(0)

    def _load_categories(self):
        result = self.api.get_categories()
        if result and "categories" in result:
            self.categories = result["categories"]

    def _build_ui(self):
        # ── 网格布局：侧边栏(左) + 内容(右) ──
        self.grid_columnconfigure(0, weight=0, minsize=220)
        self.grid_columnconfigure(1, weight=1)
        self.grid_rowconfigure(0, weight=1)

        # ── 侧边栏 ──
        self.sidebar = ctk.CTkFrame(self, width=220, corner_radius=0)
        self.sidebar.grid(row=0, column=0, sticky="nsew")
        self.sidebar.grid_rowconfigure(6, weight=1)  # 撑开

        # 用户信息区域
        user_frame = ctk.CTkFrame(self.sidebar, fg_color="transparent")
        user_frame.pack(fill="x", padx=15, pady=(20, 10))

        ctk.CTkLabel(user_frame, text="🏫", font=("Segoe UI Emoji", 28)).pack()
        ctk.CTkLabel(user_frame, text="校园论坛", font=("Microsoft YaHei", 18, "bold"),
                      text_color=("#1a73e8", "#4fc3f7")).pack(pady=(3, 2))

        # 用户信息
        avatar_text = self.user.get("display_name", "?")[0]
        ctk.CTkLabel(user_frame, text=avatar_text,
                      font=("Microsoft YaHei", 20, "bold"),
                      width=48, height=48, corner_radius=24,
                      fg_color=("#1a73e8", "#1565c0"),
                      text_color="white").pack(pady=(8, 5))
        ctk.CTkLabel(user_frame,
                      text=self.user.get("display_name", "未知用户"),
                      font=("Microsoft YaHei", 13, "bold")).pack()
        ctk.CTkLabel(user_frame,
                      text=f"@{self.user.get('username', '')}",
                      font=("Microsoft YaHei", 10),
                      text_color=("#888888", "#aaaaaa")).pack()

        # 分隔线
        ctk.CTkFrame(self.sidebar, height=1, fg_color=("#d0d0d0", "#444444")).pack(
            fill="x", padx=15, pady=10)

        # 导航按钮列表
        self.nav_items = ["首页", "发布文章", "我的文章", "草稿箱", "收藏夹", "个人资料"]
        self.nav_buttons = []
        self.nav_vars = []

        for item in self.nav_items:
            icon = NAV_ICONS.get(item, "📌")
            # 使用更现代的按钮
            var = ctk.StringVar(value=item)
            btn = ctk.CTkButton(
                self.sidebar,
                text=f"  {icon}  {item}",
                font=("Microsoft YaHei", 13),
                anchor="w",
                height=40,
                fg_color="transparent",
                text_color=("#333333", "#cccccc"),
                hover_color=("#e8eaf6", "#2a2a4a"),
                corner_radius=8,
                command=lambda v=item, idx=len(self.nav_buttons): self._on_nav_click(v, idx)
            )
            btn.pack(fill="x", padx=10, pady=2)
            self.nav_buttons.append(btn)
            self.nav_vars.append(var)

        # 通知按钮（单独）
        self.notify_btn = ctk.CTkButton(
            self.sidebar,
            text="  🔔  通知",
            font=("Microsoft YaHei", 13),
            anchor="w", height=40,
            fg_color="transparent",
            text_color=("#333333", "#cccccc"),
            hover_color=("#e8eaf6", "#2a2a4a"),
            corner_radius=8,
            command=self._show_notifications
        )
        self.notify_btn.pack(fill="x", padx=10, pady=2)

        # 底部：退出
        ctk.CTkFrame(self.sidebar, height=1, fg_color=("#d0d0d0", "#444444")).pack(
            fill="x", padx=15, pady=10)

        ctk.CTkButton(
            self.sidebar,
            text="  🚪  退出登录",
            font=("Microsoft YaHei", 13),
            anchor="w", height=40,
            fg_color="transparent",
            text_color=("#e53935", "#ef5350"),
            hover_color=("#ffebee", "#4a2020"),
            corner_radius=8,
            command=self._do_logout
        ).pack(fill="x", padx=10, pady=(2, 15))

        # ── 内容区域 ──
        self.content = ctk.CTkFrame(self, corner_radius=0, fg_color=("#f5f5f5", "#1a1a2e"))
        self.content.grid(row=0, column=1, sticky="nsew")
        self.content.grid_columnconfigure(0, weight=1)
        self.content.grid_rowconfigure(0, weight=1)

    # ── 导航处理 ────────────────────────────────────────────────

    def _select_nav(self, idx: int):
        for i, btn in enumerate(self.nav_buttons):
            if i == idx:
                btn.configure(fg_color=("#1a73e8", "#1565c0"), text_color="white")
            else:
                btn.configure(fg_color="transparent", text_color=("#333333", "#cccccc"))
        self.notify_btn.configure(fg_color="transparent", text_color=("#333333", "#cccccc"))

    def _on_nav_click(self, item: str, idx: int):
        self._select_nav(idx)
        actions = {
            "首页": self._show_home,
            "发布文章": self._show_create,
            "我的文章": self._show_my_articles,
            "草稿箱": self._show_drafts,
            "收藏夹": self._show_favorites,
            "个人资料": self._show_profile,
        }
        action = actions.get(item)
        if action:
            action()

    def _switch_content(self, frame_class, **kwargs):
        """切换内容区域"""
        if self.current_frame:
            self.current_frame.destroy()
        self.current_frame = frame_class(self.content, self.api, self.user, main_window=self, **kwargs)
        self.current_frame.pack(fill="both", expand=True)

    # ── 各个页面 ────────────────────────────────────────────────

    def _show_home(self):
        from .home_content import HomeContent
        self._switch_content(HomeContent)

    def _show_create(self):
        self._switch_content(CreateEditView)

    def _show_my_articles(self):
        from .my_articles_view import MyArticlesView
        self._switch_content(MyArticlesView, mode="published")

    def _show_drafts(self):
        from .my_articles_view import MyArticlesView
        self._switch_content(MyArticlesView, mode="drafts")

    def _show_favorites(self):
        from .favorites_view import FavoritesView
        self._switch_content(FavoritesView)

    def _show_profile(self):
        from .profile_view import ProfileView
        self._switch_content(ProfileView)

    def _show_notifications(self):
        from .notifications_view import NotificationsView
        self._switch_content(NotificationsView)
        self.notify_btn.configure(fg_color=("#1a73e8", "#1565c0"), text_color="white")
        for btn in self.nav_buttons:
            btn.configure(fg_color="transparent", text_color=("#333333", "#cccccc"))

    def show_article_detail(self, article_id: int, article_data: dict = None):
        self._switch_content(ArticleDetailView, article_id=article_id, article_data=article_data)
        self._select_nav(-1)

    def show_create_with_data(self, article_data: dict = None):
        self._switch_content(CreateEditView, edit_data=article_data)
        self._select_nav(-1)

    def show_user_page(self, user_id: int):
        from .user_page_view import UserPageView
        self._switch_content(UserPageView, target_user_id=user_id)
        self._select_nav(-1)

    def _do_logout(self):
        self.api.clear_token()
        self.destroy()
        from .auth_views import AuthWindow
        auth = AuthWindow(self.api)
        auth.mainloop()
