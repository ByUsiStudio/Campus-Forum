package miao.byusi.android.xylt;

import androidx.appcompat.app.AlertDialog;
import androidx.recyclerview.widget.LinearLayoutManager;
import androidx.recyclerview.widget.RecyclerView;
import androidx.swiperefreshlayout.widget.SwipeRefreshLayout;

import android.content.Intent;
import android.os.Bundle;
import android.view.View;

import org.json.JSONArray;
import org.json.JSONObject;

import java.util.ArrayList;
import java.util.List;

public class DraftsActivity extends BaseActivity {

    private RecyclerView recyclerView;
    private DraftsAdapter adapter;
    private List<Article> draftList;
    private SwipeRefreshLayout swipeRefreshLayout;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_drafts);

        initViews();
        loadDrafts();
    }

    @Override
    protected void onAutoRefresh() {
        loadDrafts();
    }

    private void initViews() {
        swipeRefreshLayout = findViewById(R.id.swipe_refresh);
        setupSwipeRefresh(swipeRefreshLayout);

        recyclerView = findViewById(R.id.recycler_view);
        recyclerView.setLayoutManager(new LinearLayoutManager(this));
        
        draftList = new ArrayList<>();
        adapter = new DraftsAdapter(this, draftList);
        recyclerView.setAdapter(adapter);

        adapter.setOnItemClickListener(article -> {
            Intent intent = new Intent(DraftsActivity.this, CreateArticleActivity.class);
            intent.putExtra("draft_id", article.getId());
            startActivity(intent);
        });

        adapter.setOnPublishClickListener(article -> {
            publishDraft(article.getId());
        });

        adapter.setOnDeleteClickListener(article -> {
            showDeleteDialog(article.getId());
        });
    }

    private void loadDrafts() {
        if (swipeRefreshLayout != null) {
            swipeRefreshLayout.setRefreshing(true);
        }

        ApiClient.getMyDrafts(1, 50, new ApiCallback() {
            @Override
            public void onSuccess(String response) {
                runOnUiThread(() -> {
                    try {
                        JSONObject json = new JSONObject(response);
                        JSONArray drafts = json.optJSONArray("drafts");
                        
                        draftList.clear();
                        
                        if (drafts != null) {
                            for (int i = 0; i < drafts.length(); i++) {
                                JSONObject draftJson = drafts.getJSONObject(i);
                                Article article = parseDraft(draftJson);
                                draftList.add(article);
                            }
                        }
                        
                        adapter.notifyDataSetChanged();
                        
                        if (draftList.isEmpty()) {
                            findViewById(R.id.tv_empty).setVisibility(View.VISIBLE);
                        } else {
                            findViewById(R.id.tv_empty).setVisibility(View.GONE);
                        }
                    } catch (Exception e) {
                        showToast("加载草稿失败");
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
                    showToast("加载草稿失败: " + error);
                    if (swipeRefreshLayout != null) {
                        swipeRefreshLayout.setRefreshing(false);
                    }
                });
            }
        });
    }

    private Article parseDraft(JSONObject json) {
        Article article = new Article();
        article.setId(json.optInt("id", 0));
        article.setTitle(json.optString("title", ""));
        article.setContent(json.optString("content", ""));
        article.setCreatedAt(json.optString("created_at", ""));
        return article;
    }

    private void publishDraft(int draftId) {
        ApiClient.publishDraft(draftId, new ApiCallback() {
            @Override
            public void onSuccess(String response) {
                runOnUiThread(() -> {
                    showToast("发布成功");
                    loadDrafts();
                });
            }

            @Override
            public void onError(String error) {
                runOnUiThread(() -> {
                    showToast("发布失败: " + error);
                });
            }
        });
    }

    private void showDeleteDialog(int draftId) {
        new AlertDialog.Builder(this)
                .setTitle("删除草稿")
                .setMessage("确定要删除这篇草稿吗？")
                .setPositiveButton("删除", (dialog, which) -> deleteDraft(draftId))
                .setNegativeButton("取消", null)
                .show();
    }

    private void deleteDraft(int draftId) {
        ApiClient.deleteArticle(draftId, new ApiCallback() {
            @Override
            public void onSuccess(String response) {
                runOnUiThread(() -> {
                    showToast("删除成功");
                    loadDrafts();
                });
            }

            @Override
            public void onError(String error) {
                runOnUiThread(() -> {
                    showToast("删除失败: " + error);
                });
            }
        });
    }
}
