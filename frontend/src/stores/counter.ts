/**
 * @file stores/counter.ts
 * @description 计数器状态管理 (示例)
 * 一个简单的 Pinia Store 示例，用于演示 State, Getter 和 Action 用法。
 */
import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

export const useCounterStore = defineStore('counter', () => {
  const count = ref(0)
  // Double Count 计算属性
  const doubleCount = computed(() => count.value * 2)
  
  /** 增加计数 */
  function increment() {
    count.value++
  }

  return { count, doubleCount, increment }
})
