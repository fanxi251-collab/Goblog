package front

import (
	"Goblog/internal/service"
	"net/http"

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

// Index 开发日志首页
func (h *DevlogHandler) Index(c *gin.Context) {
	devlogs, _, _ := h.devlogService.GetPublished(1, 100)
	c.HTML(http.StatusOK, "devlog.html", gin.H{
		"title":   "开发日志 - 灵序之夏",
		"devlogs": devlogs,
	})
}
