package model

// Comment 评论模型
type Comment struct {
	ID        uint   `gorm:"primarykey" json:"id"`
	Nickname  string `gorm:"not null;size:50" json:"nickname"`
	Email     string `gorm:"size:100" json:"email"`
	Content   string `gorm:"not null;type:text" json:"content"`
	PostID    uint   `gorm:"not null" json:"post_id"`               // 关联文章ID，0表示留言板
	Status    string `gorm:"size:20;default:pending" json:"status"` // pending/approved/rejected
	IP        string `gorm:"size:50" json:"ip"`
	UserAgent string `gorm:"size:255" json:"user_agent"`
	CreatedAt int64  `gorm:"autoCreateTime" json:"created_at"`
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
