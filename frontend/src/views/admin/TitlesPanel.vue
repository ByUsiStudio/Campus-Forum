<template>
  <div>
    <div class="layui-card mb-4">
      <div class="layui-card-header">
        <i class="fa-solid fa-plus-circle mr-2" style="font-size: 20px;"></i>
        添加新头衔
      </div>
      <div class="layui-card-body">
        <div class="layui-row">
          <div class="layui-col-xs12 layui-col-sm6">
            <div class="layui-form-item">
              <label class="layui-form-label">头衔名称</label>
              <div class="layui-input-block">
                <div class="layui-input-wrap">
                  <i class="fa-solid fa-star layui-input-prefix" style="color: #999;"></i>
                  <input type="text" v-model="titleForm.name" placeholder="例如：技术达人" class="layui-input" />
                </div>
              </div>
            </div>
          </div>
          <div class="layui-col-xs12 layui-col-sm6">
            <div class="layui-form-item">
              <label class="layui-form-label">描述</label>
              <div class="layui-input-block">
                <div class="layui-input-wrap">
                  <i class="fa-solid fa-file-text layui-input-prefix" style="color: #999;"></i>
                  <input type="text" v-model="titleForm.description" placeholder="头衔简介" class="layui-input" />
                </div>
              </div>
            </div>
          </div>
          <div class="layui-col-xs6 layui-col-sm3">
            <div class="layui-form-item">
              <label class="layui-form-label">颜色</label>
              <div class="layui-input-block">
                <input type="text" v-model="titleForm.color" class="layui-input" />
              </div>
            </div>
          </div>
          <div class="layui-col-xs6 layui-col-sm3">
            <div class="layui-form-item">
              <label class="layui-form-label">图标</label>
              <div class="layui-input-block">
                <input type="text" v-model="titleForm.icon" placeholder="fa-solid fa-star" class="layui-input" />
              </div>
            </div>
          </div>
          <div class="layui-col-xs12 layui-col-sm6">
            <div class="d-flex align-items-center justify-center py-2">
              <span class="title-chip" :style="{ backgroundColor: adjustColorOpacity(titleForm.color, 0.12), color: titleForm.color, border: `1px solid ${adjustColorOpacity(titleForm.color, 0.25)}`, padding: '4px 12px', borderRadius: '20px' }">
                <i v-if="titleForm.icon" :class="titleForm.icon" class="mr-1" style="font-size: 18px;"></i>
                {{ titleForm.name || '预览' }}
              </span>
            </div>
          </div>
          <div class="layui-col-xs12 layui-col-sm6">
            <button class="layui-btn layui-btn-primary layui-btn-fluid" style="height: 40px;" @click="$emit('add-title', titleForm)">
              <i class="fa-solid fa-plus mr-2"></i>
              添加头衔
            </button>
          </div>
        </div>
      </div>
    </div>

    <div class="layui-card mb-4">
      <div class="layui-card-header">
        <i class="fa-solid fa-gift mr-2" style="font-size: 20px;"></i>
        授予头衔
      </div>
      <div class="layui-card-body">
        <div class="layui-row">
          <div class="layui-col-xs12 layui-col-sm4">
            <div class="layui-form-item">
              <label class="layui-form-label">选择用户</label>
              <div class="layui-input-block">
                <select v-model="grantForm.user_id" class="layui-select">
                  <option value="">请选择用户</option>
                  <option v-for="user in users" :key="user.id" :value="user.id">{{ user.display_name || user.username }}</option>
                </select>
              </div>
            </div>
          </div>
          <div class="layui-col-xs12 layui-col-sm4">
            <div class="layui-form-item">
              <label class="layui-form-label">选择头衔</label>
              <div class="layui-input-block">
                <select v-model="grantForm.title_id" class="layui-select">
                  <option value="">请选择头衔</option>
                  <option v-for="title in titles" :key="title.id" :value="title.id">{{ title.name }}</option>
                </select>
              </div>
            </div>
          </div>
          <div class="layui-col-xs12 layui-col-sm4">
            <button class="layui-btn layui-btn-fluid" style="height: 40px; margin-top: 30px;" @click="$emit('grant', grantForm)">
              <i class="fa-solid fa-check mr-2"></i>
              授予头衔
            </button>
          </div>
        </div>
      </div>
    </div>

    <div class="layui-card mb-4">
      <div class="layui-card-header">
        <i class="fa-solid fa-medal mr-2" style="font-size: 20px;"></i>
        已创建的头衔
      </div>
      <div v-if="titles.length > 0">
        <div v-for="title in titles" :key="title.id" class="layui-card-body border-b last:border-b-0" style="padding: 12px 15px;">
          <div class="flex items-center gap-3">
            <div class="avatar" :style="{ backgroundColor: title.color || '#1E9FFF', width: '48px', height: '48px', borderRadius: '50%', display: 'flex', alignItems: 'center', justifyContent: 'center' }">
              <i v-if="title.icon" :class="title.icon" style="color: white;"></i>
              <i v-else class="fa-solid fa-medal" style="color: white;"></i>
            </div>
            <div class="flex-1">
              <span class="title-chip" :style="{ backgroundColor: adjustColorOpacity(title.color, 0.12), color: title.color, border: `1px solid ${adjustColorOpacity(title.color, 0.25)}`, padding: '2px 8px', borderRadius: '12px', fontSize: '13px' }">
                <i v-if="title.icon" :class="title.icon" class="mr-1" style="font-size: 14px;"></i>
                {{ title.name }}
              </span>
              <div class="text-muted text-sm mt-1">{{ title.description || '暂无描述' }}</div>
            </div>
            <button class="layui-btn layui-btn-danger layui-btn-sm" @click="$emit('delete-title', title)">
              <i class="fa-solid fa-trash"></i>
            </button>
          </div>
        </div>
      </div>
      <div v-else class="layui-card-body text-center py-8">
        <i class="fa-solid fa-medal" style="font-size: 48px; color: #dcdcdc;"></i>
        <div class="text-muted mt-2">暂无头衔数据</div>
      </div>
    </div>

    <div class="layui-card">
      <div class="layui-card-header">
        <i class="fa-solid fa-users mr-2" style="font-size: 20px;"></i>
        用户头衔一览
      </div>
      <div v-if="usersWithTitles.length > 0">
        <div v-for="user in usersWithTitles" :key="user.id" class="layui-card-body border-b last:border-b-0" style="padding: 12px 15px;">
          <div class="flex items-center gap-3">
            <UserAvatar :user="user" :size="48" />
            <div class="flex-1">
              <div class="font-medium">{{ user.display_name }}<span class="text-muted text-sm ml-2">@{{ user.username }}</span></div>
              <div class="flex flex-wrap gap-1 mt-1">
                <span 
                  v-for="userTitle in user.titles" 
                  :key="userTitle.id"
                  class="title-chip inline-flex items-center"
                  :style="{ backgroundColor: adjustColorOpacity(userTitle.color, 0.12), color: userTitle.color, border: `1px solid ${adjustColorOpacity(userTitle.color, 0.25)}`, padding: '2px 8px', borderRadius: '12px', fontSize: '12px' }"
                >
                  <i v-if="userTitle.icon" :class="userTitle.icon" class="mr-1" style="font-size: 12px;"></i>
                  {{ userTitle.name }}
                  <button 
                    class="ml-1"
                    style="background: none; border: none; cursor: pointer; color: #FF4D4F;"
                    @click="$emit('revoke', user.id, userTitle.id)"
                  >
                    <i class="fa-solid fa-xmark" style="font-size: 12px;"></i>
                  </button>
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div v-else class="layui-card-body text-center py-8">
        <i class="fa-solid fa-user-slash" style="font-size: 48px; color: #dcdcdc;"></i>
        <div class="text-muted mt-2">暂无用户头衔数据</div>
      </div>
    </div>
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

    const adjustColorOpacity = (color, opacity) => {
      if (!color) return 'rgba(0,0,0,0.1)'
      const hex = color.replace('#', '')
      const r = parseInt(hex.substring(0, 2), 16)
      const g = parseInt(hex.substring(2, 4), 16)
      const b = parseInt(hex.substring(4, 6), 16)
      return `rgba(${r}, ${g}, ${b}, ${opacity})`
    }

    return {
      titleForm,
      grantForm,
      usersWithTitles,
      adjustColorOpacity
    }
  }
}
</script>

<style scoped>
.last\:border-b-0:last-child {
  border-bottom: none;
}
</style>