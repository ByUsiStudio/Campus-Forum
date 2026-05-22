# 校园论坛 API 文档

## 概述

本文档描述了校园论坛后端 API 的详细接口说明。

**Base URL**: `http://localhost:3620`

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
  "password": "string",
  "email": "string"
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

获取单篇文章详情。

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 文章 ID |

**响应**:
```json
{
  "id": 1,
  "title": "string",
  "content": "string",
  "content_html": "string",
  "user": {...},
  "category": {...},
  "view_count": 100,
  "like_count": 10,
  "comments": [...],
  "created_at": "2024-01-01T00:00:00Z"
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
  "category_id": 1
}
```

**响应**:
```json
{
  "message": "文章创建成功",
  "article_id": 1
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
  "category_id": 1
}
```

---

### 删除文章

**DELETE** `/api/articles/{id}`

删除文章（需认证，仅作者或管理员）。

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

## 评论接口

### 获取文章评论

**GET** `/api/articles/{id}/comments`

获取文章的评论列表。

**响应**:
```json
{
  "comments": [
    {
      "id": 1,
      "content": "string",
      "user": {...},
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

**请求体**:
```json
{
  "content": "string"
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
  "bio": "string",
  "role": "user",
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
  "bio": "string"
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
| file | file | 图片文件 |

**响应**:
```json
{
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
| file | file | 图片文件 |

**响应**:
```json
{
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
| file | file | 视频文件 |

**响应**:
```json
{
  "url": "https://..."
}
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

## 交互式文档

启动服务后，可访问以下地址获取交互式 API 文档：

**Swagger UI**: http://localhost:3620/swagger/index.html
