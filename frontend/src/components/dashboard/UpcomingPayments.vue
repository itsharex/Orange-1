<script setup lang="ts">
/**
 * @file UpcomingPayments.vue
 * @description 仪表盘即将收款列表组件
 * 展示即将到期的回款计划，提示金额和剩余天数。
 */
import GlassCard from '@/components/common/GlassCard.vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const goToDetail = (id: number) => router.push(`/projects/${id}`)

// 前端展示用的收款项接口 (聚合了项目/客户信息)
interface PaymentDisplayItem {
  id: number
  project_id: number
  project_name: string
  client_name: string
  days_left: number
  amount: number
  status: string
}

const props = defineProps<{
  payments: PaymentDisplayItem[]
}>()

const statusColorMap: Record<string, string> = {
  danger: 'var(--color-danger)',
  warning: 'var(--color-warning)',
  success: 'var(--color-success)',
}

const statusBgMap: Record<string, string> = {
  danger: 'rgba(255, 69, 58, 0.05)',
  warning: 'rgba(255, 214, 10, 0.05)',
  success: 'rgba(50, 215, 75, 0.05)',
  pending: 'rgba(255, 255, 255, 0.05)', // 默认背景
}

const getStatusColor = (status: string) => statusColorMap[status] || statusColorMap.pending
const getStatusBg = (status: string) => statusBgMap[status] || statusBgMap.pending

</script>

<template>
  <GlassCard>
    <div class="glass-card-header">
      <h3 class="glass-card-title">即将到期收款</h3>
      <span class="status-badge status-badge--overdue">{{ props.payments.length }}笔待收</span>
    </div>

    <div class="flex flex-col gap-md">
      <div
        v-for="item in props.payments"
        :key="item.id"
        class="flex items-center justify-between p-md cursor-pointer hover:opacity-80 transition-opacity"
        :style="{
          background: getStatusBg(item.status),
          borderRadius: 'var(--radius-md)',
          borderLeft: `3px solid ${getStatusColor(item.status)}`,
        }"
        @click="goToDetail(item.project_id)"
      >
        <div>
          <div class="font-semibold">{{ item.project_name }}</div>
          <div class="text-sm text-secondary mt-sm">{{ item.client_name }}</div>
          <div class="text-sm mt-sm" :style="{ color: getStatusColor(item.status) }">
            {{ item.days_left }}天后到期
          </div>
        </div>
        <div class="text-xl font-bold">¥{{ item.amount.toLocaleString() }}</div>
      </div>
      <div v-if="props.payments.length === 0" class="text-center text-secondary py-4">
        暂无即将到期款项
      </div>
    </div>
  </GlassCard>
</template>
