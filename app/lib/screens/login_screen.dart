import 'package:flutter/material.dart';

import '../services/api_service.dart';

class LoginScreen extends StatefulWidget {
  const LoginScreen({super.key});

  @override
  State<LoginScreen> createState() => _LoginScreenState();
}

class _LoginScreenState extends State<LoginScreen> {
  final _formKey = GlobalKey<FormState>();
  final TextEditingController _usernameController = TextEditingController();
  final TextEditingController _qqController = TextEditingController();
  final TextEditingController _displayNameController = TextEditingController();
  final TextEditingController _passwordController = TextEditingController();
  bool _isRegister = false;
  bool _loading = false;
  String? _errorMessage;

  @override
  void dispose() {
    _usernameController.dispose();
    _qqController.dispose();
    _displayNameController.dispose();
    _passwordController.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    if (!_formKey.currentState!.validate()) {
      return;
    }

    setState(() {
      _loading = true;
      _errorMessage = null;
    });

    try {
      final username = _usernameController.text.trim();
      final password = _passwordController.text;
      final qqNumber = _qqController.text.trim();
      final displayName = _displayNameController.text.trim();

      if (_isRegister) {
        final success = await ApiService.register(username, qqNumber, displayName, password);
        if (success) {
          if (!mounted) return;
          setState(() {
            _isRegister = false;
          });
          ScaffoldMessenger.of(context).showSnackBar(const SnackBar(
            content: Text('注册成功，请登录。'),
          ));
        } else {
          if (!mounted) return;
          setState(() {
            _errorMessage = '注册失败，请检查输入或稍后重试。';
          });
        }
      } else {
        final success = await ApiService.login(username, password);
        if (!mounted) return;
        final currentContext = context;
        if (success) {
          Navigator.pushReplacementNamed(currentContext, '/home');
          return;
        }
        setState(() {
          _errorMessage = '登录失败，请检查用户名和密码。';
        });
      }
    } catch (error) {
      setState(() {
        _errorMessage = '网络错误：$error';
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
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text(_isRegister ? '注册账号' : '用户登录')),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Center(
          child: SingleChildScrollView(
            child: Form(
              key: _formKey,
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  TextFormField(
                    controller: _usernameController,
                    decoration: const InputDecoration(labelText: '用户名'),
                    validator: (value) {
                      if (value == null || value.trim().isEmpty) {
                        return '请输入用户名';
                      }
                      return null;
                    },
                  ),
                  if (_isRegister) ...[
                    const SizedBox(height: 12),
                    TextFormField(
                      controller: _qqController,
                      decoration: const InputDecoration(labelText: 'QQ号'),
                      keyboardType: TextInputType.number,
                      validator: (value) {
                        if (value == null || value.trim().isEmpty) {
                          return '请输入QQ号';
                        }
                        return null;
                      },
                    ),
                    const SizedBox(height: 12),
                    TextFormField(
                      controller: _displayNameController,
                      decoration: const InputDecoration(labelText: '显示名称'),
                      validator: (value) {
                        if (value == null || value.trim().isEmpty) {
                          return '请输入显示名称';
                        }
                        return null;
                      },
                    ),
                  ],
                  TextFormField(
                    controller: _passwordController,
                    decoration: const InputDecoration(labelText: '密码'),
                    obscureText: true,
                    validator: (value) {
                      if (value == null || value.length < 6) {
                        return '请输入至少6位密码';
                      }
                      return null;
                    },
                  ),
                  const SizedBox(height: 24),
                  if (_errorMessage != null)
                    Padding(
                      padding: const EdgeInsets.only(bottom: 12),
                      child: Text(
                        _errorMessage!,
                        style: const TextStyle(color: Colors.red),
                      ),
                    ),
                  ElevatedButton(
                    onPressed: _loading ? null : _submit,
                    child: _loading
                        ? const SizedBox(
                            height: 20,
                            width: 20,
                            child: CircularProgressIndicator(strokeWidth: 2),
                          )
                        : Text(_isRegister ? '注册' : '登录'),
                  ),
                  const SizedBox(height: 12),
                  TextButton(
                    onPressed: _loading
                        ? null
                        : () {
                            setState(() {
                              _isRegister = !_isRegister;
                              _errorMessage = null;
                            });
                          },
                    child: Text(_isRegister ? '已有账号？登录' : '没有账号？注册'),
                  ),
                  if (!_isRegister)
                    TextButton(
                      onPressed: _loading
                          ? null
                          : () {
                              Navigator.pushNamed(context, '/password-reset');
                            },
                      child: const Text('忘记密码'),
                    ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}
