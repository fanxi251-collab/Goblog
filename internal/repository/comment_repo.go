package repository

import (
	"Goblog/internal/model"

	"gorm.io/gorm"
)

// CommentRepository 评论仓库
type CommentRepository struct {
	db *gorm.DB
}

// NewCommentRepository 创建评论仓库
func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

// Create 创建评论
func (r *CommentRepository) Create(comment *model.Comment) error {
	return r.db.Create(comment).Error
}

// GetByID 根据ID获取评论
func (r *CommentRepository) GetByID(id uint) (*model.Comment, error) {
	var comment model.Comment
	err := r.db.First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// GetByPostID 根据文章ID获取评论
func (r *CommentRepository) GetByPostID(postID uint, status string, offset, limit int) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64

	query := r.db.Model(&model.Comment{})
	if postID > 0 {
		query = query.Where("post_id = ?", postID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&comments).Error
	return comments, total, err
}

// GetApproved 获取已审核评论（前台用）
func (r *CommentRepository) GetApproved(postID uint, offset, limit int) ([]model.Comment, int64, error) {
	return r.GetByPostID(postID, "approved", offset, limit)
}

// GetPending 获取待审核评论（后台用）
func (r *CommentRepository) GetPending(offset, limit int) ([]model.Comment, int64, error) {
	return r.GetByPostID(0, "pending", offset, limit)
}

// GetMessageBoard 获取留言板评论
func (r *CommentRepository) GetMessageBoard(offset, limit int) ([]model.Comment, int64, error) {
	// postID = 0 表示留言板
	var comments []model.Comment
	var total int64

	err := r.db.Model(&model.Comment{}).Where("post_id = ? AND status = ?", 0, "approved").
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Where("post_id = ? AND status = ?", 0, "approved").
		Offset(offset).Limit(limit).Order("created_at DESC").Find(&comments).Error
	return comments, total, err
}

// GetAll 获取所有评论（后台用）
func (r *CommentRepository) GetAll(offset, limit int) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64

	err := r.db.Model(&model.Comment{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&comments).Error
	return comments, total, err
}

// Update 更新评论
func (r *CommentRepository) Update(comment *model.Comment) error {
	return r.db.Save(comment).Error
}

// Delete 删除评论
func (r *CommentRepository) Delete(id uint) error {
	return r.db.Delete(&model.Comment{}, id).Error
}

// BatchUpdateStatus 批量更新状态
func (r *CommentRepository) BatchUpdateStatus(ids []uint, status string) error {
	return r.db.Model(&model.Comment{}).Where("id IN ?", ids).Update("status", status).Error
}
