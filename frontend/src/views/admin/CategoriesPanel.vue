<template>
  <div>
    <!-- 添加新分区表单 -->
    <v-card class="mb-4" variant="flat" rounded="lg">
      <v-card-title class="pb-2">
        <v-icon start size="20">mdi-plus-circle</v-icon>
        添加新分区
      </v-card-title>
      <v-card-text>
        <v-row dense>
          <v-col cols="12" sm="4">
            <v-text-field
              v-model="formData.name"
              label="分区名称"
              placeholder="例如：技术交流"
              variant="outlined"
              density="compact"
              hide-details
            />
          </v-col>
          <v-col cols="12" sm="5">
            <v-text-field
              v-model="formData.description"
              label="描述"
              placeholder="分区简介"
              variant="outlined"
              density="compact"
              hide-details
            />
          </v-col>
          <v-col cols="12" sm="3">
            <v-btn color="primary" block height="40" @click="handleAdd">
              <v-icon start>mdi-plus</v-icon>
              添加
            </v-btn>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>

    <!-- 分区列表 -->
    <v-card variant="flat" rounded="lg">
      <v-list lines="two" v-if="categories.length > 0">
        <v-list-item v-for="cat in categories" :key="cat.id" class="py-3">
          <template v-slot:prepend>
            <v-avatar size="48" :color="cat.color || 'primary'" variant="tonal">
              <span class="text-h6">{{ cat.sort_order || 0 }}</span>
            </v-avatar>
          </template>

          <v-list-item-title class="font-weight-medium mb-1">
            {{ cat.name }}
          </v-list-item-title>

          <v-list-item-subtitle>
            {{ cat.description || '暂无描述' }}
          </v-list-item-subtitle>

          <template v-slot:append>
            <v-btn-group variant="text" density="compact" divided>
              <v-btn size="small" color="primary" @click="$emit('edit', cat)">
                <v-icon>mdi-pencil</v-icon>
                <v-tooltip activator="parent">编辑</v-tooltip>
              </v-btn>
              <v-btn size="small" color="error" @click="$emit('delete', cat)">
                <v-icon>mdi-delete</v-icon>
                <v-tooltip activator="parent">删除</v-tooltip>
              </v-btn>
            </v-btn-group>
          </template>
        </v-list-item>
      </v-list>

      <v-card-text v-else class="text-center py-8">
        <v-icon size="48" color="grey-lighten-1">mdi-folder-outline</v-icon>
        <div class="text-body-1 text-medium-emphasis mt-2">
          暂无分区数据
        </div>
      </v-card-text>
    </v-card>
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