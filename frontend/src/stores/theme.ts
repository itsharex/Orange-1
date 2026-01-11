/**
 * @file stores/theme.ts
 * @description 主题状态管理
 * 支持 Light (亮色模式), Dark (暗色模式) 和 Auto (跟随系统)。
 */
import { ref } from 'vue'
import { defineStore } from 'pinia'

export const useThemeStore = defineStore('theme', () => {
  // 当前选择的主题模式: 'light' | 'dark' | 'auto'
  const theme = ref(localStorage.getItem('theme') || 'auto')
  // 实际生效的主题 (当 theme 为 auto 时，会根据系统偏好计算为 light 或 dark)
  const effectiveTheme = ref('light')

  /** 应用主题到 DOM */
  function applyTheme() {
    const root = document.documentElement
    let targetTheme = theme.value

    // 如果是自动模式，根据系统偏好决定
    if (targetTheme === 'auto') {
      const isDark = window.matchMedia('(prefers-color-scheme: dark)').matches
      targetTheme = isDark ? 'dark' : 'light'
    }

    effectiveTheme.value = targetTheme
    // 设置 html 标签的 data-theme 属性，供 CSS 变量使用
    root.setAttribute('data-theme', targetTheme)
  }

  /**
   * 手动设置主题
   * @param newTheme 主题模式
   */
  function setTheme(newTheme: string) {
    theme.value = newTheme
    localStorage.setItem('theme', newTheme)
    applyTheme()
  }

  /**
   * 切换主题 (Light/Dark 循环)
   * 跳过 Auto 模式
   */
  function toggleTheme() {
    // Cycle: Light <-> Dark (Skip Auto)
    if (theme.value === 'dark') {
      setTheme('light')
    } else {
      setTheme('dark')
    }
  }

  // 监听系统主题变化事件
  // 当设置为自动模式时，系统主题切换会立即反映到界面上
  window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', () => {
    if (theme.value === 'auto') {
      applyTheme()
    }
  })

  // 初始化应用
  applyTheme()

  return { theme, effectiveTheme, setTheme, toggleTheme }
})
