package admin

import (
	"Goblog/internal/model"
	"Goblog/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ColumnHandler 专栏处理器
type ColumnHandler struct {
	columnService *service.ColumnService
}

// NewColumnHandler 创建专栏处理器
func NewColumnHandler(colSvc *service.ColumnService) *ColumnHandler {
	return &ColumnHandler{columnService: colSvc}
}

// List 专栏列表
func (h *ColumnHandler) List(c *gin.Context) {
	columns, err := h.columnService.GetAll()
	if err != nil {
		columns = []model.Column{}
	}

	c.HTML(http.StatusOK, "column_list.html", gin.H{
		"title":   "专栏管理",
		"columns": columns,
	})
}

// Edit 编辑页面
func (h *ColumnHandler) Edit(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		id = 0
	}

	var column *model.Column
	if id > 0 {
		column, _ = h.columnService.GetByID(uint(id))
	}
	if column == nil {
		column = &model.Column{}
	}

	parents, _ := h.columnService.GetParents()

	c.HTML(http.StatusOK, "column_edit.html", gin.H{
		"title":   "编辑专栏",
		"column":  column,
		"parents": parents,
	})
}

// Save 保存专栏
func (h *ColumnHandler) Save(c *gin.Context) {
	id, _ := strconv.ParseUint(c.PostForm("id"), 10, 32)
	name := c.PostForm("name")
	slug := c.PostForm("slug")
	description := c.PostForm("description")
	parentID, _ := strconv.ParseUint(c.PostForm("parent_id"), 10, 32)
	sort, _ := strconv.Atoi(c.PostForm("sort"))

	column := &model.Column{
		Name:        name,
		Slug:        slug,
		Description: description,
		ParentID:    uint(parentID),
		Sort:        sort,
	}

	var err error
	if id > 0 {
		column.ID = uint(id)
		err = h.columnService.Update(column)
	} else {
		err = h.columnService.Create(column)
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "id": column.ID})
}

// Delete 删除专栏
func (h *ColumnHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	err = h.columnService.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
