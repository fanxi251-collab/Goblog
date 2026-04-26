package repository

import (
	"Goblog/internal/model"

	"gorm.io/gorm"
)

// VisitorRepository 访客仓库
type VisitorRepository struct {
	db *gorm.DB
}

// NewVisitorRepository 创建访客仓库
func NewVisitorRepository(db *gorm.DB) *VisitorRepository {
	return &VisitorRepository{db: db}
}

// Create 创建访客
func (r *VisitorRepository) Create(visitor *model.Visitor) error {
	return r.db.Create(visitor).Error
}

// GetByToken 根据Token获取访客
func (r *VisitorRepository) GetByToken(token string) (*model.Visitor, error) {
	var visitor model.Visitor
	err := r.db.Where("token = ?", token).First(&visitor).Error
	if err != nil {
		return nil, err
	}
	return &visitor, nil
}

// GetByIP 根据IP获取最近访问的访客
func (r *VisitorRepository) GetByIP(ip string) (*model.Visitor, error) {
	var visitor model.Visitor
	err := r.db.Where("ip = ?", ip).Order("created_at DESC").First(&visitor).Error
	if err != nil {
		return nil, err
	}
	return &visitor, nil
}

// Update 更新访客
func (r *VisitorRepository) Update(visitor *model.Visitor) error {
	return r.db.Save(visitor).Error
}

// GetLastCommentTime 获取最后评论时间（用于频率限制）
func (r *VisitorRepository) GetLastCommentTime(visitorID uint) (int64, error) {
	var comment model.Comment
	err := r.db.Where("post_id = 0").Order("created_at DESC").First(&comment).Error
	if err != nil {
		return 0, err
	}
	return comment.CreatedAt, nil
}
