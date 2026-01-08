<script setup lang="ts">
defineProps<{
  flat?: boolean
  hover?: boolean
  noPadding?: boolean
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
