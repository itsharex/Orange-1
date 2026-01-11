/**
 * @file api/dictionary.ts
 * @description 数据字典管理 API
 * 用于管理系统中的枚举值，如项目类型、支付阶段等。
 */
import api, { type ApiResponse } from './index'

// 字典定义
export interface Dictionary {
  id: number
  code: string   // 字典编码 (唯一标识)
  name: string   // 字典名称
  status: number // 状态
  remark: string // 备注
}

// 字典项定义
export interface DictionaryItem {
  id: number
  dictionary_id: number // 所属字典 ID
  label: string         // 显示名称
  value: string         // 实际值
  sort: number          // 排序 (越小越靠前)
  status: number        // 状态
}

// 创建/修改字典项请求参数
export interface CreateDictionaryItemRequest {
  label: string
  value: string
  sort?: number
}

// 字典 API 集合
export const dictionaryApi = {
  // 获取字典列表
  list: () =>
    api.get<ApiResponse<Dictionary[]>>('/dictionaries'),

  // 获取字典项
  getItems: (code: string) =>
    api.get<ApiResponse<DictionaryItem[]>>(`/dictionaries/${code}/items`, {
      params: { _t: Date.now() },
    }),

  // 创建字典项
  createItem: (code: string, data: CreateDictionaryItemRequest) =>
    api.post<ApiResponse<DictionaryItem>>(`/dictionaries/${code}/items`, data),

  // 更新字典项
  updateItem: (code: string, id: number, data: CreateDictionaryItemRequest) =>
    api.put<ApiResponse<DictionaryItem>>(`/dictionaries/${code}/items/${id}`, data),

  // 删除字典项
  deleteItem: (code: string, id: number) =>
    api.delete<ApiResponse<null>>(`/dictionaries/${code}/items/${id}`),
}
