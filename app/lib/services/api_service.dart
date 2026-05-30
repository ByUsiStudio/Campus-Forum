import 'dart:convert';
import 'dart:io';

import 'package:http/http.dart' as http;
import 'package:shared_preferences/shared_preferences.dart';

import '../models/article.dart';
import '../models/announcement.dart';
import '../models/chat_message.dart';
import '../models/chat_session.dart';
import '../models/category.dart';
import '../models/comment.dart';
import '../models/follow_status.dart';
import '../models/notification_model.dart';
import '../models/site_config.dart';
import '../models/sidebar_item.dart';
import '../models/title_model.dart';
import '../models/user.dart';

class ApiException implements Exception {
  final String message;
  ApiException(this.message);
  @override
  String toString() => message;
}

class ApiService {
  static const String baseUrl = 'https://xylt.cdifit.com';
  static const String _tokenKey = 'forum_token';
  static SharedPreferences? _prefs;
  static String? token;
  static User? currentUser;

  static Future<void> init() async {
    _prefs = await SharedPreferences.getInstance();
    token = _prefs?.getString(_tokenKey);
    if (isLoggedIn) {
      try {
        await fetchProfile();
      } catch (_) {
        await logout();
      }
    }
  }

  static bool get isLoggedIn => token != null && token!.isNotEmpty;

  static Map<String, String> _headers({bool auth = false}) {
    final headers = <String, String>{
      'Content-Type': 'application/json',
    };
    if (auth && isLoggedIn) {
      headers['Authorization'] = 'Bearer $token';
    }
    return headers;
  }

  static String _parseError(http.Response response) {
    try {
      final body = jsonDecode(response.body);
      if (body is Map<String, dynamic>) {
        return body['error']?.toString() ?? body['message']?.toString() ?? response.reasonPhrase ?? '请求失败';
      }
    } catch (_) {}
    return response.reasonPhrase ?? '请求失败';
  }

  static Future<void> _saveToken(String newToken) async {
    token = newToken;
    await _prefs?.setString(_tokenKey, newToken);
  }

  static Future<void> logout() async {
    token = null;
    currentUser = null;
    await _prefs?.remove(_tokenKey);
  }

  static Future<bool> login(String username, String password) async {
    final uri = Uri.parse('$baseUrl/api/auth/login');
    final response = await http.post(
      uri,
      headers: _headers(),
      body: jsonEncode({
        'username': username,
        'password': password,
      }),
    );

    if (response.statusCode == 200) {
      final body = jsonDecode(response.body) as Map<String, dynamic>;
      final loginToken = body['token'] as String?;
      final userData = body['user'] as Map<String, dynamic>?;
      if (loginToken != null) {
        await _saveToken(loginToken);
      }
      if (userData != null) {
        currentUser = User.fromJson(userData);
      }
      return loginToken != null;
    }

    throw ApiException(_parseError(response));
  }

  static Future<bool> register(String username, String qqNumber, String displayName, String password) async {
    final uri = Uri.parse('$baseUrl/api/auth/register');
    final response = await http.post(
      uri,
      headers: _headers(),
      body: jsonEncode({
        'username': username,
        'qq_number': qqNumber,
        'display_name': displayName,
        'password': password,
      }),
    );

    if (response.statusCode == 200 || response.statusCode == 201) {
      return true;
    }
    throw ApiException(_parseError(response));
  }

  static Future<bool> sendResetCode(String qqNumber) async {
    final uri = Uri.parse('$baseUrl/api/password/reset-code');
    final response = await http.post(
      uri,
      headers: _headers(),
      body: jsonEncode({'qq_number': qqNumber}),
    );
    if (response.statusCode == 200) {
      return true;
    }
    throw ApiException(_parseError(response));
  }

  static Future<bool> resetPassword(String qqNumber, String code, String password) async {
    final uri = Uri.parse('$baseUrl/api/password/reset');
    final response = await http.post(
      uri,
      headers: _headers(),
      body: jsonEncode({
        'qq_number': qqNumber,
        'code': code,
        'password': password,
      }),
    );
    if (response.statusCode == 200) {
      return true;
    }
    throw ApiException(_parseError(response));
  }

  static Future<bool> checkInit() async {
    final uri = Uri.parse('$baseUrl/api/auth/check-init');
    final response = await http.get(uri, headers: _headers());
    if (response.statusCode == 200) {
      final body = jsonDecode(response.body) as Map<String, dynamic>;
      return body['initialized'] == true;
    }
    throw ApiException(_parseError(response));
  }

  static Future<List<Category>> fetchCategories() async {
    final uri = Uri.parse('$baseUrl/api/categories');
    final response = await http.get(uri, headers: _headers());
    if (response.statusCode == 200) {
      final body = jsonDecode(response.body) as Map<String, dynamic>;
      final categories = body['categories'] as List<dynamic>? ?? [];
      return categories
          .map((dynamic item) => Category.fromJson(item as Map<String, dynamic>))
          .toList();
    }
    throw ApiException(_parseError(response));
  }

  static Future<Category> createCategory(String name, {String description = '', int sortOrder = 0}) async {
    final uri = Uri.parse('$baseUrl/api/categories');
    final response = await http.post(
      uri,
      headers: _headers(auth: true),
      body: jsonEncode({
        'name': name,
        'description': description,
        'sort_order': sortOrder,
      }),
    );
    if (response.statusCode == 200) {
      final body = jsonDecode(response.body) as Map<String, dynamic>;
      return Category.fromJson(body['category'] as Map<String, dynamic>);
    }
    throw ApiException(_parseError(response));
  }

  static Future<bool> updateCategory(int categoryId, {String? name, String? description, int? sortOrder}) async {
    final uri = Uri.parse('$baseUrl/api/categories/$categoryId');
    final requestBody = <String, dynamic>{};
    if (name != null) requestBody['name'] = name;
    if (description != null) requestBody['description'] = description;
    if (sortOrder != null) requestBody['sort_order'] = sortOrder;

    final response = await http.put(
      uri,
      headers: _headers(auth: true),
      body: jsonEncode(requestBody),
    );
    if (response.statusCode == 200) {
      return true;
    }
    throw ApiException(_parseError(response));
  }

  static Future<bool> deleteCategory(int categoryId) async {
    final uri = Uri.parse('$baseUrl/api/categories/$categoryId');
    final response = await http.delete(uri, headers: _headers(auth: true));
    if (response.statusCode == 200) {
      return true;
    }
    throw ApiException(_parseError(response));
  }

  static Future<List<Article>> fetchArticles({int page = 1, int pageSize = 20, int? categoryId}) async {
    final query = <String, String>{
      'page': page.toString(),
      'page_size': pageSize.toString(),
    };
    if (categoryId != null) {
      query['category_id'] = categoryId.toString();
    }
    final uri = Uri.parse('$baseUrl/api/articles').replace(queryParameters: query);
    final response = await http.get(uri, headers: _headers());
    if (response.statusCode == 200) {
      final body = jsonDecode(response.body) as Map<String, dynamic>;
      final articles = body['articles'] as List<dynamic>? ?? [];
      return articles
          .map((dynamic item) => Article.fromJson(item as Map<String, dynamic>))
          .toList();
    }
    throw ApiException(_parseError(response));
  }

  static Future<Article> fetchArticleDetail(int id) async {
    final uri = Uri.parse('$baseUrl/api/articles/$id');
    final response = await http.get(uri, headers: _headers());
    if (response.statusCode == 200) {
      final body = jsonDecode(response.body) as Map<String, dynamic>;
      final articleJson = body['article'] as Map<String, dynamic>? ?? body;
      return Article.fromJson(articleJson);
    }
    throw ApiException(_parseError(response));
  }

  static Future<bool> createArticle(String title, String content, int categoryId) async {
    final uri = Uri.parse('$baseUrl/api/articles');
    final response = await http.post(
      uri,
      headers: _headers(auth: true),
      body: jsonEncode({
        'title': title,
        'content': content,
        'category_id': categoryId,
      }),
    );
    if (response.statusCode == 200 || response.statusCode == 201) {
      return true;
    }
    throw ApiException(_parseError(response));
  }

  static Future<bool> updateArticle(int articleId, {String? title, String? content, int? categoryId}) async {
    final uri = Uri.parse('$baseUrl/api/articles/$articleId');
    final requestBody = <String, dynamic>{};
    if (title != null) requestBody['title'] = title;
    if (content != null) requestBody['content'] = content;
    if (categoryId != null) requestBody['category_id'] = categoryId;
    final response = await http.put(
      uri,
      headers: _headers(auth: true),
      body: jsonEncode(requestBody),
    );
    if (response.statusCode == 200) {
      return true;
    }
    throw ApiException(_parseError(response));
  }

  static Future<bool> deleteArticle(int articleId, {String? reason}) async {
    final uri = Uri.parse('$baseUrl/api/articles/$articleId');
    if (reason == null) {
      final response = await http.delete(uri, headers: _headers(auth: true));
      if (response.statusCode == 200) {
        return true;
      }
      throw ApiException(_parseError(response));
    }
    final response = await http.delete(
      uri,
      headers: _headers(auth: true),
      body: jsonEncode({'reason': reason}),
    );
    if (response.statusCode == 200) {
      return true;
    }
    throw ApiException(_parseError(response));
  }

  static Future<int> likeArticle(int articleId) async {
    final uri = Uri.parse('$baseUrl/api/articles/$articleId/like');
    final response = await http.post(uri, headers: _headers(auth: true));
    if (response.statusCode == 200) {
      final body = jsonDecode(response.body) as Map<String, dynamic>;
      return body['like_count'] as int? ?? 0;
    }
    throw ApiException(_parseError(response));
  }

  static Future<int> unlikeArticle(int articleId) async {
    final uri = Uri.parse('$baseUrl/api/articles/$articleId/like');
    final response = await http.delete(uri, headers: _headers(auth: true));
    if (response.statusCode == 200) {
      final body = jsonDecode(response.body) as Map<String, dynamic>;
      return body['like_count'] as int? ?? 0;
    }
    throw ApiException(_parseError(response));
  }

  static Future<Comment> createComment(int articleId, String content, {int? parentId}) async {
    final uri = Uri.parse('$baseUrl/api/articles/$articleId/comments');
    final response = await http.post(
      uri,
      headers: _headers(auth: true),
      body: jsonEncode({
        'content': content,
        if (parentId != null) 'parent_id': parentId,
      }),
    );
    if (response.statusCode == 200) {
      final body = jsonDecode(response.body) as Map<String, dynamic>;
      return Comment.fromJson(body['comment'] as Map<String, dynamic>);
    }
    throw ApiException(_parseError(response));
  }

  static Future<bool> deleteComment(int commentId) async {
    final uri = Uri.parse('$baseUrl/api/comments/$commentId');
    final response = await http.delete(uri, headers: _headers(auth: true));
    if (response.statusCode == 200) {
      return true;
    }
    throw ApiException(_parseError(response));
  }

  static Future<bool> likeComment(int commentId) async {
    final uri = Uri.parse('$baseUrl/api/comments/$commentId/like');
    final response = await http.post(uri, headers: _headers(auth: true));
    if (response.statusCode == 200) {
      return true;
    }
    throw ApiException(_parseError(response));
  }

  static Future<bool> unlikeComment(int commentId) async {
    final uri = Uri.parse('$baseUrl/api/comments/$commentId/like');
    final response = await http.delete(uri, headers: _headers(auth: true));
    if (response.statusCode == 200) {
      return true;
    }
    throw ApiException(_parseError(response));
  }

  static Future<List<Article>> fetchFavorites() async {
    final uri = Uri.parse('$baseUrl/api/favorites');
    final response = await http.get(uri, headers: _headers(auth: true));
    if (response.statusCode == 200) {
      final body = jsonDecode(response.body) as Map<String, dynamic>;
      final articles = body['articles'] as List<dynamic>? ?? [];
      return articles
          .map((dynamic item) => Article.fromJson(item as Map<String, dynamic>))
          .toList();
    }
    throw ApiException(_parseError(response));
  }

  static Future<bool> addFavorite(int articleId) async {
    final uri = Uri.parse('$baseUrl/api/articles/$articleId/favorite');
    final response = await http.post(uri, headers: _headers(auth: true));
    if (response.statusCode == 200) {
      return true;
    }
    throw ApiException(_parseError(response));
  }

  static Future<bool> removeFavorite(int articleId) async {
    final uri = Uri.parse('$baseUrl/api/articles/$articleId/favorite');
    final response = await http.delete(uri, headers: _headers(auth: true));
    if (response.statusCode == 200) {
      return true;
    }
    throw ApiException(_parseError(response));
  }

  static Future<bool> checkFavorite(int articleId) async {
    final uri = Uri.parse('$baseUrl/api/articles/$articleId/favorite/check');
    final response = await http.get(uri, headers: _headers(auth: true));
    if (response.statusCode == 200) {
      final body = jsonDecode(response.body) as Map<String, dynamic>;
      return body['favorited'] == true;
    }
    throw ApiException(_parseError(response));
  }

  static Future<bool> followUser(int userId) async {
    final uri = Uri.parse('$baseUrl/api/follow/$userId');
    final response = await http.post(uri, headers: _headers(auth: true));
    if (response.statusCode == 200) {
      return true;
    }
    throw ApiException(_parseError(response));
  }

  static Future<bool> unfollowUser(int userId) async {
    final uri = Uri.parse('$baseUrl/api/follow/$userId');
    final response = await http.delete(uri, headers: _headers(auth: true));
    if (response.statusCode == 200) {
      return true;
    }
    throw ApiException(_parseError(response));
  }

  static Future<List<User>> fetchFollowing() async {
    final uri = Uri.parse('$baseUrl/api/following');
    final response = await http.get(uri, headers: _headers(auth: true));
    if (response.statusCode == 200) {
      final body = jsonDecode(response.body) as Map<String, dynamic>;
      final following = body['following'] as List<dynamic>? ?? [];
      return following
          .map((dynamic item) => User.fromJson(item as Map<String, dynamic>))
          .toList();
    }
    throw ApiException(_parseError(response));
  }

  static Future<List<User>> fetchFollowers() async {
    final uri = Uri.parse('$baseUrl/api/followers');
    final response = await http.get(uri, headers: _headers(auth: true));
    if (response.statusCode == 200) {
      final body = jsonDecode(response.body) as Map<String, dynamic>;
      final followers = body['followers'] as List<dynamic>? ?? [];
      return followers
          .map((dynamic item) => User.fromJson(item as Map<String, dynamic>))
          .toList();
    }
    throw ApiException(_parseError(response));
  }

  static Future<FollowStatus> checkFollowStatus(int userId) async {
    final uri = Uri.parse('$baseUrl/api/follow/status/$userId');
    final response = await http.get(uri, headers: _headers(auth: true));
    if (response.statusCode == 200) {
      final body = jsonDecode(response.body) as Map<String, dynamic>;
      return FollowStatus.fromJson(body);
    }
    throw ApiException(_parseError(response));
  }

  static Future<List<User>> fetchMutualFriends() async {
    final uri = Uri.parse('$baseUrl/api/mutual');
    final response = await http.get(uri, headers: _headers(auth: true));
    if (response.statusCode == 200) {
      final body = jsonDecode(response.body) as Map<String, dynamic>;
      final mutual = body['mutual'] as List<dynamic>? ?? [];
      return mutual
          .map((dynamic item) => User.fromJson(item as Map<String, dynamic>))
          .toList();
    }
    throw ApiException(_parseError(response));
  }

  static Future<List<NotificationModel>> fetchNotifications() async {
    final uri = Uri.parse('$baseUrl/api/notifications');
    final response = await http.get(uri, headers: _headers(auth: true));
    if (response.statusCode == 200) {
      final body = jsonDecode(response.body) as Map<String, dynamic>;
      final notifications = body['notifications'] as List<dynamic>? ?? [];
      return notifications
          .map((dynamic item) => NotificationModel.fromJson(item as Map<String, dynamic>))
          .toList();
    }
    throw ApiException(_parseError(response));
  }

  static Future<int> fetchUnreadNotificationCount() async {
    final uri = Uri.parse('$baseUrl/api/notifications/unread-count');
    final response = await http.get(uri, headers: _headers(auth: true));
    if (response.statusCode == 200) {
      final body = jsonDecode(response.body) as Map<String, dynamic>;
      return body['unread_count'] as int? ?? 0;
    }
    throw ApiException(_parseError(response));
  }

  static Future<bool> markNotificationRead(int id) async {
    final uri = Uri.parse('$baseUrl/api/notifications/$id/read');
    final response = await http.post(uri, headers: _headers(auth: true));
    if (response.statusCode == 200) {
      return true;
    }
    throw ApiException(_parseError(response));
  }

  static Future<bool> markAllNotificationsRead() async {
    final uri = Uri.parse('$baseUrl/api/notifications/read-all');
    final response = await http.post(uri, headers: _headers(auth: true));
    if (response.statusCode == 200) {
      return true;
    }
    throw ApiException(_parseError(response));
  }

  static Future<List<NotificationModel>> fetchCommentReplyNotifications() async {
    final uri = Uri.parse('$baseUrl/api/comment-reply-notifications');
    final response = await http.get(uri, headers: _headers(auth: true));
    if (response.statusCode == 200) {
      final body = jsonDecode(response.body) as Map<String, dynamic>;
      final notifications = body['notifications'] as List<dynamic>? ?? [];
      return notifications
          .map((dynamic item) => NotificationModel.fromJson(item as Map<String, dynamic>))
          .toList();
    }
    throw ApiException(_parseError(response));
  }

  static Future<bool> markCommentReplyNotificationRead(int id) async {
    final uri = Uri.parse('$baseUrl/api/comment-reply-notifications/$id/read');
    final response = await http.post(uri, headers: _headers(auth: true));
    if (response.statusCode == 200) {
      return true;
    }
    throw ApiException(_parseError(response));
  }

  static Future<bool> markAllCommentReplyNotificationsRead() async {
    final uri = Uri.parse('$baseUrl/api/comment-reply-notifications/read-all');
    final response = await http.post(uri, headers: _headers(auth: true));
    if (response.statusCode == 200) {
      return true;
    }
    throw ApiException(_parseError(response));
  }

  static Future<List<ChatSession>> fetchChatSessions() async {
    final uri = Uri.parse('$baseUrl/api/chat/sessions');
    final response = await http.get(uri, headers: _headers(auth: true));
    if (response.statusCode == 200) {
      final body = jsonDecode(response.body) as Map<String, dynamic>;
      final sessions = body['sessions'] as List<dynamic>? ?? [];
      return sessions
          .map((dynamic item) => ChatSession.fromJson(item as Map<String, dynamic>))
          .toList();
    }
    throw ApiException(_parseError(response));
  }

  static Future<List<ChatMessage>> fetchChatMessages(int otherUserId) async {
    final uri = Uri.parse('$baseUrl/api/chat/messages/$otherUserId');
    final response = await http.get(uri, headers: _headers(auth: true));
    if (response.statusCode == 200) {
      final body = jsonDecode(response.body) as Map<String, dynamic>;
      final messages = body['messages'] as List<dynamic>? ?? [];
      return messages
          .map((dynamic item) => ChatMessage.fromJson(item as Map<String, dynamic>))
          .toList();
    }
    throw ApiException(_parseError(response));
  }

  static Future<bool> sendChatMessage(int receiverId, String content) async {
    final uri = Uri.parse('$baseUrl/api/chat/send');
    final response = await http.post(
      uri,
      headers: _headers(auth: true),
      body: jsonEncode({
        'receiver_id': receiverId,
        'content': content,
      }),
    );
    if (response.statusCode == 200) {
      return true;
    }
    throw ApiException(_parseError(response));
  }

  static Future<int> fetchChatUnreadCount() async {
    final uri = Uri.parse('$baseUrl/api/chat/unread-count');
    final response = await http.get(uri, headers: _headers(auth: true));
    if (response.statusCode == 200) {
      final body = jsonDecode(response.body) as Map<String, dynamic>;
      return body['unread_count'] as int? ?? 0;
    }
    throw ApiException(_parseError(response));
  }

  static Future<List<TitleModel>> fetchTitles() async {
    final uri = Uri.parse('$baseUrl/api/titles');
    final response = await http.get(uri, headers: _headers(auth: true));
    if (response.statusCode == 200) {
      final body = jsonDecode(response.body) as Map<String, dynamic>;
      final titles = body['titles'] as List<dynamic>? ?? [];
      return titles
          .map((dynamic item) => TitleModel.fromJson(item as Map<String, dynamic>))
          .toList();
    }
    throw ApiException(_parseError(response));
  }

  static Future<bool> createTitle(String name, {String? description, String? color, String? icon}) async {
    final uri = Uri.parse('$baseUrl/api/titles');
    final response = await http.post(
      uri,
      headers: _headers(auth: true),
      body: jsonEncode({
        'name': name,
        if (description != null) 'description': description,
        if (color != null) 'color': color,
        if (icon != null) 'icon': icon,
      }),
    );
    if (response.statusCode == 200) {
      return true;
    }
    throw ApiException(_parseError(response));
  }

  static Future<bool> updateTitle(int titleId, {String? name, String? description, String? color, String? icon, bool? isActive}) async {
    final uri = Uri.parse('$baseUrl/api/titles/$titleId');
    final body = <String, dynamic>{};
    if (name != null) body['name'] = name;
    if (description != null) body['description'] = description;
    if (color != null) body['color'] = color;
    if (icon != null) body['icon'] = icon;
    if (isActive != null) body['is_active'] = isActive;
    final response = await http.put(uri, headers: _headers(auth: true), body: jsonEncode(body));
    if (response.statusCode == 200) {
      return true;
    }
    throw ApiException(_parseError(response));
  }

  static Future<bool> deleteTitle(int titleId) async {
    final uri = Uri.parse('$baseUrl/api/titles/$titleId');
    final response = await http.delete(uri, headers: _headers(auth: true));
    if (response.statusCode == 200) {
      return true;
    }
    throw ApiException(_parseError(response));
  }

  static Future<bool> grantTitle(int userId, int titleId, {String? reason}) async {
    final uri = Uri.parse('$baseUrl/api/titles/grant');
    final response = await http.post(
      uri,
      headers: _headers(auth: true),
      body: jsonEncode({
        'user_id': userId,
        'title_id': titleId,
        if (reason != null) 'reason': reason,
      }),
    );
    if (response.statusCode == 200) {
      return true;
    }
    throw ApiException(_parseError(response));
  }

  static Future<bool> revokeTitle(int userId, int titleId) async {
    final uri = Uri.parse('$baseUrl/api/titles/revoke');
    final response = await http.post(
      uri,
      headers: _headers(auth: true),
      body: jsonEncode({
        'user_id': userId,
        'title_id': titleId,
      }),
    );
    if (response.statusCode == 200) {
      return true;
    }
    throw ApiException(_parseError(response));
  }

  static Future<List<TitleModel>> fetchUserTitles(int userId) async {
    final uri = Uri.parse('$baseUrl/api/users/$userId/titles');
    final response = await http.get(uri, headers: _headers(auth: true));
    if (response.statusCode == 200) {
      final body = jsonDecode(response.body) as Map<String, dynamic>;
      final titles = body['titles'] as List<dynamic>? ?? [];
      return titles
          .map((dynamic item) => TitleModel.fromJson(item as Map<String, dynamic>))
          .toList();
    }
    throw ApiException(_parseError(response));
  }

  static Future<User> fetchUserById(int userId) async {
    final uri = Uri.parse('$baseUrl/api/users/$userId');
    final response = await http.get(uri, headers: _headers());
    if (response.statusCode == 200) {
      final body = jsonDecode(response.body) as Map<String, dynamic>;
      return User.fromJson(body);
    }
    throw ApiException(_parseError(response));
  }

  static Future<List<Article>> fetchUserArticles(int userId) async {
    final uri = Uri.parse('$baseUrl/api/users/$userId/articles');
    final response = await http.get(uri, headers: _headers());
    if (response.statusCode == 200) {
      final body = jsonDecode(response.body) as Map<String, dynamic>;
      final articles = body['articles'] as List<dynamic>? ?? [];
      return articles
          .map((dynamic item) => Article.fromJson(item as Map<String, dynamic>))
          .toList();
    }
    throw ApiException(_parseError(response));
  }

  static Future<Announcement> fetchAnnouncement() async {
    final uri = Uri.parse('$baseUrl/api/announcement');
    final response = await http.get(uri, headers: _headers());
    if (response.statusCode == 200) {
      final body = jsonDecode(response.body) as Map<String, dynamic>;
      return Announcement.fromJson(body);
    }
    throw ApiException(_parseError(response));
  }

  static Future<bool> updateAnnouncement(String content) async {
    final uri = Uri.parse('$baseUrl/api/announcement');
    final response = await http.put(
      uri,
      headers: _headers(auth: true),
      body: jsonEncode({'content': content}),
    );
    if (response.statusCode == 200) {
      return true;
    }
    throw ApiException(_parseError(response));
  }

  static Future<List<SidebarItem>> fetchSidebarConfig() async {
    final uri = Uri.parse('$baseUrl/api/sidebar-config');
    final response = await http.get(uri, headers: _headers());
    if (response.statusCode == 200) {
      final body = jsonDecode(response.body) as Map<String, dynamic>;
      final items = body['items'] as List<dynamic>? ?? [];
      return items
          .map((dynamic item) => SidebarItem.fromJson(item as Map<String, dynamic>))
          .toList();
    }
    throw ApiException(_parseError(response));
  }

  static Future<bool> updateSidebarConfig(List<SidebarItem> items) async {
    final uri = Uri.parse('$baseUrl/api/sidebar-config');
    final response = await http.put(
      uri,
      headers: _headers(auth: true),
      body: jsonEncode({'items': items.map((e) => e.toJson()).toList()}),
    );
    if (response.statusCode == 200) {
      return true;
    }
    throw ApiException(_parseError(response));
  }

  static Future<SiteConfig> fetchSiteConfig() async {
    final uri = Uri.parse('$baseUrl/api/site-config');
    final response = await http.get(uri, headers: _headers());
    if (response.statusCode == 200) {
      final body = jsonDecode(response.body) as Map<String, dynamic>;
      return SiteConfig.fromJson(body);
    }
    throw ApiException(_parseError(response));
  }

  static Future<bool> updateSiteConfig({
    String? siteTitle,
    String? smtpHost,
    int? smtpPort,
    String? smtpUsername,
    String? smtpPassword,
    String? smtpFrom,
    String? smtpFromName,
  }) async {
    final uri = Uri.parse('$baseUrl/api/site-config');
    final requestBody = <String, dynamic>{};
    if (siteTitle != null) requestBody['site_title'] = siteTitle;
    if (smtpHost != null) requestBody['smtp_host'] = smtpHost;
    if (smtpPort != null) requestBody['smtp_port'] = smtpPort;
    if (smtpUsername != null) requestBody['smtp_username'] = smtpUsername;
    if (smtpPassword != null) requestBody['smtp_password'] = smtpPassword;
    if (smtpFrom != null) requestBody['smtp_from'] = smtpFrom;
    if (smtpFromName != null) requestBody['smtp_from_name'] = smtpFromName;
    final response = await http.put(
      uri,
      headers: _headers(auth: true),
      body: jsonEncode(requestBody),
    );
    if (response.statusCode == 200) {
      return true;
    }
    throw ApiException(_parseError(response));
  }

  static Future<bool> testSMTPConfig({
    required String smtpHost,
    required int smtpPort,
    required String smtpUsername,
    required String smtpPassword,
    required String smtpFrom,
    required String smtpTo,
  }) async {
    final uri = Uri.parse('$baseUrl/api/site-config/test-smtp');
    final response = await http.post(
      uri,
      headers: _headers(auth: true),
      body: jsonEncode({
        'smtp_host': smtpHost,
        'smtp_port': smtpPort,
        'smtp_username': smtpUsername,
        'smtp_password': smtpPassword,
        'smtp_from': smtpFrom,
        'smtp_to': smtpTo,
      }),
    );
    if (response.statusCode == 200) {
      return true;
    }
    throw ApiException(_parseError(response));
  }

  static Future<String> uploadAvatar(File file) async {
    final uri = Uri.parse('$baseUrl/api/upload/avatar');
    final request = http.MultipartRequest('POST', uri);
    request.headers.addAll(_headers(auth: true));
    request.files.add(await http.MultipartFile.fromPath('avatar', file.path));
    final response = await request.send();
    final body = await response.stream.bytesToString();
    if (response.statusCode == 200) {
      final data = jsonDecode(body) as Map<String, dynamic>;
      return data['url'] as String? ?? '';
    }
    throw ApiException(_parseError(http.Response(body, response.statusCode)));
  }

  static Future<String> uploadImage(File file) async {
    final uri = Uri.parse('$baseUrl/api/upload/image');
    final request = http.MultipartRequest('POST', uri);
    request.headers.addAll(_headers(auth: true));
    request.files.add(await http.MultipartFile.fromPath('image', file.path));
    final response = await request.send();
    final body = await response.stream.bytesToString();
    if (response.statusCode == 200) {
      final data = jsonDecode(body) as Map<String, dynamic>;
      return data['url'] as String? ?? '';
    }
    throw ApiException(_parseError(http.Response(body, response.statusCode)));
  }

  static Future<String> uploadVideo(File file) async {
    final uri = Uri.parse('$baseUrl/api/upload/video');
    final request = http.MultipartRequest('POST', uri);
    request.headers.addAll(_headers(auth: true));
    request.files.add(await http.MultipartFile.fromPath('video', file.path));
    final response = await request.send();
    final body = await response.stream.bytesToString();
    if (response.statusCode == 200) {
      final data = jsonDecode(body) as Map<String, dynamic>;
      return data['url'] as String? ?? '';
    }
    throw ApiException(_parseError(http.Response(body, response.statusCode)));
  }

  static Future<User?> fetchProfile() async {
    final uri = Uri.parse('$baseUrl/api/profile');
    final response = await http.get(uri, headers: _headers(auth: true));
    if (response.statusCode == 200) {
      final body = jsonDecode(response.body) as Map<String, dynamic>;
      currentUser = User.fromJson(body);
      return currentUser;
    }
    throw ApiException(_parseError(response));
  }

  static Future<bool> updateProfile({String? displayName, String? avatar, String? signature}) async {
    final uri = Uri.parse('$baseUrl/api/profile');
    final body = <String, dynamic>{};
    if (displayName != null) body['display_name'] = displayName;
    if (avatar != null) body['avatar'] = avatar;
    if (signature != null) body['signature'] = signature;
    final response = await http.put(uri, headers: _headers(auth: true), body: jsonEncode(body));
    if (response.statusCode == 200) {
      return true;
    }
    throw ApiException(_parseError(response));
  }
}
