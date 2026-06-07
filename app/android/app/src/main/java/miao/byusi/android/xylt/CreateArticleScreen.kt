package miao.byusi.android.xylt

import android.webkit.CookieManager
import android.webkit.WebView
import android.webkit.WebViewClient
import androidx.compose.foundation.layout.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.viewinterop.AndroidView
import androidx.navigation.NavHostController
import miuix.compose.scaffold.Scaffold
import miuix.compose.surface.Surface
import miuix.compose.text.Text
import miuix.compose.button.Button
import miuix.compose.progress.CircularProgressIndicator
import miuix.compose.icon.Icon
import miuix.compose.icon.icons.MiuixIcons
import miuix.compose.icon.icons.outlined.Close

@Composable
fun CreateArticleScreen(navController: NavHostController) {
    var isLoading by remember { mutableStateOf(true) }
    
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
                            imageVector = MiuixIcons.Outlined.Close,
                            contentDescription = "关闭"
                        )
                    }
                    Text(
                        text = "创建文章",
                        fontSize = 20.dp
                    )
                }
            }
        }
    ) { padding ->
        Box(modifier = Modifier.padding(padding)) {
            if (isLoading) {
                CircularProgressIndicator(
                    modifier = Modifier.align(Alignment.Center)
                )
            }
            
            AndroidView(
                modifier = Modifier.fillMaxSize(),
                factory = { context ->
                    WebView(context).apply {
                        settings.javaScriptEnabled = true
                        settings.domStorageEnabled = true
                        settings.allowFileAccess = true
                        settings.mixedContentMode = android.webkit.WebSettings.MIXED_CONTENT_ALWAYS_ALLOW
                        
                        // 设置 Cookie 保持登录状态
                        val cookieManager = CookieManager.getInstance()
                        cookieManager.setAcceptCookie(true)
                        cookieManager.setAcceptThirdPartyCookies(this, true)
                        
                        val token = ApiClient.getAuthToken()
                        if (!token.isEmpty()) {
                            cookieManager.setCookie("https://xylt.cdifit.cn", "auth_token=$token")
                        }
                        
                        webViewClient = object : WebViewClient() {
                            override fun onPageFinished(view: WebView?, url: String?) {
                                isLoading = false
                                
                                // 文章创建成功后自动关闭
                                if (url != null && url.matches(Regex("https://xylt.cdifit.cn/articles/\\d+"))) {
                                    navController.popBackStack()
                                }
                                
                                // 返回主页时关闭
                                if (url == "https://xylt.cdifit.cn/" || url == "https://xylt.cdifit.cn") {
                                    navController.popBackStack()
                                }
                            }
                            
                            override fun shouldOverrideUrlLoading(view: WebView?, request: android.webkit.WebResourceRequest?): Boolean {
                                val url = request?.url?.toString() ?: return false
                                
                                // 文章创建成功后自动关闭
                                if (url.matches(Regex("https://xylt.cdifit.cn/articles/\\d+"))) {
                                    navController.popBackStack()
                                    return true
                                }
                                
                                return false
                            }
                        }
                        
                        loadUrl("https://xylt.cdifit.cn/create")
                    }
                }
            )
        }
    }
}