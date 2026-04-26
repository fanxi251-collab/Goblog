package admin

import (
	"Goblog/internal/model"
	"Goblog/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// CommentHandler 评论处理器
type CommentHandler struct {
	commentService *service.CommentService
}

// NewCommentHandler 创建评论处理器
func NewCommentHandler(commentSvc *service.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentSvc}
}

// List 留言板评论列表（PostID = 0）
func (h *CommentHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize := 20
	status := c.Query("status")

	var comments []model.Comment
	var total int64
	var err error

	// 根据状态筛选（留言板评论）
	if status == "" || status == "pending" {
		comments, total, err = h.commentService.GetMessageBoardComments("pending", page, pageSize)
	} else if status == "approved" {
		comments, total, err = h.commentService.GetMessageBoardComments("approved", page, pageSize)
	} else if status == "rejected" {
		comments, total, err = h.commentService.GetMessageBoardComments("rejected", page, pageSize)
	} else {
		// 获取全部留言板评论
		comments, total, err = h.commentService.GetMessageBoardComments("", page, pageSize)
	}

	if err != nil {
		comments = []model.Comment{}
		total = 0
	}

	// 为回复加载父评论昵称和格式化时间
	for i := range comments {
		if comments[i].ParentID > 0 {
			parent, _ := h.commentService.GetByID(comments[i].ParentID)
			if parent != nil {
				comments[i].ParentNickname = parent.Nickname
			}
		}
		comments[i].FormatTime = comments[i].FormatCreatedAt()
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	c.HTML(http.StatusOK, "comment_list.html", gin.H{
		"title":       "留言管理",
		"comments":    comments,
		"page":        page,
		"totalPages":  totalPages,
		"total":       total,
		"status":      status,
		"commentType": "message",
	})
}

// ArticleComments 文章评论列表（PostID > 0）
func (h *CommentHandler) ArticleComments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize := 20
	status := c.Query("status")

	var comments []model.Comment
	var total int64
	var err error

	// 根据状态筛选（文章评论）
	if status == "" || status == "pending" {
		comments, total, err = h.commentService.GetArticleComments("pending", page, pageSize)
	} else if status == "approved" {
		comments, total, err = h.commentService.GetArticleComments("approved", page, pageSize)
	} else if status == "rejected" {
		comments, total, err = h.commentService.GetArticleComments("rejected", page, pageSize)
	} else {
		// 获取全部文章评论
		comments, total, err = h.commentService.GetArticleComments("", page, pageSize)
	}

	if err != nil {
		comments = []model.Comment{}
		total = 0
	}

	// 为回复加载父评论昵称和格式化时间
	for i := range comments {
		if comments[i].ParentID > 0 {
			parent, _ := h.commentService.GetByID(comments[i].ParentID)
			if parent != nil {
				comments[i].ParentNickname = parent.Nickname
			}
		}
		comments[i].FormatTime = comments[i].FormatCreatedAt()
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	c.HTML(http.StatusOK, "article_comment_list.html", gin.H{
		"title":       "文章评论管理",
		"comments":    comments,
		"page":        page,
		"totalPages":  totalPages,
		"total":       total,
		"status":      status,
		"commentType": "article",
	})
}

// Approve 审核通过
func (h *CommentHandler) Approve(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	err = h.commentService.Approve(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// Reject 拒绝
func (h *CommentHandler) Reject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	err = h.commentService.Reject(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// BatchApprove 批量审核通过
func (h *CommentHandler) BatchApprove(c *gin.Context) {
	ids := c.PostForm("ids")
	if ids == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择要审核的评论"})
		return
	}

	var idSlice []uint
	for _, idStr := range SplitComma(ids) {
		if id, err := strconv.ParseUint(idStr, 10, 32); err == nil {
			idSlice = append(idSlice, uint(id))
		}
	}

	if len(idSlice) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	err := h.commentService.BatchApprove(idSlice)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// BatchReject 批量拒绝
func (h *CommentHandler) BatchReject(c *gin.Context) {
	ids := c.PostForm("ids")
	if ids == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择要拒绝的评论"})
		return
	}

	var idSlice []uint
	for _, idStr := range SplitComma(ids) {
		if id, err := strconv.ParseUint(idStr, 10, 32); err == nil {
			idSlice = append(idSlice, uint(id))
		}
	}

	if len(idSlice) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	err := h.commentService.BatchReject(idSlice)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// Delete 删除评论（级联删除回复）
func (h *CommentHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	err = h.commentService.DeleteWithReplies(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// SplitComma 分割逗号分隔的ID
func SplitComma(s string) []string {
	if s == "" {
		return []string{}
	}
	var result []string
	var current []byte
	for _, c := range s {
		if c == ',' {
			if len(current) > 0 {
				result = append(result, string(current))
				current = nil
			}
		} else {
			current = append(current, byte(c))
		}
	}
	if len(current) > 0 {
		result = append(result, string(current))
	}
	return result
}

// FormatTime 格式化时间戳
func FormatTime(unix int64) string {
	if unix == 0 {
		return ""
	}
	t := time.Unix(unix, 0)
	return t.Format("2006-01-02 15:04:05")
}
