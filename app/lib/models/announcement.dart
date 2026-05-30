class Announcement {
  final String content;
  final String contentHtml;

  Announcement({
    required this.content,
    required this.contentHtml,
  });

  factory Announcement.fromJson(Map<String, dynamic> json) {
    return Announcement(
      content: json['content'] ?? '',
      contentHtml: json['content_html'] ?? json['contentHtml'] ?? '',
    );
  }
}
