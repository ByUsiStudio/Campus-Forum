<template>
  <div v-if="isInitialized && isAdmin" class="admin-layout">
    <v-navigation-drawer
      v-model="drawer"
      temporary
      fixed
      app
      class="admin-sidebar"
    >
      <v-list-item class="px-4">
        <v-list-item-content>
          <v-list-item-title class="text-h6" style="color: rgb(var(--v-theme-primary));">管理后台</v-list-item-title>
        </v-list-item-content>
      </v-list-item>
      <v-divider></v-divider>
      <v-list nav>
        <v-list-item
          v-for="item in menuItems"
          :key="item.value"
          :value="item.value"
          :active="activeTab === item.value"
          @click="activeTab = item.value"
        >
          <template v-slot:prepend>
            <v-icon>{{ item.icon }}</v-icon>
          </template>
          <v-list-item-title>{{ item.title }}</v-list-item-title>
          <v-badge
            v-if="item.badge"
            :content="item.badge"
            color="error"
            :model-value="item.badge > 0"
          ></v-badge>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>

    <v-app-bar app fixed class="admin-app-bar">
      <v-app-bar-nav-icon @click="drawer = !drawer"></v-app-bar-nav-icon>
      <v-toolbar-title style="color: rgb(var(--v-theme-primary));">管理后台</v-toolbar-title>
    </v-app-bar>

    <v-main class="admin-main">
      <div class="admin-content">
        <v-window-item value="overview">
          <v-row>
            <v-col cols="12" sm="6" md="3">
              <v-card color="primary" class="pa-4">
                <div class="text-h3 text-white">{{ statistics.user_count }}</div>
                <div class="text-body-1 text-white opacity-80">用户总数</div>
              </v-card>
            </v-col>
            <v-col cols="12" sm="6" md="3">
              <v-card color="success" class="pa-4">
                <div class="text-h3 text-white">{{ statistics.article_count }}</div>
                <div class="text-body-1 text-white opacity-80">文章总数</div>
              </v-card>
            </v-col>
            <v-col cols="12" sm="6" md="3">
              <v-card color="info" class="pa-4">
                <div class="text-h3 text-white">{{ statistics.comment_count }}</div>
                <div class="text-body-1 text-white opacity-80">评论总数</div>
              </v-card>
            </v-col>
            <v-col cols="12" sm="6" md="3">
              <v-card color="warning" class="pa-4">
                <div class="text-h3 text-white">{{ statistics.view_count }}</div>
                <div class="text-body-1 text-white opacity-80">总浏览量</div>
              </v-card>
            </v-col>
          </v-row>
        </v-window-item>

        <v-window-item value="users">
          <v-card variant="outlined" class="pa-4 mb-4">
            <div class="text-body-2 text-medium-emphasis mb-2">用户总数：{{ users.length }}</div>
          </v-card>

          <v-table>
            <thead>
              <tr>
                <th>ID</th>
                <th>用户</th>
                <th>QQ号码</th>
                <th>角色</th>
                <th>状态</th>
                <th>注册时间</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="user in users" :key="user.id">
                <td>{{ user.id }}</td>
                <td>
                  <UserAvatar :user="user" :size="32" />
                </td>
                <td>{{ user.qq_number }}</td>
                <td>
                  <v-chip size="small" :color="user.role === 'admin' ? 'error' : user.role === 'system' ? 'warning' : 'default'">
                    {{ user.role === 'admin' ? '管理员' : user.role === 'system' ? '系统管理员' : '用户' }}
                  </v-chip>
                </td>
                <td>
                  <v-chip size="small" :color="user.status === 'banned' ? 'error' : 'success'">
                    {{ user.status === 'banned' ? '已封禁' : '正常' }}
                  </v-chip>
                </td>
                <td>{{ formatDate(user.created_at) }}</td>
                <td>
                  <v-btn variant="text" size="small" color="primary" @click="showEditRoleDialog(user)" v-if="currentUserId && user.id !== currentUserId">
                    修改角色
                  </v-btn>
                  <v-btn variant="text" size="small" color="warning" @click="showBanDialog(user)" v-if="currentUserId && user.id !== currentUserId && user.status !== 'banned'">
                    封禁
                  </v-btn>
                  <v-btn variant="text" size="small" color="success" @click="handleUnban(user)" v-if="currentUserId && user.id !== currentUserId && user.status === 'banned'">
                    解封
                  </v-btn>
                  <v-btn variant="text" size="small" color="error" @click="handleDeleteUser(user)" v-if="currentUserId && user.id !== currentUserId">
                    删除
                  </v-btn>
                </td>
              </tr>
            </tbody>
          </v-table>
        </v-window-item>

        <v-window-item value="articles">
          <v-card variant="outlined" class="pa-4 mb-4">
            <v-row align="center">
              <v-col cols="12" md="4">
                <v-select
                  v-model="articleFilter"
                  :items="articleStatusOptions"
                  label="筛选状态"
                  variant="outlined"
                  density="compact"
                  hide-details
                ></v-select>
              </v-col>
              <v-col cols="12" md="8" class="text-right">
                <v-pagination
                  v-model="articlePage"
                  :length="articleTotalPages"
                  :total-visible="5"
                  density="compact"
                ></v-pagination>
              </v-col>
            </v-row>
          </v-card>

          <v-table>
            <thead>
              <tr>
                <th>ID</th>
                <th>标题</th>
                <th>作者</th>
                <th>分区</th>
                <th>点赞</th>
                <th>浏览</th>
                <th>状态</th>
                <th>发布时间</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="article in articles" :key="article.id">
                <td>{{ article.id }}</td>
                <td class="text-truncate" style="max-width: 200px;">{{ article.title }}</td>
                <td>{{ article.User?.display_name }}</td>
                <td>{{ article.Category?.name }}</td>
                <td>{{ article.like_count }}</td>
                <td>{{ article.view_count }}</td>
                <td>
                  <v-chip size="small" :color="getStatusColor(article.status)">
                    {{ getStatusText(article.status) }}
                  </v-chip>
                </td>
                <td>{{ formatDate(article.created_at) }}</td>
                <td>
                  <v-btn variant="text" size="small" color="primary" @click="showStatusDialog(article)">
                    修改状态
                  </v-btn>
                </td>
              </tr>
            </tbody>
          </v-table>
        </v-window-item>

        <v-window-item value="comments">
          <v-card variant="outlined" class="pa-4 mb-4">
            <v-row align="center">
              <v-col cols="12" md="6">
                <div class="text-body-2">评论总数：{{ commentTotal }}</div>
              </v-col>
              <v-col cols="12" md="6" class="text-right">
                <v-pagination
                  v-model="commentPage"
                  :length="commentTotalPages"
                  :total-visible="5"
                  density="compact"
                ></v-pagination>
              </v-col>
            </v-row>
          </v-card>

          <v-list lines="three">
            <v-list-item v-for="comment in allComments" :key="comment.id" class="px-0">
              <template v-slot:prepend>
                <UserAvatar :user="comment.User || {}" :size="40" />
              </template>

              <v-list-item-title>
                {{ comment.User?.display_name }}
                <span class="text-caption text-medium-emphasis ml-2">回复文章：{{ comment.Article?.title }}</span>
              </v-list-item-title>
              <v-list-item-subtitle class="mt-1">
                {{ comment.content }}
              </v-list-item-subtitle>
              <v-list-item-subtitle class="mt-1">
                {{ formatDate(comment.created_at) }}
              </v-list-item-subtitle>

              <template v-slot:append>
                <v-btn variant="text" size="small" color="error" @click="handleDeleteComment(comment.id)">
                  删除
                </v-btn>
              </template>
            </v-list-item>
          </v-list>
        </v-window-item>

        <v-window-item value="categories">
          <v-card variant="outlined" class="pa-4 mb-4">
            <v-card-title class="text-subtitle-1 pa-0 mb-4">添加新分区</v-card-title>
            <v-form @submit.prevent="addCategory">
              <v-row>
                <v-col cols="12" md="4">
                  <v-text-field
                    v-model="categoryForm.name"
                    label="分区名称"
                    variant="outlined"
                    density="compact"
                    hide-details
                  ></v-text-field>
                </v-col>
                <v-col cols="12" md="4">
                  <v-text-field
                    v-model="categoryForm.description"
                    label="描述"
                    variant="outlined"
                    density="compact"
                    hide-details
                  ></v-text-field>
                </v-col>
                <v-col cols="12" md="2">
                  <v-btn type="submit" color="primary" block>添加</v-btn>
                </v-col>
              </v-row>
            </v-form>
          </v-card>

          <v-list>
            <v-list-item v-for="cat in categories" :key="cat.id">
              <template v-slot:prepend>
                <v-avatar color="primary" size="40">
                  <span>{{ cat.sort_order || 0 }}</span>
                </v-avatar>
              </template>

              <v-list-item-title class="font-weight-bold">{{ cat.name }}</v-list-item-title>
              <v-list-item-subtitle>{{ cat.description || '无描述' }}</v-list-item-subtitle>

              <template v-slot:append>
                <v-btn variant="text" size="small" color="primary" @click="showEditCategoryDialog(cat)">
                  编辑
                </v-btn>
                <v-btn variant="text" size="small" color="error" @click="handleDeleteCategory(cat.id)">
                  删除
                </v-btn>
              </template>
            </v-list-item>
          </v-list>
        </v-window-item>

        <v-window-item value="titles">
          <v-card variant="outlined" class="pa-4 mb-4">
            <v-card-title class="text-subtitle-1 pa-0 mb-4">添加新头衔</v-card-title>
            <v-form @submit.prevent="addTitle">
              <v-row>
                <v-col cols="12" md="3">
                  <v-text-field
                    v-model="titleForm.name"
                    label="头衔名称"
                    variant="outlined"
                    density="compact"
                    hide-details
                  ></v-text-field>
                </v-col>
                <v-col cols="12" md="3">
                  <v-text-field
                    v-model="titleForm.description"
                    label="描述"
                    variant="outlined"
                    density="compact"
                    hide-details
                  ></v-text-field>
                </v-col>
                <v-col cols="12" md="2">
                  <v-text-field
                    v-model="titleForm.color"
                    label="颜色"
                    variant="outlined"
                    density="compact"
                    hide-details
                  ></v-text-field>
                </v-col>
                <v-col cols="12" md="2">
                  <v-text-field
                    v-model="titleForm.icon"
                    label="图标"
                    variant="outlined"
                    density="compact"
                    hide-details
                  ></v-text-field>
                </v-col>
                <v-col cols="12" md="2">
                  <v-btn type="submit" color="primary" block>添加</v-btn>
                </v-col>
              </v-row>
            </v-form>
          </v-card>

          <v-card variant="outlined" class="pa-4 mb-4">
            <v-card-title class="text-subtitle-1 pa-0 mb-4">授予头衔</v-card-title>
            <v-form @submit.prevent="grantTitle">
              <v-row>
                <v-col cols="12" md="4">
                  <v-select
                    v-model="grantForm.user_id"
                    :items="users"
                    item-title="display_name"
                    item-value="id"
                    label="选择用户"
                    variant="outlined"
                    density="compact"
                    hide-details
                  ></v-select>
                </v-col>
                <v-col cols="12" md="4">
                  <v-select
                    v-model="grantForm.title_id"
                    :items="titles"
                    item-title="name"
                    item-value="id"
                    label="选择头衔"
                    variant="outlined"
                    density="compact"
                    hide-details
                  ></v-select>
                </v-col>
                <v-col cols="12" md="2">
                  <v-text-field
                    v-model="grantForm.reason"
                    label="授予原因"
                    variant="outlined"
                    density="compact"
                    hide-details
                  ></v-text-field>
                </v-col>
                <v-col cols="12" md="2">
                  <v-btn type="submit" color="primary" block>授予</v-btn>
                </v-col>
              </v-row>
            </v-form>
          </v-card>

          <v-list>
            <v-list-item v-for="title in titles" :key="title.id">
              <template v-slot:prepend>
                <v-chip :color="title.color" size="small">
                  <v-icon v-if="title.icon" start size="x-small">{{ title.icon }}</v-icon>
                  {{ title.name }}
                </v-chip>
              </template>

              <v-list-item-title>{{ title.name }}</v-list-item-title>
              <v-list-item-subtitle>{{ title.description || '无描述' }}</v-list-item-subtitle>

              <template v-slot:append>
                <v-btn variant="text" size="small" color="error" @click="handleDeleteTitle(title.id)">
                  删除
                </v-btn>
              </template>
            </v-list-item>
          </v-list>
        </v-window-item>

        <v-window-item value="sidebar">
          <v-card-text class="pa-0">
            <p class="text-body-2 text-medium-emphasis mb-4">
              配置侧边栏链接列表
            </p>

            <div v-for="(item, index) in sidebarItems" :key="index" class="d-flex gap-2 mb-3 align-center">
              <v-text-field
                v-model="item.title"
                label="标题"
                variant="outlined"
                density="compact"
                hide-details
                class="flex-grow-1"
              ></v-text-field>
              <v-text-field
                v-model="item.link"
                label="链接"
                variant="outlined"
                density="compact"
                hide-details
                class="flex-grow-1"
              ></v-text-field>
              <v-text-field
                v-model="item.icon"
                label="图标"
                variant="outlined"
                density="compact"
                hide-details
                style="max-width: 120px;"
              ></v-text-field>
              <v-btn icon variant="text" color="error" @click="removeSidebarItem(index)">
                <v-icon>mdi-delete</v-icon>
              </v-btn>
            </div>

            <div class="d-flex gap-2 mt-4">
              <v-btn variant="outlined" @click="addSidebarItem">
                <v-icon start>mdi-plus</v-icon>
                添加链接
              </v-btn>
              <v-btn color="primary" @click="saveSidebarConfig">
                <v-icon start>mdi-content-save</v-icon>
                保存配置
              </v-btn>
            </div>
          </v-card-text>
        </v-window-item>

        <v-window-item value="deletions">
          <div v-if="deletionRequests.length === 0" class="text-center pa-8 text-medium-emphasis">
            暂无待审核申请
          </div>

          <v-card v-for="req in deletionRequests" :key="req.id" class="mb-4 pa-4" variant="outlined">
            <v-card-text>
              <div class="d-flex justify-space-between align-start">
                <div>
                  <div class="text-h6 mb-2">{{ req.article?.title || '文章已删除' }}</div>
                  <v-list-item-subtitle class="mb-1">
                    <v-icon size="small">mdi-account</v-icon>
                    申请人：{{ req.user?.display_name }}
                  </v-list-item-subtitle>
                  <v-list-item-subtitle class="mb-1">
                    <v-icon size="small">mdi-delete</v-icon>
                    删除原因：{{ req.reason }}
                  </v-list-item-subtitle>
                  <v-list-item-subtitle>
                    <v-icon size="small">mdi-clock</v-icon>
                    申请时间：{{ formatDate(req.created_at) }}
                  </v-list-item-subtitle>
                </div>
                <div class="d-flex gap-2">
                  <v-btn color="primary" variant="flat" @click="approveDeletion(req.id)">
                    批准删除
                  </v-btn>
                  <v-btn color="secondary" variant="outlined" @click="rejectDeletion(req.id)">
                    拒绝
                  </v-btn>
                </div>
              </div>
            </v-card-text>
          </v-card>
        </v-window-item>

        <v-window-item value="announcement">
          <v-card-text class="pa-0">
            <v-textarea
              v-model="announcementContent"
              label="公告内容（支持Markdown）"
              variant="outlined"
              rows="10"
              placeholder="输入公告内容..."
              class="mb-4"
            ></v-textarea>

            <v-btn color="primary" @click="saveAnnouncement" class="mb-4">
              <v-icon start>mdi-content-save</v-icon>
              保存公告
            </v-btn>
          </v-card-text>
        </v-window-item>

        <v-window-item value="siteconfig">
          <v-card-text class="pa-0">
            <v-alert type="info" variant="tonal" class="mb-4">
              <v-icon start>mdi-information</v-icon>
              修改网站标题后，用户将在浏览器标签页和PWA安装提示中看到新的标题
            </v-alert>

            <v-text-field
              v-model="siteConfigForm.siteTitle"
              label="网站标题"
              variant="outlined"
              hint="显示在浏览器标签页和PWA安装提示中的标题"
              persistent-hint
              class="mb-4"
              prepend-inner-icon="mdi-format-title"
            ></v-text-field>

            <v-btn color="primary" @click="saveSiteConfig">
              <v-icon start>mdi-content-save</v-icon>
              保存网站配置
            </v-btn>
          </v-card-text>
        </v-window-item>

        <v-window-item value="smtpconfig">
          <v-card-text class="pa-0">
            <v-alert type="info" variant="tonal" class="mb-4">
              <v-icon start>mdi-information</v-icon>
              SMTP配置用于发送找回密码等邮件通知，配置完成后请测试发送
            </v-alert>

            <v-text-field
              v-model="smtpConfigForm.host"
              label="SMTP服务器地址"
              variant="outlined"
              hint="如: smtp.qq.com"
              persistent-hint
              class="mb-4"
              prepend-inner-icon="mdi-server"
            ></v-text-field>

            <v-text-field
              v-model="smtpConfigForm.port"
              label="SMTP端口"
              variant="outlined"
              type="number"
              hint="QQ邮箱默认: 465"
              persistent-hint
              class="mb-4"
              prepend-inner-icon="mdi-network-port"
            ></v-text-field>

            <v-text-field
              v-model="smtpConfigForm.username"
              label="邮箱账号"
              variant="outlined"
              hint="完整邮箱地址"
              persistent-hint
              class="mb-4"
              prepend-inner-icon="mdi-email"
            ></v-text-field>

            <v-text-field
              v-model="smtpConfigForm.password"
              label="授权码/密码"
              variant="outlined"
              type="password"
              hint="QQ邮箱请使用授权码"
              persistent-hint
              class="mb-4"
              prepend-inner-icon="mdi-lock"
            ></v-text-field>

            <v-text-field
              v-model="smtpConfigForm.from"
              label="发件人邮箱"
              variant="outlined"
              hint="显示给收件人的邮箱地址"
              persistent-hint
              class="mb-4"
              prepend-inner-icon="mdi-mail-send"
            ></v-text-field>

            <v-text-field
              v-model="smtpConfigForm.fromName"
              label="发件人名称"
              variant="outlined"
              hint="显示给收件人的名称"
              persistent-hint
              class="mb-4"
              prepend-inner-icon="mdi-user"
            ></v-text-field>

            <v-switch
              v-model="smtpConfigForm.ssl"
              label="启用SSL"
              class="mb-4"
              prepend-icon="mdi-lock-check"
            ></v-switch>

            <div class="d-flex gap-4">
              <v-btn color="primary" @click="saveSmtpConfig">
                <v-icon start>mdi-content-save</v-icon>
                保存SMTP配置
              </v-btn>
              <v-btn color="info" @click="testSmtpConfig">
                <v-icon start>mdi-mail-check</v-icon>
                测试发送
              </v-btn>
            </div>
          </v-card-text>
        </v-window-item>

        <v-window-item value="notifications">
          <v-card-text class="pa-0">
            <v-card variant="outlined" class="pa-4 mb-4">
              <v-card-title class="text-subtitle-1 pa-0 mb-4">发送系统通知</v-card-title>
              <v-form @submit.prevent="handleSendNotification">
                <v-row>
                  <v-col cols="12" md="3">
                    <v-select
                      v-model="notificationForm.type"
                      :items="notificationTypes"
                      label="通知类型"
                      variant="outlined"
                      density="compact"
                      hide-details
                    ></v-select>
                  </v-col>
                  <v-col cols="12" md="3">
                    <v-select
                      v-model="notificationForm.user_id"
                      :items="users"
                      item-title="display_name"
                      item-value="id"
                      label="目标用户"
                      variant="outlined"
                      density="compact"
                      hide-details
                    ></v-select>
                  </v-col>
                  <v-col cols="12" md="4">
                    <v-text-field
                      v-model="notificationForm.title"
                      label="通知标题"
                      variant="outlined"
                      density="compact"
                      hide-details
                    ></v-text-field>
                  </v-col>
                  <v-col cols="12" md="2">
                    <v-btn type="submit" color="primary" block>发送</v-btn>
                  </v-col>
                </v-row>
                <v-row class="mt-2">
                  <v-col cols="12">
                    <v-textarea
                      v-model="notificationForm.content"
                      label="通知内容"
                      variant="outlined"
                      rows="3"
                      hide-details
                    ></v-textarea>
                  </v-col>
                </v-row>
              </v-form>
            </v-card>

            <v-card variant="outlined" class="pa-4">
              <v-card-title class="text-subtitle-1 pa-0 mb-4">历史通知</v-card-title>
              <v-list>
                <v-list-item v-for="notif in notifications" :key="notif.id">
                  <template v-slot:prepend>
                    <v-icon :color="notif.type === 'system' ? 'primary' : notif.type === 'warning' ? 'warning' : 'success'">
                      {{ notif.type === 'system' ? 'mdi-bell' : notif.type === 'warning' ? 'mdi-alert' : 'mdi-check-circle' }}
                    </v-icon>
                  </template>
                  <v-list-item-title>{{ notif.title }}</v-list-item-title>
                  <v-list-item-subtitle>{{ notif.content }}</v-list-item-subtitle>
                  <template v-slot:append>
                    <span class="text-caption text-medium-emphasis">{{ formatDate(notif.created_at) }}</span>
                  </template>
                </v-list-item>
              </v-list>
            </v-card>
          </v-card-text>
        </v-window-item>
      </div>
    </v-main>

    <v-dialog v-model="editRoleDialog.show" max-width="400">
      <v-card>
        <v-card-title>修改用户角色</v-card-title>
        <v-card-text>
          <div class="mb-4">用户：{{ editRoleDialog.user?.display_name }}</div>
          <v-select
            v-model="editRoleDialog.role"
            :items="[
              { title: '普通用户', value: 'user' },
              { title: '系统管理员', value: 'system' },
              { title: '管理员', value: 'admin' }
            ]"
            label="选择角色"
            variant="outlined"
          ></v-select>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn variant="text" @click="editRoleDialog.show = false">取消</v-btn>
          <v-btn color="primary" @click="handleEditRole">确认</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog v-model="banDialog.show" max-width="400">
      <v-card>
        <v-card-title>封禁用户</v-card-title>
        <v-card-text>
          <div class="mb-4">用户：{{ banDialog.user?.display_name }}</div>
          <v-textarea
            v-model="banDialog.reason"
            label="封禁原因"
            variant="outlined"
            rows="3"
          ></v-textarea>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn variant="text" @click="banDialog.show = false">取消</v-btn>
          <v-btn color="error" @click="handleBan">封禁</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog v-model="statusDialog.show" max-width="400">
      <v-card>
        <v-card-title>修改文章状态</v-card-title>
        <v-card-text>
          <div class="mb-4">文章：{{ statusDialog.article?.title }}</div>
          <v-select
            v-model="statusDialog.status"
            :items="[
              { title: '待审核', value: 'pending' },
              { title: '已发布', value: 'published' },
              { title: '已拒绝', value: 'rejected' }
            ]"
            label="选择状态"
            variant="outlined"
          ></v-select>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn variant="text" @click="statusDialog.show = false">取消</v-btn>
          <v-btn color="primary" @click="handleEditStatus">确认</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog v-model="editCategoryDialog.show" max-width="400">
      <v-card>
        <v-card-title>编辑分区</v-card-title>
        <v-card-text>
          <v-text-field
            v-model="editCategoryDialog.name"
            label="分区名称"
            variant="outlined"
            class="mb-4"
          ></v-text-field>
          <v-text-field
            v-model="editCategoryDialog.description"
            label="描述"
            variant="outlined"
          ></v-text-field>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn variant="text" @click="editCategoryDialog.show = false">取消</v-btn>
          <v-btn color="primary" @click="handleEditCategory">保存</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>

  <v-card v-else-if="isInitialized && !isAdmin" class="pa-8 text-center">
    <v-icon size="64" color="error" class="mb-4">mdi-lock</v-icon>
    <div class="text-h6">您没有权限访问此页面</div>
  </v-card>
</template>

<script>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import api from '../api'
import UserAvatar from '../components/UserAvatar.vue'
import { confirm as showConfirm, success as showSuccess, error as showError } from '../utils/modal'

export default {
  name: 'Admin',
  components: {
    UserAvatar
  },
  setup() {
    const router = useRouter()
    const activeTab = ref('overview')
    const drawer = ref(false)
    const isAdmin = ref(false)
    const currentUserId = ref(null)
    const isInitialized = ref(false)

    const menuItems = computed(() => [
      { title: '概览', value: 'overview', icon: 'mdi-view-dashboard' },
      { title: '用户管理', value: 'users', icon: 'mdi-account' },
      { title: '文章管理', value: 'articles', icon: 'mdi-file-text' },
      { title: '评论管理', value: 'comments', icon: 'mdi-message-text' },
      { title: '分区管理', value: 'categories', icon: 'mdi-folder' },
      { title: '头衔管理', value: 'titles', icon: 'mdi-award' },
      { title: '侧边栏配置', value: 'sidebar', icon: 'mdi-sidebar' },
      { title: '删除审核', value: 'deletions', icon: 'mdi-delete', badge: deletionRequests.value.length },
      { title: '公告管理', value: 'announcement', icon: 'mdi-bullhorn' },
      { title: '网站配置', value: 'siteconfig', icon: 'mdi-settings' },
      { title: 'SMTP配置', value: 'smtpconfig', icon: 'mdi-email' },
      { title: '通知管理', value: 'notifications', icon: 'mdi-bell' }
    ])

    const statistics = ref({
      user_count: 0,
      article_count: 0,
      comment_count: 0,
      view_count: 0
    })

    const users = ref([])
    const articles = ref([])
    const articlePage = ref(1)
    const articleTotalPages = ref(1)
    const articleFilter = ref('')
    const articleStatusOptions = [
      { title: '全部', value: '' },
      { title: '待审核', value: 'pending' },
      { title: '已发布', value: 'published' },
      { title: '已拒绝', value: 'rejected' }
    ]

    const comments = ref([])
    const commentPage = ref(1)
    const commentTotalPages = ref(1)
    const commentTotal = ref(0)
    const allComments = computed(() => comments.value)

    const categories = ref([])
    const titles = ref([])
    const sidebarItems = ref([])
    const announcementContent = ref('')
    const siteConfigForm = ref({ siteTitle: '' })
    const smtpConfigForm = ref({
      host: '',
      port: 465,
      username: '',
      password: '',
      from: '',
      fromName: '',
      ssl: true
    })
    const notifications = ref([])
    const notificationForm = ref({
      type: 'system',
      user_id: null,
      title: '',
      content: ''
    })
    const notificationTypes = [
      { title: '系统通知', value: 'system' },
      { title: '警告', value: 'warning' },
      { title: '成功', value: 'success' }
    ]
    const deletionRequests = ref([])
    const categoryForm = ref({
      name: '',
      description: '',
      color: '#6750A4'
    })
    const titleForm = ref({
      name: '',
      description: '',
      color: '#6750A4',
      icon: ''
    })
    const grantForm = ref({
      user_id: null,
      title_id: null,
      reason: ''
    })

    const editRoleDialog = ref({
      show: false,
      user: null,
      role: 'user'
    })
    const banDialog = ref({
      show: false,
      user: null,
      reason: ''
    })
    const statusDialog = ref({
      show: false,
      article: null,
      status: 'published'
    })
    const editCategoryDialog = ref({
      show: false,
      category: null,
      name: '',
      description: ''
    })

    const checkAdmin = async () => {
      try {
        const token = localStorage.getItem('token')
        if (!token) {
          router.push('/login')
          return
        }
        const response = await api.get('/auth/me')
        if (response.data.user.role !== 'admin' && response.data.user.role !== 'system') {
          router.push('/')
          return
        }
        currentUserId.value = response.data.user.id
        isAdmin.value = true
        isInitialized.value = true
      } catch (error) {
        console.error('检查管理员权限失败', error)
        localStorage.removeItem('token')
        router.push('/login')
      }
    }

    const loadStatistics = async () => {
      try {
        const response = await api.get('/admin/statistics')
        statistics.value = response.data
      } catch (error) {
        console.error('加载统计失败', error)
      }
    }

    const loadUsers = async () => {
      try {
        const response = await api.get('/admin/users')
        users.value = response.data.users
      } catch (error) {
        console.error('加载用户失败', error)
      }
    }

    const loadArticles = async () => {
      try {
        const params = { page: articlePage.value, page_size: 20 }
        if (articleFilter.value) {
          params.status = articleFilter.value
        }
        const response = await api.get('/admin/articles', { params })
        articles.value = response.data.articles
        articleTotalPages.value = response.data.total_pages
      } catch (error) {
        console.error('加载文章失败', error)
      }
    }

    const loadComments = async () => {
      try {
        const params = { page: commentPage.value, page_size: 20 }
        const response = await api.get('/admin/comments', { params })
        comments.value = response.data.comments
        commentTotal.value = response.data.total
        commentTotalPages.value = response.data.total_pages
      } catch (error) {
        console.error('加载评论失败', error)
      }
    }

    const loadCategories = async () => {
      try {
        const response = await api.get('/categories')
        categories.value = response.data
      } catch (error) {
        console.error('加载分区失败', error)
      }
    }

    const loadTitles = async () => {
      try {
        const response = await api.get('/titles')
        titles.value = response.data
      } catch (error) {
        console.error('加载头衔失败', error)
      }
    }

    const loadSidebarConfig = async () => {
      try {
        const response = await api.get('/sidebar-config')
        sidebarItems.value = response.data.items || []
      } catch (error) {
        console.error('加载侧边栏配置失败', error)
      }
    }

    const loadAnnouncement = async () => {
      try {
        const response = await api.get('/announcement')
        announcementContent.value = response.data.content || ''
      } catch (error) {
        console.error('加载公告失败', error)
      }
    }

    const loadSiteConfig = async () => {
      try {
        const response = await api.get('/site-config')
        siteConfigForm.value.siteTitle = response.data.site_title || response.data.SiteTitle || ''
      } catch (error) {
        console.error('加载网站配置失败', error)
      }
    }

    const loadSmtpConfig = async () => {
      try {
        const response = await api.get('/site-config')
        smtpConfigForm.value = {
          host: response.data.smtp_host || response.data.SMTPHost || '',
          port: response.data.smtp_port || response.data.SMTPPort || 587,
          username: response.data.smtp_username || response.data.SMTPUsername || '',
          password: response.data.smtp_password || response.data.SMTPPassword || '',
          from: response.data.smtp_from || response.data.SMTPFrom || '',
          fromName: response.data.smtp_from_name || response.data.SMTPFromName || '',
          ssl: response.data.smtp_port == 465
        }
      } catch (error) {
        console.error('加载SMTP配置失败', error)
      }
    }

    const loadNotifications = async () => {
      try {
        const response = await api.get('/notifications/admin')
        notifications.value = response.data
      } catch (error) {
        console.error('加载通知失败', error)
      }
    }

    const loadDeletionRequests = async () => {
      try {
        const response = await api.get('/deletion-requests')
        deletionRequests.value = response.data
      } catch (error) {
        console.error('加载删除请求失败', error)
      }
    }

    const showEditRoleDialog = (user) => {
      editRoleDialog.value = {
        show: true,
        user,
        role: user.role
      }
    }

    const showBanDialog = (user) => {
      banDialog.value = {
        show: true,
        user,
        reason: ''
      }
    }

    const showStatusDialog = (article) => {
      statusDialog.value = {
        show: true,
        article,
        status: article.status
      }
    }

    const showEditCategoryDialog = (category) => {
      editCategoryDialog.value = {
        show: true,
        category,
        name: category.name,
        description: category.description || ''
      }
    }

    const handleEditRole = async () => {
      try {
        await api.put(`/admin/users/${editRoleDialog.value.user.id}/role`, { role: editRoleDialog.value.role })
        showSuccess('修改成功')
        editRoleDialog.value.show = false
        loadUsers()
      } catch (error) {
        console.error('修改角色失败', error)
        showError(error.response?.data?.error || '修改失败')
      }
    }

    const handleBan = async () => {
      if (!banDialog.value.reason) {
        showError('请输入封禁原因')
        return
      }
      try {
        await api.post(`/admin/users/${banDialog.value.user.id}/ban`, { reason: banDialog.value.reason })
        showSuccess('封禁成功')
        banDialog.value.show = false
        loadUsers()
      } catch (error) {
        console.error('封禁用户失败', error)
        showError(error.response?.data?.error || '封禁失败')
      }
    }

    const handleUnban = async (user) => {
      const confirmed = await showConfirm(`确定要解封用户 "${user.display_name}" 吗？`)
      if (!confirmed) return
      try {
        await api.post(`/admin/users/${user.id}/unban`)
        showSuccess('解封成功')
        loadUsers()
      } catch (error) {
        console.error('解封用户失败', error)
        showError(error.response?.data?.error || '解封失败')
      }
    }

    const handleDeleteUser = async (user) => {
      const confirmed = await showConfirm(`确定要删除用户 "${user.display_name}" 吗？`)
      if (!confirmed) return
      try {
        await api.delete(`/admin/users/${user.id}`)
        showSuccess('删除成功')
        loadUsers()
      } catch (error) {
        console.error('删除用户失败', error)
        showError(error.response?.data?.error || '删除失败')
      }
    }

    const handleEditStatus = async () => {
      try {
        await api.put(`/admin/articles/${statusDialog.value.article.id}/status`, { status: statusDialog.value.status })
        showSuccess('修改成功')
        statusDialog.value.show = false
        loadArticles()
      } catch (error) {
        console.error('修改状态失败', error)
        showError(error.response?.data?.error || '修改失败')
      }
    }

    const handleDeleteComment = async (commentId) => {
      const confirmed = await showConfirm('确定要删除这条评论吗？')
      if (!confirmed) return
      try {
        await api.delete(`/admin/comments/${commentId}`)
        showSuccess('删除成功')
        loadComments()
      } catch (error) {
        console.error('删除评论失败', error)
        showError(error.response?.data?.error || '删除失败')
      }
    }

    const addCategory = async () => {
      if (!categoryForm.value.name) {
        showError('请输入分区名称')
        return
      }
      try {
        await api.post('/categories', categoryForm.value)
        showSuccess('添加成功')
        categoryForm.value = { name: '', description: '', color: '#6750A4' }
        loadCategories()
      } catch (error) {
        console.error('添加分区失败', error)
        showError(error.response?.data?.error || '添加失败')
      }
    }

    const handleEditCategory = async () => {
      try {
        await api.put(`/categories/${editCategoryDialog.value.category.id}`, {
          name: editCategoryDialog.value.name,
          description: editCategoryDialog.value.description
        })
        showSuccess('保存成功')
        editCategoryDialog.value.show = false
        loadCategories()
      } catch (error) {
        console.error('保存分区失败', error)
        showError(error.response?.data?.error || '保存失败')
      }
    }

    const handleDeleteCategory = async (id) => {
      const confirmed = await showConfirm('确定要删除这个分区吗？')
      if (!confirmed) return
      try {
        await api.delete(`/categories/${id}`)
        showSuccess('删除成功')
        loadCategories()
      } catch (error) {
        console.error('删除分区失败', error)
        showError(error.response?.data?.error || '删除失败')
      }
    }

    const addTitle = async () => {
      if (!titleForm.value.name) {
        showError('请输入头衔名称')
        return
      }
      try {
        await api.post('/titles', titleForm.value)
        showSuccess('添加成功')
        titleForm.value = { name: '', description: '', color: '#6750A4', icon: '' }
        loadTitles()
      } catch (error) {
        console.error('添加头衔失败', error)
        showError(error.response?.data?.error || '添加失败')
      }
    }

    const grantTitle = async () => {
      if (!grantForm.value.user_id || !grantForm.value.title_id) {
        showError('请选择用户和头衔')
        return
      }
      try {
        await api.post('/titles/grant', grantForm.value)
        showSuccess('授予成功')
        grantForm.value = { user_id: null, title_id: null, reason: '' }
      } catch (error) {
        console.error('授予头衔失败', error)
        showError(error.response?.data?.error || '授予失败')
      }
    }

    const handleDeleteTitle = async (id) => {
      const confirmed = await showConfirm('确定要删除这个头衔吗？')
      if (!confirmed) return
      try {
        await api.delete(`/titles/${id}`)
        showSuccess('删除成功')
        loadTitles()
      } catch (error) {
        console.error('删除头衔失败', error)
        showError(error.response?.data?.error || '删除失败')
      }
    }

    const addSidebarItem = () => {
      sidebarItems.value.push({ title: '', link: '', icon: 'mdi-link' })
    }

    const removeSidebarItem = (index) => {
      sidebarItems.value.splice(index, 1)
    }

    const saveSidebarConfig = async () => {
      try {
        await api.put('/sidebar-config', { items: sidebarItems.value })
        showSuccess('保存成功')
      } catch (error) {
        console.error('保存侧边栏配置失败', error)
        showError(error.response?.data?.error || '保存失败')
      }
    }

    const saveAnnouncement = async () => {
      try {
        await api.put('/announcement', { content: announcementContent.value })
        showSuccess('保存成功')
      } catch (error) {
        console.error('保存公告失败', error)
        showError(error.response?.data?.error || '保存失败')
      }
    }

    const saveSiteConfig = async () => {
      try {
        await api.put('/site-config', { site_title: siteConfigForm.value.siteTitle })
        showSuccess('保存成功')
        if (siteConfigForm.value.siteTitle) {
          document.title = siteConfigForm.value.siteTitle
        }
      } catch (error) {
        console.error('保存网站配置失败', error)
        showError(error.response?.data?.error || '保存失败')
      }
    }

    const saveSmtpConfig = async () => {
      try {
        await api.put('/site-config', {
          smtp_host: smtpConfigForm.value.host,
          smtp_port: smtpConfigForm.value.port,
          smtp_username: smtpConfigForm.value.username,
          smtp_password: smtpConfigForm.value.password,
          smtp_from: smtpConfigForm.value.from,
          smtp_from_name: smtpConfigForm.value.fromName
        })
        showSuccess('保存成功')
      } catch (error) {
        console.error('保存SMTP配置失败', error)
        showError(error.response?.data?.error || '保存失败')
      }
    }

    const testSmtpConfig = async () => {
      try {
        await api.post('/site-config/test-smtp', {
          smtp_host: smtpConfigForm.value.host,
          smtp_port: smtpConfigForm.value.port,
          smtp_username: smtpConfigForm.value.username,
          smtp_password: smtpConfigForm.value.password,
          smtp_from: smtpConfigForm.value.from,
          smtp_to: smtpConfigForm.value.from
        })
        showSuccess('测试邮件发送成功')
      } catch (error) {
        console.error('测试SMTP配置失败', error)
        showError(error.response?.data?.error || '测试失败')
      }
    }

    const handleSendNotification = async () => {
      if (!notificationForm.value.title || !notificationForm.value.content) {
        showError('请填写通知标题和内容')
        return
      }
      try {
        await api.post('/notifications', notificationForm.value)
        showSuccess('发送成功')
        notificationForm.value = { type: 'system', user_id: null, title: '', content: '' }
        loadNotifications()
      } catch (error) {
        console.error('发送通知失败', error)
        showError(error.response?.data?.error || '发送失败')
      }
    }

    const approveDeletion = async (id) => {
      const confirmed = await showConfirm('确定要批准此删除申请吗？')
      if (!confirmed) return
      try {
        await api.post(`/admin/deletion-requests/${id}/approve`)
        showSuccess('已批准删除')
        loadDeletionRequests()
      } catch (error) {
        console.error('批准删除请求失败', error)
        showError(error.response?.data?.error || '操作失败')
      }
    }

    const rejectDeletion = async (id) => {
      const confirmed = await showConfirm('确定要拒绝此删除申请吗？')
      if (!confirmed) return
      try {
        await api.post(`/admin/deletion-requests/${id}/reject`)
        showSuccess('已拒绝')
        loadDeletionRequests()
      } catch (error) {
        console.error('拒绝删除请求失败', error)
        showError(error.response?.data?.error || '操作失败')
      }
    }

    const formatDate = (dateString) => {
      if (!dateString) return ''
      const date = new Date(dateString)
      return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
      })
    }

    const getStatusColor = (status) => {
      const colors = {
        pending: 'warning',
        published: 'success',
        rejected: 'error'
      }
      return colors[status] || 'default'
    }

    const getStatusText = (status) => {
      const texts = {
        pending: '待审核',
        published: '已发布',
        rejected: '已拒绝'
      }
      return texts[status] || status
    }

    watch(activeTab, (newTab) => {
      if (newTab === 'overview') loadStatistics()
      if (newTab === 'users') loadUsers()
      if (newTab === 'articles') loadArticles()
      if (newTab === 'comments') loadComments()
      if (newTab === 'titles') { loadTitles(); loadUsers() }
      if (newTab === 'sidebar') loadSidebarConfig()
      if (newTab === 'announcement') loadAnnouncement()
      if (newTab === 'siteconfig') loadSiteConfig()
      if (newTab === 'smtpconfig') loadSmtpConfig()
      if (newTab === 'notifications') loadNotifications()
    })

    watch(articlePage, () => {
      loadArticles()
    })

    watch(commentPage, () => {
      loadComments()
    })

    watch(articleFilter, () => {
      articlePage.value = 1
      loadArticles()
    })

    onMounted(async () => {
      await checkAdmin()
      if (isAdmin.value) {
        loadStatistics()
        loadUsers()
        loadDeletionRequests()
        loadCategories()
        loadTitles()
        loadSidebarConfig()
        loadAnnouncement()
      }
    })

    return {
      activeTab,
      drawer,
      menuItems,
      isAdmin,
      isInitialized,
      currentUserId,
      statistics,
      users,
      articles,
      articlePage,
      articleTotalPages,
      articleFilter,
      articleStatusOptions,
      allComments,
      commentPage,
      categories,
      titles,
      sidebarItems,
      announcementContent,
      siteConfigForm,
      smtpConfigForm,
      notifications,
      notificationForm,
      notificationTypes,
      deletionRequests,
      categoryForm,
      titleForm,
      grantForm,
      editRoleDialog,
      banDialog,
      statusDialog,
      editCategoryDialog,
      showEditRoleDialog,
      showBanDialog,
      showStatusDialog,
      showEditCategoryDialog,
      handleEditRole,
      handleBan,
      handleUnban,
      handleDeleteUser,
      handleEditStatus,
      handleDeleteComment,
      addCategory,
      handleEditCategory,
      handleDeleteCategory,
      addTitle,
      grantTitle,
      handleDeleteTitle,
      addSidebarItem,
      removeSidebarItem,
      saveSidebarConfig,
      saveAnnouncement,
      saveSiteConfig,
      saveSmtpConfig,
      testSmtpConfig,
      handleSendNotification,
      approveDeletion,
      rejectDeletion,
      formatDate,
      getStatusColor,
      getStatusText
    }
  }
}
</script>

<style scoped>
.admin-layout {
  min-height: 100vh;
}

.admin-sidebar {
  width: 250px;
}

.admin-app-bar {
  background: rgb(var(--v-theme-surface));
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.admin-main {
  padding-top: 64px;
}

.admin-content {
  padding: 24px;
  max-width: 1400px;
  margin: 0 auto;
}

@media (max-width: 600px) {
  .admin-sidebar {
    width: 240px;
  }

  .admin-content {
    padding: 16px;
  }
}
</style>