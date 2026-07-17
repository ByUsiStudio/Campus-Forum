<template>
  <div>
    <div class="layui-card mb-4">
      <div class="layui-card-header">
        <i class="fa-solid fa-plus-circle mr-2" style="font-size: 20px;"></i>
        添加新分区
      </div>
      <div class="layui-card-body">
        <div class="layui-row">
          <div class="layui-col-xs12 layui-col-sm4">
            <div class="layui-form-item">
              <label class="layui-form-label">分区名称</label>
              <div class="layui-input-block">
                <input type="text" v-model="formData.name" placeholder="例如：技术交流" class="layui-input" />
              </div>
            </div>
          </div>
          <div class="layui-col-xs12 layui-col-sm5">
            <div class="layui-form-item">
              <label class="layui-form-label">描述</label>
              <div class="layui-input-block">
                <input type="text" v-model="formData.description" placeholder="分区简介" class="layui-input" />
              </div>
            </div>
          </div>
          <div class="layui-col-xs12 layui-col-sm3">
            <button class="layui-btn layui-btn-fluid" style="height: 40px; margin-top: 30px;" @click="handleAdd">
              <i class="fa-solid fa-plus mr-2"></i>
              添加
            </button>
          </div>
        </div>
      </div>
    </div>

    <div class="layui-card">
      <div v-if="categories.length > 0">
        <div v-for="cat in categories" :key="cat.id" class="layui-card-body border-b last:border-b-0" style="padding: 12px 15px;">
          <div class="flex items-center gap-3">
            <div class="avatar" :style="{ backgroundColor: cat.color || '#1E9FFF', width: '48px', height: '48px', borderRadius: '50%', display: 'flex', alignItems: 'center', justifyContent: 'center' }">
              <span style="color: white; font-size: 18px; font-weight: 600;">{{ cat.sort_order || 0 }}</span>
            </div>
            <div class="flex-1">
              <div class="font-medium">{{ cat.name }}</div>
              <div class="text-muted text-sm mt-1">{{ cat.description || '暂无描述' }}</div>
            </div>
            <div class="flex gap-2">
              <button class="layui-btn layui-btn-sm" @click="$emit('edit', cat)">
                <i class="fa-solid fa-pencil"></i>
              </button>
              <button class="layui-btn layui-btn-danger layui-btn-sm" @click="$emit('delete', cat)">
                <i class="fa-solid fa-trash"></i>
              </button>
            </div>
          </div>
        </div>
      </div>
      <div v-else class="layui-card-body text-center py-8">
        <i class="fa-solid fa-folder-open" style="font-size: 48px; color: #dcdcdc;"></i>
        <div class="text-muted mt-2">暂无分区数据</div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'

export default {
  name: 'CategoriesPanel',
  props: {
    categories: {
      type: Array,
      default: () => []
    },
    loading: {
      type: Boolean,
      default: false
    }
  },
  emits: ['add', 'edit', 'delete', 'refresh'],
  setup(props, { emit }) {
    const formData = ref({
      name: '',
      description: '',
      color: '#6750A4'
    })

    const handleAdd = () => {
      if (!formData.value.name) return
      emit('add', formData.value)
      formData.value = {
        name: '',
        description: '',
        color: '#6750A4'
      }
    }

    return {
      formData,
      handleAdd
    }
  }
}
</script>

<style scoped>
.last\:border-b-0:last-child {
  border-bottom: none;
}
</style>