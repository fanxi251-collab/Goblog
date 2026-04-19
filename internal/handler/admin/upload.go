package admin

import (
	"Goblog/internal/service"
	"net/http"

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

	// 返回 Vditor 需要的格式：markdown 格式的图片链接
	imageURL := "![](/static/uploads/" + path + ")"

	c.JSON(http.StatusOK, gin.H{
		"success": 1,
		"message": "上传成功",
		"data": gin.H{
			"url":  imageURL,
			"path": path,
		},
	})
}
