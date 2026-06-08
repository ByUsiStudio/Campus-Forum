package miao.byusi.android.xylt

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.outlined.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.compose.ui.graphics.Color
import androidx.navigation.NavHostController
import kotlinx.coroutines.delay
import coil.compose.AsyncImage

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
                            imageVector = Icons.Outlined.ArrowBack,
                            contentDescription = "返回"
                        )
                    }
                    Text(
                        text = "文章详情",
                        fontSize = 20.sp
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
                    Text(text = "加载失败: $error")
                    Spacer(modifier = Modifier.height(16.dp))
                    Button(onClick = refreshData) {
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
                            fontSize = 16.sp
                        )
                        Text(
                            text = article?.createdAt ?: "",
                            fontSize = 12.sp,
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
                    fontSize = 24.sp,
                    modifier = Modifier.padding(horizontal = 16.dp)
                )
                
                Spacer(modifier = Modifier.height(16.dp))
                
                // 文章内容：使用 WebView 加载后端文章页，完整保留 MD / HTML 渲染
                ArticleContentView(
                    articleId = articleId,
                    modifier = Modifier
                        .fillMaxWidth()
                        .heightIn(min = 100.dp)
                        .padding(horizontal = 16.dp)
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
                    Button(
                        onClick = {
                            toggleLike(articleId, isLiked) { newLiked ->
                                isLiked = newLiked
                            }
                        }
                    ) {
                        Icon(
                            imageVector = if (isLiked) Icons.Outlined.Favorite else Icons.Outlined.FavoriteBorder,
                            contentDescription = null
                        )
                        Spacer(modifier = Modifier.width(4.dp))
                        Text(if (isLiked) "已点赞" else "点赞")
                    }
                    Button(onClick = { 
                        replyToComment = null
                        commentText = ""
                        showCommentDialog = true 
                    }) {
                        Icon(
                            imageVector = Icons.Outlined.Comment,
                            contentDescription = null
                        )
                        Spacer(modifier = Modifier.width(4.dp))
                        Text("评论")
                    }
                    Button(onClick = { 
                        toggleFavorite(articleId, isFavorited) { newFavorited ->
                            isFavorited = newFavorited
                        }
                    }) {
                        Icon(
                            imageVector = if (isFavorited) Icons.Outlined.Bookmark else Icons.Outlined.BookmarkBorder,
                            contentDescription = null
                        )
                        Spacer(modifier = Modifier.width(4.dp))
                        Text(if (isFavorited) "已收藏" else "收藏")
                    }
                    Button(onClick = { shareArticle(articleId) }) {
                        Icon(
                            imageVector = Icons.Outlined.Share,
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
                            imageVector = Icons.Outlined.Report,
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
                    fontSize = 18.sp,
                    modifier = Modifier.padding(horizontal = 16.dp)
                )
                
                Spacer(modifier = Modifier.height(8.dp))
                
                if (comments.isEmpty()) {
                    Text(
                        text = "暂无评论",
                        fontSize = 14.sp,
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
        AlertDialog(
            onDismissRequest = { showCommentDialog = false },
            title = { Text(if (replyToComment != null) "回复评论" else "发表评论") },
            text = {
                Column(modifier = Modifier.padding(vertical = 8.dp)) {
                    if (replyToComment != null) {
                        Text(
                            text = "回复 @${replyToComment?.authorName}",
                            fontSize = 14.sp,
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
                }
            },
            dismissButton = {
                TextButton(onClick = { showCommentDialog = false }) {
                    Text("取消")
                }
            },
            confirmButton = {
                Button(
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
        )
    }
    
    // 举报对话框
    if (showReportDialog) {
        var reportReason by remember { mutableStateOf("") }
        var reportDescription by remember { mutableStateOf("") }
        
        AlertDialog(
            onDismissRequest = { showReportDialog = false },
            title = { Text("举报") },
            text = {
                Column(modifier = Modifier.padding(vertical = 8.dp)) {
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
                }
            },
            dismissButton = {
                TextButton(onClick = { showReportDialog = false }) {
                    Text("取消")
                }
            },
            confirmButton = {
                Button(
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
        )
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
                    text = comment.authorName ?: "匿名",
                    fontSize = 14.sp,
                    color = Color(0xFF007AFF)
                )
            }
            Spacer(modifier = Modifier.height(4.dp))
            Text(
                text = comment.content,
                fontSize = 14.sp
            )
            Spacer(modifier = Modifier.height(4.dp))
            Row(
                horizontalArrangement = Arrangement.spacedBy(16.dp)
            ) {
                Text(
                    text = comment.createdAt,
                    fontSize = 12.sp,
                    color = Color(0xFF666666)
                )
                Button(onClick = { onLike(comment.id) }) {
                    Icon(
                        imageVector = Icons.Outlined.FavoriteBorder,
                        contentDescription = null,
                        modifier = Modifier.size(16.dp)
                    )
                    Spacer(modifier = Modifier.width(4.dp))
                    Text("${comment.likeCount}", fontSize = 12.sp)
                }
                Button(onClick = { onReply(comment) }) {
                    Icon(
                        imageVector = Icons.Outlined.Reply,
                        contentDescription = null,
                        modifier = Modifier.size(16.dp)
                    )
                    Spacer(modifier = Modifier.width(4.dp))
                    Text("回复", fontSize = 12.sp)
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
                                    fontSize = 12.sp,
                                    color = Color(0xFF007AFF)
                                )
                                Spacer(modifier = Modifier.height(2.dp))
                                Text(
                                    text = reply.content,
                                    fontSize = 12.sp
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
    ApiClient.getArticleDetail(articleId, 1, 20, object : ApiCallback {
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
        displayName = json.optString("display_name", json.optString("nickname", null)),
        avatar = json.optString("avatar", json.optString("avatar_url", null)),
        role = json.optString("role", null),
        signature = json.optString("signature", null)
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
    return Comment(
        id = json.optInt("id", 0),
        content = json.optString("content", ""),
        articleId = json.optInt("article_id", 0),
        parentId = json.optInt("parent_id", 0).takeIf { it > 0 },
        isAnonymous = json.optBoolean("is_anonymous", false),
        likeCount = json.optInt("like_count", 0),
        createdAt = json.optString("created_at", ""),
        user = parseUser(json.optJSONObject("user")),
        authorName = json.optJSONObject("user")?.optString("display_name")
            ?: json.optJSONObject("user")?.optString("nickname")
            ?: "匿名",
        replies = json.optJSONArray("replies")?.let { array ->
            val list = mutableListOf<Comment>()
            for (i in 0 until array.length()) {
                list.add(parseComment(array.getJSONObject(i)))
            }
            list
        }
    )
}
