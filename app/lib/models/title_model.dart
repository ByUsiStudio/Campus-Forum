class TitleModel {
  final int id;
  final String name;
  final String description;
  final String color;
  final String icon;
  final bool isActive;

  TitleModel({
    required this.id,
    required this.name,
    required this.description,
    required this.color,
    required this.icon,
    required this.isActive,
  });

  factory TitleModel.fromJson(Map<String, dynamic> json) {
    return TitleModel(
      id: json['id'] ?? 0,
      name: json['name'] ?? '',
      description: json['description'] ?? '',
      color: json['color'] ?? '#6750A4',
      icon: json['icon'] ?? '',
      isActive: json['is_active'] ?? json['isActive'] ?? false,
    );
  }
}
