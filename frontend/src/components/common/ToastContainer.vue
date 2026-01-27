<script setup lang="ts">
/**
 * @file ToastContainer.vue
 * @description 全局 Toast 消息容器组件
 * 负责渲染和管理所有激活状态的 Toast 弹窗队列。
 * 通过 TransitionGroup 实现列表动画。
 */
import { useToast } from '@/composables/useToast'

const { toasts, remove } = useToast()
</script>

<template>
  <div class="toast-container">
    <TransitionGroup name="toast">
      <div
        v-for="toast in toasts"
        :key="toast.id"
        class="toast-item"
        :class="toast.type"
        @click="remove(toast.id)"
      >
        <div class="icon">
          <i v-if="toast.type === 'success'" class="ri-checkbox-circle-fill"></i>
          <i v-else-if="toast.type === 'error'" class="ri-close-circle-fill"></i>
          <i v-else-if="toast.type === 'warning'" class="ri-error-warning-fill"></i>
          <i v-else class="ri-information-fill"></i>
        </div>
        <span class="message">{{ toast.message }}</span>
      </div>
    </TransitionGroup>
  </div>
</template>

<style scoped>
/* ============================================
   Toast - Liquid Glass Capsule
   ============================================ */
.toast-container {
  position: fixed;
  top: 32px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 9999;
  display: flex;
  flex-direction: column;
  gap: 16px;
  pointer-events: none;
  /* Fix animation jump: maintain width when items become absolute */
  min-width: 340px;
  max-width: 90vw;
}

.toast-item {
  pointer-events: auto;
  width: 100%; /* Fill container */
  /* min-width moved to container to prevent collapse */
  padding: 14px 24px;
  border-radius: 50px; /* Full Capsule */
  
  /* Deep Liquid Glass */
  background: var(--bg-elevated);
  backdrop-filter: blur(24px) saturate(180%);
  -webkit-backdrop-filter: blur(24px) saturate(180%);
    
  border: 1px solid rgba(255, 255, 255, 0.3);
  box-shadow: 
    0 10px 30px -5px rgba(0, 0, 0, 0.1),
    0 4px 10px -2px rgba(0, 0, 0, 0.05),
    var(--glass-shadow-inner);
  
  color: var(--text-primary);
  display: flex;
  align-items: center;
  gap: 16px;
  cursor: pointer;
  font-family: var(--font-text);
  font-size: 14px;
  font-weight: 500;
  letter-spacing: -0.01em;
  
  position: relative;
  overflow: hidden;
  transition: all 0.4s cubic-bezier(0.25, 0.8, 0.25, 1);
  will-change: transform, box-shadow;
}

/* Hover Levitation */
.toast-item:hover {
  transform: translateY(-4px) scale(1.02);
  box-shadow: 
    0 20px 40px -5px rgba(0, 0, 0, 0.15),
    0 8px 16px -4px rgba(0, 0, 0, 0.1),
    var(--glass-shadow-inner);
}

/* Inner Glow */
.toast-item::after {
  content: '';
  position: absolute;
  inset: 0;
  border-radius: inherit;
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.5);
  pointer-events: none;
  z-index: 1;
}

[data-theme="dark"] .toast-item {
  border-color: rgba(255, 255, 255, 0.15);
  box-shadow: 
    0 10px 40px -10px rgba(0, 0, 0, 0.5),
    inset 0 1px 0 rgba(255, 255, 255, 0.1);
}

/* Type-specific Neon Glow */
.toast-item.success {
  background: linear-gradient(90deg, rgba(50, 215, 75, 0.1), rgba(50, 215, 75, 0.05));
  border-color: rgba(50, 215, 75, 0.3);
  box-shadow: 0 4px 15px rgba(50, 215, 75, 0.2), var(--glass-shadow-inner);
}
.toast-item.success .icon { color: #32D74B; }

.toast-item.error {
  background: linear-gradient(90deg, rgba(255, 69, 58, 0.1), rgba(255, 69, 58, 0.05));
  border-color: rgba(255, 69, 58, 0.3);
  box-shadow: 0 4px 15px rgba(255, 69, 58, 0.2), var(--glass-shadow-inner);
}
.toast-item.error .icon { color: #FF453A; }

.toast-item.warning {
  background: linear-gradient(90deg, rgba(255, 214, 10, 0.1), rgba(255, 214, 10, 0.05));
  border-color: rgba(255, 214, 10, 0.3);
  box-shadow: 0 4px 15px rgba(255, 214, 10, 0.2), var(--glass-shadow-inner);
}
.toast-item.warning .icon { color: #FFD60A; }

.toast-item.info {
  background: linear-gradient(90deg, rgba(10, 132, 255, 0.1), rgba(10, 132, 255, 0.05));
  border-color: rgba(10, 132, 255, 0.3);
  box-shadow: 0 4px 15px rgba(10, 132, 255, 0.2), var(--glass-shadow-inner);
}
.toast-item.info .icon { color: #0A84FF; }

.icon {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  filter: drop-shadow(0 2px 4px rgba(0,0,0,0.1));
}

.message {
  flex: 1;
  line-height: 1.5;
}

/* Transitions - Bouncy Spring */
.toast-enter-active {
  transition: all 0.6s cubic-bezier(0.34, 1.56, 0.64, 1);
}
.toast-leave-active {
  transition: all 0.5s cubic-bezier(0.25, 1, 0.5, 1);
  position: absolute;
  left: 0; /* Anchor to container */
  /* Ensure it doesn't shrink when absolute */
  width: 100%; 
  z-index: -1;
}
.toast-move {
  transition: all 0.5s cubic-bezier(0.25, 1, 0.5, 1);
}

.toast-enter-from {
  opacity: 0;
  transform: translateY(-30px) scale(0.9) rotateX(-10deg);
}

.toast-leave-to {
  opacity: 0;
  transform: translateY(-30px) scale(0.9);
  /* filter: blur(8px); Removed blur to improve performance/glitchiness */
}
</style>
