package miao.byusi.android.xylt

import android.annotation.SuppressLint
import android.webkit.CookieManager
import android.webkit.WebSettings
import android.webkit.WebView
import android.webkit.WebViewClient
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.heightIn
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.compose.ui.viewinterop.AndroidView


@SuppressLint("SetJavaScriptEnabled")
@Composable
fun ArticleContentView(
    articleId: Int,
    modifier: Modifier = Modifier
) {
    val webViewRef = remember { mutableStateOf<WebView?>(null) }
    val baseApi = remember { ConfigManager.getBaseApi() }

    AndroidView(
        modifier = modifier.fillMaxWidth().heightIn(min = 100.dp),
        factory = { ctx ->
            WebView(ctx).apply {
                settings.apply {
                    javaScriptEnabled = true
                    javaScriptCanOpenWindowsAutomatically = false
                    domStorageEnabled = true
                    // 允许 HTTPS 页面加载 HTTP CDN 资源（marked.js）
                    mixedContentMode = WebSettings.MIXED_CONTENT_COMPATIBILITY_MODE
                    useWideViewPort = true
                    loadWithOverviewMode = true
                    cacheMode = WebSettings.LOAD_DEFAULT
                    isVerticalScrollBarEnabled = false
                    isHorizontalScrollBarEnabled = false
                    setSupportZoom(true)
                    builtInZoomControls = true
                    displayZoomControls = false
                }
                setBackgroundColor(0)

                // 注入登录 Token（与 CreateArticleScreen 一致）
                val cookieManager = CookieManager.getInstance()
                cookieManager.setAcceptCookie(true)
                cookieManager.setAcceptThirdPartyCookies(this, true)
                val token = ApiClient.getAuthToken()
                if (token.isNotEmpty()) {
                    cookieManager.setCookie("https://xylt.cdifit.cn", "auth_token=$token")
                    cookieManager.setCookie("https://xylt.cdifit.cn", "token=$token")
                }

                webViewClient = object : WebViewClient() {
                    override fun shouldOverrideUrlLoading(
                        view: WebView?,
                        request: android.webkit.WebResourceRequest?
                    ): Boolean {
                        val url = request?.url?.toString().orEmpty()
                        // 同域 / articles 路径不拦截
                        if (url.startsWith("https://xylt.cdifit.cn/")) return false
                        // 外链跳转系统浏览器
                        try {
                            val intent = android.content.Intent(
                                android.content.Intent.ACTION_VIEW,
                                android.net.Uri.parse(url)
                            )
                            intent.addFlags(android.content.Intent.FLAG_ACTIVITY_NEW_TASK)
                            ctx.startActivity(intent)
                        } catch (_: Exception) { }
                        return true
                    }
                }

                webViewRef.value = this
            }
        },
        update = { webView ->
            // 加载本地 HTML 模板，传递 articleId 和 baseApi 作为 URL 参数
            val htmlUrl = "file:///android_asset/ui/articles.html?articleId=$articleId&baseApi=$baseApi"
            webView.loadUrl(htmlUrl)
        }
    )

    // articleId 变化时重新加载
    LaunchedEffect(articleId, baseApi) {
        webViewRef.value?.let { wv ->
            wv.loadUrl("file:///android_asset/ui/articles.html?articleId=$articleId&baseApi=$baseApi")
        }
    }
}
