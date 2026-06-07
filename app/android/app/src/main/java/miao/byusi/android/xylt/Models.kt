package miao.byusi.android.xylt

data class Article(
    val id: Int,
    val title: String,
    val content: String,
    val authorName: String?,
    val viewCount: Int,
    val likeCount: Int,
    val createdAt: String
)