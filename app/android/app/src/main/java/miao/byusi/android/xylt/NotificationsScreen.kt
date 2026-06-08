package miao.byusi.android.xylt

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
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
import coil.compose.AsyncImage

@Composable
fun NotificationsScreen(navController: NavHostController) {
    var notifications by remember { mutableStateOf<List<Notification>>(emptyList()) }
    var isLoading by remember { mutableStateOf(true) }
    var error by remember { mutableStateOf<String?>(null) }
    var unreadCount by remember { mutableStateOf(0) }
    
    val loadData = {
        isLoading = true
        error = null
        loadUnreadCount { count ->
            unreadCount = count
        }
        loadNotifications { result ->
            result.onSuccess { list ->
                notifications = list
                isLoading = false
            }.onFailure { e ->
                error = e.message
                isLoading = false
            }
        }
    }
    
    LaunchedEffect(Unit) {
        loadData()
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
                        text = "通知",
                        fontSize = 20.sp,
                        modifier = Modifier.weight(1f)
                    )
                    if (unreadCount > 0) {
                        Button(onClick = { 
                            markAllAsRead {
                                loadData()
                            }
                        }) {
                            Icon(
                                imageVector = Icons.Outlined.MarkEmailRead,
                                contentDescription = null
                            )
                            Spacer(modifier = Modifier.width(4.dp))
                            Text("全部已读")
                        }
                    }
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
                    Button(onClick = loadData) {
                        Text("重试")
                    }
                }
            }
        } else {
            if (notifications.isEmpty()) {
                Box(
                    modifier = Modifier
                        .fillMaxSize()
                        .padding(padding),
                    contentAlignment = Alignment.Center
                ) {
                    Column(horizontalAlignment = Alignment.CenterHorizontally) {
                        Icon(
                            imageVector = Icons.Outlined.Notifications,
                            contentDescription = null,
                            modifier = Modifier.size(64.dp),
                            tint = Color(0xFFCCCCCC)
                        )
                        Spacer(modifier = Modifier.height(16.dp))
                        Text(
                            text = "暂无通知",
                            fontSize = 16.sp,
                            color = Color(0xFF666666)
                        )
                    }
                }
            } else {
                LazyColumn(
                    modifier = Modifier
                        .fillMaxSize()
                        .padding(padding),
                    contentPadding = PaddingValues(16.dp),
                    verticalArrangement = Arrangement.spacedBy(12.dp)
                ) {
                    items(notifications) { notification ->
                        NotificationItem(
                            notification = notification,
                            onMarkRead = { 
                                markAsRead(it) {
                                    loadData()
                                }
                            },
                            onClick = {
                                if (!notification.isRead) {
                                    markAsRead(notification.id) {}
                                }
                                notification.articleId?.let { articleId ->
                                    navController.navigate("article/$articleId")
                                }
                            }
                        )
                    }
                }
            }
        }
    }
}

@Composable
fun NotificationItem(
    notification: Notification,
    onMarkRead: (Int) -> Unit,
    onClick: () -> Unit
) {
    Surface(
        modifier = Modifier.fillMaxWidth(),
        onClick = onClick
    ) {
        Row(
            modifier = Modifier.padding(16.dp),
            verticalAlignment = Alignment.CenterVertically
        ) {
            if (!notification.isRead) {
                Box(
                    modifier = Modifier
                        .size(8.dp)
                        .padding(end = 8.dp),
                    contentAlignment = Alignment.Center
                ) {
                    Surface(
                        modifier = Modifier.size(8.dp),
                        color = Color(0xFFFF5722)
                    ) {}
                }
            } else {
                Spacer(modifier = Modifier.width(16.dp))
            }
            
            AsyncImage(
                model = notification.sender?.avatar,
                contentDescription = "头像",
                modifier = Modifier.size(40.dp)
            )
            
            Spacer(modifier = Modifier.width(12.dp))
            
            Column(modifier = Modifier.weight(1f)) {
                Text(
                    text = getNotificationText(notification),
                    fontSize = 14.sp
                )
                Spacer(modifier = Modifier.height(4.dp))
                Text(
                    text = notification.createdAt,
                    fontSize = 12.sp,
                    color = Color(0xFF666666)
                )
            }
            
            if (!notification.isRead) {
                Button(onClick = { onMarkRead(notification.id) }) {
                    Icon(
                        imageVector = Icons.Outlined.MarkEmailRead,
                        contentDescription = "标记已读",
                        modifier = Modifier.size(20.dp)
                    )
                }
            }
        }
    }
}

private fun getNotificationText(notification: Notification): String {
    val senderName = notification.sender?.displayName ?: notification.sender?.username ?: "某人"
    return when (notification.type) {
        "like" -> "$senderName 赞了你的文章"
        "comment" -> "$senderName 评论了你的文章"
        "reply" -> "$senderName 回复了你的评论"
        "follow" -> "$senderName 关注了你"
        "mention" -> "$senderName 在文章中提到了你"
        else -> "$senderName 有新动态"
    }
}

private fun loadUnreadCount(callback: (Int) -> Unit) {
    ApiClient.getUnreadNotificationCount(object : ApiCallback {
        override fun onSuccess(response: String) {
            try {
                val json = org.json.JSONObject(response)
                callback(json.optInt("count", 0))
            } catch (e: Exception) {
                callback(0)
            }
        }
        override fun onError(error: String) {
            callback(0)
        }
    })
}

private fun loadNotifications(callback: (Result<List<Notification>>) -> Unit) {
    ApiClient.getNotifications(1, 50, object : ApiCallback {
        override fun onSuccess(response: String) {
            try {
                val json = org.json.JSONObject(response)
                val notificationsArray = json.optJSONArray("notifications") ?: org.json.JSONArray()
                val list = mutableListOf<Notification>()
                for (i in 0 until notificationsArray.length()) {
                    val notificationJson = notificationsArray.getJSONObject(i)
                    list.add(
                        Notification(
                            id = notificationJson.optInt("id", 0),
                            userId = notificationJson.optInt("user_id", 0),
                            senderId = notificationJson.optInt("sender_id").takeIf { it != 0 },
                            articleId = notificationJson.optInt("article_id").takeIf { it != 0 },
                            type = notificationJson.optString("type", ""),
                            isRead = notificationJson.optBoolean("is_read", false),
                            sender = parseUser(notificationJson.optJSONObject("sender")),
                            article = parseArticle(notificationJson.optJSONObject("article")),
                            createdAt = notificationJson.optString("created_at", "")
                        )
                    )
                }
                callback(Result.success(list))
            } catch (e: Exception) {
                callback(Result.failure(e))
            }
        }
        override fun onError(error: String) {
            callback(Result.failure(Exception(error)))
        }
    })
}

private fun markAsRead(notificationId: Int, onSuccess: () -> Unit) {
    ApiClient.markNotificationAsRead(notificationId, object : ApiCallback {
        override fun onSuccess(response: String) {
            onSuccess()
        }
        override fun onError(error: String) {}
    })
}

private fun markAllAsRead(onSuccess: () -> Unit) {
    ApiClient.markAllNotificationsAsRead(object : ApiCallback {
        override fun onSuccess(response: String) {
            onSuccess()
        }
        override fun onError(error: String) {}
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

private fun parseArticle(json: org.json.JSONObject?): Article? {
    if (json == null) return null
    return Article(
        id = json.optInt("id", 0),
        title = json.optString("title", ""),
        content = json.optString("content", ""),
        contentHtml = json.optString("content_html"),
        user = parseUser(json.optJSONObject("user")),
        category = parseCategory(json.optJSONObject("category")),
        viewCount = json.optInt("view_count", 0),
        likeCount = json.optInt("like_count", 0),
        commentCount = json.optInt("comment_count", 0),
        shareCount = json.optInt("share_count", 0),
        voiceUrl = json.optString("voice_url"),
        isAnonymous = json.optBoolean("is_anonymous", false),
        createdAt = json.optString("created_at", ""),
        status = json.optString("status")
    )
}

private fun parseCategory(json: org.json.JSONObject?): Category? {
    if (json == null) return null
    return Category(
        id = json.optInt("id", 0),
        name = json.optString("name", "")
    )
}