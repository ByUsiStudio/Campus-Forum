<template>
  <div class="user-avatar">
    <div class="avatar-wrap" :style="{ width: size + 'px', height: size + 'px' }">
      <img
        v-if="avatarSrc"
        :src="avatarSrc"
        :alt="displayName"
        class="avatar-img"
        loading="lazy"
      />
      <div v-else class="avatar-fallback" :style="{ fontSize: Math.round(size * 0.45) + 'px', ...avatarStyle }">
        {{ initial }}
      </div>
    </div>
    <div v-if="showUsername" class="user-info">
      <div class="user-name" :style="nameStyle">{{ displayName }}</div>
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

const defaultTitle = { id: 0, name: '注册用户', color: 'var(--campus-primary)' }

const displayName = computed(() =>
  (props.user?.display_name || props.user?.username || '未知用户').trim()
)

const avatarSrc = computed(() => props.user?.avatar || '')

const displayTitles = computed(() => {
  const titles = props.user?.titles || []
  const visible = titles.slice(0, props.maxVisibleTitles)
  return visible.length ? visible : [defaultTitle]
})

const initial = computed(() => (props.user?.display_name || props.user?.username || 'U').trim().charAt(0).toUpperCase())

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
  const base = color || 'var(--campus-primary)'
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

.avatar-wrap {
  flex-shrink: 0;
  border-radius: 50%;
  overflow: hidden;
  line-height: 0;
  background: var(--campus-surface-2);
  border: 1px solid var(--campus-border);
}

.avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 50%;
}

.avatar-fallback {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  line-height: 1;
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
