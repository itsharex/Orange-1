import { ref } from 'vue'
import { defineStore } from 'pinia'

export const useThemeStore = defineStore('theme', () => {
  const theme = ref(localStorage.getItem('theme') || 'auto')
  const effectiveTheme = ref('light')

  function applyTheme() {
    const root = document.documentElement
    let targetTheme = theme.value

    if (targetTheme === 'auto') {
      const isDark = window.matchMedia('(prefers-color-scheme: dark)').matches
      targetTheme = isDark ? 'dark' : 'light'
    }

    effectiveTheme.value = targetTheme
    root.setAttribute('data-theme', targetTheme)
  }

  function setTheme(newTheme: string) {
    theme.value = newTheme
    localStorage.setItem('theme', newTheme)
    applyTheme()
  }

  function toggleTheme() {
    // Cycle: Light <-> Dark (Skip Auto)
    if (theme.value === 'dark') {
      setTheme('light')
    } else {
      setTheme('dark')
    }
  }

  // Listen for system theme changes
  window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', () => {
    if (theme.value === 'auto') {
      applyTheme()
    }
  })

  // Initialize
  applyTheme()

  return { theme, effectiveTheme, setTheme, toggleTheme }
})
