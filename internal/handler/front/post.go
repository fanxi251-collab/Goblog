package front

import (
	"Goblog/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PostHandler 文章处理器（前台）
type PostHandler struct {
	postService     *service.PostService
	commentService  *service.CommentService
	visitorService  *service.VisitorService
	postLikeService *service.PostLikeService
}

// NewPostHandler 创建文章处理器（前台）
func NewPostHandler(
	postSvc *service.PostService,
	commentSvc *service.CommentService,
	visitorSvc *service.VisitorService,
	postLikeSvc *service.PostLikeService,
) *PostHandler {
	return &PostHandler{
		postService:     postSvc,
		commentService:  commentSvc,
		visitorService:  visitorSvc,
		postLikeService: postLikeSvc,
	}
}

// Like 点赞/取消点赞
func (h *PostHandler) Like(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	// 获取访客Token
	token := c.GetHeader("X-Visitor-Token")
	if token == "" {
		token, _ = c.Cookie("visitor_token")
	}

	ip := c.ClientIP()

	if token != "" {
		// 已登录访客，使用VisitorID点赞
		visitor, _ := h.visitorService.CheckToken(token)
		if visitor != nil {
			liked, err := h.postLikeService.Like(uint(postID), visitor.ID, ip)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"liked":   liked,
			})
			return
		}
	}

	// 未登录，使用IP点赞（备用方案）
	liked, err := h.postLikeService.LikeByIP(uint(postID), ip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"liked":   liked,
	})
}

// CheckLike 检查是否已点赞
func (h *PostHandler) CheckLike(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	// 获取访客Token
	token := c.GetHeader("X-Visitor-Token")
	if token == "" {
		token, _ = c.Cookie("visitor_token")
	}

	if token != "" {
		visitor, _ := h.visitorService.CheckToken(token)
		if visitor != nil {
			hasLiked, _ := h.postLikeService.HasLiked(uint(postID), visitor.ID)
			c.JSON(http.StatusOK, gin.H{"hasLiked": hasLiked})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"hasLiked": false})
}
