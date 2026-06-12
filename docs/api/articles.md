# 文章接口

## 获取文章列表

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

## 获取文章详情

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
  "comments": [...],
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

## 搜索文章

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
  "articles": [...],
  "total": 100,
  "page": 1,
  "page_size": 20,
  "total_pages": 5
}
```

---

## 创建文章

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

## 更新文章

**PUT** `/api/articles/{id}`

更新文章（需认证，仅作者可操作）。

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
  "title": "string",
  "content": "string",
  "category_id": 1,
  "voice_url": "string",
  "is_anonymous": false
}
```

**响应**:
```json
{
  "message": "更新成功",
  "article": {...}
}
```

---

## 删除文章

**DELETE** `/api/articles/{id}`

删除文章（需认证，仅作者或管理员可操作）。

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
  "message": "文章已删除"
}
```

---

## 恢复文章

**POST** `/api/articles/{id}/restore`

恢复已删除的文章（需认证，管理员）。

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
  "message": "文章已恢复"
}
```

---

## 点赞文章

**POST** `/api/articles/{id}/like`

点赞文章（需认证）。

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
  "message": "点赞成功"
}
```

---

## 取消点赞

**DELETE** `/api/articles/{id}/like`

取消点赞文章（需认证）。

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
  "message": "取消点赞成功"
}
```

---

## 分享文章

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

## 获取我的文章

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
  "articles": [...],
  "total": 20,
  "page": 1,
  "page_size": 20,
  "total_pages": 1
}
```

---

## 获取我的草稿

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

## 发布草稿

**POST** `/api/articles/{id}/publish`

发布草稿文章（需认证）。

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
  "message": "发布成功",
  "article": {...}
}
```

---

## 置顶文章

**POST** `/api/articles/{id}/pin`

置顶文章（需认证，管理员）。

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
  "message": "置顶成功",
  "article": {...}
}
```

---

## 取消置顶

**DELETE** `/api/articles/{id}/pin`

取消置顶文章（需认证，管理员）。

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
  "message": "取消置顶成功",
  "article": {...}
}
```