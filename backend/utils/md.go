package utils

import (
	"bytes"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
)

var forumPolicy = func() *bluemonday.Policy {

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

func MarkdownToHTML(md string) string {
	var buf bytes.Buffer

	if err := goldmark.Convert([]byte(md), &buf); err != nil {
		return md
	}

	return forumPolicy.Sanitize(buf.String())
}