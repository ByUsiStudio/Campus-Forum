<template>
  <div>
    <!-- 添加新分区表单 -->
    <el-card class="mb-4" shadow="never" body-class="categories-panel-card">
      <div slot="header" class="categories-panel-card__header">
        <el-icon :size="20"><Plus /></el-icon>
        <span>添加新分区</span>
      </div>
      <el-form :model="formData" inline class="categories-panel-form">
        <el-form-item label="分区名称" class="form-item-narrow">
          <el-input
            v-model="formData.name"
            placeholder="例如：技术交流"
            clearable
          />
        </el-form-item>
        <el-form-item label="描述" class="form-item-wide">
          <el-input
            v-model="formData.description"
            placeholder="分区简介"
            clearable
          />
        </el-form-item>
        <el-form-item class="form-item-btn">
          <el-button type="primary" @click="handleAdd">
            <el-icon class="mr-1"><Plus /></el-icon>
            添加
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 分区列表 -->
    <el-card shadow="never" body-class="categories-panel-card">
      <el-empty
        v-if="categories.length === 0"
        :image-size="120"
        description="暂无分区数据"
      />
      <el-list v-else class="categories-panel-list">
        <el-list-item v-for="cat in categories" :key="cat.id" class="py-3">
          <div class="d-flex align-center">
            <el-avatar
              :size="48"
              :style="avatarStyle(cat.color)"
              class="mr-3"
            >
              <span>{{ cat.sort_order || 0 }}</span>
            </el-avatar>

            <div class="flex-grow-1">
              <div class="font-weight-medium mb-1">{{ cat.name }}</div>
              <div class="list-item-subtitle">
                {{ cat.description || '暂无描述' }}
              </div>
            </div>

            <div class="flex-shrink-0">
              <el-tooltip content="编辑" :show-after="300">
                <el-button size="small" type="primary" text @click="$emit('edit', cat)">
                  <el-icon><Edit /></el-icon>
                </el-button>
              </el-tooltip>
              <el-tooltip content="删除" :show-after="300">
                <el-button size="small" type="danger" text @click="$emit('delete', cat)">
                  <el-icon><Delete /></el-icon>
                </el-button>
              </el-tooltip>
            </div>
          </div>
        </el-list-item>
      </el-list>
    </el-card>
  </div>
</template>

<script>
import { ref } from 'vue'
import { Plus, Edit, Delete } from '@element-plus/icons-vue'

export default {
  name: 'CategoriesPanel',
  components: {},
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

    const avatarStyle = (color) => {
      const c = color || '#6750A4'
      return {
        backgroundColor: c + '22',
        color: c,
        fontWeight: 600
      }
    }

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
      avatarStyle,
      handleAdd
    }
  }
}
</script>

<style scoped>
.categories-panel-card {
  padding: 4px;
}

.categories-panel-card__header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 1.05rem;
  font-weight: 600;
}

.categories-panel-form {
  display: flex;
  flex-wrap: wrap;
  align-items: flex-start;
  gap: 4px 16px;
}

.categories-panel-form .el-form-item {
  margin-bottom: 4px;
}

.form-item-narrow {
  width: 240px;
}

.form-item-wide {
  width: 300px;
}

.categories-panel-list {
  padding: 0;
}

.list-item-subtitle {
  color: var(--el-text-color-secondary);
  font-size: 0.9rem;
}

.font-weight-medium {
  font-weight: 500;
}
</style>
