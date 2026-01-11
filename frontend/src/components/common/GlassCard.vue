<script setup lang="ts">
/**
 * @file GlassCard.vue
 * @description 毛玻璃风格卡片组件
 * 封装了项目中统一的 Glassmorphism 样式容器。
 */
defineProps<{
  flat?: boolean      // 是否扁平化 (无阴影/边框减弱)
  hover?: boolean     // 是否启用悬停浮动效果
  noPadding?: boolean // 是否去除内边距
}>()
</script>

<template>
  <div 
    class="glass-card" 
    :class="{ 
      'glass-card--flat': flat, 
      'p-0': noPadding
    }"
  >
    <slot name="header"></slot>
    <slot></slot>
  </div>
</template>

<style scoped>
/* Inherits from liquid-glass.css */
</style>
