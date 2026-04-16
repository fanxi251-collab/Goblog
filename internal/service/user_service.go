package service

import (
	"Goblog/internal/model"
	"Goblog/internal/repository"

	"errors"
)

// UserService 用户服务
type UserService struct {
	userRepo *repository.UserRepository
}

// NewUserService 创建用户服务
func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// Login 用户登录
func (s *UserService) Login(username, password string) (*model.User, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	if !user.ValidatePassword(password) {
		return nil, errors.New("密码错误")
	}
	return user, nil
}

// GetByID 根据ID获取用户
func (s *UserService) GetByID(id uint) (*model.User, error) {
	return s.userRepo.GetByID(id)
}

// GetByUsername 根据用户名获取用户
func (s *UserService) GetByUsername(username string) (*model.User, error) {
	return s.userRepo.GetByUsername(username)
}

// Update 更新用户
func (s *UserService) Update(user *model.User) error {
	return s.userRepo.Update(user)
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(id uint, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}
	if !user.ValidatePassword(oldPassword) {
		return errors.New("原密码错误")
	}
	user.Password = newPassword
	return s.userRepo.Update(user)
}
