package front

import (
	"Goblog/internal/model"
	"Goblog/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"html/template"
)

// ColumnHandler 专栏处理器
type ColumnHandler struct {
	postService    *service.PostService
	columnService  *service.ColumnService
	commentService *service.CommentService
}

// NewColumnHandler 创建专栏处理器
func NewColumnHandler(postSvc *service.PostService, colSvc *service.ColumnService, commentSvc *service.CommentService) *ColumnHandler {
	return &ColumnHandler{
		postService:    postSvc,
		columnService:  colSvc,
		commentService: commentSvc,
	}
}

// Index 专栏首页 - 显示所有专栏 + 第一个专栏的文章
func (h *ColumnHandler) Index(c *gin.Context) {
	columns, _ := h.columnService.GetAll()

	var currentColumn interface{}
	var posts interface{}
	pageSize := 10
	highlightID := c.Query("highlight")

	// 获取统计数据
	columnsTotal := int64(len(columns))
	postsTotal, _, _, _ := h.postService.GetStats()

	// 获取最后更新时间
	lastUpdateTime, err := h.postService.GetLatestUpdateTime()
	lastUpdate := "暂无更新"
	if err == nil && lastUpdateTime > 0 {
		lastUpdate = time.Unix(lastUpdateTime, 0).Format("2006-01-02")
	}

	// 如果有专栏，获取第一个专栏的文章
	if len(columns) > 0 {
		currentColumn = columns[0]
		postsData, total, _ := h.postService.GetByColumn(columns[0].ID, "published", 1, pageSize)
		posts = postsData

		totalPages := int(total) / pageSize
		if int(total)%pageSize > 0 {
			totalPages++
		}
		c.HTML(http.StatusOK, "column.html", gin.H{
			"title":         "专栏 - 灵序之夏",
			"columns":       columns,
			"currentColumn": currentColumn,
			"posts":         posts,
			"page":          1,
			"totalPages":    totalPages,
			"columnsTotal":  columnsTotal,
			"postsTotal":    postsTotal,
			"lastUpdate":    lastUpdate,
			"highlightID":   highlightID,
		})
		return
	}

	totalPages := 0
	c.HTML(http.StatusOK, "column.html", gin.H{
		"title":         "专栏 - 灵序之夏",
		"columns":       columns,
		"currentColumn": currentColumn,
		"posts":         posts,
		"page":          1,
		"totalPages":    totalPages,
		"columnsTotal":  columnsTotal,
		"postsTotal":    postsTotal,
		"lastUpdate":    lastUpdate,
		"highlightID":   highlightID,
	})
}

// List 专栏文章列表 - 显示所有专栏 + 指定专栏的文章
func (h *ColumnHandler) List(c *gin.Context) {
	slug := c.Param("slug")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize := 10
	keyword := c.Query("keyword")
	highlightID := c.Query("highlight") // 高亮文章ID

	// 获取所有专栏
	columns, _ := h.columnService.GetAll()

	// 获取统计数据
	columnsTotal := int64(len(columns))
	postsTotal, _, _, _ := h.postService.GetStats()

	// 获取最后更新时间
	lastUpdateTime, err := h.postService.GetLatestUpdateTime()
	lastUpdate := "暂无更新"
	if err == nil && lastUpdateTime > 0 {
		lastUpdate = time.Unix(lastUpdateTime, 0).Format("2006-01-02")
	}

	// 根据slug获取当前专栏
	var currentColumn interface{}
	var posts interface{}
	var totalPages int

	if keyword != "" {
		// 搜索模式 - 在当前专栏内搜索
		var columnModel *model.Column
		if slug != "" {
			columnModel, _ = h.columnService.GetBySlug(slug)
		}
		var postsData []model.Post
		var total int64
		if columnModel != nil {
			currentColumn = columnModel
			postsData, total, _ = h.postService.SearchInColumn(columnModel.ID, keyword, "published", page, pageSize)
		} else {
			postsData, total, _ = h.postService.Search(keyword, "published", page, pageSize)
		}
		posts = postsData

		totalPages = int(total) / pageSize
		if int(total)%pageSize > 0 {
			totalPages++
		}
		c.HTML(http.StatusOK, "column.html", gin.H{
			"title":         "搜索: " + keyword + " - 灵序之夏",
			"columns":       columns,
			"currentColumn": currentColumn,
			"posts":         posts,
			"page":          page,
			"totalPages":    totalPages,
			"keyword":       keyword,
			"columnsTotal":  columnsTotal,
			"postsTotal":    postsTotal,
			"lastUpdate":    lastUpdate,
			"highlightID":   highlightID,
		})
		return
	}

	if slug != "" {
		column, err := h.columnService.GetBySlug(slug)
		if err == nil {
			currentColumn = column
			postsData, total, _ := h.postService.GetByColumn(column.ID, "published", page, pageSize)
			posts = postsData

			totalPages = int(total) / pageSize
			if int(total)%pageSize > 0 {
				totalPages++
			}
			c.HTML(http.StatusOK, "column.html", gin.H{
				"title":         column.Name + " - 灵序之夏",
				"columns":       columns,
				"currentColumn": currentColumn,
				"posts":         posts,
				"page":          page,
				"totalPages":    totalPages,
				"columnsTotal":  columnsTotal,
				"postsTotal":    postsTotal,
				"lastUpdate":    lastUpdate,
				"highlightID":   highlightID,
			})
			return
		}
	}

	// 如果没有指定slug或slug无效，显示第一个专栏
	if len(columns) > 0 {
		currentColumn = columns[0]
		postsData, total, _ := h.postService.GetByColumn(columns[0].ID, "published", 1, pageSize)
		posts = postsData

		totalPages = int(total) / pageSize
		if int(total)%pageSize > 0 {
			totalPages++
		}
		c.HTML(http.StatusOK, "column.html", gin.H{
			"title":         "专栏 - 灵序之夏",
			"columns":       columns,
			"currentColumn": currentColumn,
			"posts":         posts,
			"page":          1,
			"totalPages":    totalPages,
			"columnsTotal":  columnsTotal,
			"postsTotal":    postsTotal,
			"lastUpdate":    lastUpdate,
		})
	}
}

// Post 文章详情
func (h *ColumnHandler) Post(c *gin.Context) {
	slug := c.Param("slug")

	post, err := h.postService.GetBySlug(slug)
	if err != nil {
		c.String(404, "文章不存在")
		return
	}

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

	// 获取所有专栏
	columns, _ := h.columnService.GetAll()

	// 获取当前专栏的所有文章（不排除当前文章）
	var columnPosts []model.Post
	if post.ColumnID > 0 {
		var err error
		columnPosts, _, err = h.postService.GetByColumn(post.ColumnID, "published", 1, 1000)
		if err != nil {
			columnPosts = []model.Post{}
		}
	}

	// 渲染 Markdown
	contentHTML := ""
	if post.Content != "" {
		contentHTML = service.RenderMarkdown(post.Content)
		if contentHTML == post.Content || contentHTML == "" {
			// 恢复失败，显示原始文本
			contentHTML = "<pre style='white-space:pre-wrap;font-family:inherit;font-size:inherit;line-height:1.8;'>" + post.Content + "</pre>"
		}
	}

	// 获取文章评论（只显示已审核的）
	var comments []model.Comment
	if h.commentService != nil {
		commentsData, _, _ := h.commentService.GetApproved(post.ID, 1, 100)
		comments = commentsData
	}

	c.HTML(http.StatusOK, "post.html", gin.H{
		"title":         post.Title + " - 灵序之夏",
		"post":          post,
		"contentHTML":   template.HTML(contentHTML),
		"column":        column,
		"columns":       columns,
		"columnPosts":   columnPosts,
		"postsTotal":    postsTotal,
		"likesTotal":    likesTotal,
		"commentsTotal": commentsTotal,
		"comments":      comments,
	})
}
