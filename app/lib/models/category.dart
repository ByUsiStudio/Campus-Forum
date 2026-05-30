class Category {
  final int id;
  final String name;
  final String description;
  final int sortOrder;

  Category({
    required this.id,
    required this.name,
    required this.description,
    required this.sortOrder,
  });

  factory Category.fromJson(Map<String, dynamic> json) {
    return Category(
      id: json['id'] ?? 0,
      name: json['name'] ?? '',
      description: json['description'] ?? '',
      sortOrder: json['sort_order'] ?? json['sortOrder'] ?? 0,
    );
  }
}
