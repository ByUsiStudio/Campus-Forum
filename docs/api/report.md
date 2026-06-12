# 举报接口

## 创建举报

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
  "message": "举报成功"
}
```

---

## 获取举报列表（管理员）

**GET** `/api/reports`

获取举报列表（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "reports": [...]
}
```

---

## 获取举报详情（管理员）

**GET** `/api/reports/{id}`

获取举报详情（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 举报 ID |

**响应**:
```json
{
  "report": {...}
}
```

---

## 处理举报（管理员）

**PUT** `/api/reports/{id}/handle`

处理举报（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 举报 ID |

**请求体**:
```json
{
  "status": "resolved",
  "result": "string",
  "action": "string"
}
```

**响应**:
```json
{
  "message": "处理成功"
}
```