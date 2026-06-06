<template>
  <div>
    <!-- 统计卡片 -->
    <v-row class="mb-4">
      <v-col cols="6" sm="6" md="3">
        <v-card class="pa-4" variant="flat" color="primary" rounded="lg">
          <div class="d-flex align-center mb-2">
            <v-icon size="32" color="white">mdi-account-group</v-icon>
          </div>
          <div class="text-h4 font-weight-bold text-white">{{ statistics.user_count || 0 }}</div>
          <div class="text-body-2 text-white opacity-80">用户总数</div>
        </v-card>
      </v-col>
      
      <v-col cols="6" sm="6" md="3">
        <v-card class="pa-4" variant="flat" color="pink" rounded="lg">
          <div class="d-flex align-center mb-2">
            <v-icon size="32" color="white">mdi-file-document</v-icon>
          </div>
          <div class="text-h4 font-weight-bold text-white">{{ statistics.article_count || 0 }}</div>
          <div class="text-body-2 text-white opacity-80">文章总数</div>
        </v-card>
      </v-col>
      
      <v-col cols="6" sm="6" md="3">
        <v-card class="pa-4" variant="flat" color="cyan" rounded="lg">
          <div class="d-flex align-center mb-2">
            <v-icon size="32" color="white">mdi-comment-text</v-icon>
          </div>
          <div class="text-h4 font-weight-bold text-white">{{ statistics.comment_count || 0 }}</div>
          <div class="text-body-2 text-white opacity-80">评论总数</div>
        </v-card>
      </v-col>
      
      <v-col cols="6" sm="6" md="3">
        <v-card class="pa-4" variant="flat" color="green" rounded="lg">
          <div class="d-flex align-center mb-2">
            <v-icon size="32" color="white">mdi-eye</v-icon>
          </div>
          <div class="text-h4 font-weight-bold text-white">{{ formatNumber(statistics.view_count || 0) }}</div>
          <div class="text-body-2 text-white opacity-80">总浏览量</div>
        </v-card>
      </v-col>
    </v-row>

    <!-- 快捷操作 -->
    <v-card class="mb-4" variant="flat" rounded="lg">
      <v-card-title class="pb-2">
        <v-icon start size="20">mdi-flash</v-icon>
        快捷操作
      </v-card-title>
      <v-card-text>
        <v-row dense>
          <v-col cols="6" sm="4" md="3">
            <v-btn variant="tonal" color="primary" to="/" block prepend-icon="mdi-home">返回首页</v-btn>
          </v-col>
          <v-col cols="6" sm="4" md="3">
            <v-btn variant="tonal" color="success" to="/create" block prepend-icon="mdi-plus">创建文章</v-btn>
          </v-col>
          <v-col cols="6" sm="4" md="3">
            <v-btn variant="tonal" color="info" to="/admin/users" block prepend-icon="mdi-account">用户管理</v-btn>
          </v-col>
          <v-col cols="6" sm="4" md="3">
            <v-btn variant="tonal" color="warning" to="/admin/articles" block prepend-icon="mdi-file">文章管理</v-btn>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>

    <!-- 系统信息 -->
    <v-card variant="flat" rounded="lg">
      <v-card-title class="pb-2">
        <v-icon start size="20">mdi-information</v-icon>
        系统信息
      </v-card-title>
      <v-card-text>
        <v-list density="compact">
          <v-list-item>
            <template v-slot:prepend><v-icon>mdi-tag</v-icon></template>
            <v-list-item-title>系统版本</v-list-item-title>
            <v-list-item-subtitle>v1.3.14</v-list-item-subtitle>
          </v-list-item>
          <v-list-item>
            <template v-slot:prepend><v-icon color="success">mdi-database</v-icon></template>
            <v-list-item-title>数据库状态</v-list-item-title>
            <v-list-item-subtitle class="text-success">正常运行</v-list-item-subtitle>
          </v-list-item>
          <v-list-item>
            <template v-slot:prepend><v-icon color="info">mdi-server</v-icon></template>
            <v-list-item-title>服务器状态</v-list-item-title>
            <v-list-item-subtitle class="text-info">运行中</v-list-item-subtitle>
          </v-list-item>
        </v-list>
      </v-card-text>
    </v-card>
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