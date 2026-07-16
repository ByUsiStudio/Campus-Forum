package utils

import (
	"strconv"
	"time"

	"github.com/russross/blackfriday/v2"
	"gorm.io/gorm"
)

// Atoi 字符串转整数
func Atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

// Atouint 字符串转uint
func Atouint(s string) uint {
	i, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0
	}
	return uint(i)
}

// GetDate 获取当前日期 YYYY-MM-DD格式
func GetDate() string {
	return time.Now().Format("2006-01-02")
}

// GetWeek 获取当前周数
func GetWeek() int {
	_, week := time.Now().ISOWeek()
	return week
}

// GetMonth 获取当前月份 YYYY-MM格式
func GetMonth() string {
	return time.Now().Format("2006-01")
}

// GetYear 获取当前年份
func GetYear() int {
	return time.Now().Year()
}

// ParseDate 解析日期字符串
func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

// FormatDate 格式化日期
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatDateTime 格式化日期时间
func FormatDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// GetStartOfDay 获取一天的开始时间
func GetStartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// GetEndOfDay 获取一天的结束时间
func GetEndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
}

// DaysBetween 计算两个日期之间的天数
func DaysBetween(start, end time.Time) int {
	return int(end.Sub(start).Hours() / 24)
}

// IsToday 检查日期是否是今天
func IsToday(t time.Time) bool {
	now := time.Now()
	return t.Year() == now.Year() && t.Month() == now.Month() && t.Day() == now.Day()
}

// IsYesterday 检查日期是否是昨天
func IsYesterday(t time.Time) bool {
	yesterday := time.Now().AddDate(0, 0, -1)
	return t.Year() == yesterday.Year() && t.Month() == yesterday.Month() && t.Day() == yesterday.Day()
}

// IsThisWeek 检查日期是否在本周
func IsThisWeek(t time.Time) bool {
	now := time.Now()
	_, currentWeek := now.ISOWeek()
	_, targetWeek := t.ISOWeek()
	return currentWeek == targetWeek && t.Year() == now.Year()
}

// IsThisMonth 检查日期是否在本月
func IsThisMonth(t time.Time) bool {
	now := time.Now()
	return t.Year() == now.Year() && t.Month() == now.Month()
}

// IsThisYear 检查日期是否在今年
func IsThisYear(t time.Time) bool {
	return t.Year() == time.Now().Year()
}

// Increment GORM自增表达式
func Increment(field string) interface{} {
	return gorm.Expr(field + " + 1")
}

// Decrement GORM自减表达式
func Decrement(field string) interface{} {
	return gorm.Expr(field + " - 1")
}

// RenderMarkdown 将Markdown转换为HTML
func RenderMarkdown(content string) string {
	html := blackfriday.Run([]byte(content))
	return string(html)
}
