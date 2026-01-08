package handler

import (
	"strconv"

	"github.com/FruitsAI/Orange/internal/models"
	"github.com/FruitsAI/Orange/internal/pkg/response"
	"github.com/FruitsAI/Orange/internal/service"
	"github.com/gin-gonic/gin"
)

// NotificationHandler 通知处理器
type NotificationHandler struct {
	notificationService *service.NotificationService
}

// NewNotificationHandler 创建通知处理器
func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{
		notificationService: service.NewNotificationService(),
	}
}

// CreateRequest 创建通知请求
type CreateNotificationRequest struct {
	Title        string `json:"title" binding:"required"`
	Content      string `json:"content" binding:"required"`
	Type         string `json:"type"`
	TargetUserID int64  `json:"target_user_id"` // 0 = 全员通知
}

// Create 创建通知（仅管理员）
// POST /api/v1/notifications
func (h *NotificationHandler) Create(c *gin.Context) {
	userID := c.GetInt64("user_id")
	role := c.GetString("role")

	// 检查权限
	if role != "admin" {
		response.Forbidden(c)
		return
	}

	var req CreateNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	notification, err := h.notificationService.Create(userID, req.Title, req.Content, req.Type, req.TargetUserID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, notification)
}

// Update 更新通知（仅管理员）
// PUT /api/v1/notifications/:id
func (h *NotificationHandler) Update(c *gin.Context) {
	role := c.GetString("role")

	// 检查权限
	if role != "admin" {
		response.Forbidden(c)
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ParamError(c, "无效的通知ID")
		return
	}

	var req CreateNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	notification, err := h.notificationService.Update(id, req.Title, req.Content, req.Type, req.TargetUserID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, notification)
}

// List 获取通知列表
// GET /api/v1/notifications?page=1&page_size=10
func (h *NotificationHandler) List(c *gin.Context) {
	userID := c.GetInt64("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var notifications []models.Notification
	var total int64
	var err error

	// 统一使用 ListByUser，确保能获取到 is_read 状态
	// 这里不再区分 Admin，让 Admin 也拥有正常的"收件箱"体验
	notifications, total, err = h.notificationService.ListByUser(userID, page, pageSize)

	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessPage(c, notifications, total, page, pageSize)
}

// MarkAsRead 标记为已读
// PUT /api/v1/notifications/:id/read
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userID := c.GetInt64("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ParamError(c, "无效的通知ID")
		return
	}

	if err := h.notificationService.MarkAsRead(id, userID); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "已标记为已读", nil)
}

// Delete 删除通知（仅管理员）
// DELETE /api/v1/notifications/:id
func (h *NotificationHandler) Delete(c *gin.Context) {
	role := c.GetString("role")

	// 检查权限
	if role != "admin" {
		response.Forbidden(c)
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ParamError(c, "无效的通知ID")
		return
	}

	if err := h.notificationService.Delete(id); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

// UnreadCount 获取未读通知数量
// GET /api/v1/notifications/unread-count
func (h *NotificationHandler) UnreadCount(c *gin.Context) {
	userID := c.GetInt64("user_id")

	count, err := h.notificationService.GetUnreadCount(userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"count": count})
}

// ListUsers 获取用户列表（用于选择发送目标）
// GET /api/v1/notifications/users
func (h *NotificationHandler) ListUsers(c *gin.Context) {
	role := c.GetString("role")

	// 检查权限
	if role != "admin" {
		response.Forbidden(c)
		return
	}

	users, err := h.notificationService.ListUsers()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, users)
}

// Get 获取单个通知
// GET /api/v1/notifications/:id
func (h *NotificationHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ParamError(c, "无效的通知ID")
		return
	}

	notification, err := h.notificationService.Get(id)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, notification)
}
