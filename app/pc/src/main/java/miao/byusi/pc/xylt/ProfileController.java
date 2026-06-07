package miao.byusi.pc.xylt;

import javafx.fxml.FXML;
import javafx.geometry.Insets;
import javafx.scene.control.*;
import javafx.scene.layout.*;
import javafx.scene.text.Text;
import org.json.JSONArray;
import org.json.JSONObject;

import java.io.IOException;
import java.util.ArrayList;
import java.util.List;

public class ProfileController {

    @FXML
    private VBox rootBox;

    @FXML
    private Label usernameLabel;

    @FXML
    private Label nicknameLabel;

    @FXML
    private Label bioLabel;

    @FXML
    private Label articleCountLabel;

    @FXML
    private Label followCountLabel;

    @FXML
    private Label followerCountLabel;

    @FXML
    private TabPane contentTabPane;

    @FXML
    private VBox articlesList;

    @FXML
    private VBox favoritesList;

    @FXML
    private VBox followingList;

    @FXML
    private VBox notificationsList;

    private int currentPage = 1;
    private final int pageSize = 20;

    @FXML
    private void initialize() {
        loadProfile();
    }

    private void loadProfile() {
        ApiClient.getProfile(new ApiCallback() {
            @Override
            public void onSuccess(String response) {
                javafx.application.Platform.runLater(() -> {
                    try {
                        JSONObject json = new JSONObject(response);
                        updateProfileUI(json);
                        loadMyArticles();
                    } catch (Exception e) {
                        showError("加载用户信息失败");
                    }
                });
            }

            @Override
            public void onError(String error) {
                javafx.application.Platform.runLater(() -> {
                    showError("获取用户信息失败: " + error);
                });
            }
        });
    }

    private void updateProfileUI(JSONObject json) {
        String username = json.optString("username", "");
        String displayName = json.optString("display_name", username);
        String signature = json.optString("signature", "");
        String avatarUrl = json.optString("avatar", "");

        usernameLabel.setText(username);
        nicknameLabel.setText(displayName);
        bioLabel.setText(signature.isEmpty() ? "暂无个人简介" : signature);

        int articles = json.optInt("article_count", 0);
        int follows = json.optInt("following_count", 0);
        int followers = json.optInt("follower_count", 0);

        articleCountLabel.setText(String.valueOf(articles));
        followCountLabel.setText(String.valueOf(follows));
        followerCountLabel.setText(String.valueOf(followers));
    }

    @FXML
    private void loadMyArticles() {
        ApiClient.getMyArticles(currentPage, pageSize, new ApiCallback() {
            @Override
            public void onSuccess(String response) {
                javafx.application.Platform.runLater(() -> {
                    try {
                        JSONObject json = new JSONObject(response);
                        JSONArray articles = json.optJSONArray("articles");
                        articlesList.getChildren().clear();

                        if (articles != null && articles.length() > 0) {
                            for (int i = 0; i < articles.length(); i++) {
                                JSONObject article = articles.getJSONObject(i);
                                VBox card = createArticleCard(article);
                                articlesList.getChildren().add(card);
                            }
                        } else {
                            articlesList.getChildren().add(new Label("暂无文章"));
                        }
                    } catch (Exception e) {
                        articlesList.getChildren().add(new Label("加载文章失败"));
                    }
                });
            }

            @Override
            public void onError(String error) {
                javafx.application.Platform.runLater(() -> {
                    articlesList.getChildren().add(new Label("加载文章失败: " + error));
                });
            }
        });
    }

    @FXML
    private void loadFavorites() {
        ApiClient.getFavorites(currentPage, pageSize, new ApiCallback() {
            @Override
            public void onSuccess(String response) {
                javafx.application.Platform.runLater(() -> {
                    try {
                        JSONObject json = new JSONObject(response);
                        JSONArray articles = json.optJSONArray("articles");
                        favoritesList.getChildren().clear();

                        if (articles != null && articles.length() > 0) {
                            for (int i = 0; i < articles.length(); i++) {
                                JSONObject article = articles.getJSONObject(i);
                                VBox card = createArticleCard(article);
                                favoritesList.getChildren().add(card);
                            }
                        } else {
                            favoritesList.getChildren().add(new Label("暂无收藏"));
                        }
                    } catch (Exception e) {
                        favoritesList.getChildren().add(new Label("加载收藏失败"));
                    }
                });
            }

            @Override
            public void onError(String error) {
                javafx.application.Platform.runLater(() -> {
                    favoritesList.getChildren().add(new Label("加载收藏失败: " + error));
                });
            }
        });
    }

    @FXML
    private void loadFollowing() {
        ApiClient.getFollowingList(currentPage, pageSize, new ApiCallback() {
            @Override
            public void onSuccess(String response) {
                javafx.application.Platform.runLater(() -> {
                    try {
                        JSONObject json = new JSONObject(response);
                        JSONArray users = json.optJSONArray("users");
                        followingList.getChildren().clear();

                        if (users != null && users.length() > 0) {
                            for (int i = 0; i < users.length(); i++) {
                                JSONObject user = users.getJSONObject(i);
                                VBox card = createUserCard(user);
                                followingList.getChildren().add(card);
                            }
                        } else {
                            followingList.getChildren().add(new Label("暂无关注"));
                        }
                    } catch (Exception e) {
                        followingList.getChildren().add(new Label("加载关注列表失败"));
                    }
                });
            }

            @Override
            public void onError(String error) {
                javafx.application.Platform.runLater(() -> {
                    followingList.getChildren().add(new Label("加载关注列表失败: " + error));
                });
            }
        });
    }

    @FXML
    private void loadNotifications() {
        ApiClient.getNotifications(currentPage, pageSize, new ApiCallback() {
            @Override
            public void onSuccess(String response) {
                javafx.application.Platform.runLater(() -> {
                    try {
                        JSONObject json = new JSONObject(response);
                        JSONArray notifications = json.optJSONArray("notifications");
                        notificationsList.getChildren().clear();

                        if (notifications != null && notifications.length() > 0) {
                            for (int i = 0; i < notifications.length(); i++) {
                                JSONObject notification = notifications.getJSONObject(i);
                                VBox card = createNotificationCard(notification);
                                notificationsList.getChildren().add(card);
                            }
                        } else {
                            notificationsList.getChildren().add(new Label("暂无通知"));
                        }
                    } catch (Exception e) {
                        notificationsList.getChildren().add(new Label("加载通知失败"));
                    }
                });
            }

            @Override
            public void onError(String error) {
                javafx.application.Platform.runLater(() -> {
                    notificationsList.getChildren().add(new Label("加载通知失败: " + error));
                });
            }
        });
    }

    private VBox createArticleCard(JSONObject article) {
        VBox card = new VBox();
        card.setSpacing(8);
        card.setPadding(new Insets(12));
        card.getStyleClass().add("profile-article-card");

        Label titleLabel = new Label(article.optString("title", ""));
        titleLabel.getStyleClass().add("article-title");

        Label contentLabel = new Label(article.optString("content", ""));
        contentLabel.getStyleClass().add("article-content");
        contentLabel.setWrapText(true);
        contentLabel.setMaxWidth(600);

        HBox statsBox = new HBox();
        statsBox.setSpacing(16);
        statsBox.getChildren().addAll(
            new Label("浏览: " + article.optInt("view_count", 0)),
            new Label("点赞: " + article.optInt("like_count", 0)),
            new Label("评论: " + article.optInt("comment_count", 0))
        );

        card.getChildren().addAll(titleLabel, contentLabel, statsBox);

        card.setOnMouseClicked(event -> {
            if (event.getClickCount() == 2) {
                int articleId = article.optInt("id", 0);
                if (articleId > 0) {
                    MainController mainController = MainApp.getMainController();
                    if (mainController != null) {
                        mainController.showArticleDetail(articleId);
                    }
                }
            }
        });

        return card;
    }

    private VBox createUserCard(JSONObject user) {
        VBox card = new VBox();
        card.setSpacing(8);
        card.setPadding(new Insets(12));
        card.getStyleClass().add("profile-user-card");

        String nickname = user.optString("nickname", user.optString("username", "用户"));
        String bio = user.optString("bio", "");

        Label nameLabel = new Label(nickname);
        nameLabel.getStyleClass().add("user-name");

        Label bioLabel = new Label(bio.isEmpty() ? "暂无个人简介" : bio);
        bioLabel.getStyleClass().add("user-bio");
        bioLabel.setWrapText(true);

        card.getChildren().addAll(nameLabel, bioLabel);

        return card;
    }

    private VBox createNotificationCard(JSONObject notification) {
        VBox card = new VBox();
        card.setSpacing(8);
        card.setPadding(new Insets(12));
        card.getStyleClass().add("profile-notification-card");

        String type = notification.optString("type", "system");
        String content = notification.optString("content", "");
        String createdAt = notification.optString("created_at", "");
        boolean isRead = notification.optBoolean("is_read", false);

        Label typeLabel = new Label(getNotificationTypeName(type));
        typeLabel.getStyleClass().add("notification-type");

        if (!isRead) {
            typeLabel.setStyle("-fx-font-weight: bold;");
        }

        Label contentLabel = new Label(content);
        contentLabel.getStyleClass().add("notification-content");
        contentLabel.setWrapText(true);

        Label dateLabel = new Label(createdAt);
        dateLabel.getStyleClass().add("notification-date");

        card.getChildren().addAll(typeLabel, contentLabel, dateLabel);

        return card;
    }

    private String getNotificationTypeName(String type) {
        switch (type) {
            case "like":
                return "点赞通知";
            case "comment":
                return "评论通知";
            case "follow":
                return "关注通知";
            case "system":
                return "系统通知";
            default:
                return "通知";
        }
    }

    @FXML
    private void showEditProfileDialog() {
        Dialog<ButtonType> dialog = new Dialog<>();
        dialog.setTitle("编辑资料");
        dialog.setHeaderText("修改个人资料");

        GridPane grid = new GridPane();
        grid.setHgap(10);
        grid.setVgap(10);
        grid.setPadding(new Insets(20));

        TextField displayNameField = new TextField();
        displayNameField.setPromptText("显示名称");
        displayNameField.setText(nicknameLabel.getText());

        TextArea signatureField = new TextArea();
        signatureField.setPromptText("个人简介");
        signatureField.setText(bioLabel.getText());
        signatureField.setPrefRowCount(3);

        grid.add(new Label("显示名称:"), 0, 0);
        grid.add(displayNameField, 1, 0);
        grid.add(new Label("简介:"), 0, 1);
        grid.add(signatureField, 1, 1);

        dialog.getDialogPane().setContent(grid);
        dialog.getDialogPane().getButtonTypes().addAll(ButtonType.OK, ButtonType.CANCEL);

        dialog.showAndWait().ifPresent(result -> {
            if (result == ButtonType.OK) {
                String displayName = displayNameField.getText();
                String signature = signatureField.getText();
                updateProfile(displayName, signature);
            }
        });
    }

    private void updateProfile(String displayName, String signature) {
        ApiClient.updateProfile(displayName, signature, new ApiCallback() {
            @Override
            public void onSuccess(String response) {
                javafx.application.Platform.runLater(() -> {
                    nicknameLabel.setText(displayName);
                    bioLabel.setText(signature);
                    showSuccess("资料更新成功");
                });
            }

            @Override
            public void onError(String error) {
                javafx.application.Platform.runLater(() -> {
                    showError("更新失败: " + error);
                });
            }
        });
    }

    @FXML
    private void refreshData() {
        loadProfile();
    }

    private void showError(String message) {
        Alert alert = new Alert(Alert.AlertType.ERROR);
        alert.setTitle("错误");
        alert.setHeaderText(null);
        alert.setContentText(message);
        alert.showAndWait();
    }

    private void showSuccess(String message) {
        Alert alert = new Alert(Alert.AlertType.INFORMATION);
        alert.setTitle("成功");
        alert.setHeaderText(null);
        alert.setContentText(message);
        alert.showAndWait();
    }
}
