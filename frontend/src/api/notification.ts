import api, { type ApiResponse } from './index'

// 通知类型
export interface Notification {
  id: number
  title: string
  content: string
  type: number // 1:System, 2:Activity, 3:Private
  sender_id: number
  is_global: number // 0:No, 1:Yes
  is_read: boolean
  create_time: string
  sender?: {
    id: number
    name: string
    username: string
  }
}

// 用户简要信息（用于选择发送目标）
export interface UserBrief {
  id: number
  name: string
  username: string
}

// 创建通知请求
export interface CreateNotificationRequest {
  title: string
  content: string
  type?: string
  target_user_id?: number // 0 或不填表示全员通知
}

// 通知 API
export const notificationApi = {
  // 获取通知列表
  list: (page = 1, pageSize = 10) =>
    api.get<ApiResponse<{ list: Notification[]; total: number }>>('/notifications', {
      params: { _t: Date.now(), page, page_size: pageSize },
    }),

  // 创建通知（管理员）
  create: (data: CreateNotificationRequest) =>
    api.post<ApiResponse<Notification>>('/notifications', data),

  // 更新通知（管理员）
  update: (id: number, data: CreateNotificationRequest) => {
    return api.put<ApiResponse<Notification>>(`/notifications/${id}`, data)
  },

  get: (id: number) => {
    return api.get<ApiResponse<Notification>>(`/notifications/${id}`)
  },

  // 标记为已读
  markAsRead: (id: number) =>
    api.put<ApiResponse<null>>(`/notifications/${id}/read`),

  // 删除通知（管理员）
  delete: (id: number) =>
    api.delete<ApiResponse<null>>(`/notifications/${id}`),

  // 获取未读数量
  getUnreadCount: () =>
    api.get<ApiResponse<{ count: number }>>('/notifications/unread-count', {
      params: { _t: Date.now() },
    }),

  // 获取用户列表（管理员用于选择发送目标）
  getUsers: () =>
    api.get<ApiResponse<UserBrief[]>>('/notifications/users', {
      params: { _t: Date.now() },
    }),
}
