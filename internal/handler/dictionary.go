package handler

import (
	"strconv"

	"github.com/FruitsAI/Orange/internal/dto"
	"github.com/FruitsAI/Orange/internal/pkg/response"
	"github.com/FruitsAI/Orange/internal/service"
	"github.com/gin-gonic/gin"
)

// DictionaryHandler 字典处理器
type DictionaryHandler struct {
	dictService *service.DictionaryService
}

// NewDictionaryHandler 创建字典处理器
func NewDictionaryHandler() *DictionaryHandler {
	return &DictionaryHandler{
		dictService: service.NewDictionaryService(),
	}
}

// List 获取字典列表
// GET /api/v1/dictionaries
func (h *DictionaryHandler) List(c *gin.Context) {
	dictionaries, err := h.dictService.List()
	if err != nil {
		response.InternalError(c, "获取字典列表失败")
		return
	}

	response.Success(c, dictionaries)
}

// GetItems 获取字典项
// GET /api/v1/dictionaries/:code/items
func (h *DictionaryHandler) GetItems(c *gin.Context) {
	code := c.Param("code")

	items, err := h.dictService.GetItems(code)
	if err != nil {
		response.NotFound(c, "字典不存在")
		return
	}

	response.Success(c, items)
}

// CreateItem 创建字典项
// POST /api/v1/dictionaries/:code/items
func (h *DictionaryHandler) CreateItem(c *gin.Context) {
	if c.GetString("role") != "admin" {
		response.Forbidden(c, "无权操作")
		return
	}
	code := c.Param("code")

	var req dto.CreateDictionaryItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	item, err := h.dictService.CreateItem(code, req.Label, req.Value, req.Sort)
	if err != nil {
		response.InternalError(c, "创建字典项失败")
		return
	}

	response.Success(c, item)
}

// UpdateItem 更新字典项
// PUT /api/v1/dictionaries/:code/items/:id
func (h *DictionaryHandler) UpdateItem(c *gin.Context) {
	if c.GetString("role") != "admin" {
		response.Forbidden(c, "无权操作")
		return
	}
	// code := c.Param("code") // Not used in service update, but part of URL
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ParamError(c, "无效的字典项ID")
		return
	}

	var req dto.CreateDictionaryItemRequest // Reuse create request wrapper or make new one
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	// Assuming sort is optional in request, handling 0 might be tricky if 0 is valid sort.
	// But dto probably defines it.

	item, err := h.dictService.UpdateItem(id, req.Label, req.Value, req.Sort)
	if err != nil {
		response.InternalError(c, "更新字典项失败")
		return
	}

	response.Success(c, item)
}

// DeleteItem 删除字典项
// DELETE /api/v1/dictionaries/:code/items/:id
func (h *DictionaryHandler) DeleteItem(c *gin.Context) {
	if c.GetString("role") != "admin" {
		response.Forbidden(c, "无权操作")
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ParamError(c, "无效的字典项ID")
		return
	}

	if err := h.dictService.DeleteItem(id); err != nil {
		response.InternalError(c, "删除字典项失败")
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}
