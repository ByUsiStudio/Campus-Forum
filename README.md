# Campus Forum 校园论坛

Campus Forum 是一个功能完善的校园社区论坛系统，采用前后端分离架构开发。后端基于 Go 语言 + Gin 框架构建，提供高效稳定的 RESTful API 服务；前端采用 Vue 3 构建单页应用，为用户提供流畅的浏览体验。支持文章发布、评论互动、用户认证等完整论坛功能。

> [**更新日志**](CHANGELOG.md)

## 技术栈

### 后端
- **语言**：Go 1.21
- **框架**：Gin
- **数据库**：MySQL + GORM
- **认证**：JWT
- **文件存储**：WebDAV
- **即时通讯**：IM Server

### 前端
- **框架**：Vue 3
- **UI 组件**：Vuetify 3
- **路由**：Vue Router
- **构建工具**：Vite

### 部署
- **容器化**：Docker + Docker Compose
- **Web服务器**：Nginx

## 功能特性

### 用户系统
- 用户注册与登录
- JWT 身份认证
- 个人资料管理
- 头像上传

### 文章系统
- 创建、编辑、删除文章
- 文章分类浏览
- Markdown 富文本支持
- 语音朗读功能
- 匿名发布
- 点赞/取消点赞
- 收藏文章
- 评论互动（支持嵌套回复）
- 匿名评论

### 好友系统
- 双向好友关系（替代原关注系统）
- 好友请求发送、同意、拒绝
- 好友备注名管理
- 好友列表查看
- 共同好友查看

### 即时聊天
- 实时消息发送与接收
- 会话列表管理
- 未读消息计数
- 消息已读标记

### 内容管理
- 文章分类管理
- 站点公告发布
- 侧边栏配置
- 删除请求审批流程

### 文件管理
- 头像上传
- 图片上传
- 视频上传
- 语音上传
- WebDAV 云存储支持

### 消息系统
- 评论回复通知
- 粉丝通知
- 系统通知
- 好友请求通知

### 签到系统
- 每日签到功能
- 连续签到天数统计
- 签到历史记录查询

### 管理员功能
- 用户封禁/解封
- 文章置顶/取消置顶
- 举报管理系统

## API 文档

- **[API 文档](docs/api/README.md)**

## 项目结构

```
campus-forum/
├── backend/              # Go 后端服务
│   ├── controllers/      # 控制器层
│   ├── database/         # 数据库初始化
│   ├── middleware/        # 中间件
│   ├── models/           # 数据模型
│   ├── utils/            # 工具函数
│   ├── sdk/              # 第三方 SDK（IM Server）
│   ├── main.go          # 入口文件
│   └── config.json       # 配置文件
├── frontend/             # Vue 前端应用
│   ├── src/
│   │   ├── api/         # API 接口
│   │   ├── components/  # 公共组件
│   │   ├── views/       # 页面视图
│   │   ├── utils/       # 工具函数
│   │   ├── App.vue
│   │   └── main.js
│   └── index.html
├── nginx/               # Nginx 配置
├── docker-compose.yml   # 容器编排
├── build.sh             # Linux 构建脚本
├── build.bat            # Windows 构建脚本
└── CHANGELOG.md         # 更新日志
```

## 快速开始

### 环境要求

- Docker >= 20.10
- Docker Compose >= 2.0

### 使用 Docker Compose 部署

```bash
# 克隆项目
git clone https://gitee.com/byusistudio/campus-forum.git
cd campus-forum

# 启动所有服务
docker-compose up -d
```

服务启动后将提供以下端口：
- 前端：http://localhost
- 后端 API：http://localhost:3620

### 手动部署

#### 后端启动

```bash
cd backend

# 修改配置文件 config.json
# 配置数据库连接、WebDAV 等参数

# 运行服务
go run main.go
```

#### 前端启动

```bash
cd frontend

# 安装依赖
npm install

# 开发模式
npm run dev

# 生产构建
npm run build
```

## 多平台编译

### Linux/macOS
```bash
./build.sh
```

### Windows
```batch
build.bat
```

## 使用已经打包好的镜像

> Github仓库更新提交后自动更新 docker 镜像

### 1. 拉取镜像

```bash
docker pull beiane1/campus-forum:latest
```

### 2. 创建 Docker 网络

```bash
docker network create forum-network
```

### 3. 启动 MySQL（独立部署）

```bash
docker run -d \
  --name forum-mysql \
  --network forum-network \
  -e MYSQL_ROOT_PASSWORD=your_root_password \
  -e MYSQL_DATABASE=forum \
  -e MYSQL_USER=forum_user \
  -e MYSQL_PASSWORD=your_forum_password \
  -v mysql-data:/var/lib/mysql \
  -p 3306:3306 \
  mysql:8.0
```

### 4. 启动 Forum 应用

```bash
docker run -d \
  --name campus-forum \
  --network forum-network \
  -e DB_HOST=forum-mysql \
  -e DB_PORT=3306 \
  -e DB_USER=forum_user \
  -e DB_PASSWORD=your_forum_password \
  -e DB_NAME=forum \
  -p 80:80 \
  -p 3620:3620 \
  beiane1/campus-forum:latest
```

## 配置说明

### 后端配置 (backend/config.json)

```json
{
    "port": 3620,
    "jwt_secret": "your-secret-key-change-this-in-production",
    "webdav": {
        "url": "http://localhost:5005",
        "username": "webdav_user",
        "password": "webdav_pass"
    },
    "database": {
        "host": "localhost",
        "port": "3306",
        "username": "root",
        "password": "123456",
        "dbname": "forum"
    }
}
```

## 开源说明

本项目基于 AGPL-3.0 许可证开源，详见 [LICENSE](./LICENSE) 文件。

开源仓库：
- GitHub：https://github.com/ByUsiStudio/Campus-Forum
- Gitee：https://gitee.com/byusistudio/campus-forum

## 贡献者

欢迎提交 Issue 和 Pull Request 来完善这个项目。
