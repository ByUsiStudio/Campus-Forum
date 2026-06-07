package miao.byusi.pc.xylt;

import javafx.scene.control.Button;
import javafx.scene.control.Label;
import javafx.scene.layout.HBox;

public class ArticleCard extends HBox {

    private Label titleLabel;
    private Label authorLabel;
    private Label contentLabel;
    private Label statsLabel;
    private Button likeButton;
    private Button commentButton;
    private Button shareButton;

    public ArticleCard() {
        getStyleClass().add("article-card");
        setSpacing(10);
        setPrefWidth(800);

        titleLabel = new Label();
        titleLabel.getStyleClass().add("article-title");

        authorLabel = new Label();
        authorLabel.getStyleClass().add("article-author");

        contentLabel = new Label();
        contentLabel.getStyleClass().add("article-content");
        contentLabel.setWrapText(true);

        statsLabel = new Label();
        statsLabel.getStyleClass().add("article-stats");

        likeButton = new Button("点赞");
        commentButton = new Button("评论");
        shareButton = new Button("分享");

        getChildren().addAll(titleLabel, authorLabel, contentLabel, statsLabel, likeButton, commentButton, shareButton);
    }

    public void setTitle(String title) {
        titleLabel.setText(title);
    }

    public void setAuthor(String author) {
        authorLabel.setText(author);
    }

    public void setContent(String content) {
        contentLabel.setText(content);
    }

    public void setStats(int views, int likes, int comments) {
        statsLabel.setText("浏览: " + views + " | 点赞: " + likes + " | 评论: " + comments);
    }
}
