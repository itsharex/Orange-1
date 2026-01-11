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
</script>

<template>
  <!-- 路由出口 -->
  <RouterView />
  <!-- 全局 Toast 消息容器 -->
  <ToastContainer />
</template>
