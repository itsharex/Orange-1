<script setup lang="ts">
import { ref } from 'vue'

const visible = ref(false)
const title = ref('确认')
const message = ref('')
const resolvePromise = ref<((value: boolean) => void) | null>(null)

const open = (options: { title?: string; message: string }) => {
  title.value = options.title || '确认'
  message.value = options.message
  visible.value = true

  return new Promise<boolean>((resolve) => {
    resolvePromise.value = resolve
  })
}

const handleConfirm = () => {
  visible.value = false
  resolvePromise.value?.(true)
}

const handleCancel = () => {
  visible.value = false
  resolvePromise.value?.(false)
}

defineExpose({ open })
</script>

<template>
  <Teleport to="body">
    <Transition name="fade">
      <div v-if="visible" class="confirm-overlay" @click.self="handleCancel">
        <div class="confirm-modal">
          <h3 class="confirm-title">{{ title }}</h3>
          <p class="confirm-message">{{ message }}</p>
          <div class="confirm-actions">
            <button class="btn btn-ghost" @click="handleCancel">取消</button>
            <button class="btn btn-primary" @click="handleConfirm">确认</button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.confirm-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10000;
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  transition: all 0.5s cubic-bezier(0.19, 1, 0.22, 1);
}

[data-theme='dark'] .confirm-overlay {
  background: rgba(0, 0, 0, 0.5);
}

.confirm-modal {
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
  border-radius: var(--radius-xl);
  box-shadow: 
    var(--glass-shadow-outer),
    0 20px 40px rgba(0,0,0,0.1); /* Extra depth for modal */
    
  padding: 32px;
  min-width: 320px;
  max-width: 400px;
  text-align: center;
  position: relative;
  overflow: hidden;
  isolation: isolate;
  
  transform-origin: center center;
}

/* Specular Highlight */
.confirm-modal::before {
  content: '';
  position: absolute;
  inset: 0;
  background: var(--glass-specular);
  pointer-events: none;
  z-index: 0;
  border-radius: inherit;
  opacity: 0.6;
  mix-blend-mode: overlay;
}

/* Inner Glow */
.confirm-modal::after {
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

[data-theme='dark'] .confirm-modal::after {
  box-shadow:
    inset 0 0 0 1px rgba(255, 255, 255, 0.08),
    inset 0 1px 2px rgba(255, 255, 255, 0.05);
}

/* Content z-index fix */
.confirm-title,
.confirm-message,
.confirm-actions {
  position: relative;
  z-index: 2;
}

.confirm-title {
  font-family: var(--font-display);
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 12px;
  color: var(--text-primary);
  letter-spacing: -0.01em;
}

.confirm-message {
  font-family: var(--font-text);
  font-size: 15px;
  color: var(--text-secondary);
  margin-bottom: 28px;
  line-height: 1.6;
}

.confirm-actions {
  display: flex;
  gap: 16px;
  justify-content: center;
}

.confirm-actions .btn {
  min-width: 90px;
  font-weight: 500;
  transition: all 0.2s cubic-bezier(0.2, 0, 0, 1);
}

.confirm-actions .btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
}

.confirm-actions .btn:active {
  transform: translateY(0) scale(0.98);
}

/* Transitions - Apple Spring */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.4s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.fade-enter-active .confirm-modal,
.fade-leave-active .confirm-modal {
  transition: transform 0.5s cubic-bezier(0.19, 1, 0.22, 1);
}

.fade-enter-from .confirm-modal {
  transform: scale(0.9) translateY(20px);
}

.fade-leave-to .confirm-modal {
  transform: scale(0.95) translateY(10px);
}
</style>
