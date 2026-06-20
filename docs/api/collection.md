# 收藏夹管理接口

## 获取用户收藏夹列表

**GET** `/api/collections`

获取当前登录用户的所有收藏夹。

**响应**:
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "name": "技术文章",
      "description": "收藏的技术相关文章",
      "is_public": false,
      "cover_image": "https://example.com/cover.jpg",
      "article_count": 10,
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

---

## 获取收藏夹详情

**GET** `/api/collections/{id}`

获取指定收藏夹的详细信息和其中的文章列表。

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 收藏夹 ID |

**响应**:
```json
{
  "success": true,
  "data": {
    "collection": {
      "id": 1,
      "name": "技术文章",
      "description": "收藏的技术相关文章",
      "is_public": false,
      "article_count": 2
    },
    "articles": [
      {
        "id": 1,
        "collection_id": 1,
        "article_id": 5,
        "article": {
          "id": 5,
          "title": "文章标题",
          "user": {
            "id": 2,
            "username": "author",
            "display_name": "作者"
          }
        },
        "note": "这篇写得很好",
        "created_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

---

## 创建收藏夹

**POST** `/api/collections`

创建新的收藏夹。

**请求体**:
```json
{
  "name": "学习资料",
  "description": "学习相关的文章收藏",
  "is_public": false,
  "cover_image": "https://example.com/cover.jpg"
}
```

**响应**:
```json
{
  "success": true,
  "data": {
    "id": 2,
    "user_id": 1,
    "name": "学习资料",
    "description": "学习相关的文章收藏",
    "is_public": false,
    "cover_image": "https://example.com/cover.jpg",
    "article_count": 0,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

## 更新收藏夹

**PUT** `/api/collections/{id}`

更新收藏夹信息。

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 收藏夹 ID |

**请求体**:
```json
{
  "name": "更新后的名称",
  "description": "更新后的描述",
  "is_public": true
}
```

**响应**:
```json
{
  "success": true,
  "data": {...}
}
```

---

## 删除收藏夹

**DELETE** `/api/collections/{id}`

删除收藏夹及其中的所有文章关联。

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 收藏夹 ID |

**响应**:
```json
{
  "success": true,
  "message": "删除成功"
}
```

---

## 添加文章到收藏夹

**POST** `/api/collections/{id}/articles`

将文章添加到指定收藏夹。

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 收藏夹 ID |

**请求体**:
```json
{
  "article_id": 5,
  "note": "这篇写得很好，值得反复阅读"
}
```

**响应**:
```json
{
  "success": true,
  "message": "添加成功"
}
```

**错误响应**:
- `400` - 文章已在收藏夹中
- `404` - 收藏夹不存在或文章不存在

---

## 从收藏夹移除文章

**DELETE** `/api/collections/{id}/articles/{article_id}`

从收藏夹中移除指定文章。

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 收藏夹 ID |
| article_id | int | 文章 ID |

**响应**:
```json
{
  "success": true,
  "message": "移除成功"
}
```

---

## 获取文章版本历史

**GET** `/api/articles/{id}/versions`

获取指定文章的所有版本历史。

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 文章 ID |

**响应**:
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "article_id": 5,
      "title": "文章标题 v1",
      "version": 1,
      "change_log": "初始版本",
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

---

## 获取指定版本内容

**GET** `/api/articles/{id}/versions/{version}`

获取文章指定版本的完整内容。

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 文章 ID |
| version | int | 版本号 |

**响应**:
```json
{
  "success": true,
  "data": {
    "id": 1,
    "article_id": 5,
    "title": "文章标题 v1",
    "content": "文章内容（Markdown）",
    "content_html": "<p>文章内容（HTML）</p>",
    "version": 1,
    "change_log": "初始版本",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

## 恢复文章版本

**POST** `/api/articles/{id}/versions/{version}/restore`

将文章恢复到指定版本。

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 文章 ID |
| version | int | 版本号 |

**响应**:
```json
{
  "success": true,
  "message": "恢复成功"
}
```

---

## 数据模型

### Collection
| 字段 | 类型 | 描述 |
|------|------|------|
| id | int | 收藏夹 ID |
| name | string | 收藏夹名称 |
| description | string | 收藏夹描述 |
| is_public | bool | 是否公开 |
| cover_image | string | 封面图片 URL |
| article_count | int | 文章数量 |

### CollectionArticle
| 字段 | 类型 | 描述 |
|------|------|------|
| collection_id | int | 收藏夹 ID |
| article_id | int | 文章 ID |
| note | string | 收藏笔记 |
