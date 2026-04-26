package model

import (
	"time"
)

// Devlog 开发日志模型
type Devlog struct {
	ID          uint   `gorm:"primarykey" json:"id"`
	Title       string `gorm:"not null;size:200" json:"title"`
	Description string `gorm:"type:text" json:"description"`        // 日志描述
	Date        int64  `gorm:"not null" json:"date"`                // 发布日期（日期，不带时间）
	Status      string `gorm:"size:20;default:draft" json:"status"` // 状态：draft/published
	CreatedAt   int64  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   int64  `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 表名
func (Devlog) TableName() string {
	return "devlogs"
}

// FormatDate 格式化日期
func (d *Devlog) FormatDate() string {
	if d.Date == 0 {
		return ""
	}
	return time.Unix(d.Date, 0).Format("2006-01-02")
}

// IsPublished 检查是否已发布
func (d *Devlog) IsPublished() bool {
	return d.Status == "published"
}
