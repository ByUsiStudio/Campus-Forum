package miao.byusi.pc.xylt;

import javafx.animation.KeyFrame;
import javafx.animation.Timeline;
import javafx.fxml.FXML;
import javafx.geometry.Insets;
import javafx.scene.control.*;
import javafx.scene.layout.*;
import javafx.scene.text.Text;
import javafx.util.Duration;
import org.json.JSONArray;
import org.json.JSONObject;

import java.util.ArrayList;
import java.util.List;

public class MainController {

    @FXML
    private VBox articleList;

    @FXML
    private TextField searchField;

    @FXML
    private Label statusLabel;

    @FXML
    private Label versionLabel;

    private List<Article> currentArticles = new ArrayList<>();
    private int currentPage = 1;
    private final int pageSize = 20;
    private boolean isLoggedIn = false;

    @FXML
    private void initialize() {
        loadArticles();
        loadVersion();
        startStatusUpdate();
    }

    @FXML
    private void loadArticles() {
        ApiClient.getArticles(currentPage, pageSize, new ApiCallback() {
            @Override
            public void onSuccess(String response) {
                javafx.application.Platform.runLater(() -> {
                    try {
                        JSONObject json = new JSONObject(response);
                        JSONArray articles = json.getJSONArray("articles");
                        currentArticles.clear();
                        articleList.getChildren().clear();

                        for (int i = 0; i < articles.length(); i++) {
                            JSONObject articleJson = articles.getJSONObject(i);
                            Article article = parseArticle(articleJson);
                            currentArticles.add(article);
                            addArticleCard(article);
                        }

                        updateStatus("文章加载成功，共 " + articles.length() + " 篇");
                    } catch (Exception e) {
                        updateStatus("解析文章数据失败");
                    }
                });
            }

            @Override
            public void onError(String error) {
                javafx.application.Platform.runLater(() -> {
                    updateStatus("加载文章失败: " + error);
                });
            }
        });
    }

    @FXML
    private void searchArticles() {
        String keyword = searchField.getText();
        if (keyword.isEmpty()) {
            loadArticles();
            return;
        }

        ApiClient.searchArticles(keyword, 1, pageSize, new ApiCallback() {
            @Override
            public void onSuccess(String response) {
                javafx.application.Platform.runLater(() -> {
                    try {
                        JSONObject json = new JSONObject(response);
                        JSONArray articles = json.getJSONArray("articles");
                        currentArticles.clear();
                        articleList.getChildren().clear();

                        for (int i = 0; i < articles.length(); i++) {
                            JSONObject articleJson = articles.getJSONObject(i);
                            Article article = parseArticle(articleJson);
                            currentArticles.add(article);
                            addArticleCard(article);
                        }

                        updateStatus("搜索到 " + articles.length() + " 篇文章");
                    } catch (Exception e) {
                        updateStatus("搜索解析失败");
                    }
                });
            }

            @Override
            public void onError(String error) {
                javafx.application.Platform.runLater(() -> {
                    updateStatus("搜索失败: " + error);
                });
            }
        });
    }

    @FXML
    private void showLoginDialog() {
        Dialog<ButtonType> dialog = new Dialog<>();
        dialog.setTitle("登录");
        dialog.setHeaderText("请输入用户名和密码");

        GridPane grid = new GridPane();
        grid.setHgap(10);
        grid.setVgap(10);
        grid.setPadding(new Insets(20, 150, 10, 10));

        TextField usernameField = new TextField();
        usernameField.setPromptText("用户名");
        PasswordField passwordField = new PasswordField();
        passwordField.setPromptText("密码");

        grid.add(new Label("用户名:"), 0, 0);
        grid.add(usernameField, 1, 0);
        grid.add(new Label("密码:"), 0, 1);
        grid.add(passwordField, 1, 1);

        dialog.getDialogPane().setContent(grid);
        dialog.getDialogPane().getButtonTypes().addAll(ButtonType.OK, ButtonType.CANCEL);

        dialog.showAndWait().ifPresent(result -> {
            if (result == ButtonType.OK) {
                String username = usernameField.getText();
                String password = passwordField.getText();
                performLogin(username, password);
            }
        });
    }

    private void performLogin(String username, String password) {
        ApiClient.login(username, password, new ApiCallback() {
            @Override
            public void onSuccess(String response) {
                javafx.application.Platform.runLater(() -> {
                    try {
                        JSONObject json = new JSONObject(response);
                        if (json.has("token")) {
                            String token = json.getString("token");
                            ApiClient.setAuthToken(token);
                            isLoggedIn = true;
                            updateStatus("登录成功");
                            loadUserProfile();
                        } else {
                            updateStatus("登录失败：未获取到令牌");
                        }
                    } catch (Exception e) {
                        updateStatus("登录解析失败");
                    }
                });
            }

            @Override
            public void onError(String error) {
                javafx.application.Platform.runLater(() -> {
                    updateStatus("登录失败: " + error);
                });
            }
        });
    }

    private void loadUserProfile() {
        ApiClient.getProfile(new ApiCallback() {
            @Override
            public void onSuccess(String response) {
                javafx.application.Platform.runLater(() -> {
                    try {
                        JSONObject json = new JSONObject(response);
                        String nickname = json.optString("nickname", json.optString("username", "用户"));
                        updateStatus("欢迎, " + nickname);
                    } catch (Exception e) {
                        updateStatus("获取用户信息失败");
                    }
                });
            }

            @Override
            public void onError(String error) {
                javafx.application.Platform.runLater(() -> {
                    updateStatus("获取用户信息失败: " + error);
                });
            }
        });
    }

    @FXML
    private void showRegisterDialog() {
        Dialog<ButtonType> dialog = new Dialog<>();
        dialog.setTitle("注册");
        dialog.setHeaderText("请填写注册信息");

        GridPane grid = new GridPane();
        grid.setHgap(10);
        grid.setVgap(10);
        grid.setPadding(new Insets(20, 150, 10, 10));

        TextField usernameField = new TextField();
        usernameField.setPromptText("用户名");
        TextField qqNumberField = new TextField();
        qqNumberField.setPromptText("QQ号码");
        TextField displayNameField = new TextField();
        displayNameField.setPromptText("显示名称");
        PasswordField passwordField = new PasswordField();
        passwordField.setPromptText("密码");

        grid.add(new Label("用户名:"), 0, 0);
        grid.add(usernameField, 1, 0);
        grid.add(new Label("QQ号码:"), 0, 1);
        grid.add(qqNumberField, 1, 1);
        grid.add(new Label("显示名称:"), 0, 2);
        grid.add(displayNameField, 1, 2);
        grid.add(new Label("密码:"), 0, 3);
        grid.add(passwordField, 1, 3);

        dialog.getDialogPane().setContent(grid);
        dialog.getDialogPane().getButtonTypes().addAll(ButtonType.OK, ButtonType.CANCEL);

        dialog.showAndWait().ifPresent(result -> {
            if (result == ButtonType.OK) {
                String username = usernameField.getText();
                String qqNumber = qqNumberField.getText();
                String displayName = displayNameField.getText();
                String password = passwordField.getText();
                performRegister(username, qqNumber, displayName, password);
            }
        });
    }

    private void performRegister(String username, String qqNumber, String displayName, String password) {
        ApiClient.register(username, qqNumber, displayName, password, new ApiCallback() {
            @Override
            public void onSuccess(String response) {
                javafx.application.Platform.runLater(() -> {
                    updateStatus("注册成功，请登录");
                });
            }

            @Override
            public void onError(String error) {
                javafx.application.Platform.runLater(() -> {
                    updateStatus("注册失败: " + error);
                });
            }
        });
    }

    @FXML
    private void createArticle() {
        if (!isLoggedIn) {
            showLoginDialog();
            return;
        }

        Dialog<ButtonType> dialog = new Dialog<>();
        dialog.setTitle("发布文章");
        dialog.setHeaderText("请填写文章信息");

        GridPane grid = new GridPane();
        grid.setHgap(10);
        grid.setVgap(10);
        grid.setPadding(new Insets(20, 150, 10, 10));

        TextField titleField = new TextField();
        titleField.setPromptText("文章标题");
        TextArea contentArea = new TextArea();
        contentArea.setPromptText("文章内容");
        contentArea.setPrefRowCount(8);

        grid.add(new Label("标题:"), 0, 0);
        grid.add(titleField, 1, 0);
        grid.add(new Label("内容:"), 0, 1);
        grid.add(contentArea, 1, 1);

        dialog.getDialogPane().setContent(grid);
        dialog.getDialogPane().getButtonTypes().addAll(ButtonType.OK, ButtonType.CANCEL);

        dialog.showAndWait().ifPresent(result -> {
            if (result == ButtonType.OK) {
                String title = titleField.getText();
                String content = contentArea.getText();
                submitArticle(title, content);
            }
        });
    }

    private void submitArticle(String title, String content) {
        ApiClient.createArticle(title, content, 1, false, new ApiCallback() {
            @Override
            public void onSuccess(String response) {
                javafx.application.Platform.runLater(() -> {
                    updateStatus("文章发布成功");
                    loadArticles();
                });
            }

            @Override
            public void onError(String error) {
                javafx.application.Platform.runLater(() -> {
                    updateStatus("发布失败: " + error);
                });
            }
        });
    }

    @FXML
    private void refreshArticles() {
        currentPage = 1;
        loadArticles();
    }

    private void loadVersion() {
        ApiClient.getVersion(new ApiCallback() {
            @Override
            public void onSuccess(String response) {
                javafx.application.Platform.runLater(() -> {
                    try {
                        JSONObject json = new JSONObject(response);
                        String frontend = json.optString("frontend_version", "1.0.0");
                        String backend = json.optString("backend_version", "1.0.0");
                        versionLabel.setText("版本: 前端 " + frontend + " | 后端 " + backend);
                    } catch (Exception e) {
                        versionLabel.setText("版本: 1.0.0");
                    }
                });
            }

            @Override
            public void onError(String error) {
                javafx.application.Platform.runLater(() -> {
                    versionLabel.setText("版本: 1.0.0");
                });
            }
        });
    }

    private void startStatusUpdate() {
        Timeline timeline = new Timeline(new KeyFrame(Duration.seconds(30), event -> {
            updateStatus("在线");
        }));
        timeline.setCycleCount(Timeline.INDEFINITE);
        timeline.play();
    }

    private void updateStatus(String message) {
        statusLabel.setText(message);
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

        return article;
    }

    private void addArticleCard(Article article) {
        VBox card = new VBox();
        card.getStyleClass().add("article-card");
        card.setSpacing(12);
        card.setPadding(new Insets(16));
        card.setPrefWidth(800);

        Label titleLabel = new Label(article.getTitle());
        titleLabel.getStyleClass().add("article-title");

        Label authorLabel = new Label(article.getAuthor());
        authorLabel.getStyleClass().add("article-author");

        Label contentLabel = new Label(article.getContent());
        contentLabel.getStyleClass().add("article-content");
        contentLabel.setWrapText(true);

        HBox statsBox = new HBox();
        statsBox.setSpacing(24);
        Label viewsLabel = new Label("浏览: " + article.getViewCount());
        Label likesLabel = new Label("点赞: " + article.getLikeCount());
        Label commentsLabel = new Label("评论: " + article.getCommentCount());
        statsBox.getChildren().addAll(viewsLabel, likesLabel, commentsLabel);

        HBox actionBox = new HBox();
        actionBox.setSpacing(8);
        Button likeBtn = new Button("点赞");
        likeBtn.getStyleClass().add("action-button");
        likeBtn.setOnAction(e -> likeArticle(article.getId()));

        Button commentBtn = new Button("评论");
        commentBtn.getStyleClass().add("action-button");

        Button shareBtn = new Button("分享");
        shareBtn.getStyleClass().add("action-button");

        Button reportBtn = new Button("举报");
        reportBtn.getStyleClass().add("action-button");
        reportBtn.setOnAction(e -> showReportDialog(article.getId()));

        actionBox.getChildren().addAll(likeBtn, commentBtn, shareBtn, reportBtn);

        card.getChildren().addAll(titleLabel, authorLabel, contentLabel, statsBox, actionBox);

        card.setOnMouseClicked(event -> {
            if (event.getClickCount() == 2) {
                showArticleDetail(article.getId());
            }
        });

        articleList.getChildren().add(card);
    }

    private void likeArticle(int articleId) {
        if (!isLoggedIn) {
            showLoginDialog();
            return;
        }

        ApiClient.likeArticle(articleId, new ApiCallback() {
            @Override
            public void onSuccess(String response) {
                javafx.application.Platform.runLater(() -> {
                    updateStatus("点赞成功");
                    refreshArticles();
                });
            }

            @Override
            public void onError(String error) {
                javafx.application.Platform.runLater(() -> {
                    updateStatus("点赞失败: " + error);
                });
            }
        });
    }

    private void showReportDialog(int articleId) {
        Dialog<ButtonType> dialog = new Dialog<>();
        dialog.setTitle("举报文章");
        dialog.setHeaderText("请选择举报原因");

        VBox content = new VBox();
        content.setSpacing(10);
        content.setPadding(new Insets(20));

        ComboBox<String> reasonCombo = new ComboBox<>();
        reasonCombo.getItems().addAll(
            "垃圾广告",
            "色情低俗",
            "暴力血腥",
            "政治敏感",
            "违法犯罪",
            "谣言虚假",
            "侵犯隐私",
            "其他违规"
        );
        reasonCombo.setValue("垃圾广告");

        TextArea descriptionArea = new TextArea();
        descriptionArea.setPromptText("请详细描述举报原因...");
        descriptionArea.setPrefRowCount(4);

        content.getChildren().addAll(
            new Label("举报原因:"),
            reasonCombo,
            new Label("详细说明:"),
            descriptionArea
        );

        dialog.getDialogPane().setContent(content);
        dialog.getDialogPane().getButtonTypes().addAll(ButtonType.OK, ButtonType.CANCEL);

        dialog.showAndWait().ifPresent(result -> {
            if (result == ButtonType.OK) {
                String reason = reasonCombo.getValue();
                String description = descriptionArea.getText();
                submitReport("article", articleId, reason, description);
            }
        });
    }

    private void submitReport(String targetType, int targetId, String reason, String description) {
        ApiClient.submitReport(targetType, targetId, reason, description, new ApiCallback() {
            @Override
            public void onSuccess(String response) {
                javafx.application.Platform.runLater(() -> {
                    updateStatus("举报已提交，感谢您的反馈");
                });
            }

            @Override
            public void onError(String error) {
                javafx.application.Platform.runLater(() -> {
                    updateStatus("举报失败: " + error);
                });
            }
        });
    }

    private void showArticleDetail(int articleId) {
        showArticleDetailPublic(articleId);
    }

    public void showArticleDetailPublic(int articleId) {
        ApiClient.getArticleDetail(articleId, new ApiCallback() {
            @Override
            public void onSuccess(String response) {
                javafx.application.Platform.runLater(() -> {
                    try {
                        JSONObject json = new JSONObject(response);
                        showArticleDetailDialog(json);
                    } catch (Exception e) {
                        updateStatus("解析文章详情失败");
                    }
                });
            }

            @Override
            public void onError(String error) {
                javafx.application.Platform.runLater(() -> {
                    updateStatus("获取文章详情失败: " + error);
                });
            }
        });
    }

    private void showArticleDetailDialog(JSONObject articleJson) {
        Dialog<ButtonType> dialog = new Dialog<>();
        dialog.setTitle("文章详情");
        dialog.setHeaderText(articleJson.optString("title", ""));
        dialog.setWidth(800);

        ScrollPane scrollPane = new ScrollPane();
        VBox contentBox = new VBox();
        contentBox.setSpacing(16);
        contentBox.setPadding(new Insets(20));

        // 文章信息
        Label titleLabel = new Label(articleJson.optString("title", ""));
        titleLabel.getStyleClass().add("detail-title");

        JSONObject user = articleJson.optJSONObject("user");
        String author = user != null ? user.optString("nickname", user.optString("username", "匿名")) : "匿名";
        Label authorLabel = new Label("作者: " + author);
        authorLabel.getStyleClass().add("detail-author");

        String createdAt = articleJson.optString("created_at", "");
        Label dateLabel = new Label("发布时间: " + createdAt);
        dateLabel.getStyleClass().add("detail-date");

        Label contentLabel = new Label(articleJson.optString("content", ""));
        contentLabel.getStyleClass().add("detail-content");
        contentLabel.setWrapText(true);

        // 统计信息
        HBox statsBox = new HBox();
        statsBox.setSpacing(24);
        int viewCount = articleJson.optInt("view_count", 0);
        int likeCount = articleJson.optInt("like_count", 0);
        int commentCount = articleJson.optInt("comment_count", 0);
        statsBox.getChildren().addAll(
            new Label("浏览: " + viewCount),
            new Label("点赞: " + likeCount),
            new Label("评论: " + commentCount)
        );

        // 操作按钮
        HBox actionBox = new HBox();
        actionBox.setSpacing(12);

        Button likeBtn = new Button("点赞");
        likeBtn.getStyleClass().add("action-button");
        likeBtn.setOnAction(e -> likeArticle(articleJson.optInt("id", 0)));

        Button favoriteBtn = new Button("收藏");
        favoriteBtn.getStyleClass().add("action-button");
        favoriteBtn.setOnAction(e -> addFavorite(articleJson.optInt("id", 0)));

        Button reportBtn = new Button("举报");
        reportBtn.getStyleClass().add("action-button");
        reportBtn.setOnAction(e -> showReportDialog(articleJson.optInt("id", 0)));

        actionBox.getChildren().addAll(likeBtn, favoriteBtn, reportBtn);

        // 评论区标题
        Label commentTitle = new Label("评论");
        commentTitle.getStyleClass().add("comment-title");

        // 评论列表容器
        VBox commentListBox = new VBox();
        commentListBox.setSpacing(12);
        commentListBox.setId("commentListBox");

        // 加载评论
        loadCommentsForDetail(articleJson.optInt("id", 0), commentListBox);

        // 发表评论
        HBox commentInputBox = new HBox();
        commentInputBox.setSpacing(10);

        TextArea commentInput = new TextArea();
        commentInput.setPromptText("写下你的评论...");
        commentInput.setPrefRowCount(3);
        HBox.setHgrow(commentInput, Priority.ALWAYS);

        Button submitCommentBtn = new Button("发表");
        submitCommentBtn.getStyleClass().add("submit-button");
        submitCommentBtn.setOnAction(e -> {
            String commentContent = commentInput.getText();
            if (!commentContent.isEmpty()) {
                submitComment(articleJson.optInt("id", 0), commentContent, commentListBox);
                commentInput.clear();
            }
        });

        commentInputBox.getChildren().addAll(commentInput, submitCommentBtn);

        contentBox.getChildren().addAll(
            titleLabel, authorLabel, dateLabel, contentLabel,
            statsBox, actionBox,
            new Separator(),
            commentTitle, commentListBox, commentInputBox
        );

        scrollPane.setContent(contentBox);
        scrollPane.setFitToWidth(true);
        scrollPane.setPrefHeight(600);

        dialog.getDialogPane().setContent(scrollPane);
        dialog.getDialogPane().getButtonTypes().addAll(ButtonType.CLOSE);

        dialog.showAndWait();
    }

    private void loadCommentsForDetail(int articleId, VBox commentListBox) {
        ApiClient.getComments(articleId, 1, 50, new ApiCallback() {
            @Override
            public void onSuccess(String response) {
                javafx.application.Platform.runLater(() -> {
                    try {
                        JSONObject json = new JSONObject(response);
                        JSONArray comments = json.optJSONArray("comments");
                        if (comments != null) {
                            commentListBox.getChildren().clear();
                            for (int i = 0; i < comments.length(); i++) {
                                JSONObject comment = comments.getJSONObject(i);
                                VBox commentCard = createCommentCard(comment);
                                commentListBox.getChildren().add(commentCard);
                            }
                        }
                    } catch (Exception e) {
                        Label errorLabel = new Label("加载评论失败");
                        commentListBox.getChildren().add(errorLabel);
                    }
                });
            }

            @Override
            public void onError(String error) {
                javafx.application.Platform.runLater(() -> {
                    Label errorLabel = new Label("加载评论失败: " + error);
                    commentListBox.getChildren().add(errorLabel);
                });
            }
        });
    }

    private VBox createCommentCard(JSONObject comment) {
        VBox card = new VBox();
        card.setSpacing(8);
        card.setPadding(new Insets(12));
        card.getStyleClass().add("comment-card");

        JSONObject user = comment.optJSONObject("user");
        String author = user != null ? user.optString("nickname", user.optString("username", "匿名")) : "匿名";
        String content = comment.optString("content", "");
        String createdAt = comment.optString("created_at", "");

        Label authorLabel = new Label(author);
        authorLabel.getStyleClass().add("comment-author");

        Label contentLabel = new Label(content);
        contentLabel.getStyleClass().add("comment-content");
        contentLabel.setWrapText(true);

        Label dateLabel = new Label(createdAt);
        dateLabel.getStyleClass().add("comment-date");

        HBox actionBox = new HBox();
        actionBox.setSpacing(16);

        Button likeBtn = new Button("赞 " + comment.optInt("like_count", 0));
        likeBtn.getStyleClass().add("comment-action-button");

        actionBox.getChildren().add(likeBtn);

        card.getChildren().addAll(authorLabel, contentLabel, dateLabel, actionBox);

        return card;
    }

    private void submitComment(int articleId, String content, VBox commentListBox) {
        if (!isLoggedIn) {
            showLoginDialog();
            return;
        }

        ApiClient.createComment(articleId, content, null, new ApiCallback() {
            @Override
            public void onSuccess(String response) {
                javafx.application.Platform.runLater(() -> {
                    updateStatus("评论发表成功");
                    loadCommentsForDetail(articleId, commentListBox);
                });
            }

            @Override
            public void onError(String error) {
                javafx.application.Platform.runLater(() -> {
                    updateStatus("发表评论失败: " + error);
                });
            }
        });
    }

    private void addFavorite(int articleId) {
        if (!isLoggedIn) {
            showLoginDialog();
            return;
        }

        ApiClient.addFavorite(articleId, new ApiCallback() {
            @Override
            public void onSuccess(String response) {
                javafx.application.Platform.runLater(() -> {
                    updateStatus("收藏成功");
                });
            }

            @Override
            public void onError(String error) {
                javafx.application.Platform.runLater(() -> {
                    updateStatus("收藏失败: " + error);
                });
            }
        });
    }
}
