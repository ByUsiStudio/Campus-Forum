<template>
  <v-container>
    <v-card>
      <v-card-title class="d-flex align-center">
        <v-icon left>mdi-folder-star</v-icon>
        我的收藏夹
        <v-btn color="primary" class="ml-auto" @click="createDialog = true">
          <v-icon left>mdi-plus</v-icon>
          新建收藏夹
        </v-btn>
      </v-card-title>

      <v-card-text>
        <v-row>
          <v-col cols="12" sm="6" md="4" v-for="collection in collections" :key="collection.id">
            <v-card outlined hover>
              <v-card-text>
                <div class="d-flex align-center">
                  <v-avatar size="40" color="primary">
                    <v-icon dark>mdi-folder</v-icon>
                  </v-avatar>
                  <div class="ml-3">
                    <div class="title">{{ collection.name }}</div>
                    <div class="caption">{{ collection.description }}</div>
                  </div>
                </div>
                
                <v-row class="mt-2">
                  <v-col cols="6">
                    <v-chip small>
                      <v-icon small left>mdi-file-document</v-icon>
                      {{ collection.article_count }} 文章
                    </v-chip>
                  </v-col>
                  <v-col cols="6">
                    <v-chip small :color="collection.is_public ? 'success' : 'grey'">
                      {{ collection.is_public ? '公开' : '私密' }}
                    </v-chip>
                  </v-col>
                </v-row>
              </v-card-text>

              <v-card-actions>
                <v-btn text small color="primary" @click="viewCollection(collection.id)">
                  查看
                </v-btn>
                <v-btn text small color="warning" @click="editCollection(collection)">
                  编辑
                </v-btn>
                <v-btn text small color="error" @click="deleteCollection(collection.id)">
                  删除
                </v-btn>
              </v-card-actions>
            </v-card>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>

    <!-- 创建收藏夹对话框 -->
    <v-dialog v-model="createDialog" max-width="500">
      <v-card>
        <v-card-title>新建收藏夹</v-card-title>
        <v-card-text>
          <v-form ref="createForm">
            <v-text-field
              v-model="newCollection.name"
              label="收藏夹名称"
              outlined
              :rules="[v => !!v || '请输入名称']"
            ></v-text-field>
            <v-textarea
              v-model="newCollection.description"
              label="描述"
              outlined
              rows="3"
            ></v-textarea>
            <v-switch
              v-model="newCollection.is_public"
              label="公开收藏夹"
            ></v-switch>
          </v-form>
        </v-card-text>
        <v-card-actions>
          <v-btn text @click="createDialog = false">取消</v-btn>
          <v-btn color="primary" @click="createCollection">创建</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- 编辑收藏夹对话框 -->
    <v-dialog v-model="editDialog" max-width="500">
      <v-card>
        <v-card-title>编辑收藏夹</v-card-title>
        <v-card-text>
          <v-form ref="editForm">
            <v-text-field
              v-model="editCollectionData.name"
              label="收藏夹名称"
              outlined
              :rules="[v => !!v || '请输入名称']"
            ></v-text-field>
            <v-textarea
              v-model="editCollectionData.description"
              label="描述"
              outlined
              rows="3"
            ></v-textarea>
            <v-switch
              v-model="editCollectionData.is_public"
              label="公开收藏夹"
            ></v-switch>
          </v-form>
        </v-card-text>
        <v-card-actions>
          <v-btn text @click="editDialog = false">取消</v-btn>
          <v-btn color="primary" @click="updateCollection">保存</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- 收藏夹详情对话框 -->
    <v-dialog v-model="viewDialog" max-width="800">
      <v-card>
        <v-card-title>{{ currentCollection.name }}</v-card-title>
        <v-card-text>
          <v-list three-line>
            <v-list-item v-for="item in collectionArticles" :key="item.id">
              <v-list-item-content>
                <v-list-item-title>{{ item.article.title }}</v-list-item-title>
                <v-list-item-subtitle>
                  {{ item.note || '无备注' }}
                </v-list-item-subtitle>
              </v-list-item-content>
              <v-list-item-action>
                <v-btn text small color="primary" :to="`/article/${item.article.id}`">
                  查看
                </v-btn>
                <v-btn text small color="error" @click="removeArticle(item.article.id)">
                  移除
                </v-btn>
              </v-list-item-action>
            </v-list-item>
          </v-list>
        </v-card-text>
        <v-card-actions>
          <v-btn text @click="viewDialog = false">关闭</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
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