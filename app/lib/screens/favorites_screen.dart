import 'package:flutter/material.dart';

import '../models/article.dart';
import '../services/api_service.dart';

class FavoritesScreen extends StatefulWidget {
  const FavoritesScreen({super.key});

  @override
  State<FavoritesScreen> createState() => _FavoritesScreenState();
}

class _FavoritesScreenState extends State<FavoritesScreen> {
  late Future<List<Article>> _favoritesFuture;

  @override
  void initState() {
    super.initState();
    _favoritesFuture = ApiService.fetchFavorites();
  }

  Future<void> _refresh() async {
    setState(() {
      _favoritesFuture = ApiService.fetchFavorites();
    });
  }

  @override
  Widget build(BuildContext context) {
    return FutureBuilder<List<Article>>(
      future: _favoritesFuture,
      builder: (context, snapshot) {
        if (snapshot.connectionState == ConnectionState.waiting) {
          return const Center(child: CircularProgressIndicator());
        }
        if (snapshot.hasError) {
          return Center(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                const Text('加载收藏失败，请稍后重试。'),
                const SizedBox(height: 12),
                Text(snapshot.error.toString()),
                const SizedBox(height: 16),
                ElevatedButton(onPressed: _refresh, child: const Text('刷新')),
              ],
            ),
          );
        }

        final favorites = snapshot.data ?? [];
        if (favorites.isEmpty) {
          return const Center(child: Text('暂无收藏内容。'));
        }

        return RefreshIndicator(
          onRefresh: _refresh,
          child: ListView.separated(
            itemCount: favorites.length,
            separatorBuilder: (_, __) => const Divider(height: 1),
            itemBuilder: (context, index) {
              final article = favorites[index];
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
                  Navigator.pushNamed(context, '/article', arguments: article.id);
                },
              );
            },
          ),
        );
      },
    );
  }
}
