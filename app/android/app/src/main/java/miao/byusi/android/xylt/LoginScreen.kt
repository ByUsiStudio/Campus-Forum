package miao.byusi.android.xylt

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.unit.dp
import androidx.compose.ui.graphics.Color
import androidx.navigation.NavHostController
import miuix.compose.scaffold.Scaffold
import miuix.compose.surface.Surface
import miuix.compose.text.Text
import miuix.compose.button.Button
import miuix.compose.button.FilledButton
import miuix.compose.button.TextButton
import miuix.compose.progress.CircularProgressIndicator
import miuix.compose.icon.Icon
import miuix.compose.icon.icons.MiuixIcons
import miuix.compose.icon.icons.outlined.ArrowBack
import miuix.compose.icon.icons.outlined.Person
import miuix.compose.icon.icons.outlined.Lock
import miuix.compose.textfield.TextField

@Composable
fun LoginScreen(navController: NavHostController) {
    var username by remember { mutableStateOf("") }
    var password by remember { mutableStateOf("") }
    var isLoading by remember { mutableStateOf(false) }
    var error by remember { mutableStateOf<String?>(null) }
    
    Scaffold(
        topBar = {
            Surface(
                modifier = Modifier.fillMaxWidth()
            ) {
                Row(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(16.dp),
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Button(onClick = { navController.popBackStack() }) {
                        Icon(
                            imageVector = MiuixIcons.Outlined.ArrowBack,
                            contentDescription = "返回"
                        )
                    }
                    Text(
                        text = "登录",
                        fontSize = 20.dp
                    )
                }
            }
        }
    ) { padding ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(padding)
                .padding(32.dp),
            horizontalAlignment = Alignment.CenterHorizontally
        ) {
            Spacer(modifier = Modifier.height(32.dp))
            
            Text(
                text = "校园论坛",
                fontSize = 28.dp,
                color = Color(0xFF007AFF)
            )
            
            Spacer(modifier = Modifier.height(8.dp))
            
            Text(
                text = "登录您的账户",
                fontSize = 14.dp,
                color = Color(0xFF666666)
            )
            
            Spacer(modifier = Modifier.height(32.dp))
            
            // 用户名输入
            TextField(
                value = username,
                onValueChange = { username = it },
                label = "用户名",
                modifier = Modifier.fillMaxWidth(),
                leadingIcon = {
                    Icon(
                        imageVector = MiuixIcons.Outlined.Person,
                        contentDescription = null
                    )
                }
            )
            
            Spacer(modifier = Modifier.height(16.dp))
            
            // 密码输入
            TextField(
                value = password,
                onValueChange = { password = it },
                label = "密码",
                visualTransformation = PasswordVisualTransformation(),
                keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Password),
                modifier = Modifier.fillMaxWidth(),
                leadingIcon = {
                    Icon(
                        imageVector = MiuixIcons.Outlined.Lock,
                        contentDescription = null
                    )
                }
            )
            
            if (error != null) {
                Spacer(modifier = Modifier.height(8.dp))
                Text(
                    text = error ?: "",
                    color = Color(0xFFFF3B30),
                    fontSize = 12.dp
                )
            }
            
            Spacer(modifier = Modifier.height(24.dp))
            
            FilledButton(
                onClick = {
                    if (username.isEmpty()) {
                        error = "请输入用户名"
                        return
                    }
                    if (password.isEmpty()) {
                        error = "请输入密码"
                        return
                    }
                    isLoading = true
                    error = null
                    login(username, password) { result ->
                        isLoading = false
                        result.onSuccess {
                            navController.popBackStack()
                        }.onFailure { e ->
                            error = e.message
                        }
                    }
                },
                modifier = Modifier.fillMaxWidth(),
                enabled = !isLoading
            ) {
                if (isLoading) {
                    CircularProgressIndicator(
                        modifier = Modifier.size(24.dp)
                    )
                } else {
                    Text("登录")
                }
            }
            
            Spacer(modifier = Modifier.height(16.dp))
            
            Row(
                horizontalArrangement = Arrangement.Center
            ) {
                Text(
                    text = "没有账户？",
                    fontSize = 14.dp,
                    color = Color(0xFF666666)
                )
                TextButton(onClick = { navController.navigate("register") }) {
                    Text("立即注册")
                }
            }
        }
    }
}

private fun login(username: String, password: String, callback: (Result<Unit>) -> Unit) {
    ApiClient.login(username, password, object : ApiCallback {
        override fun onSuccess(response: String) {
            try {
                val json = org.json.JSONObject(response)
                val token = json.optString("token", "")
                if (!token.isEmpty()) {
                    ApiClient.setAuthToken(token)
                }
                callback(Result.success(Unit))
            } catch (e: Exception) {
                callback(Result.failure(e))
            }
        }
        override fun onError(error: String) {
            callback(Result.failure(Exception(error)))
        }
    })
}