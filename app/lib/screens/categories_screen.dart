import 'package:flutter/material.dart';

import '../models/category.dart';
import '../services/api_service.dart';
import 'article_list_screen.dart';

class CategoriesScreen extends StatefulWidget {
  const CategoriesScreen({super.key});

  @override
  State<CategoriesScreen> createState() => _CategoriesScreenState();
}

class _CategoriesScreenState extends State<CategoriesScreen> {
  late Future<List<Category>> _categoriesFuture;

  @override
  void initState() {
    super.initState();
    _categoriesFuture = ApiService.fetchCategories();
  }

  Future<void> _refresh() async {
    setState(() {
      _categoriesFuture = ApiService.fetchCategories();
    });
  }

  @override
  Widget build(BuildContext context) {
    return FutureBuilder<List<Category>>(
      future: _categoriesFuture,
      builder: (context, snapshot) {
        if (snapshot.connectionState == ConnectionState.waiting) {
          return const Center(child: CircularProgressIndicator());
        }
        if (snapshot.hasError) {
          return Center(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                const Text('加载分类失败，请稍后重试。'),
                const SizedBox(height: 12),
                Text(snapshot.error.toString()),
                const SizedBox(height: 16),
                ElevatedButton(onPressed: _refresh, child: const Text('刷新')),
              ],
            ),
          );
        }

        final categories = snapshot.data ?? [];
        if (categories.isEmpty) {
          return const Center(child: Text('当前没有分类。'));
        }

        return RefreshIndicator(
          onRefresh: _refresh,
          child: ListView.separated(
            itemCount: categories.length,
            separatorBuilder: (_, __) => const Divider(height: 1),
            itemBuilder: (context, index) {
              final category = categories[index];
              return ListTile(
                title: Text(category.name),
                subtitle: Text(category.description),
                onTap: () {
                  Navigator.push(
                    context,
                    MaterialPageRoute<void>(
                      builder: (context) => ArticleListScreen(categoryId: category.id),
                    ),
                  );
                },
              );
            },
          ),
        );
      },
    );
  }
}
