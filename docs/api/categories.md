# 分类接口

## 获取分类列表

**GET** `/api/categories`

获取所有分类列表。

**响应**:
```json
{
  "categories": [
    {
      "id": 1,
      "name": "string",
      "description": "string",
      "icon": "string",
      "sort_order": 1,
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

---

## 创建分类（管理员）

**POST** `/api/categories`

创建新分类（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "name": "string",
  "description": "string",
  "icon": "string",
  "sort_order": 1
}
```

**响应**:
```json
{
  "message": "创建成功",
  "category": {...}
}
```

---

## 更新分类（管理员）

**PUT** `/api/categories/{id}`

更新分类信息（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 分类 ID |

**请求体**:
```json
{
  "name": "string",
  "description": "string",
  "icon": "string",
  "sort_order": 1
}
```

**响应**:
```json
{
  "message": "更新成功",
  "category": {...}
}
```

---

## 删除分类（管理员）

**DELETE** `/api/categories/{id}`

删除分类（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 分类 ID |

**响应**:
```json
{
  "message": "删除成功"
}
```