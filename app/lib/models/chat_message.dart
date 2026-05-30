import 'user.dart';

class ChatMessage {
  final int id;
  final int senderId;
  final int receiverId;
  final String content;
  final bool isRead;
  final String createdAt;
  final String senderName;

  ChatMessage({
    required this.id,
    required this.senderId,
    required this.receiverId,
    required this.content,
    required this.isRead,
    required this.createdAt,
    required this.senderName,
  });

  factory ChatMessage.fromJson(Map<String, dynamic> json) {
    return ChatMessage(
      id: json['id'] ?? 0,
      senderId: json['sender_id'] ?? json['senderId'] ?? 0,
      receiverId: json['receiver_id'] ?? json['receiverId'] ?? 0,
      content: json['content'] ?? '',
      isRead: json['is_read'] ?? json['isRead'] ?? false,
      createdAt: json['created_at'] ?? json['createdAt'] ?? '',
      senderName: json['sender_name'] ?? json['senderName'] ?? '',
    );
  }
}
