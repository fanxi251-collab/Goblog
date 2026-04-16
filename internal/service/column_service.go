package service

import (
	"Goblog/internal/model"
	"Goblog/internal/repository"

	"errors"
)

// ColumnService 专栏服务
type ColumnService struct {
	columnRepo *repository.ColumnRepository
}

// NewColumnService 创建专栏服务
func NewColumnService(columnRepo *repository.ColumnRepository) *ColumnService {
	return &ColumnService{columnRepo: columnRepo}
}

// Create 创建专栏
func (s *ColumnService) Create(column *model.Column) error {
	// 检查slug唯一性
	existing, _ := s.columnRepo.GetBySlug(column.Slug)
	if existing != nil {
		return errors.New("slug已存在")
	}
	return s.columnRepo.Create(column)
}

// GetByID 根据ID获取专栏
func (s *ColumnService) GetByID(id uint) (*model.Column, error) {
	return s.columnRepo.GetByID(id)
}

// GetBySlug 根据slug获取专栏
func (s *ColumnService) GetBySlug(slug string) (*model.Column, error) {
	return s.columnRepo.GetBySlug(slug)
}

// GetAll 获取所有专栏
func (s *ColumnService) GetAll() ([]model.Column, error) {
	return s.columnRepo.GetAll()
}

// GetParents 获取父级专栏
func (s *ColumnService) GetParents() ([]model.Column, error) {
	return s.columnRepo.GetParents()
}

// GetChildren 获取子级专栏
func (s *ColumnService) GetChildren(parentID uint) ([]model.Column, error) {
	return s.columnRepo.GetChildren(parentID)
}

// Update 更新专栏
func (s *ColumnService) Update(column *model.Column) error {
	// 检查slug唯一性（排除自己）
	existing, _ := s.columnRepo.GetBySlug(column.Slug)
	if existing != nil && existing.ID != column.ID {
		return errors.New("slug已存在")
	}
	return s.columnRepo.Update(column)
}

// Delete 删除专栏
func (s *ColumnService) Delete(id uint) error {
	return s.columnRepo.Delete(id)
}
