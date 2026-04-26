package model

// Column 专栏模型
type Column struct {
	ID          uint   `gorm:"primarykey" json:"id"`
	Name        string `gorm:"not null;size:50" json:"name"`
	Slug        string `gorm:"unique;not null;size:50" json:"slug"`
	Description string `gorm:"size:255" json:"description"`
	Sort        int    `gorm:"default:0" json:"sort"`
	CreatedAt   int64  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   int64  `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 表名
func (Column) TableName() string {
	return "columns"
}
