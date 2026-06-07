package miao.byusi.android.xylt;

import android.os.Bundle;
import android.os.Handler;
import android.os.Looper;
import android.view.View;
import android.widget.Toast;

import androidx.annotation.Nullable;
import androidx.appcompat.app.AppCompatActivity;
import androidx.swiperefreshlayout.widget.SwipeRefreshLayout;

public abstract class BaseActivity extends AppCompatActivity {

    protected Handler handler = new Handler(Looper.getMainLooper());
    protected Runnable refreshRunnable;
    protected static final long AUTO_REFRESH_INTERVAL = 30000; // 30秒自动刷新

    @Override
    protected void onCreate(@Nullable Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setupAutoRefresh();
    }

    private void setupAutoRefresh() {
        refreshRunnable = new Runnable() {
            @Override
            public void run() {
                if (isAutoRefreshEnabled()) {
                    onAutoRefresh();
                }
                handler.postDelayed(this, AUTO_REFRESH_INTERVAL);
            }
        };
    }

    @Override
    protected void onResume() {
        super.onResume();
        startAutoRefresh();
    }

    @Override
    protected void onPause() {
        super.onPause();
        stopAutoRefresh();
    }

    @Override
    protected void onDestroy() {
        super.onDestroy();
        stopAutoRefresh();
        handler.removeCallbacksAndMessages(null);
    }

    protected void startAutoRefresh() {
        if (refreshRunnable != null) {
            handler.removeCallbacks(refreshRunnable);
            handler.postDelayed(refreshRunnable, AUTO_REFRESH_INTERVAL);
        }
    }

    protected void stopAutoRefresh() {
        if (refreshRunnable != null) {
            handler.removeCallbacks(refreshRunnable);
        }
    }

    protected void triggerRefresh() {
        stopAutoRefresh();
        onAutoRefresh();
        startAutoRefresh();
    }

    protected boolean isAutoRefreshEnabled() {
        return true;
    }

    protected void onAutoRefresh() {
    }

    protected void setupSwipeRefresh(SwipeRefreshLayout swipeRefreshLayout) {
        swipeRefreshLayout.setColorSchemeResources(R.color.colorPrimary);
        swipeRefreshLayout.setOnRefreshListener(() -> {
            triggerRefresh();
            swipeRefreshLayout.setRefreshing(false);
        });
    }

    protected void showToast(String message) {
        Toast.makeText(this, message, Toast.LENGTH_SHORT).show();
    }

    protected void showLongToast(String message) {
        Toast.makeText(this, message, Toast.LENGTH_LONG).show();
    }
}
