import 'package:flutter/material.dart';

import '../models/article.dart';
import '../services/api_service.dart';

class ArticleListScreen extends StatefulWidget {
  const ArticleListScreen({super.key, this.categoryId});

  final int? categoryId;

  @override
  State<ArticleListScreen> createState() => _ArticleListScreenState();
}

class _ArticleListScreenState extends State<ArticleListScreen> {
  late Future<List<Article>> _articlesFuture;

  @override
  void initState() {
    super.initState();
    _articlesFuture = ApiService.fetchArticles(categoryId: widget.categoryId);
  }

  Future<void> _refresh() async {
    setState(() {
      _articlesFuture = ApiService.fetchArticles(categoryId: widget.categoryId);
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('文章列表'),
        actions: [
          IconButton(
            icon: const Icon(Icons.logout),
            tooltip: '退出登录',
            onPressed: () async {
              await ApiService.logout();
              if (!mounted) return;
              final currentContext = context;
              Navigator.pushReplacementNamed(currentContext, '/');
            },
          ),
        ],
      ),
      body: FutureBuilder<List<Article>>(
        future: _articlesFuture,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return const Center(child: CircularProgressIndicator());
          }
          if (snapshot.hasError) {
            return Center(
              child: Padding(
                padding: const EdgeInsets.all(16.0),
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    const Text('加载文章失败，请稍后重试。'),
                    const SizedBox(height: 12),
                    Text(snapshot.error.toString()),
                    const SizedBox(height: 16),
                    ElevatedButton(onPressed: _refresh, child: const Text('刷新')),
                  ],
                ),
              ),
            );
          }
          final articles = snapshot.data ?? [];
          if (articles.isEmpty) {
            return const Center(child: Text('当前没有文章。'));
          }
          return RefreshIndicator(
            onRefresh: _refresh,
            child: ListView.separated(
              itemCount: articles.length,
              separatorBuilder: (_, index) => const Divider(height: 1),
              itemBuilder: (context, index) {
                final article = articles[index];
                return ListTile(
                  title: Text(article.title),
                  subtitle: Text(article.category?['name']?.toString() ?? ''),
                  trailing: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Text('点赞 ${article.likeCount}'),
                      Text('评论 ${article.commentCount}'),
                    ],
                  ),
                  onTap: () {
                    Navigator.pushNamed(
                      context,
                      '/article',
                      arguments: article.id,
                    );
                  },
                );
              },
            ),
          );
        },
      ),
    );
  }
}
