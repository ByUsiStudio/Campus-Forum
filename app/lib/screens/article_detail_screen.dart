import 'package:flutter/material.dart';

import '../models/article.dart';
import '../services/api_service.dart';

class ArticleDetailScreen extends StatefulWidget {
  const ArticleDetailScreen({super.key});

  @override
  State<ArticleDetailScreen> createState() => _ArticleDetailScreenState();
}

class _ArticleDetailScreenState extends State<ArticleDetailScreen> {
  late final int articleId;
  late Future<Article> _detailFuture;

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    articleId = ModalRoute.of(context)?.settings.arguments as int;
    _detailFuture = ApiService.fetchArticleDetail(articleId);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('文章详情')),
      body: FutureBuilder<Article>(
        future: _detailFuture,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return const Center(child: CircularProgressIndicator());
          }
          if (snapshot.hasError) {
            return Center(child: Text('加载失败：${snapshot.error}'));
          }
          final article = snapshot.data!;
          return Padding(
            padding: const EdgeInsets.all(16.0),
            child: SingleChildScrollView(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(article.title, style: const TextStyle(fontSize: 22, fontWeight: FontWeight.bold)),
                  const SizedBox(height: 12),
                  Row(
                    children: [
                      Text('分类: ${article.category?['name'] ?? '未分类'}'),
                      const SizedBox(width: 16),
                      Text('点赞 ${article.likeCount}'),
                    ],
                  ),
                  const SizedBox(height: 16),
                  Text(article.content, style: const TextStyle(fontSize: 16, height: 1.5)),
                  const SizedBox(height: 24),
                  Text('作者: ${article.user?['display_name'] ?? article.user?['username'] ?? '匿名'}'),
                  const SizedBox(height: 4),
                  Text('创建时间: ${article.createdAt}'),
                ],
              ),
            ),
          );
        },
      ),
    );
  }
}
