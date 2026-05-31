<template>
  <div class="titles-panel">
    <div class="panel-header">
      <div class="header-left">
        <h2 class="panel-title">头衔管理</h2>
        <p class="panel-subtitle">创建头衔并授予给用户</p>
      </div>
      <div class="header-actions">
        <v-btn variant="outlined" color="primary" @click="$emit('refresh')" :loading="loading" size="small">
          <v-icon start size="18">mdi-refresh</v-icon>
          <span class="d-none d-sm-inline">刷新</span>
        </v-btn>
      </div>
    </div>

    <v-row>
      <v-col cols="12" lg="6">
        <v-card class="form-card">
          <v-card-title class="form-title">
            <v-icon class="title-icon">mdi-plus-circle</v-icon>
            添加新头衔
          </v-card-title>
          <v-card-text>
            <v-text-field
              v-model="titleForm.name"
              label="头衔名称"
              placeholder="例如：技术达人"
              variant="outlined"
              density="comfortable"
              prepend-inner-icon="mdi-star"
              class="mb-3"
              hide-details="auto"
            ></v-text-field>
            <v-text-field
              v-model="titleForm.description"
              label="描述"
              placeholder="头衔简介"
              variant="outlined"
              density="comfortable"
              prepend-inner-icon="mdi-text"
              class="mb-3"
              hide-details="auto"
            ></v-text-field>
            <v-row no-gutters>
              <v-col cols="5">
                <v-text-field
                  v-model="titleForm.color"
                  label="颜色"
                  variant="outlined"
                  density="comfortable"
                  hide-details="auto"
                >
                </v-text-field>
              </v-col>
              <v-col cols="7" class="pl-2">
                <v-text-field
                  v-model="titleForm.icon"
                  label="图标"
                  placeholder="mdi-star"
                  variant="outlined"
                  density="comfortable"
                  hide-details="auto"
                ></v-text-field>
              </v-col>
            </v-row>
            <div class="color-preview my-3">
              <v-chip :color="titleForm.color" size="medium">
                <v-icon v-if="titleForm.icon" start size="18">{{ titleForm.icon }}</v-icon>
                {{ titleForm.name || '预览' }}
              </v-chip>
            </div>
            <v-btn color="primary" block size="large" @click="$emit('add-title', titleForm)" class="mt-2">
              <v-icon start>mdi-plus</v-icon>
              添加头衔
            </v-btn>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="12" lg="6">
        <v-card class="form-card">
          <v-card-title class="form-title">
            <v-icon class="title-icon">mdi-gift</v-icon>
            授予头衔
          </v-card-title>
          <v-card-text>
            <v-select
              v-model="grantForm.user_id"
              :items="users"
              item-title="display_name"
              item-value="id"
              label="选择用户"
              variant="outlined"
              density="comfortable"
              class="mb-3"
              hide-details="auto"
            ></v-select>
            <v-select
              v-model="grantForm.title_id"
              :items="titles"
              item-title="name"
              item-value="id"
              label="选择头衔"
              variant="outlined"
              density="comfortable"
              class="mb-3"
              hide-details="auto"
            ></v-select>
            <v-text-field
              v-model="grantForm.reason"
              label="授予原因"
              placeholder="为什么授予此头衔"
              variant="outlined"
              density="comfortable"
              class="mb-3"
              hide-details="auto"
            ></v-text-field>
            <v-btn color="success" block size="large" @click="$emit('grant', grantForm)">
              <v-icon start>mdi-check</v-icon>
              授予头衔
            </v-btn>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <v-card class="mt-6">
      <v-card-title class="form-title">
        <v-icon class="title-icon">mdi-medal</v-icon>
        已创建的头衔
      </v-card-title>
      <v-card-text>
        <div v-if="titles.length === 0" class="empty-state">
          <v-icon size="48" color="grey">mdi-medal-outline</v-icon>
          <div class="empty-text">暂无头衔数据</div>
        </div>

        <div v-else class="titles-list">
          <div v-for="title in titles" :key="title.id" class="title-item">
            <div class="title-item-left">
              <v-chip :color="title.color" size="medium">
                <v-icon v-if="title.icon" start size="16">{{ title.icon }}</v-icon>
                {{ title.name }}
              </v-chip>
            </div>
            <div class="title-info">
              <div class="title-name d-none d-sm-block">{{ title.name }}</div>
              <div class="title-desc">{{ title.description || '无描述' }}</div>
            </div>
            <v-btn variant="text" size="small" color="error" @click="$emit('delete-title', title)">
              <v-icon size="18">mdi-delete</v-icon>
              <span class="d-none d-sm-inline ml-1">删除</span>
            </v-btn>
          </div>
        </div>
      </v-card-text>
    </v-card>

    <v-card class="mt-6" variant="outlined">
      <v-card-title class="form-title">
        <v-icon class="title-icon">mdi-account-multiple</v-icon>
        用户头衔一览
      </v-card-title>
      <v-card-text>
        <div class="user-titles-container">
          <div v-if="usersWithTitles.length === 0" class="empty-state">
            <v-icon size="48" color="grey">mdi-account-off</v-icon>
            <div class="empty-text">暂无用户头衔数据</div>
          </div>

          <div v-else class="user-titles-list">
            <div v-for="user in usersWithTitles" :key="user.id" class="user-title-card">
              <div class="user-header">
                <UserAvatar :user="user" :size="36" />
                <div class="user-info">
                  <div class="user-name">{{ user.display_name }}</div>
                  <div class="user-username">@{{ user.username }}</div>
                </div>
              </div>
              <div class="user-titles-chips">
                <v-chip
                  v-for="userTitle in user.titles"
                  :key="userTitle.id"
                  :color="userTitle.color"
                  size="small"
                  class="mr-1 mb-1"
                >
                  <v-icon v-if="userTitle.icon" start size="14">{{ userTitle.icon }}</v-icon>
                  {{ userTitle.name }}
                  <v-btn
                    variant="text"
                    size="x-small"
                    icon
                    color="error"
                    @click="$emit('revoke', user.id, userTitle.id)"
                    class="ml-1"
                  >
                    <v-icon size="14">mdi-close-circle</v-icon>
                  </v-btn>
                </v-chip>
              </div>
            </div>
          </div>
        </div>
      </v-card-text>
    </v-card>
  </div>
</template>

<script>
import { computed, reactive } from 'vue'
import UserAvatar from '../../components/UserAvatar.vue'

export default {
  name: 'TitlesPanel',
  components: {
    UserAvatar
  },
  props: {
    titles: {
      type: Array,
      default: () => []
    },
    users: {
      type: Array,
      default: () => []
    },
    loading: {
      type: Boolean,
      default: false
    }
  },
  emits: ['add-title', 'grant', 'revoke', 'delete-title', 'refresh'],
  setup(props) {
    const titleForm = reactive({
      name: '',
      description: '',
      color: '#6750A4',
      icon: ''
    })

    const grantForm = reactive({
      user_id: null,
      title_id: null,
      reason: ''
    })

    const usersWithTitles = computed(() => {
      return props.users.filter(user => user.titles && user.titles.length > 0)
    })

    return {
      titleForm,
      grantForm,
      usersWithTitles
    }
  }
}
</script>

<style scoped>
.titles-panel {
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 16px;
}

.panel-title {
  font-size: 1.5rem;
  font-weight: 700;
  color: #1a1a2e;
  margin: 0 0 4px 0;
}

.panel-subtitle {
  font-size: 0.9rem;
  color: #6b7280;
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 8px;
}

.form-card {
  border-radius: 16px;
  height: 100%;
}

.form-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 1rem;
  font-weight: 600;
  padding: 16px 20px !important;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
}

.title-icon {
  color: #6750A4 !important;
}

.color-preview {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 40px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 40px 20px;
}

.empty-text {
  margin-top: 12px;
  font-size: 0.9rem;
  color: #9ca3af;
}

.titles-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.title-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: #f8f9ff;
  border-radius: 10px;
  transition: background 0.2s;
  flex-wrap: wrap;
}

.title-item:hover {
  background: #f0f2ff;
}

.title-item-left {
  flex-shrink: 0;
}

.title-info {
  flex: 1;
  min-width: 0;
  min-width: 100px;
}

.title-name {
  font-weight: 600;
  color: #1a1a2e;
}

.title-desc {
  font-size: 0.85rem;
  color: #6b7280;
}

.user-titles-table {
  overflow-x: auto;
}

.user-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.user-titles-container {
  width: 100%;
}

.user-titles-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.user-title-card {
  padding: 12px 16px;
  background: #f8f9ff;
  border-radius: 12px;
  border: 1px solid rgba(103, 80, 164, 0.1);
}

.user-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 10px;
}

.user-info {
  flex: 1;
  min-width: 0;
}

.user-name {
  font-weight: 600;
  color: #1a1a2e;
  font-size: 0.95rem;
}

.user-username {
  font-size: 0.8rem;
  color: #6b7280;
}

.user-titles-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

@media (max-width: 599px) {
  .panel-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .header-actions {
    width: 100%;
    justify-content: flex-end;
  }

  .form-card {
    height: auto;
  }

  .form-title {
    font-size: 0.9rem;
    padding: 12px 16px !important;
  }

  .title-item {
    flex-direction: column;
    align-items: flex-start;
  }

  .title-item-left {
    width: 100%;
  }

  .title-info {
    width: 100%;
    margin: 8px 0;
  }

  .user-title-card {
    padding: 10px 12px;
  }
}

@media (max-width: 960px) {
  .form-card {
    height: auto;
  }
}

@media (min-width: 600px) and (max-width: 960px) {
  .title-item {
    flex-wrap: nowrap;
  }

  .title-info {
    min-width: 150px;
  }
}
</style>