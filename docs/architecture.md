# 架构设计文档

## 概述

校园论坛采用标准的三层架构模式（Repository-Service-Controller），实现了关注点分离，提高了代码的可维护性、可测试性和可扩展性。

---

## 架构分层

### Controller 层（HTTP 控制层）

**职责**:
- 处理 HTTP 请求和响应
- 参数校验和绑定
- 调用 Service 层处理业务逻辑
- 返回统一的响应格式

**设计原则**:
- 不包含业务逻辑
- 不直接访问数据库
- 只处理 HTTP 相关的操作
- 统一错误处理和响应格式

**目录结构**:
```
backend/controllers/
├── auth.go           # 认证相关接口
├── article.go        # 文章相关接口
├── comment.go        # 评论相关接口
├── category.go       # 分类相关接口
├── follow.go         # 关注相关接口
├── leaderboard.go    # 排行榜相关接口
├── notification.go   # 系统通知接口
├── user_notification.go # 用户个人通知接口
├── admin.go          # 管理员操作接口
├── siteconfig.go     # 站点配置接口
├── user.go           # 用户资料接口
├── favorite.go       # 收藏相关接口
├── friends.go        # 好友系统接口
└── ...
```

**示例**:
```go
func CreateArticle(c *gin.Context) {
    userID := c.GetUint("user_id")
    
    var input struct {
        Title       string `json:"title" binding:"required"`
        Content     string `json:"content" binding:"required"`
        CategoryID  uint   `json:"category_id"`
        IsAnonymous bool   `json:"is_anonymous"`
    }
    
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    article, err := service.Article.CreateArticle(userID, input.Title, input.Content, input.CategoryID, input.IsAnonymous)
    if err != nil {
        if appErr, ok := utils.IsAppError(err); ok {
            c.JSON(appErr.Code, gin.H{"error": appErr.Message})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "创建成功", "article": article})
}
```

---

### Service 层（业务逻辑层）

**职责**:
- 实现核心业务逻辑
- 调用 Repository 层进行数据访问
- 统一错误处理和业务规则验证
- 事务管理

**设计原则**:
- 不处理 HTTP 请求/响应
- 不直接操作数据库（通过 Repository）
- 业务逻辑与数据访问解耦
- 可独立测试

**目录结构**:
```
backend/service/
├── auth_service.go          # 认证业务逻辑
├── article_service.go       # 文章管理逻辑
├── comment_service.go       # 评论管理逻辑
├── follow_service.go        # 关注和好友逻辑
├── leaderboard_service.go   # 排行榜逻辑
├── user_notification_service.go # 用户通知逻辑
├── admin_service.go         # 管理员操作逻辑
├── siteconfig_service.go    # 站点配置逻辑
└── service.go               # ServiceContainer 初始化
```

**示例**:
```go
func (s *ArticleService) CreateArticle(userID uint, title, content string, categoryID uint, isAnonymous bool) (*models.Article, error) {
    contentHTML := utils.RenderMarkdown(content)
    
    article := models.Article{
        UserID:      userID,
        Title:       title,
        Content:     content,
        ContentHTML: contentHTML,
        CategoryID:  categoryID,
        IsAnonymous: isAnonymous,
        Status:      "published",
    }
    
    if result := database.DB.Create(&article); result.Error != nil {
        return nil, utils.NewError("创建文章失败", 500)
    }
    
    s.updateUserExperience(userID, "article_create", 10)
    
    return &article, nil
}
```

---

### Repository 层（数据访问层）

**职责**:
- 封装数据库操作
- 提供通用 CRUD 接口
- 各模型专用查询方法
- 事务管理支持

**设计原则**:
- 只负责数据访问
- 不包含业务逻辑
- 提供通用接口，便于替换实现
- 支持事务

**目录结构**:
```
backend/repository/
├── repository.go            # 通用 BaseRepository 接口和实现
├── user.go                  # 用户数据访问
├── article.go               # 文章数据访问
├── comment.go               # 评论数据访问
├── category.go              # 分类数据访问
├── follow.go                # 关注数据访问
├── notification.go          # 系统通知数据访问
├── personal_notification.go # 个人通知数据访问
└── ...
```

**通用接口**:
```go
type BaseRepository interface {
    Create(model interface{}) error
    Update(model interface{}) error
    Delete(id uint) error
    GetByID(id uint, model interface{}) error
    FindAll(model interface{}) error
    FindByCondition(model interface{}, condition interface{}) error
    Count(model interface{}) (int64, error)
}
```

**示例**:
```go
type ArticleRepository struct {
    *repository.BaseRepository
}

func NewArticleRepository() *ArticleRepository {
    return &ArticleRepository{
        BaseRepository: repository.NewBaseRepository(),
    }
}

func (r *ArticleRepository) GetArticlesByUser(userID uint, page, pageSize int) ([]models.Article, int, error) {
    var articles []models.Article
    var total int64
    
    query := database.DB.Model(&models.Article{}).Where("user_id = ?", userID)
    query.Count(&total)
    
    offset := (page - 1) * pageSize
    err := query.Preload("Category").Order("created_at DESC").
        Offset(offset).Limit(pageSize).Find(&articles).Error
    
    totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
    
    return articles, totalPages, err
}
```

---

## 数据流向

```
HTTP Request
     │
     ▼
Controller Layer
     │  参数校验
     │  调用 Service
     ▼
Service Layer
     │  业务逻辑处理
     │  调用 Repository
     ▼
Repository Layer
     │  数据库操作
     ▼
Database
```

---

## 错误处理

### 自定义错误类型

```go
type AppError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}

func NewError(message string, code int) error {
    return &AppError{
        Code:    code,
        Message: message,
    }
}

func IsAppError(err error) (*AppError, bool) {
    appErr, ok := err.(*AppError)
    return appErr, ok
}
```

### 错误处理流程

1. Repository 层返回原始错误或自定义错误
2. Service 层捕获错误，转换为 AppError 或直接传递
3. Controller 层使用 `utils.IsAppError()` 判断错误类型
4. 根据错误类型返回相应的 HTTP 状态码

---

## 依赖注入

### ServiceContainer

通过 ServiceContainer 管理所有 Service 和 Repository 的实例：

```go
type ServiceContainer struct {
    Auth           *AuthService
    Article        *ArticleService
    Comment        *CommentService
    Follow         *FollowService
    Leaderboard    *LeaderboardService
    UserNotification *UserNotificationService
    Admin          *AdminService
    AdminConfig    *SiteConfigService
}

var Services *ServiceContainer

func InitServices() {
    Services = &ServiceContainer{
        Auth:           &AuthService{},
        Article:        &ArticleService{},
        Comment:        &CommentService{},
        Follow:         &FollowService{},
        Leaderboard:    &LeaderboardService{},
        UserNotification: &UserNotificationService{},
        Admin:          &AdminService{},
        AdminConfig:    &SiteConfigService{},
    }
}
```

---

## 状态管理（前端）

### 集中式 Store

前端使用 Pinia-like 的集中式状态管理：

```js
import { reactive, computed } from 'vue'

const state = reactive({
    user: null,
    token: '',
    refreshToken: '',
    siteConfig: null,
    unreadCount: 0,
})

const store = {
    state,
    user: computed(() => state.user),
    token: computed(() => state.token),
    
    setUser(user) {
        state.user = user
    },
    
    setToken(token) {
        state.token = token
    },
    
    // ... 其他方法
}

export default store
```

### API 层

使用 Axios 封装，支持请求/响应拦截器：

- 请求拦截器：自动添加 Token
- 响应拦截器：统一错误处理，Token 自动刷新

---

## 权限系统

### 权限组层级

| 权限组 | 级别 | 描述 |
|--------|------|------|
| system | 100 | 系统管理员，拥有所有权限 |
| admin | 80 | 管理员，管理普通用户内容 |
| 版主 | 50 | 管理板块内容 |
| 普通用户 | 10 | 基础权限 |
| 新人 | 1 | 新注册用户默认权限 |

### 权限中间件

```go
func RequireMinLevel(level int) gin.HandlerFunc {
    return func(c *gin.Context) {
        userLevel := c.GetInt("user_level")
        if userLevel < level {
            c.JSON(http.StatusForbidden, gin.H{"error": "权限不足"})
            c.Abort()
            return
        }
        c.Next()
    }
}

func RequireSystemAdmin() gin.HandlerFunc {
    return func(c *gin.Context) {
        role := c.GetString("role")
        if role != "system" {
            c.JSON(http.StatusForbidden, gin.H{"error": "需要系统管理员权限"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```

---

## 项目结构

```
campus-forum/
├── backend/
│   ├── controllers/      # HTTP 控制层
│   ├── service/          # 业务逻辑层
│   ├── repository/       # 数据访问层
│   ├── models/           # 数据模型
│   ├── middleware/       # 中间件
│   ├── database/         # 数据库初始化
│   ├── utils/            # 工具函数
│   ├── init/             # 系统初始化
│   ├── main.go           # 入口文件
│   └── config.json       # 配置文件
├── frontend/
│   ├── src/
│   │   ├── api/          # API 接口
│   │   ├── stores/       # 状态管理
│   │   ├── hooks/        # 自定义 hooks
│   │   ├── utils/        # 工具函数
│   │   ├── components/   # 公共组件
│   │   ├── views/        # 页面视图
│   │   ├── App.vue
│   │   └── main.js
│   └── index.html
├── docs/
│   ├── api/              # API 文档
│   └── architecture.md   # 架构文档
├── docker-compose.yml
├── build.sh
├── build.bat
├── README.md
└── CHANGELOG.md
```

---

## 设计模式总结

| 模式 | 应用位置 | 目的 |
|------|----------|------|
| Repository | repository/ | 数据访问抽象，解耦业务逻辑与数据库操作 |
| Service | service/ | 业务逻辑封装，提高可测试性 |
| Controller | controllers/ | HTTP 请求处理，参数校验 |
| Dependency Injection | service/service.go | 统一管理服务实例 |
| Facade | 前端 API 层 | 统一的 API 调用接口 |
| Observer | 前端 Store | 响应式状态更新 |

---

## 开发规范

### 新增功能流程

1. **创建/更新 Model** - 定义数据结构
2. **创建/更新 Repository** - 实现数据访问方法
3. **创建/更新 Service** - 实现业务逻辑
4. **创建/更新 Controller** - 实现 HTTP 接口
5. **注册路由** - 在 main.go 中注册新路由
6. **更新 API 文档** - 在 docs/api/ 中添加接口说明

### 代码风格

- 使用 Go 标准格式化工具：`gofmt`
- 使用 Vue 3 Composition API
- 保持代码简洁，避免重复
- 添加适当的注释（中文）
- 遵循命名规范：
  - Go：驼峰命名（CamelCase）
  - JavaScript：驼峰命名（camelCase）
  - 文件名：小写加下划线（snake_case）

---

## 测试建议

### 单元测试

- **Repository 层**：测试数据库 CRUD 操作
- **Service 层**：测试业务逻辑，使用 Mock Repository
- **Controller 层**：测试 HTTP 请求响应，使用 Gin 测试工具

### 集成测试

- 测试完整的 API 调用流程
- 测试权限控制
- 测试错误处理

### 性能测试

- 测试高并发场景
- 测试数据库查询性能
- 测试 API 响应时间
