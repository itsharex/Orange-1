/**
 * @file api/token.ts
 * @description 个人访问令牌相关 API
 */
import api, { type ApiResponse } from './index'

export interface PersonalAccessToken {
  id: number
  name: string
  scopes: string
  status: number // 1: Active, 0: Revoked
  last_used_at: string | null
  expires_at: string | null
  create_time: string
  // Note: Token hash is not returned here
}

export interface CreateTokenRequest {
  name: string
  expires_in: number // 天数
}

export interface CreateTokenResponse {
  token: string // The raw token, shown only once
  data: PersonalAccessToken
}

export const tokenApi = {
  // 获取令牌列表
  list: () =>
    api.get<ApiResponse<PersonalAccessToken[]>>('/tokens', { params: { _t: new Date().getTime() } }),

  // 创建令牌
  create: (data: CreateTokenRequest) =>
    api.post<ApiResponse<CreateTokenResponse>>('/tokens', data),

  // 撤销令牌
  revoke: (id: number) =>
    api.post<ApiResponse<null>>(`/tokens/${id}/revoke`),
  
  // 删除令牌
  delete: (id: number) =>
    api.delete<ApiResponse<null>>(`/tokens/${id}`),
}
