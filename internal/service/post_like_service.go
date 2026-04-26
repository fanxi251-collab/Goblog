package service

import (
	"Goblog/internal/model"
	"Goblog/internal/repository"
)

// PostLikeService 文章点赞服务
type PostLikeService struct {
	postLikeRepo *repository.PostLikeRepository
	postRepo     *repository.PostRepository
}

// NewPostLikeService 创建文章点赞服务
func NewPostLikeService(postLikeRepo *repository.PostLikeRepository, postRepo *repository.PostRepository) *PostLikeService {
	return &PostLikeService{
		postLikeRepo: postLikeRepo,
		postRepo:     postRepo,
	}
}

// Like 点赞/取消点赞
func (s *PostLikeService) Like(postID uint, visitorID uint, ip string) (liked bool, err error) {
	// 检查是否已点赞
	existing, err := s.postLikeRepo.GetByPostIDAndVisitor(postID, visitorID)
	if err != nil {
		return false, err
	}

	if existing != nil {
		// 已点赞，取消点赞
		err = s.postLikeRepo.Delete(existing.ID)
		if err != nil {
			return false, err
		}
		// 减少点赞数
		s.postRepo.DecrLikeCount(postID)
		return false, nil
	}

	// 未点赞，添加点赞
	like := &model.PostLike{
		PostID:    postID,
		VisitorID: visitorID,
		IP:        ip,
	}
	err = s.postLikeRepo.Create(like)
	if err != nil {
		return false, err
	}
	// 增加点赞数
	s.postRepo.IncrLikeCount(postID)
	return true, nil
}

// HasLiked 检查是否已点赞
func (s *PostLikeService) HasLiked(postID uint, visitorID uint) (bool, error) {
	existing, err := s.postLikeRepo.GetByPostIDAndVisitor(postID, visitorID)
	if err != nil {
		return false, err
	}
	return existing != nil, nil
}

// LikeByIP IP点赞（备用方案，不需要登录）
func (s *PostLikeService) LikeByIP(postID uint, ip string) (liked bool, err error) {
	existing, err := s.postLikeRepo.GetByPostIDAndIP(postID, ip)
	if err != nil {
		return false, err
	}

	if existing != nil {
		// 已点赞，取消点赞
		err = s.postLikeRepo.Delete(existing.ID)
		if err != nil {
			return false, err
		}
		s.postRepo.DecrLikeCount(postID)
		return false, nil
	}

	// 未点赞，添加点赞
	like := &model.PostLike{
		PostID: postID,
		IP:     ip,
	}
	err = s.postLikeRepo.Create(like)
	if err != nil {
		return false, err
	}
	s.postRepo.IncrLikeCount(postID)
	return true, nil
}
