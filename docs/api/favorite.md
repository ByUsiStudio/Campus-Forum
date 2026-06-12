# 收藏接口

## 添加收藏

**POST** `/api/articles/{id}/favorite`

收藏文章（需认证）。

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
  "message": "收藏成功"
}
```

---

## 取消收藏

**DELETE** `/api/articles/{id}/favorite`

取消收藏文章（需认证）。

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
  "message": "取消收藏成功"
}
```

---

## 获取收藏列表

**GET** `/api/favorites`

获取当前用户的收藏列表（需认证）。

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

## 检查收藏状态

**GET** `/api/articles/{id}/favorite/check`

检查是否收藏了指定文章（需认证）。

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
  "favorited": true
}
```