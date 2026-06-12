# 评论接口

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
  "parent_id": 1 (可选，回复评论时使用),
  "is_anonymous": false
}
```

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