class User {
  final int id;
  final String username;
  final String displayName;
  final String avatar;
  final String role;
  final String qqNumber;
  final String createdAt;

  User({
    required this.id,
    required this.username,
    required this.displayName,
    required this.avatar,
    required this.role,
    required this.qqNumber,
    required this.createdAt,
  });

  factory User.fromJson(Map<String, dynamic> json) {
    return User(
      id: json['id'] ?? 0,
      username: json['username'] ?? '',
      displayName: json['display_name'] ?? json['displayName'] ?? '',
      avatar: json['avatar'] ?? '',
      role: json['role'] ?? '',
      qqNumber: json['qq_number'] ?? json['qqNumber'] ?? '',
      createdAt: json['created_at'] ?? json['createdAt'] ?? '',
    );
  }
}
