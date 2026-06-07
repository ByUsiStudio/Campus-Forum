package miao.byusi.android.xylt;

public class Article {
    private int id;
    private String title;
    private String content;
    private String author;
    private String avatarUrl;
    private int viewCount;
    private int likeCount;
    private int commentCount;
    private String createdAt;
    private boolean isTop;

    public Article() {}

    public Article(int id, String title, String content, String author, 
                   String avatarUrl, int viewCount, int likeCount, int commentCount, String createdAt) {
        this.id = id;
        this.title = title;
        this.content = content;
        this.author = author;
        this.avatarUrl = avatarUrl;
        this.viewCount = viewCount;
        this.likeCount = likeCount;
        this.commentCount = commentCount;
        this.createdAt = createdAt;
        this.isTop = false;
    }

    public int getId() { return id; }
    public void setId(int id) { this.id = id; }
    public String getTitle() { return title; }
    public void setTitle(String title) { this.title = title; }
    public String getContent() { return content; }
    public void setContent(String content) { this.content = content; }
    public String getAuthor() { return author; }
    public void setAuthor(String author) { this.author = author; }
    public String getAvatarUrl() { return avatarUrl; }
    public void setAvatarUrl(String avatarUrl) { this.avatarUrl = avatarUrl; }
    public int getViewCount() { return viewCount; }
    public void setViewCount(int viewCount) { this.viewCount = viewCount; }
    public int getLikeCount() { return likeCount; }
    public void setLikeCount(int likeCount) { this.likeCount = likeCount; }
    public int getCommentCount() { return commentCount; }
    public void setCommentCount(int commentCount) { this.commentCount = commentCount; }
    public String getCreatedAt() { return createdAt; }
    public void setCreatedAt(String createdAt) { this.createdAt = createdAt; }
    public boolean isTop() { return isTop; }
    public void setTop(boolean top) { isTop = top; }
}
