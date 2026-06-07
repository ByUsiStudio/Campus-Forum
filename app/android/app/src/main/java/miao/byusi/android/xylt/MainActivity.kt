package miao.byusi.android.xylt

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.animation.Crossfade
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.rememberCoroutineScope
import androidx.compose.runtime.setValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.toArgb
import androidx.compose.ui.platform.LocalContext
import androidx.core.view.WindowCompat
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import kotlinx.coroutines.launch
import top.yukonga.miuix.kmp.theme.ColorSchemeMode
import top.yukonga.miuix.kmp.theme.MiuixTheme
import top.yukonga.miuix.kmp.theme.ThemeController
import top.yukonga.miuix.kmp.utils.platformCompositionLocals

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        // 初始化配置管理器
        ConfigManager.init(this)
        enableEdgeToEdge()
        // 让系统状态栏/导航栏背景透明，背景由 Compose 主题色绘制
        WindowCompat.setDecorFitsSystemWindows(window, false)
        window.statusBarColor = Color.Transparent.toArgb()
        window.navigationBarColor = Color.Transparent.toArgb()
        setContent {
            AppTheme {
                AppRoot()
            }
        }
    }
}

/**
 * 应用主题：基于 Miuix 0.8.x。
 *  - [ColorSchemeMode.System] 跟随系统深色模式
 *  - [themeColor] 取校园论坛主色（紫）
 *  - [platformCompositionLocals] 提供 MIUI/HyperOS 平台本地化支持
 */
@Composable
fun AppTheme(
    content: @Composable () -> Unit
) {
    val controller = remember { ThemeController(ColorSchemeMode.System) }
    MiuixTheme(
        controller = controller,
        enableDynamicColor = false,
        basicColors = null,
        colorScheme = null,
        content = {
            platformCompositionLocals {
                content()
            }
        }
    )
}

/**
 * 根 Composable：先显示启动校验页面，校验通过后再渲染主导航。
 *
 * 启动校验是阻塞的：服务器不在线时不允许进入主界面，给用户明确的"重试"入口，
 * 避免在没数据时进入主界面造成大面积的"网络错误"提示。
 */
@Composable
fun AppRoot() {
    var splashState by remember { mutableStateOf<SplashState>(SplashState.Checking) }
    var navigated by remember { mutableStateOf(false) }
    val scope = rememberCoroutineScope()
    val context = LocalContext.current

    val doCheck: () -> Unit = {
        if (splashState !is SplashState.Checking) {
            splashState = SplashState.Checking
        }
        scope.launch {
            val result = HealthClient.check()
            // 校验成功时把版本信息持久化，便于"我的"页面展示
            if (result.healthy) {
                ServerConfigStore.saveFromResult(context, result)
            }
            splashState = if (result.healthy) {
                SplashState.Success(
                    backendVersion = result.backendVersion,
                    frontendVersion = result.frontendVersion
                )
            } else {
                SplashState.Failed(
                    message = result.message,
                    baseUrl = ApiClient.getBaseUrl(),
                    showUrl = true
                )
            }
        }
    }

    LaunchedEffect(Unit) {
        doCheck()
    }

    Crossfade(
        targetState = navigated,
        label = "splash-to-main"
    ) { showMain ->
        if (showMain) {
            AppNavigation()
        } else {
            SplashScreen(
                state = splashState,
                onRetry = doCheck,
                onSuccess = { navigated = true }
            )
        }
    }
}

/**
 * 主页导航图
 */
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
        // composable("drafts") {
        //     DraftsScreen(navController)
        // }

        // 个人资料
        composable("profile") {
            ProfileScreen(navController)
        }

        // 我的收藏
        composable("favorites") {
            FavoritesScreen(navController)
        }

        // 我的通知
        composable("notifications") {
            NotificationsScreen(navController)
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