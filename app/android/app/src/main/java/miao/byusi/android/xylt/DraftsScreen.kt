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
import kotlinx.coroutines.delay

@Composable
fun DraftsScreen(navController: NavHostController) {
    var drafts by remember { mutableStateOf<List<Draft>>(emptyList()) }
    var isLoading by remember { mutableStateOf(true) }
    var error by remember { mutableStateOf<String?>(null) }
    
    // 自动刷新：每30秒刷新一次
    LaunchedEffect(Unit) {
        while (true) {
            loadDrafts { result ->
                result.onSuccess { data ->
                    drafts = data
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
                            imageVector = Icons.Outlined.ArrowBack,
                            contentDescription = "返回"
                        )
                    }
                    Text(
                        text = "草稿箱",
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
                Text("加载失败: $error")
            }
        } else if (drafts.isEmpty()) {
            Box(
                modifier = Modifier
                    .fillMaxSize()
                    .padding(padding),
                contentAlignment = Alignment.Center
            ) {
                Text("暂无草稿")
            }
        } else {
            LazyColumn(
                modifier = Modifier
                    .fillMaxSize()
                    .padding(padding),
                contentPadding = PaddingValues(16.dp),
                verticalArrangement = Arrangement.spacedBy(8.dp)
            ) {
                items(drafts) { draft ->
                    DraftCard(
                        draft = draft,
                        onPublish = { publishDraft(draft.id) },
                        onDelete = { deleteDraft(draft.id) }
                    )
                }
            }
        }
    }
}

@Composable
fun DraftCard(draft: Draft, onPublish: () -> Unit, onDelete: () -> Unit) {
    Card(
        modifier = Modifier.fillMaxWidth()
    ) {
        Column(
            modifier = Modifier.padding(16.dp)
        ) {
            Text(
                text = draft.title.ifEmpty { "无标题草稿" },
                fontSize = 16.sp
            )
            Spacer(modifier = Modifier.height(8.dp))
            Text(
                text = draft.content.take(100) + if (draft.content.length > 100) "..." else "",
                fontSize = 14.sp,
                color = Color(0xFF666666)
            )
            Spacer(modifier = Modifier.height(8.dp))
            Text(
                text = draft.createdAt,
                fontSize = 12.sp,
                color = Color(0xFF666666)
            )
            Spacer(modifier = Modifier.height(16.dp))
            Row(
                horizontalArrangement = Arrangement.spacedBy(8.dp)
            ) {
                Button(onClick = onPublish) {
                    Text("发布")
                }
                OutlinedButton(onClick = onDelete) {
                    Text("删除")
                }
            }
        }
    }
}

private fun loadDrafts(callback: (Result<List<Draft>>) -> Unit) {
    ApiClient.getMyDrafts(1, 50, object : ApiCallback {
        override fun onSuccess(response: String) {
            try {
                val json = org.json.JSONObject(response)
                val draftsArray = json.optJSONArray("drafts") ?: return callback(Result.success(emptyList()))
                val list = mutableListOf<Draft>()
                for (i in 0 until draftsArray.length()) {
                    val draftJson = draftsArray.getJSONObject(i)
                    val draft = Draft(
                        id = draftJson.optInt("id", 0),
                        title = draftJson.optString("title", ""),
                        content = draftJson.optString("content", ""),
                        createdAt = draftJson.optString("created_at", "")
                    )
                    list.add(draft)
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

private fun publishDraft(draftId: Int) {
    ApiClient.publishDraft(draftId, object : ApiCallback {
        override fun onSuccess(response: String) {}
        override fun onError(error: String) {}
    })
}

private fun deleteDraft(draftId: Int) {
    ApiClient.deleteArticle(draftId, object : ApiCallback {
        override fun onSuccess(response: String) {}
        override fun onError(error: String) {}
    })
}
// Draft 数据类统一在 Models.kt 声明。