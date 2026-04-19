package model

import (
	"time"
)

// Post 文章模型
type Post struct {
	ID           uint   `gorm:"primarykey" json:"id"`
	Title        string `gorm:"not null;size:200" json:"title"`
	Slug         string `gorm:"unique;size:200" json:"slug"`
	Content      string `gorm:"type:text" json:"content"` // Markdown内容
	CoverImage   string `gorm:"size:255" json:"cover_image"`
	Excerpt      string `gorm:"type:text" json:"excerpt"` // 摘要
	ColumnID     uint   `gorm:"not null" json:"column_id"`
	Status       string `gorm:"size:20;default:draft" json:"status"` // draft/published
	ViewCount    int    `gorm:"default:0" json:"view_count"`
	LikeCount    int    `gorm:"default:0" json:"like_count"`
	CommentCount int    `gorm:"default:0" json:"comment_count"`
	IsTop        bool   `gorm:"default:false" json:"is_top"`
	CreatedAt    int64  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    int64  `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 表名
func (Post) TableName() string {
	return "posts"
}

// IsPublished 是否已发布
func (p *Post) IsPublished() bool {
	return p.Status == "published"
}

// FormatDate 格式化日期
func (p *Post) FormatDate() string {
	if p.CreatedAt == 0 {
		return ""
	}
	return time.Unix(p.CreatedAt, 0).Format("2006-01-02")
}
