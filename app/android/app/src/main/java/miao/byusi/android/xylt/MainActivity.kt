package miao.byusi.android.xylt

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import top.yukonga.miuix.kmp.theme.MiuixTheme
import top.yukonga.miuix.kmp.theme.ThemeController
import top.yukonga.miuix.kmp.theme.ColorSchemeMode

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        setContent {
            AppTheme {
                AppNavigation()
            }
        }
    }
}

@Composable
fun AppTheme(
    content: @Composable () -> Unit
) {
    val controller = ThemeController(ColorSchemeMode.System)
    MiuixTheme(
        controller = controller,
        content = content
    )
}

@Composable
fun AppNavigation() {
    val navController = rememberNavController()
    
    NavHost(
        navController = navController,
        startDestination = "home", // 直接启动到主页，不要求登录
        modifier = Modifier.fillMaxSize()
    ) {
        // 主页 - 文章列表
        composable("home") {
            HomeScreen(navController)
        }
        
        // 文章详情
        composable("article/{articleId}") { backStackEntry ->
            val articleId = backStackEntry.arguments?.getString("articleId")?.toIntOrNull() ?: 0
            ArticleDetailScreen(navController, articleId)
        }
        
        // 创建文章 (WebView)
        composable("create") {
            CreateArticleScreen(navController)
        }
        
        // 草稿箱
        composable("drafts") {
            DraftsScreen(navController)
        }
        
        // 个人资料
        composable("profile") {
            ProfileScreen(navController)
        }
        
        // 登录
        composable("login") {
            LoginScreen(navController)
        }
        
        // 注册
        composable("register") {
            RegisterScreen(navController)
        }
    }
}