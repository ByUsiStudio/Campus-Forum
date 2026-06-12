# 认证接口

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
  "token": "jwt_token_string",
  "user": {
    "id": 1,
    "username": "string",
    "display_name": "string",
    "avatar": "string",
    "role": "user"
  }
}
```

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