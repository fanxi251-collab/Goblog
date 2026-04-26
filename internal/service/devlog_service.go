package service

import (
	"Goblog/internal/model"
	"Goblog/internal/repository"
)

// DevlogService 开发日志服务
type DevlogService struct {
	devlogRepo *repository.DevlogRepository
}

// NewDevlogService 创建开发日志服务
func NewDevlogService(devlogRepo *repository.DevlogRepository) *DevlogService {
	return &DevlogService{devlogRepo: devlogRepo}
}

// Create 创建日志
func (s *DevlogService) Create(devlog *model.Devlog) error {
	return s.devlogRepo.Create(devlog)
}

// GetByID 根据ID获取日志
func (s *DevlogService) GetByID(id uint) (*model.Devlog, error) {
	return s.devlogRepo.GetByID(id)
}

// GetAll 获取所有日志
func (s *DevlogService) GetAll(page, pageSize int) ([]model.Devlog, int64, error) {
	offset := (page - 1) * pageSize
	return s.devlogRepo.GetAll(offset, pageSize)
}

// GetPublished 获取已发布日志
func (s *DevlogService) GetPublished(page, pageSize int) ([]model.Devlog, int64, error) {
	offset := (page - 1) * pageSize
	return s.devlogRepo.GetPublished(offset, pageSize)
}

// Update 更新日志
func (s *DevlogService) Update(devlog *model.Devlog) error {
	return s.devlogRepo.Update(devlog)
}

// Delete 删除日志
func (s *DevlogService) Delete(id uint) error {
	return s.devlogRepo.Delete(id)
}

// Publish 发布日志
func (s *DevlogService) Publish(id uint) error {
	devlog, err := s.devlogRepo.GetByID(id)
	if err != nil {
		return err
	}
	devlog.Status = "published"
	return s.devlogRepo.Update(devlog)
}

// Unpublish 下架日志
func (s *DevlogService) Unpublish(id uint) error {
	devlog, err := s.devlogRepo.GetByID(id)
	if err != nil {
		return err
	}
	devlog.Status = "draft"
	return s.devlogRepo.Update(devlog)
}
