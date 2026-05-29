<template>
  <div class="user-avatar d-flex align-center">
    <v-avatar :size="size" class="flex-shrink-0">
      <v-img :src="user.avatar" :alt="user.display_name" />
    </v-avatar>
    <div class="user-info">
      <div class="user-name" :class="nameSizeClass">{{ user.display_name }}</div>
      <div class="titles-row">
        <template v-for="title in displayTitles" :key="title.id || 'default'">
          <v-chip
            :size="chipSize"
            :style="{ backgroundColor: adjustColorOpacity(title.color, 0.12), color: title.color, border: `1px solid ${adjustColorOpacity(title.color, 0.25)}` }"
            class="title-chip"
          >
            <v-icon v-if="title.icon" start size="x-small">{{ title.icon }}</v-icon>
            {{ title.name }}
          </v-chip>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  user: {
    type: Object,
    required: true
  },
  size: {
    type: Number,
    default: 40
  },
  showUsername: {
    type: Boolean,
    default: false
  },
  maxVisibleTitles: {
    type: Number,
    default: 2
  }
})

const defaultTitle = {
  id: 0,
  name: '注册用户',
  color: '#625B71',
  icon: 'mdi-account',
  description: '系统默认头衔'
}

const displayTitles = computed(() => {
  const titles = props.user.titles || []
  const visibleTitles = titles.slice(0, props.maxVisibleTitles)
  if (visibleTitles.length === 0) {
    return [defaultTitle]
  }
  return visibleTitles
})

const nameSizeClass = computed(() => {
  if (props.size >= 64) return 'text-h6 font-weight-bold'
  if (props.size >= 48) return 'text-body-1 font-weight-bold'
  return 'text-body-2 font-weight-medium'
})

const chipSize = computed(() => {
  if (props.size >= 64) return 'small'
  if (props.size >= 48) return 'x-small'
  return 'x-small'
})

const adjustColorOpacity = (color, opacity) => {
  if (!color) return `rgba(103, 80, 164, ${opacity})`
  if (color.startsWith('#')) {
    const r = parseInt(color.slice(1, 3), 16)
    const g = parseInt(color.slice(3, 5), 16)
    const b = parseInt(color.slice(5, 7), 16)
    return `rgba(${r}, ${g}, ${b}, ${opacity})`
  }
  return color
}
</script>

<style scoped>
.user-avatar {
  min-width: 0;
}

.user-info {
  min-width: 0;
  flex-shrink: 1;
  margin-left: 12px;
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
  font-weight: 500;
  backdrop-filter: blur(2px);
  height: 20px;
  margin: 0;
}

.title-chip :deep(.v-icon) {
  opacity: 0.9;
}

.title-chip :deep(.v-chip__content) {
  gap: 2px;
}

.flex-shrink-0 {
  flex-shrink: 0;
}
</style>