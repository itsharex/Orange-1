<script setup lang="ts">
/**
 * @file StatusBadge.vue
 * @description 状态徽章组件
 * 用于显示项目或款项的不同状态 (颜色/文案映射)。
 */
import { computed } from 'vue'

const props = defineProps<{
  status: 'active' | 'completed' | 'pending' | 'overdue' | 'notstarted' | 'archived'
}>()

// 状态配置映射 (key -> {label, class})
const statusConfig = {
  active: { label: '进行中', class: 'status-badge--active' },
  completed: { label: '已完成', class: 'status-badge--completed' },
  pending: { label: '即将交付', class: 'status-badge--pending' },
  overdue: { label: '已逾期', class: 'status-badge--danger' }, // 逾期状态使用红色样式
  notstarted: { label: '未开始', class: 'status-badge--overdue' }, // 复用 overdue 样式 (灰色/橙色)
  archived: { label: '已归档', class: 'status-badge--overdue' }
}

const config = computed(() => statusConfig[props.status] || statusConfig.active)
</script>

<template>
  <span class="status-badge" :class="config.class">
    {{ config.label }}
  </span>
</template>

<style scoped>
/* Inherits from liquid-glass.css .status-badge etc. */
/* If mapped correctly, no extra styles needed */
.status-badge--danger {
    background: rgba(255, 69, 58, 0.15);
    color: var(--color-danger);
}
</style>
