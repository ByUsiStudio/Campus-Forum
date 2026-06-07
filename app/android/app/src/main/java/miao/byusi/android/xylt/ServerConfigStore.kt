package miao.byusi.android.xylt

import android.content.Context
import android.content.SharedPreferences

/**
 * 启动校验后保存的服务器快照（后端版本 / API 地址 / 上次校验时间）。
 *
 * 用于：
 *  1. "我的"页面展示当前连接的后端版本、API 地址
 *  2. 客户端判断上次成功的健康状态（避免离线时仍展示"已连接"）
 */
object ServerConfigStore {
    private const val PREFS = "xylt_server_config"
    private const val KEY_API_URL = "api_url"
    private const val KEY_BACKEND_VERSION = "backend_version"
    private const val KEY_FRONTEND_VERSION = "frontend_version"
    private const val KEY_LAST_CHECK_AT = "last_check_at"

    private fun prefs(ctx: Context): SharedPreferences =
        ctx.getSharedPreferences(PREFS, Context.MODE_PRIVATE)

    fun saveSuccess(
        ctx: Context,
        apiUrl: String,
        backendVersion: String?,
        frontendVersion: String?,
        timestampMs: Long = System.currentTimeMillis()
    ) {
        prefs(ctx).edit().apply {
            putString(KEY_API_URL, apiUrl)
            putString(KEY_BACKEND_VERSION, backendVersion.orEmpty())
            putString(KEY_FRONTEND_VERSION, frontendVersion.orEmpty())
            putLong(KEY_LAST_CHECK_AT, timestampMs)
            apply()
        }
    }

    /** 在启动校验成功后调用 */
    fun saveFromResult(ctx: Context, result: HealthClient.HealthResult) {
        if (result.healthy) {
            saveSuccess(
                ctx = ctx,
                apiUrl = ApiClient.getBaseUrl(),
                backendVersion = result.backendVersion,
                frontendVersion = result.frontendVersion
            )
        }
    }

    fun apiUrl(ctx: Context): String =
        prefs(ctx).getString(KEY_API_URL, null) ?: ApiClient.getBaseUrl()

    fun backendVersion(ctx: Context): String? =
        prefs(ctx).getString(KEY_BACKEND_VERSION, null)?.takeIf { it.isNotBlank() }

    fun frontendVersion(ctx: Context): String? =
        prefs(ctx).getString(KEY_FRONTEND_VERSION, null)?.takeIf { it.isNotBlank() }

    fun lastCheckAt(ctx: Context): Long =
        prefs(ctx).getLong(KEY_LAST_CHECK_AT, 0L)
}
