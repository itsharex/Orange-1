/**
 * @file composables/useToast.ts
 * @description 全局轻量提示 (Toast) Hook
 * 提供成功、失败、警告、信息等不同类型的 Toast 消息推送功能。
 */
import { ref } from 'vue'

export type ToastType = 'success' | 'error' | 'warning' | 'info'

export interface Toast {
  id: number
  message: string
  type: ToastType
  duration?: number // 显示时长 (ms)
}

// 全局响应式状态，确保多个组件使用的 Toast 队列一致
const toasts = ref<Toast[]>([])
let nextId = 0

export function useToast() {
  /** 移除指定 ID 的 Toast */
  const remove = (id: number) => {
    const index = toasts.value.findIndex((t) => t.id === id)
    if (index !== -1) {
      toasts.value.splice(index, 1)
    }
  }

  /**
   * 添加 Toast
   * @param message 消息内容
   * @param type 类型 (默认 info)
   * @param duration 持续时间 (默认 3000ms)
   */
  const add = (message: string, type: ToastType = 'info', duration = 3000) => {
    const id = nextId++
    const toast: Toast = { id, message, type, duration }
    toasts.value.push(toast)

    // 自动移除定时器
    if (duration > 0) {
      setTimeout(() => {
        remove(id)
      }, duration)
    }
  }

  // 快捷方法
  const success = (message: string, duration?: number) => add(message, 'success', duration)
  const error = (message: string, duration?: number) => add(message, 'error', duration)
  const warning = (message: string, duration?: number) => add(message, 'warning', duration)
  const info = (message: string, duration?: number) => add(message, 'info', duration)

  return {
    toasts,
    add,
    remove,
    success,
    error,
    warning,
    info,
  }
}
