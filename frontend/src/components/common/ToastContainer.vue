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
}

.toast-item {
  pointer-events: auto;
  min-width: 320px;
  max-width: 90vw;
  padding: 14px 20px;
  border-radius: var(--radius-xl); /* Capsule shape */
  
  /* Liquid Glass Base */
  background: var(--bg-elevated);
  backdrop-filter: 
    blur(var(--glass-blur)) 
    saturate(var(--glass-saturation))
    brightness(var(--glass-brightness));
  -webkit-backdrop-filter: 
    blur(var(--glass-blur)) 
    saturate(var(--glass-saturation))
    brightness(var(--glass-brightness));
    
  border: var(--glass-border);
  box-shadow: var(--glass-shadow-outer), var(--glass-shadow-inner);
  
  color: var(--text-primary);
  display: flex;
  align-items: center;
  gap: 14px;
  cursor: pointer;
  font-family: var(--font-text);
  font-size: 14px;
  font-weight: 500;
  letter-spacing: -0.01em;
  
  position: relative;
  overflow: hidden;
  transition: all 0.4s cubic-bezier(0.2, 0, 0, 1);
}

/* Specular Highlight (High Gloss) */
.toast-item::before {
  content: '';
  position: absolute;
  inset: 0;
  background: var(--glass-specular);
  pointer-events: none;
  z-index: 0;
  border-radius: inherit;
  opacity: 0.8;
  mix-blend-mode: overlay;
}

/* Inner Glow / Edge Light */
.toast-item::after {
  content: '';
  position: absolute;
  inset: 0;
  border-radius: inherit;
  box-shadow:
    inset 0 0 0 1px rgba(255, 255, 255, 0.15),
    inset 0 1px 2px rgba(255, 255, 255, 0.2);
  pointer-events: none;
  z-index: 1;
}

[data-theme="dark"] .toast-item::after {
  box-shadow:
    inset 0 0 0 1px rgba(255, 255, 255, 0.08),
    inset 0 1px 2px rgba(255, 255, 255, 0.05);
}

/* Icon & Content styling ensure they are above effects */
.icon, .message {
  position: relative;
  z-index: 2;
}

/* Type-specific styling - using subtle tints instead of heavy borders */
.toast-item.success {
  background: rgba(50, 215, 75, 0.15); /* --color-success tinted */
  border-color: rgba(50, 215, 75, 0.3);
}
.toast-item.success .icon {
  color: var(--color-success);
}

.toast-item.error {
  background: rgba(255, 69, 58, 0.15); /* --color-danger tinted */
  border-color: rgba(255, 69, 58, 0.3);
}
.toast-item.error .icon {
  color: var(--color-danger);
}

.toast-item.warning {
  background: rgba(255, 214, 10, 0.15); /* --color-warning tinted */
  border-color: rgba(255, 214, 10, 0.3);
}
.toast-item.warning .icon {
  color: var(--color-warning);
}

.toast-item.info {
  background: rgba(10, 132, 255, 0.15); /* --color-info tinted */
  border-color: rgba(10, 132, 255, 0.3);
}
.toast-item.info .icon {
  color: var(--color-info);
}

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

/* Hover Effect */
.toast-item:hover {
  transform: translateY(-2px) scale(1.02);
  box-shadow: 
    0 12px 24px rgba(0, 0, 0, 0.12),
    var(--glass-shadow-inner);
}

/* Transitions */
.toast-enter-active,
.toast-leave-active {
  transition: all 0.5s cubic-bezier(0.19, 1, 0.22, 1); /* Apple-like spring */
}

.toast-enter-from {
  opacity: 0;
  transform: translateY(-40px) scale(0.9);
  filter: blur(10px);
}

.toast-leave-to {
  opacity: 0;
  transform: translateY(-20px) scale(0.95);
  filter: blur(10px);
}
</style>
