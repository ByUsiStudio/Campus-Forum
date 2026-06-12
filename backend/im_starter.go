package main

import (
	"fmt"
	"forum/utils"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

var (
	imProcess *exec.Cmd
	imMutex   sync.Mutex
)

// StartIMServer 启动IM服务器子进程
func StartIMServer(config Config) error {
	imMutex.Lock()
	defer imMutex.Unlock()

	// 检查是否已运行
	if imProcess != nil && imProcess.Process != nil {
		return nil
	}

	// IM服务可执行文件路径
	imLauncherPath := filepath.Join("sdk", "im-server", "launcher")
	imPath := filepath.Join(imLauncherPath, "im-server.exe")

	// 检查是否已编译
	compiled := true
	if _, err := os.Stat(imPath); os.IsNotExist(err) {
		compiled = false
		utils.Warn("IM服务可执行文件不存在，将使用 go run 方式启动")
	}

	// 生成IM配置文件
	imConfigPath := filepath.Join(imLauncherPath, "conf", "config.yml")
	if err := generateIMConfig(config, imConfigPath); err != nil {
		return fmt.Errorf("生成IM配置文件失败: %v", err)
	}

	utils.Info("IM配置文件已生成: %s", imConfigPath)

	var cmd *exec.Cmd
	if !compiled {
		// 使用 go run 方式运行
		mainGoPath := filepath.Join(imLauncherPath, "main.go")
		cmd = exec.Command("go", "run", mainGoPath, "-config="+imConfigPath)
		cmd.Dir = imLauncherPath
	} else {
		// 使用编译后的可执行文件
		cmd = exec.Command(imPath, "-config="+imConfigPath)
		cmd.Dir = imLauncherPath
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 启动进程
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动IM服务失败: %v", err)
	}

	imProcess = cmd
	utils.Info("IM服务已启动，PID: %d", cmd.Process.Pid)

	// 等待进程结束
	go func() {
		cmd.Wait()
		imMutex.Lock()
		imProcess = nil
		imMutex.Unlock()
		utils.Info("IM服务进程已结束")
	}()

	return nil
}

// generateIMConfig 生成IM服务配置
func generateIMConfig(config Config, configPath string) error {
	// 确保目录存在
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// 构建YAML配置内容
	var sb strings.Builder
	sb.WriteString("nodeName: \"\"\n")
	sb.WriteString("nodeHost: \"127.0.0.1\"\n")
	sb.WriteString("msgStoreEngine: \"mysql\"\n")
	sb.WriteString("\n")
	sb.WriteString("log:\n")
	sb.WriteString("  logPath: \"logs/\"\n")
	sb.WriteString("  logName: \"imserver.log\"\n")
	sb.WriteString("  logExpireHours: 168\n")
	sb.WriteString("\n")
	sb.WriteString("kvdb:\n")
	sb.WriteString("  isOpen: false\n")
	sb.WriteString("  dataPath: \"kvdb/\"\n")
	sb.WriteString("\n")
	sb.WriteString("msgLogs:\n")
	sb.WriteString("  logPath: \"logs/\"\n")
	sb.WriteString("  maxBackups: 7\n")
	sb.WriteString("  isCompress: true\n")
	sb.WriteString("\n")
	sb.WriteString("mysql:\n")
	sb.WriteString(fmt.Sprintf("  user: %q\n", config.Database.Username))
	sb.WriteString(fmt.Sprintf("  password: %q\n", config.Database.Password))
	sb.WriteString(fmt.Sprintf("  address: %q\n", config.Database.Host+":"+config.Database.Port))
	sb.WriteString(fmt.Sprintf("  name: %q\n", config.Database.DBName))
	sb.WriteString("  debug: false\n")
	sb.WriteString("\n")
	sb.WriteString("mongodb:\n")
	sb.WriteString("  address: \"localhost:27017\"\n")
	sb.WriteString("  name: \"imserver\"\n")
	sb.WriteString("\n")
	sb.WriteString("connectManager:\n")

	wsPort := config.IM.WsPort
	if wsPort == 0 {
		wsPort = 9003
	}
	sb.WriteString(fmt.Sprintf("  wsPort: %d\n", wsPort))
	sb.WriteString("\n")
	sb.WriteString("apiGateway:\n")

	apiPort := config.IM.ApiGatewayPort
	if apiPort == 0 {
		apiPort = 9001
	}
	sb.WriteString(fmt.Sprintf("  httpPort: %d\n", apiPort))
	sb.WriteString("\n")
	sb.WriteString("navGateway:\n")
	sb.WriteString("  httpPort: 9000\n")
	sb.WriteString("\n")
	sb.WriteString("adminGateway:\n")

	adminPort := config.IM.AdminPort
	if adminPort == 0 {
		adminPort = 8090
	}
	sb.WriteString(fmt.Sprintf("  httpPort: %d\n", adminPort))
	sb.WriteString("\n")
	sb.WriteString("adminSecret: \"admin_secret_change_this\"\n")

	return os.WriteFile(configPath, []byte(sb.String()), 0644)
}

// StopIMServer 停止IM服务器子进程
func StopIMServer() {
	imMutex.Lock()
	defer imMutex.Unlock()

	if imProcess != nil && imProcess.Process != nil {
		utils.Info("正在停止IM服务...")
		imProcess.Process.Kill()
		imProcess.Wait()
		imProcess = nil
		utils.Info("IM服务已停止")
	}
}

// IsIMServerRunning 检查IM服务是否运行中
func IsIMServerRunning() bool {
	imMutex.Lock()
	defer imMutex.Unlock()

	if imProcess == nil || imProcess.Process == nil {
		return false
	}

	// 检查进程是否存活
	process, err := os.FindProcess(imProcess.Process.Pid)
	if err != nil {
		return false
	}

	// 发送空信号检测进程是否存活
	err = process.Signal(os.Signal(nil))
	return err == nil
}
