# Campus Forum 前端（Element Plus 重构版）

基于 **Vue 3 + Vite + Vue Router + Pinia + Element Plus** 的校园论坛前端。

## 技术栈

- **构建**：Vite 4
- **框架**：Vue 3（`<script setup>` / Options API 混用，建议新代码用 script setup）
- **路由**：Vue Router 4（懒加载）
- **状态管理**：Pinia
- **UI 组件**：Element Plus + `@element-plus/icons-vue`

## 目录结构

```
src/
├── api/            # 统一 API 层
│   ├── http.js     #   axios 实例（token 刷新 / 错误处理 / 并发去重）
│   ├── index.js    #   默认导出 http + 透出全部业务模块
│   └── *.js        #   业务模块（article / user / admin / ...）
├── router/         # 路由定义与守卫（index.js）
├── stores/         # Pinia 状态（user / app / notification）
├── components/     # 公共组件（UserAvatar / ArticleList / MdEditor / ...）
├── views/          # 页面（含 admin/ 后台管理）
├── utils/          # message.js 等工具
└── styles/         # 全局样式（index.css + 设计变量）
```

## 常用命令

```bash
cd frontend
yarn install
yarn dev        # 开发
yarn build      # 生产构建
```

## 关键约定

- **API 调用**：
  - 直接用 axios 实例：`import http from '@/api/http'` → `http.get('/xxx')`
  - 或用具名模块：`import { articleApi } from '@/api'`
  - 默认导出 `import api from '@/api'` 等价于 `http`（axios 实例）。
  - **注意**：`http` 的 baseURL 已是 `/api`，路径不要重复写 `/api/...`。
- **提示 / 确认框**：统一用 `import { success, error, confirm, prompt, warning } from '@/utils/message'`。
- **鉴权状态**：用 `useUserStore()`（`isLoggedIn` / `user` / `isAdmin` / `logout`），不再直接依赖 `localStorage` 跨页传状态。
- **路由**：只读页面（文章/分类/排行榜/话题/搜索）未登录可浏览，登录后自动回跳。

## 迁移自 Vuetify

- 原 UI 框架为 Vuetify 3，现已整体迁移到 Element Plus。
- 主要映射：`v-card → el-card`、`v-btn(color) → el-button(type)`、`v-chip → el-tag`、`v-dialog → el-dialog`、`v-text-field → el-input`、`v-select → el-select`、`v-pagination → el-pagination`、`v-icon(mdi) → el-icon(@element-plus/icons-vue)`。

## 部署

由 `nginx` 反向代理 `/api` 到后端，前端构建产物放置于站点根目录。
