<template>
  <div class="user-avatar">
    <el-avatar :size="size" :src="user?.avatar || ''" class="avatar" :style="avatarStyle">
      {{ initial }}
    </el-avatar>
    <div v-if="showUsername" class="user-info">
      <div class="user-name" :style="nameStyle">{{ user?.display_name || user?.username || '未知用户' }}</div>
      <div v-if="displayTitles.length" class="titles-row">
        <span
          v-for="title in displayTitles"
          :key="title.id || 'default'"
          class="title-chip"
          :style="titleStyle(title.color)"
        >
          {{ title.name }}
        </span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  user: { type: Object, default: null },
  size: { type: Number, default: 40 },
  showUsername: { type: Boolean, default: true },
  maxVisibleTitles: { type: Number, default: 2 }
})

const defaultTitle = { id: 0, name: '注册用户', color: '#625b71' }

const displayTitles = computed(() => {
  const titles = props.user?.titles || []
  const visible = titles.slice(0, props.maxVisibleTitles)
  return visible.length ? visible : [defaultTitle]
})

const initial = computed(() => (props.user?.display_name || props.user?.username || 'U')[0].toUpperCase())

const nameStyle = computed(() => ({
  fontWeight: props.size >= 48 ? 700 : 500,
  fontSize: props.size >= 64 ? '18px' : props.size >= 48 ? '16px' : '14px'
}))

const avatarStyle = computed(() => ({
  backgroundColor: 'var(--campus-primary)',
  color: '#fff',
  fontWeight: 600
}))

const titleStyle = (color) => {
  const base = color || '#625b71'
  return {
    color: base,
    border: `1px solid ${base}`,
    backgroundColor: 'transparent'
  }
}
</script>

<style scoped>
.user-avatar {
  display: flex;
  align-items: center;
  min-width: 0;
}
.avatar {
  flex-shrink: 0;
}
.user-info {
  margin-left: 12px;
  min-width: 0;
  overflow: hidden;
}
.user-name {
  line-height: 1.3;
  word-break: break-word;
}
.titles-row {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  margin-top: 4px;
}
.title-chip {
  font-size: 12px;
  line-height: 1;
  padding: 2px 6px;
  border-radius: 999px;
  white-space: nowrap;
}
</style>
