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
  "total_sign_ins": 12,
  "reward_coins": 3,
  "bonus_coins": 0,
  "max_continuous_days": 30,
  "total_coins": 100
}
```

**响应说明**:
| 字段 | 类型 | 描述 |
|------|------|------|
| message | string | 签到结果信息 |
| sign_in_days | int | 连续签到天数 |
| total_sign_ins | int | 累计签到次数 |
| reward_coins | int | 当次签到获得的币 |
| bonus_coins | int | 额外奖励币（连续7天+5，连续30天+15，连续365天+50） |
| max_continuous_days | int | 历史最长连续签到天数 |
| total_coins | int | 用户累计总币数 |

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
  "total_sign_ins": 12,
  "max_continuous_days": 30,
  "total_coins": 100,
  "month_sign_in_count": 3,
  "week_sign_in_count": 2,
  "config": {
    "daily_coins": 1,
    "weekly_bonus": 7,
    "monthly_bonus": 30,
    "yearly_bonus": 365
  }
}
```

**响应说明**:
| 字段 | 类型 | 描述 |
|------|------|------|
| has_signed_in | bool | 今日是否已签到 |
| sign_in_days | int | 当前连续签到天数 |
| total_sign_ins | int | 累计签到次数 |
| max_continuous_days | int | 历史最长连续签到天数 |
| total_coins | int | 用户累计总币数 |
| month_sign_in_count | int | 本月签到天数 |
| week_sign_in_count | int | 本周签到天数 |
| config | object | 签到配置信息 |
| config.daily_coins | int | 每日基础币数 |
| config.weekly_bonus | int | 连续7天奖励触发天数 |
| config.monthly_bonus | int | 连续30天奖励触发天数 |
| config.yearly_bonus | int | 连续365天奖励触发天数 |

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
  "records": [
    {
      "id": 1,
      "user_id": 10,
      "sign_in_at": "2024-01-01T12:00:00Z",
      "sign_in_date": "2024-01-01",
      "continuous_day": 5,
      "reward_points": 3,
      "ip_address": "127.0.0.1"
    }
  ],
  "total": 100,
  "page": 1,
  "page_size": 30,
  "total_pages": 4
}
```

**签到记录字段说明**:
| 字段 | 类型 | 描述 |
|------|------|------|
| id | int | 记录ID |
| user_id | int | 用户ID |
| sign_in_at | string | 答到时间（ISO格式） |
| sign_in_date | string | 答到日期（YYYY-MM-DD格式） |
| continuous_day | int | 当次答到后的连续天数 |
| reward_points | int | 当次答到获得的币数 |
| ip_address | string | 答到IP |

---

## 获取签到排行榜

**GET** `/api/signin/rankings`

获取签到排行榜（需认证）。

**Headers**:
```
Authorization: Bearer <token>
```

**Query 参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| limit | int | 返回数量，默认 10，最大 100 |

**响应**:
```json
{
  "continuous_rankings": [
    {
      "id": 1,
      "username": "user1",
      "display_name": "用户1",
      "avatar": "...",
      "sign_in_days": 30,
      "total_sign_ins": 50,
      "total_coins": 100
    }
  ],
  "points_rankings": [
    {
      "id": 2,
      "username": "user2",
      "display_name": "用户2",
      "avatar": "...",
      "sign_in_days": 10,
      "total_sign_ins": 30,
      "total_coins": 200
    }
  ]
}
```

**排行榜说明**:
| 排行榜 | 排序规则 |
|--------|----------|
| continuous_rankings | 按连续签到天数降序，相同则按累计签到次数降序 |
| points_rankings | 按累计币数降序，相同则按累计签到次数降序 |

---

## 获取签到配置

**GET** `/api/signin/config`

获取签到配置信息。

**Headers**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "config": {
    "id": 1,
    "daily_points": 1,
    "weekly_bonus": 7,
    "monthly_bonus": 30,
    "yearly_bonus": 365,
    "enabled": true
  }
}
```

---

## 更新签到配置（管理员）

**PUT** `/api/signin/config`

更新签到配置（需认证，管理员）。

**Headers**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "daily_points": 2,
  "weekly_bonus": 7,
  "monthly_bonus": 30,
  "yearly_bonus": 365,
  "enabled": true
}
```

**请求参数说明**:
| 参数 | 类型 | 描述 |
|------|------|------|
| daily_points | int | 每日签到基础币数（选填，大于0） |
| weekly_bonus | int | 连续签到触发周奖励的天数（选填，大于0） |
| monthly_bonus | int | 连续签到触发月奖励的天数（选填，大于0） |
| yearly_bonus | int | 连续签到触发年奖励的天数（选填，大于0） |
| enabled | bool | 是否启用签到功能（选填） |

**响应**:
```json
{
  "message": "更新成功",
  "config": {
    "id": 1,
    "daily_points": 2,
    "weekly_bonus": 7,
    "monthly_bonus": 30,
    "yearly_bonus": 365,
    "enabled": true
  }
}
```

---

## 签到奖励机制

### 币奖励规则

| 签到类型 | 奖励币数 |
|----------|----------|
| 每日基础签到 | `daily_points`（默认1币） |
| 连续7天签到 | 额外奖励5币 |
| 连续30天签到 | 额外奖励15币 |
| 连续365天签到 | 额外奖励50币 |

### 示例

用户连续签到第7天时：
- 基础币数：1币
- 周奖励：5币
- 当次总计：6币
- 连续签到天数：7天

---

## 签到记录数据表结构

### sign_in_records 表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| user_id | uint | 用户ID（索引） |
| sign_in_at | datetime | 答到时间（索引） |
| sign_in_date | string | 答到日期 YYYY-MM-DD（索引） |
| continuous_day | int | 当次答到后的连续天数 |
| reward_points | int | 当次答到获得的币数 |
| ip_address | string | 答到IP地址 |

### sign_in_config 表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| daily_points | int | 每日基础币数 |
| weekly_bonus | int | 连续答到触发周奖励的天数 |
| monthly_bonus | int | 连续答到触发月奖励的天数 |
| yearly_bonus | int | 连续答到触发年奖励的天数 |
| enabled | bool | 是否启用 |

### users 表新增字段

| 字段 | 类型 | 说明 |
|------|------|------|
| sign_in_days | int | 当前连续答到天数 |
| total_sign_ins | int | 累计答到次数 |
| max_continuous_days | int | 历史最长连续答到天数 |
| total_coins | int | 累计币数 |
