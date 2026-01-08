import api, { type ApiResponse } from './index'
import type { Project, Payment } from './project'

// 统计数据
export interface DashboardStats {
  total_amount: number
  paid_amount: number
  pending_amount: number
  overdue_amount: number
  total_trend: number
  paid_trend: number
  pending_trend: number
  overdue_trend: number
  avg_collection_days: number
  avg_collection_days_trend: number
}

// 收入趋势
// 收入趋势
export interface IncomeTrend {
  labels: string[]
  actual_values: number[]
  expected_values: number[]
}

// 仪表盘 API
export const dashboardApi = {
  // 获取统计数据
  getStats: (period?: 'week' | 'month' | 'quarter' | 'year') =>
    api.get<ApiResponse<DashboardStats>>('/dashboard/stats', { params: { period, _t: Date.now() } }),

  // 获取收入趋势
  getIncomeTrend: (period?: 'week' | 'month' | 'quarter' | 'year') =>
    api.get<ApiResponse<IncomeTrend>>('/dashboard/income-trend', { params: { period, _t: Date.now() } }),

  // 获取近期项目
  getRecentProjects: () =>
    api.get<ApiResponse<Project[]>>('/dashboard/recent-projects', { params: { _t: Date.now() } }),

  // 获取即将到期收款
  getUpcomingPayments: () =>
    api.get<ApiResponse<Payment[]>>('/dashboard/upcoming-payments', { params: { _t: Date.now() } }),
}
