<script setup lang="ts">
/**
 * @file App.vue
 * @description 应用根组件
 * 负责全局主题初始化、路由视图渲染及全局 Toast 容器挂载。
 */
import { RouterView } from 'vue-router'
import { useThemeStore } from '@/stores/theme'
import ToastContainer from '@/components/common/ToastContainer.vue'

// 初始化主题 (读取本地存储或系统偏好)
useThemeStore()

// Living Light Interaction
import { onMounted, onUnmounted } from 'vue'

const updateLightPosition = (e: MouseEvent) => {
  const x = e.clientX / window.innerWidth
  const y = e.clientY / window.innerHeight
  
  document.body.style.setProperty('--light-x', x.toString())
  document.body.style.setProperty('--light-y', y.toString())
  
  // Calculate specular reflection angle (opposite to light)
  document.body.style.setProperty('--specular-x', ((0.5 - x) * 20).toString() + 'deg')
  document.body.style.setProperty('--specular-y', ((0.5 - y) * 20).toString() + 'deg')
}

onMounted(() => {
  window.addEventListener('mousemove', updateLightPosition)
})

onUnmounted(() => {
  window.removeEventListener('mousemove', updateLightPosition)
})
</script>

<template>
  <!-- 路由出口 -->
  <RouterView />
  <!-- 全局 Toast 消息容器 -->
  <ToastContainer />
</template>
