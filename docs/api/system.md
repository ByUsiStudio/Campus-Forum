# 系统接口

## 获取公告

**GET** `/api/announcement`

获取系统公告。

**响应**:
```json
{
  "id": 1,
  "content": "string",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

---

## 更新公告（管理员）

**PUT** `/api/announcement`

更新系统公告（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "content": "string"
}
```

**响应**:
```json
{
  "message": "更新成功"
}
```

---

## 获取侧边栏配置

**GET** `/api/sidebar-config`

获取侧边栏配置。

**响应**:
```json
{
  "id": 1,
  "content": {...},
  "updated_at": "2024-01-01T00:00:00Z"
}
```

---

## 更新侧边栏配置（管理员）

**PUT** `/api/sidebar-config`

更新侧边栏配置（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "content": {...}
}
```

**响应**:
```json
{
  "message": "更新成功"
}
```

---

## 获取网站配置

**GET** `/api/site-config`

获取网站配置。

**响应**:
```json
{
  "id": 1,
  "site_name": "string",
  "site_description": "string",
  "smtp_enabled": true,
  "updated_at": "2024-01-01T00:00:00Z"
}
```

---

## 更新网站配置（管理员）

**PUT** `/api/site-config`

更新网站配置（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "site_name": "string",
  "site_description": "string",
  "smtp_enabled": true,
  "smtp_host": "string",
  "smtp_port": 587,
  "smtp_username": "string",
  "smtp_password": "string"
}
```

**响应**:
```json
{
  "message": "更新成功"
}
```

---

## 测试SMTP配置（管理员）

**POST** `/api/site-config/test-smtp`

测试SMTP配置（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "message": "测试成功"
}
```

---

## 获取系统日志（管理员）

**GET** `/api/system-logs`

获取系统日志（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "logs": [...]
}
```

---

## 获取日志模块列表（管理员）

**GET** `/api/system-logs/modules`

获取日志模块列表（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "modules": [...]
}
```

---

## 删除旧日志（管理员）

**DELETE** `/api/system-logs/old`

删除旧日志（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "message": "删除成功"
}
```

---

## 获取我的操作日志

**GET** `/api/my-logs`

获取当前用户的操作日志（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "logs": [...]
}
```

---

## 获取版本信息

**GET** `/api/version`

获取系统版本信息。

**响应**:
```json
{
  "version": "1.0.0"
}
```