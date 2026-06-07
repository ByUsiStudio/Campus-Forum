package miao.byusi.android.xylt

import androidx.compose.foundation.layout.*
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

@Composable
fun ProfileScreen(navController: NavHostController) {
    var profile by remember { mutableStateOf<UserProfile?>(null) }
    var isLoading by remember { mutableStateOf(true) }
    var error by remember { mutableStateOf<String?>(null) }
    
    // 自动刷新：每30秒刷新一次
    LaunchedEffect(Unit) {
        while (true) {
            loadProfile { result ->
                result.onSuccess { data ->
                    profile = data
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
                        text = "个人资料",
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
                    Text("请登录查看个人资料")
                    Spacer(modifier = Modifier.height(16.dp))
                    FilledButton(onClick = { navController.navigate("login") }) {
                        Text("登录")
                    }
                }
            }
        } else if (profile != null) {
            Column(
                modifier = Modifier
                    .fillMaxSize()
                    .padding(padding),
                horizontalAlignment = Alignment.CenterHorizontally
            ) {
                Spacer(modifier = Modifier.height(32.dp))
                
                AsyncImage(
                    model = profile?.avatar,
                    contentDescription = "头像",
                    modifier = Modifier.size(96.dp)
                )
                
                Spacer(modifier = Modifier.height(16.dp))
                
                Text(
                    text = profile?.displayName ?: profile?.username ?: "",
                    fontSize = 20.dp
                )
                
                Text(
                    text = "@${profile?.username ?: ""}",
                    fontSize = 14.dp,
                    color = Color(0xFF666666)
                )
                
                Spacer(modifier = Modifier.height(8.dp))
                
                Text(
                    text = profile?.signature ?: "暂无个人简介",
                    fontSize = 14.dp,
                    color = Color(0xFF666666)
                )
                
                Spacer(modifier = Modifier.height(24.dp))
                
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.SpaceEvenly
                ) {
                    Column(horizontalAlignment = Alignment.CenterHorizontally) {
                        Text(
                            text = "${profile?.articleCount ?: 0}",
                            fontSize = 18.dp
                        )
                        Text(
                            text = "文章",
                            fontSize = 12.dp,
                            color = Color(0xFF666666)
                        )
                    }
                    Column(horizontalAlignment = Alignment.CenterHorizontally) {
                        Text(
                            text = "${profile?.followerCount ?: 0}",
                            fontSize = 18.dp
                        )
                        Text(
                            text = "粉丝",
                            fontSize = 12.dp,
                            color = Color(0xFF666666)
                        )
                    }
                    Column(horizontalAlignment = Alignment.CenterHorizontally) {
                        Text(
                            text = "${profile?.followingCount ?: 0}",
                            fontSize = 18.dp
                        )
                        Text(
                            text = "关注",
                            fontSize = 12.dp,
                            color = Color(0xFF666666)
                        )
                    }
                }
            }
        }
    }
}

private fun loadProfile(callback: (Result<UserProfile>) -> Unit) {
    ApiClient.getProfile(object : ApiCallback {
        override fun onSuccess(response: String) {
            try {
                val json = org.json.JSONObject(response)
                val profile = UserProfile(
                    id = json.optInt("id", 0),
                    username = json.optString("username", ""),
                    displayName = json.optString("display_name", json.optString("nickname", "")),
                    avatar = json.optString("avatar", json.optString("avatar_url", "")),
                    signature = json.optString("signature", ""),
                    articleCount = json.optInt("article_count", 0),
                    followerCount = json.optInt("follower_count", 0),
                    followingCount = json.optInt("following_count", 0)
                )
                callback(Result.success(profile))
            } catch (e: Exception) {
                callback(Result.failure(e))
            }
        }
        override fun onError(error: String) {
            callback(Result.failure(Exception(error)))
        }
    })
}

data class UserProfile(
    val id: Int,
    val username: String,
    val displayName: String?,
    val avatar: String?,
    val signature: String?,
    val articleCount: Int,
    val followerCount: Int,
    val followingCount: Int
)