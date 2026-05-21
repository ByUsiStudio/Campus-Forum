# Campus Forum 校园论坛

Campus Forum 是一个功能完善的校园社区论坛系统，采用前后端分离架构开发。后端基于 Go 语言 + Gin 框架构建，提供高效稳定的 RESTful API 服务；前端采用 Vue 3 构建单页应用，为用户提供流畅的浏览体验。支持文章发布、评论互动、用户认证等完整论坛功能。

## 技术栈

### 后端
- **语言**：Go 1.21
- **框架**：Gin
- **数据库**：MySQL + GORM
- **认证**：JWT
- **文件存储**：WebDAV

### 前端
- **框架**：Vue 3
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
- 点赞/取消点赞
- 评论互动

### 内容管理
- 文章分类管理
- 站点公告发布
- 侧边栏配置
- 删除请求审批流程

### 文件管理
- 头像上传
- 图片上传
- 视频上传
- WebDAV 云存储支持

## 项目结构

```
campus-forum/
├── backend/              # Go 后端服务
│   ├── controllers/      # 控制器层
│   ├── database/         # 数据库初始化
│   ├── middleware/       # 中间件
│   ├── models/           # 数据模型
│   ├── utils/           # 工具函数
│   ├── main.go          # 入口文件
│   └── config.json       # 配置文件
├── frontend/             # Vue 前端应用
│   ├── src/
│   │   ├── api/         # API 接口
│   │   ├── components/ # 公共组件
│   │   ├── views/      # 页面视图
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
- 后端 API：http://localhost:8080

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

## API 接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/auth/register | 用户注册 |
| POST | /api/auth/login | 用户登录 |
| GET | /api/auth/profile | 获取用户资料 |
| PUT | /api/auth/profile | 更新用户资料 |
| GET | /api/articles | 获取文章列表 |
| POST | /api/articles | 创建文章 |
| GET | /api/articles/:id | 获取文章详情 |
| PUT | /api/articles/:id | 更新文章 |
| DELETE | /api/articles/:id | 删除文章 |
| POST | /api/articles/:id/like | 点赞文章 |
| DELETE | /api/articles/:id/like | 取消点赞 |
| POST | /api/articles/:id/comments | 添加评论 |
| DELETE | /api/comments/:id | 删除评论 |
| GET | /api/categories | 获取分类列表 |
| POST | /api/categories | 创建分类 |
| GET | /api/announcements | 获取公告 |
| PUT | /api/announcements | 更新公告 |

## 配置说明

### 后端配置 (backend/config.json)

```json
{
  "server": {
    "port": "8080"
  },
  "database": {
    "dsn": "user:password@tcp(host:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
  },
  "jwt": {
    "secret": "your-jwt-secret-key"
  },
  "webdav": {
    "url": "https://your-webdav-server.com",
    "username": "your-username",
    "password": "your-password"
  }
}
```

## 许可证

本项目基于 AGPL-3.0 许可证开源，详见 [LICENSE](./LICENSE) 文件。

## 贡献者

欢迎提交 Issue 和 Pull Request 来完善这个项目。