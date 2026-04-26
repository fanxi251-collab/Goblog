package service

import (
	"Goblog/internal/config"
	"Goblog/internal/model"
	"Goblog/internal/repository"
	"strings"
	"time"
)

// VisitorService 访客服务
type VisitorService struct {
	visitorRepo *repository.VisitorRepository
	commentRepo *repository.CommentRepository
}

// NewVisitorService 创建访客服务
func NewVisitorService(visitorRepo *repository.VisitorRepository, commentRepo *repository.CommentRepository) *VisitorService {
	return &VisitorService{
		visitorRepo: visitorRepo,
		commentRepo: commentRepo,
	}
}

// Register 注册新访客
func (s *VisitorService) Register(nickname, email, ip string) (*model.Visitor, error) {
	// 生成唯一Token
	token, err := model.GenerateToken()
	if err != nil {
		return nil, err
	}

	// XSS 清洗
	nickname = model.SanitizeComment(nickname)
	email = model.SanitizeComment(email)

	visitor := &model.Visitor{
		Token:     token,
		Nickname:  nickname,
		Email:     email,
		IP:        ip,
		CreatedAt: time.Now().Unix(),
	}

	if err := s.visitorRepo.Create(visitor); err != nil {
		return nil, err
	}

	return visitor, nil
}

// CheckToken 检查Token是否有效
func (s *VisitorService) CheckToken(token string) (*model.Visitor, error) {
	if token == "" {
		return nil, nil
	}

	visitor, err := s.visitorRepo.GetByToken(token)
	if err != nil {
		return nil, nil // Token无效返回nil
	}

	// 检查是否在30天有效期内
	if !visitor.IsValid() {
		return nil, nil
	}

	return visitor, nil
}

// UpdateInfo 更新访客信息
func (s *VisitorService) UpdateInfo(token, nickname, email string) (*model.Visitor, error) {
	visitor, err := s.visitorRepo.GetByToken(token)
	if err != nil {
		return nil, err
	}

	// XSS 清洗
	nickname = model.SanitizeComment(nickname)
	email = model.SanitizeComment(email)

	visitor.Nickname = nickname
	visitor.Email = email

	if err := s.visitorRepo.Update(visitor); err != nil {
		return nil, err
	}

	return visitor, nil
}

// CheckRateLimit 检查频率限制
func (s *VisitorService) CheckRateLimit(token, ip string) error {
	cfg := config.Get()
	rateLimit := cfg.Comment.RateLimit
	if rateLimit <= 0 {
		rateLimit = 3 // 默认3秒
	}

	// 检查最后评论时间
	var lastTime int64 = 0

	// 优先检查Token对应的访客
	if token != "" {
		if visitor, err := s.visitorRepo.GetByToken(token); err == nil {
			lastTime, _ = s.commentRepo.GetLastCommentTime(visitor.IP)
		}
	}

	// 如果没有Token或Token无记录，检查IP
	if lastTime == 0 {
		lastTime, _ = s.commentRepo.GetLastCommentTime(ip)
	}

	if lastTime > 0 {
		elapsed := time.Now().Unix() - lastTime
		if elapsed < int64(rateLimit) {
			return &RateLimitError{
				WaitSeconds: rateLimit - int(elapsed),
			}
		}
	}

	return nil
}

// RateLimitError 频率限制错误
type RateLimitError struct {
	WaitSeconds int
}

func (e *RateLimitError) Error() string {
	return "操作太频繁，请稍后再试"
}

// CheckBlockedWords 检查敏感词
func (s *VisitorService) CheckBlockedWords(content string) bool {
	cfg := config.Get()
	blockedWords := cfg.Comment.BlockedWords

	if len(blockedWords) == 0 {
		return false
	}

	content = strings.ToLower(content)
	for _, word := range blockedWords {
		if strings.Contains(content, strings.ToLower(word)) {
			return true // 包含敏感词
		}
	}

	return false
}

// ShouldAutoApprove 判断是否应该自动审核通过
func (s *VisitorService) ShouldAutoApprove(content string) bool {
	// 如果包含敏感词，不自动通过
	if s.CheckBlockedWords(content) {
		return false
	}

	cfg := config.Get()
	return cfg.Comment.AutoApprove
}
