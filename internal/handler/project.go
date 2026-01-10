package handler

import (
	"strconv"

	"github.com/FruitsAI/Orange/internal/dto"
	"github.com/FruitsAI/Orange/internal/pkg/response"
	"github.com/FruitsAI/Orange/internal/service"
	"github.com/gin-gonic/gin"
)

// ProjectHandler 项目处理器
type ProjectHandler struct {
	projectService *service.ProjectService
}

// NewProjectHandler 创建项目处理器
func NewProjectHandler() *ProjectHandler {
	return &ProjectHandler{
		projectService: service.NewProjectService(),
	}
}

// List 获取项目列表
// GET /api/v1/projects
func (h *ProjectHandler) List(c *gin.Context) {
	userID := c.GetInt64("user_id")
	status := c.Query("status")
	keyword := c.Query("keyword")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	result, err := h.projectService.List(userID, status, keyword, page, pageSize)
	if err != nil {
		response.InternalError(c, "获取项目列表失败")
		return
	}

	response.SuccessPage(c, result.List, result.Total, result.Page, result.PageSize)
}

// Get 获取项目详情
// GET /api/v1/projects/:id
func (h *ProjectHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ParamError(c, "无效的项目ID")
		return
	}

	project, err := h.projectService.Get(id)
	if err != nil {
		response.NotFound(c, "项目不存在")
		return
	}

	response.Success(c, project)
}

// Create 创建项目
// POST /api/v1/projects
func (h *ProjectHandler) Create(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var req dto.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	req.UserID = userID // 手动设置 UserID

	project, err := h.projectService.Create(req)
	if err != nil {
		response.InternalError(c, "创建项目失败")
		return
	}

	response.Success(c, project)
}

// Update 更新项目
// PUT /api/v1/projects/:id
func (h *ProjectHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ParamError(c, "无效的项目ID")
		return
	}

	var req dto.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	project, err := h.projectService.Update(id, req)
	if err != nil {
		response.InternalError(c, "更新项目失败")
		return
	}

	response.Success(c, project)
}

// Delete 删除项目
// DELETE /api/v1/projects/:id
func (h *ProjectHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ParamError(c, "无效的项目ID")
		return
	}

	if err := h.projectService.Delete(id); err != nil {
		response.InternalError(c, "删除项目失败")
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

// Archive 归档项目
// POST /api/v1/projects/:id/archive
func (h *ProjectHandler) Archive(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ParamError(c, "无效的项目ID")
		return
	}

	if err := h.projectService.Archive(id); err != nil {
		response.InternalError(c, "归档项目失败")
		return
	}

	response.SuccessWithMessage(c, "归档成功", nil)
}

// CheckContractNumber 检查合同编号是否已存在
// GET /api/v1/projects/check-contract-number?contract_number=xxx&exclude_id=xxx
func (h *ProjectHandler) CheckContractNumber(c *gin.Context) {
	userID := c.GetInt64("user_id")
	contractNumber := c.Query("contract_number")
	if contractNumber == "" {
		response.ParamError(c, "合同编号不能为空")
		return
	}

	excludeID, _ := strconv.ParseInt(c.DefaultQuery("exclude_id", "0"), 10, 64)

	exists, err := h.projectService.CheckContractNumberExists(userID, contractNumber, excludeID)
	if err != nil {
		response.InternalError(c, "检查合同编号失败")
		return
	}

	response.Success(c, gin.H{"exists": exists})
}

// GenerateContractNumber 生成下一个合同编号
// GET /api/v1/projects/generate-contract-number?date=2026-01-10
func (h *ProjectHandler) GenerateContractNumber(c *gin.Context) {
	userID := c.GetInt64("user_id")
	date := c.Query("date")
	if date == "" {
		response.ParamError(c, "日期不能为空")
		return
	}

	contractNumber, err := h.projectService.GenerateNextContractNumber(userID, date)
	if err != nil {
		response.InternalError(c, "生成合同编号失败")
		return
	}

	response.Success(c, gin.H{"contract_number": contractNumber})
}
