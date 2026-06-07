# 校园论坛 API 文档

## 概述

本文档描述了校园论坛后端 API 的详细接口说明。

**Base URL**: `http://0.0.0.0:3620`

**认证方式**: JWT Bearer Token

---

## 认证接口

### 用户注册

**POST** `/api/auth/register`

注册新用户账号。

**请求体**:
```json
{
  "username": "string",
  "qq_number": "string",
  "display_name": "string",
  "password": "string"
}
```

**响应**:
```json
{
  "message": "注册成功"
}
```

---

### 用户登录

**POST** `/api/auth/login`

用户登录获取 Token。

**请求体**:
```json
{
  "username": "string",
  "password": "string"
}
```

**响应**:
```json
{
  "token": "jwt_token_string",
  "user": {
    "id": 1,
    "username": "string",
    "display_name": "string",
    "avatar": "string",
    "role": "user"
  }
}
```

---

### 发送密码重置验证码

**POST** `/api/password/reset-code`

发送密码重置验证码到用户QQ邮箱。

**请求体**:
```json
{
  "qq_number": "string"
}
```

**响应**:
```json
{
  "message": "验证码已发送到您的QQ邮箱",
  "identifier": "string (标识token，用于后续密码重置请求)"
}
```

---

### 重置密码

**POST** `/api/password/reset`

使用验证码和标识token重置密码。

**请求体**:
```json
{
  "qq_number": "string",
  "code": "string (验证码)",
  "identifier": "string (发送验证码时返回的标识token)",
  "password": "string (新密码，至少6位)"
}
```

**响应**:
```json
{
  "message": "密码重置成功"
}
```

---

## 文章接口

### 获取文章列表

**GET** `/api/articles`

获取文章列表，支持分页和分类筛选。

**Query 参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| page | int | 页码，默认 1 |
| page_size | int | 每页数量，默认 20 |
| category_id | int | 分类 ID |

**响应**:
```json
{
  "articles": [
    {
      "id": 1,
      "title": "string",
      "content": "string",
      "content_html": "string",
      "user": {...},
      "category": {...},
      "view_count": 100,
      "like_count": 10,
      "comment_count": 5,
      "voice_url": "string",
      "is_anonymous": false,
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 100,
  "page": 1,
  "total_pages": 5
}
```

---

### 获取文章详情

**GET** `/api/articles/{id}`

获取单篇文章详情，包含评论列表。

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 文章 ID |

**Query 参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| page | int | 评论页码，默认 1 |
| page_size | int | 每页评论数量，默认 20 |

**响应**:
```json
{
  "article": {
    "id": 1,
    "title": "string",
    "content": "string",
    "content_html": "string",
    "user": {...},
    "category": {...},
    "view_count": 100,
    "like_count": 10,
    "comment_count": 5,
    "voice_url": "string",
    "is_anonymous": false,
    "created_at": "2024-01-01T00:00:00Z"
  },
  "comments": [
    {
      "id": 1,
      "content": "string",
      "user": {...},
      "article_id": 1,
      "parent_id": null,
      "like_count": 0,
      "reply_count": 2,
      "replies": [...],
      "is_anonymous": false,
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 10,
  "page": 1,
  "page_size": 20,
  "total_pages": 1,
  "liked": false,
  "comment_liked": {
    "1": true
  }
}
```

---

### 搜索文章

**GET** `/api/articles/search`

搜索文章（支持按标题或内容搜索）。

**Query 参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| keyword | string | 搜索关键词 |
| page | int | 页码，默认 1 |
| page_size | int | 每页数量，默认 20 |

**响应**:
```json
{
  "articles": [
    {
      "id": 1,
      "title": "string",
      "content": "string",
      "content_html": "string",
      "user": {...},
      "category": {...},
      "view_count": 100,
      "like_count": 10,
      "comment_count": 5,
      "share_count": 0,
      "voice_url": "string",
      "is_anonymous": false,
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 100,
  "page": 1,
  "page_size": 20,
  "total_pages": 5
}
```

---

### 分享文章

**POST** `/api/articles/{id}/share`

分享文章（需认证），增加分享次数。

**Headers**:
```
Authorization: Bearer <token>
```

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 文章 ID |

**响应**:
```json
{
  "message": "分享成功",
  "share_count": 1
}
```

---

### 获取我的文章

**GET** `/api/my/articles`

获取当前用户的文章列表（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**Query 参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| page | int | 页码，默认 1 |
| page_size | int | 每页数量，默认 20 |

**响应**:
```json
{
  "articles": [
    {
      "id": 1,
      "title": "string",
      "content": "string",
      "content_html": "string",
      "user": {...},
      "category": {...},
      "view_count": 100,
      "like_count": 10,
      "comment_count": 5,
      "share_count": 0,
      "voice_url": "string",
      "is_anonymous": false,
      "status": "published",
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 20,
  "page": 1,
  "page_size": 20,
  "total_pages": 1
}
```

---

### 获取我的草稿

**GET** `/api/my/drafts`

获取当前用户的草稿列表（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**Query 参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| page | int | 页码，默认 1 |
| page_size | int | 每页数量，默认 20 |

**响应**:
```json
{
  "articles": [...],
  "total": 10,
  "page": 1,
  "page_size": 20,
  "total_pages": 1
}
```

---

### 创建文章

**POST** `/api/articles`

创建新文章（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "title": "string",
  "content": "string (Markdown)",
  "category_id": 1,
  "voice_url": "string (可选，语音文件URL)",
  "is_anonymous": false,
  "status": "published (published/draft)"
}
```

**响应**:
```json
{
  "message": "创建成功",
  "article": {...}
}
```

---

### 更新文章

**PUT** `/api/articles/{id}`

更新文章（需认证，仅作者或管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "title": "string",
  "content": "string",
  "category_id": 1,
  "voice_url": "string",
  "is_anonymous": false
}
```

---

### 删除文章

**DELETE** `/api/articles/{id}`

删除文章（软删除）（需认证，仅作者或管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

---

### 恢复文章

**POST** `/api/articles/{id}/restore`

恢复已删除的文章（需管理员权限）。

**Headers**:
```
Authorization: Bearer <token>
```

---

### 置顶文章

**POST** `/api/articles/{id}/pin`

置顶文章（需管理员权限）。

**Headers**:
```
Authorization: Bearer <token>
```

---

### 取消置顶

**DELETE** `/api/articles/{id}/pin`

取消文章置顶（需管理员权限）。

**Headers**:
```
Authorization: Bearer <token>
```

---

### 发布草稿

**POST** `/api/articles/{id}/publish`

发布草稿（需认证，仅作者）。

**Headers**:
```
Authorization: Bearer <token>
```

---

### 点赞文章

**POST** `/api/articles/{id}/like`

点赞文章（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

---

### 取消点赞

**DELETE** `/api/articles/{id}/like`

取消文章点赞（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

---

### 收藏文章

**POST** `/api/articles/{id}/favorite`

收藏文章（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

---

### 取消收藏

**DELETE** `/api/articles/{id}/favorite`

取消收藏（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

---

### 检查收藏状态

**GET** `/api/articles/{id}/favorite/check`

检查文章收藏状态（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

---

### 获取收藏列表

**GET** `/api/favorites`

获取当前用户的收藏列表（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

---

## 评论接口

### 获取文章评论

**GET** `/api/articles/{id}/comments`

获取文章的评论列表。

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 文章 ID |

**响应**:
```json
{
  "comments": [
    {
      "id": 1,
      "content": "string",
      "user": {
        "id": 1,
        "username": "string",
        "display_name": "string",
        "avatar": "string"
      },
      "article_id": 1,
      "parent_id": null,
      "like_count": 3,
      "reply_count": 5,
      "replies": [
        {
          "id": 2,
          "content": "string",
          "user": {...},
          "parent_id": 1,
          "is_anonymous": false,
          "created_at": "2024-01-01T00:00:00Z"
        }
      ],
      "is_anonymous": false,
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

---

### 添加评论

**POST** `/api/articles/{id}/comments`

添加评论（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 文章 ID |

**请求体**:
```json
{
  "content": "string (评论内容，支持HTML过滤)",
  "parent_id": null (回复评论时填写父评论ID，顶级评论传null或不传),
  "is_anonymous": false (是否匿名评论)
}
```

**响应**:
```json
{
  "message": "评论成功",
  "comment": {
    "id": 1,
    "content": "string",
    "user": {...},
    "is_anonymous": false,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

### 删除评论

**DELETE** `/api/comments/{id}`

删除评论（需认证，仅作者或管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 评论 ID |

**响应**:
```json
{
  "message": "删除成功"
}
```

---

### 点赞评论

**POST** `/api/comments/{id}/like`

点赞评论（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 评论 ID |

**响应**:
```json
{
  "message": "点赞成功"
}
```

---

### 取消评论点赞

**DELETE** `/api/comments/{id}/like`

取消评论点赞（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 评论 ID |

**响应**:
```json
{
  "message": "取消点赞成功"
}
```

---

## 关注接口

### 关注用户

**POST** `/api/follow/{id}`

关注用户（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 被关注用户 ID |

---

### 取消关注

**DELETE** `/api/follow/{id}`

取消关注（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

---

### 获取关注列表

**GET** `/api/following`

获取当前用户的关注列表（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "following": [
    {
      "id": 1,
      "username": "string",
      "display_name": "string",
      "avatar": "string",
      "role": "user",
      "signature": "string"
    }
  ]
}
```

---

### 获取粉丝列表

**GET** `/api/followers`

获取当前用户的粉丝列表（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "followers": [...]
}
```

---

### 检查关注状态

**GET** `/api/follow/status/{id}`

检查与指定用户的关注状态（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "is_following": true,
  "is_followed": false,
  "mutual": false,
  "following_user": {...}
}
```

---

### 获取互相关注好友

**GET** `/api/mutual`

获取互相关注的好友列表（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

---

### 获取用户关注列表

**GET** `/api/users/{id}/following`

获取指定用户的关注列表（公开接口）。

**响应**:
```json
{
  "following": [...]
}
```

---

### 获取用户粉丝列表

**GET** `/api/users/{id}/followers`

获取指定用户的粉丝列表（公开接口）。

**响应**:
```json
{
  "followers": [...]
}
```

---

## 分区接口

### 获取分区列表

**GET** `/api/categories`

获取所有分区。

**响应**:
```json
{
  "categories": [
    {
      "id": 1,
      "name": "技术"
    }
  ]
}
```

---

### 创建分区

**POST** `/api/categories`

创建新分区（需管理员权限）。

**Headers**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "name": "string"
}
```

---

### 更新分区

**PUT** `/api/categories/{id}`

更新分区（需管理员权限）。

**Headers**:
```
Authorization: Bearer <token>
```

---

### 删除分区

**DELETE** `/api/categories/{id}`

删除分区（需管理员权限）。

**Headers**:
```
Authorization: Bearer <token>
```

---

## 用户接口

### 获取个人资料

**GET** `/api/profile`

获取当前登录用户信息（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "id": 1,
  "username": "string",
  "display_name": "string",
  "email": "string",
  "avatar": "string",
  "signature": "string",
  "role": "user",
  "qq_number": "string",
  "created_at": "2024-01-01T00:00:00Z"
}
```

---

### 更新个人资料

**PUT** `/api/profile`

更新当前用户信息（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "display_name": "string",
  "signature": "string"
}
```

---

### 获取用户公开信息

**GET** `/api/users/{id}`

获取指定用户的公开信息。

**响应**:
```json
{
  "id": 1,
  "username": "string",
  "display_name": "string",
  "avatar": "string",
  "role": "user",
  "signature": "string",
  "titles": [...],
  "created_at": "2024-01-01T00:00:00Z"
}
```

---

### 获取用户文章列表

**GET** `/api/users/{id}/articles`

获取指定用户的文章列表（支持分页）。

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 用户 ID |

**Query 参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| page | int | 页码，默认 1 |
| page_size | int | 每页数量，默认 20 |

**响应**:
```json
{
  "articles": [
    {
      "id": 1,
      "title": "string",
      "content": "string",
      "content_html": "string",
      "user": {...},
      "category": {...},
      "view_count": 100,
      "like_count": 10,
      "comment_count": 5,
      "share_count": 0,
      "voice_url": "string",
      "is_anonymous": false,
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 50,
  "page": 1,
  "page_size": 20,
  "total_pages": 3
}
```

---

## 上传接口

### 上传头像

**POST** `/api/upload/avatar`

上传用户头像（需认证）。

**Headers**:
```
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**Form Data**:
| 字段 | 类型 | 描述 |
|------|------|------|
| avatar | file | 图片文件 |

**响应**:
```json
{
  "message": "上传成功",
  "url": "https://..."
}
```

---

### 上传图片

**POST** `/api/upload/image`

上传文章图片（需认证）。

**Headers**:
```
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**Form Data**:
| 字段 | 类型 | 描述 |
|------|------|------|
| image | file | 图片文件 |

**响应**:
```json
{
  "message": "上传成功",
  "url": "https://..."
}
```

---

### 上传视频

**POST** `/api/upload/video`

上传视频（需认证）。

**Headers**:
```
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**Form Data**:
| 字段 | 类型 | 描述 |
|------|------|------|
| video | file | 视频文件 |

**响应**:
```json
{
  "message": "上传成功",
  "url": "https://..."
}
```

---

### 上传语音

**POST** `/api/upload/voice`

上传语音文件（需认证）。

**Headers**:
```
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**Form Data**:
| 字段 | 类型 | 描述 |
|------|------|------|
| voice | file | 语音文件 |

**响应**:
```json
{
  "message": "上传成功",
  "url": "https://..."
}
```

---

## 通知接口

### 获取通知列表

**GET** `/api/notifications`

获取当前用户的通知列表（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

---

### 获取未读通知数量

**GET** `/api/notifications/unread-count`

获取未读通知数量（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

---

### 标记通知为已读

**POST** `/api/notifications/{id}/read`

标记单条通知为已读（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

---

### 标记所有通知为已读

**POST** `/api/notifications/read-all`

标记所有通知为已读（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

---

### 获取评论回复通知

**GET** `/api/comment-reply-notifications`

获取评论回复通知列表（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

---

### 获取粉丝通知

**GET** `/api/follow-notifications`

获取粉丝通知列表（关注对象发新内容时的通知）（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "notifications": [
    {
      "id": 1,
      "user_id": 123,
      "sender_id": 456,
      "article_id": 789,
      "type": "new_article",
      "is_read": false,
      "sender": {...},
      "article": {...},
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

---

### 标记粉丝通知为已读

**POST** `/api/follow-notifications/{id}/read`

标记单条粉丝通知为已读（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

---

### 标记所有粉丝通知为已读

**POST** `/api/follow-notifications/read-all`

标记所有粉丝通知为已读（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

---

### 获取粉丝通知未读数量

**GET** `/api/follow-notifications/unread-count`

获取粉丝通知未读数量（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

---

## 举报接口

### 创建举报

**POST** `/api/reports`

创建举报（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "target_type": "article",
  "target_id": 1,
  "reason": "string",
  "description": "string"
}
```

**响应**:
```json
{
  "message": "举报已提交"
}
```

---

### 获取举报列表

**GET** `/api/reports`

获取举报列表（需管理员权限）。

**Headers**:
```
Authorization: Bearer <token>
```

**Query 参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| page | int | 页码，默认 1 |
| page_size | int | 每页数量，默认 20 |
| status | string | 状态筛选（pending/resolved/dismissed） |

**响应**:
```json
{
  "reports": [...],
  "total": 100,
  "page": 1,
  "total_pages": 5
}
```

---

### 获取举报详情

**GET** `/api/reports/:id`

获取举报详情（需管理员权限）。

**Headers**:
```
Authorization: Bearer <token>
```

---

### 处理举报

**PUT** `/api/reports/:id/handle`

处理举报（需管理员权限）。

**Headers**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "status": "resolved",
  "action_taken": "string"
}
```

---

## 用户通知接口

### 获取用户通知列表

**GET** `/api/user-notifications`

获取当前用户的通知列表（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**Query 参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| page | int | 页码，默认 1 |
| page_size | int | 每页数量，默认 20 |
| type | string | 通知类型筛选 |

---

### 获取单个通知

**GET** `/api/user-notifications/:id`

获取单个通知详情（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

---

### 标记通知为已读

**POST** `/api/user-notifications/:id/read`

标记通知为已读（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

---

### 标记所有通知为已读

**POST** `/api/user-notifications/read-all`

标记所有通知为已读（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

---

### 删除通知

**DELETE** `/api/user-notifications/:id`

删除通知（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

---

### 清空所有通知

**DELETE** `/api/user-notifications/clear`

清空所有通知（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

---

### 发送单独通知

**POST** `/api/user-notifications/send`

发送通知给指定用户（需管理员权限）。

**Headers**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "user_id": 1,
  "title": "string",
  "content": "string",
  "type": "system",
  "priority": "normal",
  "link": "/article/123"
}
```

---

### 批量发送通知

**POST** `/api/user-notifications/send-batch`

批量发送通知（需管理员权限）。

**Headers**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "user_ids": [1, 2, 3],
  "title": "string",
  "content": "string",
  "type": "system",
  "priority": "normal",
  "link": "/article/123"
}
```

---

### 管理员获取用户通知

**GET** `/api/admin/user-notifications/:user_id`

管理员获取指定用户的通知（需管理员权限）。

**Headers**:
```
Authorization: Bearer <token>
```

---

## 权限组接口

### 获取所有权限组

**GET** `/api/permission-groups`

获取所有权限组。

**响应**:
```json
{
  "groups": [
    {
      "id": 1,
      "name": "string",
      "description": "string",
      "level": 1,
      "is_default": true,
      "permissions": ["permission1", "permission2"],
      "user_count": 100,
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

---

### 获取单个权限组

**GET** `/api/permission-groups/:id`

获取单个权限组详情。

---

### 创建权限组

**POST** `/api/permission-groups`

创建权限组（需管理员权限）。

**Headers**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "name": "string",
  "description": "string",
  "level": 1,
  "permissions": ["permission1", "permission2"]
}
```

---

### 更新权限组

**PUT** `/api/permission-groups/:id`

更新权限组（需管理员权限）。

**Headers**:
```
Authorization: Bearer <token>
```

---

### 删除权限组

**DELETE** `/api/permission-groups/:id`

删除权限组（需管理员权限）。

**Headers**:
```
Authorization: Bearer <token>
```

---

### 授予用户权限组

**POST** `/api/permission-groups/grant`

授予用户权限组（需管理员权限）。

**Headers**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "user_id": 1,
  "group_id": 1
}
```

---

### 撤销用户权限组

**DELETE** `/api/permission-groups/:id/revoke-user/:user_id`

撤销用户权限组（需管理员权限）。

**Headers**:
```
Authorization: Bearer <token>
```

---

### 初始化默认权限组

**POST** `/api/permission-groups/init`

初始化默认权限组（需管理员权限）。

**Headers**:
```
Authorization: Bearer <token>
```

---

### 获取用户权限组

**GET** `/api/users/:id/permission-groups`

获取指定用户的权限组。

---

### 检查用户权限

**GET** `/api/permissions/check`

检查当前用户权限。

**Headers**:
```
Authorization: Bearer <token>
```

**Query 参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| permissions | string | 权限列表，逗号分隔 |

**响应**:
```json
{
  "has_permissions": true
}
```

---

## 用户在线状态接口

### 获取用户状态

**GET** `/api/user/status/:id`

获取用户在线状态。

**响应**:
```json
{
  "user_id": 1,
  "is_online": true,
  "last_active": "2024-01-01T00:00:00Z"
}
```

---

### 更新用户状态

**POST** `/api/user/status`

更新当前用户在线状态（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

---

### 获取所有用户状态

**GET** `/api/users/status`

获取所有用户状态（需管理员权限）。

**Headers**:
```
Authorization: Bearer <token>
```

---

### 获取在线用户

**GET** `/api/users/online`

获取在线用户列表（需管理员权限）。

**Headers**:
```
Authorization: Bearer <token>
```

---

### 清理离线用户状态

**POST** `/api/users/status/cleanup`

清理离线用户状态（需管理员权限）。

**Headers**:
```
Authorization: Bearer <token>
```

---

## 系统日志接口

### 获取系统日志

**GET** `/api/system-logs`

获取系统日志（需管理员权限）。

**Headers**:
```
Authorization: Bearer <token>
```

**Query 参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| page | int | 页码，默认 1 |
| page_size | int | 每页数量，默认 20 |
| module | string | 模块筛选 |
| user_id | int | 用户ID筛选 |

---

### 获取日志模块列表

**GET** `/api/system-logs/modules`

获取日志模块列表（需管理员权限）。

**Headers**:
```
Authorization: Bearer <token>
```

---

### 删除旧日志

**DELETE** `/api/system-logs/old`

删除旧日志（需管理员权限）。

**Headers**:
```
Authorization: Bearer <token>
```

**Query 参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| days | int | 保留天数，默认30 |

---

### 获取我的操作日志

**GET** `/api/my-logs`

获取当前用户的操作日志（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

---

## 配置接口

### 获取侧边栏配置

**GET** `/api/sidebar-config`

获取侧边栏导航配置。

**响应**:
```json
{
  "items": [
    {
      "title": "首页",
      "link": "/",
      "icon": "mdi-home"
    }
  ]
}
```

---

### 更新侧边栏配置

**PUT** `/api/sidebar-config`

更新侧边栏配置（需管理员权限）。

**Headers**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "items": [
    {
      "title": "string",
      "link": "string",
      "icon": "string"
    }
  ]
}
```

---

### 获取公告

**GET** `/api/announcement`

获取网站公告。

**响应**:
```json
{
  "content": "string",
  "content_html": "string"
}
```

---

### 更新公告

**PUT** `/api/announcement`

更新网站公告（需管理员权限）。

**Headers**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "content": "string (Markdown)"
}
```

---

### 获取网站配置

**GET** `/api/site-config`

获取网站配置（公开接口）。

**响应**:
```json
{
  "site_title": "string",
  "site_description": "string",
  "icp_number": "string",
  "public_security_number": "string",
  "smtp_host": "string",
  "smtp_port": 587,
  "smtp_username": "string",
  "smtp_from": "string",
  "smtp_from_name": "string",
  "smtp_password_set": true
}
```

**说明**: SMTP密码不会直接返回，`smtp_password_set`字段表示是否已设置密码。

---

### 更新网站配置

**PUT** `/api/site-config`

更新网站配置（需管理员权限）。

**Headers**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "site_title": "string",
  "site_description": "string",
  "icp_number": "string",
  "public_security_number": "string",
  "smtp_host": "string",
  "smtp_port": 587,
  "smtp_username": "string",
  "smtp_password": "string (仅在需要修改密码时提供)",
  "smtp_from": "string",
  "smtp_from_name": "string"
}
```

**说明**: 仅在需要修改SMTP密码时提供`smtp_password`字段，为空或不提供则保持原密码不变。

---

### 获取版本信息

**GET** `/api/version`

获取前后端版本信息（公开接口）。

**响应**:
```json
{
  "frontend": {
    "version": "string"
  },
  "backend": {
    "version": "string"
  }
}
```

---

## 头衔接口

### 获取所有头衔

**GET** `/api/titles`

获取所有头衔列表。

---

### 授予用户头衔

**POST** `/api/titles/grant`

授予用户头衔（需管理员权限）。

**Headers**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "user_id": 1,
  "title_id": 1
}
```

---

### 撤销用户头衔

**POST** `/api/titles/revoke`

撤销用户头衔（需管理员权限）。

**Headers**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "user_id": 1,
  "title_id": 1
}
```

---

### 获取用户头衔

**GET** `/api/users/{id}/titles`

获取指定用户的头衔列表。

---

## 管理接口

### 获取删除请求列表

**GET** `/api/deletion-requests`

获取待处理的删除请求列表（需管理员权限）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "requests": [
    {
      "id": 1,
      "article_id": 1,
      "user": {...},
      "reason": "string",
      "status": "pending",
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

---

### 批准删除请求

**POST** `/api/deletion-requests/{id}/approve`

批准删除请求（需管理员权限）。

**Headers**:
```
Authorization: Bearer <token>
```

---

### 拒绝删除请求

**POST** `/api/deletion-requests/{id}/reject`

拒绝删除请求（需管理员权限）。

**Headers**:
```
Authorization: Bearer <token>
```

---

### 获取统计数据

**GET** `/api/admin/statistics`

获取后台统计数据（需管理员权限）。

**Headers**:
```
Authorization: Bearer <token>
```

---

### 获取所有用户

**GET** `/api/admin/users`

获取用户列表（需管理员权限）。

**Headers**:
```
Authorization: Bearer <token>
```

---

### 更新用户

**PUT** `/api/admin/users/{id}`

更新用户信息（需管理员权限）。

**Headers**:
```
Authorization: Bearer <token>
```

---

### 封禁用户

**POST** `/api/admin/users/{id}/ban`

封禁用户（需管理员权限）。

**Headers**:
```
Authorization: Bearer <token>
```

---

### 解封用户

**POST** `/api/admin/users/{id}/unban`

解封用户（需管理员权限）。

**Headers**:
```
Authorization: Bearer <token>
```

---

### 删除用户

**DELETE** `/api/admin/users/{id}`

删除用户（需管理员权限）。

**Headers**:
```
Authorization: Bearer <token>
```

---

### 更新文章状态

**PUT** `/api/admin/articles/:id/status`

更新文章状态（需管理员权限）。

**Headers**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "status": "published/hidden/deleted"
}
```

---

### 管理员删除评论

**DELETE** `/api/admin/comments/:id`

管理员删除评论（需管理员权限）。

**Headers**:
```
Authorization: Bearer <token>
```

---

## 错误响应

所有接口错误时返回以下格式：

```json
{
  "error": "错误信息描述"
}
```

**状态码**:
- `200` - 成功
- `400` - 请求参数错误
- `401` - 未授权（未登录或 Token 过期）
- `403` - 无权限（需要管理员权限）
- `404` - 资源不存在
- `500` - 服务器内部错误

---
