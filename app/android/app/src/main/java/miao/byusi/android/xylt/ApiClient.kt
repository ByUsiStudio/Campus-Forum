package miao.byusi.android.xylt

import okhttp3.*
import okhttp3.MediaType.Companion.toMediaType
import okhttp3.RequestBody.Companion.toRequestBody
import java.io.IOException

object ApiClient {
    private const val BASE_URL = "https://xylt.cdifit.cn"
    private val client = OkHttpClient.Builder().build()
    private const val JSON_MEDIA_TYPE = "application/json; charset=utf-8".toMediaType()
    
    private var authToken: String = ""
    
    fun setAuthToken(token: String) {
        authToken = token
    }
    
    fun getAuthToken(): String = authToken
    
    private fun buildRequest(url: String, body: RequestBody? = null): Request {
        val builder = Request.Builder().url(url)
        if (authToken.isNotEmpty()) {
            builder.addHeader("Authorization", "Bearer $authToken")
        }
        if (body != null) {
            builder.post(body)
        }
        return builder.build()
    }
    
    // 获取文章列表
    fun getArticles(page: Int, limit: Int, callback: ApiCallback) {
        val url = "$BASE_URL/api/articles?page=$page&limit=$limit"
        client.newCall(buildRequest(url)).enqueue(object : Callback {
            override fun onFailure(call: Call, e: IOException) {
                callback.onError(e.message ?: "网络错误")
            }
            override fun onResponse(call: Call, response: Response) {
                if (response.isSuccessful) {
                    callback.onSuccess(response.body?.string() ?: "")
                } else {
                    callback.onError("请求失败: ${response.code}")
                }
            }
        })
    }
    
    // 获取文章详情
    fun getArticleDetail(articleId: Int, callback: ApiCallback) {
        val url = "$BASE_URL/api/articles/$articleId"
        client.newCall(buildRequest(url)).enqueue(object : Callback {
            override fun onFailure(call: Call, e: IOException) {
                callback.onError(e.message ?: "网络错误")
            }
            override fun onResponse(call: Call, response: Response) {
                if (response.isSuccessful) {
                    callback.onSuccess(response.body?.string() ?: "")
                } else {
                    callback.onError("请求失败: ${response.code}")
                }
            }
        })
    }
    
    // 点赞文章
    fun likeArticle(articleId: Int, callback: ApiCallback) {
        val url = "$BASE_URL/api/articles/$articleId/like"
        val body = "".toRequestBody(JSON_MEDIA_TYPE)
        client.newCall(buildRequest(url, body)).enqueue(object : Callback {
            override fun onFailure(call: Call, e: IOException) {
                callback.onError(e.message ?: "网络错误")
            }
            override fun onResponse(call: Call, response: Response) {
                if (response.isSuccessful) {
                    callback.onSuccess(response.body?.string() ?: "")
                } else {
                    callback.onError("请求失败: ${response.code}")
                }
            }
        })
    }
    
    // 获取评论
    fun getComments(articleId: Int, page: Int, limit: Int, callback: ApiCallback) {
        val url = "$BASE_URL/api/articles/$articleId/comments?page=$page&limit=$limit"
        client.newCall(buildRequest(url)).enqueue(object : Callback {
            override fun onFailure(call: Call, e: IOException) {
                callback.onError(e.message ?: "网络错误")
            }
            override fun onResponse(call: Call, response: Response) {
                if (response.isSuccessful) {
                    callback.onSuccess(response.body?.string() ?: "")
                } else {
                    callback.onError("请求失败: ${response.code}")
                }
            }
        })
    }
    
    // 创建评论
    fun createComment(articleId: Int, content: String, parentId: Int?, callback: ApiCallback) {
        val url = "$BASE_URL/api/articles/$articleId/comments"
        val jsonBody = org.json.JSONObject().apply {
            put("content", content)
            if (parentId != null) {
                put("parent_id", parentId)
            }
        }.toString()
        val body = jsonBody.toRequestBody(JSON_MEDIA_TYPE)
        client.newCall(buildRequest(url, body)).enqueue(object : Callback {
            override fun onFailure(call: Call, e: IOException) {
                callback.onError(e.message ?: "网络错误")
            }
            override fun onResponse(call: Call, response: Response) {
                if (response.isSuccessful) {
                    callback.onSuccess(response.body?.string() ?: "")
                } else {
                    callback.onError("请求失败: ${response.code}")
                }
            }
        })
    }
    
    // 获取草稿列表
    fun getMyDrafts(page: Int, limit: Int, callback: ApiCallback) {
        val url = "$BASE_URL/api/articles/drafts?page=$page&limit=$limit"
        client.newCall(buildRequest(url)).enqueue(object : Callback {
            override fun onFailure(call: Call, e: IOException) {
                callback.onError(e.message ?: "网络错误")
            }
            override fun onResponse(call: Call, response: Response) {
                if (response.isSuccessful) {
                    callback.onSuccess(response.body?.string() ?: "")
                } else {
                    callback.onError("请求失败: ${response.code}")
                }
            }
        })
    }
    
    // 发布草稿
    fun publishDraft(draftId: Int, callback: ApiCallback) {
        val url = "$BASE_URL/api/articles/$draftId/publish"
        val body = "".toRequestBody(JSON_MEDIA_TYPE)
        client.newCall(buildRequest(url, body)).enqueue(object : Callback {
            override fun onFailure(call: Call, e: IOException) {
                callback.onError(e.message ?: "网络错误")
            }
            override fun onResponse(call: Call, response: Response) {
                if (response.isSuccessful) {
                    callback.onSuccess(response.body?.string() ?: "")
                } else {
                    callback.onError("请求失败: ${response.code}")
                }
            }
        })
    }
    
    // 删除文章
    fun deleteArticle(articleId: Int, callback: ApiCallback) {
        val url = "$BASE_URL/api/articles/$articleId"
        val request = Request.Builder().url(url).delete().apply {
            if (authToken.isNotEmpty()) {
                addHeader("Authorization", "Bearer $authToken")
            }
        }.build()
        client.newCall(request).enqueue(object : Callback {
            override fun onFailure(call: Call, e: IOException) {
                callback.onError(e.message ?: "网络错误")
            }
            override fun onResponse(call: Call, response: Response) {
                if (response.isSuccessful) {
                    callback.onSuccess(response.body?.string() ?: "")
                } else {
                    callback.onError("请求失败: ${response.code}")
                }
            }
        })
    }
    
    // 获取用户资料
    fun getProfile(callback: ApiCallback) {
        val url = "$BASE_URL/api/user/profile"
        client.newCall(buildRequest(url)).enqueue(object : Callback {
            override fun onFailure(call: Call, e: IOException) {
                callback.onError(e.message ?: "网络错误")
            }
            override fun onResponse(call: Call, response: Response) {
                if (response.isSuccessful) {
                    callback.onSuccess(response.body?.string() ?: "")
                } else {
                    callback.onError("请求失败: ${response.code}")
                }
            }
        })
    }
    
    // 更新用户资料
    fun updateProfile(displayName: String, signature: String, callback: ApiCallback) {
        val url = "$BASE_URL/api/user/profile"
        val jsonBody = org.json.JSONObject().apply {
            put("display_name", displayName)
            put("signature", signature)
        }.toString()
        val body = jsonBody.toRequestBody(JSON_MEDIA_TYPE)
        client.newCall(buildRequest(url, body)).enqueue(object : Callback {
            override fun onFailure(call: Call, e: IOException) {
                callback.onError(e.message ?: "网络错误")
            }
            override fun onResponse(call: Call, response: Response) {
                if (response.isSuccessful) {
                    callback.onSuccess(response.body?.string() ?: "")
                } else {
                    callback.onError("请求失败: ${response.code}")
                }
            }
        })
    }
    
    // 登录
    fun login(username: String, password: String, callback: ApiCallback) {
        val url = "$BASE_URL/api/auth/login"
        val jsonBody = org.json.JSONObject().apply {
            put("username", username)
            put("password", password)
        }.toString()
        val body = jsonBody.toRequestBody(JSON_MEDIA_TYPE)
        client.newCall(buildRequest(url, body)).enqueue(object : Callback {
            override fun onFailure(call: Call, e: IOException) {
                callback.onError(e.message ?: "网络错误")
            }
            override fun onResponse(call: Call, response: Response) {
                if (response.isSuccessful) {
                    callback.onSuccess(response.body?.string() ?: "")
                } else {
                    callback.onError("请求失败: ${response.code}")
                }
            }
        })
    }
    
    // 注册
    fun register(username: String, qqNumber: String, displayName: String, password: String, callback: ApiCallback) {
        val url = "$BASE_URL/api/auth/register"
        val jsonBody = org.json.JSONObject().apply {
            put("username", username)
            put("qq_number", qqNumber)
            put("display_name", displayName)
            put("password", password)
        }.toString()
        val body = jsonBody.toRequestBody(JSON_MEDIA_TYPE)
        client.newCall(buildRequest(url, body)).enqueue(object : Callback {
            override fun onFailure(call: Call, e: IOException) {
                callback.onError(e.message ?: "网络错误")
            }
            override fun onResponse(call: Call, response: Response) {
                if (response.isSuccessful) {
                    callback.onSuccess(response.body?.string() ?: "")
                } else {
                    callback.onError("请求失败: ${response.code}")
                }
            }
        })
    }
    
    // 提交举报
    fun submitReport(targetType: String, targetId: Int, reason: String, description: String, callback: ApiCallback) {
        val url = "$BASE_URL/api/reports"
        val jsonBody = org.json.JSONObject().apply {
            put("target_type", targetType)
            put("target_id", targetId)
            put("reason", reason)
            put("description", description)
        }.toString()
        val body = jsonBody.toRequestBody(JSON_MEDIA_TYPE)
        client.newCall(buildRequest(url, body)).enqueue(object : Callback {
            override fun onFailure(call: Call, e: IOException) {
                callback.onError(e.message ?: "网络错误")
            }
            override fun onResponse(call: Call, response: Response) {
                if (response.isSuccessful) {
                    callback.onSuccess(response.body?.string() ?: "")
                } else {
                    callback.onError("请求失败: ${response.code}")
                }
            }
        })
    }
}

interface ApiCallback {
    fun onSuccess(response: String)
    fun onError(error: String)
}