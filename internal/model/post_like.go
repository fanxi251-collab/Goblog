package model

// PostLike 文章点赞记录
type PostLike struct {
	ID        uint   `gorm:"primarykey" json:"id"`
	PostID    uint   `gorm:"not null" json:"post_id"`    // 文章ID
	VisitorID uint   `gorm:"not null" json:"visitor_id"` // 访客ID
	IP        string `gorm:"size:50" json:"ip"`          // IP地址
	CreatedAt int64  `gorm:"autoCreateTime" json:"created_at"`
}

// TableName 表名
func (PostLike) TableName() string {
	return "post_likes"
}

// HasLiked 检查是否已点赞（根据IP）
func (pl *PostLike) HasLiked() bool {
	return pl != nil && pl.ID > 0
}
