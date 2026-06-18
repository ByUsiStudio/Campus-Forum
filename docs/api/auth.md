# 认证接口

## Token 说明

系统使用双令牌机制：
- **访问令牌 (access token)**：有效期 1 小时，用于 API 认证
- **刷新令牌 (refresh token)**：有效期 7 天，用于刷新访问令牌

前端会在访问令牌即将过期时（5 分钟内）自动刷新令牌。

---

## 用户注册

**POST** `/api/auth/register`

注册新用户账号。

**请求体**:
```json
{
  "username": "string",
  "qq_number": "string",
  "display_name": "string",
  "password": "string"
}
```

**响应**:
```json
{
  "message": "注册成功"
}
```

---

## 用户登录

**POST** `/api/auth/login`

用户登录获取 Token。

**请求体**:
```json
{
  "username": "string",
  "password": "string"
}
```

**响应**:
```json
{
  "token": "access_token_string",
  "refresh_token": "refresh_token_string",
  "expires_in": 3600,
  "user": {
    "id": 1,
    "username": "string",
    "display_name": "string",
    "avatar": "string",
    "role": "user",
    "status": "active"
  }
}
```

**字段说明**:
| 字段 | 类型 | 描述 |
|------|------|------|
| token | string | 访问令牌，有效期 1 小时 |
| refresh_token | string | 刷新令牌，有效期 7 天 |
| expires_in | int | 访问令牌过期时间（秒） |
| user | object | 用户信息 |

---

## 刷新令牌

**POST** `/api/auth/refresh-token`

使用刷新令牌获取新的访问令牌。

**请求体**:
```json
{
  "refresh_token": "refresh_token_string"
}
```

**响应**:
```json
{
  "token": "new_access_token_string",
  "refresh_token": "new_refresh_token_string",
  "expires_in": 3600
}
```

**说明**：
- 刷新令牌时会同时生成新的访问令牌和刷新令牌
- 原刷新令牌失效，需使用新的刷新令牌
- 如果刷新令牌无效或已过期，返回 401 错误

---

## 初始化管理员

**POST** `/api/auth/init-admin`

初始化管理员账号（仅在系统未初始化时可用）。

**请求体**:
```json
{
  "username": "string",
  "password": "string",
  "qq_number": "string"
}
```

**响应**:
```json
{
  "message": "管理员初始化成功"
}
```

---

## 检查初始化状态

**GET** `/api/auth/check-init`

检查系统是否已初始化。

**响应**:
```json
{
  "initialized": true
}
```