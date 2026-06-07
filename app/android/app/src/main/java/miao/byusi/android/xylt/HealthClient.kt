package miao.byusi.android.xylt

import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import okhttp3.OkHttpClient
import okhttp3.Request
import org.json.JSONObject
import java.io.IOException
import java.util.concurrent.TimeUnit

/**
 * 启动时服务器健康校验客户端
 *
 * 通过访问公开的 /api/version 接口来判断服务器是否正常在线：
 *   - HTTP 200 + status=ok  : 服务器在线
 *   - HTTP 503              : 服务降级（数据库异常等）
 *   - 连接失败 / 超时 / 解析错误 : 服务器离线
 *
 * 该接口为后端已有公开接口，无需鉴权，无需修改后端。
 */
object HealthClient {

    private const val HEALTH_PATH = "/api/version"

    /** 健康校验超时时间（秒） */
    private const val TIMEOUT_SECONDS = 5L

    private val client: OkHttpClient = OkHttpClient.Builder()
        .connectTimeout(TIMEOUT_SECONDS, TimeUnit.SECONDS)
        .readTimeout(TIMEOUT_SECONDS, TimeUnit.SECONDS)
        .writeTimeout(TIMEOUT_SECONDS, TimeUnit.SECONDS)
        .callTimeout(TIMEOUT_SECONDS + 1, TimeUnit.SECONDS)
        .build()

    /**
     * 健康校验结果
     *
     * @param online 服务器是否在线
     * @param healthy 服务器是否健康（online + 业务返回 ok）
     * @param message 给用户展示的提示信息
     * @param backendVersion 后端版本（若可获取）
     * @param frontendVersion 前端版本（若可获取）
     */
    data class HealthResult(
        val online: Boolean,
        val healthy: Boolean,
        val message: String,
        val backendVersion: String? = null,
        val frontendVersion: String? = null
    )

    /**
     * 同步执行健康校验（协程版本，请在协程中调用）
     */
    suspend fun check(baseUrl: String = ApiClient.getBaseUrl()): HealthResult =
        withContext(Dispatchers.IO) {
            val url = baseUrl.trimEnd('/') + HEALTH_PATH
            val request = Request.Builder()
                .url(url)
                .get()
                .build()
            try {
                client.newCall(request).execute().use { response ->
                    val body = response.body?.string().orEmpty()
                    if (!response.isSuccessful) {
                        return@withContext HealthResult(
                            online = true,
                            healthy = false,
                            message = "服务器返回异常: HTTP ${response.code}"
                        )
                    }
                    val json = runCatching { JSONObject(body) }.getOrNull()
                    val backend = json?.optJSONObject("backend")?.optString("version")
                    val frontend = json?.optJSONObject("frontend")?.optString("version")
                    if (backend != null || frontend != null) {
                        HealthResult(
                            online = true,
                            healthy = true,
                            message = "服务器在线",
                            backendVersion = backend,
                            frontendVersion = frontend
                        )
                    } else {
                        HealthResult(
                            online = true,
                            healthy = false,
                            message = "服务器响应格式异常"
                        )
                    }
                }
            } catch (e: IOException) {
                HealthResult(
                    online = false,
                    healthy = false,
                    message = friendlyError(e)
                )
            } catch (e: Exception) {
                HealthResult(
                    online = false,
                    healthy = false,
                    message = friendlyError(e)
                )
            }
        }

    private fun friendlyError(e: Throwable): String {
        val raw = e.message ?: e.javaClass.simpleName
        return when {
            raw.contains("Unable to resolve host", ignoreCase = true) ->
                "无法解析服务器地址，请检查网络或 API 地址配置"
            raw.contains("Failed to connect", ignoreCase = true) ||
                raw.contains("ConnectException", ignoreCase = true) ->
                "无法连接服务器，请确认服务已启动"
            raw.contains("timeout", ignoreCase = true) ||
                raw.contains("Timeout", ignoreCase = true) ->
                "连接服务器超时，请稍后重试"
            raw.contains("SSL", ignoreCase = true) ||
                raw.contains("Trust", ignoreCase = true) ->
                "安全连接失败: $raw"
            else -> "无法连接服务器: $raw"
        }
    }
}
