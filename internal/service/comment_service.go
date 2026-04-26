package service

import (
	"Goblog/internal/model"
	"Goblog/internal/repository"
)

// CommentService 评论服务
type CommentService struct {
	commentRepo *repository.CommentRepository
}

// NewCommentService 创建评论服务
func NewCommentService(commentRepo *repository.CommentRepository) *CommentService {
	return &CommentService{commentRepo: commentRepo}
}

// Create 创建评论（带XSS清洗）
func (s *CommentService) Create(comment *model.Comment) error {
	// XSS清洗评论内容
	comment.Content = model.SanitizeComment(comment.Content)
	comment.Nickname = model.SanitizeComment(comment.Nickname)
	return s.commentRepo.Create(comment)
}

// GetByID 根据ID获取评论
func (s *CommentService) GetByID(id uint) (*model.Comment, error) {
	return s.commentRepo.GetByID(id)
}

// GetByPostID 根据文章ID获取评论
func (s *CommentService) GetByPostID(postID uint, status string, page, pageSize int) ([]model.Comment, int64, error) {
	offset := (page - 1) * pageSize
	return s.commentRepo.GetByPostID(postID, status, offset, pageSize)
}

// GetApproved 获取已审核评论
func (s *CommentService) GetApproved(postID uint, page, pageSize int) ([]model.Comment, int64, error) {
	return s.commentRepo.GetApproved(postID, (page-1)*pageSize, pageSize)
}

// GetPending 获取待审核评论
func (s *CommentService) GetPending(page, pageSize int) ([]model.Comment, int64, error) {
	return s.commentRepo.GetPending((page-1)*pageSize, pageSize)
}

// GetMessageBoard 获取留言板评论
func (s *CommentService) GetMessageBoard(page, pageSize int) ([]model.Comment, int64, error) {
	return s.commentRepo.GetMessageBoard((page-1)*pageSize, pageSize)
}

// GetAll 获取所有评论
func (s *CommentService) GetAll(page, pageSize int) ([]model.Comment, int64, error) {
	return s.commentRepo.GetAll((page-1)*pageSize, pageSize)
}

// Update 更新评论
func (s *CommentService) Update(comment *model.Comment) error {
	// XSS清洗评论内容
	comment.Content = model.SanitizeComment(comment.Content)
	return s.commentRepo.Update(comment)
}

// Delete 删除评论
func (s *CommentService) Delete(id uint) error {
	return s.commentRepo.Delete(id)
}

// Approve 审核通过
func (s *CommentService) Approve(id uint) error {
	comment, err := s.commentRepo.GetByID(id)
	if err != nil {
		return err
	}
	comment.Status = "approved"
	return s.commentRepo.Update(comment)
}

// Reject 拒绝
func (s *CommentService) Reject(id uint) error {
	comment, err := s.commentRepo.GetByID(id)
	if err != nil {
		return err
	}
	comment.Status = "rejected"
	return s.commentRepo.Update(comment)
}

// BatchApprove 批量审核通过
func (s *CommentService) BatchApprove(ids []uint) error {
	return s.commentRepo.BatchUpdateStatus(ids, "approved")
}

// BatchReject 批量拒绝
func (s *CommentService) BatchReject(ids []uint) error {
	return s.commentRepo.BatchUpdateStatus(ids, "rejected")
}

// GetReplies 获取评论的所有回复
func (s *CommentService) GetReplies(parentID uint) ([]model.Comment, error) {
	return s.commentRepo.GetReplies(parentID)
}

// GetMessageBoardWithReplies 获取留言板评论（含回复树）
func (s *CommentService) GetMessageBoardWithReplies(page, pageSize int) ([]model.Comment, int64, error) {
	offset := (page - 1) * pageSize
	// 获取顶层评论
	comments, total, err := s.commentRepo.GetTopLevelComments(0, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// 为每个顶层评论加载回复
	for i := range comments {
		replies, _ := s.commentRepo.GetReplies(comments[i].ID)
		comments[i].Replies = replies
	}

	return comments, total, nil
}

// DeleteWithReplies 级联删除评论及其回复
func (s *CommentService) DeleteWithReplies(id uint) error {
	return s.commentRepo.DeleteWithReplies(id)
}

// GetArticleComments 获取文章评论（后台）
func (s *CommentService) GetArticleComments(status string, page, pageSize int) ([]model.Comment, int64, error) {
	offset := (page - 1) * pageSize
	return s.commentRepo.GetArticleComments(status, offset, pageSize)
}

// GetMessageBoardComments 获取留言板评论（后台）
func (s *CommentService) GetMessageBoardComments(status string, page, pageSize int) ([]model.Comment, int64, error) {
	offset := (page - 1) * pageSize
	return s.commentRepo.GetMessageBoardComments(status, offset, pageSize)
}
