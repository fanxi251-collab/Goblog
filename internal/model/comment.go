package model

import (
	"time"
)

// Comment 评论模型
type Comment struct {
	ID        uint   `gorm:"primarykey" json:"id"`
	ParentID  uint   `gorm:"default:0" json:"parent_id"` // 父评论ID，0表示顶层留言
	Nickname  string `gorm:"not null;size:50" json:"nickname"`
	Email     string `gorm:"size:100" json:"email"`
	Content   string `gorm:"not null;type:text" json:"content"`
	PostID    uint   `gorm:"not null" json:"post_id"`               // 关联文章ID，0表示留言板
	Status    string `gorm:"size:20;default:pending" json:"status"` // pending/approved/rejected
	IP        string `gorm:"size:50" json:"ip"`
	UserAgent string `gorm:"size:255" json:"user_agent"`
	CreatedAt int64  `gorm:"autoCreateTime" json:"created_at"`

	// 运行时字段（不存数据库）
	Replies        []Comment `gorm:"-" json:"replies,omitempty"`         // 回复列表
	ParentNickname string    `gorm:"-" json:"parent_nickname,omitempty"` // 父评论昵称（用于显示）
	FormatTime     string    `gorm:"-" json:"format_time,omitempty"`     // 格式化后的时间
}

// FormatCreatedAt 格式化创建时间
func (c *Comment) FormatCreatedAt() string {
	if c.CreatedAt == 0 {
		return ""
	}
	return time.Unix(c.CreatedAt, 0).Format("2006-01-02 15:04:05")
}

// IsReply 是否是回复
func (c *Comment) IsReply() bool {
	return c.ParentID > 0
}

// IsTopLevel 是否是顶层留言
func (c *Comment) IsTopLevel() bool {
	return c.ParentID == 0
}

// TableName 表名
func (Comment) TableName() string {
	return "comments"
}

// IsApproved 是否已审核
func (c *Comment) IsApproved() bool {
	return c.Status == "approved"
}

// IsPending 是否待审核
func (c *Comment) IsPending() bool {
	return c.Status == "pending"
}
