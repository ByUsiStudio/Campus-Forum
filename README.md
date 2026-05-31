# Campus Forum 校园论坛

Campus Forum 是一个功能完善的校园社区论坛系统，采用前后端分离架构开发。后端基于 Go 语言 + Gin 框架构建，提供高效稳定的 RESTful API 服务；前端采用 Vue 3 构建单页应用，为用户提供流畅的浏览体验。支持文章发布、评论互动、用户认证等完整论坛功能。

> [**更新计划**](./待办计划.MD)

## 技术栈

### 后端
- **语言**：Go 1.21
- **框架**：Gin
- **数据库**：MySQL + GORM
- **认证**：JWT
- **文件存储**：WebDAV

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
- 关注/粉丝系统

### 文章系统
- 创建、编辑、删除文章
- 文章分类浏览
- Markdown 富文本支持
- 语音朗读功能
- 匿名发布
- 点赞/取消点赞
- 收藏文章
- 评论互动
- 匿名评论

### 关注系统
- 关注/取消关注用户
- 粉丝列表查看
- 关注列表查看
- 关注对象发新内容时提醒

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

## API 文档

- **[交互式 API 文档 (Swagger)](http://localhost:3620/swagger/index.html)** - 启动服务后访问
- **[Markdown API 文档](backend/docs/API.md)** - 离线查看

## 项目结构

```
campus-forum/
├── backend/              # Go 后端服务
│   ├── controllers/      # 控制器层
│   ├── database/         # 数据库初始化
│   ├── middleware/        # 中间件
│   ├── models/           # 数据模型
│   ├── utils/            # 工具函数
│   ├── main.go          # 入口文件
│   └── config.json       # 配置文件
├── frontend/             # Vue 前端应用
│   ├── src/
│   │   ├── api/         # API 接口
│   │   ├── components/  # 公共组件
│   │   ├── views/       # 页面视图
│   │   ├── App.vue
│   │   └── main.js
│   └── index.html
├── nginx/               # Nginx 配置
├── docker-compose.yml   # 容器编排
└── build.sh            # 构建脚本
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

> **API 文档**：启动后访问 http://localhost:3620/swagger/index.html 查看完整的 API 文档

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
  -p 8080:8080 \
  beiane1/campus-forum:1.3.11
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
