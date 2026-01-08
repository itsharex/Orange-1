package handler

import (
	"strconv"

	"github.com/FruitsAI/Orange/internal/dto"
	"github.com/FruitsAI/Orange/internal/pkg/response"
	"github.com/FruitsAI/Orange/internal/service"
	"github.com/gin-gonic/gin"
)

// PaymentHandler 收款处理器
type PaymentHandler struct {
	paymentService *service.PaymentService
}

// NewPaymentHandler 创建收款处理器
func NewPaymentHandler() *PaymentHandler {
	return &PaymentHandler{
		paymentService: service.NewPaymentService(),
	}
}

// GetByProject 获取项目的收款列表
// GET /api/v1/projects/:id/payments
func (h *PaymentHandler) GetByProject(c *gin.Context) {
	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ParamError(c, "无效的项目ID")
		return
	}

	payments, err := h.paymentService.ListByProject(projectID)
	if err != nil {
		response.InternalError(c, "获取收款列表失败")
		return
	}

	response.Success(c, payments)
}

// List 获取收款列表
// GET /api/v1/payments
func (h *PaymentHandler) List(c *gin.Context) {
	userID := c.GetInt64("user_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate != "" && endDate != "" {
		payments, err := h.paymentService.ListByDateRange(userID, startDate, endDate)
		if err != nil {
			response.InternalError(c, "获取收款列表失败")
			return
		}
		response.Success(c, payments)
		return
	}

	projectIDStr := c.Query("project_id")
	if projectIDStr != "" {
		projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
		if err != nil {
			response.ParamError(c, "无效的项目ID")
			return
		}
		payments, err := h.paymentService.ListByProject(projectID)
		if err != nil {
			response.InternalError(c, "获取收款列表失败")
			return
		}
		response.Success(c, payments)
		return
	}

	// 如果没有指定项目或日期，返回空列表
	response.Success(c, []interface{}{})
}

// Create 创建收款
// POST /api/v1/payments
func (h *PaymentHandler) Create(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var req dto.PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	req.UserID = userID // 手动设置 UserID

	payment, err := h.paymentService.Create(req)
	if err != nil {
		response.InternalError(c, "创建收款失败")
		return
	}

	response.Success(c, payment)
}

// Update 更新收款
// PUT /api/v1/payments/:id
func (h *PaymentHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ParamError(c, "无效的收款ID")
		return
	}

	var req dto.PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	payment, err := h.paymentService.Update(id, req)
	if err != nil {
		response.InternalError(c, "更新收款失败")
		return
	}

	response.Success(c, payment)
}

// Delete 删除收款
// DELETE /api/v1/payments/:id
func (h *PaymentHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ParamError(c, "无效的收款ID")
		return
	}

	if err := h.paymentService.Delete(id); err != nil {
		response.InternalError(c, "删除收款失败")
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

// Confirm 确认收款
// POST /api/v1/payments/:id/confirm
func (h *PaymentHandler) Confirm(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ParamError(c, "无效的收款ID")
		return
	}

	var req dto.ConfirmPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	if err := h.paymentService.Confirm(id, req.ActualDate, req.Method); err != nil {
		response.InternalError(c, "确认收款失败")
		return
	}

	response.SuccessWithMessage(c, "确认成功", nil)
}
