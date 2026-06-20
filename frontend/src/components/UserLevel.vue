<template>
  <v-card class="user-level-card">
    <v-card-title class="d-flex align-center">
      <v-icon left>mdi-star</v-icon>
      用户等级
    </v-card-title>
    <v-card-text>
      <div class="level-info">
        <div class="level-badge">
          <v-avatar size="80" color="primary">
            <span class="white--text headline">{{ level.level }}</span>
          </v-avatar>
          <div class="level-title mt-2">{{ level.title }}</div>
        </div>
        
        <div class="experience-bar mt-4">
          <v-progress-linear
            :value="experiencePercent"
            color="primary"
            height="10"
            rounded
          ></v-progress-linear>
          <div class="experience-text mt-2">
            <span>{{ level.experience }} / {{ level.next_level }}</span>
            <span class="ml-2">经验值</span>
          </div>
        </div>
      </div>

      <v-divider class="my-4"></v-divider>

      <div class="achievements-section">
        <div class="d-flex align-center mb-3">
          <v-icon left>mdi-trophy</v-icon>
          <span class="subtitle-1">成就徽章</span>
          <v-chip small color="primary" class="ml-2">{{ unlockedCount }} / {{ totalCount }}</v-chip>
        </div>

        <v-row dense>
          <v-col cols="12" sm="6" md="4" v-for="achievement in achievements" :key="achievement.id">
            <v-card outlined class="achievement-card">
              <v-card-text class="text-center">
                <v-icon :color="getRarityColor(achievement.achievement.rarity)" size="40">
                  {{ achievement.achievement.icon || 'mdi-medal' }}
                </v-icon>
                <div class="achievement-name mt-2">{{ achievement.achievement.name }}</div>
                <div class="achievement-desc mt-1">{{ achievement.achievement.description }}</div>
                <v-chip small :color="getRarityColor(achievement.achievement.rarity)" class="mt-2">
                  {{ getRarityText(achievement.achievement.rarity) }}
                </v-chip>
                <div v-if="achievement.unlocked_at" class="mt-2 text-caption">
                  解锁于: {{ formatDate(achievement.unlocked_at) }}
                </div>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>
      </div>
    </v-card-text>
  </v-card>
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

    const getRarityColor = (rarity) => {
      const colors = {
        common: 'grey',
        rare: 'blue',
        epic: 'purple',
        legendary: 'orange'
      }
      return colors[rarity] || 'grey'
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
      getRarityText,
      formatDate
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

.achievement-card {
  transition: transform 0.2s;
}

.achievement-card:hover {
  transform: translateY(-2px);
}

.achievement-name {
  font-weight: bold;
}

.achievement-desc {
  font-size: 12px;
  color: #666;
}
</style>