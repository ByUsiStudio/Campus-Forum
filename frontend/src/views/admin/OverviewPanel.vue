<template>
  <div>
    <!-- 面板头部 -->
    <div class="d-flex justify-space-between align-center mb-4">
      <div class="d-flex align-center">
        <h2 class="text-h5 font-weight-bold">数据概览</h2>
        <v-chip size="small" variant="tonal" color="primary" class="ml-3">
          实时统计
        </v-chip>
      </div>
      <v-btn 
        variant="tonal" 
        color="primary" 
        @click="$emit('refresh')" 
        :loading="loading"
        prepend-icon="mdi-refresh"
      >
        刷新数据
      </v-btn>
    </div>

    <!-- 统计卡片网格 -->
    <v-row>
      <v-col cols="12" sm="6" lg="3">
        <v-card elevation="2" hover class="transition-swing">
          <v-card-item>
            <div class="d-flex justify-space-between align-start mb-2">
              <v-avatar color="primary" size="48" class="elevation-2">
                <v-icon color="white" size="24">mdi-account-group</v-icon>
              </v-avatar>
              <v-chip size="x-small" color="success" variant="tonal">+12%</v-chip>
            </div>
            <div class="mt-4">
              <div class="text-h4 font-weight-bold">{{ statistics.user_count || 0 }}</div>
              <div class="text-body-2 text-medium-emphasis mt-1">用户总数</div>
            </div>
          </v-card-item>
        </v-card>
      </v-col>

      <v-col cols="12" sm="6" lg="3">
        <v-card elevation="2" hover class="transition-swing">
          <v-card-item>
            <div class="d-flex justify-space-between align-start mb-2">
              <v-avatar color="pink" size="48" class="elevation-2">
                <v-icon color="white" size="24">mdi-file-document</v-icon>
              </v-avatar>
              <v-chip size="x-small" color="success" variant="tonal">+8%</v-chip>
            </div>
            <div class="mt-4">
              <div class="text-h4 font-weight-bold">{{ statistics.article_count || 0 }}</div>
              <div class="text-body-2 text-medium-emphasis mt-1">文章总数</div>
            </div>
          </v-card-item>
        </v-card>
      </v-col>

      <v-col cols="12" sm="6" lg="3">
        <v-card elevation="2" hover class="transition-swing">
          <v-card-item>
            <div class="d-flex justify-space-between align-start mb-2">
              <v-avatar color="cyan" size="48" class="elevation-2">
                <v-icon color="white" size="24">mdi-comment-text</v-icon>
              </v-avatar>
              <v-chip size="x-small" color="success" variant="tonal">+15%</v-chip>
            </div>
            <div class="mt-4">
              <div class="text-h4 font-weight-bold">{{ statistics.comment_count || 0 }}</div>
              <div class="text-body-2 text-medium-emphasis mt-1">评论总数</div>
            </div>
          </v-card-item>
        </v-card>
      </v-col>

      <v-col cols="12" sm="6" lg="3">
        <v-card elevation="2" hover class="transition-swing">
          <v-card-item>
            <div class="d-flex justify-space-between align-start mb-2">
              <v-avatar color="green" size="48" class="elevation-2">
                <v-icon color="white" size="24">mdi-eye</v-icon>
              </v-avatar>
              <v-chip size="x-small" color="success" variant="tonal">+23%</v-chip>
            </div>
            <div class="mt-4">
              <div class="text-h4 font-weight-bold">{{ formatNumber(statistics.view_count || 0) }}</div>
              <div class="text-body-2 text-medium-emphasis mt-1">总浏览量</div>
            </div>
          </v-card-item>
        </v-card>
      </v-col>
    </v-row>

    <!-- 内容区域 -->
    <v-row>
      <v-col cols="12" lg="8">
        <v-card elevation="2">
          <v-card-title class="d-flex align-center">
            <v-icon start size="20" color="primary">mdi-flash</v-icon>
            快捷操作
          </v-card-title>
          <v-card-text>
            <div class="d-flex flex-wrap ga-2">
              <v-btn 
                variant="tonal" 
                color="primary" 
                to="/"
                prepend-icon="mdi-home"
                class="mb-2"
              >
                返回首页
              </v-btn>
              <v-btn 
                variant="tonal" 
                color="success" 
                to="/create"
                prepend-icon="mdi-plus"
                class="mb-2"
              >
                创建文章
              </v-btn>
              <v-btn 
                variant="tonal" 
                color="info" 
                to="/admin/users"
                prepend-icon="mdi-account"
                class="mb-2"
              >
                用户管理
              </v-btn>
              <v-btn 
                variant="tonal" 
                color="warning" 
                to="/admin/articles"
                prepend-icon="mdi-file"
                class="mb-2"
              >
                文章管理
              </v-btn>
            </div>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="12" lg="4">
        <v-card elevation="2">
          <v-card-title class="d-flex align-center">
            <v-icon start size="20" color="primary">mdi-information</v-icon>
            系统信息
          </v-card-title>
          <v-card-text>
            <v-list density="compact">
              <v-list-item>
                <template v-slot:prepend>
                  <v-icon color="primary">mdi-version</v-icon>
                </template>
                <v-list-item-title>系统版本</v-list-item-title>
                <template v-slot:append>
                  <v-chip size="x-small" variant="tonal">v1.3.14</v-chip>
                </template>
              </v-list-item>
              <v-divider />
              <v-list-item>
                <template v-slot:prepend>
                  <v-icon color="success">mdi-database</v-icon>
                </template>
                <v-list-item-title>数据库</v-list-item-title>
                <template v-slot:append>
                  <v-chip size="x-small" color="success" variant="tonal">正常</v-chip>
                </template>
              </v-list-item>
              <v-divider />
              <v-list-item>
                <template v-slot:prepend>
                  <v-icon color="info">mdi-server</v-icon>
                </template>
                <v-list-item-title>服务器</v-list-item-title>
                <template v-slot:append>
                  <v-chip size="x-small" color="success" variant="tonal">运行中</v-chip>
                </template>
              </v-list-item>
            </v-list>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>

<script>
export default {
  name: 'OverviewPanel',
  props: {
    statistics: {
      type: Object,
      default: () => ({})
    },
    loading: {
      type: Boolean,
      default: false
    }
  },
  methods: {
    formatNumber(num) {
      return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',')
    }
  }
}
</script>

<style scoped>
/* 完全使用 Vuetify 内置样式 */
</style>
