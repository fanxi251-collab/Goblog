package model

import (
	"github.com/microcosm-cc/bluemonday"
)

var (
	// 评论XSS策略：严格，只允许基本标签
	commentPolicy = bluemonday.NewPolicy()
	// 文章XSS策略：允许更多Markdown标签
	articlePolicy = bluemonday.NewPolicy()
	// XSS开关
	xssEnabled = true
)

func init() {
	// 评论策略：只允许纯文本和换行
	commentPolicy.AllowElements("p", "br", "strong", "em")
	commentPolicy.AllowAttrs("").OnElements("p", "br", "strong", "em")

	// 文章策略：允许基本Markdown渲染后的HTML
	articlePolicy.AllowElements("p", "br", "strong", "em", "a", "code", "pre", "blockquote")
	articlePolicy.AllowAttrs("href").OnElements("a")
	articlePolicy.AllowAttrs("").OnElements("p", "br", "strong", "em", "code", "pre", "blockquote")
	// 文章策略：允许链接
	articlePolicy.AllowAttrs("href").OnElements("a")
}

// SetXSSEnabled 设置XSS开关
func SetXSSEnabled(enabled bool) {
	xssEnabled = enabled
}

// GetXSSEnabled 获取XSS开关状态
func GetXSSEnabled() bool {
	return xssEnabled
}

// SanitizeComment 清洗评论内容
func SanitizeComment(content string) string {
	if !xssEnabled {
		return content
	}
	return commentPolicy.Sanitize(content)
}

// SanitizeArticle 清洗文章内容
func SanitizeArticle(content string) string {
	if !xssEnabled {
		return content
	}
	return articlePolicy.Sanitize(content)
}
