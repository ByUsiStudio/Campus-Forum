<template>
  <div class="user-level-card layui-card">
    <div class="layui-card-header d-flex align-center gap-3">
      <i class="fa-solid fa-star"></i>
      <span style="font-size: 18px; font-weight: 600;">用户等级</span>
    </div>
    
    <div class="layui-card-body">
      <div class="level-info">
        <div class="level-badge text-center">
          <div class="level-avatar" style="width: 80px; height: 80px; background: var(--primary); border-radius: 50%; display: flex; align-items: center; justify-content: center; margin: 0 auto;">
            <span style="color: white; font-size: 36px; font-weight: 700;">{{ level.level }}</span>
          </div>
          <div class="level-title mt-2">{{ level.title }}</div>
        </div>
        
        <div class="experience-bar mt-4">
          <div class="progress-bar">
            <div class="progress-fill" :style="{ width: experiencePercent + '%' }"></div>
          </div>
          <div class="experience-text mt-2">
            <span>{{ level.experience }} / {{ level.next_level }}</span>
            <span class="ml-2">经验值</span>
          </div>
        </div>
      </div>

      <hr class="layui-bg-gray mt-4 mb-4" />

      <div class="achievements-section">
        <div class="d-flex align-center mb-3">
          <i class="fa-solid fa-trophy mr-2"></i>
          <span style="font-size: 16px; font-weight: 600;">成就徽章</span>
          <span class="layui-badge layui-bg-primary ml-2">{{ unlockedCount }} / {{ totalCount }}</span>
        </div>

        <div class="achievements-grid">
          <div
            v-for="achievement in achievements"
            :key="achievement.id"
            class="achievement-card"
          >
            <div class="achievement-content text-center">
              <i :class="getIconClass(achievement.achievement.icon)" :style="{ color: getRarityColor(achievement.achievement.rarity) }" style="font-size: 40px;"></i>
              <div class="achievement-name mt-2">{{ achievement.achievement.name }}</div>
              <div class="achievement-desc mt-1">{{ achievement.achievement.description }}</div>
              <span class="achievement-rarity mt-2" :style="{ background: getRarityBgColor(achievement.achievement.rarity), color: getRarityColor(achievement.achievement.rarity) }">{{ getRarityText(achievement.achievement.rarity) }}</span>
              <div v-if="achievement.unlocked_at" class="mt-2 text-muted" style="font-size: 12px;">
                解锁于: {{ formatDate(achievement.unlocked_at) }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { levelApi } from '@/api'

export default {
  name: 'UserLevel',
  setup() {
    const level = ref({
      level: 1,
      experience: 0,
      next_level: 100,
      title: '新手'
    })
    const achievements = ref([])
    const allAchievements = ref([])

    const experiencePercent = computed(() => {
      return (level.value.experience / level.value.next_level) * 100
    })

    const unlockedCount = computed(() => {
      return achievements.value.filter(a => a.unlocked_at).length
    })

    const totalCount = computed(() => {
      return allAchievements.value.length
    })

    const mdiToFa = {
      'mdi-star': 'fa-solid fa-star',
      'mdi-trophy': 'fa-solid fa-trophy',
      'mdi-medal': 'fa-solid fa-medal',
      'mdi-award': 'fa-solid fa-award',
      'mdi-crown': 'fa-solid fa-crown',
      'mdi-badge': 'fa-solid fa-award',
      'mdi-flame': 'fa-solid fa-flame',
      'mdi-rocket': 'fa-solid fa-rocket',
      'mdi-code': 'fa-solid fa-code',
      'mdi-book': 'fa-solid fa-book',
      'mdi-book-open': 'fa-solid fa-book-open',
      'mdi-pencil': 'fa-solid fa-pencil',
      'mdi-message': 'fa-solid fa-message',
      'mdi-comment': 'fa-solid fa-comment',
      'mdi-eye': 'fa-solid fa-eye',
      'mdi-globe': 'fa-solid fa-globe',
      'mdi-flag': 'fa-solid fa-flag',
      'mdi-heart': 'fa-solid fa-heart',
      'mdi-heart-outline': 'fa-regular fa-heart',
      'mdi-lightbulb': 'fa-solid fa-lightbulb',
      'mdi-lightbulb-outline': 'fa-regular fa-lightbulb',
      'mdi-smile': 'fa-solid fa-face-smile',
      'mdi-frown': 'fa-solid fa-face-frown',
      'mdi-meh': 'fa-solid fa-face-meh',
      'mdi-hand': 'fa-solid fa-hand',
      'mdi-hand-heart': 'fa-solid fa-heart-handshake',
      'mdi-handshake': 'fa-solid fa-handshake',
      'mdi-thumbs-up': 'fa-solid fa-thumbs-up',
      'mdi-thumbs-down': 'fa-solid fa-thumbs-down',
      'mdi-share': 'fa-solid fa-share-nodes',
      'mdi-link': 'fa-solid fa-link',
      'mdi-home': 'fa-solid fa-house',
      'mdi-settings': 'fa-solid fa-gear',
      'mdi-cog': 'fa-solid fa-gear',
      'mdi-lock': 'fa-solid fa-lock',
      'mdi-unlock': 'fa-solid fa-unlock',
      'mdi-bell': 'fa-solid fa-bell',
      'mdi-bell-ring': 'fa-solid fa-bell',
      'mdi-calendar': 'fa-solid fa-calendar',
      'mdi-clock': 'fa-solid fa-clock',
      'mdi-image': 'fa-solid fa-image',
      'mdi-video': 'fa-solid fa-video',
      'mdi-file': 'fa-solid fa-file',
      'mdi-file-text': 'fa-solid fa-file-lines',
      'mdi-folder': 'fa-solid fa-folder',
      'mdi-tag': 'fa-solid fa-tag',
      'mdi-users': 'fa-solid fa-users',
      'mdi-user': 'fa-solid fa-user',
      'mdi-user-circle': 'fa-solid fa-user-circle',
      'mdi-account': 'fa-solid fa-user',
      'mdi-account-circle': 'fa-solid fa-user-circle',
      'mdi-account-group': 'fa-solid fa-users',
      'mdi-group': 'fa-solid fa-users',
      'mdi-upload': 'fa-solid fa-upload',
      'mdi-download': 'fa-solid fa-download',
      'mdi-save': 'fa-solid fa-floppy-disk',
      'mdi-delete': 'fa-solid fa-trash',
      'mdi-edit': 'fa-solid fa-pencil',
      'mdi-plus': 'fa-solid fa-plus',
      'mdi-minus': 'fa-solid fa-minus',
      'mdi-check': 'fa-solid fa-check',
      'mdi-close': 'fa-solid fa-xmark',
      'mdi-x': 'fa-solid fa-xmark',
      'mdi-alert': 'fa-solid fa-circle-exclamation',
      'mdi-alert-circle': 'fa-solid fa-circle-exclamation',
      'mdi-info': 'fa-solid fa-circle-info',
      'mdi-help': 'fa-solid fa-circle-question',
      'mdi-help-circle': 'fa-solid fa-circle-question',
      'mdi-search': 'fa-solid fa-search',
      'mdi-magnify': 'fa-solid fa-magnifying-glass',
      'mdi-filter': 'fa-solid fa-filter',
      'mdi-sort': 'fa-solid fa-arrow-down-wide-short',
      'mdi-refresh': 'fa-solid fa-rotate-right',
      'mdi-reload': 'fa-solid fa-rotate-right',
      'mdi-play': 'fa-solid fa-play',
      'mdi-pause': 'fa-solid fa-pause',
      'mdi-volume-high': 'fa-solid fa-volume-high',
      'mdi-volume-mute': 'fa-solid fa-volume-xmark',
      'mdi-fullscreen': 'fa-solid fa-expand',
      'mdi-list': 'fa-solid fa-list',
      'mdi-grid': 'fa-solid fa-grid-3x3',
      'mdi-chart-bar': 'fa-solid fa-chart-bar',
      'mdi-chart-line': 'fa-solid fa-chart-line',
      'mdi-chart-pie': 'fa-solid fa-chart-pie',
      'mdi-wallet': 'fa-solid fa-wallet',
      'mdi-coins': 'fa-solid fa-coins',
      'mdi-gift': 'fa-solid fa-gift',
      'mdi-shield': 'fa-solid fa-shield',
      'mdi-shield-check': 'fa-solid fa-shield-check',
      'mdi-key': 'fa-solid fa-key',
      'mdi-message-circle': 'fa-solid fa-comment-circle',
      'mdi-bookmark': 'fa-solid fa-bookmark',
      'mdi-bookmark-outline': 'fa-regular fa-bookmark',
      'mdi-rss': 'fa-solid fa-rss',
      'mdi-cloud': 'fa-solid fa-cloud',
      'mdi-sun': 'fa-solid fa-sun',
      'mdi-moon': 'fa-solid fa-moon',
      'mdi-snowflake': 'fa-solid fa-snowflake',
      'mdi-droplet': 'fa-solid fa-droplet',
      'mdi-fire': 'fa-solid fa-flame',
      'mdi-leaf': 'fa-solid fa-leaf',
      'mdi-tree': 'fa-solid fa-tree',
      'mdi-mountain': 'fa-solid fa-mountain',
      'mdi-earth': 'fa-solid fa-globe',
      'mdi-map': 'fa-solid fa-map',
      'mdi-map-marker': 'fa-solid fa-map-pin',
      'mdi-compass': 'fa-solid fa-compass',
      'mdi-circle': 'fa-solid fa-circle',
      'mdi-circle-outline': 'fa-regular fa-circle',
      'mdi-square': 'fa-solid fa-square',
      'mdi-square-outline': 'fa-regular fa-square',
      'mdi-triangle': 'fa-solid fa-triangle',
      'mdi-diamond': 'fa-solid fa-diamond',
      'mdi-octagon': 'fa-solid fa-octagon',
      'mdi-hexagon': 'fa-solid fa-hexagon',
      'mdi-hexagon-outline': 'fa-regular fa-hexagon',
      'mdi-play-circle': 'fa-solid fa-circle-play',
      'mdi-pause-circle': 'fa-solid fa-circle-pause',
      'mdi-stop-circle': 'fa-solid fa-circle-stop',
      'mdi-skip-next': 'fa-solid fa-forward-step',
      'mdi-skip-previous': 'fa-solid fa-backward-step',
      'mdi-fast-forward': 'fa-solid fa-fast-forward',
      'mdi-fast-rewind': 'fa-solid fa-fast-backward',
      'mdi-shuffle': 'fa-solid fa-shuffle',
      'mdi-repeat': 'fa-solid fa-repeat',
      'mdi-repeat-once': 'fa-solid fa-repeat-1',
      'mdi-repeat-variant': 'fa-solid fa-repeat',
      'mdi-volume-1': 'fa-solid fa-volume-low',
      'mdi-volume-2': 'fa-solid fa-volume-high',
      'mdi-volume-medium': 'fa-solid fa-volume-medium',
      'mdi-volume-minus': 'fa-solid fa-volume-down',
      'mdi-volume-plus': 'fa-solid fa-volume-up',
      'mdi-volume-mute': 'fa-solid fa-volume-xmark',
      'mdi-subtitles': 'fa-solid fa-closed-captioning',
      'mdi-subtitles-outline': 'fa-solid fa-closed-captioning',
      'mdi-captions': 'fa-solid fa-closed-captioning',
      'mdi-captions-off': 'fa-solid fa-closed-captioning',
      'mdi-aspect-ratio': 'fa-solid fa-square',
      'mdi-aspect-ratio-box': 'fa-solid fa-square',
      'mdi-crop': 'fa-solid fa-crop',
      'mdi-crop-16-9': 'fa-solid fa-square',
      'mdi-crop-4-3': 'fa-solid fa-square',
      'mdi-crop-free': 'fa-solid fa-square',
      'mdi-rotate': 'fa-solid fa-rotate',
      'mdi-rotate-3d': 'fa-solid fa-rotate',
      'mdi-rotate-ccw': 'fa-solid fa-rotate-left',
      'mdi-rotate-cw': 'fa-solid fa-rotate-right',
      'mdi-flip': 'fa-solid fa-arrows-up-down-left-right',
      'mdi-flip-horizontal': 'fa-solid fa-rotate-horizontal',
      'mdi-flip-vertical': 'fa-solid fa-rotate-vertical',
      'mdi-mirror': 'fa-solid fa-rotate-horizontal',
      'mdi-step-forward': 'fa-solid fa-step-forward',
      'mdi-step-backward': 'fa-solid fa-step-backward',
      'mdi-eject': 'fa-solid fa-eject',
      'mdi-music': 'fa-solid fa-music',
      'mdi-music-box': 'fa-solid fa-music',
      'mdi-music-circle': 'fa-solid fa-music',
      'mdi-disc': 'fa-solid fa-disc',
      'mdi-cd': 'fa-solid fa-compact-disc',
      'mdi-bluetooth': 'fa-brands fa-bluetooth',
      'mdi-bluetooth-audio': 'fa-brands fa-bluetooth-b',
      'mdi-headphones-bluetooth': 'fa-solid fa-headphones',
      'mdi-speaker-bluetooth': 'fa-solid fa-volume-high',
      'mdi-cpu': 'fa-solid fa-cpu',
      'mdi-harddrive': 'fa-solid fa-hard-drive',
      'mdi-memory': 'fa-solid fa-memory-stick',
      'mdi-network': 'fa-solid fa-network-wired',
      'mdi-bluetooth-off': 'fa-solid fa-bluetooth-b',
      'mdi-usb': 'fa-solid fa-usb',
      'mdi-printer': 'fa-solid fa-printer',
      'mdi-scanner': 'fa-solid fa-scanner',
      'mdi-monitor': 'fa-solid fa-monitor',
      'mdi-laptop': 'fa-solid fa-laptop',
      'mdi-tablet': 'fa-solid fa-tablet-screen-button',
      'mdi-smartphone': 'fa-solid fa-mobile-screen-button',
      'mdi-watch': 'fa-solid fa-watch',
      'mdi-calculator': 'fa-solid fa-calculator',
      'mdi-keyboard': 'fa-solid fa-keyboard',
      'mdi-gamepad': 'fa-solid fa-gamepad',
      'mdi-joystick': 'fa-solid fa-gamepad-2',
      'mdi-controller': 'fa-solid fa-gamepad',
      'mdi-beaker': 'fa-solid fa-flask-conical',
      'mdi-flask': 'fa-solid fa-flask-vial',
      'mdi-dna': 'fa-solid fa-dna',
      'mdi-molecule': 'fa-solid fa-atom',
      'mdi-atom': 'fa-solid fa-atom',
      'mdi-robot': 'fa-solid fa-robot'
    }

    const getIconClass = (icon) => {
      if (!icon) return 'fa-solid fa-medal'
      if (icon.startsWith('fa-')) return icon
      return mdiToFa[icon] || 'fa-solid fa-medal'
    }

    const getRarityColor = (rarity) => {
      const colors = {
        common: '#666',
        rare: '#1E9FFF',
        epic: '#9B59B6',
        legendary: '#FFA500'
      }
      return colors[rarity] || '#666'
    }

    const getRarityBgColor = (rarity) => {
      const colors = {
        common: 'rgba(102, 102, 102, 0.1)',
        rare: 'rgba(30, 159, 255, 0.1)',
        epic: 'rgba(155, 89, 182, 0.1)',
        legendary: 'rgba(255, 165, 0, 0.1)'
      }
      return colors[rarity] || 'rgba(102, 102, 102, 0.1)'
    }

    const getRarityText = (rarity) => {
      const texts = {
        common: '普通',
        rare: '稀有',
        epic: '史诗',
        legendary: '传说'
      }
      return texts[rarity] || '普通'
    }

    const formatDate = (date) => {
      return new Date(date).toLocaleDateString('zh-CN')
    }

    const loadData = async () => {
      try {
        const levelRes = await levelApi.getLevel()
        if (levelRes.data.success) {
          level.value = levelRes.data.data
        }

        const achievementsRes = await levelApi.getUserAchievements()
        if (achievementsRes.data.success) {
          achievements.value = achievementsRes.data.data
        }

        const allRes = await levelApi.getAllAchievements()
        if (allRes.data.success) {
          allAchievements.value = allRes.data.data
        }
      } catch (error) {
        console.error('加载等级数据失败:', error)
      }
    }

    onMounted(() => {
      loadData()
    })

    return {
      level,
      achievements,
      experiencePercent,
      unlockedCount,
      totalCount,
      getRarityColor,
      getRarityBgColor,
      getRarityText,
      formatDate,
      getIconClass
    }
  }
}
</script>

<style scoped>
.user-level-card {
  max-width: 800px;
}

.level-badge {
  text-align: center;
}

.level-title {
  font-size: 18px;
  font-weight: bold;
}

.experience-text {
  display: flex;
  justify-content: space-between;
  font-size: 14px;
}

.progress-bar {
  height: 10px;
  background: #f0f0f0;
  border-radius: 5px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: var(--primary);
  border-radius: 5px;
  transition: width 0.3s ease;
}

.achievements-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 15px;
}

.achievement-card {
  background: white;
  border: 1px solid #f0f0f0;
  border-radius: 8px;
  padding: 15px;
  transition: transform 0.2s;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }
}

.achievement-content {
  text-align: center;
}

.achievement-name {
  font-weight: bold;
  font-size: 14px;
}

.achievement-desc {
  font-size: 12px;
  color: #666;
}

.achievement-rarity {
  display: inline-block;
  padding: 2px 10px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.text-muted {
  color: #999;
}

.d-flex {
  display: flex;
}

.align-center {
  align-items: center;
}

.mr-2 {
  margin-right: 8px;
}

.ml-2 {
  margin-left: 8px;
}

.mt-2 {
  margin-top: 8px;
}

.mt-4 {
  margin-top: 16px;
}

.mb-3 {
  margin-bottom: 12px;
}

.mb-4 {
  margin-bottom: 16px;
}
</style>