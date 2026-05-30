import 'user.dart';

class Comment {
  final int id;
  final String content;
  final int userId;
  final User? user;
  final int articleId;
  final int? parentId;
  final int likeCount;
  final int replyCount;
  final List<Comment> replies;
  final String createdAt;

  Comment({
    required this.id,
    required this.content,
    required this.userId,
    this.user,
    required this.articleId,
    this.parentId,
    required this.likeCount,
    required this.replyCount,
    required this.replies,
    required this.createdAt,
  });

  factory Comment.fromJson(Map<String, dynamic> json) {
    return Comment(
      id: json['id'] ?? 0,
      content: json['content'] ?? '',
      userId: json['user_id'] ?? json['userId'] ?? 0,
      user: json['user'] != null ? User.fromJson(json['user'] as Map<String, dynamic>) : null,
      articleId: json['article_id'] ?? json['articleId'] ?? 0,
      parentId: json['parent_id'] != null ? (json['parent_id'] as num).toInt() : null,
      likeCount: json['like_count'] ?? json['likeCount'] ?? 0,
      replyCount: json['reply_count'] ?? json['replyCount'] ?? 0,
      replies: (json['replies'] as List<dynamic>?)
              ?.map((item) => Comment.fromJson(item as Map<String, dynamic>))
              .toList() ??
          [],
      createdAt: json['created_at'] ?? json['createdAt'] ?? '',
    );
  }
}
