class FollowStatus {
  final bool isFollowing;
  final bool isFollowed;
  final bool mutual;

  FollowStatus({
    required this.isFollowing,
    required this.isFollowed,
    required this.mutual,
  });

  factory FollowStatus.fromJson(Map<String, dynamic> json) {
    return FollowStatus(
      isFollowing: json['is_following'] ?? json['isFollowing'] ?? false,
      isFollowed: json['is_followed'] ?? json['isFollowed'] ?? false,
      mutual: json['mutual'] ?? false,
    );
  }
}
