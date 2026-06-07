package miao.byusi.android.xylt;

import androidx.appcompat.app.AlertDialog;
import androidx.swiperefreshlayout.widget.SwipeRefreshLayout;

import android.os.Bundle;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;
import android.widget.ImageView;
import android.widget.TextView;

import com.bumptech.glide.Glide;

import org.json.JSONArray;
import org.json.JSONObject;

public class ArticleDetailActivity extends BaseActivity {

    private int articleId;
    private int userId;
    private SwipeRefreshLayout swipeRefreshLayout;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_article_detail);

        articleId = getIntent().getIntExtra("article_id", 0);
        
        initViews();
        loadArticleDetail(articleId);
        loadComments(articleId);
        setupListeners();
    }

    @Override
    protected void onAutoRefresh() {
        loadArticleDetail(articleId);
        loadComments(articleId);
    }

    private void initViews() {
        swipeRefreshLayout = findViewById(R.id.swipe_refresh);
        setupSwipeRefresh(swipeRefreshLayout);
    }

    private void loadArticleDetail(int articleId) {
        if (swipeRefreshLayout != null) {
            swipeRefreshLayout.setRefreshing(true);
        }

        ApiClient.getArticleDetail(articleId, new ApiCallback() {
            @Override
            public void onSuccess(String response) {
                runOnUiThread(() -> {
                    try {
                        JSONObject json = new JSONObject(response);
                        displayArticleDetail(json);
                    } catch (Exception e) {
                        showToast("加载文章详情失败");
                    } finally {
                        if (swipeRefreshLayout != null) {
                            swipeRefreshLayout.setRefreshing(false);
                        }
                    }
                });
            }

            @Override
            public void onError(String error) {
                runOnUiThread(() -> {
                    showToast("加载文章详情失败: " + error);
                    if (swipeRefreshLayout != null) {
                        swipeRefreshLayout.setRefreshing(false);
                    }
                });
            }
        });
    }

    private void displayArticleDetail(JSONObject json) {
        TextView tvTitle = findViewById(R.id.tv_article_title);
        TextView tvAuthor = findViewById(R.id.tv_author);
        TextView tvContent = findViewById(R.id.tv_content);
        TextView tvDate = findViewById(R.id.tv_date);
        TextView tvViews = findViewById(R.id.tv_views);
        TextView tvLikes = findViewById(R.id.tv_likes);
        TextView tvComments = findViewById(R.id.tv_comments);
        ImageView ivAvatar = findViewById(R.id.iv_author_avatar);

        tvTitle.setText(json.optString("title", ""));
        tvContent.setText(json.optString("content", ""));
        tvDate.setText(json.optString("created_at", ""));
        tvViews.setText(String.valueOf(json.optInt("view_count", 0)));
        tvLikes.setText(String.valueOf(json.optInt("like_count", 0)));
        tvComments.setText(String.valueOf(json.optInt("comment_count", 0)));

        JSONObject user = json.optJSONObject("user");
        if (user != null) {
            String author = user.optString("nickname", user.optString("username", "匿名"));
            tvAuthor.setText(author);
            userId = user.optInt("id", 0);
            
            String avatarUrl = user.optString("avatar_url", "");
            if (!avatarUrl.isEmpty()) {
                Glide.with(this).load(avatarUrl).into(ivAvatar);
            }
        }
    }

    private void loadComments(int articleId) {
        ApiClient.getComments(articleId, 1, 50, new ApiCallback() {
            @Override
            public void onSuccess(String response) {
                runOnUiThread(() -> {
                    try {
                        JSONObject json = new JSONObject(response);
                        JSONArray comments = json.optJSONArray("comments");
                        displayComments(comments);
                    } catch (Exception e) {
                        showToast("加载评论失败");
                    }
                });
            }

            @Override
            public void onError(String error) {
                runOnUiThread(() -> {
                    showToast("加载评论失败: " + error);
                });
            }
        });
    }

    private void displayComments(JSONArray comments) {
        TextView tvCommentsList = findViewById(R.id.tv_comments_list);
        if (comments == null || comments.length() == 0) {
            tvCommentsList.setText("暂无评论");
            return;
        }

        StringBuilder sb = new StringBuilder();
        for (int i = 0; i < comments.length(); i++) {
            try {
                JSONObject comment = comments.getJSONObject(i);
                String author = "匿名";
                JSONObject commentUser = comment.optJSONObject("user");
                if (commentUser != null) {
                    author = commentUser.optString("nickname", commentUser.optString("username", "匿名"));
                }
                String content = comment.optString("content", "");
                String createdAt = comment.optString("created_at", "");
                
                sb.append(author).append(": ").append(content).append("\n");
                sb.append(createdAt).append("\n\n");
            } catch (Exception e) {
                // 跳过无效评论
            }
        }
        tvCommentsList.setText(sb.toString());
    }

    private void setupListeners() {
        Button btnLike = findViewById(R.id.btn_like);
        Button btnComment = findViewById(R.id.btn_comment);
        Button btnShare = findViewById(R.id.btn_share);
        Button btnReport = findViewById(R.id.btn_report);

        btnLike.setOnClickListener(v -> likeArticle(articleId));

        btnComment.setOnClickListener(v -> showCommentDialog());

        btnShare.setOnClickListener(v -> {
            showToast("分享功能开发中");
        });

        btnReport.setOnClickListener(v -> showReportDialog());
    }

    private void likeArticle(int articleId) {
        ApiClient.likeArticle(articleId, new ApiCallback() {
            @Override
            public void onSuccess(String response) {
                runOnUiThread(() -> {
                    showToast("点赞成功");
                    loadArticleDetail(articleId);
                });
            }

            @Override
            public void onError(String error) {
                runOnUiThread(() -> {
                    showToast("点赞失败: " + error);
                });
            }
        });
    }

    private void showCommentDialog() {
        EditText commentInput = new EditText(this);
        commentInput.setHint("写下你的评论...");
        
        new AlertDialog.Builder(this)
                .setTitle("发表评论")
                .setView(commentInput)
                .setPositiveButton("发表", (dialog, which) -> {
                    String content = commentInput.getText().toString();
                    if (!content.isEmpty()) {
                        submitComment(articleId, content);
                    }
                })
                .setNegativeButton("取消", null)
                .show();
    }

    private void submitComment(int articleId, String content) {
        ApiClient.createComment(articleId, content, null, new ApiCallback() {
            @Override
            public void onSuccess(String response) {
                runOnUiThread(() -> {
                    showToast("评论发表成功");
                    loadComments(articleId);
                });
            }

            @Override
            public void onError(String error) {
                runOnUiThread(() -> {
                    showToast("发表评论失败: " + error);
                });
            }
        });
    }

    private void showReportDialog() {
        String[] reasons = {"垃圾广告", "色情低俗", "暴力血腥", "政治敏感", "违法犯罪", "谣言虚假", "侵犯隐私", "其他违规"};
        
        new AlertDialog.Builder(this)
                .setTitle("举报文章")
                .setItems(reasons, (dialog, which) -> {
                    String reason = reasons[which];
                    showReportDescriptionDialog(reason);
                })
                .setNegativeButton("取消", null)
                .show();
    }

    private void showReportDescriptionDialog(String reason) {
        EditText descInput = new EditText(this);
        descInput.setHint("请详细描述举报原因...");
        
        new AlertDialog.Builder(this)
                .setTitle("举报原因: " + reason)
                .setView(descInput)
                .setPositiveButton("提交", (dialog, which) -> {
                    String description = descInput.getText().toString();
                    submitReport(reason, description);
                })
                .setNegativeButton("取消", null)
                .show();
    }

    private void submitReport(String reason, String description) {
        ApiClient.submitReport("article", articleId, reason, description, new ApiCallback() {
            @Override
            public void onSuccess(String response) {
                runOnUiThread(() -> {
                    showToast("举报已提交，感谢您的反馈");
                });
            }

            @Override
            public void onError(String error) {
                runOnUiThread(() -> {
                    showToast("举报失败: " + error);
                });
            }
        });
    }
}
