# 更新日志

所有重要的项目变更都将记录在此文件中。

格式遵循 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，版本号遵循 [语义化版本控制](https://semver.org/lang/zh-CN/)。

---

## [3.2.1] - 2026-08-22

## [Unreleased]

### 优化

- **手机端顶部栏**：搜索框改为搜索按钮，点击跳转到独立搜索页 `/search` 展示搜索框；桌面端仍保留行内输入框
- **权限精简**：用户体系统一为「管理员 admin」与「普通用户」两种角色。后端不再允许把用户设为 `system` 角色，移除贯穿各接口的 `system` 专属保护（更新/删除/封禁/通知等均归管理员管理）；`system` 角色用户在数据库中的旧数据仍保留管理员等价权限，前端管理员判断仅认 `admin`
- **文章列表展示**：收紧文章卡片的行边距与内边距（卡片 padding、标题/摘要行高与间距、页脚上边距），移动端显示更紧凑
- **版本号展示**：修复「unknown」与前后端版本不一致问题。后端 `/version` 在未通过构建脚本注入版本时回退读取项目根目录 `version.json`；CI / build.sh / build.bat 均注入前后端版本号，前端构建通过 `VITE_APP_VERSION` 传入本次发布版本号（此前用 `--version` 会被 Vite CLI 误判导致构建失败）；后台侧边栏版本信息优化为分栏标签展示

### 修复

- 移除构建脚本中无效的 `-X forum/controllers.SwaggerVersion=...` ldflags（该符号不存在，属于死参数）

---

## [3.0.1] - 2026-08-22

在 3.0.0（Element Plus 重构）基础上完成的移动端适配、前后端对接修复与安全加固。

### 新增

- **手机端全局适配**：所有 `el-dialog` 在 `≤575px` 屏自动收窄为 `92vw`，消除 15+ 处固定宽度弹窗在手机端横向溢出（App 公告、文章分享/举报、写文章语音、收藏夹/话题详情、后台各面板、编辑器、上传弹窗等）
- 补齐项目中使用但此前缺失定义的响应式显示工具类 `d-none / d-sm-* / d-md-* / d-lg-* / d-xl-*`（修复后台外壳 `d-lg-none`、搜索框 `d-none d-sm-block` 等无效类）
- **App 顶部导航**手机端新增汉堡入口与下拉导航面板（登录/注册/写文章/收藏/签到/话题/个人中心/排行榜/通知/管理后台），此前导航在手机上被隐藏后不可达；路由切换自动收起
- **管理后台外壳**移动端侧边栏改为离屏抽屉（遮罩 + 平移过渡），汉堡按钮真正控制开合，路由切换自动关闭
- Markdown 正文中的表格包裹进可横向滚动容器，避免窄屏截断/溢出

### 修复

- **对接（D1–D6 / M1–M4）**：
  - 删除文章审核流：普通用户删除文章改为提交带 `reason` 的待审核 `DeletionRequest`，不再被静默丢弃；管理员仍直接软删
  - 个人中心好友/粉丝计数读取正确字段（`mutual_friends`），不再恒为 0
  - 举报管理（AdminReports）状态统计卡片改为读取后端返回的 `pending_count/resolved_count/rejected_count`；`GET /reports` 后端支持 `search`/`target_type` 过滤
  - 侧边栏配置字段统一 `url → link`，与后端模型一致
  - 修复 `adminUserNotificationApi` 重复具名导出：统一从 `admin.js`（完整转发 `priority`/`link`）导出一处，避免发送通知信息丢失
  - `/my/articles` 仅返回已发布文章，与草稿接口互斥
- **安全**：`UpdateSiteConfig` 不再返回 SMTP 明文密码，改为返回 `smtp_password_set` 标记（此前明文密码随响应泄露）
- 修复登录失败时 `error` 被同名 ref 遮蔽导致的运行时 TypeError；登录后自动拉取 `/profile` 补全用户资料字段
- 修复管理后台点击菜单切换到设置项（网站配置 `siteconfig` / 邮件配置 `smtpconfig`）命中前端通配路由导致 **404** 的问题（导航目标 `/<path>` → `/admin/<path>`）
- 修复 `GetTopics` 分页 `total` 未跟随 `is_hot` 过滤导致热门话题页总页数不准
- 修复 `frontend/src/api/` 下 7 个文件（`article / user / common / comment / auth / notification / upload`）共 77 处中文注释编码乱码（含多行注释吞并方法名已拆分），不改变任何逻辑

### 移动端适配（其余视图）

- `Category` 页改用 `el-row` 响应式栅格并在移动端隐藏侧栏（此前无样式栅格 + 侧栏不隐藏）
- `Home` 分区选择器、`SignIn` 周/月统计卡、`Leaderboard` 分段选择器、`CollectionList` 标题区、`NotificationBell` 弹层宽度、`AdminSystemLogs` 分页/筛选在窄屏的适配

### 文档

- 同步 `docs/api/`：`report.md`、`deletion.md`、`articles.md`、`system.md`、`topic.md` 至实际响应结构

### 注意

- 沙箱环境无法执行真实 `vite build` / `go build`（esbuild 子进程、Go 模块缓存写入被拦截），本次以逐文件 `node --check` 语法校验、图标具名导入校验、后端改动大括号配平作为静态验证依据；请在本地 `cd frontend && yarn build`、`cd backend && go build ./...` 做最终构建确认

---

## [3.0.0] - 2026-08-22

### 重大重构

将前端从 **Vuetify 3** 彻底重构为 **Element Plus** 技术栈，采用 Vue 3 + Vite + Vue Router + **Pinia**。

### 新增

- 引入 **Pinia 状态管理**：`src/stores/`（`user` / `app` / `notification`），替代散落的 `localStorage` 与跨组件事件
- 引入 **Element Plus** 组件体系及 `@element-plus/icons-vue`，全局组件/图标注册，中文本地化
- 新增统一交互反馈工具 `src/utils/message.js`（`success / error / confirm / prompt / warning`，基于 `ElMessage` / `ElMessageBox`）

### 重构

- **目录分层**：路由独立为 `src/router/`（含守卫/懒加载），状态独立为 `src/stores/`，全局样式为 `src/styles/index.css`
- **API 数据层**：新增统一 `src/api/http.js`（axios 实例，token 自动刷新、并发去重、统一错误提示）；`src/api/` 各业务模块统一从 `./http` 复用，`src/api/index.js` 默认导出 `http` 并透出全部具名模块
- **路由守卫**：文章详情/分类/排行榜/话题等只读页面未登录可浏览，登录后自动回跳原目标页
- **全部视图与组件**迁移到 Element Plus（`Home / Article / Profile / CreateArticle / 登录注册 / 排行榜 / 话题 / 后台管理全部页面 / 编辑器 MdEditor / MarkdownViewer / 上传 / 头像 等）

### 修复

- 修复 7 个 API 模块（`admin / collection / friend / leaderboard / level / statistics / topic`）因编码损坏导致注释吞并方法名、产生 `Unexpected token` 语法错误的问题（重写为正确的 UTF-8 内容）
- 清理依赖与死代码：移除 `vuetify`、`@mdi/font`、`vue3-markdown-it`、`marked` 等旧依赖（公告/编辑器统一改用 `markdown-it` 渲染）；删除 `src/utils/modal.js` 与 `src/styles/style.css`

### 注意

- 沙箱环境无法执行真实 `vite build`（esbuild 子进程被沙箱拦截）；本次以 `@vue/compiler-sfc` 全量编译校验（55 个 `.vue` 全部通过）+ Babel ESM 语法校验（23 个 `.js` 全部通过）作为验证依据，请在本地 `cd frontend && yarn install && yarn dev`（或 `yarn build`）做最终构建与浏览器确认

---

## [2.2.0] - 2026-08-22

### 重构

- 保持现有界面 UI（Vuetify 3 布局/配色/组件）不变的前提下，全量重写前端代码实现：统一 APIs 对接层、规范视图与组件逻辑
- 重写 API 请求层（`api/index.js`）为统一封装：修复请求拦截器中 `getErrorMessage` 未定义的潜在崩溃，改进 token 自动刷新（并发去重、队列按序放行、刷新失败跳登录），统一错误解析与提示
- 重写路由守卫（`main.js`）：文章详情/分类/排行榜/话题等只读页面在未登录时也可浏览（与后端公开接口一致），登录页支持登录后回跳原目标页
- 逐文件重写全部视图的 `<script>` 逻辑（模板/UI 保持逐字不变），清理死代码与调试日志，统一 API 调用规范
- 重写 `api/friend.js` 并与后端对齐：好友/关注请求字段由 `user_id` 修正为后端要求的 `friend_id`

### 修复

- 修复"关注/添加好友"功能因请求体字段名（`user_id` → `friend_id`）与后端不匹配而完全失效的问题（涉及 `Article`、`Profile`、`FollowList`）
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