package miao.byusi.android.xylt

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.compose.ui.graphics.Color
import androidx.navigation.NavHostController
import kotlinx.coroutines.delay
import coil.compose.AsyncImage
import top.yukonga.miuix.kmp.basic.Scaffold
import top.yukonga.miuix.kmp.basic.Surface
import top.yukonga.miuix.kmp.basic.Text
import top.yukonga.miuix.kmp.basic.Button
import top.yukonga.miuix.kmp.basic.FilledButton
import top.yukonga.miuix.kmp.basic.OutlinedButton
import top.yukonga.miuix.kmp.basic.CircularProgressIndicator
import top.yukonga.miuix.kmp.basic.Icon
import top.yukonga.miuix.kmp.basic.TextField
import top.yukonga.miuix.kmp.basic.Dialog
import top.yukonga.miuix.kmp.icon.MiuixIcons
import top.yukonga.miuix.kmp.icon.icons.outlined.ArrowBack
import top.yukonga.miuix.kmp.icon.icons.outlined.Favorite
import top.yukonga.miuix.kmp.icon.icons.outlined.FavoriteBorder
import top.yukonga.miuix.kmp.icon.icons.outlined.Comment
import top.yukonga.miuix.kmp.icon.icons.outlined.Report
import top.yukonga.miuix.kmp.icon.icons.outlined.Bookmark
import top.yukonga.miuix.kmp.icon.icons.outlined.BookmarkBorder
import top.yukonga.miuix.kmp.icon.icons.outlined.Share
import top.yukonga.miuix.kmp.icon.icons.outlined.Reply

@Composable
fun ArticleDetailScreen(navController: NavHostController, articleId: Int) {
    var article by remember { mutableStateOf<ArticleDetail?>(null) }
    var comments by remember { mutableStateOf<List<Comment>>(emptyList()) }
    var isLoading by remember { mutableStateOf(true) }
    var error by remember { mutableStateOf<String?>(null) }
    var isLiked by remember { mutableStateOf(false) }
    var isFavorited by remember { mutableStateOf(false) }
    var showCommentDialog by remember { mutableStateOf(false) }
    var showReportDialog by remember { mutableStateOf(false) }
    var commentText by remember { mutableStateOf("") }
    var replyToComment by remember { mutableStateOf<Comment?>(null) }
    var refreshing by remember { mutableStateOf(false) }
    
    // 自动刷新：每30秒刷新一次
    LaunchedEffect(articleId) {
        while (true) {
            loadArticleDetail(articleId) { result ->
                result.onSuccess { (a, c, liked) ->
                    article = a
                    comments = c
                    isLiked = liked
                    isLoading = false
                    error = null
                }.onFailure { e ->
                    error = e.message
                    isLoading = false
                }
            }
            checkFavoriteStatus(articleId) { favorited ->
                isFavorited = favorited
            }
            delay(30000)
        }
    }
    
    val refreshData = {
        refreshing = true
        loadArticleDetail(articleId) { result ->
            result.onSuccess { (a, c, liked) ->
                article = a
                comments = c
                isLiked = liked
                error = null
            }.onFailure { e ->
                error = e.message
            }
            checkFavoriteStatus(articleId) { favorited ->
                isFavorited = favorited
            }
            refreshing = false
        }
    }
    
    Scaffold(
        topBar = {
            Surface(
                modifier = Modifier.fillMaxWidth()
            ) {
                Row(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(16.dp),
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Button(onClick = { navController.popBackStack() }) {
                        Icon(
                            imageVector = MiuixIcons.Outlined.ArrowBack,
                            contentDescription = "返回"
                        )
                    }
                    Text(
                        text = "文章详情",
                        fontSize = 20.dp
                    )
                }
            }
        }
    ) { padding ->
        if (isLoading) {
            Box(
                modifier = Modifier
                    .fillMaxSize()
                    .padding(padding),
                contentAlignment = Alignment.Center
            ) {
                CircularProgressIndicator()
            }
        } else if (error != null) {
            Box(
                modifier = Modifier
                    .fillMaxSize()
                    .padding(padding),
                contentAlignment = Alignment.Center
            ) {
                Column(horizontalAlignment = Alignment.CenterHorizontally) {
                    Text("加载失败: $error")
                    Spacer(modifier = Modifier.height(16.dp))
                    FilledButton(onClick = refreshData) {
                        Text("重试")
                    }
                }
            }
        } else if (article != null) {
            Column(
                modifier = Modifier
                    .fillMaxSize()
                    .padding(padding)
                    .verticalScroll(rememberScrollState())
            ) {
                // 作者信息
                Row(
                    modifier = Modifier.padding(16.dp),
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    AsyncImage(
                        model = article?.authorAvatar,
                        contentDescription = "头像",
                        modifier = Modifier.size(48.dp)
                    )
                    Spacer(modifier = Modifier.width(12.dp))
                    Column(modifier = Modifier.weight(1f)) {
                        Text(
                            text = article?.authorName ?: "匿名",
                            fontSize = 16.dp
                        )
                        Text(
                            text = article?.createdAt ?: "",
                            fontSize = 12.dp,
                            color = Color(0xFF666666)
                        )
                    }
                    article?.user?.id?.let { userId ->
                        OutlinedButton(onClick = { toggleFollow(userId) }) {
                            Text("关注")
                        }
                    }
                }
                
                // 文章标题
                Text(
                    text = article?.title ?: "",
                    fontSize = 24.dp,
                    modifier = Modifier.padding(horizontal = 16.dp)
                )
                
                Spacer(modifier = Modifier.height(16.dp))
                
                // 文章内容
                Text(
                    text = article?.content ?: "",
                    fontSize = 16.dp,
                    modifier = Modifier.padding(horizontal = 16.dp)
                )
                
                Spacer(modifier = Modifier.height(16.dp))
                
                // 统计信息
                Row(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(horizontal = 16.dp),
                    horizontalArrangement = Arrangement.spacedBy(24.dp)
                ) {
                    Text("${article?.viewCount ?: 0} 浏览")
                    Text("${article?.likeCount ?: 0} 点赞")
                    Text("${article?.commentCount ?: 0} 评论")
                    Text("${article?.shareCount ?: 0} 分享")
                }
                
                Spacer(modifier = Modifier.height(8.dp))
                
                // 操作按钮
                Row(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(horizontal = 8.dp),
                    horizontalArrangement = Arrangement.spacedBy(8.dp)
                ) {
                    FilledButton(
                        onClick = {
                            toggleLike(articleId, isLiked) { newLiked ->
                                isLiked = newLiked
                            }
                        }
                    ) {
                        Icon(
                            imageVector = if (isLiked) MiuixIcons.Outlined.Favorite else MiuixIcons.Outlined.FavoriteBorder,
                            contentDescription = null
                        )
                        Spacer(modifier = Modifier.width(4.dp))
                        Text(if (isLiked) "已点赞" else "点赞")
                    }
                    FilledButton(onClick = { 
                        replyToComment = null
                        commentText = ""
                        showCommentDialog = true 
                    }) {
                        Icon(
                            imageVector = MiuixIcons.Outlined.Comment,
                            contentDescription = null
                        )
                        Spacer(modifier = Modifier.width(4.dp))
                        Text("评论")
                    }
                    FilledButton(onClick = { 
                        toggleFavorite(articleId, isFavorited) { newFavorited ->
                            isFavorited = newFavorited
                        }
                    }) {
                        Icon(
                            imageVector = if (isFavorited) MiuixIcons.Outlined.Bookmark else MiuixIcons.Outlined.BookmarkBorder,
                            contentDescription = null
                        )
                        Spacer(modifier = Modifier.width(4.dp))
                        Text(if (isFavorited) "已收藏" else "收藏")
                    }
                    FilledButton(onClick = { shareArticle(articleId) }) {
                        Icon(
                            imageVector = MiuixIcons.Outlined.Share,
                            contentDescription = null
                        )
                        Spacer(modifier = Modifier.width(4.dp))
                        Text("分享")
                    }
                }
                
                Spacer(modifier = Modifier.height(8.dp))
                
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.Center
                ) {
                    OutlinedButton(onClick = { showReportDialog = true }) {
                        Icon(
                            imageVector = MiuixIcons.Outlined.Report,
                            contentDescription = null
                        )
                        Spacer(modifier = Modifier.width(4.dp))
                        Text("举报")
                    }
                }
                
                Spacer(modifier = Modifier.height(24.dp))
                
                // 评论列表
                Text(
                    text = "评论 (${comments.size})",
                    fontSize = 18.dp,
                    modifier = Modifier.padding(horizontal = 16.dp)
                )
                
                Spacer(modifier = Modifier.height(8.dp))
                
                if (comments.isEmpty()) {
                    Text(
                        text = "暂无评论",
                        fontSize = 14.dp,
                        color = Color(0xFF666666),
                        modifier = Modifier.padding(horizontal = 16.dp)
                    )
                } else {
                    Column(
                        modifier = Modifier.padding(horizontal = 16.dp),
                        verticalArrangement = Arrangement.spacedBy(12.dp)
                    ) {
                        comments.forEach { comment ->
                            CommentItem(
                                comment = comment,
                                onLike = { likeComment(it) },
                                onReply = { 
                                    replyToComment = it
                                    commentText = ""
                                    showCommentDialog = true
                                }
                            )
                        }
                    }
                }
            }
        }
    }
    
    // 评论对话框
    if (showCommentDialog) {
        Dialog(
            onDismissRequest = { showCommentDialog = false },
            title = { Text(if (replyToComment != null) "回复评论" else "发表评论") }
        ) {
            Column(modifier = Modifier.padding(16.dp)) {
                if (replyToComment != null) {
                    Text(
                        text = "回复 @${replyToComment?.authorName}",
                        fontSize = 14.dp,
                        color = Color(0xFF666666)
                    )
                    Spacer(modifier = Modifier.height(8.dp))
                }
                TextField(
                    value = commentText,
                    onValueChange = { commentText = it },
                    placeholder = { Text("说点什么...") },
                    modifier = Modifier.fillMaxWidth(),
                    minLines = 3,
                    maxLines = 5
                )
                Spacer(modifier = Modifier.height(16.dp))
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.End
                ) {
                    OutlinedButton(onClick = { showCommentDialog = false }) {
                        Text("取消")
                    }
                    Spacer(modifier = Modifier.width(8.dp))
                    FilledButton(
                        onClick = {
                            submitComment(
                                articleId,
                                commentText,
                                replyToComment?.id
                            ) {
                                showCommentDialog = false
                                replyToComment = null
                                commentText = ""
                                refreshData()
                            }
                        },
                        enabled = commentText.isNotBlank()
                    ) {
                        Text("发送")
                    }
                }
            }
        }
    }
    
    // 举报对话框
    if (showReportDialog) {
        var reportReason by remember { mutableStateOf("") }
        var reportDescription by remember { mutableStateOf("") }
        
        Dialog(
            onDismissRequest = { showReportDialog = false },
            title = { Text("举报") }
        ) {
            Column(modifier = Modifier.padding(16.dp)) {
                TextField(
                    value = reportReason,
                    onValueChange = { reportReason = it },
                    placeholder = { Text("举报原因") },
                    modifier = Modifier.fillMaxWidth()
                )
                Spacer(modifier = Modifier.height(8.dp))
                TextField(
                    value = reportDescription,
                    onValueChange = { reportDescription = it },
                    placeholder = { Text("详细描述") },
                    modifier = Modifier.fillMaxWidth(),
                    minLines = 3,
                    maxLines = 5
                )
                Spacer(modifier = Modifier.height(16.dp))
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.End
                ) {
                    OutlinedButton(onClick = { showReportDialog = false }) {
                        Text("取消")
                    }
                    Spacer(modifier = Modifier.width(8.dp))
                    FilledButton(
                        onClick = {
                            submitReport("article", articleId, reportReason, reportDescription) {
                                showReportDialog = false
                            }
                        },
                        enabled = reportReason.isNotBlank()
                    ) {
                        Text("提交")
                    }
                }
            }
        }
    }
}

@Composable
fun CommentItem(
    comment: Comment,
    onLike: (Int) -> Unit,
    onReply: (Comment) -> Unit
) {
    Surface(
        modifier = Modifier.fillMaxWidth()
    ) {
        Column(modifier = Modifier.padding(12.dp)) {
            Row(verticalAlignment = Alignment.CenterVertically) {
                AsyncImage(
                    model = comment.user?.avatar,
                    contentDescription = "头像",
                    modifier = Modifier.size(32.dp)
                )
                Spacer(modifier = Modifier.width(8.dp))
                Text(
                    text = comment.authorName,
                    fontSize = 14.dp,
                    color = Color(0xFF007AFF)
                )
            }
            Spacer(modifier = Modifier.height(4.dp))
            Text(
                text = comment.content,
                fontSize = 14.dp
            )
            Spacer(modifier = Modifier.height(4.dp))
            Row(
                horizontalArrangement = Arrangement.spacedBy(16.dp)
            ) {
                Text(
                    text = comment.createdAt,
                    fontSize = 12.dp,
                    color = Color(0xFF666666)
                )
                Button(onClick = { onLike(comment.id) }) {
                    Icon(
                        imageVector = MiuixIcons.Outlined.FavoriteBorder,
                        contentDescription = null,
                        modifier = Modifier.size(16.dp)
                    )
                    Spacer(modifier = Modifier.width(4.dp))
                    Text("${comment.likeCount}", fontSize = 12.dp)
                }
                Button(onClick = { onReply(comment) }) {
                    Icon(
                        imageVector = MiuixIcons.Outlined.Reply,
                        contentDescription = null,
                        modifier = Modifier.size(16.dp)
                    )
                    Spacer(modifier = Modifier.width(4.dp))
                    Text("回复", fontSize = 12.dp)
                }
            }
            // 回复列表
            if (!comment.replies.isNullOrEmpty()) {
                Spacer(modifier = Modifier.height(8.dp))
                Column(
                    modifier = Modifier.padding(start = 16.dp),
                    verticalArrangement = Arrangement.spacedBy(8.dp)
                ) {
                    comment.replies?.forEach { reply ->
                        Surface {
                            Column(modifier = Modifier.padding(8.dp)) {
                                Text(
                                    text = "@${reply.authorName}",
                                    fontSize = 12.dp,
                                    color = Color(0xFF007AFF)
                                )
                                Spacer(modifier = Modifier.height(2.dp))
                                Text(
                                    text = reply.content,
                                    fontSize = 12.dp
                                )
                            }
                        }
                    }
                }
            }
        }
    }
}

private fun toggleLike(articleId: Int, isCurrentlyLiked: Boolean, callback: (Boolean) -> Unit) {
    val apiCallback = object : ApiCallback {
        override fun onSuccess(response: String) {
            callback(!isCurrentlyLiked)
        }
        override fun onError(error: String) {}
    }
    if (isCurrentlyLiked) {
        ApiClient.unlikeArticle(articleId, apiCallback)
    } else {
        ApiClient.likeArticle(articleId, apiCallback)
    }
}

private fun toggleFavorite(articleId: Int, isCurrentlyFavorited: Boolean, callback: (Boolean) -> Unit) {
    val apiCallback = object : ApiCallback {
        override fun onSuccess(response: String) {
            callback(!isCurrentlyFavorited)
        }
        override fun onError(error: String) {}
    }
    if (isCurrentlyFavorited) {
        ApiClient.unfavoriteArticle(articleId, apiCallback)
    } else {
        ApiClient.favoriteArticle(articleId, apiCallback)
    }
}

private fun checkFavoriteStatus(articleId: Int, callback: (Boolean) -> Unit) {
    ApiClient.checkFavoriteStatus(articleId, object : ApiCallback {
        override fun onSuccess(response: String) {
            try {
                val json = org.json.JSONObject(response)
                callback(json.optBoolean("favorited", false))
            } catch (e: Exception) {
                callback(false)
            }
        }
        override fun onError(error: String) {
            callback(false)
        }
    })
}

private fun shareArticle(articleId: Int) {
    ApiClient.shareArticle(articleId, object : ApiCallback {
        override fun onSuccess(response: String) {}
        override fun onError(error: String) {}
    })
}

private fun toggleFollow(userId: Int) {
    ApiClient.checkFollowStatus(userId, object : ApiCallback {
        override fun onSuccess(response: String) {
            try {
                val json = org.json.JSONObject(response)
                val isFollowing = json.optBoolean("is_following", false)
                val callback = object : ApiCallback {
                    override fun onSuccess(response: String) {}
                    override fun onError(error: String) {}
                }
                if (isFollowing) {
                    ApiClient.unfollowUser(userId, callback)
                } else {
                    ApiClient.followUser(userId, callback)
                }
            } catch (e: Exception) {}
        }
        override fun onError(error: String) {}
    })
}

private fun likeComment(commentId: Int) {
    ApiClient.likeComment(commentId, object : ApiCallback {
        override fun onSuccess(response: String) {}
        override fun onError(error: String) {}
    })
}

private fun submitComment(articleId: Int, content: String, parentId: Int?, onSuccess: () -> Unit) {
    ApiClient.createComment(articleId, content, parentId, false, object : ApiCallback {
        override fun onSuccess(response: String) {
            onSuccess()
        }
        override fun onError(error: String) {}
    })
}

private fun submitReport(targetType: String, targetId: Int, reason: String, description: String, onSuccess: () -> Unit) {
    ApiClient.submitReport(targetType, targetId, reason, description, object : ApiCallback {
        override fun onSuccess(response: String) {
            onSuccess()
        }
        override fun onError(error: String) {}
    })
}

private fun loadArticleDetail(articleId: Int, callback: (Result<Triple<ArticleDetail, List<Comment>, Boolean>>) -> Unit) {
    ApiClient.getArticleDetail(articleId, object : ApiCallback {
        override fun onSuccess(response: String) {
            try {
                val json = org.json.JSONObject(response)
                val articleJson = json.getJSONObject("article")
                val commentsArray = json.optJSONArray("comments") ?: org.json.JSONArray()
                val liked = json.optBoolean("liked", false)

                val article = ArticleDetail(
                    id = articleJson.optInt("id", 0),
                    title = articleJson.optString("title", ""),
                    content = articleJson.optString("content", ""),
                    contentHtml = articleJson.optString("content_html"),
                    authorName = articleJson.optJSONObject("user")?.optString("display_name")
                        ?: articleJson.optJSONObject("user")?.optString("nickname")
                        ?: "匿名",
                    authorAvatar = articleJson.optJSONObject("user")?.optString("avatar"),
                    user = parseUser(articleJson.optJSONObject("user")),
                    category = parseCategory(articleJson.optJSONObject("category")),
                    viewCount = articleJson.optInt("view_count", 0),
                    likeCount = articleJson.optInt("like_count", 0),
                    commentCount = articleJson.optInt("comment_count", 0),
                    shareCount = articleJson.optInt("share_count", 0),
                    voiceUrl = articleJson.optString("voice_url"),
                    isAnonymous = articleJson.optBoolean("is_anonymous", false),
                    createdAt = articleJson.optString("created_at", "")
                )

                val commentsList = mutableListOf<Comment>()
                for (i in 0 until commentsArray.length()) {
                    commentsList.add(parseComment(commentsArray.getJSONObject(i)))
                }

                callback(Result.success(Triple(article, commentsList, liked)))
            } catch (e: Exception) {
                callback(Result.failure(e))
            }
        }
        override fun onError(error: String) {
            callback(Result.failure(Exception(error)))
        }
    })
}

private fun parseUser(json: org.json.JSONObject?): User? {
    if (json == null) return null
    return User(
        id = json.optInt("id", 0),
        username = json.optString("username", ""),
        displayName = json.optString("display_name"),
        avatar = json.optString("avatar"),
        role = json.optString("role"),
        signature = json.optString("signature")
    )
}

private fun parseCategory(json: org.json.JSONObject?): Category? {
    if (json == null) return null
    return Category(
        id = json.optInt("id", 0),
        name = json.optString("name", "")
    )
}

private fun parseComment(json: org.json.JSONObject): Comment {
    val repliesArray = json.optJSONArray("replies")
    val replies = if (repliesArray != null) {
        mutableListOf<Comment>().apply {
            for (i in 0 until repliesArray.length()) {
                add(parseComment(repliesArray.getJSONObject(i)))
            }
        }
    } else null
    
    return Comment(
        id = json.optInt("id", 0),
        content = json.optString("content", ""),
        user = parseUser(json.optJSONObject("user")),
        authorName = json.optJSONObject("user")?.optString("display_name")
            ?: json.optJSONObject("user")?.optString("nickname")
            ?: "匿名",
        articleId = json.optInt("article_id", 0),
        parentId = json.optInt("parent_id").takeIf { it != 0 },
        likeCount = json.optInt("like_count", 0),
        replyCount = json.optInt("reply_count", 0),
        replies = replies,
        isAnonymous = json.optBoolean("is_anonymous", false),
        createdAt = json.optString("created_at", "")
    )
}
// 数据类统一在 Models.kt，避免分散声明与重名冲突。