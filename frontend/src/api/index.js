/**
 * 统一 API 入口
 * --------------------------------------------------------------------------
 * - 默认导出 `http`（axios 实例）：`import api from '@/api'` 仍可直接 api.get/post...
 * - 按业务域具名导出：articleApi / userApi / ...
 */
import http from './http'

export default http

export { authApi } from './auth'
export { articleApi, reportApi } from './article'
export { commentApi } from './comment'
export { userApi, signinApi } from './user'
export { friendApi } from './friend'
export { notificationApi, userNotificationApi, adminUserNotificationApi } from './notification'
export { commonApi } from './common'
export { uploadApi } from './upload'
export { default as adminApi, adminUserApi, adminArticleApi, adminCommentApi, adminStatsApi, adminSidebarApi, adminDeletionApi, adminAnnouncementApi, adminSiteConfigApi, adminNotificationApi, adminTitleApi, adminCategoryApi, adminAuthApi, permissionGroupApi } from './admin'
export { levelApi } from './level'
export { statisticsApi } from './statistics'
export { collectionApi } from './collection'
export { topicApi } from './topic'
export { leaderboardApi } from './leaderboard'
