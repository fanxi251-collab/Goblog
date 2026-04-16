package front

import (
	"Goblog/internal/model"
	"Goblog/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// MessageHandler 留言板处理器
type MessageHandler struct {
	commentService *service.CommentService
}

// NewMessageHandler 创建留言板处理器
func NewMessageHandler(commentSvc *service.CommentService) *MessageHandler {
	return &MessageHandler{commentService: commentSvc}
}

// Index 留言板页面
func (h *MessageHandler) Index(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize := 20

	// 获取已审核的留言（post_id = 0 表示留言板）
	comments, total, _ := h.commentService.GetMessageBoard(page, pageSize)

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	c.HTML(http.StatusOK, "message.html", gin.H{
		"title":      "留言板 - 灵序之夏",
		"comments":   comments,
		"page":       page,
		"totalPages": totalPages,
		"total":      total,
	})
}

// Submit 提交留言
func (h *MessageHandler) Submit(c *gin.Context) {
	nickname := c.PostForm("nickname")
	email := c.PostForm("email")
	content := c.PostForm("content")

	// 验证
	if nickname == "" || content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "昵称和内容不能为空"})
		return
	}

	// 创建评论（XSS清洗在Service层自动处理）
	comment := &model.Comment{
		Nickname:  nickname,
		Email:     email,
		Content:   content,
		PostID:    0,         // 0表示留言板
		Status:    "pending", // 需要审核
		IP:        c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
	}

	err := h.commentService.Create(comment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "留言已提交，等待审核后显示",
	})
}
