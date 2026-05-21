package utils

import (
    "bytes"

    "github.com/microcosm-cc/bluemonday"
    "github.com/yuin/goldmark"
)

func MarkdownToHTML(md string) string {
    // 将Markdown转换为HTML
    var buf bytes.Buffer
    if err := goldmark.Convert([]byte(md), &buf); err != nil {
        return md
    }
    
    unsafeHTML := buf.Bytes()
    
    // 清理HTML，防止XSS攻击
    policy := bluemonday.UGCPolicy()
    // 允许图片、视频等标签
    policy.AllowAttrs("src").OnElements("img", "video", "source")
    policy.AllowAttrs("controls", "width", "height").OnElements("video")
    policy.AllowAttrs("type", "src").OnElements("source")
    policy.AllowAttrs("href", "target").OnElements("a")
    
    safeHTML := policy.SanitizeBytes(unsafeHTML)
    
    return string(safeHTML)
}