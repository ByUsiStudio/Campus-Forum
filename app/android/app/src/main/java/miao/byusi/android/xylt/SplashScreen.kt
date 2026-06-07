package miao.byusi.android.xylt

import androidx.compose.animation.core.animateFloatAsState
import androidx.compose.animation.core.tween
import androidx.compose.foundation.Image
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableFloatStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.alpha
import androidx.compose.ui.draw.clip
import androidx.compose.ui.draw.scale
import androidx.compose.ui.graphics.Brush
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import kotlinx.coroutines.delay
import top.yukonga.miuix.kmp.basic.Button
import top.yukonga.miuix.kmp.basic.Card
import top.yukonga.miuix.kmp.basic.CircularProgressIndicator
import top.yukonga.miuix.kmp.basic.Surface
import top.yukonga.miuix.kmp.basic.Text
import top.yukonga.miuix.kmp.theme.MiuixTheme

/**
 * 启动校验页面（Miuix 风格）
 *
 * 状态机由外层 [SplashState] 决定：
 *  - Checking: 显示加载圈 + 提示
 *  - Success : 短暂显示对号后通过 [onSuccess] 跳转
 *  - Failed  : 在 [Card] 卡片中显示错误信息 + [Button] 重试按钮
 *
 * 整体外观遵循 Miuix 设计语言：圆角卡片、主题色按钮、状态栏沉浸。
 */
@Composable
fun SplashScreen(
    state: SplashState,
    onRetry: () -> Unit,
    onSuccess: () -> Unit
) {
    // 启动 logo 呼吸缩放动画
    var pulse by remember { mutableFloatStateOf(0f) }
    val logoScale by animateFloatAsState(
        targetValue = 0.95f + 0.05f * pulse,
        animationSpec = tween(durationMillis = 1200),
        label = "logo-pulse"
    )

    LaunchedEffect(Unit) {
        while (true) {
            pulse = 1f
            delay(1200)
            pulse = 0f
            delay(1200)
        }
    }

    if (state is SplashState.Success) {
        LaunchedEffect(Unit) {
            delay(500)
            onSuccess()
        }
    }

    Surface(
        modifier = Modifier.fillMaxSize()
    ) {
        Box(
            modifier = Modifier
                .fillMaxSize()
                .background(
                    Brush.verticalGradient(
                        colors = listOf(
                            MiuixTheme.colorScheme.primary,
                            MiuixTheme.colorScheme.primaryContainer
                        )
                    )
                ),
            contentAlignment = Alignment.Center
        ) {
            Column(
                horizontalAlignment = Alignment.CenterHorizontally,
                verticalArrangement = Arrangement.Center,
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(horizontal = 32.dp)
            ) {
                // 应用图标（带轻微呼吸动画）
                Image(
                    painter = painterResource(id = R.mipmap.ic_launcher),
                    contentDescription = "App Logo",
                    modifier = Modifier
                        .size(112.dp)
                        .scale(logoScale)
                        .clip(RoundedCornerShape(28.dp))
                )

                Spacer(modifier = Modifier.height(20.dp))

                Text(
                    text = "校园论坛",
                    color = MiuixTheme.colorScheme.onPrimary,
                    fontSize = 26.sp,
                    fontWeight = FontWeight.Bold
                )

                Spacer(modifier = Modifier.height(6.dp))

                Text(
                    text = "分享与交流的社区",
                    color = MiuixTheme.colorScheme.onPrimary.copy(alpha = 0.85f),
                    fontSize = 13.sp
                )

                Spacer(modifier = Modifier.height(40.dp))

                when (state) {
                    is SplashState.Checking -> CheckingBlock()
                    is SplashState.Success -> SuccessBlock(state)
                    is SplashState.Failed -> FailedBlock(state, onRetry)
                }
            }
        }
    }
}

@Composable
private fun CheckingBlock() {
    CircularProgressIndicator(
        modifier = Modifier.size(32.dp)
    )
    Spacer(modifier = Modifier.height(14.dp))
    Text(
        text = "正在校验服务器…",
        color = MiuixTheme.colorScheme.onPrimary,
        fontSize = 14.sp
    )
}

@Composable
private fun SuccessBlock(state: SplashState.Success) {
    StatusDot(
        color = MiuixTheme.colorScheme.secondary,
        modifier = Modifier.size(18.dp)
    )
    Spacer(modifier = Modifier.height(14.dp))
    Text(
        text = "服务器在线，正在进入…",
        color = MiuixTheme.colorScheme.onPrimary,
        fontSize = 14.sp
    )
    val ver = state.backendVersion
    if (!ver.isNullOrBlank()) {
        Spacer(modifier = Modifier.height(4.dp))
        Text(
            text = "后端版本: $ver",
            color = MiuixTheme.colorScheme.onPrimary.copy(alpha = 0.7f),
            fontSize = 12.sp
        )
    }
}

@Composable
private fun FailedBlock(state: SplashState.Failed, onRetry: () -> Unit) {
    // 使用 Miuix Card 承载错误信息
    Card(
        modifier = Modifier
            .fillMaxWidth()
            .padding(horizontal = 4.dp)
    ) {
        Column(
            horizontalAlignment = Alignment.CenterHorizontally,
            modifier = Modifier
                .fillMaxWidth()
                .padding(horizontal = 20.dp, vertical = 18.dp)
        ) {
            StatusDot(
                color = Color(0xFFFF5252),
                modifier = Modifier.size(16.dp)
            )
            Spacer(modifier = Modifier.height(10.dp))
            Text(
                text = state.message,
                color = MiuixTheme.colorScheme.onSurface,
                fontSize = 14.sp,
                textAlign = TextAlign.Center
            )
            if (state.showUrl) {
                Spacer(modifier = Modifier.height(6.dp))
                Text(
                    text = state.baseUrl,
                    color = MiuixTheme.colorScheme.onSurface.copy(alpha = 0.65f),
                    fontSize = 12.sp,
                    textAlign = TextAlign.Center
                )
            }
            Spacer(modifier = Modifier.height(16.dp))
            // Miuix 风格主按钮：跟随主题色，带圆角
            Button(
                onClick = onRetry,
                modifier = Modifier
                    .fillMaxWidth()
                    .height(44.dp)
            ) {
                Text(
                    text = "重  试",
                    color = MiuixTheme.colorScheme.onPrimary,
                    fontWeight = FontWeight.SemiBold
                )
            }
        }
    }
}

@Composable
private fun StatusDot(color: Color, modifier: Modifier = Modifier) {
    Box(
        modifier = modifier
            .clip(RoundedCornerShape(50))
            .background(color)
    )
}

/** 启动页面状态 */
sealed class SplashState {
    /** 正在校验中 */
    data object Checking : SplashState()
    /** 校验成功 */
    data class Success(
        val backendVersion: String?,
        val frontendVersion: String?
    ) : SplashState()
    /** 校验失败 */
    data class Failed(
        val message: String,
        val baseUrl: String,
        val showUrl: Boolean = true
    ) : SplashState()
}
