# 删除审核接口

## 提交删除申请（普通用户，审核流）

当普通用户删除自己的一篇文章时，并不直接删除，而是提交一条待审核的删除申请。

**DELETE** `/api/articles/{id}`

**请求体**:
```json
{
  "reason": "申请删除的原因"
}
```

**说明**:
- 需要认证，且为文章作者本人。
- 管理员（admin/system）删除时直接软删除文章，不走审核流，返回 `{ "message": "文章已删除" }`。
- 普通用户删除时创建一条 `pending` 状态的删除申请，文章保留，返回 `{ "message": "删除申请已提交，等待管理员审核" }`。
- 若该文章已存在待审核的删除申请，重复提交会返回相同提示（不重复创建）。

---

## 获取删除申请列表（管理员）

**GET** `/api/deletion-requests`

获取文章删除申请列表（需认证，管理员）。

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