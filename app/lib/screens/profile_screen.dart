import 'package:flutter/material.dart';

import '../services/api_service.dart';

class ProfileScreen extends StatefulWidget {
  const ProfileScreen({super.key});

  @override
  State<ProfileScreen> createState() => _ProfileScreenState();
}

class _ProfileScreenState extends State<ProfileScreen> {
  bool _loading = false;
  String? _errorMessage;

  Future<void> _refreshProfile() async {
    setState(() {
      _loading = true;
      _errorMessage = null;
    });
    try {
      await ApiService.fetchProfile();
    } catch (error) {
      setState(() {
        _errorMessage = '加载用户资料失败：$error';
      });
    } finally {
      if (mounted) {
        setState(() {
          _loading = false;
        });
      }
    }
  }

  @override
  void initState() {
    super.initState();
    _refreshProfile();
  }

  @override
  Widget build(BuildContext context) {
    final user = ApiService.currentUser;
    return Padding(
      padding: const EdgeInsets.all(16.0),
      child: _loading
          ? const Center(child: CircularProgressIndicator())
          : Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                if (_errorMessage != null)
                  Padding(
                    padding: const EdgeInsets.symmetric(vertical: 12),
                    child: Text(_errorMessage!, style: const TextStyle(color: Colors.red)),
                  ),
                if (user != null) ...[
                  CircleAvatar(
                    radius: 42,
                    backgroundImage: user.avatar.isNotEmpty ? NetworkImage(user.avatar) : null,
                    child: user.avatar.isEmpty ? Text(user.displayName.isNotEmpty ? user.displayName[0] : '我') : null,
                  ),
                  const SizedBox(height: 16),
                  Text('昵称: ${user.displayName}', style: const TextStyle(fontSize: 18)),
                  const SizedBox(height: 8),
                  Text('用户名: ${user.username}'),
                  const SizedBox(height: 8),
                  Text('QQ号: ${user.qqNumber}'),
                  const SizedBox(height: 8),
                  Text('角色: ${user.role}'),
                  const SizedBox(height: 8),
                  Text('注册时间: ${user.createdAt.isNotEmpty ? user.createdAt : '未知'}'),
                  const SizedBox(height: 24),
                  ElevatedButton(
                    onPressed: () async {
                      await ApiService.logout();
                      if (!mounted) return;
                      Navigator.pushReplacementNamed(context, '/');
                    },
                    child: const Text('退出登录'),
                  ),
                ] else ...[
                  const Center(child: Text('当前未登录，请重新登录。')),
                  const SizedBox(height: 16),
                  ElevatedButton(
                    onPressed: () {
                      Navigator.pushReplacementNamed(context, '/');
                    },
                    child: const Text('返回登录'),
                  ),
                ]
              ],
            ),
    );
  }
}
