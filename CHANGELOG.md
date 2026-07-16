# 更新日志

所有重要的项目变更都将记录在此文件中。

格式遵循 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，版本号遵循 [语义化版本控制](https://semver.org/lang/zh-CN/)。

---

## [3.0.0] - 2026-07-16

### 重构

#### 后端架构重构

- **引入三层架构模式（Repository-Service-Controller）**：
  - 新增 `repository/` 目录，创建通用 BaseRepository 接口和各模型专用 Repository
  - 新增 `service/` 目录，创建服务层处理业务逻辑，替代直接在 Controller 中操作数据库
  - 重构 `controllers/` 目录，所有 Controller 仅负责 HTTP 请求处理，调用 Service 层处理业务

- **Repository 层**：
  - 创建 `repository/repository.go` - 通用 CRUD 操作接口和实现
  - 创建 `repository/user.go`, `article.go`, `comment.go`, `category.go`, `follow.go`, `notification.go`, `personal_notification.go` 等模型专用 Repository

- **Service 层**：
  - 创建 `service/auth_service.go` - 认证业务逻辑
  - 创建 `service/article_service.go` - 文章管理逻辑
  - 创建 `service/comment_service.go` - 评论管理逻辑
  - 创建 `service/follow_service.go` - 关注和好友系统逻辑
  - 创建 `service/leaderboard_service.go` - 排行榜逻辑
  - 创建 `service/user_notification_service.go` - 用户通知逻辑
  - 创建 `service/admin_service.go` - 管理员操作逻辑
  - 创建 `service/siteconfig_service.go` - 站点配置逻辑

- **Controller 层重构**：
  - 所有 Controller 调用 Service 层处理业务，不再直接操作数据库
  - 新增 `controllers/friends.go` - 好友系统控制器
  - 更新 `controllers/user_notification.go` - 个人通知控制器
  - 更新 `controllers/admin.go` - 管理员控制器，新增多个管理接口

- **统一错误处理**：
  - 创建 `utils/errors.go` - 自定义 AppError 类型和错误处理工具
  - 所有错误返回统一格式

- **基础设施改进**：
  - 修复 `main.go` 的 graceful shutdown 问题，使用 `http.Server` 和信号处理
  - 添加 Markdown 渲染工具函数 `utils.RenderMarkdown`

#### 前端架构重构

- **状态管理**：
  - 创建 `src/stores/index.js` - 集中式状态管理（Pinia-like Store）
  - 替代直接 localStorage 操作，支持响应式状态
  - 管理用户认证状态、站点配置、通知状态

- **API 层重构**：
  - 更新 `src/api/index.js` - 使用 Axios 拦截器处理 Token
  - 统一错误处理
  - Token 自动刷新机制

- **自定义 Hooks**：
  - 创建 `src/hooks/index.js` - 分页、防抖、节流、滚动监听、尺寸监听等 hooks

- **工具函数**：
  - 创建 `src/utils/index.js` - 日期格式化、文件大小、UUID 生成等工具函数

- **UI 组件**：
  - 更新 `Login.vue` - 使用 Store 进行认证
  - 更新 `App.vue` - 使用 Store 值进行用户状态管理和导航

### 新增

- 新增好友系统完整实现（请求发送、接受、拒绝、删除、共同好友查看）
- 新增个人通知系统（发送、批量发送、标记已读、删除、清空）
- 新增管理员徽章管理功能（授予、撤销、显示状态更新）
- 新增用户排行榜查询功能
- 新增 SMTP 配置测试接口

### 修复

- 修复重复路由注册问题
- 修复控制器函数名不匹配问题
- 修复模型字段名不匹配问题
- 修复构建脚本中的路径问题

### 改进

- 更新项目结构文档，反映新的三层架构
- 更新技术栈描述，添加 Repository-Service-Controller 架构模式
- 更新构建脚本，增加版本号参数说明

---

## [2.0.4] - 2026-06-20

### 修复

- 修复重复路由注册（`/api/leaderboard`、`/api/achievements/all`、`/api/level/config` 在认证和公开路由组中重复定义）
- 调整路由配置，将话题、排行榜、成就配置等公开接口统一移至 `api` 路由组

---

## [2.0.3] - 2026-06-20

### 修复

- 修复路由参数名冲突（`/articles/:article_id` 与 `/articles/:id` 冲突）
- 修复重复路由注册（`/api/topics` 在认证和公开路由组中重复定义）
- 修复控制器函数名不匹配（`GetLevelConfig`、`GetUserExperienceRecords`）
- 修复模型字段名不匹配（`NextLevelExp` → `NextLevel`）
- 修复模型类型名不匹配（`TopicArticle` → `ArticleTopic`）
- 修复文件编码错误（重建 `collection.go`、`level.go`、`leaderboard.go`、`topic.go`）
- 修复 API 文档中的路由参数与代码不一致

---

## [2.0.2] - 2026-06-20

### 新增

- **用户等级与成就系统**：
  - 新增用户等级系统，支持经验值积累和等级提升
  - 新增成就系统，包含多种成就类型和稀有度等级
  - 新增等级配置管理，管理员可自定义等级称号和特权
  - 新增成就管理功能，管理员可创建和管理成就
  - 新增经验值记录查询功能

- **数据统计与分析**：
  - 新增用户统计数据模型，记录用户活跃度和贡献数据
  - 新增每日统计数据模型，记录系统每日运营数据
  - 新增文章统计数据模型，记录文章浏览和互动数据
  - 新增系统概览数据模型，提供系统整体运营数据
  - 新增统计仪表板功能，管理员可查看系统运营数据
  - 新增用户活跃度分析功能

- **内容管理增强**：
  - 新增收藏夹功能，用户可创建和管理多个收藏夹
  - 新增收藏夹文章关联功能，支持添加收藏笔记
  - 新增文章版本历史功能，记录文章修改历史
  - 新增文章版本恢复功能，支持恢复到历史版本
  - 新增话题标签系统，支持文章话题分类
  - 新增话题关注功能，用户可关注感兴趣的话题
  - 新增热门话题功能，自动计算和展示热门话题

- **社交互动增强**：
  - 新增排行榜系统，支持多种排行榜类型和统计周期
  - 新增用户徽章系统，支持多种徽章类型和展示管理
  - 新增用户活跃度模型，记录用户每日活跃数据
  - 新增排行榜查询功能，用户可查看各类排行榜
  - 新增徽章授予和管理功能，管理员可授予徽章

### 改进

- 更新数据库迁移，新增所有新功能的数据表
- 更新API路由，新增所有新功能的API接口
- 更新前端API模块，新增所有新功能的API调用
- 更新前端组件，新增用户等级、排行榜、话题、收藏夹等组件

---

## [1.5.9] - 2026-06-18

### 新增

- **Token 自动刷新机制**：
  - 后端新增 `RefreshToken` 接口，支持刷新访问令牌
  - 登录时返回 `access_token`（1小时）和 `refresh_token`（7天）
  - 前端 API 拦截器自动检测并刷新即将过期的令牌

- **权限组系统优化**：
  - system 权限组权限大于 admin（级别 100 vs 80）
  - admin 只能管理普通用户内容（删除评论文章、标记删除用户、更新公告）
  - system 组拥有所有权限（包括系统配置、权限组初始化、系统日志）
  - admin 无法对 system 组的通知进行任何更改

- **构建脚本优化**：
  - 编译时自动复制 `config.json` 到 build 目录
  - 服务端支持自动生成默认配置文件

- **Markdown 链接优化**：
  - 文章查看界面的外链自动在新标签打开
  - 添加 `rel="noopener noreferrer"` 提高安全性
  - 本站链接保持原行为

- **公告弹窗展示**：
  - 公告使用弹窗展示，支持 Markdown 语法自动解析
  - 支持"不再显示"功能

### 修复

- 修复文章查看界面 404 接口问题（关注接口改为好友接口）
- 修复 `CommentReply.vue` 中 v-model 绑定 prop 的问题

### 改进

- 更新 API 文档（auth.md、permission.md）
- 权限中间件实际应用到所有管理后台路由

---

## [1.5.8] - 2026-06-14

### 移除

- 删除了 IM 相关代码（聊天功能）：
  - 删除 `backend/im_starter.go`
  - 删除 `frontend/src/composables/useGoIM.js`
  - 删除 `frontend/src/api/im.js`
  - 删除 `frontend/src/views/ChatList.vue`
  - 删除 `frontend/src/views/ChatRoom.vue`
  - 移除路由配置中的聊天路由
  - 移除配置文件中的 IM 配置项
  - 更新构建脚本移除 IM 服务器编译步骤

### 新增

- 新增统一错误处理：
  - 后端 `backend/utils/errors.go` 提供统一错误响应格式和便捷函数
  - 前端 `NotFound.vue` - 404 页面，显示浮动动画图标和返回按钮
  - 前端 `Forbidden.vue` - 403 页面，显示警告锁图标和返回按钮
  - 更新路由配置添加 `/403` 和通配符路由
  - 更新 API 拦截器支持解析统一错误响应格式

### 修复

- 修复前端构建错误：
  - 添加 Vite 路径别名配置 `@`
  - 修复 ChatRoom.vue 中 `getSenderAvatar` 和 `messagesContainer` 重复定义
  - 修复 API 导出重复问题（`adminNotificationApi`、`adminUserNotificationApi`）
  - 修复 Go modules 模式冲突问题

### 改进

- 更新评论 API 文档（`docs/api/comments.md`），补充嵌套回复说明（最大5层嵌套）
- 减小 Markdown 文档行距（从 1.6 调整为 1.4）
- 改进文章页面错误处理，404 时显示友好错误页面
