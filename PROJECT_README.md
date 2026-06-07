# Campus Forum 校园论坛

## 项目架构

```
backend/
├── controllers/      # 控制器层
├── database/         # 数据库层
├── init/            # 系统初始化
├── middleware/      # 中间件
├── models/          # 数据模型
├── repository/      # Repository 数据访问层 (新增)
│   ├── repository.go              # 通用基类
│   ├── user.go                   # 用户 Repository
│   ├── article.go                # 文章 Repository
│   ├── personal_notification.go  # 个人通知 Repository
│   ├── permission_group.go       # 权限组 Repository
│   └── system_log.go            # 系统日志 Repository
├── service/         # 服务容器 (新增)
└── utils/           # 工具函数

frontend/
├── src/
│   ├── views/
│   │   └── admin/
│   │       └── AdminUserNotifications.vue  # 用户通知和权限管理
│   └── api/
│       └── index.js    # API 模块导出 (更新)
```

## 新增功能

### 1. 用户单独通知系统
- 向单个用户发送通知
- 支持批量发送通知
- 支持通知优先级（high/normal/low）
- 支持通知类型（system/warning/promotion/reminder/activity）
- 支持已读状态追踪

### 2. 完整权限组系统
- 自定义权限组
- 权限级别（Level）管理
- JSON 格式权限列表
- 权限组过期时间
- 默认权限组设置
- 4个默认权限组：
  - 新人（Level 1）
  - 普通用户（Level 10）
  - 版主（Level 50）
  - 内容审核员（Level 60）

### 3. Repository 数据访问层
- 统一的 CRUD 操作
- 泛型支持 `BaseRepository[T]`
- 业务 Repository 封装
- 链式调用：`Where`, `Preload`, `Order`
- 分页查询支持

### 4. 系统初始化
- 自动创建默认权限组
- 自动创建默认用户头衔
- 自动为无权限组用户分配默认组

## 权限系统详解

### 权限列表示例

```javascript
[
  "article.view",
  "article.create",
  "article.edit.own",
  "article.edit.all",
  "article.delete.own",
  "article.delete.all",
  "article.pin",
  "comment.view",
  "comment.create",
  "comment.edit.all",
  "comment.delete.all",
  "user.view.all",
  "user.ban",
  "report.view",
  "report.handle",
  "deletion.view",
  "deletion.handle"
]
```

### 权限检查中间件使用

```go
// 检查单个权限
protected.POST("/articles/:id/pin", middleware.RequirePermission("article.pin"), handler)

// 检查任一权限
protected.GET("/admin/comments", middleware.RequireAnyPermission("comment.view", "comment.delete.all"), handler)
```

## API 接口说明

### 用户通知 API

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| POST | `/api/user-notifications/send` | 发送单独通知 | 管理员 |
| POST | `/api/user-notifications/send-batch` | 批量发送通知 | 管理员 |
| GET | `/api/user-notifications` | 获取用户通知 | 登录用户 |
| POST | `/api/user-notifications/:id/read` | 标记已读 | 登录用户 |
| POST | `/api/user-notifications/read-all` | 全部已读 | 登录用户 |
| DELETE | `/api/user-notifications/:id` | 删除通知 | 登录用户 |
| DELETE | `/api/user-notifications/clear` | 清空通知 | 登录用户 |

### 权限组 API

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| GET | `/api/permission-groups` | 获取权限组列表 | 登录用户 |
| GET | `/api/permission-groups/:id` | 获取权限组详情 | 登录用户 |
| POST | `/api/permission-groups` | 创建权限组 | 管理员 |
| PUT | `/api/permission-groups/:id` | 更新权限组 | 管理员 |
| DELETE | `/api/permission-groups/:id` | 删除权限组 | 管理员 |
| POST | `/api/permission-groups/grant` | 授予用户权限组 | 管理员 |
| DELETE | `/api/permission-groups/:group_id/users/:user_id` | 撤销用户权限组 | 管理员 |
| GET | `/api/users/:id/permission-groups` | 获取用户权限组 | 登录用户 |
| GET | `/api/permissions/check` | 检查用户权限 | 登录用户 |

## 数据库模型

### PersonalNotification（用户个人通知）
```go
type PersonalNotification struct {
    ID          uint
    SenderID    uint
    UserID      uint
    Type        string
    Title       string
    Content     string
    RelatedType string
    RelatedID   uint
    Link        string
    IsRead      bool
    ReadAt      *time.Time
    Priority    string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

### PermissionGroup（权限组）
```go
type PermissionGroup struct {
    ID          uint
    Name        string
    Description string
    Level       int
    IsDefault   bool
    IsActive    bool
    Permissions string // JSON 格式
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

## 构建和部署

```bash
# 构建项目
cd "d:\ByUsi\Projects\Campus Forum"
.\build 1.5.2

# 或者单独构建后端
cd backend
go build -o ../campus_forum.exe
```

## 配置文件说明

`config.json`
```json
{
  "port": "8080",
  "jwt_secret": "your-secret-key",
  "webdav": {
    "url": "https://webdav.example.com",
    "username": "user",
    "password": "pass"
  },
  "database": {
    "host": "localhost",
    "port": "3306",
    "username": "root",
    "password": "",
    "dbname": "forum"
  }
}
```

## 生产级特性

1. **分层架构** - Controller → Service → Repository → Database
2. **统一数据访问** - BaseRepository 提供标准 CRUD
3. **权限控制** - 完善的 RBAC 权限系统
4. **系统初始化** - 自动创建默认权限组和头衔
5. **外键约束处理** - 迁移时禁用外键检查，避免错误
6. **日志记录** - 完整的操作日志和系统日志
7. **速率限制** - API 接口访问频率控制
8. **CORS 支持** - 跨域资源共享

## 迁移指南

从原有代码迁移到 Repository 模式：

1. 创建对应业务的 Repository
2. 将 `database.DB` 调用替换为 Repository 方法
3. 保持业务逻辑不变，只改变数据访问方式
4. 例如 `user_notification_v2.go` 是完整示例

## 开发规范

- Repository 层负责所有数据库操作
- Controller 层只处理请求响应
- Service 层负责业务逻辑编排
- 所有查询尽量使用 Repository 方法
- 复杂查询可使用 `repo.DB()` 获取原始连接
