class SidebarItem {
  final String title;
  final String link;
  final String icon;

  SidebarItem({
    required this.title,
    required this.link,
    required this.icon,
  });

  factory SidebarItem.fromJson(Map<String, dynamic> json) {
    return SidebarItem(
      title: json['title'] ?? '',
      link: json['link'] ?? '',
      icon: json['icon'] ?? '',
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'title': title,
      'link': link,
      'icon': icon,
    };
  }
}
