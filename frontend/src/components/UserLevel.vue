<template>
  <el-card class="user-level-card" shadow="never">
    <template #header>
      <div class="level-header">
        <el-icon :size="20"><Star /></el-icon>
        <span>用户等级</span>
      </div>
    </template>

    <div class="level-info">
      <div class="level-badge">
        <el-avatar :size="80" class="level-avatar">
          <span class="level-avatar-text">{{ level.level }}</span>
        </el-avatar>
        <div class="level-title mt-2">{{ level.title }}</div>
      </div>

      <div class="experience-bar mt-4">
        <div class="experience-progress">
          <el-progress
            :percentage="experiencePercent"
            :show-text="false"
            :stroke-width="10"
          />
        </div>
        <div class="experience-text mt-2">
          <span>{{ level.experience }} / {{ level.next_level }}</span>
          <span class="ml-2">经验值</span>
        </div>
      </div>
    </div>

    <el-divider />

    <div class="achievements-section">
      <div class="achievements-header mb-3">
        <el-icon :size="18"><Trophy /></el-icon>
        <span class="achievements-title">成就徽章</span>
        <el-tag size="small" type="primary" class="ml-2">{{ unlockedCount }} / {{ totalCount }}</el-tag>
      </div>

      <div class="achievements-grid">
        <div
          v-for="achievement in achievements"
          :key="achievement.id"
          class="achievement-card card-surface"
        >
          <div class="achievement-card-body text-center">
            <div class="achievement-icon" :style="{ color: getRarityColor(achievement.achievement.rarity) }">
              <el-icon :size="40" v-if="!achievement.achievement.icon || !achievement.achievement.icon.startsWith('http')">
                <Medal />
              </el-icon>
              <img
                v-else
                :src="achievement.achievement.icon"
                class="achievement-icon-img"
                alt=""
              />
            </div>
            <div class="achievement-name mt-2">{{ achievement.achievement.name }}</div>
            <div class="achievement-desc mt-1">{{ achievement.achievement.description }}</div>
            <el-tag
              size="small"
              effect="plain"
              class="mt-2"
              :style="getTagStyle(achievement.achievement.rarity)"
            >
              {{ getRarityText(achievement.achievement.rarity) }}
            </el-tag>
            <div v-if="achievement.unlocked_at" class="mt-2 text-secondary achievement-unlock">
              解锁于: {{ formatDate(achievement.unlocked_at) }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </el-card>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { Star, Trophy, Medal } from '@element-plus/icons-vue'
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

    const getTagStyle = (rarity) => {
      const styles = {
        common: { color: '#909399', borderColor: '#909399', backgroundColor: 'rgba(144,147,153,0.1)' },
        rare: { color: '#409eff', borderColor: '#409eff', backgroundColor: 'rgba(64,158,255,0.1)' },
        epic: { color: '#9c27b0', borderColor: '#9c27b0', backgroundColor: 'rgba(156,39,176,0.1)' },
        legendary: { color: '#ff9800', borderColor: '#ff9800', backgroundColor: 'rgba(255,152,0,0.1)' }
      }
      return styles[rarity] || styles.common
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
      getTagStyle,
      formatDate,
      Star,
      Trophy,
      Medal
    }
  }
}
</script>

<style scoped>
.user-level-card {
  max-width: 800px;
}

.level-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
  color: var(--campus-text);
}

.level-badge {
  text-align: center;
}

.level-avatar {
  background-color: var(--campus-primary);
  color: #fff;
  font-weight: 700;
}

.level-avatar-text {
  font-size: 28px;
}

.level-title {
  font-size: 18px;
  font-weight: bold;
  color: var(--campus-text);
}

.experience-bar {
  width: 100%;
}

.experience-progress {
  width: 100%;
  border-radius: 6px;
  overflow: hidden;
}

.experience-progress :deep(.el-progress-bar__outer) {
  border-radius: 5px;
}

.experience-progress :deep(.el-progress-bar__inner) {
  background-color: var(--campus-primary);
}

.experience-text {
  display: flex;
  justify-content: space-between;
  font-size: 14px;
}

.achievements-section {
  min-height: 60px;
}

.achievements-header {
  display: flex;
  align-items: center;
  color: var(--campus-text);
}

.achievements-title {
  font-size: 15px;
  font-weight: 600;
}

.achievements-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 12px;
}

.achievement-card {
  transition: transform 0.2s;
}

.achievement-card:hover {
  transform: translateY(-2px);
}

.achievement-card-body {
  padding: 14px;
  background: transparent;
}

.achievement-icon {
  line-height: 0;
}

.achievement-icon-img {
  width: 40px;
  height: 40px;
  object-fit: cover;
}

.achievement-name {
  font-weight: bold;
  color: var(--campus-text);
}

.achievement-desc {
  font-size: 12px;
  color: #666;
}

.achievement-unlock {
  font-size: 12px;
}

.ml-2 {
  margin-left: 8px;
}
</style>
