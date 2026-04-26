package admin

import (
	"Goblog/internal/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// UploadHandler 上传处理器
type UploadHandler struct {
	fileService service.FileService
}

// NewUploadHandler 创建上传处理器
func NewUploadHandler(fileSvc service.FileService) *UploadHandler {
	return &UploadHandler{fileService: fileSvc}
}

// Upload 上传文件
func (h *UploadHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择文件"})
		return
	}

	path, err := h.fileService.Upload(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回 HTML img 标签（更可靠）
	path = strings.ReplaceAll(path, "\\", "/")
	imageURL := "<img src=\"/static/uploads/" + path + "\" alt=\"\">"

	c.JSON(http.StatusOK, gin.H{
		"success": 1,
		"message": "上传成功",
		"data": gin.H{
			"url":  imageURL,
			"path": path,
		},
	})
}

// UploadCover 上传封面图片
func (h *UploadHandler) UploadCover(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择图片"})
		return
	}

	// 获取 postID
	postIDStr := c.PostForm("post_id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil || postID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	path, err := h.fileService.UploadCover(file, uint(postID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": 1,
		"message": "封面上传成功",
		"data": gin.H{
			"url":  "/static/" + path,
			"path": path,
		},
	})
}
