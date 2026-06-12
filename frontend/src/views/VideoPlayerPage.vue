<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()

const videoUrl = ref('')

onMounted(() => {
  videoUrl.value = route.query.url || ''
})
</script>

<template>
  <v-app>
    <v-app-bar app>
      <v-btn icon @click="router.go(-1)">
        <v-icon>mdi-arrow-left</v-icon>
      </v-btn>
      <v-toolbar-title>视频播放</v-toolbar-title>
    </v-app-bar>
    
    <v-container class="py-6">
      <v-card v-if="videoUrl">
        <video 
          :src="videoUrl" 
          controls 
          class="w-full"
          style="max-height: 60vh;"
        />
      </v-card>
      
      <v-card v-else class="text-center py-12">
        <v-icon size="64" color="grey">mdi-video-off</v-icon>
        <p class="mt-4 text-grey">视频地址无效</p>
      </v-card>
    </v-container>
  </v-app>
</template>
