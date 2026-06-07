package miao.byusi.android.xylt;

import androidx.appcompat.app.AppCompatActivity;

import android.os.Bundle;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;
import android.widget.Toast;

public class CreateArticleActivity extends AppCompatActivity {

    private EditText etTitle;
    private EditText etContent;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_create_article);

        initViews();
        setupListeners();
    }

    private void initViews() {
        etTitle = findViewById(R.id.et_title);
        etContent = findViewById(R.id.et_content);
    }

    private void setupListeners() {
        Button btnSubmit = findViewById(R.id.btn_submit);
        Button btnCancel = findViewById(R.id.btn_cancel);

        btnSubmit.setOnClickListener(v -> {
            String title = etTitle.getText().toString().trim();
            String content = etContent.getText().toString().trim();

            if (title.isEmpty()) {
                Toast.makeText(this, "请输入标题", Toast.LENGTH_SHORT).show();
                return;
            }

            if (content.isEmpty()) {
                Toast.makeText(this, "请输入内容", Toast.LENGTH_SHORT).show();
                return;
            }

            // 提交文章
            submitArticle(title, content);
        });

        btnCancel.setOnClickListener(v -> {
            finish();
        });
    }

    private void submitArticle(String title, String content) {
        // 调用API提交文章
        Toast.makeText(this, "文章发布成功", Toast.LENGTH_SHORT).show();
        finish();
    }
}
