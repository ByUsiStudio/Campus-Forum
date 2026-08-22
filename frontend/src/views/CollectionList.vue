<template>
  <div class="page-container">
    <el-card class="card-surface">
      <template #header>
        <div class="collection-card-title">
          <el-icon class="mr-2"><Folder /></el-icon>
          <span class="font-medium">我的收藏夹</span>
          <el-button type="primary" class="ml-auto" @click="openCreateDialog">
            <el-icon class="mr-1"><Plus /></el-icon>
            新建收藏夹
          </el-button>
        </div>
      </template>

      <div v-if="loading" class="loading"></div>
      <el-alert
        v-else-if="errorMsg"
        :title="errorMsg"
        type="error"
        :closable="false"
        show-icon
      ></el-alert>

      <div v-else class="collection-grid">
        <el-card
          v-for="collection in collections"
          :key="collection.id"
          shadow="hover"
          class="collection-item"
        >
          <div class="d-flex align-center">
            <el-avatar :size="40" class="collection-avatar">
              <el-icon><Folder /></el-icon>
            </el-avatar>
            <div class="ml-3 collection-info">
              <div class="collection-name">{{ collection.name }}</div>
              <div class="collection-desc text-secondary">{{ collection.description }}</div>
            </div>
          </div>

          <div class="collection-chips mt-2">
            <el-tag size="small" type="info" effect="plain">
              <el-icon class="mr-1"><Document /></el-icon>
              {{ collection.article_count }} 文章
            </el-tag>
            <el-tag size="small" :type="collection.is_public ? 'success' : 'info'" effect="plain">
              {{ collection.is_public ? '公开' : '私密' }}
            </el-tag>
          </div>

          <div class="collection-actions">
            <el-button type="primary" link size="small" @click="viewCollection(collection.id)">
              查看
            </el-button>
            <el-button type="warning" link size="small" @click="openEditDialog(collection)">
              编辑
            </el-button>
            <el-button type="danger" link size="small" @click="deleteCollection(collection.id)">
              删除
            </el-button>
          </div>
        </el-card>
      </div>
    </el-card>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { Folder, Plus, Document } from '@element-plus/icons-vue'
import { collectionApi } from '@/api'
import { success, error, confirm } from '@/utils/message'
import { jcOpenHtml, jcFieldsConfig, jcCloseAll } from '@/utils/jcu'

export default {
  name: 'CollectionList',
  setup() {
    const loading = ref(true)
    const errorMsg = ref('')
    const collections = ref([])

    const loadCollections = async () => {
      loading.value = true
      errorMsg.value = ''
      try {
        const res = await collectionApi.getCollections()
        if (res.data.success) {
          collections.value = res.data.data
        } else {
          errorMsg.value = res.data.message || '加载收藏夹失败'
        }
      } catch (err) {
        errorMsg.value = '加载收藏夹失败'
        console.error('加载收藏夹失败:', err)
      } finally {
        loading.value = false
      }
    }

    const openCreateDialog = () => {
      const cfg = jcFieldsConfig([
        { name: 'name', label: '收藏夹名称', placeholder: '请输入收藏夹名称', required: true },
        { name: 'description', label: '描述', type: 'textarea', placeholder: '请输入描述' },
        { name: 'is_public', label: '公开收藏夹', type: 'text', placeholder: '填写 1 公开 / 0 私密（默认公开）' }
      ])
      jcOpenHtml({
        title: '新建收藏夹',
        content: cfg.html,
        width: 480,
        size: 'sm',
        buttons: [
          { text: '取消', type: 'default', action: () => jcCloseAll() },
          {
            text: '创建',
            type: 'primary',
            action: () => {
              if (!cfg.validate(document)) return
              const v = cfg.collect(document)
              createCollection({ name: v.name, description: v.description || '', is_public: v.is_public !== '0' })
              jcCloseAll()
            }
          }
        ]
      })
    }

    const createCollection = async (payload) => {
      try {
        const res = await collectionApi.createCollection(payload)
        if (res.data.success) {
          collections.value.push(res.data.data)
          success('创建成功')
        }
      } catch (err) {
        error('创建收藏夹失败')
        console.error('创建收藏夹失败:', err)
      }
    }

    const openEditDialog = (collection) => {
      const cfg = jcFieldsConfig([
        { name: 'name', label: '收藏夹名称', placeholder: '请输入收藏夹名称', value: collection.name, required: true },
        { name: 'description', label: '描述', type: 'textarea', placeholder: '请输入描述', value: collection.description || '' },
        { name: 'is_public', label: '是否公开（填 1/0）', type: 'text', value: collection.is_public ? '1' : '0' }
      ])
      jcOpenHtml({
        title: '编辑收藏夹',
        content: cfg.html,
        width: 480,
        size: 'sm',
        buttons: [
          { text: '取消', type: 'default', action: () => jcCloseAll() },
          {
            text: '保存',
            type: 'primary',
            action: () => {
              if (!cfg.validate(document)) return
              const v = cfg.collect(document)
              updateCollection({ ...collection, name: v.name, description: v.description || '', is_public: v.is_public !== '0' })
              jcCloseAll()
            }
          }
        ]
      })
    }

    const updateCollection = async (payload) => {
      try {
        const res = await collectionApi.updateCollection(payload.id, payload)
        if (res.data.success) {
          const index = collections.value.findIndex(c => c.id === payload.id)
          if (index !== -1) collections.value[index] = res.data.data
          success('保存成功')
        }
      } catch (err) {
        error('更新收藏夹失败')
        console.error('更新收藏夹失败:', err)
      }
    }

    const deleteCollection = async (collectionId) => {
      try {
        await confirm('确定要删除该收藏夹吗？', { type: 'warning' })
      } catch {
        return
      }
      try {
        await collectionApi.deleteCollection(collectionId)
        collections.value = collections.value.filter(c => c.id !== collectionId)
        success('删除成功')
      } catch (err) {
        error('删除收藏夹失败')
        console.error('删除收藏夹失败:', err)
      }
    }

    const viewCollection = async (collectionId) => {
      let currentCollection = null
      let collectionArticles = []
      try {
        const res = await collectionApi.getCollection(collectionId)
        if (res.data.success) {
          currentCollection = res.data.data.collection
          collectionArticles = res.data.data.articles || []
        }
      } catch (err) {
        error('加载收藏夹详情失败')
        console.error('加载收藏夹详情失败:', err)
        return
      }

      if (!collectionArticles.length) {
        jcOpenHtml({
          title: currentCollection?.name || '收藏夹详情',
          content: '<div style="text-align:center;color:var(--jc-text-light,#888);padding:20px 0;">暂无收藏文章</div>',
          width: 640,
          size: 'md',
          buttons: [{ text: '关闭', type: 'primary', action: () => jcCloseAll() }]
        })
        return
      }

      const rows = collectionArticles.map((item) => {
        const id = item.article.id
        const title = (item.article.title || '').replace(/</g, '&lt;')
        const note = (item.note || '无备注').replace(/</g, '&lt;')
        return `<div style="display:flex;align-items:center;justify-content:space-between;gap:10px;padding:10px 0;border-bottom:1px solid var(--jc-border,#eef1f7);">
          <div style="min-width:0;flex:1;">
            <div style="font-weight:600;">${title}</div>
            <div style="font-size:12px;color:var(--jc-text-light,#888);">${note}</div>
          </div>
          <div style="flex-shrink:0;">
            <a href="/article/${id}" style="margin-right:10px;">查看</a>
            <span data-jc-remove="${id}" style="cursor:pointer;color:var(--jc-danger,#ef4444);">移除</span>
          </div>
        </div>`
      }).join('')
      const containerId = 'jc-collection-view'
      const content = `<div id="${containerId}">${rows}</div>`

      jcOpenHtml({
        title: currentCollection?.name || '收藏夹详情',
        content,
        width: 680,
        size: 'md',
        onMount: (root) => {
          root.addEventListener('click', (ev) => {
            const target = ev.target.closest('[data-jc-remove]')
            if (!target) return
            const articleId = Number(target.getAttribute('data-jc-remove'))
            confirm('确定要从此收藏夹移除该文章吗？').then((ok) => {
              if (!ok) return
              removeArticle(currentCollection.id, articleId)
                .then(() => {
                  success('移除成功')
                  target.closest('[style*="border-bottom"]').remove()
                })
                .catch(() => error('移除文章失败'))
            })
          })
        },
        buttons: [{ text: '关闭', type: 'primary', action: () => jcCloseAll() }]
      })
    }

    const removeArticle = async (collectionId, articleId) => {
      await collectionApi.removeArticleFromCollection(collectionId, articleId)
    }

    onMounted(() => {
      loadCollections()
    })

    return {
      loading,
      errorMsg,
      collections,
      loadCollections,
      openCreateDialog,
      openEditDialog,
      updateCollection,
      deleteCollection,
      viewCollection,
      removeArticle,
      Folder,
      Plus,
      Document
    }
  }
}
</script>

<style scoped>
.collection-card-title {
  display: flex;
  align-items: center;
  font-weight: 500;
}

.ml-auto {
  margin-left: auto;
}

.ml-3 {
  margin-left: 12px;
}

.mr-1 {
  margin-right: 4px;
}

.mr-2 {
  margin-right: 8px;
}

.mt-2 {
  margin-top: 8px;
}

.collection-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
}

.collection-item {
  display: flex;
  flex-direction: column;
}

.collection-avatar {
  background-color: var(--campus-primary);
  color: #fff;
}

.collection-info {
  flex: 1;
  min-width: 0;
}

.collection-name {
  font-weight: 600;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.collection-desc {
  font-size: 12px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.collection-chips {
  display: flex;
  gap: 8px;
}

.collection-actions {
  display: flex;
  justify-content: flex-end;
  gap: 4px;
  margin-top: 12px;
}

.article-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  gap: 8px;
}

.article-row > div:first-child {
  min-width: 0;
  flex: 1;
}

.article-title {
  font-weight: 500;
  word-break: break-word;
}

.article-note {
  font-size: 12px;
}

.article-actions {
  display: flex;
  gap: 8px;
  margin-left: 16px;
  flex-shrink: 0;
}

.loading {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--campus-primary);
}

.loading::after {
  content: '';
  width: 24px;
  height: 24px;
  border: 2px solid var(--campus-primary);
  border-top-color: transparent;
  border-radius: 50%;
  animation: loading-spin 0.8s linear infinite;
}

@keyframes loading-spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
