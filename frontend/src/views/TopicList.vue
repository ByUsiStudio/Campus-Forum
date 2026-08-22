<template>
  <div class="page-container topic-list-page">
    <el-card class="mb-4 card-surface">
      <template #header>
        <div class="topic-header">
          <el-icon :size="22" class="topic-header-icon"><Collection /></el-icon>
          <span class="topic-header-title">话题广场</span>
          <el-button
            type="primary"
            plain
            class="topic-header-btn"
            @click="showHotTopics"
          >
            <el-icon class="mr-1"><TrendCharts /></el-icon>
            热门话题
          </el-button>
        </div>
      </template>

      <!-- 话题卡片网格 -->
      <div v-if="topics.length" class="topic-grid">
        <div
          v-for="topic in topics"
          :key="topic.id"
          class="topic-grid-item"
        >
          <el-card
            shadow="hover"
            class="topic-card cursor-pointer"
            :body-style="{ padding: '16px' }"
            @click="viewTopic(topic.id)"
          >
            <div class="topic-main">
              <el-avatar :size="40" class="topic-avatar">
                <el-icon :size="22"><PriceTag /></el-icon>
              </el-avatar>
              <div class="topic-info">
                <div class="topic-name">{{ topic.display_name || topic.name }}</div>
                <div class="topic-desc text-secondary">{{ topic.description }}</div>
              </div>
            </div>

            <div class="topic-stats mt-2">
              <el-tag size="small" effect="plain" class="topic-stat-tag">
                <el-icon class="mr-1"><Document /></el-icon>
                {{ topic.article_count }} 文章
              </el-tag>
              <el-tag size="small" effect="plain" class="topic-stat-tag">
                <el-icon class="mr-1"><User /></el-icon>
                {{ topic.follow_count }} 关注
              </el-tag>
            </div>

            <div v-if="topic.is_hot" class="mt-2">
              <el-tag size="small" type="warning" effect="light">
                <el-icon class="mr-1"><TrendCharts /></el-icon>
                热门
              </el-tag>
            </div>
          </el-card>
        </div>
      </div>

      <el-empty v-else description="暂无话题" />

      <!-- 分页 -->
      <div v-if="totalPages > 1" class="text-center mt-4">
        <el-pagination
          layout="prev, pager, next"
          v-model:current-page="page"
          :page-size="limit"
          :total="total"
          @current-change="loadTopics"
        />
      </div>
    </el-card>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  Collection,
  TrendCharts,
  Document,
  User,
  PriceTag
} from '@element-plus/icons-vue'
import { topicApi } from '@/api'
import { error } from '@/utils/message'
import { jcOpenHtml, jcCloseAll } from '@/utils/jcu'

export default {
  name: 'TopicList',
  setup() {
    const router = useRouter()
    const topics = ref([])
    const page = ref(1)
    const limit = ref(20)
    const total = ref(0)
    const isFollowing = ref(false)

    const totalPages = computed(() => {
      return Math.ceil(total.value / limit.value)
    })

    const loadTopics = async () => {
      try {
        const res = await topicApi.getTopics(page.value, limit.value)
        if (res.data.success) {
          topics.value = res.data.data.topics || []
          total.value = res.data.data.total || 0
        }
      } catch (error) {
        console.error('加载话题失败:', error)
      }
    }

    const escapeHtml = (s) => String(s ?? '').replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;').replace(/'/g, '&#39;')

    const viewTopic = async (topicId) => {
      try {
        const res = await topicApi.getTopic(topicId)
        if (res.data.success) {
          const topic = res.data.data.topic
          const articles = res.data.data.articles || []

          let following = false
          try {
            const followedRes = await topicApi.getFollowedTopics()
            const followedList = (followedRes.data.success && followedRes.data.data) || []
            following = followedList.some(f => f.topic_id === topicId)
          } catch (err) {
            following = false
          }
          isFollowing.value = following

          const name = topic.display_name || topic.name
          const desc = topic.description || ''

          const articlesHtml = articles.length
            ? articles.map((at) => {
                const art = at.article || {}
                const title = escapeHtml(art.title || '')
                const author = escapeHtml(art.user?.display_name || art.user?.username || '')
                const date = escapeHtml(art.created_at ? formatDate(art.created_at) : '')
                return `<div data-jc-article="${art.id ?? ''}" style="display:flex;align-items:center;justify-content:space-between;gap:10px;padding:10px 0;border-bottom:1px solid var(--jc-border,#eef1f7);cursor:pointer;">
                  <div style="min-width:0;flex:1;">
                    <div style="font-weight:600;color:var(--jc-text,#333);">${title}</div>
                    <div style="font-size:12px;color:var(--jc-text-light,#888);">${author} · ${date}</div>
                  </div>
                  <span style="flex-shrink:0;color:var(--jc-primary,#4f6ef7);font-size:13px;">查看</span>
                </div>`
              }).join('')
            : '<div style="text-align:center;color:var(--jc-text-light,#888);padding:20px 0;">暂无相关文章</div>'

          const content = `<div>
            <div style="display:flex;align-items:center;justify-content:space-between;gap:12px;margin-bottom:12px;">
              <div style="font-size:1.1rem;font-weight:700;color:var(--jc-text,#333);">${escapeHtml(name)}</div>
              <button data-jc-follow data-following="${following ? '1' : '0'}" style="flex-shrink:0;cursor:pointer;border:none;border-radius:var(--jc-radius,10px);padding:7px 16px;font-size:13px;font-weight:600;background:${following ? 'var(--jc-surface,#fff)' : 'var(--jc-primary,#4f6ef7)'};color:${following ? 'var(--jc-primary,#4f6ef7)' : '#fff'};border:1px solid var(--jc-primary,#4f6ef7);">${following ? '已关注' : '关注'}</button>
            </div>
            ${desc ? `<div style="color:var(--jc-text-light,#888);font-size:13px;margin-bottom:14px;line-height:1.6;">${escapeHtml(desc)}</div>` : ''}
            <div style="height:1px;background:var(--jc-border,#eef1f7);margin:4px 0 12px;"></div>
            <div style="font-weight:600;font-size:14px;margin-bottom:4px;color:var(--jc-text,#333);">相关文章</div>
            <div>${articlesHtml}</div>
          </div>`

          jcOpenHtml({
            title: name || '话题详情',
            content,
            width: 680,
            size: 'md',
            buttons: [{ text: '关闭', type: 'primary', action: () => jcCloseAll() }],
            onMount: (root) => {
              const followBtn = root.querySelector('[data-jc-follow]')
              if (followBtn) {
                followBtn.addEventListener('click', async () => {
                  const current = followBtn.getAttribute('data-following') === '1'
                  try {
                    if (current) {
                      await topicApi.unfollowTopic(topicId)
                      isFollowing.value = false
                      followBtn.setAttribute('data-following', '0')
                      followBtn.textContent = '关注'
                      followBtn.style.background = 'var(--jc-primary,#4f6ef7)'
                      followBtn.style.color = '#fff'
                      followBtn.style.border = '1px solid var(--jc-primary,#4f6ef7)'
                    } else {
                      await topicApi.followTopic(topicId)
                      isFollowing.value = true
                      followBtn.setAttribute('data-following', '1')
                      followBtn.textContent = '已关注'
                      followBtn.style.background = 'var(--jc-surface,#fff)'
                      followBtn.style.color = 'var(--jc-primary,#4f6ef7)'
                      followBtn.style.border = '1px solid var(--jc-primary,#4f6ef7)'
                    }
                  } catch (err) {
                    error('关注操作失败，请稍后再试')
                  }
                })
              }
              root.addEventListener('click', (ev) => {
                const target = ev.target.closest('[data-jc-article]')
                if (!target) return
                const articleId = target.getAttribute('data-jc-article')
                if (articleId) router.push('/article/' + articleId)
              })
            }
          })
        }
      } catch (error) {
        console.error('加载话题详情失败:', error)
      }
    }

    const showHotTopics = async () => {
      try {
        const res = await topicApi.getTopics(1, 20, true)
        if (res.data.success) {
          topics.value = res.data.data.topics || []
          total.value = res.data.data.total || 0
        }
      } catch (error) {
        console.error('加载热门话题失败:', error)
      }
    }

    const formatDate = (date) => {
      return new Date(date).toLocaleDateString('zh-CN')
    }

    onMounted(() => {
      loadTopics()
    })

    return {
      topics,
      page,
      total,
      totalPages,
      isFollowing,
      Collection,
      TrendCharts,
      Document,
      User,
      PriceTag,
      loadTopics,
      viewTopic,
      showHotTopics,
      formatDate
    }
  }
}
</script>

<style scoped>
.topic-list-page {
  padding-top: 16px;
}

.topic-card {
  border-radius: var(--campus-radius);
  transition: transform 0.2s;
}

.topic-card:hover {
  transform: translateY(-4px);
}

.topic-header {
  display: flex;
  align-items: center;
  font-size: 1.25rem;
  font-weight: 600;
  gap: 8px;
}

.topic-header-icon {
  color: var(--campus-primary);
}

.topic-header-btn {
  margin-left: auto;
}

.topic-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
}

.topic-main {
  display: flex;
  align-items: center;
  gap: 12px;
}

.topic-avatar {
  flex-shrink: 0;
  color: #fff;
  font-weight: 600;
}

.topic-info {
  flex: 1;
  min-width: 0;
}

.topic-name {
  font-weight: 600;
  font-size: 1.05rem;
  color: var(--campus-text);
}

.topic-desc {
  font-size: 0.85rem;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.topic-stats {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.topic-stat-tag {
  display: inline-flex;
  align-items: center;
}
</style>
