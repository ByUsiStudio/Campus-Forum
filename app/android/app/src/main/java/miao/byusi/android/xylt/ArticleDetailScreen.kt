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
import miuix.compose.scaffold.Scaffold
import miuix.compose.surface.Surface
import miuix.compose.text.Text
import miuix.compose.button.Button
import miuix.compose.button.FilledButton
import miuix.compose.progress.CircularProgressIndicator
import miuix.compose.icon.Icon
import miuix.compose.icon.icons.MiuixIcons
import miuix.compose.icon.icons.outlined.ArrowBack
import miuix.compose.icon.icons.outlined.Favorite
import miuix.compose.icon.icons.outlined.Comment
import miuix.compose.icon.icons.outlined.Report

@Composable
fun ArticleDetailScreen(navController: NavHostController, articleId: Int) {
    var article by remember { mutableStateOf<ArticleDetail?>(null) }
    var comments by remember { mutableStateOf<List<Comment>>(emptyList()) }
    var isLoading by remember { mutableStateOf(true) }
    var error by remember { mutableStateOf<String?>(null) }
    
    // 自动刷新：每30秒刷新一次
    LaunchedEffect(articleId) {
        while (true) {
            loadArticleDetail(articleId) { result ->
                result.onSuccess { (a, c) ->
                    article = a
                    comments = c
                    isLoading = false
                    error = null
                }.onFailure { e ->
                    error = e.message
                    isLoading = false
                }
            }
            delay(30000)
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
                Text("加载失败: $error")
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
                    Column {
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
                    Text("${comments.size} 评论")
                }
                
                Spacer(modifier = Modifier.height(8.dp))
                
                // 操作按钮
                Row(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(horizontal = 8.dp),
                    horizontalArrangement = Arrangement.spacedBy(8.dp)
                ) {
                    FilledButton(onClick = { likeArticle(articleId) }) {
                        Icon(
                            imageVector = MiuixIcons.Outlined.Favorite,
                            contentDescription = null
                        )
                        Spacer(modifier = Modifier.width(4.dp))
                        Text("点赞")
                    }
                    FilledButton(onClick = { /* 显示评论对话框 */ }) {
                        Icon(
                            imageVector = MiuixIcons.Outlined.Comment,
                            contentDescription = null
                        )
                        Spacer(modifier = Modifier.width(4.dp))
                        Text("评论")
                    }
                    FilledButton(onClick = { /* 显示举报对话框 */ }) {
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
                    text = "评论",
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
                            CommentItem(comment = comment)
                        }
                    }
                }
            }
        }
    }
}

@Composable
fun CommentItem(comment: Comment) {
    Surface(
        modifier = Modifier.fillMaxWidth()
    ) {
        Column(modifier = Modifier.padding(12.dp)) {
            Text(
                text = comment.authorName,
                fontSize = 14.dp,
                color = Color(0xFF007AFF)
            )
            Spacer(modifier = Modifier.height(4.dp))
            Text(
                text = comment.content,
                fontSize = 14.dp
            )
            Spacer(modifier = Modifier.height(4.dp))
            Text(
                text = comment.createdAt,
                fontSize = 12.dp,
                color = Color(0xFF666666)
            )
        }
    }
}

private fun likeArticle(articleId: Int) {
    ApiClient.likeArticle(articleId, object : ApiCallback {
        override fun onSuccess(response: String) {}
        override fun onError(error: String) {}
    })
}

private fun loadArticleDetail(articleId: Int, callback: (Result<Pair<ArticleDetail, List<Comment>>>) -> Unit) {
    ApiClient.getArticleDetail(articleId, object : ApiCallback {
        override fun onSuccess(response: String) {
            try {
                val json = org.json.JSONObject(response)
                val articleJson = json.getJSONObject("article")
                val commentsArray = json.optJSONArray("comments") ?: org.json.JSONArray()
                
                val article = ArticleDetail(
                    id = articleJson.optInt("id", 0),
                    title = articleJson.optString("title", ""),
                    content = articleJson.optString("content", ""),
                    authorName = articleJson.optJSONObject("user")?.optString("display_name")
                        ?: articleJson.optJSONObject("user")?.optString("nickname")
                        ?: "匿名",
                    authorAvatar = articleJson.optJSONObject("user")?.optString("avatar"),
                    viewCount = articleJson.optInt("view_count", 0),
                    likeCount = articleJson.optInt("like_count", 0),
                    createdAt = articleJson.optString("created_at", "")
                )
                
                val commentsList = mutableListOf<Comment>()
                for (i in 0 until commentsArray.length()) {
                    val commentJson = commentsArray.getJSONObject(i)
                    val comment = Comment(
                        id = commentJson.optInt("id", 0),
                        content = commentJson.optString("content", ""),
                        authorName = commentJson.optJSONObject("user")?.optString("display_name")
                            ?: commentJson.optJSONObject("user")?.optString("nickname")
                            ?: "匿名",
                        createdAt = commentJson.optString("created_at", "")
                    )
                    commentsList.add(comment)
                }
                
                callback(Result.success(Pair(article, commentsList)))
            } catch (e: Exception) {
                callback(Result.failure(e))
            }
        }
        override fun onError(error: String) {
            callback(Result.failure(Exception(error)))
        }
    })
}

// 数据类
data class ArticleDetail(
    val id: Int,
    val title: String,
    val content: String,
    val authorName: String?,
    val authorAvatar: String?,
    val viewCount: Int,
    val likeCount: Int,
    val createdAt: String
)

data class Comment(
    val id: Int,
    val content: String,
    val authorName: String?,
    val createdAt: String
)