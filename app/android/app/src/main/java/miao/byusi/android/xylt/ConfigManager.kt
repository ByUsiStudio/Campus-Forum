package miao.byusi.android.xylt

import android.content.Context
import android.util.Log
import org.json.JSONObject
import java.io.BufferedReader
import java.io.InputStreamReader

object ConfigManager {
    private const val TAG = "ConfigManager"
    private const val CONFIG_FILE = "config.json5"
    private const val DEFAULT_BASE_API = "https://xylt.cdifit.cn"
    
    private var appConfig: AppConfig? = null
    
    fun init(context: Context) {
        try {
            appConfig = loadConfigFromAssets(context)
            Log.d(TAG, "Config loaded: ${appConfig?.baseApi}")
        } catch (e: Exception) {
            Log.e(TAG, "Failed to load config, using default", e)
            appConfig = AppConfig(DEFAULT_BASE_API)
        }
    }
    
    fun getBaseApi(): String = appConfig?.baseApi ?: DEFAULT_BASE_API
    
    private fun loadConfigFromAssets(context: Context): AppConfig {
        val inputStream = context.assets.open(CONFIG_FILE)
        val reader = BufferedReader(InputStreamReader(inputStream))
        val content = reader.readText()
        reader.close()
        
        // 简单的JSON5解析 - 移除注释并使用标准JSON解析
        val cleanContent = cleanJson5(content)
        val json = JSONObject(cleanContent)
        
        return AppConfig(
            baseApi = json.optString("baseApi", DEFAULT_BASE_API)
        )
    }
    
    private fun cleanJson5(content: String): String {
        var result = content
        
        // 移除单行注释 // ...
        result = result.replace(Regex("//.*$", RegexOption.MULTILINE), "")
        
        // 移除多行注释 /* ... */
        result = result.replace(Regex("/\\*[\\s\\S]*?\\*/"), "")
        
        // 移除尾随逗号
        result = result.replace(Regex(",(\\s*[}\\]])"), "$1")
        
        return result
    }
}
