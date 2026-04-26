package repository

import (
	"Goblog/internal/model"

	"gorm.io/gorm"
)

// PostLikeRepository 文章点赞仓库
type PostLikeRepository struct {
	db *gorm.DB
}

// NewPostLikeRepository 创建文章点赞仓库
func NewPostLikeRepository(db *gorm.DB) *PostLikeRepository {
	return &PostLikeRepository{db: db}
}

// Create 创建点赞记录
func (r *PostLikeRepository) Create(like *model.PostLike) error {
	return r.db.Create(like).Error
}

// GetByPostIDAndVisitor 根据文章ID和访客ID获取点赞记录
func (r *PostLikeRepository) GetByPostIDAndVisitor(postID, visitorID uint) (*model.PostLike, error) {
	var like model.PostLike
	err := r.db.Where("post_id = ? AND visitor_id = ?", postID, visitorID).First(&like).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &like, nil
}

// GetByPostIDAndIP 根据文章ID和IP获取点赞记录
func (r *PostLikeRepository) GetByPostIDAndIP(postID uint, ip string) (*model.PostLike, error) {
	var like model.PostLike
	err := r.db.Where("post_id = ? AND ip = ?", postID, ip).First(&like).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &like, nil
}

// Delete 删除点赞记录
func (r *PostLikeRepository) Delete(id uint) error {
	return r.db.Delete(&model.PostLike{}, id).Error
}

// DeleteByPostIDAndVisitor 根据文章ID和访客ID删除点赞记录
func (r *PostLikeRepository) DeleteByPostIDAndVisitor(postID, visitorID uint) error {
	return r.db.Where("post_id = ? AND visitor_id = ?", postID, visitorID).Delete(&model.PostLike{}).Error
}
