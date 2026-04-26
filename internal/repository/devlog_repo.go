package repository

import (
	"Goblog/internal/model"
	"time"

	"gorm.io/gorm"
)

// DevlogRepository 开发日志仓库
type DevlogRepository struct {
	db *gorm.DB
}

// NewDevlogRepository 创建开发日志仓库
func NewDevlogRepository(db *gorm.DB) *DevlogRepository {
	return &DevlogRepository{db: db}
}

// Create 创建日志
func (r *DevlogRepository) Create(devlog *model.Devlog) error {
	return r.db.Create(devlog).Error
}

// GetByID 根据ID获取日志
func (r *DevlogRepository) GetByID(id uint) (*model.Devlog, error) {
	var devlog model.Devlog
	err := r.db.First(&devlog, id).Error
	if err != nil {
		return nil, err
	}
	return &devlog, nil
}

// GetAll 获取所有日志（按日期倒序）
func (r *DevlogRepository) GetAll(offset, limit int) ([]model.Devlog, int64, error) {
	var devlogs []model.Devlog
	var total int64

	err := r.db.Model(&model.Devlog{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Offset(offset).Limit(limit).Order("date DESC").Find(&devlogs).Error
	return devlogs, total, err
}

// GetPublished 获取已发布日志
func (r *DevlogRepository) GetPublished(offset, limit int) ([]model.Devlog, int64, error) {
	var devlogs []model.Devlog
	var total int64

	err := r.db.Model(&model.Devlog{}).Where("status = ?", "published").Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Where("status = ?", "published").Offset(offset).Limit(limit).Order("date DESC").Find(&devlogs).Error
	return devlogs, total, err
}

// Update 更新日志
func (r *DevlogRepository) Update(devlog *model.Devlog) error {
	devlog.UpdatedAt = time.Now().Unix()
	return r.db.Save(devlog).Error
}

// Delete 删除日志
func (r *DevlogRepository) Delete(id uint) error {
	return r.db.Delete(&model.Devlog{}, id).Error
}
