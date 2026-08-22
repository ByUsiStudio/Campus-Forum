<template>
  <div class="user-level-card card-surface">
    <div class="level-head">
      <div class="level-head-left">
        <div class="level-icon-chip">
          <el-icon :size="18"><Star /></el-icon>
        </div>
        <span class="level-title-text">用户等级</span>
      </div>
      <el-tag size="small" effect="plain" class="level-count-tag">
        <el-icon class="mr-1" :size="13"><Trophy /></el-icon>
        {{ unlockedCount }} / {{ totalCount }}
      </el-tag>
    </div>

    <div class="level-body">
      <div class="level-badge">
        <div class="level-ring">
          <div class="level-num">{{ level.level }}</div>
          <div class="level-label">LV</div>
        </div>
        <div class="level-badge-info">
          <div class="level-premium-title">{{ level.title }}</div>
          <div class="level-caption">继续互动提升等级</div>
        </div>
      </div>

      <div class="experience-block">
        <div class="experience-top">
          <span class="experience-label">成长经验</span>
          <span class="experience-values">
            <span class="exp-now">{{ level.experience }}</span>
            <span class="exp-sep"> / </span>
            <span class="exp-next">{{ level.next_level }}</span>
          </span>
        </div>
        <div class="experience-progress">
          <el-progress
            :percentage="experiencePercent"
            :show-text="false"
            :stroke-width="10"
            class="exp-progress"
          />
        </div>
        <div class="experience-hint text-secondary">距离下一级还差 {{ remainingExp }} 经验</div>
      </div>
    </div>

    <div class="achievements-section">
      <div class="achievements-header">
        <el-icon :size="18" class="trophy-icon"><Trophy /></el-icon>
        <span class="achievements-title">成就徽章</span>
      </div>

      <div class="achievements-grid">
        <div
          v-for="achievement in achievements"
          :key="achievement.id"
          class="achievement-card"
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
  </div>
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

    const remainingExp = computed(() => {
      return Math.max(0, (level.value.next_level || 0) - (level.value.experience || 0))
    })

    const unlockedCount = computed(() => {
      return achievements.value.filter(a => a.unlocked_at).length
    })

    const totalCount = computed(() => {
      return allAchievements.value.length
    })

    const getRarityColor = (rarity) => {
      const colors = {
        common: '#909399',
        rare: '#409eff',
        epic: '#9c27b0',
        legendary: '#ff9800'
      }
      return colors[rarity] || '#909399'
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
        common: { color: '#909399', borderColor: '#909399', backgroundColor: 'rgba(144,147,153,0.08)' },
        rare: { color: '#409eff', borderColor: '#409eff', backgroundColor: 'rgba(64,158,255,0.08)' },
        epic: { color: '#9c27b0', borderColor: '#9c27b0', backgroundColor: 'rgba(156,39,176,0.08)' },
        legendary: { color: '#ff9800', borderColor: '#ff9800', backgroundColor: 'rgba(255,152,0,0.08)' }
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
      remainingExp,
      unlockedCount,
      totalCount,
      getRarityColor,
      getRarityText,
      getTagStyle,
      formatDate
    }
  }
}
</script>

<style scoped>
.user-level-card {
  max-width: 800px;
  padding: 20px;
}

.level-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 12px;
  margin-bottom: 18px;
}

.level-head-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.level-icon-chip {
  width: 34px;
  height: 34px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 9px;
  background: linear-gradient(135deg, var(--campus-primary-light), var(--campus-primary));
  color: #fff;
  box-shadow: 0 4px 10px rgba(79, 110, 247, 0.3);
}

.level-title-text {
  font-size: 16px;
  font-weight: 700;
  color: var(--campus-text);
  letter-spacing: -0.01em;
}

.level-count-tag {
  border-radius: 999px;
  background: var(--campus-primary-soft);
  border-color: transparent;
  color: var(--campus-primary);
  font-weight: 500;
}

.level-body {
  display: flex;
  align-items: center;
  gap: 24px;
  flex-wrap: wrap;
  padding: 20px;
  background: linear-gradient(135deg, rgba(79, 110, 247, 0.08), rgba(109, 139, 251, 0.05));
  border: 1px solid var(--campus-border);
  border-radius: var(--campus-radius);
}

.level-badge {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-shrink: 0;
}

.level-ring {
  width: 84px;
  height: 84px;
  border-radius: 50%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #fff;
  background: conic-gradient(from 180deg, var(--campus-primary-light), var(--campus-primary-dark), var(--campus-primary));
  box-shadow: 0 8px 20px rgba(79, 110, 247, 0.35);
  position: relative;
}

.level-ring::before {
  content: '';
  position: absolute;
  inset: 5px;
  border-radius: 50%;
  background: var(--campus-surface);
  opacity: 0.12;
}

.level-num {
  font-size: 30px;
  font-weight: 800;
  line-height: 1;
  letter-spacing: -0.02em;
}

.level-label {
  font-size: 12px;
  font-weight: 600;
  margin-top: 2px;
  letter-spacing: 0.1em;
  opacity: 0.9;
}

.level-badge-info {
  min-width: 0;
}

.level-premium-title {
  font-size: 20px;
  font-weight: 800;
  color: var(--campus-text);
  background: linear-gradient(90deg, var(--campus-primary), var(--campus-primary-dark));
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.level-caption {
  font-size: 13px;
  color: var(--campus-text-secondary);
  margin-top: 4px;
}

.experience-block {
  flex: 1;
  min-width: 220px;
}

.experience-top {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  font-size: 13px;
}

.experience-label {
  font-weight: 600;
  color: var(--campus-text);
}

.experience-values {
  font-variant-numeric: tabular-nums;
}

.exp-now {
  font-weight: 700;
  color: var(--campus-primary);
  font-size: 15px;
}

.exp-sep,
.exp-next {
  color: var(--campus-text-muted);
}

.experience-progress {
  margin-top: 8px;
}

.exp-progress :deep(.el-progress-bar__outer) {
  border-radius: 999px;
  background: rgba(79, 110, 247, 0.12);
}

.exp-progress :deep(.el-progress-bar__inner) {
  background: linear-gradient(90deg, var(--campus-primary-light), var(--campus-primary));
  border-radius: 999px;
}

.experience-hint {
  font-size: 12px;
  margin-top: 6px;
}

.achievements-section {
  margin-top: 20px;
  min-height: 60px;
}

.achievements-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 14px;
  color: var(--campus-text);
}

.trophy-icon {
  color: var(--campus-warning);
}

.achievements-title {
  font-size: 15px;
  font-weight: 700;
}

.achievements-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(190px, 1fr));
  gap: 12px;
}

.achievement-card {
  background: var(--campus-surface-2);
  border-radius: var(--campus-radius-sm);
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.achievement-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--campus-shadow-sm);
}

.achievement-card-body {
  padding: 14px;
}

.achievement-icon {
  line-height: 0;
}

.achievement-icon-img {
  width: 40px;
  height: 40px;
  object-fit: cover;
  border-radius: 8px;
}

.achievement-name {
  font-weight: 700;
  color: var(--campus-text);
}

.achievement-desc {
  font-size: 12px;
  color: var(--campus-text-secondary);
}

.achievement-unlock {
  font-size: 12px;
}

.mt-1 { margin-top: 4px; }
.mt-2 { margin-top: 8px; }
.mr-1 { margin-right: 4px; }
</style>
