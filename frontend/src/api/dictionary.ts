import api, { type ApiResponse } from './index'

// 字典类型
export interface Dictionary {
  id: number
  code: string
  name: string
  status: number
  remark: string
}

// 字典项类型
export interface DictionaryItem {
  id: number
  dictionary_id: number
  label: string
  value: string
  sort: number
  status: number
}

// 创建字典项请求
export interface CreateDictionaryItemRequest {
  label: string
  value: string
  sort?: number
}

// 字典 API
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
