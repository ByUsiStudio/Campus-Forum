import 'package:flutter/material.dart';

import 'screens/article_detail_screen.dart';
import 'screens/home_screen.dart';
import 'screens/login_screen.dart';
import 'screens/password_reset_screen.dart';
import 'services/api_service.dart';

Future<void> main() async {
  WidgetsFlutterBinding.ensureInitialized();
  await ApiService.init();
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: '麻阳综合性校园论坛',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: const Color.fromARGB(255, 210, 32, 255)),
        useMaterial3: true,
      ),
      initialRoute: ApiService.isLoggedIn ? '/home' : '/',
      routes: {
        '/': (context) => const LoginScreen(),
        '/home': (context) => const HomeScreen(),
        '/article': (context) => const ArticleDetailScreen(),
        '/password-reset': (context) => const PasswordResetScreen(),
      },
    );
  }
}
