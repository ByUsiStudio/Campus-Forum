<template>
  <div>
    <!-- 添加新头衔表单 -->
    <el-card class="mb-4" shadow="never" body-style="padding: 20px">
      <template #header>
        <div class="panel-header">
          <el-icon :size="20" class="header-icon"><Plus /></el-icon>
          <span>添加新头衔</span>
        </div>
      </template>

      <el-row :gutter="16">
        <el-col :xs="24" :sm="12">
          <el-form label-position="top">
            <el-form-item label="头衔名称">
              <el-input
                v-model="titleForm.name"
                placeholder="例如：技术达人"
                clearable
              >
                <template #prefix>
                  <el-icon><Star /></el-icon>
                </template>
              </el-input>
            </el-form-item>
          </el-form>
        </el-col>
        <el-col :xs="24" :sm="12">
          <el-form label-position="top">
            <el-form-item label="描述">
              <el-input
                v-model="titleForm.description"
                placeholder="头衔简介"
                clearable
              >
                <template #prefix>
                  <el-icon><EditPen /></el-icon>
                </template>
              </el-input>
            </el-form-item>
          </el-form>
        </el-col>
        <el-col :xs="12" :sm="6">
          <el-form label-position="top">
            <el-form-item label="颜色">
              <el-input v-model="titleForm.color" clearable />
            </el-form-item>
          </el-form>
        </el-col>
        <el-col :xs="12" :sm="6">
          <el-form label-position="top">
            <el-form-item label="图标">
              <el-input v-model="titleForm.icon" placeholder="Plus / Star" clearable />
            </el-form-item>
          </el-form>
        </el-col>
        <el-col :xs="24" :sm="6">
          <div class="preview-cell">
            <el-tag
              :color="titleForm.color"
              effect="dark"
              style="border: none; color: #fff"
            >
              <el-icon v-if="titleForm.icon" class="mr-1" :style="{ verticalAlign: '-0.2em' }">
                <component :is="titleForm.icon" />
              </el-icon>
              {{ titleForm.name || '预览' }}
            </el-tag>
          </div>
        </el-col>
        <el-col :xs="24" :sm="6">
          <el-button type="primary" style="width: 100%; height: 40px" @click="$emit('add-title', titleForm)">
            <el-icon class="mr-1"><Plus /></el-icon>
            添加头衔
          </el-button>
        </el-col>
      </el-row>
    </el-card>

    <!-- 授予头衔表单 -->
    <el-card class="mb-4" shadow="never" body-style="padding: 20px">
      <template #header>
        <div class="panel-header">
          <el-icon :size="20" class="header-icon"><Gift /></el-icon>
          <span>授予头衔</span>
        </div>
      </template>

      <el-row :gutter="16">
        <el-col :xs="24" :sm="8">
          <el-form label-position="top">
            <el-form-item label="选择用户">
              <el-select
                v-model="grantForm.user_id"
                placeholder="请选择用户"
                filterable
                clearable
                style="width: 100%"
              >
                <el-option
                  v-for="user in users"
                  :key="user.id"
                  :label="user.display_name"
                  :value="user.id"
                >
                  <div class="d-flex align-center">
                    <UserAvatar :user="user" :size="28" />
                    <span class="ml-2">@{{ user.username }}</span>
                  </div>
                </el-option>
              </el-select>
            </el-form-item>
          </el-form>
        </el-col>
        <el-col :xs="24" :sm="8">
          <el-form label-position="top">
            <el-form-item label="选择头衔">
              <el-select
                v-model="grantForm.title_id"
                :placeholder="'请选择头衔'"
                filterable
                clearable
                style="width: 100%"
              >
                <el-option
                  v-for="title in titles"
                  :key="title.id"
                  :label="title.name"
                  :value="title.id"
                />
              </el-select>
            </el-form-item>
          </el-form>
        </el-col>
        <el-col :xs="24" :sm="8">
          <el-button type="success" style="width: 100%; height: 40px" @click="$emit('grant', grantForm)">
            <el-icon class="mr-1"><Check /></el-icon>
            授予头衔
          </el-button>
        </el-col>
      </el-row>
    </el-card>

    <!-- 已创建的头衔列表 -->
    <el-card class="mb-4" shadow="never" body-style="padding: 20px">
      <template #header>
        <div class="panel-header">
          <el-icon :size="20" class="header-icon"><Medal /></el-icon>
          <span>已创建的头衔</span>
        </div>
      </template>

      <div
        v-if="titles.length > 0"
        v-loading="loading"
        class="title-list"
      >
        <div v-for="title in titles" :key="title.id" class="title-list-item">
          <div class="title-avatar" :style="{ backgroundColor: title.color || 'var(--el-color-primary)' }">
            <el-icon v-if="title.icon" :size="24" style="color: #fff">
              <component :is="title.icon" />
            </el-icon>
            <el-icon v-else :size="24" style="color: #fff">
              <Medal />
            </el-icon>
          </div>

          <div class="title-info">
            <div class="title-name">
              <el-tag
                :color="title.color"
                effect="dark"
                size="small"
                style="border: none; color: #fff"
              >
                <el-icon v-if="title.icon" class="mr-1" :style="{ verticalAlign: '-0.2em' }">
                  <component :is="title.icon" />
                </el-icon>
                {{ title.name }}
              </el-tag>
            </div>
            <div class="title-desc">{{ title.description || '暂无描述' }}</div>
          </div>

          <div class="title-action">
            <el-button size="small" type="danger" text @click="$emit('delete-title', title)">
              <el-icon :size="16"><Delete /></el-icon>
              <span class="ml-1">删除</span>
            </el-button>
          </div>
        </div>
      </div>

      <el-empty v-else-if="!loading" :image-size="64" description="暂无头衔数据" />

      <div v-else v-loading="loading" style="min-height: 80px" />
    </el-card>

    <!-- 用户头衔一览 -->
    <el-card shadow="never" body-style="padding: 20px">
      <template #header>
        <div class="panel-header">
          <el-icon :size="20" class="header-icon"><User /></el-icon>
          <span>用户头衔一览</span>
        </div>
      </template>

      <div v-if="usersWithTitles.length > 0" class="user-title-list">
        <div v-for="user in usersWithTitles" :key="user.id" class="user-title-item">
          <UserAvatar :user="user" :size="48" />

          <div class="user-title-info">
            <div class="user-title-name">
              {{ user.display_name }}
              <span class="user-title-username">@{{ user.username }}</span>
            </div>
            <div class="user-title-chips">
              <el-tag
                v-for="userTitle in user.titles"
                :key="userTitle.id"
                :color="userTitle.color"
                effect="dark"
                size="small"
                class="mr-1 mb-1"
                style="border: none; color: #fff"
              >
                <el-icon v-if="userTitle.icon" class="mr-1" :style="{ verticalAlign: '-0.2em' }">
                  <component :is="userTitle.icon" />
                </el-icon>
                {{ userTitle.name }}
                <el-button
                  type="danger"
                  text
                  size="small"
                  class="revoke-btn"
                  @click="$emit('revoke', user.id, userTitle.id)"
                >
                  <el-icon :size="14"><CircleClose /></el-icon>
                </el-button>
              </el-tag>
            </div>
          </div>
        </div>
      </div>

      <el-empty v-else :image-size="64" description="暂无用户头衔数据" />
    </el-card>
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
.panel-header {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 1.1rem;
  font-weight: 700;
}

.header-icon {
  color: var(--el-color-primary);
}

.preview-cell {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 8px;
  height: 100%;
}

.title-list {
  display: flex;
  flex-direction: column;
}

.title-list-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 0;
  border-bottom: 1px solid var(--el-border-color-lighter);
}

.title-list-item:last-child {
  border-bottom: none;
}

.title-avatar {
  flex-shrink: 0;
  width: 48px;
  height: 48px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.title-info {
  flex: 1;
  min-width: 0;
}

.title-name {
  font-weight: 600;
  margin-bottom: 4px;
}

.title-desc {
  color: var(--el-text-color-secondary);
  font-size: 0.875rem;
}

.title-action {
  flex-shrink: 0;
}

.user-title-list {
  display: flex;
  flex-direction: column;
}

.user-title-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 12px 0;
  border-bottom: 1px solid var(--el-border-color-lighter);
}

.user-title-item:last-child {
  border-bottom: none;
}

.user-title-info {
  flex: 1;
  min-width: 0;
}

.user-title-name {
  font-weight: 600;
  margin-bottom: 4px;
}

.user-title-username {
  color: var(--el-text-color-secondary);
  font-size: 0.75rem;
  margin-left: 8px;
}

.user-title-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.revoke-btn {
  margin-left: 6px;
  padding: 0 2px;
}

.mb-1 { margin-bottom: 4px; }
.mr-1 { margin-right: 4px; }
.mb-4 { margin-bottom: 16px; }
</style>
