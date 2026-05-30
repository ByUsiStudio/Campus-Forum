<template>
  <div class="system-panel">
    <div class="panel-header">
      <div class="header-left">
        <h2 class="panel-title">系统设置</h2>
        <p class="panel-subtitle">网站配置、邮件服务与通知管理</p>
      </div>
    </div>

    <v-tabs v-model="systemTab" class="system-tabs" color="primary">
      <v-tab value="site">
        <v-icon start>mdi-web</v-icon>
        网站配置
      </v-tab>
      <v-tab value="smtp">
        <v-icon start>mdi-email</v-icon>
        邮件配置
      </v-tab>
      <v-tab value="notifications">
        <v-icon start>mdi-bell</v-icon>
        通知管理
      </v-tab>
    </v-tabs>

    <v-window v-model="systemTab" class="system-window">
      <v-window-item value="site">
        <v-card class="config-card">
          <v-card-title class="config-title">
            <v-icon class="title-icon">mdi-web</v-icon>
            网站配置
          </v-card-title>
          <v-card-text>
            <v-text-field
              v-model="localSiteConfig.siteTitle"
              label="网站标题"
              placeholder="校园论坛"
              variant="outlined"
              class="mb-4"
            ></v-text-field>
            <v-btn color="primary" size="large" @click="$emit('save-site', localSiteConfig)">
              <v-icon start>mdi-content-save</v-icon>
              保存配置
            </v-btn>
          </v-card-text>
        </v-card>
      </v-window-item>

      <v-window-item value="smtp">
        <v-card class="config-card">
          <v-card-title class="config-title">
            <v-icon class="title-icon">mdi-email</v-icon>
            SMTP 邮件配置
          </v-card-title>
          <v-card-text>
            <v-row>
              <v-col cols="12" sm="6">
                <v-text-field
                  v-model="localSmtpConfig.host"
                  label="SMTP 服务器"
                  placeholder="smtp.example.com"
                  variant="outlined"
                  class="mb-4"
                ></v-text-field>
              </v-col>
              <v-col cols="12" sm="6">
                <v-text-field
                  v-model="localSmtpConfig.port"
                  label="端口"
                  type="number"
                  variant="outlined"
                  class="mb-4"
                ></v-text-field>
              </v-col>
              <v-col cols="12" sm="6">
                <v-text-field
                  v-model="localSmtpConfig.username"
                  label="用户名"
                  variant="outlined"
                  class="mb-4"
                ></v-text-field>
              </v-col>
              <v-col cols="12" sm="6">
                <v-text-field
                  v-model="localSmtpConfig.password"
                  label="密码"
                  variant="outlined"
                  type="password"
                  class="mb-4"
                ></v-text-field>
              </v-col>
              <v-col cols="12" sm="6">
                <v-text-field
                  v-model="localSmtpConfig.from"
                  label="发件人邮箱"
                  variant="outlined"
                  class="mb-4"
                ></v-text-field>
              </v-col>
              <v-col cols="12" sm="6">
                <v-text-field
                  v-model="localSmtpConfig.fromName"
                  label="发件人名称"
                  variant="outlined"
                  class="mb-4"
                ></v-text-field>
              </v-col>
            </v-row>
            <div class="smtp-actions">
              <v-btn color="primary" size="large" @click="$emit('save-smtp', localSmtpConfig)">
                <v-icon start>mdi-content-save</v-icon>
                保存配置
              </v-btn>
              <v-btn variant="outlined" size="large" @click="$emit('test-smtp', localSmtpConfig)">
                <v-icon start>mdi-send</v-icon>
                发送测试邮件
              </v-btn>
            </div>
          </v-card-text>
        </v-card>
      </v-window-item>

      <v-window-item value="notifications">
        <v-card class="config-card">
          <v-card-title class="config-title">
            <v-icon class="title-icon">mdi-bell</v-icon>
            发送通知
          </v-card-title>
          <v-card-text>
            <v-row>
              <v-col cols="12" sm="6">
                <v-select
                  v-model="localNotificationForm.type"
                  :items="notificationTypes"
                  item-title="title"
                  item-value="value"
                  label="通知类型"
                  variant="outlined"
                  class="mb-4"
                ></v-select>
              </v-col>
              <v-col cols="12" sm="6">
                <v-text-field
                  v-model="localNotificationForm.title"
                  label="通知标题"
                  variant="outlined"
                  class="mb-4"
                ></v-text-field>
              </v-col>
              <v-col cols="12">
                <v-textarea
                  v-model="localNotificationForm.content"
                  label="通知内容"
                  variant="outlined"
                  rows="4"
                  class="mb-4"
                ></v-textarea>
              </v-col>
            </v-row>
            <v-btn color="primary" size="large" @click="$emit('send-notification')">
              <v-icon start>mdi-send</v-icon>
              发送通知
            </v-btn>
          </v-card-text>
        </v-card>

        <v-card class="mt-6">
          <v-card-title class="config-title">
            <v-icon class="title-icon">mdi-history</v-icon>
            历史通知
          </v-card-title>
          <v-card-text>
            <div v-if="notifications.length === 0" class="empty-state">
              <v-icon size="40" color="grey">mdi-bell-off</v-icon>
              <div class="empty-text">暂无通知记录</div>
            </div>

            <v-list v-else class="notification-list">
              <v-list-item v-for="notif in notifications" :key="notif.id" class="notification-item">
                <template #prepend>
                  <v-icon :color="getNotificationColor(notif.type)">{{ getNotificationIcon(notif.type) }}</v-icon>
                </template>
                <v-list-item-title class="notification-title">{{ notif.title }}</v-list-item-title>
                <v-list-item-subtitle class="notification-content">{{ notif.content }}</v-list-item-subtitle>
                <template #append>
                  <v-btn variant="text" size="small" color="error" @click="$emit('delete-notification', notif.id)">
                    <v-icon>mdi-delete</v-icon>
                  </v-btn>
                </template>
              </v-list-item>
            </v-list>
          </v-card-text>
        </v-card>
      </v-window-item>
    </v-window>
  </div>
</template>

<script>
import { ref, watch } from 'vue'

export default {
  name: 'SystemPanel',
  props: {
    siteConfig: {
      type: Object,
      default: () => ({ siteTitle: '' })
    },
    smtpConfig: {
      type: Object,
      default: () => ({})
    },
    notifications: {
      type: Array,
      default: () => []
    },
    notificationForm: {
      type: Object,
      default: () => ({})
    },
    notificationTypes: {
      type: Array,
      default: () => []
    }
  },
  emits: [
    'save-site',
    'save-smtp',
    'test-smtp',
    'send-notification',
    'delete-notification',
    'refresh-site',
    'refresh-smtp',
    'refresh-notifications'
  ],
  setup(props) {
    const systemTab = ref('site')
    const localSiteConfig = ref({ ...props.siteConfig })
    const localSmtpConfig = ref({ ...props.smtpConfig })
    const localNotificationForm = ref({ ...props.notificationForm })

    watch(() => props.siteConfig, (val) => {
      localSiteConfig.value = { ...val }
    }, { deep: true })

    watch(() => props.smtpConfig, (val) => {
      localSmtpConfig.value = { ...val }
    }, { deep: true })

    watch(() => props.notificationForm, (val) => {
      localNotificationForm.value = { ...val }
    }, { deep: true })

    const getNotificationColor = (type) => {
      const colors = {
        system: 'primary',
        activity: 'success',
        update: 'info',
        warning: 'warning'
      }
      return colors[type] || 'default'
    }

    const getNotificationIcon = (type) => {
      const icons = {
        system: 'mdi-information',
        activity: 'mdi-calendar',
        update: 'mdi-update',
        warning: 'mdi-alert'
      }
      return icons[type] || 'mdi-bell'
    }

    return {
      systemTab,
      localSiteConfig,
      localSmtpConfig,
      localNotificationForm,
      getNotificationColor,
      getNotificationIcon
    }
  }
}
</script>

<style scoped>
.system-panel {
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 16px;
}

.panel-title {
  font-size: 1.5rem;
  font-weight: 700;
  color: #1a1a2e;
  margin: 0 0 4px 0;
}

.panel-subtitle {
  font-size: 0.9rem;
  color: #6b7280;
  margin: 0;
}

.system-tabs {
  margin-bottom: 24px;
  background: white;
  border-radius: 12px;
  padding: 8px;
}

.system-window {
  margin-top: 0;
}

.config-card {
  border-radius: 16px;
}

.config-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 1rem;
  font-weight: 600;
  padding: 20px 24px !important;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
}

.title-icon {
  width: 36px;
  height: 36px;
  padding: 6px;
  border-radius: 8px;
  background: rgba(103, 80, 164, 0.1);
  color: #6750A4 !important;
}

.smtp-actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 40px 20px;
}

.empty-text {
  margin-top: 12px;
  font-size: 0.9rem;
  color: #9ca3af;
}

.notification-list {
  padding: 0;
}

.notification-item {
  border-radius: 10px;
  margin-bottom: 8px;
  background: #f8f9ff;
}

.notification-title {
  font-weight: 600;
}

.notification-content {
  margin-top: 4px;
}

@media (max-width: 600px) {
  .smtp-actions {
    flex-direction: column;
  }

  .smtp-actions .v-btn {
    width: 100%;
  }
}
</style>
