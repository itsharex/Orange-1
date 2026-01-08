import api, { type ApiResponse } from './index'

// 用户类型
export interface User {
  id: number
  username: string
  name: string
  email: string
  phone: string
  avatar: string
  role: string
  department: string
  position: string
  status: number
}

// 登录请求
export interface LoginRequest {
  username: string
  password: string
}

// 登录响应
export interface LoginResponse {
  token: string
  user: User
}

// 注册请求
export interface RegisterRequest {
  username: string
  name: string
  email?: string
  phone?: string
  password: string
}

// 修改密码请求
export interface ChangePasswordRequest {
  old_password: string
  new_password: string
}

// 更新个人信息请求
export interface UpdateProfileRequest {
  name?: string
  email?: string
  phone?: string
  department?: string
  position?: string
}

// 认证 API
export const authApi = {
  // 登录
  login: (data: LoginRequest) =>
    api.post<ApiResponse<LoginResponse>>('/auth/login', data),

  // 注册
  register: (data: RegisterRequest) =>
    api.post<ApiResponse<null>>('/auth/register', data),

  // 退出登录
  logout: () =>
    api.post<ApiResponse<null>>('/auth/logout'),

  // 获取当前用户
  getCurrentUser: () =>
    api.get<ApiResponse<User>>('/users/me'),

  // 更新个人信息
  updateProfile: (data: UpdateProfileRequest) =>
    api.put<ApiResponse<User>>('/users/me', data),

  // 修改密码
  changePassword: (data: ChangePasswordRequest) =>
    api.put<ApiResponse<null>>('/users/me/password', data),
}
