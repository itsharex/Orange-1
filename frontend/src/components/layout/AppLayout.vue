<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import AppSidebar from './AppSidebar.vue'
import AppHeader from './AppHeader.vue'
import { useLayoutStore } from '@/stores/layout'

const layoutStore = useLayoutStore()
const appBackground = ref<HTMLElement | null>(null)

// Liquid Glass 动态效果逻辑
const handleMouseMove = (e: MouseEvent) => {
  // 更新背景光效位置
  if (appBackground.value) {
    const xPercent = (e.clientX / window.innerWidth) * 100
    const yPercent = (e.clientY / window.innerHeight) * 100
    
    appBackground.value.style.setProperty('--light-x', `${xPercent}%`)
    appBackground.value.style.setProperty('--light-y', `${yPercent}%`)
  }

  // 更新玻璃卡片的高光效果
  updateSpecularHighlights(e)
}

const updateSpecularHighlights = (e: MouseEvent) => {
  const glassCards = document.querySelectorAll('.glass-card, .liquid-glass')
  
  glassCards.forEach((card) => {
    const el = card as HTMLElement
    const rect = el.getBoundingClientRect()
    const x = ((e.clientX - rect.left) / rect.width) * 100
    const y = ((e.clientY - rect.top) / rect.height) * 100

    // 只在卡片附近时更新，优化性能
    if (x >= -20 && x <= 120 && y >= -20 && y <= 120) {
      el.style.setProperty('--specular-x', `${x}%`)
      el.style.setProperty('--specular-y', `${y}%`)
    }
  })
}

// 侧边栏滚动效果
const mainContent = ref<HTMLElement | null>(null)
const handleScroll = (e: Event) => {
  const target = e.target as HTMLElement
  const sidebar = document.querySelector('.sidebar')
  if (sidebar) {
    if (target.scrollTop > 50) {
      sidebar.classList.add('scrolled')
    } else {
      sidebar.classList.remove('scrolled')
    }
  }
}

onMounted(() => {
  document.addEventListener('mousemove', handleMouseMove)
  // 注意：如果是 window 滚动则用 window，如果是 mainContent 滚动则用 mainContent
  // 原型是 mainContent overflow-y: auto
  if (mainContent.value) {
    mainContent.value.addEventListener('scroll', handleScroll)
  }
})

onUnmounted(() => {
  document.removeEventListener('mousemove', handleMouseMove)
  if (mainContent.value) {
    mainContent.value.removeEventListener('scroll', handleScroll)
  }
})
</script>

<template>
  <div class="app-container">
    <!-- SVG Filters -->
    <svg class="svg-filters" style="position: absolute; width: 0; height: 0; overflow: hidden">
      <defs>
        <filter id="glass-refraction" x="-20%" y="-20%" width="140%" height="140%">
          <feTurbulence type="fractalNoise" baseFrequency="0.015" numOctaves="3" result="noise" seed="1" />
          <feDisplacementMap in="SourceGraphic" in2="noise" scale="3" xChannelSelector="R" yChannelSelector="G" />
        </filter>
        <filter id="glass-blur">
          <feGaussianBlur in="SourceGraphic" stdDeviation="20" />
        </filter>
        <filter id="glass-glow" x="-50%" y="-50%" width="200%" height="200%">
          <feGaussianBlur in="SourceGraphic" stdDeviation="8" result="blur" />
          <feColorMatrix in="blur" type="matrix" values="1 0 0 0 0  0 1 0 0 0  0 0 1 0 0  0 0 0 18 -7" result="glow" />
          <feBlend in="SourceGraphic" in2="glow" mode="screen" />
        </filter>
      </defs>
    </svg>

    <!-- Background -->
    <div class="app-background" ref="appBackground"></div>

    <!-- Sidebar -->
    <AppSidebar />

    <!-- Main Content -->
    <main 
      class="main-content" 
      id="mainContent"
      ref="mainContent"
      :class="{ 'ml-[76px]': layoutStore.sidebarCollapsed }"
    >
      <AppHeader />
      
      <div class="view-content animate-fade-in">
        <RouterView v-slot="{ Component }">
          <Transition name="fade" mode="out-in">
            <component :is="Component" />
          </Transition>
        </RouterView>
      </div>
    </main>
  </div>
</template>

<style scoped>
.main-content {
  transition: margin-left var(--transition-spring);
  height: 100vh;
  overflow-y: auto;
  overflow-x: hidden;
}

/* 简单的淡入淡出过渡 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
