package handler

import (
	"github.com/FruitsAI/Orange/internal/pkg/response"
	"github.com/FruitsAI/Orange/internal/service"
	"github.com/gin-gonic/gin"
)

// DashboardHandler 仪表盘处理器
type DashboardHandler struct {
	dashboardService *service.DashboardService
}

// NewDashboardHandler 创建仪表盘处理器
func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{
		dashboardService: service.NewDashboardService(),
	}
}

// Stats 获取统计数据
// GET /api/v1/dashboard/stats
func (h *DashboardHandler) Stats(c *gin.Context) {
	userID := c.GetInt64("user_id")

	period := c.Query("period") // 不传则返回全局统计
	stats, err := h.dashboardService.GetStats(userID, period)
	if err != nil {
		response.InternalError(c, "获取统计数据失败")
		return
	}

	response.Success(c, stats)
}

// IncomeTrend 获取收入趋势
// GET /api/v1/dashboard/income-trend
func (h *DashboardHandler) IncomeTrend(c *gin.Context) {
	userID := c.GetInt64("user_id")
	period := c.DefaultQuery("period", "month")

	trend, err := h.dashboardService.GetIncomeTrend(userID, period)
	if err != nil {
		response.InternalError(c, "获取收入趋势失败")
		return
	}

	response.Success(c, trend)
}

// RecentProjects 获取最近项目
// GET /api/v1/dashboard/recent-projects
func (h *DashboardHandler) RecentProjects(c *gin.Context) {
	userID := c.GetInt64("user_id")

	projects, err := h.dashboardService.GetRecentProjects(userID)
	if err != nil {
		response.InternalError(c, "获取最近项目失败")
		return
	}

	response.Success(c, projects)
}

// UpcomingPayments 获取即将到期收款
// GET /api/v1/dashboard/upcoming-payments
func (h *DashboardHandler) UpcomingPayments(c *gin.Context) {
	userID := c.GetInt64("user_id")

	payments, err := h.dashboardService.GetUpcomingPayments(userID)
	if err != nil {
		response.InternalError(c, "获取即将到期收款失败")
		return
	}

	response.Success(c, payments)
}
