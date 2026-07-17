<template>
  <div>
    <div class="stats-grid">
      <div class="stat-card primary">
        <div class="stat-icon">
          <i class="fa-solid fa-users"></i>
        </div>
        <div class="stat-value">{{ statistics.user_count || 0 }}</div>
        <div class="stat-label">用户总数</div>
      </div>
      
      <div class="stat-card pink">
        <div class="stat-icon">
          <i class="fa-solid fa-file-lines"></i>
        </div>
        <div class="stat-value">{{ statistics.article_count || 0 }}</div>
        <div class="stat-label">文章总数</div>
      </div>
      
      <div class="stat-card cyan">
        <div class="stat-icon">
          <i class="fa-solid fa-comment"></i>
        </div>
        <div class="stat-value">{{ statistics.comment_count || 0 }}</div>
        <div class="stat-label">评论总数</div>
      </div>
      
      <div class="stat-card green">
        <div class="stat-icon">
          <i class="fa-solid fa-eye"></i>
        </div>
        <div class="stat-value">{{ formatNumber(statistics.view_count || 0) }}</div>
        <div class="stat-label">总浏览量</div>
      </div>
    </div>

    <div class="layui-card mb-4">
      <div class="layui-card-header d-flex align-center gap-3">
        <i class="fa-solid fa-bolt"></i>
        <span>快捷操作</span>
      </div>
      <div class="layui-card-body">
        <div class="action-grid">
          <router-link to="/" class="layui-btn layui-btn-normal">
            <i class="fa-solid fa-house mr-1"></i>返回首页
          </router-link>
          <router-link to="/create" class="layui-btn">
            <i class="fa-solid fa-plus mr-1"></i>创建文章
          </router-link>
          <router-link to="/admin/users" class="layui-btn layui-btn-warm">
            <i class="fa-solid fa-user mr-1"></i>用户管理
          </router-link>
          <router-link to="/admin/articles" class="layui-btn layui-btn-danger">
            <i class="fa-solid fa-file mr-1"></i>文章管理
          </router-link>
        </div>
      </div>
    </div>

    <div class="layui-card">
      <div class="layui-card-header d-flex align-center gap-3">
        <i class="fa-solid fa-circle-info"></i>
        <span>系统信息</span>
      </div>
      <div class="layui-card-body">
        <ul class="system-info-list">
          <li>
            <i class="fa-solid fa-tag"></i>
            <span class="info-title">系统版本</span>
            <span class="info-value">v1.3.14</span>
          </li>
          <li>
            <i class="fa-solid fa-database" style="color: #52C41A;"></i>
            <span class="info-title">数据库状态</span>
            <span class="info-value" style="color: #52C41A;">正常运行</span>
          </li>
          <li>
            <i class="fa-solid fa-server" style="color: #1E9FFF;"></i>
            <span class="info-title">服务器状态</span>
            <span class="info-value" style="color: #1E9FFF;">运行中</span>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'OverviewPanel',
  props: {
    statistics: { type: Object, default: () => ({}) },
    loading: { type: Boolean, default: false }
  },
  methods: {
    formatNumber(num) {
      return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',')
    }
  }
}
</script>

<style scoped>
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
  margin-bottom: 16px;
}

.stat-card {
  padding: 20px;
  border-radius: 12px;
  color: white;
}

.stat-card.primary {
  background: #1E9FFF;
}

.stat-card.pink {
  background: #E91E63;
}

.stat-card.cyan {
  background: #00BCD4;
}

.stat-card.green {
  background: #4CAF50;
}

.stat-icon {
  font-size: 32px;
  margin-bottom: 8px;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
}

.stat-label {
  font-size: 14px;
  opacity: 0.8;
}

.action-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 12px;
}

.system-info-list {
  list-style: none;
  padding: 0;
}

.system-info-list li {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 0;
  border-bottom: 1px solid #f0f0f0;
}

.system-info-list li:last-child {
  border-bottom: none;
}

.system-info-list li i {
  width: 24px;
  text-align: center;
  color: #999;
}

.info-title {
  flex: 1;
  color: #333;
}

.info-value {
  color: #666;
}

.d-flex {
  display: flex;
}

.align-center {
  align-items: center;
}

.gap-3 {
  gap: 12px;
}

.mr-1 {
  margin-right: 4px;
}

.mb-4 {
  margin-bottom: 16px;
}
</style>