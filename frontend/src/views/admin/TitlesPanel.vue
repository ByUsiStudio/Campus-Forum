<template>
  <div>
    <!-- 添加新头衔表单 -->
    <v-card class="mb-4" variant="flat" rounded="lg">
      <v-card-title class="pb-2">
        <v-icon start size="20">mdi-plus-circle</v-icon>
        添加新头衔
      </v-card-title>
      <v-card-text>
        <v-row dense>
          <v-col cols="12" sm="6">
            <v-text-field
              v-model="titleForm.name"
              label="头衔名称"
              placeholder="例如：技术达人"
              variant="outlined"
              density="compact"
              prepend-inner-icon="mdi-star"
              hide-details
            />
          </v-col>
          <v-col cols="12" sm="6">
            <v-text-field
              v-model="titleForm.description"
              label="描述"
              placeholder="头衔简介"
              variant="outlined"
              density="compact"
              prepend-inner-icon="mdi-text"
              hide-details
            />
          </v-col>
          <v-col cols="6" sm="3">
            <v-text-field
              v-model="titleForm.color"
              label="颜色"
              variant="outlined"
              density="compact"
              hide-details
            />
          </v-col>
          <v-col cols="6" sm="3">
            <v-text-field
              v-model="titleForm.icon"
              label="图标"
              placeholder="mdi-star"
              variant="outlined"
              density="compact"
              hide-details
            />
          </v-col>
          <v-col cols="12" sm="6">
            <div class="d-flex align-center justify-center pa-2">
              <v-chip :color="titleForm.color" size="medium">
                <v-icon v-if="titleForm.icon" start size="18">{{ titleForm.icon }}</v-icon>
                {{ titleForm.name || '预览' }}
              </v-chip>
            </div>
          </v-col>
          <v-col cols="12" sm="6">
            <v-btn color="primary" block height="40" @click="$emit('add-title', titleForm)">
              <v-icon start>mdi-plus</v-icon>
              添加头衔
            </v-btn>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>

    <!-- 授予头衔表单 -->
    <v-card class="mb-4" variant="flat" rounded="lg">
      <v-card-title class="pb-2">
        <v-icon start size="20">mdi-gift</v-icon>
        授予头衔
      </v-card-title>
      <v-card-text>
        <v-row dense>
          <v-col cols="12" sm="4">
            <v-select
              v-model="grantForm.user_id"
              :items="users"
              item-title="display_name"
              item-value="id"
              label="选择用户"
              variant="outlined"
              density="compact"
              hide-details
            />
          </v-col>
          <v-col cols="12" sm="4">
            <v-select
              v-model="grantForm.title_id"
              :items="titles"
              item-title="name"
              item-value="id"
              label="选择头衔"
              variant="outlined"
              density="compact"
              hide-details
            />
          </v-col>
          <v-col cols="12" sm="4">
            <v-btn color="success" block height="40" @click="$emit('grant', grantForm)">
              <v-icon start>mdi-check</v-icon>
              授予头衔
            </v-btn>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>

    <!-- 已创建的头衔列表 -->
    <v-card class="mb-4" variant="flat" rounded="lg">
      <v-card-title class="pb-2">
        <v-icon start size="20">mdi-medal</v-icon>
        已创建的头衔
      </v-card-title>
      <v-list lines="two" v-if="titles.length > 0">
        <v-list-item v-for="title in titles" :key="title.id" class="py-3">
          <template v-slot:prepend>
            <v-avatar size="48" :color="title.color || 'primary'" variant="tonal">
              <v-icon v-if="title.icon">{{ title.icon }}</v-icon>
              <v-icon v-else>mdi-medal</v-icon>
            </v-avatar>
          </template>

          <v-list-item-title class="font-weight-medium mb-1">
            <v-chip :color="title.color" size="small" class="mr-2">
              <v-icon v-if="title.icon" start size="16">{{ title.icon }}</v-icon>
              {{ title.name }}
            </v-chip>
          </v-list-item-title>

          <v-list-item-subtitle>
            {{ title.description || '暂无描述' }}
          </v-list-item-subtitle>

          <template v-slot:append>
            <v-btn size="small" color="error" variant="text" @click="$emit('delete-title', title)">
              <v-icon>mdi-delete</v-icon>
              <v-tooltip activator="parent">删除</v-tooltip>
            </v-btn>
          </template>
        </v-list-item>
      </v-list>

      <v-card-text v-else class="text-center py-8">
        <v-icon size="48" color="grey-lighten-1">mdi-medal-outline</v-icon>
        <div class="text-body-1 text-medium-emphasis mt-2">
          暂无头衔数据
        </div>
      </v-card-text>
    </v-card>

    <!-- 用户头衔一览 -->
    <v-card variant="flat" rounded="lg">
      <v-card-title class="pb-2">
        <v-icon start size="20">mdi-account-multiple</v-icon>
        用户头衔一览
      </v-card-title>
      <v-list lines="two" v-if="usersWithTitles.length > 0">
        <v-list-item v-for="user in usersWithTitles" :key="user.id" class="py-3">
          <template v-slot:prepend>
            <UserAvatar :user="user" :size="48" />
          </template>

          <v-list-item-title class="font-weight-medium mb-1">
            {{ user.display_name }}
            <span class="text-caption text-medium-emphasis ml-2">@{{ user.username }}</span>
          </v-list-item-title>

          <v-list-item-subtitle>
            <div class="d-flex flex-wrap ga-1">
              <v-chip
                v-for="userTitle in user.titles"
                :key="userTitle.id"
                :color="userTitle.color"
                size="small"
                class="mr-1 mb-1"
              >
                <v-icon v-if="userTitle.icon" start size="14">{{ userTitle.icon }}</v-icon>
                {{ userTitle.name }}
                <v-btn
                  variant="text"
                  size="x-small"
                  icon
                  color="error"
                  @click="$emit('revoke', user.id, userTitle.id)"
                  class="ml-1"
                >
                  <v-icon size="14">mdi-close-circle</v-icon>
                  <v-tooltip activator="parent">撤销</v-tooltip>
                </v-btn>
              </v-chip>
            </div>
          </v-list-item-subtitle>
        </v-list-item>
      </v-list>

      <v-card-text v-else class="text-center py-8">
        <v-icon size="48" color="grey-lighten-1">mdi-account-off</v-icon>
        <div class="text-body-1 text-medium-emphasis mt-2">
          暂无用户头衔数据
        </div>
      </v-card-text>
    </v-card>
  </div>
</template>

<script>
import { computed, reactive } from 'vue'
import UserAvatar from '../../components/UserAvatar.vue'

export default {
  name: 'TitlesPanel',
  components: {
    UserAvatar
  },
  props: {
    titles: {
      type: Array,
      default: () => []
    },
    users: {
      type: Array,
      default: () => []
    },
    loading: {
      type: Boolean,
      default: false
    }
  },
  emits: ['add-title', 'grant', 'revoke', 'delete-title', 'refresh'],
  setup(props) {
    const titleForm = reactive({
      name: '',
      description: '',
      color: '#6750A4',
      icon: ''
    })

    const grantForm = reactive({
      user_id: null,
      title_id: null,
      reason: ''
    })

    const usersWithTitles = computed(() => {
      return props.users.filter(user => user.titles && user.titles.length > 0)
    })

    return {
      titleForm,
      grantForm,
      usersWithTitles
    }
  }
}
</script>