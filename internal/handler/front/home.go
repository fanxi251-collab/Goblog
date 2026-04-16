package front

import (
	"Goblog/internal/model"
	"Goblog/internal/service"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// HomeHandler 首页处理器
type HomeHandler struct {
	postService   *service.PostService
	columnService *service.ColumnService
}

// NewHomeHandler 创建首页处理器
func NewHomeHandler(postSvc *service.PostService, colSvc *service.ColumnService) *HomeHandler {
	return &HomeHandler{
		postService:   postSvc,
		columnService: colSvc,
	}
}

// Index 首页
func (h *HomeHandler) Index(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize := 10
	search := c.Query("search")

	var posts []model.Post
	var total int64
	var err error

	columns, _ := h.columnService.GetAll()

	// 根据是否有搜索参数，决定获取方式
	if search != "" {
		// 搜索模式 - 搜索无专栏的文章
		posts, total, err = h.postService.SearchNoColumn(search, page, pageSize)
	} else {
		// 首页 - 获取无专栏的已发布文章
		posts, total, err = h.postService.GetPublishedNoColumn(page, pageSize)
	}

	if err != nil {
		posts = []model.Post{}
		total = 0
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":      "灵序之夏",
		"posts":      posts,
		"columns":    columns,
		"page":       page,
		"totalPages": totalPages,
		"total":      total,
		"search":     search,
	})
}

// About 关于页面
func (h *HomeHandler) About(c *gin.Context) {
	c.HTML(http.StatusOK, "about.html", gin.H{
		"title": "关于我 - 灵序之夏",
	})
}
