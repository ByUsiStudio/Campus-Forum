# 更新日志

所有重要的项目变更都将记录在此文件中。

格式遵循 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，版本号遵循 [语义化版本控制](https://semver.org/lang/zh-CN/)。

---

## [2.2.0] - 2026-08-22

### 修复

- 修复系统日志页（AdminSystemLogs）请求带上了双重 `/api` 前缀（`/api/api/system-logs`）导致接口 404 的问题
- 清理指向不存在后端接口的死代码：移除 `api/user.js` 中的 `getUserFollowing`/`getUserFollowers`（`/users/:id/following|followers`）、`api/notification.js` 中的 `/follow-notifications*` 系列封装
- 修复管理后台底部导航与顶栏引用了不存在的路由名（`AdminDashboard`→`AdminIndex`、`AdminSettings`→`AdminSiteConfig`），导航跳转恢复正常
- 移除移动端应用抽屉中指向不存在页面的「聊天」入口及对应无效样式
- 修复排行榜「查看」按钮指向不存在的 `/user/:id` 路由，改跳转至个人中心
- 修复 `Profile` 与 `FollowList` 中未导入的 `showSuccess` 导致的关注操作运行时报错
- 修复全局弹窗被 Esc/点击遮罩关闭后 Promise 永不 resolve、`await confirm/alert/prompt` 永久挂起的阻塞问题
- 修复话题列表分页 `@input` 事件在 Vuetify 3 中不触发导致的分页失效
- 修复评论嵌套回复在 `v-if`/`v-for` 同元素下（Vue 3 优先规则）导致的回复不渲染
- 修复上传弹窗进度条使用失效的 `:value` 且无真实上传进度的问题，改用 XHR 实现真实进度、仅成功后置为 100%
- 修复 Markdown 编辑器自动补全后残留光标占位符 `|` 的问题
- 修复签到页签到后「签到记录」列表不刷新、以及 `v-icon left` 等 Vuetify 2 遗留写法失效的问题
- 修复视频播放器退出全屏 `scrollTo` 使用无效 `behavior` 值的问题
- 修复个人中心「已发布文章」在本人主页重复显示两遍的排版问题
- 修复 404/403 错误页使用无效的 `min-h-screen`（Tailwind）类导致不垂直居中的问题
- 兼容 Vuetify 3：修正 `UserLevel` 组件中 `white--text`/`headline`、`v-progress-linear :value`、`v-row dense`、`v-card outlined` 等失效写法
- 清理文章详情与侧边栏中刷屏的调试日志

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