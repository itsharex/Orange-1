/**
 * @file stores/layout.ts
 * @description 布局状态管理
 * 控制侧边栏折叠状态等 UI 布局相关的响应式状态。
 */
import { ref } from 'vue'
import { defineStore } from 'pinia'

export const useLayoutStore = defineStore('layout', () => {
  // 侧边栏是否折叠 (通过 v-model 或 toggle 切换)
  const sidebarCollapsed = ref(false)

  /** 切换侧边栏状态 */
  function toggleSidebar() {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  /** 设置侧边栏折叠状态 */
  function setSidebarCollapsed(value: boolean) {
    sidebarCollapsed.value = value
  }

  return { sidebarCollapsed, toggleSidebar, setSidebarCollapsed }
})
