class Article {
  final int id;
  final String title;
  final String content;
  final String contentHtml;
  final Map<String, dynamic>? user;
  final Map<String, dynamic>? category;
  final int viewCount;
  final int likeCount;
  final int commentCount;
  final String createdAt;

  Article({
    required this.id,
    required this.title,
    required this.content,
    required this.contentHtml,
    this.user,
    this.category,
    required this.viewCount,
    required this.likeCount,
    required this.commentCount,
    required this.createdAt,
  });

  factory Article.fromJson(Map<String, dynamic> json) {
    return Article(
      id: json['id'] ?? 0,
      title: json['title'] ?? '',
      content: json['content'] ?? '',
      contentHtml: json['content_html'] ?? json['contentHtml'] ?? '',
      user: json['user'] as Map<String, dynamic>?,
      category: json['category'] as Map<String, dynamic>?,
      viewCount: json['view_count'] ?? json['viewCount'] ?? 0,
      likeCount: json['like_count'] ?? json['likeCount'] ?? 0,
      commentCount: json['comment_count'] ?? json['commentCount'] ?? 0,
      createdAt: json['created_at'] ?? json['createdAt'] ?? '',
    );
  }
}
