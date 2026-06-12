# 删除审核接口

## 获取删除申请列表（管理员）

**GET** `/api/deletion-requests`

获取账号删除申请列表（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "requests": [...]
}
```

---

## 批准删除申请（管理员）

**POST** `/api/deletion-requests/{id}/approve`

批准账号删除申请（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 删除申请 ID |

**响应**:
```json
{
  "message": "批准成功"
}
```

---

## 拒绝删除申请（管理员）

**POST** `/api/deletion-requests/{id}/reject`

拒绝账号删除申请（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 删除申请 ID |

**响应**:
```json
{
  "message": "拒绝成功"
}
```