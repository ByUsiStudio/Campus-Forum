class NotificationModel {
  final int id;
  final String type;
  final String title;
  final String content;
  final String target;
  final bool isRead;
  final String createdAt;

  NotificationModel({
    required this.id,
    required this.type,
    required this.title,
    required this.content,
    required this.target,
    required this.isRead,
    required this.createdAt,
  });

  factory NotificationModel.fromJson(Map<String, dynamic> json) {
    return NotificationModel(
      id: json['id'] ?? 0,
      type: json['type'] ?? '',
      title: json['title'] ?? '',
      content: json['content'] ?? '',
      target: json['target'] ?? '',
      isRead: json['is_read'] ?? json['isRead'] ?? false,
      createdAt: json['created_at'] ?? json['createdAt'] ?? '',
    );
  }
}
