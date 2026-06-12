# 签到接口

## 用户签到

**POST** `/api/signin`

用户签到（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "message": "签到成功",
  "sign_in_days": 5,
  "total_sign_ins": 12
}
```

---

## 获取签到状态

**GET** `/api/signin/status`

获取用户签到状态（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "has_signed_in": false,
  "sign_in_days": 5,
  "total_sign_ins": 12
}
```

---

## 获取签到历史

**GET** `/api/signin/history`

获取用户签到历史（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**Query 参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| page | int | 页码，默认 1 |
| page_size | int | 每页数量，默认 30 |

**响应**:
```json
{
  "records": [...],
  "total": 100,
  "page": 1,
  "page_size": 30,
  "total_pages": 4
}
```