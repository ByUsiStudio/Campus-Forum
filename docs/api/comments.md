# 评论接口

## 获取文章评论列表

**GET** `/api/articles/{id}`

获取文章评论列表，包含嵌套回复（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 文章 ID |

**查询参数**:
| 参数 | 类型 | 默认值 | 描述 |
|------|------|--------|------|
| page | int | 1 | 页码 |
| page_size | int | 20 | 每页评论数（不含回复） |

**响应**:
```json
{
  "comments": [
    {
      "id": 1,
      "content": "评论内容",
      "user_id": 1,
      "user": {
        "id": 1,
        "username": "user1",
        "display_name": "显示名称",
        "avatar": "头像URL"
      },
      "article_id": 1,
      "parent_id": null,
      "like_count": 5,
      "reply_count": 3,
      "is_anonymous": false,
      "created_at": "2024-01-01T00:00:00Z",
      "replies": [
        {
          "id": 2,
          "content": "回复内容",
          "user_id": 2,
          "user": {...},
          "parent_id": 1,
          "like_count": 1,
          "reply_count": 0,
          "is_anonymous": false,
          "created_at": "2024-01-01T01:00:00Z",
          "replies": []
        }
      ]
    }
  ],
  "comment_liked": {
    "1": true,
    "2": false
  },
  "total": 10,
  "page": 1,
  "page_size": 20
}
```

**嵌套回复说明**:
- 评论支持多级嵌套，**最大嵌套深度为 5 层**
- 回复列表通过 `replies` 字段返回，**仅顶级评论**包含 `replies`
- 每条评论都有 `reply_count` 字段表示直接子回复数量
- 匿名评论的 `user` 字段会被替换为 `{ id: 0, username: "anonymous", display_name: "匿名用户", avatar: "" }`

---

## 创建评论

**POST** `/api/articles/{id}/comments`

创建评论（需认证）。

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
  "content": "string",
  "parent_id": 1 (可选，回复评论时使用，指定父评论ID),
  "is_anonymous": false
}
```

**嵌套回复规则**:
- `parent_id` 为空时创建顶级评论
- `parent_id` 有值时创建该评论的子回复
- 嵌套层级限制：**最多 5 层**，超过后 `parent_id` 将被忽略
- 回复时 `reply_count` 会自动累加到父评论

**响应**:
```json
{
  "message": "评论成功",
  "comment": {...}
}
```

---

## 删除评论

**DELETE** `/api/comments/{id}`

删除评论（需认证，仅作者或管理员可操作）。

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

## 点赞评论

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

## 取消点赞评论

**DELETE** `/api/comments/{id}/like`

取消点赞评论（需认证）。

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