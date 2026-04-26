package admin

import (
	"Goblog/internal/model"
	"Goblog/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// DevlogHandler 开发日志处理器
type DevlogHandler struct {
	devlogService *service.DevlogService
}

// NewDevlogHandler 创建开发日志处理器
func NewDevlogHandler(devlogSvc *service.DevlogService) *DevlogHandler {
	return &DevlogHandler{
		devlogService: devlogSvc,
	}
}

// List 日志列表（显示所有：草稿+已发布）
func (h *DevlogHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize := 20

	devlogs, total, _ := h.devlogService.GetAll(page, pageSize)

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	c.HTML(http.StatusOK, "devlogs.html", gin.H{
		"title":      "开发日志 - 灵序之夏",
		"devlogs":    devlogs,
		"page":       page,
		"totalPages": totalPages,
		"total":      total,
	})
}

// Edit 日志编辑（新建/编辑）
func (h *DevlogHandler) Edit(c *gin.Context) {
	idStr := c.Param("id")
	var devlog *model.Devlog
	var isNew bool

	if idStr == "" {
		// 新建
		isNew = true
		devlog = &model.Devlog{
			Date:   time.Now().Unix(),
			Status: "draft",
		}
	} else {
		// 编辑
		id, _ := strconv.ParseUint(idStr, 10, 32)
		var err error
		devlog, err = h.devlogService.GetByID(uint(id))
		if err != nil {
			c.String(http.StatusNotFound, "日志不存在")
			return
		}
		isNew = false
	}

	c.HTML(http.StatusOK, "devlog_edit.html", gin.H{
		"title":  "编辑日志 - 灵序之夏",
		"devlog": devlog,
		"isNew":  isNew,
	})
}

// Save 保存日志
func (h *DevlogHandler) Save(c *gin.Context) {
	idStr := c.PostForm("id")
	title := c.PostForm("title")
	description := c.PostForm("description")
	dateStr := c.PostForm("date")
	action := c.PostForm("action")

	// 解析日期
	var dateUnix int64
	if dateStr != "" {
		t, err := time.Parse("2006-01-02", dateStr)
		if err == nil {
			dateUnix = t.Unix()
		}
	}

	// 根据action设置状态
	status := "draft"
	if action == "publish" {
		status = "published"
	}

	// 判断是新建还是更新：id 为空或为 "0" 时是新建
	isNew := idStr == "" || idStr == "0"

	if !isNew {
		// 更新
		id, _ := strconv.ParseUint(idStr, 10, 32)
		devlog, err := h.devlogService.GetByID(uint(id))
		if err != nil {
			c.String(http.StatusBadRequest, "日志不存在")
			return
		}
		devlog.Title = title
		devlog.Description = description
		devlog.Date = dateUnix
		devlog.Status = status

		if err := h.devlogService.Update(devlog); err != nil {
			c.String(http.StatusInternalServerError, "保存失败")
			return
		}
	} else {
		// 新建
		devlog := &model.Devlog{
			Title:       title,
			Description: description,
			Date:        dateUnix,
			Status:      status,
		}
		if err := h.devlogService.Create(devlog); err != nil {
			c.String(http.StatusInternalServerError, "创建失败")
			return
		}
	}

	c.Redirect(http.StatusFound, "/admin/devlogs")
}

// Delete 删除日志
func (h *DevlogHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效ID"})
		return
	}

	if err := h.devlogService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// Publish 发布日志
func (h *DevlogHandler) Publish(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效ID"})
		return
	}

	if err := h.devlogService.Publish(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发布失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// Unpublish 下架日志
func (h *DevlogHandler) Unpublish(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效ID"})
		return
	}

	if err := h.devlogService.Unpublish(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "下架失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
