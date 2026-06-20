# 排行榜与徽章接口

## 获取排行榜

**GET** `/api/leaderboard`

获取各类排行榜数据（公开接口）。

**Query 参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| type | string | 排行榜类型，默认 `experience`。可选：`experience`, `articles`, `likes`, `comments`, `sign_in` |
| period | string | 统计周期，默认 `all_time`。可选：`daily`, `weekly`, `monthly`, `all_time` |
| limit | int | 返回数量，默认 50 |

**响应**:
```json
{
  "success": true,
  "data": {
    "leaderboard": [
      {
        "id": 1,
        "type": "experience",
        "period": "all_time",
        "user_id": 1,
        "user": {
          "id": 1,
          "username": "user1",
          "display_name": "用户1",
          "avatar": "https://example.com/avatar.jpg"
        },
        "score": 1250.5,
        "rank": 1,
        "date": "2024-01-01",
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 100,
    "type": "experience",
    "period": "all_time"
  }
}
```

---

## 获取用户排名

**GET** `/api/leaderboard/rank`

获取当前用户在排行榜中的排名。

**Query 参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| type | string | 排行榜类型，默认 `experience` |
| period | string | 统计周期，默认 `all_time` |

**响应**:
```json
{
  "success": true,
  "data": {
    "id": 1,
    "user_id": 1,
    "type": "experience",
    "period": "all_time",
    "score": 1250.5,
    "rank": 5,
    "date": "2024-01-01"
  }
}
```

**无排名记录时**:
```json
{
  "success": true,
  "data": {
    "rank": 0,
    "score": 0
  }
}
```

---

## 获取用户徽章

**GET** `/api/badges`

获取当前用户获得的所有徽章。

**响应**:
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "badge_type": "top_author",
      "badge_name": "优秀作者",
      "badge_icon": "mdi-trophy",
      "description": "发表文章数量排名前10",
      "is_displayed": true,
      "earned_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

---

## 更新徽章显示状态

**PUT** `/api/badges/{id}/display`

设置徽章是否在个人资料中展示。

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 徽章记录 ID |

**请求体**:
```json
{
  "is_displayed": false
}
```

**响应**:
```json
{
  "success": true,
  "data": {
    "id": 1,
    "is_displayed": false
  }
}
```

---

## 授予用户徽章（管理员）

**POST** `/api/badges/grant`

手动授予用户徽章。

**权限**: 管理员（等级 >= 80）

**请求体**:
```json
{
  "user_id": 2,
  "badge_type": "special_contributor",
  "badge_name": "特殊贡献者",
  "badge_icon": "mdi-medal",
  "description": "为社区做出特殊贡献"
}
```

**响应**:
```json
{
  "success": true,
  "data": {
    "id": 2,
    "user_id": 2,
    "badge_type": "special_contributor",
    "badge_name": "特殊贡献者",
    "badge_icon": "mdi-medal",
    "description": "为社区做出特殊贡献",
    "is_displayed": true,
    "earned_at": "2024-01-01T00:00:00Z"
  }
}
```

---

## 撤销用户徽章（管理员）

**DELETE** `/api/badges/{id}`

撤销用户的徽章。

**权限**: 管理员（等级 >= 80）

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 徽章记录 ID |

**响应**:
```json
{
  "success": true,
  "message": "撤销成功"
}
```

---

## 排行榜类型说明

| 类型 | 描述 |
|------|------|
| experience | 经验值排行 |
| articles | 文章数量排行 |
| likes | 获得点赞数排行 |
| comments | 评论数量排行 |
| sign_in | 签到天数排行 |

## 统计周期说明

| 周期 | 描述 |
|------|------|
| daily | 日榜（当日数据） |
| weekly | 周榜（本周数据） |
| monthly | 月榜（本月数据） |
| all_time | 总榜（全部历史数据） |

## 数据模型

### Leaderboard
| 字段 | 类型 | 描述 |
|------|------|------|
| type | string | 排行榜类型 |
| period | string | 统计周期 |
| user_id | int | 用户 ID |
| score | float64 | 分数 |
| rank | int | 排名 |
| date | string | 统计日期 |
| week | int | 周数 |
| month | string | 月份 YYYY-MM |

### UserBadge
| 字段 | 类型 | 描述 |
|------|------|------|
| badge_type | string | 徽章类型 |
| badge_name | string | 徽章名称 |
| badge_icon | string | 徽章图标 |
| is_displayed | bool | 是否展示 |
| earned_at | datetime | 获得时间 |
