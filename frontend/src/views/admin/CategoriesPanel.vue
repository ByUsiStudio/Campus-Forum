<template>
  <div class="categories-panel">
    <div class="panel-header">
      <div class="header-left">
        <h2 class="panel-title">分区管理</h2>
        <p class="panel-subtitle">添加、编辑与删除论坛分区</p>
      </div>
      <v-btn variant="outlined" color="primary" @click="$emit('refresh')" :loading="loading">
        <v-icon start>mdi-refresh</v-icon>
        刷新
      </v-btn>
    </div>

    <v-card class="add-form-card">
      <v-card-title class="form-title">
        <v-icon class="title-icon">mdi-plus-circle</v-icon>
        添加新分区
      </v-card-title>
      <v-card-text>
        <v-row>
          <v-col cols="12" sm="4">
            <v-text-field
              v-model="formData.name"
              label="分区名称"
              placeholder="例如：技术交流"
              variant="outlined"
              density="comfortable"
              hide-details
            ></v-text-field>
          </v-col>
          <v-col cols="12" sm="5">
            <v-text-field
              v-model="formData.description"
              label="描述"
              placeholder="分区简介"
              variant="outlined"
              density="comfortable"
              hide-details
            ></v-text-field>
          </v-col>
          <v-col cols="12" sm="3" class="d-flex align-center">
            <v-btn color="primary" size="large" block @click="handleAdd">
              <v-icon start>mdi-plus</v-icon>
              添加
            </v-btn>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>

    <div v-if="categories.length === 0" class="empty-state">
      <v-icon size="64" color="grey-lighten-1">mdi-folder-outline</v-icon>
      <div class="empty-text">暂无分区数据</div>
    </div>

    <div v-else class="categories-grid">
      <v-card
        v-for="cat in categories"
        :key="cat.id"
        class="category-card"
        variant="outlined"
      >
        <div class="category-main">
          <div class="category-icon" :style="{ backgroundColor: cat.color || '#6750A4' }">
            {{ cat.sort_order || 0 }}
          </div>

          <div class="category-content">
            <div class="category-name">{{ cat.name }}</div>
            <div class="category-desc">{{ cat.description || '暂无描述' }}</div>
          </div>

          <div class="category-actions">
            <v-btn variant="text" size="small" color="primary" @click="$emit('edit', cat)">
              <v-icon start size="18">mdi-pencil</v-icon>
              编辑
            </v-btn>
            <v-btn variant="text" size="small" color="error" @click="$emit('delete', cat.id)">
              <v-icon start size="18">mdi-delete</v-icon>
              删除
            </v-btn>
          </div>
        </div>
      </v-card>
    </div>
  </div>
</template>

<script>
import { ref, reactive } from 'vue'

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
    }

    return {
      formData,
      handleAdd
    }
  }
}
</script>

<style scoped>
.categories-panel {
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

.add-form-card {
  border-radius: 16px;
  margin-bottom: 24px;
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

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  background: white;
  border-radius: 16px;
}

.empty-text {
  margin-top: 16px;
  font-size: 1rem;
  color: #9ca3af;
}

.categories-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 16px;
}

.category-card {
  border-radius: 12px;
  transition: box-shadow 0.2s, transform 0.2s;
}

.category-card:hover {
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  transform: translateY(-2px);
}

.category-main {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
}

.category-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 700;
  font-size: 1.1rem;
  flex-shrink: 0;
}

.category-content {
  flex: 1;
  min-width: 0;
}

.category-name {
  font-weight: 600;
  color: #1a1a2e;
  font-size: 1rem;
}

.category-desc {
  font-size: 0.85rem;
  color: #6b7280;
  margin-top: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.category-actions {
  display: flex;
  gap: 4px;
  flex-shrink: 0;
}

@media (max-width: 600px) {
  .categories-grid {
    grid-template-columns: 1fr;
  }

  .category-main {
    flex-wrap: wrap;
  }

  .category-actions {
    width: 100%;
    justify-content: flex-end;
    margin-top: 8px;
  }
}
</style>
