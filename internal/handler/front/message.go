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
	postService    *service.PostService
	commentService *service.CommentService
	visitorService *service.VisitorService
}

// NewMessageHandler 创建留言板处理器
func NewMessageHandler(postSvc *service.PostService, commentSvc *service.CommentService, visitorSvc *service.VisitorService) *MessageHandler {
	return &MessageHandler{
		postService:    postSvc,
		commentService: commentSvc,
		visitorService: visitorSvc,
	}
}

// Index 留言板页面
func (h *MessageHandler) Index(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize := 20

	// 获取已审核的留言（含回复）
	comments, total, _ := h.commentService.GetMessageBoardWithReplies(page, pageSize)

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	// 尝试获取当前访客信息
	token := c.GetHeader("X-Visitor-Token")
	if token == "" {
		token, _ = c.Cookie("visitor_token")
	}

	var visitor *model.Visitor
	if token != "" {
		visitor, _ = h.visitorService.CheckToken(token)
	}

	c.HTML(http.StatusOK, "message.html", gin.H{
		"title":      "留言板 - 灵序之夏",
		"comments":   comments,
		"page":       page,
		"totalPages": totalPages,
		"total":      total,
		"visitor":    visitor,
	})
}

// Register 访客注册（首次访问弹窗后调用）
func (h *MessageHandler) Register(c *gin.Context) {
	nickname := c.PostForm("nickname")
	email := c.PostForm("email")

	// 验证昵称必填
	if nickname == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "昵称不能为空"})
		return
	}

	// 检查昵称是否已存在
	nicknameExists, _ := h.visitorService.CheckNicknameExists(nickname)
	if nicknameExists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "昵称已存在，请更换"})
		return
	}

	// 检查邮箱是否已存在（包含管理员邮箱）
	if email != "" {
		emailExists, _ := h.visitorService.CheckEmailExists(email)
		if emailExists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱已存在，请更换"})
			return
		}
	}

	ip := c.ClientIP()

	visitor, err := h.visitorService.Register(nickname, email, ip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败，请重试"})
		return
	}

	// 设置Cookie（30天有效期）
	c.SetCookie("visitor_token", visitor.Token, 30*24*60*60, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"token":   visitor.Token,
		"visitor": visitor,
	})
}

// Check 检查访客Token
func (h *MessageHandler) Check(c *gin.Context) {
	token := c.GetHeader("X-Visitor-Token")
	if token == "" {
		token, _ = c.Cookie("visitor_token")
	}

	if token == "" {
		c.JSON(http.StatusOK, gin.H{"is_login": false})
		return
	}

	visitor, err := h.visitorService.CheckToken(token)
	if err != nil || visitor == nil {
		// Token无效，删除Cookie
		c.SetCookie("visitor_token", "", -1, "/", "", false, true)
		c.JSON(http.StatusOK, gin.H{"is_login": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"is_login": true,
		"visitor":  visitor,
	})
}

// CheckExist 检查昵称/邮箱是否已存在（预检查接口）
func (h *MessageHandler) CheckExist(c *gin.Context) {
	nickname := c.Query("nickname")
	email := c.Query("email")

	result := gin.H{
		"nickname_exists": false,
		"email_exists":    false,
	}

	// 检查昵称
	if nickname != "" {
		exists, err := h.visitorService.CheckNicknameExists(nickname)
		if err == nil {
			result["nickname_exists"] = exists
		}
	}

	// 检查邮箱
	if email != "" {
		exists, err := h.visitorService.CheckEmailExists(email)
		if err == nil {
			result["email_exists"] = exists
		}
	}

	c.JSON(http.StatusOK, result)
}

// DeleteAllVisitors 删除所有访客（测试用）
func (h *MessageHandler) DeleteAllVisitors(c *gin.Context) {
	if err := h.visitorService.DeleteAllVisitors(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// UpdateInfo 更新访客信息
func (h *MessageHandler) UpdateInfo(c *gin.Context) {
	token := c.GetHeader("X-Visitor-Token")
	if token == "" {
		token, _ = c.Cookie("visitor_token")
	}

	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请先注册"})
		return
	}

	nickname := c.PostForm("nickname")
	email := c.PostForm("email")

	if nickname == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "昵称不能为空"})
		return
	}

	visitor, err := h.visitorService.UpdateInfo(token, nickname, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"visitor": visitor,
	})
}

// SubmitPostComment 提交文章评论（带Token验证、频率限制、敏感词过滤）
func (h *MessageHandler) SubmitPostComment(c *gin.Context) {
	var req struct {
		PostID   uint   `json:"post_id" binding:"required"`
		Content  string `json:"content" binding:"required"`
		ParentID uint   `json:"parent_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请输入评论内容"})
		return
	}

	// 验证文章存在（直接用 ID 查询）
	post, err := h.postService.GetByID(req.PostID)
	if err != nil || post == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文章不存在"})
		return
	}

	token := c.GetHeader("X-Visitor-Token")
	if token == "" {
		token, _ = c.Cookie("visitor_token")
	}

	ip := c.ClientIP()

	// 1. 频率限制检查
	if err := h.visitorService.CheckRateLimit(token, ip); err != nil {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
		return
	}

	// 2. 获取或创建访客（必须登录才能评论）
	var nickname, email string
	var visitor *model.Visitor

	if token != "" {
		visitor, _ = h.visitorService.CheckToken(token)
	}

	if visitor == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}

	nickname = visitor.Nickname
	email = visitor.Email

	// 3. 敏感词检查
	hasBlockedWords := h.visitorService.CheckBlockedWords(req.Content)
	autoApprove := h.visitorService.ShouldAutoApprove(req.Content)

	// 状态：敏感词→rejected，其他根据配置
	status := "pending"
	if !hasBlockedWords && autoApprove {
		status = "approved"
	}

	// 创建评论
	comment := &model.Comment{
		ParentID:  req.ParentID,
		Nickname:  nickname,
		Email:     email,
		Content:   req.Content,
		PostID:    req.PostID,
		Status:    status,
		IP:        ip,
		UserAgent: c.Request.UserAgent(),
	}

	err = h.commentService.Create(comment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 根据状态返回不同消息
	if status == "approved" {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "评论已发布",
			"status":  status,
		})
	} else if hasBlockedWords {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "评论包含不适当内容，已被拒绝",
			"status":  status,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "评论已提交，等待审核后显示",
			"status":  status,
		})
	}
}

// Submit 提交留言（带Token验证、频率限制、敏感词过滤）
func (h *MessageHandler) Submit(c *gin.Context) {
	content := c.PostForm("content")
	token := c.GetHeader("X-Visitor-Token")
	if token == "" {
		token, _ = c.Cookie("visitor_token")
	}

	ip := c.ClientIP()

	// 1. 频率限制检查
	if err := h.visitorService.CheckRateLimit(token, ip); err != nil {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
		return
	}

	// 验证内容
	if content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "内容不能为空"})
		return
	}

	// 解析 parent_id（回复功能）
	parentID := uint(0)
	if pidStr := c.PostForm("parent_id"); pidStr != "" {
		if pid, err := strconv.ParseUint(pidStr, 10, 32); err == nil {
			parentID = uint(pid)
		}
	}

	// 如果是回复，验证父评论存在且已审核
	if parentID > 0 {
		parentComment, err := h.commentService.GetByID(parentID)
		if err != nil || parentComment == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "父评论不存在"})
			return
		}
		if !parentComment.IsApproved() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "不能回复未审核的评论"})
			return
		}
		// 确保不是回复的回复的回复（限制3层）
		if parentComment.ParentID > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "回复层级已达上限"})
			return
		}
	}

	var nickname, email string
	var visitor *model.Visitor

	// 2. 获取或创建访客
	if token != "" {
		visitor, _ = h.visitorService.CheckToken(token)
	}

	if visitor == nil {
		// 没有Token，创建临时访客信息（昵称邮箱从表单获取）
		nickname = c.PostForm("nickname")
		email = c.PostForm("email")
		if nickname == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请先设置昵称"})
			return
		}
	} else {
		// 有Token，使用已保存的信息
		nickname = visitor.Nickname
		email = visitor.Email
	}

	// 3. 敏感词检查
	hasBlockedWords := h.visitorService.CheckBlockedWords(content)
	autoApprove := h.visitorService.ShouldAutoApprove(content)

	// 状态：敏感词→rejected，其他根据配置
	status := "pending"
	if !hasBlockedWords && autoApprove {
		status = "approved"
	}

	// 创建评论（XSS清洗在Service层自动处理）
	comment := &model.Comment{
		ParentID:  parentID,
		Nickname:  nickname,
		Email:     email,
		Content:   content,
		PostID:    0, // 0表示留言板
		Status:    status,
		IP:        ip,
		UserAgent: c.Request.UserAgent(),
	}

	err := h.commentService.Create(comment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 根据状态返回不同消息
	if status == "approved" {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "留言已发布",
			"status":  status,
		})
	} else if hasBlockedWords {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "留言包含不适当内容，已被拒绝",
			"status":  status,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "留言已提交，等待审核后显示",
			"status":  status,
		})
	}
}
