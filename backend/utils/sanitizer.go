package utils

import (
	"html"
	"regexp"
	"strings"
)

// SanitizeHTML 清理HTML内容，防止XSS攻击
func SanitizeHTML(input string) string {
	// 首先转义所有HTML实体
	sanitized := html.EscapeString(input)

	// 允许的安全标签白名单
	allowedTags := map[string]bool{
		"b":      true,
		"i":      true,
		"u":      true,
		"strong": true,
		"em":     true,
		"p":      true,
		"br":     true,
		"a":      true,
		"img":    true,
		"code":   true,
		"pre":    true,
		"h1":     true,
		"h2":     true,
		"h3":     true,
		"h4":     true,
		"h5":     true,
		"h6":     true,
		"ul":     true,
		"ol":     true,
		"li":     true,
		"blockquote": true,
	}

	// 恢复允许的标签
	for tag := range allowedTags {
		openPattern := regexp.MustCompile(`&lt;` + tag + `(&gt;|\s[^&]*&gt;)`)
		closePattern := regexp.MustCompile(`&lt;/` + tag + `&gt;`)

		sanitized = openPattern.ReplaceAllStringFunc(sanitized, func(match string) string {
			return strings.Replace(match, "&lt;", "<", 1)
		})
		sanitized = closePattern.ReplaceAllStringFunc(sanitized, func(match string) string {
			return strings.Replace(match, "&lt;", "<", 1)
		})
		sanitized = strings.ReplaceAll(sanitized, "&gt;", ">")
	}

	// 清理a标签的href属性，只允许安全协议
	hrefPattern := regexp.MustCompile(`<a\s+[^>]*href=["']([^"']*)["'][^>]*>`)
	sanitized = hrefPattern.ReplaceAllStringFunc(sanitized, func(match string) string {
		if strings.Contains(match, "javascript:") || strings.Contains(match, "data:") {
			return ""
		}
		return match
	})

	// 清理img标签的src属性
	imgPattern := regexp.MustCompile(`<img\s+[^>]*src=["']([^"']*)["'][^>]*>`)
	sanitized = imgPattern.ReplaceAllStringFunc(sanitized, func(match string) string {
		if strings.Contains(match, "javascript:") || strings.Contains(match, "data:") {
			return ""
		}
		return match
	})

	return sanitized
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
