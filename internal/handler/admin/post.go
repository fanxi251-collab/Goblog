package admin

import (
	"Goblog/internal/model"
	"Goblog/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// PostHandler 文章处理器
type PostHandler struct {
	postService   *service.PostService
	columnService *service.ColumnService
}

// NewPostHandler 创建文章处理器
func NewPostHandler(postSvc *service.PostService, colSvc *service.ColumnService) *PostHandler {
	return &PostHandler{
		postService:   postSvc,
		columnService: colSvc,
	}
}

// List 文章列表（已发布）
func (h *PostHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize := 20
	keyword := c.Query("keyword")

	var posts []model.Post
	var total int64
	var columns, _ = h.columnService.GetAll()

	if keyword != "" {
		// 搜索模式
		posts, total, _ = h.postService.Search(keyword, "", page, pageSize)
	} else {
		// 正常列表
		posts, total, _ = h.postService.GetByStatus("published", page, pageSize)
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	c.HTML(http.StatusOK, "post_list.html", gin.H{
		"title":      "文章管理",
		"posts":      posts,
		"columns":    columns,
		"page":       page,
		"totalPages": totalPages,
		"total":      total,
		"keyword":    keyword,
	})
}

// Drafts 草稿箱
func (h *PostHandler) Drafts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize := 20
	keyword := c.Query("keyword")

	var posts []model.Post
	var total int64
	var columns, _ = h.columnService.GetAll()

	if keyword != "" {
		// 搜索模式
		posts, total, _ = h.postService.Search(keyword, "", page, pageSize)
	} else {
		// 草稿列表
		posts, total, _ = h.postService.GetByStatus("draft", page, pageSize)
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	c.HTML(http.StatusOK, "draft_list.html", gin.H{
		"title":      "草稿箱",
		"posts":      posts,
		"columns":    columns,
		"page":       page,
		"totalPages": totalPages,
		"total":      total,
		"keyword":    keyword,
	})
}

// Edit 编辑页面
func (h *PostHandler) Edit(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		id = 0
	}

	var post *model.Post
	if id > 0 {
		post, _ = h.postService.GetByID(uint(id))
	}
	if post == nil {
		post = &model.Post{
			Status:   "draft",
			ColumnID: 0,
		}
	}

	columns, _ := h.columnService.GetAll()

	c.HTML(http.StatusOK, "post_edit.html", gin.H{
		"title":   "编辑文章",
		"post":    post,
		"columns": columns,
	})
}

// Save 保存文章
func (h *PostHandler) Save(c *gin.Context) {
	id, _ := strconv.ParseUint(c.PostForm("id"), 10, 32)
	title := c.PostForm("title")
	slug := c.PostForm("slug")
	content := c.PostForm("content")
	excerpt := c.PostForm("excerpt")
	columnID, _ := strconv.ParseUint(c.PostForm("column_id"), 10, 32)
	coverImage := c.PostForm("cover_image")
	action := c.PostForm("action")

	// 验证专栏必选
	if columnID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择专栏"})
		return
	}

	// 根据action设置状态
	status := "draft"
	if action == "publish" {
		status = "published"
	}

	// 处理 slug
	if id > 0 {
		// 更新时：如果 slug 为空，获取原文章的 slug
		if slug == "" {
			oldPost, _ := h.postService.GetByID(uint(id))
			if oldPost != nil {
				slug = oldPost.Slug
			}
		}
	} else {
		// 新建时：如果 slug 为空，自动生成
		if slug == "" {
			slug = generateSlug(title)
		}
		// 检查 slug 是否已存在，存在则添加后缀
		for {
			exists, _ := h.postService.GetBySlug(slug)
			if exists == nil {
				break
			}
			slug = slug + "-" + randomString(4)
		}
	}

	post := &model.Post{
		Title:      title,
		Slug:       slug,
		Content:    content,
		Excerpt:    excerpt,
		CoverImage: coverImage,
		ColumnID:   uint(columnID),
		Status:     status,
		CreatedAt:  time.Now().Unix(),
	}

	var err error
	if id > 0 {
		post.ID = uint(id)
		// 获取原有文章数据，保留 CreatedAt
		oldPost, _ := h.postService.GetByID(uint(id))
		if oldPost != nil {
			post.CreatedAt = oldPost.CreatedAt
			post.UpdatedAt = oldPost.UpdatedAt
		}
		err = h.postService.Update(post)
	} else {
		err = h.postService.Create(post)
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "id": post.ID})
}

// Delete 删除文章
func (h *PostHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	err = h.postService.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// Publish 发布文章
func (h *PostHandler) Publish(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	err = h.postService.Publish(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// generateSlug 根据标题自动生成slug
func generateSlug(title string) string {
	if title == "" {
		return ""
	}

	// 简单转换：转小写，空格和特殊字符替换为短横线
	var result []rune
	for _, r := range title {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			result = append(result, r)
		} else if r == ' ' || r == '-' || r == '_' {
			if len(result) > 0 && result[len(result)-1] != '-' {
				result = append(result, '-')
			}
		}
	}

	slug := string(result)
	// 移除开头和结尾的短横线
	for len(slug) > 0 && slug[0] == '-' {
		slug = slug[1:]
	}
	for len(slug) > 0 && slug[len(slug)-1] == '-' {
		slug = slug[:len(slug)-1]
	}

	// 如果为空，生成随机字符串
	if slug == "" {
		return randomString(8)
	}

	return slug
}

// randomString 生成随机字符串
func randomString(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]rune, length)
	for i := range b {
		b[i] = rune(letters[i%len(letters)])
	}
	return string(b)
}
