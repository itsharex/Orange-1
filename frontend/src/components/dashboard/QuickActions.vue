<script setup lang="ts">
import GlassCard from '@/components/common/GlassCard.vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const handleAction = (name: string) => {
  if (name === '新建项目') {
    router.push('/projects/create')
  } else if (name === '收款日历') {
    router.push('/calendar')
  } else if (name === '生成报表') {
    router.push('/analytics')
  } else if (name === '添加收款') {
    router.push('/payment/create')
  }
}

const actions = [
  {
    name: '新建项目',
    icon: 'ri-add-circle-line',
    color: '#3b82f6',
    bg: 'linear-gradient(135deg, rgba(59, 130, 246, 0.2), rgba(59, 130, 246, 0.05))',
  },
  {
    name: '添加收款',
    icon: 'ri-money-dollar-box-line',
    color: '#10b981',
    bg: 'linear-gradient(135deg, rgba(16, 185, 129, 0.2), rgba(16, 185, 129, 0.05))',
  },
  {
    name: '生成报表',
    icon: 'ri-file-chart-line',
    color: '#8b5cf6',
    bg: 'linear-gradient(135deg, rgba(139, 92, 246, 0.2), rgba(139, 92, 246, 0.05))',
  },
  {
    name: '收款日历',
    icon: 'ri-calendar-check-line',
    color: '#f59e0b',
    bg: 'linear-gradient(135deg, rgba(245, 158, 11, 0.2), rgba(245, 158, 11, 0.05))',
  },
]
</script>

<template>
  <GlassCard>
    <div class="glass-card-header">
      <h3 class="glass-card-title">快捷操作</h3>
    </div>
    <div class="quick-actions-grid">
      <div
        v-for="action in actions"
        :key="action.name"
        class="quick-action cursor-pointer"
        @click="handleAction(action.name)"
      >
        <div class="quick-action-icon" :style="{ background: action.bg, color: action.color }">
          <i :class="action.icon"></i>
        </div>
        <span class="text-sm font-medium">{{ action.name }}</span>
      </div>
    </div>
  </GlassCard>
</template>

<style scoped>
.quick-actions-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 8px;
}

/* 768px-1024px 时变为4列（快捷操作独立一行，宽度充足） */
@media (min-width: 768px) and (max-width: 1024px) {
  .quick-actions-grid {
    grid-template-columns: repeat(4, 1fr);
  }
}
</style>
