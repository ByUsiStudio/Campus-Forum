import 'user.dart';

class ChatSession {
  final int sessionId;
  final User otherUser;
  final int unreadCount;
  final String lastMessage;
  final String lastMessageAt;

  ChatSession({
    required this.sessionId,
    required this.otherUser,
    required this.unreadCount,
    required this.lastMessage,
    required this.lastMessageAt,
  });

  factory ChatSession.fromJson(Map<String, dynamic> json) {
    return ChatSession(
      sessionId: json['session_id'] ?? 0,
      otherUser: User.fromJson(json['other_user'] as Map<String, dynamic>),
      unreadCount: json['unread_count'] ?? json['unreadCount'] ?? 0,
      lastMessage: json['last_message'] ?? json['lastMessage'] ?? '',
      lastMessageAt: json['last_message_at'] ?? json['lastMessageAt'] ?? '',
    );
  }
}
