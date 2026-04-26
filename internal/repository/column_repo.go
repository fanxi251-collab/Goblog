package repository

import (
	"Goblog/internal/model"

	"gorm.io/gorm"
)

// ColumnRepository 专栏仓库
type ColumnRepository struct {
	db *gorm.DB
}

// NewColumnRepository 创建专栏仓库
func NewColumnRepository(db *gorm.DB) *ColumnRepository {
	return &ColumnRepository{db: db}
}

// Create 创建专栏
func (r *ColumnRepository) Create(column *model.Column) error {
	return r.db.Create(column).Error
}

// GetByID 根据ID获取专栏
func (r *ColumnRepository) GetByID(id uint) (*model.Column, error) {
	var column model.Column
	err := r.db.First(&column, id).Error
	if err != nil {
		return nil, err
	}
	return &column, nil
}

// GetBySlug 根据slug获取专栏
func (r *ColumnRepository) GetBySlug(slug string) (*model.Column, error) {
	var column model.Column
	err := r.db.Where("slug = ?", slug).First(&column).Error
	if err != nil {
		return nil, err
	}
	return &column, nil
}

// GetAll 获取所有专栏
func (r *ColumnRepository) GetAll() ([]model.Column, error) {
	var columns []model.Column
	err := r.db.Order("sort ASC, created_at DESC").Find(&columns).Error
	return columns, err
}

// Update 更新专栏
func (r *ColumnRepository) Update(column *model.Column) error {
	return r.db.Save(column).Error
}

// Delete 删除专栏
func (r *ColumnRepository) Delete(id uint) error {
	return r.db.Delete(&model.Column{}, id).Error
}
