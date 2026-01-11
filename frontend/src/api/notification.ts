/**
 * @file api/notification.ts
 * @description 系统通知 API
 * 处理通知的获取、创建、标记已读及删除操作。
 */
import api, { type ApiResponse } from './index'

// 通知数据模型
export interface Notification {
  id: number
  title: string
  content: string
  type: number // 通知类型: 1:System(系统), 2:Activity(活动), 3:Private(私信)
  sender_id: number
  is_global: number // 是否全员通知: 0:否, 1:是
  is_read: boolean  // 当前用户是否已读
  create_time: string
  sender?: {        // 发送者信息
    id: number
    name: string
    username: string
  }
}

// 用户简略信息 (用于通知发送目标选择)
export interface UserBrief {
  id: number
  name: string
  username: string
}

// 创建通知请求参数
export interface CreateNotificationRequest {
  title: string
  content: string
  type?: string
  target_user_id?: number // 目标用户ID (0 或不填表示全员通知)
}

// 通知 API 集合
export const notificationApi = {
  /**
   * 获取通知列表
   * @param page 页码
   * @param pageSize 每页数量
   */
  list: (page = 1, pageSize = 10) =>
    api.get<ApiResponse<{ list: Notification[]; total: number }>>('/notifications', {
      params: { _t: Date.now(), page, page_size: pageSize },
    }),

  // 创建通知 (管理员权限)
  create: (data: CreateNotificationRequest) =>
    api.post<ApiResponse<Notification>>('/notifications', data),

  // 更新通知 (管理员权限)
  update: (id: number, data: CreateNotificationRequest) => {
    return api.put<ApiResponse<Notification>>(`/notifications/${id}`, data)
  },

  // 获取单条通知详情
  get: (id: number) => {
    return api.get<ApiResponse<Notification>>(`/notifications/${id}`)
  },

  // 标记通知为已读
  markAsRead: (id: number) =>
    api.put<ApiResponse<null>>(`/notifications/${id}/read`),

  // 删除通知 (管理员权限)
  delete: (id: number) =>
    api.delete<ApiResponse<null>>(`/notifications/${id}`),

  // 获取当前用户的未读通知数量
  getUnreadCount: () =>
    api.get<ApiResponse<{ count: number }>>('/notifications/unread-count', {
      params: { _t: Date.now() },
    }),

  // 获取可发送通知的用户列表 (管理员用)
  getUsers: () =>
    api.get<ApiResponse<UserBrief[]>>('/notifications/users', {
      params: { _t: Date.now() },
    }),
}
