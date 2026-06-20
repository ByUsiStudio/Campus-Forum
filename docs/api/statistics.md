# 数据统计与分析接口

## 获取用户个人统计

**GET** `/api/statistics`

获取当前登录用户的个人统计数据。

**响应**:
```json
{
  "success": true,
  "data": {
    "id": 1,
    "user_id": 1,
    "total_articles": 10,
    "total_comments": 50,
    "total_likes": 200,
    "total_views": 1500,
    "total_shares": 5,
    "total_favorites": 30,
    "total_followers": 20,
    "total_following": 15,
    "active_days": 30,
    "last_active_date": "2024-01-01T00:00:00Z",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

## 获取每日统计数据

**GET** `/api/statistics/daily`

获取系统每日统计数据，默认返回最近30天。

**Query 参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| start_date | string | 开始日期，格式 YYYY-MM-DD，默认30天前 |
| end_date | string | 结束日期，格式 YYYY-MM-DD，默认今天 |

**响应**:
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "date": "2024-01-01",
      "new_users": 5,
      "active_users": 120,
      "new_articles": 10,
      "new_comments": 50,
      "total_views": 800,
      "total_likes": 200,
      "total_shares": 10,
      "total_sign_ins": 80,
      "peak_online": 50
    }
  ]
}
```

---

## 获取系统概览

**GET** `/api/statistics/overview`

获取系统整体概览数据（需要登录）。

**响应**:
```json
{
  "success": true,
  "data": {
    "id": 1,
    "total_users": 500,
    "total_articles": 1000,
    "total_comments": 5000,
    "total_categories": 10,
    "online_users": 30,
    "today_active_users": 120,
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

---

## 获取用户活跃度

**GET** `/api/statistics/activity`

获取当前用户的活跃度历史数据。

**Query 参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| start_date | string | 开始日期，格式 YYYY-MM-DD，默认30天前 |
| end_date | string | 结束日期，格式 YYYY-MM-DD，默认今天 |

**响应**:
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "date": "2024-01-01",
      "login_count": 3,
      "post_count": 2,
      "comment_count": 5,
      "like_count": 10,
      "view_count": 50,
      "active_score": 25.5
    }
  ]
}
```

---

## 获取文章统计

**GET** `/api/articles/{id}/statistics`

获取指定文章的统计数据。

**路径参数**:
| 参数 | 类型 | 描述 |
|------|------|------|
| id | int | 文章 ID |

**响应**:
```json
{
  "success": true,
  "data": {
    "id": 1,
    "article_id": 5,
    "daily_views": 50,
    "daily_likes": 10,
    "daily_comments": 5,
    "daily_shares": 2,
    "weekly_views": 300,
    "monthly_views": 1200,
    "last_reset_at": "2024-01-01T00:00:00Z"
  }
}
```

---

## 获取统计仪表板（管理员）

**GET** `/api/statistics/dashboard`

获取管理员统计仪表板数据，包含系统概览、最近7天统计、热门文章和活跃用户。

**权限**: 管理员（等级 >= 80）

**响应**:
```json
{
  "success": true,
  "data": {
    "overview": {
      "total_users": 500,
      "total_articles": 1000,
      "total_comments": 5000,
      "online_users": 30
    },
    "recent_stats": [
      {
        "date": "2024-01-01",
        "new_users": 5,
        "active_users": 120,
        "new_articles": 10,
        "new_comments": 50
      }
    ],
    "hot_articles": [
      {
        "id": 1,
        "title": "热门文章标题",
        "view_count": 1000,
        "like_count": 100,
        "comment_count": 50
      }
    ],
    "active_users": [
      {
        "id": 1,
        "username": "user1",
        "display_name": "用户1",
        "avatar": "https://example.com/avatar.jpg",
        "online_status": "online"
      }
    ]
  }
}
```

---

## 数据模型

### DailyStatistics
| 字段 | 类型 | 描述 |
|------|------|------|
| date | string | 日期 YYYY-MM-DD |
| new_users | int | 新增用户数 |
| active_users | int | 活跃用户数 |
| new_articles | int | 新增文章数 |
| new_comments | int | 新增评论数 |
| total_views | int | 总浏览量 |
| total_likes | int | 总点赞数 |
| total_shares | int | 总分享数 |
| total_sign_ins | int | 总签到数 |
| peak_online | int | 峰值在线人数 |

### ArticleStatistics
| 字段 | 类型 | 描述 |
|------|------|------|
| daily_views | int | 当日浏览量 |
| weekly_views | int | 本周浏览量 |
| monthly_views | int | 本月浏览量 |
