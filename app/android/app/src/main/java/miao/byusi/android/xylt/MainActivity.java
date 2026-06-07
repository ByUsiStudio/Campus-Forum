package miao.byusi.android.xylt;

import androidx.recyclerview.widget.LinearLayoutManager;
import androidx.recyclerview.widget.RecyclerView;
import androidx.swiperefreshlayout.widget.SwipeRefreshLayout;

import android.content.Intent;
import android.os.Bundle;
import android.view.View;

import com.google.android.material.floatingactionbutton.FloatingActionButton;

import org.json.JSONArray;
import org.json.JSONObject;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

public class MainActivity extends BaseActivity {

    private RecyclerView articleRecyclerView;
    private ArticleAdapter articleAdapter;
    private List<Article> articleList;
    private int currentPage = 1;
    private final int pageSize = 20;
    private SwipeRefreshLayout swipeRefreshLayout;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);

        initViews();
        loadArticles();
        setupListeners();
    }

    @Override
    protected void onAutoRefresh() {
        loadArticles();
    }

    private void initViews() {
        swipeRefreshLayout = findViewById(R.id.swipe_refresh);
        setupSwipeRefresh(swipeRefreshLayout);

        articleRecyclerView = findViewById(R.id.recycler_view);
        articleRecyclerView.setLayoutManager(new LinearLayoutManager(this));
        
        articleList = new ArrayList<>();
        articleAdapter = new ArticleAdapter(this, articleList);
        articleRecyclerView.setAdapter(articleAdapter);
    }

    private void loadArticles() {
        if (swipeRefreshLayout != null) {
            swipeRefreshLayout.setRefreshing(true);
        }

        ApiClient.getArticles(currentPage, pageSize, new ApiCallback() {
            @Override
            public void onSuccess(String response) {
                runOnUiThread(() -> {
                    try {
                        JSONObject json = new JSONObject(response);
                        JSONArray articles = json.optJSONArray("articles");
                        
                        articleList.clear();
                        
                        if (articles != null) {
                            for (int i = 0; i < articles.length(); i++) {
                                JSONObject articleJson = articles.getJSONObject(i);
                                Article article = parseArticle(articleJson);
                                articleList.add(article);
                            }
                        }
                        
                        Collections.sort(articleList, (a, b) -> {
                            if (a.isTop() && !b.isTop()) return -1;
                            if (!a.isTop() && b.isTop()) return 1;
                            return b.getCreatedAt().compareTo(a.getCreatedAt());
                        });
                        
                        articleAdapter.notifyDataSetChanged();
                    } catch (Exception e) {
                        showToast("加载文章失败");
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
                    showToast("加载文章失败: " + error);
                    if (swipeRefreshLayout != null) {
                        swipeRefreshLayout.setRefreshing(false);
                    }
                });
            }
        });
    }

    private Article parseArticle(JSONObject json) {
        Article article = new Article();
        article.setId(json.optInt("id", 0));
        article.setTitle(json.optString("title", ""));
        article.setContent(json.optString("content", ""));
        
        JSONObject user = json.optJSONObject("user");
        if (user != null) {
            article.setAuthor(user.optString("nickname", user.optString("username", "匿名")));
            article.setAvatarUrl(user.optString("avatar_url", ""));
        } else {
            article.setAuthor(json.optString("author", "匿名"));
        }
        
        article.setViewCount(json.optInt("view_count", 0));
        article.setLikeCount(json.optInt("like_count", 0));
        article.setCommentCount(json.optInt("comment_count", 0));
        article.setCreatedAt(json.optString("created_at", ""));
        article.setTop(json.optBoolean("is_top", false));
        
        return article;
    }

    private void setupListeners() {
        FloatingActionButton createBtn = findViewById(R.id.fab_create);
        createBtn.setOnClickListener(v -> {
            Intent intent = new Intent(MainActivity.this, CreateArticleActivity.class);
            startActivity(intent);
        });

        View draftsBtn = findViewById(R.id.btn_drafts);
        draftsBtn.setOnClickListener(v -> {
            Intent intent = new Intent(MainActivity.this, DraftsActivity.class);
            startActivity(intent);
        });

        articleAdapter.setOnItemClickListener(article -> {
            Intent intent = new Intent(MainActivity.this, ArticleDetailActivity.class);
            intent.putExtra("article_id", article.getId());
            startActivity(intent);
        });
    }
}
