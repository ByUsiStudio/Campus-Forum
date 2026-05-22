package utils

import (
	"fmt"
	"runtime"
	"time"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[90m"
)

type Logger struct {
	prefix string
	color  string
}

func NewLogger(prefix string, color string) *Logger {
	return &Logger{
		prefix: prefix,
		color:  color,
	}
}

func (l *Logger) print(level string, levelColor string, format string, args ...interface{}) {
	now := time.Now().Format("2006-01-02 15:04:05")
	_, file, line, _ := runtime.Caller(2)

	if len(args) > 0 {
		fmt.Printf("%s[%s]%s %s%s%s %s:%d %s\n",
			Gray, now, Reset,
			levelColor, level, Reset,
			file, line,
			fmt.Sprintf(format, args...))
	} else {
		fmt.Printf("%s[%s]%s %s%s%s %s:%d %s\n",
			Gray, now, Reset,
			levelColor, level, Reset,
			file, line,
			format)
	}
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.print("INFO", Cyan, format, args...)
}

func (l *Logger) Success(format string, args ...interface{}) {
	l.print("SUCCESS", Green, format, args...)
}

func (l *Logger) Warn(format string, args ...interface{}) {
	l.print("WARN", Yellow, format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.print("ERROR", Red, format, args...)
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.print("DEBUG", Gray, format, args...)
}

func Info(format string, args ...interface{}) {
	logger := NewLogger("", Cyan)
	logger.Info(format, args...)
}

func Success(format string, args ...interface{}) {
	logger := NewLogger("", Green)
	logger.Success(format, args...)
}

func Warn(format string, args ...interface{}) {
	logger := NewLogger("", Yellow)
	logger.Warn(format, args...)
}

func Error(format string, args ...interface{}) {
	logger := NewLogger("", Red)
	logger.Error(format, args...)
}

func Debug(format string, args ...interface{}) {
	logger := NewLogger("", Gray)
	logger.Debug(format, args...)
}

func Request(method, path string, statusCode int, duration time.Duration) {
	statusColor := Green
	if statusCode >= 400 {
		statusColor = Red
	} else if statusCode >= 300 {
		statusColor = Yellow
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("%s[%s]%s %s%-7s%s %s%s%s %s%d%s %s%dms%s\n",
		Gray, now, Reset,
		Cyan, method, Reset,
		Yellow, path, Reset,
		statusColor, statusCode, Reset,
		Gray, duration.Milliseconds(), Reset)
}

func Section(title string) {
	fmt.Printf("\n%s━━━ %s %s━━━%s\n", Purple, title, Purple, Reset)
}

func Separator() {
	fmt.Printf("%s%s%s\n", Gray, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━", Reset)
}
