<script setup lang="ts">
import GlassCard from '@/components/common/GlassCard.vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const goToDetail = (id: number) => router.push(`/projects/${id}`)



// Extend Payment type for display if necessary, but for now assuming backend returns needed fields
// Actually backend Payment struct: ID, ProjectID, Phase, Amount, Status, PaidTime, DueDate
// The UI needs: id, project (name), client (name), daysLeft (calculated), amount, status
// We might need to transform data in parent or here. Let's look at api/project.ts Project and Payment definitions.
// Wait, dashboardApi returns specific structs potentially?
// dashboardApi.getUpcomingPayments returns Payment[].
// Let's assume the parent passes the array and we iterate.
// However, the UI displays "Project Name", "Client Name", "Days Left".
// The standard Payment struct might not have Project Name or Client Name joined.
// Let's check api/project.ts first to be sure.
// But for now, I will define a UI specific interface for the prop or use `any` temporarily if I can't check.
// Better: Check api/project.ts. I haven't checked it yet.
// For safety, I'll define an interface locally that matches what we need to display.

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
  pending: 'rgba(255, 255, 255, 0.05)', // Default for others
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
