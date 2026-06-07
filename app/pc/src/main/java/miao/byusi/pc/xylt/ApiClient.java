package miao.byusi.pc.xylt;

import okhttp3.*;
import org.json.JSONObject;

import java.io.IOException;
import java.util.concurrent.TimeUnit;

public class ApiClient {

    private static final String BASE_URL = "https://xylt.cdifit.cn/api";
    private static final MediaType JSON = MediaType.parse("application/json; charset=utf-8");
    private static String authToken = "";

    private static final OkHttpClient client = new OkHttpClient.Builder()
            .connectTimeout(30, TimeUnit.SECONDS)
            .readTimeout(30, TimeUnit.SECONDS)
            .writeTimeout(30, TimeUnit.SECONDS)
            .build();

    public static void setAuthToken(String token) {
        authToken = token;
    }

    public static String getAuthToken() {
        return authToken;
    }

    // ==================== 认证接口 ====================

    public static void login(String username, String password, ApiCallback callback) {
        String json = "{\"username\":\"" + username + "\",\"password\":\"" + password + "\"}";
        RequestBody body = RequestBody.create(json, JSON);
        Request request = new Request.Builder()
                .url(BASE_URL + "/auth/login")
                .post(body)
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("登录失败");
                }
            }
        });
    }

    public static void register(String username, String qqNumber, String displayName, String password, ApiCallback callback) {
        String json = "{\"username\":\"" + username + "\",\"qq_number\":\"" + qqNumber + "\",\"display_name\":\"" + displayName + "\",\"password\":\"" + password + "\"}";
        RequestBody body = RequestBody.create(json, JSON);
        Request request = new Request.Builder()
                .url(BASE_URL + "/auth/register")
                .post(body)
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("注册失败");
                }
            }
        });
    }

    public static void sendPasswordResetCode(String qqNumber, ApiCallback callback) {
        String json = "{\"qq_number\":\"" + qqNumber + "\"}";
        RequestBody body = RequestBody.create(json, JSON);
        Request request = new Request.Builder()
                .url(BASE_URL + "/password/reset-code")
                .post(body)
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("发送验证码失败");
                }
            }
        });
    }

    public static void resetPassword(String qqNumber, String code, String identifier, String password, ApiCallback callback) {
        String json = "{\"qq_number\":\"" + qqNumber + "\",\"code\":\"" + code + "\",\"identifier\":\"" + identifier + "\",\"password\":\"" + password + "\"}";
        RequestBody body = RequestBody.create(json, JSON);
        Request request = new Request.Builder()
                .url(BASE_URL + "/password/reset")
                .post(body)
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("密码重置失败");
                }
            }
        });
    }

    public static void getProfile(ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/profile")
                .addHeader("Authorization", "Bearer " + authToken)
                .get()
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("获取用户信息失败");
                }
            }
        });
    }

    public static void updateProfile(String displayName, String signature, ApiCallback callback) {
        String json = "{\"display_name\":\"" + displayName + "\",\"signature\":\"" + signature + "\"}";
        RequestBody body = RequestBody.create(json, JSON);
        Request request = new Request.Builder()
                .url(BASE_URL + "/profile")
                .addHeader("Authorization", "Bearer " + authToken)
                .put(body)
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("更新用户信息失败");
                }
            }
        });
    }

    // ==================== 文章接口 ====================

    public static void getArticles(int page, int pageSize, ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/articles?page=" + page + "&page_size=" + pageSize)
                .get()
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("获取文章列表失败");
                }
            }
        });
    }

    public static void getArticleDetail(int articleId, ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/articles/" + articleId)
                .get()
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("获取文章详情失败");
                }
            }
        });
    }

    public static void searchArticles(String keyword, int page, int pageSize, ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/articles/search?q=" + keyword + "&page=" + page + "&page_size=" + pageSize)
                .get()
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("搜索文章失败");
                }
            }
        });
    }

    public static void createArticle(String title, String content, int categoryId, boolean isAnonymous, ApiCallback callback) {
        String json = "{\"title\":\"" + title + "\",\"content\":\"" + content + "\",\"category_id\":" + categoryId + ",\"is_anonymous\":" + isAnonymous + "}";
        RequestBody body = RequestBody.create(json, JSON);
        Request request = new Request.Builder()
                .url(BASE_URL + "/articles")
                .addHeader("Authorization", "Bearer " + authToken)
                .post(body)
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("发布文章失败");
                }
            }
        });
    }

    public static void updateArticle(int articleId, String title, String content, int categoryId, boolean isAnonymous, ApiCallback callback) {
        String json = "{\"title\":\"" + title + "\",\"content\":\"" + content + "\",\"category_id\":" + categoryId + ",\"is_anonymous\":" + isAnonymous + "}";
        RequestBody body = RequestBody.create(json, JSON);
        Request request = new Request.Builder()
                .url(BASE_URL + "/articles/" + articleId)
                .addHeader("Authorization", "Bearer " + authToken)
                .put(body)
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("更新文章失败");
                }
            }
        });
    }

    public static void deleteArticle(int articleId, ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/articles/" + articleId)
                .addHeader("Authorization", "Bearer " + authToken)
                .delete()
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("删除文章失败");
                }
            }
        });
    }

    public static void shareArticle(int articleId, ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/articles/" + articleId + "/share")
                .addHeader("Authorization", "Bearer " + authToken)
                .post(RequestBody.create("", JSON))
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("分享文章失败");
                }
            }
        });
    }

    public static void pinArticle(int articleId, ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/articles/" + articleId + "/pin")
                .addHeader("Authorization", "Bearer " + authToken)
                .post(RequestBody.create("", JSON))
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("置顶文章失败");
                }
            }
        });
    }

    public static void unpinArticle(int articleId, ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/articles/" + articleId + "/pin")
                .addHeader("Authorization", "Bearer " + authToken)
                .delete()
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("取消置顶失败");
                }
            }
        });
    }

    public static void restoreArticle(int articleId, ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/articles/" + articleId + "/restore")
                .addHeader("Authorization", "Bearer " + authToken)
                .post(RequestBody.create("", JSON))
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("恢复文章失败");
                }
            }
        });
    }

    public static void getMyArticles(int page, int pageSize, ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/my/articles?page=" + page + "&page_size=" + pageSize)
                .addHeader("Authorization", "Bearer " + authToken)
                .get()
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("获取我的文章失败");
                }
            }
        });
    }

    public static void getMyDrafts(int page, int pageSize, ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/my/drafts?page=" + page + "&page_size=" + pageSize)
                .addHeader("Authorization", "Bearer " + authToken)
                .get()
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("获取草稿失败");
                }
            }
        });
    }

    public static void publishDraft(int articleId, ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/articles/" + articleId + "/publish")
                .addHeader("Authorization", "Bearer " + authToken)
                .post(RequestBody.create("", JSON))
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("发布草稿失败");
                }
            }
        });
    }

    public static void likeArticle(int articleId, ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/articles/" + articleId + "/like")
                .addHeader("Authorization", "Bearer " + authToken)
                .post(RequestBody.create("", JSON))
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("点赞失败");
                }
            }
        });
    }

    public static void unlikeArticle(int articleId, ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/articles/" + articleId + "/like")
                .addHeader("Authorization", "Bearer " + authToken)
                .delete()
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("取消点赞失败");
                }
            }
        });
    }

    public static void addFavorite(int articleId, ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/articles/" + articleId + "/favorite")
                .addHeader("Authorization", "Bearer " + authToken)
                .post(RequestBody.create("", JSON))
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("收藏失败");
                }
            }
        });
    }

    public static void removeFavorite(int articleId, ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/articles/" + articleId + "/favorite")
                .addHeader("Authorization", "Bearer " + authToken)
                .delete()
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("取消收藏失败");
                }
            }
        });
    }

    public static void checkFavorite(int articleId, ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/articles/" + articleId + "/favorite/check")
                .addHeader("Authorization", "Bearer " + authToken)
                .get()
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("检查收藏状态失败");
                }
            }
        });
    }

    public static void getFavorites(int page, int pageSize, ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/favorites?page=" + page + "&page_size=" + pageSize)
                .addHeader("Authorization", "Bearer " + authToken)
                .get()
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("获取收藏列表失败");
                }
            }
        });
    }

    // ==================== 分类接口 ====================

    public static void getCategories(ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/categories")
                .get()
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("获取分类失败");
                }
            }
        });
    }

    // ==================== 评论接口 ====================

    public static void getComments(int articleId, int page, int pageSize, ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/articles/" + articleId + "/comments?page=" + page + "&page_size=" + pageSize)
                .get()
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("获取评论失败");
                }
            }
        });
    }

    public static void createComment(int articleId, String content, Integer parentId, ApiCallback callback) {
        String json = parentId != null
                ? "{\"content\":\"" + content + "\",\"parent_id\":" + parentId + "}"
                : "{\"content\":\"" + content + "\"}";
        RequestBody body = RequestBody.create(json, JSON);
        Request request = new Request.Builder()
                .url(BASE_URL + "/articles/" + articleId + "/comments")
                .addHeader("Authorization", "Bearer " + authToken)
                .post(body)
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("发表评论失败");
                }
            }
        });
    }

    public static void deleteComment(int commentId, ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/comments/" + commentId)
                .addHeader("Authorization", "Bearer " + authToken)
                .delete()
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("删除评论失败");
                }
            }
        });
    }

    public static void likeComment(int commentId, ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/comments/" + commentId + "/like")
                .addHeader("Authorization", "Bearer " + authToken)
                .post(RequestBody.create("", JSON))
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("点赞评论失败");
                }
            }
        });
    }

    public static void unlikeComment(int commentId, ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/comments/" + commentId + "/like")
                .addHeader("Authorization", "Bearer " + authToken)
                .delete()
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("取消点赞失败");
                }
            }
        });
    }

    // ==================== 举报接口 ====================

    public static void submitReport(String targetType, int targetId, String reason, String description, ApiCallback callback) {
        String json = "{\"target_type\":\"" + targetType + "\",\"target_id\":" + targetId + ",\"reason\":\"" + reason + "\",\"description\":\"" + description + "\"}";
        RequestBody body = RequestBody.create(json, JSON);
        Request request = new Request.Builder()
                .url(BASE_URL + "/reports")
                .addHeader("Authorization", "Bearer " + authToken)
                .post(body)
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("提交举报失败");
                }
            }
        });
    }

    // ==================== 关注接口 ====================

    public static void followUser(int userId, ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/follow/" + userId)
                .addHeader("Authorization", "Bearer " + authToken)
                .post(RequestBody.create("", JSON))
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("关注失败");
                }
            }
        });
    }

    public static void unfollowUser(int userId, ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/follow/" + userId)
                .addHeader("Authorization", "Bearer " + authToken)
                .delete()
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("取消关注失败");
                }
            }
        });
    }

    public static void checkFollowStatus(int userId, ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/follow/status/" + userId)
                .addHeader("Authorization", "Bearer " + authToken)
                .get()
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("检查关注状态失败");
                }
            }
        });
    }

    public static void getFollowingList(int page, int pageSize, ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/following?page=" + page + "&page_size=" + pageSize)
                .addHeader("Authorization", "Bearer " + authToken)
                .get()
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("获取关注列表失败");
                }
            }
        });
    }

    public static void getFollowerList(int page, int pageSize, ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/followers?page=" + page + "&page_size=" + pageSize)
                .addHeader("Authorization", "Bearer " + authToken)
                .get()
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("获取粉丝列表失败");
                }
            }
        });
    }

    // ==================== 通知接口 ====================

    public static void getNotifications(int page, int pageSize, ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/notifications?page=" + page + "&page_size=" + pageSize)
                .addHeader("Authorization", "Bearer " + authToken)
                .get()
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("获取通知失败");
                }
            }
        });
    }

    public static void getUnreadCount(ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/notifications/unread-count")
                .addHeader("Authorization", "Bearer " + authToken)
                .get()
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("获取未读数量失败");
                }
            }
        });
    }

    public static void markNotificationRead(int notificationId, ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/notifications/" + notificationId + "/read")
                .addHeader("Authorization", "Bearer " + authToken)
                .post(RequestBody.create("", JSON))
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("标记已读失败");
                }
            }
        });
    }

    public static void markAllNotificationsRead(ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/notifications/read-all")
                .addHeader("Authorization", "Bearer " + authToken)
                .post(RequestBody.create("", JSON))
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("标记全部已读失败");
                }
            }
        });
    }

    // ==================== 用户接口 ====================

    public static void getUserById(int userId, ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/users/" + userId)
                .get()
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("获取用户信息失败");
                }
            }
        });
    }

    public static void getUserArticles(int userId, int page, int pageSize, ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/users/" + userId + "/articles?page=" + page + "&page_size=" + pageSize)
                .get()
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("获取用户文章失败");
                }
            }
        });
    }

    // ==================== 公告接口 ====================

    public static void getAnnouncement(ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/announcement")
                .get()
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("获取公告失败");
                }
            }
        });
    }

    // ==================== 版本接口 ====================

    public static void getVersion(ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/version")
                .get()
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("获取版本信息失败");
                }
            }
        });
    }

    // ==================== 上传接口 ====================

    public static void uploadAvatar(byte[] imageData, String filename, ApiCallback callback) {
        RequestBody body = new MultipartBody.Builder()
                .setType(MultipartBody.FORM)
                .addFormDataPart("avatar", filename,
                        RequestBody.create(imageData, MediaType.parse("image/*")))
                .build();

        Request request = new Request.Builder()
                .url(BASE_URL + "/upload/avatar")
                .addHeader("Authorization", "Bearer " + authToken)
                .post(body)
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("上传头像失败");
                }
            }
        });
    }

    public static void uploadImage(byte[] imageData, String filename, ApiCallback callback) {
        RequestBody body = new MultipartBody.Builder()
                .setType(MultipartBody.FORM)
                .addFormDataPart("image", filename,
                        RequestBody.create(imageData, MediaType.parse("image/*")))
                .build();

        Request request = new Request.Builder()
                .url(BASE_URL + "/upload/image")
                .addHeader("Authorization", "Bearer " + authToken)
                .post(body)
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("上传图片失败");
                }
            }
        });
    }

    public static void uploadVideo(byte[] videoData, String filename, ApiCallback callback) {
        RequestBody body = new MultipartBody.Builder()
                .setType(MultipartBody.FORM)
                .addFormDataPart("video", filename,
                        RequestBody.create(videoData, MediaType.parse("video/*")))
                .build();

        Request request = new Request.Builder()
                .url(BASE_URL + "/upload/video")
                .addHeader("Authorization", "Bearer " + authToken)
                .post(body)
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("上传视频失败");
                }
            }
        });
    }

    public static void uploadVoice(byte[] voiceData, String filename, ApiCallback callback) {
        RequestBody body = new MultipartBody.Builder()
                .setType(MultipartBody.FORM)
                .addFormDataPart("voice", filename,
                        RequestBody.create(voiceData, MediaType.parse("audio/*")))
                .build();

        Request request = new Request.Builder()
                .url(BASE_URL + "/upload/voice")
                .addHeader("Authorization", "Bearer " + authToken)
                .post(body)
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("上传语音失败");
                }
            }
        });
    }

    // ==================== 用户通知接口 ====================

    public static void getUserNotifications(int page, int pageSize, ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/user-notifications?page=" + page + "&page_size=" + pageSize)
                .addHeader("Authorization", "Bearer " + authToken)
                .get()
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("获取用户通知失败");
                }
            }
        });
    }

    public static void getUserNotification(int notificationId, ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/user-notifications/" + notificationId)
                .addHeader("Authorization", "Bearer " + authToken)
                .get()
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("获取通知详情失败");
                }
            }
        });
    }

    public static void markUserNotificationRead(int notificationId, ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/user-notifications/" + notificationId + "/read")
                .addHeader("Authorization", "Bearer " + authToken)
                .post(RequestBody.create("", JSON))
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("标记已读失败");
                }
            }
        });
    }

    public static void markAllUserNotificationsRead(ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/user-notifications/read-all")
                .addHeader("Authorization", "Bearer " + authToken)
                .post(RequestBody.create("", JSON))
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("标记全部已读失败");
                }
            }
        });
    }

    // ==================== 粉丝通知接口 ====================

    public static void getFollowNotifications(int page, int pageSize, ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/follow-notifications?page=" + page + "&page_size=" + pageSize)
                .addHeader("Authorization", "Bearer " + authToken)
                .get()
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("获取粉丝通知失败");
                }
            }
        });
    }

    public static void markFollowNotificationRead(int notificationId, ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/follow-notifications/" + notificationId + "/read")
                .addHeader("Authorization", "Bearer " + authToken)
                .post(RequestBody.create("", JSON))
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("标记已读失败");
                }
            }
        });
    }

    public static void markAllFollowNotificationsRead(ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/follow-notifications/read-all")
                .addHeader("Authorization", "Bearer " + authToken)
                .post(RequestBody.create("", JSON))
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("标记全部已读失败");
                }
            }
        });
    }

    public static void getFollowNotificationsUnreadCount(ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/follow-notifications/unread-count")
                .addHeader("Authorization", "Bearer " + authToken)
                .get()
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("获取未读数量失败");
                }
            }
        });
    }

    // ==================== 评论回复通知接口 ====================

    public static void getCommentReplyNotifications(ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/comment-reply-notifications")
                .addHeader("Authorization", "Bearer " + authToken)
                .get()
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("获取评论回复通知失败");
                }
            }
        });
    }

    // ==================== 关注/粉丝列表接口 ====================

    public static void getUserFollowing(int userId, int page, int pageSize, ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/users/" + userId + "/following?page=" + page + "&page_size=" + pageSize)
                .get()
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("获取用户关注列表失败");
                }
            }
        });
    }

    public static void getUserFollowers(int userId, int page, int pageSize, ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/users/" + userId + "/followers?page=" + page + "&page_size=" + pageSize)
                .get()
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("获取用户粉丝列表失败");
                }
            }
        });
    }

    public static void getMutualFriends(ApiCallback callback) {
        Request request = new Request.Builder()
                .url(BASE_URL + "/mutual")
                .addHeader("Authorization", "Bearer " + authToken)
                .get()
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                callback.onError(e.getMessage());
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    callback.onSuccess(response.body().string());
                } else {
                    callback.onError("获取互关好友失败");
                }
            }
        });
    }
}
