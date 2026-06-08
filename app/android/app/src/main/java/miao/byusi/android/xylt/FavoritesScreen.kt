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
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.unit.sp
import androidx.navigation.NavHostController

@Composable
fun FavoritesScreen(navController: NavHostController) {
    var articles by remember { mutableStateOf<List<Article>>(emptyList()) }
    var isLoading by remember { mutableStateOf(true) }
    var error by remember { mutableStateOf<String?>(null) }
    
    val loadData = {
        isLoading = true
        error = null
        loadFavorites { result ->
            result.onSuccess { list ->
                articles = list
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
                        text = "我的收藏",
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
                    Text("加载失败: $error", fontSize = 14.sp)
                    Spacer(modifier = Modifier.height(16.dp))
                    Button(onClick = loadData) {
                        Text("重试")
                    }
                }
            }
        } else {
            if (articles.isEmpty()) {
                Box(
                    modifier = Modifier
                        .fillMaxSize()
                        .padding(padding),
                    contentAlignment = Alignment.Center
                ) {
                    Column(horizontalAlignment = Alignment.CenterHorizontally) {
                        Icon(
                            imageVector = Icons.Outlined.Bookmark,
                            contentDescription = null,
                            modifier = Modifier.size(64.dp),
                            tint = Color(0xFFCCCCCC)
                        )
                        Spacer(modifier = Modifier.height(16.dp))
                        Text(
                            text = "暂无收藏",
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
}

private fun loadFavorites(callback: (Result<List<Article>>) -> Unit) {
    ApiClient.getFavorites(1, 50, object : ApiCallback {
        override fun onSuccess(response: String) {
            try {
                val list = parseFavoriteArticles(response)
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

private fun parseFavoriteArticles(response: String): List<Article> {
    val json = org.json.JSONObject(response)
    val articlesArray = json.optJSONArray("articles") ?: json.optJSONArray("favorites") ?: org.json.JSONArray()
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