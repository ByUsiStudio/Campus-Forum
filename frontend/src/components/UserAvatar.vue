<template>
  <div class="user-avatar d-flex align-center">
    <v-avatar :size="size" class="mr-2">
      <v-img :src="user.avatar" :alt="user.display_name"></v-img>
    </v-avatar>
    <div class="user-info">
      <div class="d-flex align-center flex-wrap">
        <span class="user-name text-body-1 font-weight-medium">{{ user.display_name }}</span>
        <v-chip
          v-for="title in user.titles"
          :key="title.id"
          size="x-small"
          :color="title.color"
          class="ml-1"
        >
          <v-icon v-if="title.icon" start size="x-small">{{ title.icon }}</v-icon>
          {{ title.name }}
        </v-chip>
      </div>
      <span v-if="showUsername" class="text-caption text-medium-emphasis">@{{ user.username }}</span>
    </div>
  </div>
</template>

<script setup>
defineProps({
  user: {
    type: Object,
    required: true
  },
  size: {
    type: Number,
    default: 32
  },
  showUsername: {
    type: Boolean,
    default: false
  }
})
</script>

<style scoped>
.user-name {
  line-height: 1.2;
}
</style>