package utils

import (
	"bytes"
	"regexp"

	"github.com/yuin/goldmark"
)

// 简单的XSS过滤
var unsafePatterns = []*regexp.Regexp{
	regexp.MustCompile(`<script[^>]*>.*?</script>`),
	regexp.MustCompile(`<iframe[^>]*>.*?</iframe>`),
	regexp.MustCompile(`<object[^>]*>.*?</object>`),
	regexp.MustCompile(`<embed[^>]*>.*?</embed>`),
	regexp.MustCompile(`on\w+\s*=\s*["']?[^"']*["']?`),
}

func MarkdownToHTML(md string) string {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(md), &buf); err != nil {
		return md
	}

	html := buf.String()

	// 简单过滤危险标签和事件
	for _, pattern := range unsafePatterns {
		html = pattern.ReplaceAllString(html, "")
	}

	return html
}
