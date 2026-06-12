# 密码重置接口

## 发送密码重置验证码

**POST** `/api/password/reset-code`

发送密码重置验证码到用户QQ邮箱。

**请求体**:
```json
{
  "qq_number": "string"
}
```

**响应**:
```json
{
  "message": "验证码已发送到您的QQ邮箱",
  "identifier": "string (标识token，用于后续密码重置请求)"
}
```

---

## 重置密码

**POST** `/api/password/reset`

使用验证码和标识token重置密码。

**请求体**:
```json
{
  "qq_number": "string",
  "code": "string (验证码)",
  "identifier": "string (发送验证码时返回的标识token)",
  "password": "string (新密码，至少6位)"
}
```

**响应**:
```json
{
  "message": "密码重置成功"
}
```