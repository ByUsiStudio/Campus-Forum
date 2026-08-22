<template>
  <div class="page-container">
    <el-card class="card-surface">
      <template #header>
        <div class="collection-card-title">
          <el-icon class="mr-2"><Folder /></el-icon>
          <span class="font-medium">我的收藏夹</span>
          <el-button type="primary" class="ml-auto" @click="createDialog = true">
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
            <el-button type="warning" link size="small" @click="editCollection(collection)">
              编辑
            </el-button>
            <el-button type="danger" link size="small" @click="deleteCollection(collection.id)">
              删除
            </el-button>
          </div>
        </el-card>
      </div>
    </el-card>

    <!-- 创建收藏夹对话框 -->
    <el-dialog v-model="createDialog" title="新建收藏夹" width="500px">
      <el-form ref="createForm" :model="newCollection">
        <el-form-item
          label="收藏夹名称"
          prop="name"
          :rules="[{ required: true, message: '请输入名称', trigger: 'blur' }]"
        >
          <el-input v-model="newCollection.name" placeholder="请输入收藏夹名称"></el-input>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="newCollection.description"
            type="textarea"
            :rows="3"
            placeholder="请输入描述"
          ></el-input>
        </el-form-item>
        <el-form-item label="公开收藏夹" prop="is_public">
          <el-switch v-model="newCollection.is_public"></el-switch>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialog = false">取消</el-button>
        <el-button type="primary" @click="createCollection">创建</el-button>
      </template>
    </el-dialog>

    <!-- 编辑收藏夹对话框 -->
    <el-dialog v-model="editDialog" title="编辑收藏夹" width="500px">
      <el-form ref="editForm" :model="editCollectionData">
        <el-form-item
          label="收藏夹名称"
          prop="name"
          :rules="[{ required: true, message: '请输入名称', trigger: 'blur' }]"
        >
          <el-input v-model="editCollectionData.name" placeholder="请输入收藏夹名称"></el-input>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="editCollectionData.description"
            type="textarea"
            :rows="3"
            placeholder="请输入描述"
          ></el-input>
        </el-form-item>
        <el-form-item label="公开收藏夹" prop="is_public">
          <el-switch v-model="editCollectionData.is_public"></el-switch>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialog = false">取消</el-button>
        <el-button type="primary" @click="updateCollection">保存</el-button>
      </template>
    </el-dialog>

    <!-- 收藏夹详情对话框 -->
    <el-dialog v-model="viewDialog" :title="currentCollection.name || '收藏夹详情'" width="800px">
      <el-empty v-if="!collectionArticles.length" description="暂无收藏文章"></el-empty>
      <el-list v-else>
        <el-list-item v-for="item in collectionArticles" :key="item.id">
          <div class="article-row">
            <div>
              <div class="article-title">{{ item.article.title }}</div>
              <div class="text-secondary article-note">{{ item.note || '无备注' }}</div>
            </div>
            <div class="article-actions">
              <el-button type="primary" link size="small" @click="$router.push(`/article/${item.article.id}`)">
                查看
              </el-button>
              <el-button type="danger" link size="small" @click="removeArticle(item.article.id)">
                移除
              </el-button>
            </div>
          </div>
        </el-list-item>
      </el-list>
      <template #footer>
        <el-button @click="viewDialog = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Folder, Plus, Document } from '@element-plus/icons-vue'
import { collectionApi } from '@/api'
import { success, error, confirm } from '@/utils/message'

export default {
  name: 'CollectionList',
  setup() {
    const loading = ref(true)
    const errorMsg = ref('')
    const collections = ref([])
    const createDialog = ref(false)
    const editDialog = ref(false)
    const viewDialog = ref(false)
    const createForm = ref(null)
    const editForm = ref(null)
    const newCollection = ref({
      name: '',
      description: '',
      is_public: false
    })
    const editCollectionData = ref({})
    const currentCollection = ref({})
    const collectionArticles = ref([])

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

    const createCollection = async () => {
      const valid = await createForm.value?.validate().catch(() => false)
      if (!valid) return
      try {
        const res = await collectionApi.createCollection(newCollection.value)
        if (res.data.success) {
          collections.value.push(res.data.data)
          createDialog.value = false
          newCollection.value = { name: '', description: '', is_public: false }
          createForm.value?.resetFields()
          success('创建成功')
        }
      } catch (err) {
        error('创建收藏夹失败')
        console.error('创建收藏夹失败:', err)
      }
    }

    const editCollection = (collection) => {
      editCollectionData.value = { ...collection }
      editDialog.value = true
    }

    const updateCollection = async () => {
      const valid = await editForm.value?.validate().catch(() => false)
      if (!valid) return
      try {
        const res = await collectionApi.updateCollection(editCollectionData.value.id, editCollectionData.value)
        if (res.data.success) {
          const index = collections.value.findIndex(c => c.id === editCollectionData.value.id)
          if (index !== -1) {
            collections.value[index] = res.data.data
          }
          editDialog.value = false
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
      try {
        const res = await collectionApi.getCollection(collectionId)
        if (res.data.success) {
          currentCollection.value = res.data.data.collection
          collectionArticles.value = res.data.data.articles
          viewDialog.value = true
        }
      } catch (err) {
        error('加载收藏夹详情失败')
        console.error('加载收藏夹详情失败:', err)
      }
    }

    const removeArticle = async (articleId) => {
      try {
        await collectionApi.removeArticleFromCollection(currentCollection.value.id, articleId)
        collectionArticles.value = collectionArticles.value.filter(a => a.article.id !== articleId)
        success('移除成功')
      } catch (err) {
        error('移除文章失败')
        console.error('移除文章失败:', err)
      }
    }

    onMounted(() => {
      loadCollections()
    })

    return {
      loading,
      errorMsg,
      collections,
      createDialog,
      editDialog,
      viewDialog,
      createForm,
      editForm,
      newCollection,
      editCollectionData,
      currentCollection,
      collectionArticles,
      loadCollections,
      createCollection,
      editCollection,
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
