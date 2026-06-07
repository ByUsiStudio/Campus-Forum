package miao.byusi.android.xylt;

import androidx.appcompat.app.AlertDialog;
import androidx.swiperefreshlayout.widget.SwipeRefreshLayout;

import android.os.Bundle;
import android.view.View;
import android.widget.EditText;
import android.widget.ImageView;
import android.widget.TextView;

import com.bumptech.glide.Glide;

import org.json.JSONObject;

public class ProfileActivity extends BaseActivity {

    private SwipeRefreshLayout swipeRefreshLayout;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_profile);

        initViews();
        loadUserProfile();
    }

    @Override
    protected void onAutoRefresh() {
        loadUserProfile();
    }

    private void initViews() {
        swipeRefreshLayout = findViewById(R.id.swipe_refresh);
        setupSwipeRefresh(swipeRefreshLayout);
    }

    private void loadUserProfile() {
        if (swipeRefreshLayout != null) {
            swipeRefreshLayout.setRefreshing(true);
        }

        ApiClient.getProfile(new ApiCallback() {
            @Override
            public void onSuccess(String response) {
                runOnUiThread(() -> {
                    try {
                        JSONObject json = new JSONObject(response);
                        displayProfile(json);
                    } catch (Exception e) {
                        showToast("加载用户信息失败");
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
                    showToast("获取用户信息失败: " + error);
                    if (swipeRefreshLayout != null) {
                        swipeRefreshLayout.setRefreshing(false);
                    }
                });
            }
        });
    }

    private void displayProfile(JSONObject json) {
        ImageView ivAvatar = findViewById(R.id.iv_avatar);
        TextView tvUsername = findViewById(R.id.tv_username);
        TextView tvEmail = findViewById(R.id.tv_email);
        TextView tvBio = findViewById(R.id.tv_bio);
        TextView tvArticlesCount = findViewById(R.id.tv_articles_count);
        TextView tvFollowersCount = findViewById(R.id.tv_followers_count);
        TextView tvFollowingCount = findViewById(R.id.tv_following_count);

        String username = json.optString("username", "");
        String displayName = json.optString("display_name", username);
        String signature = json.optString("signature", "");
        String avatarUrl = json.optString("avatar", json.optString("avatar_url", ""));

        tvUsername.setText(displayName);
        tvEmail.setText("@" + username);
        tvBio.setText(signature.isEmpty() ? "暂无个人简介" : signature);

        int articlesCount = json.optInt("article_count", 0);
        int followersCount = json.optInt("follower_count", 0);
        int followingCount = json.optInt("following_count", 0);

        tvArticlesCount.setText(String.valueOf(articlesCount));
        tvFollowersCount.setText(String.valueOf(followersCount));
        tvFollowingCount.setText(String.valueOf(followingCount));

        if (!avatarUrl.isEmpty()) {
            Glide.with(this).load(avatarUrl).into(ivAvatar);
        }
    }

    public void onEditProfileClick(View view) {
        EditText displayNameInput = new EditText(this);
        displayNameInput.setHint("请输入显示名称");

        EditText signatureInput = new EditText(this);
        signatureInput.setHint("请输入个人简介");

        TextView tvDisplayName = findViewById(R.id.tv_username);
        TextView tvSignature = findViewById(R.id.tv_bio);

        displayNameInput.setText(tvDisplayName.getText().toString().replace("@", ""));
        signatureInput.setText(tvSignature.getText().toString().equals("暂无个人简介") ? "" : tvSignature.getText().toString());

        new AlertDialog.Builder(this)
                .setTitle("编辑资料")
                .setMessage("请输入新资料")
                .setView(createDialogView(displayNameInput, signatureInput))
                .setPositiveButton("保存", (dialog, which) -> {
                    String displayName = displayNameInput.getText().toString();
                    String signature = signatureInput.getText().toString();
                    updateProfile(displayName, signature);
                })
                .setNegativeButton("取消", null)
                .show();
    }

    private View createDialogView(EditText displayNameInput, EditText signatureInput) {
        android.widget.LinearLayout layout = new android.widget.LinearLayout(this);
        layout.setOrientation(android.widget.LinearLayout.VERTICAL);
        layout.setPadding(50, 20, 50, 20);

        TextView label1 = new TextView(this);
        label1.setText("显示名称:");
        layout.addView(label1);
        layout.addView(displayNameInput);

        TextView label2 = new TextView(this);
        label2.setText("个人简介:");
        layout.addView(label2);
        layout.addView(signatureInput);

        return layout;
    }

    private void updateProfile(String displayName, String signature) {
        ApiClient.updateProfile(displayName, signature, new ApiCallback() {
            @Override
            public void onSuccess(String response) {
                runOnUiThread(() -> {
                    showToast("资料更新成功");
                    loadUserProfile();
                });
            }

            @Override
            public void onError(String error) {
                runOnUiThread(() -> {
                    showToast("更新失败: " + error);
                });
            }
        });
    }
}
