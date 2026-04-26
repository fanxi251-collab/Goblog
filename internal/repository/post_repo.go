package repository

import (
	"Goblog/internal/model"
	"time"

	"gorm.io/gorm"
)

// PostRepository 文章仓库
type PostRepository struct {
	db *gorm.DB
}

// NewPostRepository 创建文章仓库
func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

// Create 创建文章
func (r *PostRepository) Create(post *model.Post) error {
	return r.db.Create(post).Error
}

// GetByID 根据ID获取文章
func (r *PostRepository) GetByID(id uint) (*model.Post, error) {
	var post model.Post
	err := r.db.First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// GetBySlug 根据slug获取文章
func (r *PostRepository) GetBySlug(slug string) (*model.Post, error) {
	var post model.Post
	err := r.db.Where("slug = ?", slug).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// GetByColumn 根据专栏获取文章
func (r *PostRepository) GetByColumn(columnID uint, status string, offset, limit int) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64

	query := r.db.Model(&model.Post{})
	if columnID > 0 {
		query = query.Where("column_id = ?", columnID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset(offset).Limit(limit).Order("is_top DESC, created_at DESC").Find(&posts).Error
	return posts, total, err
}

// GetBySearch 根据关键词搜索文章（标题+简介）
func (r *PostRepository) GetBySearch(keyword string, status string, offset, limit int) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64

	query := r.db.Model(&model.Post{})

	// 状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 关键词搜索（标题+简介）
	if keyword != "" {
		keyword = "%" + keyword + "%"
		query = query.Where("title LIKE ? OR excerpt LIKE ?", keyword, keyword)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset(offset).Limit(limit).Order("is_top DESC, created_at DESC").Find(&posts).Error
	return posts, total, err
}

// SearchInColumn 根据关键词在指定专栏中搜索文章
func (r *PostRepository) SearchInColumn(columnID uint, keyword string, status string, offset, limit int) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64

	query := r.db.Model(&model.Post{})

	// 专栏筛选
	if columnID > 0 {
		query = query.Where("column_id = ?", columnID)
	}

	// 状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 关键词搜索（标题+简介）
	if keyword != "" {
		keyword = "%" + keyword + "%"
		query = query.Where("title LIKE ? OR excerpt LIKE ?", keyword, keyword)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset(offset).Limit(limit).Order("is_top DESC, created_at DESC").Find(&posts).Error
	return posts, total, err
}

// GetAll 获取所有文章（后台用）
func (r *PostRepository) GetAll(offset, limit int) ([]model.Post, int64, error) {
	return r.GetByColumn(0, "", offset, limit)
}

// GetByStatus 根据状态获取文章
func (r *PostRepository) GetByStatus(status string, offset, limit int) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64

	query := r.db.Model(&model.Post{}).Where("status = ?", status)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&posts).Error
	return posts, total, err
}

// Update 更新文章
func (r *PostRepository) Update(post *model.Post) error {
	post.UpdatedAt = time.Now().Unix()
	return r.db.Save(post).Error
}

// Delete 删除文章
func (r *PostRepository) Delete(id uint) error {
	return r.db.Delete(&model.Post{}, id).Error
}

// IncrLikeCount 增加点赞次数
func (r *PostRepository) IncrLikeCount(id uint) error {
	return r.db.Model(&model.Post{}).Where("id = ?", id).
		UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error
}

// DecrLikeCount 减少点赞次数
func (r *PostRepository) DecrLikeCount(id uint) error {
	return r.db.Model(&model.Post{}).Where("id = ? AND like_count > ?", id, 0).
		UpdateColumn("like_count", gorm.Expr("like_count - ?", 1)).Error
}

// IncrCommentCount 增加评论次数
func (r *PostRepository) IncrCommentCount(id uint) error {
	return r.db.Model(&model.Post{}).Where("id = ?", id).
		UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1)).Error
}

// GetStats 获取统计数据（文章总数、点赞总数、评论总数）
func (r *PostRepository) GetStats() (int64, int64, int64, error) {
	var posts []model.Post
	var total int64

	// 获取文章总数
	if err := r.db.Model(&model.Post{}).Count(&total).Error; err != nil {
		return 0, 0, 0, err
	}

	// 获取所有文章（统计用）
	if err := r.db.Find(&posts).Error; err != nil {
		return total, 0, 0, err
	}

	// 统计点赞和评论
	var totalLikes, totalComments int64
	for _, post := range posts {
		totalLikes += int64(post.LikeCount)
		totalComments += int64(post.CommentCount)
	}

	return total, totalLikes, totalComments, nil
}

// GetLatestUpdateTime 获取最后更新时间
func (r *PostRepository) GetLatestUpdateTime() (int64, error) {
	var post model.Post
	err := r.db.Where("status = ?", "published").Order("updated_at DESC").First(&post).Error
	if err != nil {
		return 0, err
	}
	return post.UpdatedAt, nil
}
