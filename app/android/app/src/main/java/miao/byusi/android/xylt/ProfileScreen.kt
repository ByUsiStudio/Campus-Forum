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
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.navigation.NavHostController
import kotlinx.coroutines.delay
import coil.compose.AsyncImage
import kotlinx.coroutines.launch

@Composable
fun ProfileScreen(navController: NavHostController) {
    var profile by remember { mutableStateOf<UserProfile?>(null) }
    var isLoading by remember { mutableStateOf(true) }
    var error by remember { mutableStateOf<String?>(null) }
    var showEditDialog by remember { mutableStateOf(false) }
    val context = LocalContext.current
    val scope = rememberCoroutineScope()
    var checking by remember { mutableStateOf(false) }
    var healthMsg by remember { mutableStateOf<String?>(null) }

    val loadData = {
        isLoading = true
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
    }

    LaunchedEffect(Unit) {
        loadData()
        while (true) {
            delay(30000)
            loadProfile { result ->
                result.onSuccess { data ->
                    profile = data
                }
            }
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
                        .padding(horizontal = 12.dp, vertical = 8.dp),
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Button(onClick = { navController.popBackStack() }) {
                        Icon(
                            imageVector = Icons.Outlined.ArrowBack,
                            contentDescription = "返回"
                        )
                    }
                    Spacer(modifier = Modifier.width(4.dp))
                    Text(
                        text = "个人资料",
                        fontSize = 20.sp,
                        fontWeight = FontWeight.SemiBold
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
                    Text("请登录查看个人资料", fontSize = 14.sp)
                    Spacer(modifier = Modifier.height(16.dp))
                    Button(onClick = { navController.navigate("login") }) {
                        Text("登录")
                    }
                }
            }
        } else if (profile != null) {
            Column(
                modifier = Modifier
                    .fillMaxSize()
                    .padding(padding)
                    .verticalScroll(rememberScrollState()),
                horizontalAlignment = Alignment.CenterHorizontally
            ) {
                Spacer(modifier = Modifier.height(24.dp))

                AsyncImage(
                    model = profile?.avatar,
                    contentDescription = "头像",
                    modifier = Modifier.size(96.dp)
                )

                Spacer(modifier = Modifier.height(12.dp))

                Text(
                    text = profile?.displayName ?: profile?.username ?: "",
                    fontSize = 20.sp,
                    fontWeight = FontWeight.Bold
                )

                Text(
                    text = "@${profile?.username ?: ""}",
                    fontSize = 13.sp,
                    color = Color(0xFF888888)
                )

                Spacer(modifier = Modifier.height(6.dp))

                Text(
                    text = profile?.signature ?: "暂无个人简介",
                    fontSize = 14.sp,
                    color = Color(0xFF666666)
                )

                Spacer(modifier = Modifier.height(16.dp))

                OutlinedButton(onClick = { showEditDialog = true }) {
                    Icon(
                        imageVector = Icons.Outlined.Edit,
                        contentDescription = null
                    )
                    Spacer(modifier = Modifier.width(4.dp))
                    Text("编辑资料")
                }

                Spacer(modifier = Modifier.height(20.dp))

                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.SpaceEvenly
                ) {
                    Column(horizontalAlignment = Alignment.CenterHorizontally) {
                        Text(text = "${profile?.articleCount ?: 0}", fontSize = 18.sp, fontWeight = FontWeight.SemiBold)
                        Text(text = "文章", fontSize = 12.sp, color = Color(0xFF888888))
                    }
                    Column(horizontalAlignment = Alignment.CenterHorizontally) {
                        Text(text = "${profile?.followerCount ?: 0}", fontSize = 18.sp, fontWeight = FontWeight.SemiBold)
                        Text(text = "粉丝", fontSize = 12.sp, color = Color(0xFF888888))
                    }
                    Column(horizontalAlignment = Alignment.CenterHorizontally) {
                        Text(text = "${profile?.followingCount ?: 0}", fontSize = 18.sp, fontWeight = FontWeight.SemiBold)
                        Text(text = "关注", fontSize = 12.sp, color = Color(0xFF888888))
                    }
                }

                Spacer(modifier = Modifier.height(24.dp))

                // 功能菜单
                Column(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(horizontal = 16.dp),
                    verticalArrangement = Arrangement.spacedBy(12.dp)
                ) {
                    MenuItem(
                        icon = Icons.Outlined.Bookmark,
                        title = "我的收藏",
                        onClick = { navController.navigate("favorites") }
                    )
                    MenuItem(
                        icon = Icons.Outlined.Notifications,
                        title = "我的通知",
                        onClick = { navController.navigate("notifications") }
                    )
                    MenuItem(
                        icon = Icons.Outlined.Drafts,
                        title = "草稿箱",
                        onClick = { navController.navigate("drafts") }
                    )
                }

                Spacer(modifier = Modifier.height(24.dp))

                // 服务器状态卡（Miuix 风格）
                ServerStatusCard(
                    apiUrl = ServerConfigStore.apiUrl(context),
                    backendVersion = ServerConfigStore.backendVersion(context),
                    frontendVersion = ServerConfigStore.frontendVersion(context),
                    lastCheckAt = ServerConfigStore.lastCheckAt(context),
                    checking = checking,
                    healthMsg = healthMsg,
                    onReCheck = {
                        if (!checking) {
                            checking = true
                            healthMsg = null
                            scope.launch {
                                val r = HealthClient.check()
                                if (r.healthy) {
                                    ServerConfigStore.saveFromResult(context, r)
                                    healthMsg = "服务器在线（${r.backendVersion ?: "unknown"}）"
                                } else {
                                    healthMsg = r.message
                                }
                                checking = false
                            }
                        }
                    }
                )

                Spacer(modifier = Modifier.height(24.dp))
            }
        }
    }

    if (showEditDialog && profile != null) {
        EditProfileDialog(
            currentDisplayName = profile?.displayName ?: "",
            currentSignature = profile?.signature ?: "",
            onDismiss = { showEditDialog = false },
            onSave = { displayName, signature ->
                updateProfile(displayName, signature) {
                    showEditDialog = false
                    loadData()
                }
            }
        )
    }
}

@Composable
private fun MenuItem(
    icon: androidx.compose.ui.graphics.vector.ImageVector,
    title: String,
    onClick: () -> Unit
) {
    Surface(
        modifier = Modifier.fillMaxWidth(),
        onClick = onClick
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(16.dp),
            verticalAlignment = Alignment.CenterVertically
        ) {
            Icon(
                imageVector = icon,
                contentDescription = null,
                tint = Color(0xFF7C4DFF)
            )
            Spacer(modifier = Modifier.width(12.dp))
            Text(
                text = title,
                fontSize = 16.sp
            )
            Spacer(modifier = Modifier.weight(1f))
            Icon(
                imageVector = Icons.Outlined.ArrowForwardIos,
                contentDescription = null,
                tint = Color(0xFFCCCCCC),
                modifier = Modifier.size(16.dp)
            )
        }
    }
}

@Composable
private fun EditProfileDialog(
    currentDisplayName: String,
    currentSignature: String,
    onDismiss: () -> Unit,
    onSave: (String, String) -> Unit
) {
    var displayName by remember { mutableStateOf(currentDisplayName) }
    var signature by remember { mutableStateOf(currentSignature) }
    
    AlertDialog(
        onDismissRequest = onDismiss,
        title = { Text("编辑资料") },
        text = {
            Column(modifier = Modifier.padding(0.dp)) {
                Text(
                    text = "昵称",
                    fontSize = 14.sp,
                    color = Color(0xFF666666)
                )
                Spacer(modifier = Modifier.height(4.dp))
                TextField(
                    value = displayName,
                    onValueChange = { displayName = it },
                    modifier = Modifier.fillMaxWidth(),
                    placeholder = { Text("请输入昵称") }
                )
                
                Spacer(modifier = Modifier.height(16.dp))
                
                Text(
                    text = "个人简介",
                    fontSize = 14.sp,
                    color = Color(0xFF666666)
                )
                Spacer(modifier = Modifier.height(4.dp))
                TextField(
                    value = signature,
                    onValueChange = { signature = it },
                    modifier = Modifier.fillMaxWidth(),
                    placeholder = { Text("请输入个人简介") },
                    minLines = 3,
                    maxLines = 5
                )
            }
        },
        dismissButton = {
            OutlinedButton(onClick = onDismiss) {
                Text("取消")
            }
        },
        confirmButton = {
            Button(
                onClick = { onSave(displayName, signature) },
                enabled = displayName.isNotBlank()
            ) {
                Text("保存")
            }
        }
    )
}

@Composable
private fun ServerStatusCard(
    apiUrl: String,
    backendVersion: String?,
    frontendVersion: String?,
    lastCheckAt: Long,
    checking: Boolean,
    healthMsg: String?,
    onReCheck: () -> Unit
) {
    Card(
        modifier = Modifier
            .fillMaxWidth()
            .padding(horizontal = 16.dp)
    ) {
        Column(
            modifier = Modifier
                .fillMaxWidth()
                .padding(horizontal = 18.dp, vertical = 16.dp)
        ) {
            Text(
                text = "服务器状态",
                fontSize = 15.sp,
                fontWeight = FontWeight.SemiBold,
                color = Color(0xFF7C4DFF)
            )
            Spacer(modifier = Modifier.height(10.dp))
            InfoLine("API 地址", apiUrl)
            InfoLine("后端版本", backendVersion ?: "—")
            InfoLine("前端版本", frontendVersion ?: "—")
            InfoLine(
                "上次校验",
                if (lastCheckAt <= 0L) "—" else formatTime(lastCheckAt)
            )
            if (healthMsg != null) {
                Spacer(modifier = Modifier.height(8.dp))
                Text(
                    text = healthMsg,
                    fontSize = 12.sp,
                    color = if (healthMsg.contains("在线")) Color(0xFF2E7D32) else Color(0xFFC62828)
                )
            }
            Spacer(modifier = Modifier.height(12.dp))
            Button(
                onClick = onReCheck,
                modifier = Modifier
                    .fillMaxWidth()
                    .height(40.dp)
            ) {
                if (checking) {
                    CircularProgressIndicator(
                        modifier = Modifier.size(18.dp)
                    )
                    Spacer(modifier = Modifier.width(8.dp))
                    Text("校验中…")
                } else {
                    Text("重新校验服务器")
                }
            }
        }
    }
}

@Composable
private fun InfoLine(label: String, value: String) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .padding(vertical = 3.dp)
    ) {
        Text(
            text = label,
            fontSize = 13.sp,
            color = Color(0xFF888888),
            modifier = Modifier.width(78.dp)
        )
        Text(
            text = value,
            fontSize = 13.sp,
            color = Color(0xFF333333),
            modifier = Modifier.weight(1f)
        )
    }
}

private fun formatTime(ms: Long): String {
    val sdf = java.text.SimpleDateFormat("yyyy-MM-dd HH:mm:ss", java.util.Locale.getDefault())
    return sdf.format(java.util.Date(ms))
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

private fun updateProfile(displayName: String, signature: String, onSuccess: () -> Unit) {
    ApiClient.updateProfile(displayName, signature, object : ApiCallback {
        override fun onSuccess(response: String) {
            onSuccess()
        }
        override fun onError(error: String) {}
    })
}
// UserProfile 数据类统一在 Models.kt 声明。