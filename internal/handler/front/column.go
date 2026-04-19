package front

import (
	"Goblog/internal/model"
	"Goblog/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ColumnHandler 专栏处理器
type ColumnHandler struct {
	postService   *service.PostService
	columnService *service.ColumnService
}

// NewColumnHandler 创建专栏处理器
func NewColumnHandler(postSvc *service.PostService, colSvc *service.ColumnService) *ColumnHandler {
	return &ColumnHandler{
		postService:   postSvc,
		columnService: colSvc,
	}
}

// Index 专栏首页
func (h *ColumnHandler) Index(c *gin.Context) {
	columns, _ := h.columnService.GetAll()

	c.HTML(http.StatusOK, "column.html", gin.H{
		"title":   "专栏 - 灵序之夏",
		"columns": columns,
	})
}

// List 专栏文章列表
func (h *ColumnHandler) List(c *gin.Context) {
	slug := c.Param("slug")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize := 10

	// 根据slug获取专栏
	column, err := h.columnService.GetBySlug(slug)
	if err != nil {
		c.String(404, "专栏不存在")
		return
	}

	// 获取专栏下的已发布文章
	posts, total, _ := h.postService.GetByColumn(column.ID, "published", page, pageSize)

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	c.HTML(http.StatusOK, "column_detail.html", gin.H{
		"title":      column.Name + " - 灵序之夏",
		"column":     column,
		"posts":      posts,
		"page":       page,
		"totalPages": totalPages,
	})
}

// Post 文章详情
func (h *ColumnHandler) Post(c *gin.Context) {
	slug := c.Param("slug")

	post, err := h.postService.GetBySlug(slug)
	if err != nil {
		c.String(404, "文章不存在")
		return
	}

	// 增加浏览次数
	h.postService.IncrViewCount(post.ID)

	// 渲染Markdown
	content := service.RenderMarkdown(post.Content)
	post.Content = content

	// 获取专栏信息（只有 ColumnID > 0 时才获取）
	var column *model.Column
	if post.ColumnID > 0 {
		column, _ = h.columnService.GetByID(post.ColumnID)
	}
	// 如果 column 为 nil，创建一个空对象避免模板报错
	if column == nil {
		column = &model.Column{}
	}

	// 获取统计数据
	postsTotal, likesTotal, commentsTotal, _ := h.postService.GetStats()

	c.HTML(http.StatusOK, "post.html", gin.H{
		"title":         post.Title + " - 灵序之夏",
		"post":          post,
		"column":        column,
		"postsTotal":    postsTotal,
		"likesTotal":    likesTotal,
		"commentsTotal": commentsTotal,
	})
}
