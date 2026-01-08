<script setup lang="ts">
import { RouterView } from 'vue-router'
import { useThemeStore } from '@/stores/theme'
import ToastContainer from '@/components/common/ToastContainer.vue'

// Initialize theme
useThemeStore()
</script>

<template>
  <RouterView />
  <ToastContainer />
</template>
