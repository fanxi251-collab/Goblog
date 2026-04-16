package service

import (
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

var markdownParser = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
)

// RenderMarkdown 将Markdown渲染为HTML
func RenderMarkdown(content string) string {
	if content == "" {
		return ""
	}
	var buf strings.Builder
	if err := markdownParser.Convert([]byte(content), &buf); err != nil {
		return content
	}
	return buf.String()
}
