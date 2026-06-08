package miao.byusi.android.xylt

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.outlined.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.compose.ui.graphics.Color
import androidx.navigation.NavHostController
import kotlinx.coroutines.delay

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun HomeScreen(navController: NavHostController) {
    var articles by remember { mutableStateOf<List<Article>>(emptyList()) }
    var isLoading by remember { mutableStateOf(true) }
    var error by remember { mutableStateOf<String?>(null) }

    // 自动刷新：每30秒刷新一次
    LaunchedEffect(Unit) {
        while (true) {
            loadArticles { result ->
                result.onSuccess { data ->
                    articles = data
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
            TopAppBar(
                title = { Text("校园论坛") },
                modifier = Modifier
            )
        },
        floatingActionButton = {
            ExtendedFloatingActionButton(
                onClick = { navController.navigate("create") },
                modifier = Modifier.padding(16.dp)
            ) {
                Text("发  布")
            }
        },
        bottomBar = {
            NavigationBar {
                NavigationBarItem(
                    icon = { Icon(Icons.Outlined.Search, contentDescription = "主页") },
                    label = { Text("主页") },
                    selected = true,
                    onClick = { }
                )
                NavigationBarItem(
                    icon = { Icon(Icons.Outlined.Drafts, contentDescription = "草稿") },
                    label = { Text("草稿") },
                    selected = false,
                    onClick = { navController.navigate("drafts") }
                )
                NavigationBarItem(
                    icon = { Icon(Icons.Outlined.Person, contentDescription = "我的") },
                    label = { Text("我的") },
                    selected = false,
                    onClick = { navController.navigate("profile") }
                )
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
                    Text("加载失败: $error", fontSize = 14.sp)
                    Spacer(modifier = Modifier.height(16.dp))
                    Button(onClick = {
                        loadArticles { result ->
                            result.onSuccess { data -> articles = data }
                            result.onFailure { e -> error = e.message }
                        }
                    }) {
                        Text("重  试")
                    }
                }
            }
        } else {
            LazyColumn(
                modifier = Modifier
                    .fillMaxSize()
                    .padding(padding),
                contentPadding = PaddingValues(16.dp),
                verticalArrangement = Arrangement.spacedBy(10.dp)
            ) {
                items(articles) { article ->
                    ArticleCard(article = article, onClick = {
                        navController.navigate("article/${article.id}")
                    })
                }
            }
        }
    }
}

@Composable
fun ArticleCard(article: Article, onClick: () -> Unit) {
    ElevatedCard(
        onClick = onClick,
        modifier = Modifier.fillMaxWidth()
    ) {
        Column(
            modifier = Modifier.padding(16.dp)
        ) {
            Text(
                text = article.title,
                fontSize = 16.sp,
                fontWeight = FontWeight.SemiBold
            )
            Spacer(modifier = Modifier.height(8.dp))
            Text(
                text = article.content.take(100) + if (article.content.length > 100) "..." else "",
                fontSize = 14.sp,
                color = MaterialTheme.colorScheme.onSurfaceVariant
            )
            Spacer(modifier = Modifier.height(10.dp))
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceBetween
            ) {
                Text(
                    text = article.user?.displayName ?: article.user?.username ?: "匿名",
                    fontSize = 12.sp,
                    color = MaterialTheme.colorScheme.primary
                )
                Text(
                    text = "${article.viewCount} 浏览 · ${article.likeCount} 点赞",
                    fontSize = 12.sp,
                    color = MaterialTheme.colorScheme.onSurfaceVariant
                )
            }
        }
    }
}

// 数据加载函数
private fun loadArticles(callback: (Result<List<Article>>) -> Unit) {
    ApiClient.getArticles(1, 20, null, object : ApiCallback {
        override fun onSuccess(response: String) {
            try {
                val articles = parseArticles(response)
                callback(Result.success(articles))
            } catch (e: Exception) {
                callback(Result.failure(e))
            }
        }
        override fun onError(error: String) {
            callback(Result.failure(Exception(error)))
        }
    })
}

private fun parseArticles(response: String): List<Article> {
    val json = org.json.JSONObject(response)
    val articlesArray = json.optJSONArray("articles") ?: return emptyList()
    val list = mutableListOf<Article>()
    for (i in 0 until articlesArray.length()) {
        val articleJson = articlesArray.getJSONObject(i)
        val userJson = articleJson.optJSONObject("user")
        val user = userJson?.let {
            User(
                id = it.optInt("id", 0),
                username = it.optString("username", ""),
                displayName = it.optString("display_name", it.optString("nickname", null)),
                avatar = it.optString("avatar_url", it.optString("avatar", null))
            )
        }
        val article = Article(
            id = articleJson.optInt("id", 0),
            title = articleJson.optString("title", ""),
            content = articleJson.optString("content", ""),
            user = user,
            viewCount = articleJson.optInt("view_count", 0),
            likeCount = articleJson.optInt("like_count", 0),
            createdAt = articleJson.optString("created_at", "")
        )
        list.add(article)
    }
    return list
}
