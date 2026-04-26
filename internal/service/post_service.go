package service

import (
	"Goblog/internal/model"
	"Goblog/internal/repository"

	"errors"
)

// PostService 文章服务
type PostService struct {
	postRepo *repository.PostRepository
}

// NewPostService 创建文章服务
func NewPostService(postRepo *repository.PostRepository) *PostService {
	return &PostService{postRepo: postRepo}
}

// Create 创建文章
func (s *PostService) Create(post *model.Post) error {
	// 检查slug唯一性
	if post.Slug != "" {
		existing, _ := s.postRepo.GetBySlug(post.Slug)
		if existing != nil {
			return errors.New("slug已存在")
		}
	}
	// XSS清洗文章内容
	if model.GetXSSEnabled() {
		post.Content = model.SanitizeArticle(post.Content)
	}
	return s.postRepo.Create(post)
}

// GetByID 根据ID获取文章
func (s *PostService) GetByID(id uint) (*model.Post, error) {
	return s.postRepo.GetByID(id)
}

// GetBySlug 根据slug获取文章
func (s *PostService) GetBySlug(slug string) (*model.Post, error) {
	return s.postRepo.GetBySlug(slug)
}

// GetByColumn 根据专栏获取文章
func (s *PostService) GetByColumn(columnID uint, status string, page, pageSize int) ([]model.Post, int64, error) {
	offset := (page - 1) * pageSize
	return s.postRepo.GetByColumn(columnID, status, offset, pageSize)
}

// Search 搜索文章（标题+内容）
func (s *PostService) Search(keyword string, status string, page, pageSize int) ([]model.Post, int64, error) {
	return s.postRepo.GetBySearch(keyword, status, (page-1)*pageSize, pageSize)
}

// SearchInColumn 在专栏内搜索文章
func (s *PostService) SearchInColumn(columnID uint, keyword string, status string, page, pageSize int) ([]model.Post, int64, error) {
	return s.postRepo.SearchInColumn(columnID, keyword, status, (page-1)*pageSize, pageSize)
}

// GetAll 获取所有文章
func (s *PostService) GetAll(page, pageSize int) ([]model.Post, int64, error) {
	return s.postRepo.GetAll((page-1)*pageSize, pageSize)
}

// GetByStatus 根据状态获取文章
func (s *PostService) GetByStatus(status string, page, pageSize int) ([]model.Post, int64, error) {
	return s.postRepo.GetByStatus(status, (page-1)*pageSize, pageSize)
}

// Update 更新文章
func (s *PostService) Update(post *model.Post) error {
	// 检查slug唯一性（排除自己）
	if post.Slug != "" {
		existing, _ := s.postRepo.GetBySlug(post.Slug)
		if existing != nil && existing.ID != post.ID {
			return errors.New("slug已存在")
		}
	}
	// XSS清洗文章内容
	if model.GetXSSEnabled() {
		post.Content = model.SanitizeArticle(post.Content)
	}
	return s.postRepo.Update(post)
}

// Delete 删除文章
func (s *PostService) Delete(id uint) error {
	return s.postRepo.Delete(id)
}

// Publish 发布文章
func (s *PostService) Publish(id uint) error {
	post, err := s.postRepo.GetByID(id)
	if err != nil {
		return err
	}
	post.Status = "published"
	return s.postRepo.Update(post)
}

// Unpublish 取消发布
func (s *PostService) Unpublish(id uint) error {
	post, err := s.postRepo.GetByID(id)
	if err != nil {
		return err
	}
	post.Status = "draft"
	return s.postRepo.Update(post)
}

// GetStats 获取统计数据
func (s *PostService) GetStats() (int64, int64, int64, error) {
	return s.postRepo.GetStats()
}

// GetLatestUpdateTime 获取最后更新时间
func (s *PostService) GetLatestUpdateTime() (int64, error) {
	return s.postRepo.GetLatestUpdateTime()
}
