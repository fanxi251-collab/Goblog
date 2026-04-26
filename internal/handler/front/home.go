package front

import (
	"Goblog/internal/service"

	"net/http"

	"github.com/gin-gonic/gin"
)

// HomeHandler 首页处理器
type HomeHandler struct {
	postService   *service.PostService
	columnService *service.ColumnService
	devlogService *service.DevlogService
}

// NewHomeHandler 创建首页处理器
func NewHomeHandler(postSvc *service.PostService, colSvc *service.ColumnService, devlogSvc *service.DevlogService) *HomeHandler {
	return &HomeHandler{
		postService:   postSvc,
		columnService: colSvc,
		devlogService: devlogSvc,
	}
}

// Index 首页
func (h *HomeHandler) Index(c *gin.Context) {
	columns, _ := h.columnService.GetAll()

	// 获取统计数据
	postsTotal, likesTotal, commentsTotal, _ := h.postService.GetStats()

	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":         "灵序之夏",
		"columns":       columns,
		"postsTotal":    postsTotal,
		"likesTotal":    likesTotal,
		"commentsTotal": commentsTotal,
	})
}

// Search 搜索文章API
func (h *HomeHandler) Search(c *gin.Context) {
	keyword := c.Query("q")
	if keyword == "" {
		c.JSON(http.StatusOK, gin.H{"success": true, "data": []interface{}{}})
		return
	}

	// 搜索已发布文章，返回标题、slug、创建时间、专栏slug
	posts, total, _ := h.postService.Search(keyword, "published", 1, 10)

	var results []interface{}
	for _, post := range posts {
		// 获取专栏slug
		columnSlug := ""
		if post.ColumnID > 0 {
			if column, err := h.columnService.GetByID(post.ColumnID); err == nil && column != nil {
				columnSlug = column.Slug
			}
		}
		results = append(results, gin.H{
			"id":          post.ID,
			"title":       post.Title,
			"slug":        post.Slug,
			"column_slug": columnSlug,
			"created_at":  post.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    results,
		"total":   total,
	})
}
