<template>
  <div class="container collection-list">
    <div class="layui-card">
      <div class="layui-card-header d-flex align-center justify-between">
        <div class="d-flex align-center">
          <i class="fa-solid fa-folder-star mr-2"></i>
          <span class="font-weight-bold text-lg">我的收藏夹</span>
        </div>
        <button class="layui-btn layui-btn-primary" @click="createDialog = true">
          <i class="fa-solid fa-plus mr-1"></i>
          新建收藏夹
        </button>
      </div>

      <div class="layui-card-body">
        <div class="collection-grid">
          <div class="collection-item" v-for="collection in collections" :key="collection.id">
            <div class="layui-card">
              <div class="layui-card-body">
                <div class="d-flex align-center">
                  <div class="avatar primary">
                    <i class="fa-solid fa-folder"></i>
                  </div>
                  <div class="ml-3 flex-1">
                    <div class="title font-weight-bold">{{ collection.name }}</div>
                    <div class="caption text-muted">{{ collection.description }}</div>
                  </div>
                </div>
                
                <div class="mt-2 d-flex gap-2">
                  <span class="layui-badge">
                    <i class="fa-solid fa-file-lines mr-1"></i>
                    {{ collection.article_count }} 文章
                  </span>
                  <span class="layui-badge" :class="collection.is_public ? 'layui-bg-green' : 'layui-bg-gray'">
                    {{ collection.is_public ? '公开' : '私密' }}
                  </span>
                </div>
              </div>

              <div class="layui-card-footer d-flex justify-end gap-2">
                <button class="layui-btn layui-btn-xs" @click="viewCollection(collection.id)">
                  查看
                </button>
                <button class="layui-btn layui-btn-xs layui-btn-warm" @click="editCollection(collection)">
                  编辑
                </button>
                <button class="layui-btn layui-btn-xs layui-btn-danger" @click="deleteCollection(collection.id)">
                  删除
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 创建收藏夹对话框 -->
    <div class="modal-overlay" v-if="createDialog" @click.self="createDialog = false">
      <div class="modal-content animate-scale-in">
        <div class="layui-card">
          <div class="layui-card-header">新建收藏夹</div>
          <div class="layui-card-body">
            <form ref="createForm" class="layui-form">
              <div class="layui-form-item">
                <label class="layui-form-label">收藏夹名称</label>
                <div class="layui-input-block">
                  <input type="text" v-model="newCollection.name" class="layui-input" placeholder="请输入名称" />
                </div>
              </div>
              <div class="layui-form-item">
                <label class="layui-form-label">描述</label>
                <div class="layui-input-block">
                  <textarea v-model="newCollection.description" class="layui-textarea" rows="3" placeholder="请输入描述"></textarea>
                </div>
              </div>
              <div class="layui-form-item">
                <div class="layui-input-block">
                  <input type="checkbox" v-model="newCollection.is_public" lay-skin="switch" title="公开收藏夹" />
                </div>
              </div>
            </form>
          </div>
          <div class="layui-card-footer d-flex justify-end gap-2">
            <button class="layui-btn layui-btn-primary" @click="createDialog = false">取消</button>
            <button class="layui-btn" @click="createCollection">创建</button>
          </div>
        </div>
      </div>
    </div>

    <!-- 编辑收藏夹对话框 -->
    <div class="modal-overlay" v-if="editDialog" @click.self="editDialog = false">
      <div class="modal-content animate-scale-in">
        <div class="layui-card">
          <div class="layui-card-header">编辑收藏夹</div>
          <div class="layui-card-body">
            <form ref="editForm" class="layui-form">
              <div class="layui-form-item">
                <label class="layui-form-label">收藏夹名称</label>
                <div class="layui-input-block">
                  <input type="text" v-model="editCollectionData.name" class="layui-input" placeholder="请输入名称" />
                </div>
              </div>
              <div class="layui-form-item">
                <label class="layui-form-label">描述</label>
                <div class="layui-input-block">
                  <textarea v-model="editCollectionData.description" class="layui-textarea" rows="3" placeholder="请输入描述"></textarea>
                </div>
              </div>
              <div class="layui-form-item">
                <div class="layui-input-block">
                  <input type="checkbox" v-model="editCollectionData.is_public" lay-skin="switch" title="公开收藏夹" />
                </div>
              </div>
            </form>
          </div>
          <div class="layui-card-footer d-flex justify-end gap-2">
            <button class="layui-btn layui-btn-primary" @click="editDialog = false">取消</button>
            <button class="layui-btn" @click="updateCollection">保存</button>
          </div>
        </div>
      </div>
    </div>

    <!-- 收藏夹详情对话框 -->
    <div class="modal-overlay" v-if="viewDialog" @click.self="viewDialog = false">
      <div class="modal-content modal-large animate-scale-in">
        <div class="layui-card">
          <div class="layui-card-header">{{ currentCollection.name }}</div>
          <div class="layui-card-body">
            <ul class="collection-detail-list">
              <li class="detail-item" v-for="item in collectionArticles" :key="item.id">
                <div class="detail-content">
                  <div class="detail-title">{{ item.article.title }}</div>
                  <div class="detail-subtitle">
                    {{ item.note || '无备注' }}
                  </div>
                </div>
                <div class="detail-actions">
                  <button class="layui-btn layui-btn-xs" :to="`/article/${item.article.id}`">
                    查看
                  </button>
                  <button class="layui-btn layui-btn-xs layui-btn-danger" @click="removeArticle(item.article.id)">
                    移除
                  </button>
                </div>
              </li>
            </ul>
          </div>
          <div class="layui-card-footer d-flex justify-end">
            <button class="layui-btn layui-btn-primary" @click="viewDialog = false">关闭</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { collectionApi } from '@/api'

export default {
  name: 'CollectionList',
  setup() {
    const collections = ref([])
    const createDialog = ref(false)
    const editDialog = ref(false)
    const viewDialog = ref(false)
    const newCollection = ref({
      name: '',
      description: '',
      is_public: false
    })
    const editCollectionData = ref({})
    const currentCollection = ref({})
    const collectionArticles = ref([])

    const loadCollections = async () => {
      try {
        const res = await collectionApi.getCollections()
        if (res.data.success) {
          collections.value = res.data.data
        }
      } catch (error) {
        console.error('加载收藏夹失败:', error)
      }
    }

    const createCollection = async () => {
      try {
        const res = await collectionApi.createCollection(newCollection.value)
        if (res.data.success) {
          collections.value.push(res.data.data)
          createDialog.value = false
          newCollection.value = { name: '', description: '', is_public: false }
        }
      } catch (error) {
        console.error('创建收藏夹失败:', error)
      }
    }

    const editCollection = (collection) => {
      editCollectionData.value = { ...collection }
      editDialog.value = true
    }

    const updateCollection = async () => {
      try {
        const res = await collectionApi.updateCollection(editCollectionData.value.id, editCollectionData.value)
        if (res.data.success) {
          const index = collections.value.findIndex(c => c.id === editCollectionData.value.id)
          if (index !== -1) {
            collections.value[index] = res.data.data
          }
          editDialog.value = false
        }
      } catch (error) {
        console.error('更新收藏夹失败:', error)
      }
    }

    const deleteCollection = async (collectionId) => {
      try {
        await collectionApi.deleteCollection(collectionId)
        collections.value = collections.value.filter(c => c.id !== collectionId)
      } catch (error) {
        console.error('删除收藏夹失败:', error)
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
      } catch (error) {
        console.error('加载收藏夹详情失败:', error)
      }
    }

    const removeArticle = async (articleId) => {
      try {
        await collectionApi.removeArticleFromCollection(currentCollection.value.id, articleId)
        collectionArticles.value = collectionArticles.value.filter(a => a.article.id !== articleId)
      } catch (error) {
        console.error('移除文章失败:', error)
      }
    }

    onMounted(() => {
      loadCollections()
    })

    return {
      collections,
      createDialog,
      editDialog,
      viewDialog,
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
      removeArticle
    }
  }
}
</script>

<style scoped>
.collection-list {
  padding: 20px;
}

.collection-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 15px;
}

.collection-item {
  transition: transform 0.3s ease;
}

.collection-item:hover {
  transform: translateY(-4px);
}

.avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  
  &.primary {
    background: var(--primary);
  }
  
  i {
    font-size: 18px;
  }
}

.collection-detail-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.detail-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 0;
  border-bottom: 1px solid #eee;
  
  &:last-child {
    border-bottom: none;
  }
}

.detail-content {
  flex: 1;
}

.detail-title {
  font-weight: 600;
  color: #333;
}

.detail-subtitle {
  color: #999;
  font-size: 14px;
  margin-top: 4px;
}

.detail-actions {
  display: flex;
  gap: 8px;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 12px;
  width: 90%;
  max-width: 500px;
  overflow: hidden;
  
  &.modal-large {
    max-width: 800px;
  }
}

.d-flex {
  display: flex;
}

.align-center {
  align-items: center;
}

.justify-between {
  justify-content: space-between;
}

.justify-end {
  justify-content: flex-end;
}

.flex-1 {
  flex: 1;
}

.mr-1 {
  margin-right: 5px;
}

.mr-2 {
  margin-right: 10px;
}

.ml-3 {
  margin-left: 15px;
}

.mt-2 {
  margin-top: 10px;
}

.gap-2 {
  gap: 10px;
}

.font-weight-bold {
  font-weight: 600;
}

.text-lg {
  font-size: 18px;
}

.text-muted {
  color: #999;
}
</style>