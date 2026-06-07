import json
import os
import requests
from typing import Optional, Dict, Any, List


class ApiClient:
    """API客户端 - 封装所有后端接口调用"""

    def __init__(self, base_url: str):
        self.base_url = base_url.rstrip("/")
        self.token: Optional[str] = None
        self._token_file = os.path.join(os.path.dirname(__file__), "..", ".token")

    # ── Token管理 ──────────────────────────────────────────────

    def load_token(self) -> Optional[str]:
        """从本地文件加载Token"""
        try:
            with open(self._token_file, "r") as f:
                self.token = f.read().strip()
                return self.token
        except (FileNotFoundError, IOError):
            self.token = None
            return None

    def save_token(self, token: str):
        """保存Token到本地"""
        self.token = token
        os.makedirs(os.path.dirname(self._token_file), exist_ok=True)
        with open(self._token_file, "w") as f:
            f.write(token)

    def clear_token(self):
        """清除本地Token"""
        self.token = None
        try:
            os.remove(self._token_file)
        except FileNotFoundError:
            pass

    # ── HTTP请求工具 ───────────────────────────────────────────

    def _headers(self, auth: bool = True) -> Dict[str, str]:
        h = {"Content-Type": "application/json"}
        if auth and self.token:
            h["Authorization"] = f"Bearer {self.token}"
        return h

    def _request(self, method: str, path: str, **kwargs) -> Optional[Dict]:
        """通用请求方法"""
        url = f"{self.base_url}{path}"
        auth = kwargs.pop("auth", True)
        files = kwargs.pop("files", None)
        try:
            if files:
                headers = {"Authorization": f"Bearer {self.token}"} if auth and self.token else {}
                resp = requests.request(method, url, headers=headers, files=files, timeout=30)
            else:
                headers = self._headers(auth=auth)
                resp = requests.request(method, url, headers=headers, json=kwargs.get("json"), params=kwargs.get("params"), timeout=30)
            if resp.status_code == 401:
                self.clear_token()
            try:
                return resp.json()
            except ValueError:
                return {"error": f"HTTP {resp.status_code}: {resp.text}"}
        except requests.RequestException as e:
            return {"error": f"网络错误: {str(e)}"}

    def _get(self, path: str, params: Optional[Dict] = None, auth: bool = True) -> Optional[Dict]:
        return self._request("GET", path, params=params, auth=auth)

    def _post(self, path: str, data: Optional[Dict] = None, auth: bool = True, files: Optional[Dict] = None) -> Optional[Dict]:
        return self._request("POST", path, json=data, auth=auth, files=files)

    def _put(self, path: str, data: Optional[Dict] = None, auth: bool = True) -> Optional[Dict]:
        return self._request("PUT", path, json=data, auth=auth)

    def _delete(self, path: str, auth: bool = True) -> Optional[Dict]:
        return self._request("DELETE", path, auth=auth)

    # ══════════════════════════════════════════════════════════════
    #  认证接口
    # ══════════════════════════════════════════════════════════════

    def register(self, username: str, qq_number: str, display_name: str, password: str) -> Optional[Dict]:
        return self._post("/api/auth/register", data={
            "username": username, "qq_number": qq_number,
            "display_name": display_name, "password": password
        }, auth=False)

    def login(self, username: str, password: str) -> Optional[Dict]:
        result = self._post("/api/auth/login", data={"username": username, "password": password}, auth=False)
        if result and "token" in result:
            self.save_token(result["token"])
        return result

    def send_reset_code(self, qq_number: str) -> Optional[Dict]:
        return self._post("/api/password/reset-code", data={"qq_number": qq_number}, auth=False)

    def reset_password(self, qq_number: str, code: str, identifier: str, password: str) -> Optional[Dict]:
        return self._post("/api/password/reset", data={
            "qq_number": qq_number, "code": code,
            "identifier": identifier, "password": password
        }, auth=False)

    # ══════════════════════════════════════════════════════════════
    #  文章接口
    # ══════════════════════════════════════════════════════════════

    def get_articles(self, page: int = 1, page_size: int = 20, category_id: Optional[int] = None) -> Optional[Dict]:
        params = {"page": page, "page_size": page_size}
        if category_id:
            params["category_id"] = category_id
        return self._get("/api/articles", params=params, auth=False)

    def get_article_detail(self, article_id: int, page: int = 1, page_size: int = 20) -> Optional[Dict]:
        return self._get(f"/api/articles/{article_id}", params={"page": page, "page_size": page_size}, auth=False)

    def search_articles(self, keyword: str, page: int = 1, page_size: int = 20) -> Optional[Dict]:
        return self._get("/api/articles/search", params={"keyword": keyword, "page": page, "page_size": page_size}, auth=False)

    def create_article(self, title: str, content: str, category_id: int,
                       is_anonymous: bool = False, status: str = "published",
                       voice_url: Optional[str] = None) -> Optional[Dict]:
        data = {
            "title": title, "content": content, "category_id": category_id,
            "is_anonymous": is_anonymous, "status": status
        }
        if voice_url:
            data["voice_url"] = voice_url
        return self._post("/api/articles", data=data)

    def update_article(self, article_id: int, title: str, content: str,
                       category_id: int, is_anonymous: bool = False,
                       voice_url: Optional[str] = None) -> Optional[Dict]:
        data = {
            "title": title, "content": content, "category_id": category_id, "is_anonymous": is_anonymous
        }
        if voice_url:
            data["voice_url"] = voice_url
        return self._put(f"/api/articles/{article_id}", data=data)

    def delete_article(self, article_id: int) -> Optional[Dict]:
        return self._delete(f"/api/articles/{article_id}")

    def get_my_articles(self, page: int = 1, page_size: int = 20) -> Optional[Dict]:
        return self._get("/api/my/articles", params={"page": page, "page_size": page_size})

    def get_my_drafts(self, page: int = 1, page_size: int = 20) -> Optional[Dict]:
        return self._get("/api/my/drafts", params={"page": page, "page_size": page_size})

    def publish_draft(self, article_id: int) -> Optional[Dict]:
        return self._post(f"/api/articles/{article_id}/publish")

    def share_article(self, article_id: int) -> Optional[Dict]:
        return self._post(f"/api/articles/{article_id}/share")

    # ── 点赞 ──

    def like_article(self, article_id: int) -> Optional[Dict]:
        return self._post(f"/api/articles/{article_id}/like")

    def unlike_article(self, article_id: int) -> Optional[Dict]:
        return self._delete(f"/api/articles/{article_id}/like")

    # ── 收藏 ──

    def favorite_article(self, article_id: int) -> Optional[Dict]:
        return self._post(f"/api/articles/{article_id}/favorite")

    def unfavorite_article(self, article_id: int) -> Optional[Dict]:
        return self._delete(f"/api/articles/{article_id}/favorite")

    def check_favorite(self, article_id: int) -> Optional[Dict]:
        return self._get(f"/api/articles/{article_id}/favorite/check")

    def get_favorites(self) -> Optional[Dict]:
        return self._get("/api/favorites")

    # ══════════════════════════════════════════════════════════════
    #  评论接口
    # ══════════════════════════════════════════════════════════════

    def get_comments(self, article_id: int) -> Optional[Dict]:
        return self._get(f"/api/articles/{article_id}/comments", auth=False)

    def add_comment(self, article_id: int, content: str,
                    parent_id: Optional[int] = None, is_anonymous: bool = False) -> Optional[Dict]:
        data = {"content": content, "is_anonymous": is_anonymous}
        if parent_id is not None:
            data["parent_id"] = parent_id
        return self._post(f"/api/articles/{article_id}/comments", data=data)

    def delete_comment(self, comment_id: int) -> Optional[Dict]:
        return self._delete(f"/api/comments/{comment_id}")

    def like_comment(self, comment_id: int) -> Optional[Dict]:
        return self._post(f"/api/comments/{comment_id}/like")

    def unlike_comment(self, comment_id: int) -> Optional[Dict]:
        return self._delete(f"/api/comments/{comment_id}/like")

    # ══════════════════════════════════════════════════════════════
    #  关注接口
    # ══════════════════════════════════════════════════════════════

    def follow_user(self, user_id: int) -> Optional[Dict]:
        return self._post(f"/api/follow/{user_id}")

    def unfollow_user(self, user_id: int) -> Optional[Dict]:
        return self._delete(f"/api/follow/{user_id}")

    def get_following(self) -> Optional[Dict]:
        return self._get("/api/following")

    def get_followers(self) -> Optional[Dict]:
        return self._get("/api/followers")

    def get_follow_status(self, user_id: int) -> Optional[Dict]:
        return self._get(f"/api/follow/status/{user_id}")

    def get_mutual(self) -> Optional[Dict]:
        return self._get("/api/mutual")

    # ══════════════════════════════════════════════════════════════
    #  分区接口
    # ══════════════════════════════════════════════════════════════

    def get_categories(self) -> Optional[Dict]:
        return self._get("/api/categories", auth=False)

    # ══════════════════════════════════════════════════════════════
    #  用户接口
    # ══════════════════════════════════════════════════════════════

    def get_profile(self) -> Optional[Dict]:
        return self._get("/api/profile")

    def update_profile(self, display_name: Optional[str] = None, signature: Optional[str] = None) -> Optional[Dict]:
        data = {}
        if display_name is not None:
            data["display_name"] = display_name
        if signature is not None:
            data["signature"] = signature
        return self._put("/api/profile", data=data)

    def get_user_info(self, user_id: int) -> Optional[Dict]:
        return self._get(f"/api/users/{user_id}", auth=False)

    def get_user_articles(self, user_id: int, page: int = 1, page_size: int = 20) -> Optional[Dict]:
        return self._get(f"/api/users/{user_id}/articles",
                         params={"page": page, "page_size": page_size}, auth=False)

    # ══════════════════════════════════════════════════════════════
    #  上传接口
    # ══════════════════════════════════════════════════════════════

    def upload_avatar(self, file_path: str) -> Optional[Dict]:
        with open(file_path, "rb") as f:
            return self._post("/api/upload/avatar", files={"avatar": f}, auth=True)

    def upload_image(self, file_path: str) -> Optional[Dict]:
        with open(file_path, "rb") as f:
            return self._post("/api/upload/image", files={"image": f}, auth=True)

    # ══════════════════════════════════════════════════════════════
    #  通知接口
    # ══════════════════════════════════════════════════════════════

    def get_notifications(self) -> Optional[Dict]:
        return self._get("/api/notifications")

    def get_unread_count(self) -> Optional[Dict]:
        return self._get("/api/notifications/unread-count")

    def read_notification(self, notification_id: int) -> Optional[Dict]:
        return self._post(f"/api/notifications/{notification_id}/read")

    def read_all_notifications(self) -> Optional[Dict]:
        return self._post("/api/notifications/read-all")

    # ══════════════════════════════════════════════════════════════
    #  配置接口
    # ══════════════════════════════════════════════════════════════

    def get_announcement(self) -> Optional[Dict]:
        return self._get("/api/announcement", auth=False)

    def get_site_config(self) -> Optional[Dict]:
        return self._get("/api/site-config", auth=False)

    def get_sidebar_config(self) -> Optional[Dict]:
        return self._get("/api/sidebar-config", auth=False)

    # ══════════════════════════════════════════════════════════════
    #  用户状态接口
    # ══════════════════════════════════════════════════════════════

    def get_user_status(self, user_id: int) -> Optional[Dict]:
        return self._get(f"/api/user/status/{user_id}", auth=False)

    def update_user_status(self) -> Optional[Dict]:
        return self._post("/api/user/status")
