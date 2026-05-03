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

	// 处理路径
	path = strings.ReplaceAll(path, "\\", "/")
	imageURL := "/static/uploads/" + path

	// 返回 JSON 格式: { code: 0, data: { url: "...", name: "..." } }
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg": "上传成功",
		"data": gin.H{
			"url":  imageURL,
			"name": file.Filename,
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

	// 获取 postID（可选，新建文章时为0）
	postIDStr := c.PostForm("post_id")
	postID := uint(0)
	if postIDStr != "" {
		id, err := strconv.ParseUint(postIDStr, 10, 32)
		if err == nil {
			postID = uint(id)
		}
	}

	path, err := h.fileService.UploadCover(file, postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "封面上传成功",
		"data": gin.H{
			"url":  "/static/" + path,
			"path": path,
		},
	})
}
