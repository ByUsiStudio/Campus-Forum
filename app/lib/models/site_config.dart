class SiteConfig {
  final int id;
  final String siteTitle;
  final String smtpHost;
  final int smtpPort;
  final String smtpUsername;
  final String smtpPassword;
  final String smtpFrom;
  final String smtpFromName;

  SiteConfig({
    required this.id,
    required this.siteTitle,
    required this.smtpHost,
    required this.smtpPort,
    required this.smtpUsername,
    required this.smtpPassword,
    required this.smtpFrom,
    required this.smtpFromName,
  });

  factory SiteConfig.fromJson(Map<String, dynamic> json) {
    return SiteConfig(
      id: json['id'] ?? 0,
      siteTitle: json['site_title'] ?? json['siteTitle'] ?? '',
      smtpHost: json['smtp_host'] ?? json['smtpHost'] ?? '',
      smtpPort: json['smtp_port'] ?? json['smtpPort'] ?? 0,
      smtpUsername: json['smtp_username'] ?? json['smtpUsername'] ?? '',
      smtpPassword: json['smtp_password'] ?? json['smtpPassword'] ?? '',
      smtpFrom: json['smtp_from'] ?? json['smtpFrom'] ?? '',
      smtpFromName: json['smtp_from_name'] ?? json['smtpFromName'] ?? '',
    );
  }
}
