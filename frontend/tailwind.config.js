/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  darkMode: ['selector', '[data-theme="dark"]'],
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: '#FF9F0A',
          light: '#FFB340',
          dark: '#cc7a00',
        },
        secondary: '#5E5CE6',
        success: '#32D74B',
        warning: '#FFD60A',
        danger: '#FF453A',
        info: '#64D2FF',
        teal: '#30B0C7',
        // 背景色引用 CSS 变量以支持动态切换
        bg: {
          base: 'var(--bg-base)',
          elevated: 'var(--bg-elevated)',
          content: 'var(--bg-content)',
          sidebar: 'var(--bg-sidebar)',
        },
        // 文本色引用 CSS 变量
        text: {
          primary: 'var(--text-primary)',
          secondary: 'var(--text-secondary)',
          tertiary: 'var(--text-tertiary)',
          quaternary: 'var(--text-quaternary)',
        }
      },
      fontFamily: {
        display: ['"SF Pro Display"', '-apple-system', 'BlinkMacSystemFont', 'system-ui', 'sans-serif'],
        text: ['"SF Pro Text"', '-apple-system', 'BlinkMacSystemFont', 'system-ui', 'sans-serif'],
        mono: ['"SF Mono"', 'Menlo', 'Monaco', 'Courier New', 'monospace'],
      },
      borderRadius: {
        xs: '6px',
        sm: '10px',
        md: '14px',
        lg: '18px',
        xl: '22px',
        '2xl': '28px',
      },
      spacing: {
        xs: '4px',
        sm: '8px',
        md: '16px',
        lg: '24px',
        xl: '32px',
        '2xl': '48px',
      }
    },
  },
  plugins: [],
}
