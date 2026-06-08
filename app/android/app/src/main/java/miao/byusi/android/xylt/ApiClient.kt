package miao.byusi.android.xylt

import okhttp3.*
import okhttp3.MediaType.Companion.toMediaType
import okhttp3.RequestBody.Companion.toRequestBody
import java.io.IOException

object ApiClient {
    private val client = OkHttpClient.Builder().build()
    private val JSON_MEDIA_TYPE = "application/json; charset=utf-8".toMediaType()

    private var authToken: String = ""

    fun setAuthToken(token: String) {
        authToken = token
    }

    fun getAuthToken(): String = authToken

    /** 暴露当前 API 基础地址，供启动校验等场景使用 */
    fun getBaseUrl(): String = ConfigManager.getBaseApi()
    
    private fun buildRequest(url: String, body: RequestBody? = null, method: String = "GET"): Request {
        val builder = Request.Builder().url(url)
        if (authToken.isNotEmpty()) {
            builder.addHeader("Authorization", "Bearer $authToken")
        }
        when (method) {
            "POST" -> if (body != null) builder.post(body) else builder.post(RequestBody.create(null, ""))
            "PUT" -> if (body != null) builder.put(body) else builder.put(RequestBody.create(null, ""))
            "DELETE" -> if (body != null) builder.delete(body) else builder.delete()
            else -> builder.get()
        }
        return builder.build()
    }
    
    // ========== 认证接口 ==========
    
    // 用户登录
    fun login(username: String, password: String, callback: ApiCallback) {
        val url = "${getBaseUrl()}/api/auth/login"
        val jsonBody = org.json.JSONObject().apply {
            put("username", username)
            put("password", password)
        }.toString()
        val body = jsonBody.toRequestBody(JSON_MEDIA_TYPE)
        client.newCall(buildRequest(url, body, "POST")).enqueue(object : Callback {
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
    
    // 用户注册
    fun register(username: String, qqNumber: String, displayName: String, password: String, callback: ApiCallback) {
        val url = "${getBaseUrl()}/api/auth/register"
        val jsonBody = org.json.JSONObject().apply {
            put("username", username)
            put("qq_number", qqNumber)
            put("display_name", displayName)
            put("password", password)
        }.toString()
        val body = jsonBody.toRequestBody(JSON_MEDIA_TYPE)
        client.newCall(buildRequest(url, body, "POST")).enqueue(object : Callback {
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
    
    // ========== 文章接口 ==========
    
    // 获取文章列表
    fun getArticles(page: Int, pageSize: Int, categoryId: Int? = null, callback: ApiCallback) {
        val urlBuilder = StringBuilder("${getBaseUrl()}/api/articles?page=$page&page_size=$pageSize")
        if (categoryId != null) {
            urlBuilder.append("&category_id=$categoryId")
        }
        client.newCall(buildRequest(urlBuilder.toString())).enqueue(object : Callback {
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
    
    // 搜索文章
    fun searchArticles(keyword: String, page: Int, pageSize: Int, callback: ApiCallback) {
        val url = "${getBaseUrl()}/api/articles/search?keyword=$keyword&page=$page&page_size=$pageSize"
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
    fun getArticleDetail(articleId: Int, page: Int = 1, pageSize: Int = 20, callback: ApiCallback) {
        val url = "${getBaseUrl()}/api/articles/$articleId?page=$page&page_size=$pageSize"
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
    
    // 获取我的文章
    fun getMyArticles(page: Int, pageSize: Int, callback: ApiCallback) {
        val url = "${getBaseUrl()}/api/my/articles?page=$page&page_size=$pageSize"
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
    
    // 获取我的草稿
    fun getMyDrafts(page: Int, pageSize: Int, callback: ApiCallback) {
        val url = "${getBaseUrl()}/api/my/drafts?page=$page&page_size=$pageSize"
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
    
    // 创建文章
    fun createArticle(title: String, content: String, categoryId: Int, voiceUrl: String? = null, isAnonymous: Boolean = false, status: String = "published", callback: ApiCallback) {
        val url = "${getBaseUrl()}/api/articles"
        val jsonBody = org.json.JSONObject().apply {
            put("title", title)
            put("content", content)
            put("category_id", categoryId)
            put("is_anonymous", isAnonymous)
            put("status", status)
            if (voiceUrl != null) {
                put("voice_url", voiceUrl)
            }
        }.toString()
        val body = jsonBody.toRequestBody(JSON_MEDIA_TYPE)
        client.newCall(buildRequest(url, body, "POST")).enqueue(object : Callback {
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
    
    // 更新文章
    fun updateArticle(articleId: Int, title: String, content: String, categoryId: Int, voiceUrl: String? = null, isAnonymous: Boolean = false, callback: ApiCallback) {
        val url = "${getBaseUrl()}/api/articles/$articleId"
        val jsonBody = org.json.JSONObject().apply {
            put("title", title)
            put("content", content)
            put("category_id", categoryId)
            put("is_anonymous", isAnonymous)
            if (voiceUrl != null) {
                put("voice_url", voiceUrl)
            }
        }.toString()
        val body = jsonBody.toRequestBody(JSON_MEDIA_TYPE)
        client.newCall(buildRequest(url, body, "PUT")).enqueue(object : Callback {
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
        val url = "${getBaseUrl()}/api/articles/$articleId"
        client.newCall(buildRequest(url, null, "DELETE")).enqueue(object : Callback {
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
        val url = "${getBaseUrl()}/api/articles/$draftId/publish"
        val body = "".toRequestBody(JSON_MEDIA_TYPE)
        client.newCall(buildRequest(url, body, "POST")).enqueue(object : Callback {
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
        val url = "${getBaseUrl()}/api/articles/$articleId/like"
        val body = "".toRequestBody(JSON_MEDIA_TYPE)
        client.newCall(buildRequest(url, body, "POST")).enqueue(object : Callback {
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
    
    // 取消点赞
    fun unlikeArticle(articleId: Int, callback: ApiCallback) {
        val url = "${getBaseUrl()}/api/articles/$articleId/like"
        client.newCall(buildRequest(url, null, "DELETE")).enqueue(object : Callback {
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
    
    // 收藏文章
    fun favoriteArticle(articleId: Int, callback: ApiCallback) {
        val url = "${getBaseUrl()}/api/articles/$articleId/favorite"
        val body = "".toRequestBody(JSON_MEDIA_TYPE)
        client.newCall(buildRequest(url, body, "POST")).enqueue(object : Callback {
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
    
    // 取消收藏
    fun unfavoriteArticle(articleId: Int, callback: ApiCallback) {
        val url = "${getBaseUrl()}/api/articles/$articleId/favorite"
        client.newCall(buildRequest(url, null, "DELETE")).enqueue(object : Callback {
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
    
    // 检查收藏状态
    fun checkFavoriteStatus(articleId: Int, callback: ApiCallback) {
        val url = "${getBaseUrl()}/api/articles/$articleId/favorite/check"
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
    
    // 获取收藏列表
    fun getFavorites(page: Int, pageSize: Int, callback: ApiCallback) {
        val url = "${getBaseUrl()}/api/favorites?page=$page&page_size=$pageSize"
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
    
    // 分享文章
    fun shareArticle(articleId: Int, callback: ApiCallback) {
        val url = "${getBaseUrl()}/api/articles/$articleId/share"
        val body = "".toRequestBody(JSON_MEDIA_TYPE)
        client.newCall(buildRequest(url, body, "POST")).enqueue(object : Callback {
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
    
    // ========== 评论接口 ==========
    
    // 获取文章评论
    fun getComments(articleId: Int, callback: ApiCallback) {
        val url = "${getBaseUrl()}/api/articles/$articleId/comments"
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
    fun createComment(articleId: Int, content: String, parentId: Int? = null, isAnonymous: Boolean = false, callback: ApiCallback) {
        val url = "${getBaseUrl()}/api/articles/$articleId/comments"
        val jsonBody = org.json.JSONObject().apply {
            put("content", content)
            put("is_anonymous", isAnonymous)
            if (parentId != null) {
                put("parent_id", parentId)
            }
        }.toString()
        val body = jsonBody.toRequestBody(JSON_MEDIA_TYPE)
        client.newCall(buildRequest(url, body, "POST")).enqueue(object : Callback {
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
    
    // 删除评论
    fun deleteComment(commentId: Int, callback: ApiCallback) {
        val url = "${getBaseUrl()}/api/comments/$commentId"
        client.newCall(buildRequest(url, null, "DELETE")).enqueue(object : Callback {
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
    
    // 点赞评论
    fun likeComment(commentId: Int, callback: ApiCallback) {
        val url = "${getBaseUrl()}/api/comments/$commentId/like"
        val body = "".toRequestBody(JSON_MEDIA_TYPE)
        client.newCall(buildRequest(url, body, "POST")).enqueue(object : Callback {
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
    
    // 取消评论点赞
    fun unlikeComment(commentId: Int, callback: ApiCallback) {
        val url = "${getBaseUrl()}/api/comments/$commentId/like"
        client.newCall(buildRequest(url, null, "DELETE")).enqueue(object : Callback {
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
    
    // ========== 关注接口 ==========
    
    // 关注用户
    fun followUser(userId: Int, callback: ApiCallback) {
        val url = "${getBaseUrl()}/api/follow/$userId"
        val body = "".toRequestBody(JSON_MEDIA_TYPE)
        client.newCall(buildRequest(url, body, "POST")).enqueue(object : Callback {
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
    
    // 取消关注
    fun unfollowUser(userId: Int, callback: ApiCallback) {
        val url = "${getBaseUrl()}/api/follow/$userId"
        client.newCall(buildRequest(url, null, "DELETE")).enqueue(object : Callback {
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
    
    // 获取关注列表
    fun getFollowing(callback: ApiCallback) {
        val url = "${getBaseUrl()}/api/following"
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
    
    // 获取粉丝列表
    fun getFollowers(callback: ApiCallback) {
        val url = "${getBaseUrl()}/api/followers"
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
    
    // 检查关注状态
    fun checkFollowStatus(userId: Int, callback: ApiCallback) {
        val url = "${getBaseUrl()}/api/follow/status/$userId"
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
    
    // ========== 分区接口 ==========
    
    // 获取分区列表
    fun getCategories(callback: ApiCallback) {
        val url = "${getBaseUrl()}/api/categories"
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
    
    // ========== 用户接口 ==========
    
    // 获取个人资料
    fun getProfile(callback: ApiCallback) {
        val url = "${getBaseUrl()}/api/profile"
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
    
    // 更新个人资料
    fun updateProfile(displayName: String, signature: String, callback: ApiCallback) {
        val url = "${getBaseUrl()}/api/profile"
        val jsonBody = org.json.JSONObject().apply {
            put("display_name", displayName)
            put("signature", signature)
        }.toString()
        val body = jsonBody.toRequestBody(JSON_MEDIA_TYPE)
        client.newCall(buildRequest(url, body, "PUT")).enqueue(object : Callback {
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
    
    // 获取用户公开信息
    fun getUserProfile(userId: Int, callback: ApiCallback) {
        val url = "${getBaseUrl()}/api/users/$userId"
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
    
    // 获取用户文章列表
    fun getUserArticles(userId: Int, page: Int, pageSize: Int, callback: ApiCallback) {
        val url = "${getBaseUrl()}/api/users/$userId/articles?page=$page&page_size=$pageSize"
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
    
    // ========== 通知接口 ==========
    
    // 获取通知列表
    fun getNotifications(page: Int, pageSize: Int, callback: ApiCallback) {
        val url = "${getBaseUrl()}/api/user-notifications?page=$page&page_size=$pageSize"
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
    
    // 获取未读通知数量
    fun getUnreadNotificationCount(callback: ApiCallback) {
        val url = "${getBaseUrl()}/api/notifications/unread-count"
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
    
    // 标记通知为已读
    fun markNotificationAsRead(notificationId: Int, callback: ApiCallback) {
        val url = "${getBaseUrl()}/api/user-notifications/$notificationId/read"
        val body = "".toRequestBody(JSON_MEDIA_TYPE)
        client.newCall(buildRequest(url, body, "POST")).enqueue(object : Callback {
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
    
    // 标记所有通知为已读
    fun markAllNotificationsAsRead(callback: ApiCallback) {
        val url = "${getBaseUrl()}/api/user-notifications/read-all"
        val body = "".toRequestBody(JSON_MEDIA_TYPE)
        client.newCall(buildRequest(url, body, "POST")).enqueue(object : Callback {
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
    
    // ========== 举报接口 ==========
    
    // 提交举报
    fun submitReport(targetType: String, targetId: Int, reason: String, description: String, callback: ApiCallback) {
        val url = "${getBaseUrl()}/api/reports"
        val jsonBody = org.json.JSONObject().apply {
            put("target_type", targetType)
            put("target_id", targetId)
            put("reason", reason)
            put("description", description)
        }.toString()
        val body = jsonBody.toRequestBody(JSON_MEDIA_TYPE)
        client.newCall(buildRequest(url, body, "POST")).enqueue(object : Callback {
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