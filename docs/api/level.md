# 用户等级与成就接口

## 获取用户等级信息

**GET** `/api/level`

获取当前登录用户的等级信息。如果没有等级记录，会自动创建初始等级（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "success": true,
  "data": {
    "id": 1,
    "user_id": 1,
    "level": 1,
    "experience": 0,
    "next_level": 100,
    "title": "新手",
    "badge": "",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

---

## 获取用户经验记录

**GET** `/api/level/experience-records`

获取当前用户的经验获取记录（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**Query 参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| page | int | 页码，默认 1 |
| limit | int | 每页数量，默认 20 |

**响应**:
```json
{
  "success": true,
  "data": {
    "records": [
      {
        "id": 1,
        "user_id": 1,
        "type": "post",
        "amount": 10,
        "description": "发布文章",
        "related_id": 5,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "limit": 20
  }
}
```

---

## 获取用户成就列表

**GET** `/api/achievements`

获取当前用户已获得的所有成就（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "achievement_id": 2,
      "achievement": {
        "id": 2,
        "code": "first_post",
        "name": "初次发帖",
        "description": "发布第一篇文章",
        "icon": "mdi-pencil",
        "category": "post",
        "reward": 50,
        "rarity": "common"
      },
      "progress": 100,
      "unlocked_at": "2024-01-01T00:00:00Z",
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

---

## 获取所有成就定义

**GET** `/api/achievements/all`

获取系统中所有成就的定义列表（公开接口）。

**响应**:
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "code": "first_post",
      "name": "初次发帖",
      "description": "发布第一篇文章",
      "icon": "mdi-pencil",
      "category": "post",
      "condition": "{\"min_posts\": 1}",
      "reward": 50,
      "rarity": "common",
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

---

## 获取等级配置

**GET** `/api/level/config`

获取所有等级配置规则（公开接口）。

**响应**:
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "level": 1,
      "min_exp": 0,
      "title": "新手",
      "badge": "mdi-star-outline",
      "privileges": "{\"daily_signin_bonus\": 1}",
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

---

## 创建等级配置（管理员）

**POST** `/api/level/config`

创建新的等级配置规则（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "level": 2,
  "min_exp": 100,
  "title": "学徒",
  "badge": "mdi-star-half-full",
  "privileges": "{\"daily_signin_bonus\": 2}"
}
```

**响应**:
```json
{
  "success": true,
  "data": {
    "id": 2,
    "level": 2,
    "min_exp": 100,
    "title": "学徒",
    "badge": "mdi-star-half-full",
    "privileges": "{\"daily_signin_bonus\": 2}"
  }
}
```

---

## 更新等级配置（管理员）

**PUT** `/api/level/config/{id}`

更新指定等级配置（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 配置 ID |

**请求体**: 同创建等级配置

**响应**:
```json
{
  "success": true,
  "data": {...}
}
```

---

## 创建成就（管理员）

**POST** `/api/achievements`

创建新的成就定义（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "code": "active_commenter",
  "name": "活跃评论者",
  "description": "累计发布100条评论",
  "icon": "mdi-comment-multiple",
  "category": "comment",
  "condition": "{\"min_comments\": 100}",
  "reward": 200,
  "rarity": "rare"
}
```

**响应**:
```json
{
  "success": true,
  "data": {...}
}
```

---

## 更新成就（管理员）

**PUT** `/api/achievements/{id}`

更新成就定义（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 成就 ID |

**请求体**: 同创建成就

**响应**:
```json
{
  "success": true,
  "data": {...}
}
```

---

## 删除成就（管理员）

**DELETE** `/api/achievements/{id}`

删除成就定义（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 成就 ID |

**响应**:
```json
{
  "success": true,
  "message": "删除成功"
}
```

---

## 数据模型

### UserLevel
| 字段 | 类型 | 描述 |
|------|------|------|
| id | int | 记录 ID |
| user_id | int | 用户 ID |
| level | int | 当前等级 |
| experience | int | 当前经验值 |
| next_level | int | 升级所需经验 |
| title | string | 等级称号 |
| badge | string | 等级徽章图标 |

### Achievement
| 字段 | 类型 | 描述 |
|------|------|------|
| code | string | 成就唯一代码 |
| name | string | 成就名称 |
| category | string | 分类：post, comment, social, special |
| rarity | string | 稀有度：common, rare, epic, legendary |
| reward | int | 奖励经验值 |

### ExperienceRecord 类型说明
| 类型 | 描述 |
|------|------|
| login | 每日登录 |
| post | 发布文章 |
| comment | 发表评论 |
| like | 获得点赞 |
| share | 分享文章 |
| achievement | 获得成就 |
