# 话题标签接口

## 获取话题列表

**GET** `/api/topics`

获取所有话题标签列表（公开接口）。

**Query 参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| page | int | 页码，默认 1 |
| limit | int | 每页数量，默认 20 |
| is_hot | bool | 是否只返回热门话题 |

**响应**:
```json
{
  "success": true,
  "data": {
    "topics": [
      {
        "id": 1,
        "name": "golang",
        "display_name": "Go语言",
        "description": "Go语言相关讨论",
        "icon": "mdi-language-go",
        "cover_image": "https://example.com/go.jpg",
        "article_count": 50,
        "follow_count": 100,
        "is_hot": true,
        "is_official": true,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 20,
    "page": 1,
    "limit": 20
  }
}
```

---

## 获取话题详情

**GET** `/api/topics/{id}`

获取指定话题的详细信息和关联文章列表。

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 话题 ID |

**Query 参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| page | int | 页码，默认 1 |
| limit | int | 每页数量，默认 20 |

**响应**:
```json
{
  "success": true,
  "data": {
    "topic": {
      "id": 1,
      "name": "golang",
      "display_name": "Go语言",
      "description": "Go语言相关讨论",
      "article_count": 50,
      "follow_count": 100
    },
    "articles": [
      {
        "id": 1,
        "article_id": 5,
        "article": {
          "id": 5,
          "title": "Go语言入门教程",
          "user": {
            "id": 1,
            "username": "author"
          },
          "category": {
            "id": 1,
            "name": "技术"
          }
        }
      }
    ],
    "total": 50,
    "page": 1,
    "limit": 20
  }
}
```

---

## 获取热门话题

**GET** `/api/topics/hot`

获取热门话题列表（公开接口）。

**响应**:
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "golang",
      "display_name": "Go语言",
      "article_count": 50,
      "follow_count": 100,
      "is_hot": true
    }
  ]
}
```

---

## 创建话题（管理员）

**POST** `/api/topics`

创建新的话题标签。

**权限**: 管理员（等级 >= 80）

**请求体**:
```json
{
  "name": "vuejs",
  "display_name": "Vue.js",
  "description": "Vue.js 前端框架讨论",
  "icon": "mdi-vuejs",
  "cover_image": "https://example.com/vue.jpg",
  "is_hot": false,
  "is_official": true
}
```

**响应**:
```json
{
  "success": true,
  "data": {
    "id": 2,
    "name": "vuejs",
    "display_name": "Vue.js",
    "description": "Vue.js 前端框架讨论",
    "icon": "mdi-vuejs",
    "article_count": 0,
    "follow_count": 0,
    "is_hot": false,
    "is_official": true
  }
}
```

---

## 更新话题（管理员）

**PUT** `/api/topics/{id}`

更新话题信息。

**权限**: 管理员（等级 >= 80）

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 话题 ID |

**请求体**: 同创建话题

**响应**:
```json
{
  "success": true,
  "data": {...}
}
```

---

## 删除话题（管理员）

**DELETE** `/api/topics/{id}`

删除话题及其所有关联。

**权限**: 管理员（等级 >= 80）

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 话题 ID |

**响应**:
```json
{
  "success": true,
  "message": "删除成功"
}
```

---

## 关注话题

**POST** `/api/topics/{id}/follow`

关注指定话题。

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 话题 ID |

**响应**:
```json
{
  "success": true,
  "message": "关注成功"
}
```

**错误响应**:
- `400` - 已关注该话题

---

## 取消关注话题

**DELETE** `/api/topics/{id}/follow`

取消关注指定话题。

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 话题 ID |

**响应**:
```json
{
  "success": true,
  "message": "取消关注成功"
}
```

---

## 获取用户关注的话题

**GET** `/api/topics/followed`

获取当前用户关注的所有话题。

**响应**:
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "topic_id": 1,
      "topic": {
        "id": 1,
        "name": "golang",
        "display_name": "Go语言",
        "article_count": 50
      },
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

---

## 为文章添加话题

**POST** `/api/articles/{id}/topics`

为当前用户的文章添加话题标签。

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 文章 ID |

**请求体**:
```json
{
  "topic_id": 1
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
- `400` - 文章已添加该话题
- `404` - 文章不存在或无权限

---

## 从文章移除话题

**DELETE** `/api/articles/{id}/topics/{topic_id}`

从文章中移除话题标签。

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 文章 ID |
| topic_id | int | 话题 ID |

**响应**:
```json
{
  "success": true,
  "message": "移除成功"
}
```

---

## 数据模型

### Topic
| 字段 | 类型 | 描述 |
|------|------|------|
| id | int | 话题 ID |
| name | string | 话题唯一名称（英文标识） |
| display_name | string | 显示名称 |
| description | string | 话题描述 |
| icon | string | 话题图标 |
| article_count | int | 关联文章数 |
| follow_count | int | 关注数 |
| is_hot | bool | 是否热门 |
| is_official | bool | 是否官方 |
