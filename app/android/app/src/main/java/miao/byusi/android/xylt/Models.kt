package miao.byusi.android.xylt

/**
 * 全局数据模型集中地。
 *
 * 之前 ArticleDetailScreen / DraftsScreen / ProfileScreen 各自在文件底部
 * 声明了 [ArticleDetail] / [Comment] / [Draft] / [UserProfile]，分散且容易
 * 重复定义。现统一收敛到本文件，避免命名冲突并提升可维护性。
 */

/** 用户（基础信息） */
data class User(
    val id: Int,
    val username: String,
    val displayName: String?,
    val avatar: String?,
    val role: String? = null,
    val signature: String? = null
)

/** 文章（列表项） */
data class Article(
    val id: Int,
    val title: String,
    val content: String,
    val contentHtml: String? = null,
    val user: User? = null,
    val category: Category? = null,
    val viewCount: Int = 0,
    val likeCount: Int = 0,
    val commentCount: Int = 0,
    val shareCount: Int = 0,
    val voiceUrl: String? = null,
    val isAnonymous: Boolean = false,
    val createdAt: String,
    val status: String? = null
)

/** 文章详情（包含作者头像） */
data class ArticleDetail(
    val id: Int,
    val title: String,
    val content: String,
    val contentHtml: String? = null,
    val authorName: String?,
    val authorAvatar: String?,
    val user: User? = null,
    val category: Category? = null,
    val viewCount: Int,
    val likeCount: Int,
    val commentCount: Int = 0,
    val shareCount: Int = 0,
    val voiceUrl: String? = null,
    val isAnonymous: Boolean = false,
    val createdAt: String
)

/** 分区 */
data class Category(
    val id: Int,
    val name: String
)

/** 评论 */
data class Comment(
    val id: Int,
    val content: String,
    val user: User? = null,
    val authorName: String?,
    val articleId: Int,
    val parentId: Int? = null,
    val likeCount: Int = 0,
    val replyCount: Int = 0,
    val replies: List<Comment>? = null,
    val isAnonymous: Boolean = false,
    val createdAt: String
)

/** 草稿 */
data class Draft(
    val id: Int,
    val title: String,
    val content: String,
    val createdAt: String,
    val status: String? = null
)

/** 用户资料 */
data class UserProfile(
    val id: Int,
    val username: String,
    val displayName: String?,
    val avatar: String?,
    val signature: String?,
    val role: String? = null,
    val articleCount: Int = 0,
    val followerCount: Int = 0,
    val followingCount: Int = 0,
    val createdAt: String? = null
)

/** 通知 */
data class Notification(
    val id: Int,
    val userId: Int,
    val senderId: Int? = null,
    val articleId: Int? = null,
    val type: String,
    val isRead: Boolean = false,
    val sender: User? = null,
    val article: Article? = null,
    val createdAt: String
)

/** 文章列表响应 */
data class ArticlesResponse(
    val articles: List<Article>,
    val total: Int,
    val page: Int,
    val totalPages: Int
)

/** 文章详情响应 */
data class ArticleDetailResponse(
    val article: ArticleDetail,
    val comments: List<Comment>,
    val total: Int,
    val page: Int,
    val pageSize: Int,
    val totalPages: Int,
    val liked: Boolean = false,
    val commentLiked: Map<String, Boolean> = emptyMap()
)

/** 登录响应 */
data class LoginResponse(
    val token: String,
    val user: User
)

/** 简单消息响应 */
data class MessageResponse(
    val message: String
)

/** 关注状态 */
data class FollowStatus(
    val isFollowing: Boolean,
    val isFollowed: Boolean,
    val mutual: Boolean,
    val followingUser: User? = null
)

/** 分区列表响应 */
data class CategoriesResponse(
    val categories: List<Category>
)

/** 关注列表响应 */
data class FollowingResponse(
    val following: List<User>
)

/** 粉丝列表响应 */
data class FollowersResponse(
    val followers: List<User>
)

/** 通知列表响应 */
data class NotificationsResponse(
    val notifications: List<Notification>,
    val total: Int,
    val page: Int,
    val pageSize: Int,
    val totalPages: Int
)

/** 未读通知数量响应 */
data class UnreadCountResponse(
    val count: Int
)

/** 创建文章响应 */
data class CreateArticleResponse(
    val message: String,
    val article: Article? = null
)

/** 分享文章响应 */
data class ShareArticleResponse(
    val message: String,
    val shareCount: Int
)
