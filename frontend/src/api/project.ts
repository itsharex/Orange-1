import api, { type ApiResponse, type PageData } from './index'
import type { User } from './auth'

// 项目类型
export interface Project {
  id: number
  name: string
  company: string
  total_amount: number
  received_amount: number
  status: 'active' | 'completed' | 'pending' | 'notstarted' | 'archived'
  type: string
  contract_number: string
  contract_date: string
  payment_method: string
  start_date: string
  end_date: string
  description: string
  create_time: string
  payments?: Payment[]
  user?: User
}

// 收款类型
export interface Payment {
  id: number
  project_id: number
  stage: string
  amount: number
  percentage: number
  plan_date: string
  status: 'paid' | 'pending' | 'overdue'
  actual_date: string
  method: string
  remark: string
  project?: Project
}

// 项目列表查询参数
export interface ProjectListParams {
  page?: number
  page_size?: number
  status?: string
  keyword?: string
}

// 创建/更新项目请求
export interface ProjectRequest {
  name: string
  company: string
  total_amount: number
  status?: string
  type: string
  contract_number?: string
  contract_date?: string
  payment_method?: string
  start_date: string
  end_date: string
  description?: string
}

// 创建收款请求
export interface PaymentRequest {
  project_id: number
  stage: string
  amount: number
  percentage?: number
  plan_date: string
  status?: string
  method?: string
  remark?: string
}

// 确认收款请求
export interface ConfirmPaymentRequest {
  actual_date: string
  method?: string
}

// 项目 API
export const projectApi = {
  // 获取项目列表
  list: (params?: ProjectListParams) =>
    api.get<ApiResponse<PageData<Project>>>('/projects', { params }),

  // 获取项目详情
  get: (id: number) =>
    api.get<ApiResponse<Project>>(`/projects/${id}`, { params: { _t: Date.now() } }),

  // 创建项目
  create: (data: ProjectRequest) =>
    api.post<ApiResponse<Project>>('/projects', data),

  // 更新项目
  update: (id: number, data: ProjectRequest) =>
    api.put<ApiResponse<Project>>(`/projects/${id}`, data),

  // 删除项目
  delete: (id: number) =>
    api.delete<ApiResponse<null>>(`/projects/${id}`),

  // 归档项目
  archive: (id: number) =>
    api.post<ApiResponse<null>>(`/projects/${id}/archive`),

  // 获取项目收款列表
  getPayments: (projectId: number) =>
    api.get<ApiResponse<Payment[]>>(`/projects/${projectId}/payments`, { params: { _t: Date.now() } }),
}

// 收款 API
export const paymentApi = {
  // 获取收款列表
  list: (params?: { project_id?: number; status?: string; start_date?: string; end_date?: string; _t?: number }) =>
    api.get<ApiResponse<Payment[]>>('/payments', { params }),

  // 创建收款
  create: (data: PaymentRequest) =>
    api.post<ApiResponse<Payment>>('/payments', data),

  // 更新收款
  update: (id: number, data: PaymentRequest) =>
    api.put<ApiResponse<Payment>>(`/payments/${id}`, data),

  // 删除收款
  delete: (id: number) =>
    api.delete<ApiResponse<null>>(`/payments/${id}`),

  // 确认收款
  confirm: (id: number, data: ConfirmPaymentRequest) =>
    api.post<ApiResponse<null>>(`/payments/${id}/confirm`, data),
}
