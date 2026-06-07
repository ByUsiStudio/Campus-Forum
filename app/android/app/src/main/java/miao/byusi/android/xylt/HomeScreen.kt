package miao.byusi.android.xylt

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.compose.ui.graphics.Color
import androidx.navigation.NavHostController
import kotlinx.coroutines.delay
import top.yukonga.miuix.kmp.basic.Scaffold
import top.yukonga.miuix.kmp.basic.SmallTopAppBar
import top.yukonga.miuix.kmp.basic.Card
import top.yukonga.miuix.kmp.basic.Text
import top.yukonga.miuix.kmp.basic.Button
import top.yukonga.miuix.kmp.basic.FilledButton
import top.yukonga.miuix.kmp.basic.CircularProgressIndicator
import top.yukonga.miuix.kmp.basic.Icon
import top.yukonga.miuix.kmp.basic.NavigationBar
import top.yukonga.miuix.kmp.basic.NavigationBarItem
import top.yukonga.miuix.kmp.icon.MiuixIcons

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
            // Miuix 风格小顶栏（带返回行为/滚动折叠/大标题样式）
            SmallTopAppBar(
                title = "校园论坛",
                modifier = Modifier
            )
        },
        floatingActionButton = {
            FilledButton(
                onClick = { navController.navigate("create") },
                modifier = Modifier.padding(16.dp)
            ) {
                Text("发  布")
            }
        },
        bottomBar = {
            NavigationBar(
                items = listOf(
                    NavigationBarItem("主页", MiuixIcons.Basic.Search),
                    NavigationBarItem("草稿", MiuixIcons.Basic.Search),
                    NavigationBarItem("我的", MiuixIcons.Basic.Person)
                ),
                selected = 0,
                onClick = { index ->
                    when (index) {
                        1 -> navController.navigate("drafts")
                        2 -> navController.navigate("profile")
                    }
                }
            )
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
                    FilledButton(onClick = {
                        loadArticles { result ->
                            result.onSuccess { data -> articles = data }
                            .onFailure { e -> error = e.message }
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
    Card(
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
                color = Color(0xFF666666)
            )
            Spacer(modifier = Modifier.height(10.dp))
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceBetween
            ) {
                Text(
                    text = article.authorName ?: "匿名",
                    fontSize = 12.sp,
                    color = Color(0xFF7C4DFF)
                )
                Text(
                    text = "${article.viewCount} 浏览 · ${article.likeCount} 点赞",
                    fontSize = 12.sp,
                    color = Color(0xFF888888)
                )
            }
        }
    }
}

// 数据加载函数
private fun loadArticles(callback: (Result<List<Article>>) -> Unit) {
    ApiClient.getArticles(1, 20, object : ApiCallback {
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
        val article = Article(
            id = articleJson.optInt("id", 0),
            title = articleJson.optString("title", ""),
            content = articleJson.optString("content", ""),
            authorName = articleJson.optJSONObject("user")?.optString("display_name")
                ?: articleJson.optJSONObject("user")?.optString("nickname")
                ?: "匿名",
            viewCount = articleJson.optInt("view_count", 0),
            likeCount = articleJson.optInt("like_count", 0),
            createdAt = articleJson.optString("created_at", "")
        )
        list.add(article)
    }
    return list
}
