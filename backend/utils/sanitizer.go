package utils

import (
	"html"
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

// 创建全局XSS过滤器实例
var (
	// SanitizePolicy 论坛专用的HTML清理策略
	SanitizePolicy = func() *bluemonday.Policy {
		p := bluemonday.UGCPolicy()

		// 图片
		p.AllowElements("img")
		p.AllowAttrs(
			"src",
			"alt",
			"title",
			"width",
			"height",
		).OnElements("img")

		// 视频
		p.AllowElements("video", "source")
		p.AllowAttrs(
			"src",
			"type",
			"controls",
			"width",
			"height",
			"poster",
		).OnElements("video", "source")

		// 音频
		p.AllowElements("audio")
		p.AllowAttrs(
			"src",
			"type",
			"controls",
		).OnElements("audio")

		// iframe（B站、YouTube）
		p.AllowElements("iframe")
		p.AllowAttrs(
			"src",
			"width",
			"height",
			"frameborder",
			"allowfullscreen",
			"allow",
		).OnElements("iframe")

		// 代码高亮
		p.AllowAttrs("class").Globally()

		// data-*
		p.AllowDataAttributes()

		// 允许协议
		p.AllowURLSchemes(
			"http",
			"https",
		)

		return p
	}()

	// StrictPolicy 严格策略（仅文本）
	StrictPolicy = bluemonday.StrictPolicy()
)

// SanitizeHTML 清理HTML内容，防止XSS攻击（使用论坛专用策略）
func SanitizeHTML(input string) string {
	return SanitizePolicy.Sanitize(input)
}

// SanitizeHTMLStrict 严格模式清理HTML（移除所有标签）
func SanitizeHTMLStrict(input string) string {
	return StrictPolicy.Sanitize(input)
}

// SanitizeText 清理纯文本输入
func SanitizeText(input string) string {
	// 去除首尾空格
	input = strings.TrimSpace(input)

	// 限制长度
	if len(input) > 10000 {
		input = input[:10000]
	}

	// 转义HTML
	return html.EscapeString(input)
}

// ValidateEmail 简单的邮箱格式验证
func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// ValidateQQNumber 验证QQ号码格式
func ValidateQQNumber(qq string) bool {
	pattern := `^[1-9]\d{4,11}$`
	matched, _ := regexp.MatchString(pattern, qq)
	return matched
}

// ValidateUsername 验证用户名格式
func ValidateUsername(username string) bool {
	// 只允许字母、数字、下划线，长度3-20
	pattern := `^[a-zA-Z0-9_]{3,20}$`
	matched, _ := regexp.MatchString(pattern, username)
	return matched
}

// ValidatePassword 验证密码强度
func ValidatePassword(password string) (bool, string) {
	if len(password) < 6 {
		return false, "密码长度至少为6位"
	}

	if len(password) > 128 {
		return false, "密码长度不能超过128位"
	}

	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	if !hasLetter && !hasNumber {
		return false, "密码必须包含字母或数字"
	}

	return true, ""
}

// SanitizeFilename 清理文件名，防止路径遍历攻击
func SanitizeFilename(filename string) string {
	// 移除路径分隔符
	filename = strings.ReplaceAll(filename, "..", "")
	filename = strings.ReplaceAll(filename, "/", "")
	filename = strings.ReplaceAll(filename, "\\", "")

	// 只保留安全的字符
	reg := regexp.MustCompile(`[^a-zA-Z0-9._-]`)
	filename = reg.ReplaceAllString(filename, "_")

	return filename
}

// TruncateString 截断字符串到指定长度
func TruncateString(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen])
}
