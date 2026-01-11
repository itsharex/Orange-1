<script setup lang="ts">
/**
 * @file ProjectList.vue
 * @description 仪表盘近期项目列表组件
 * 展示项目概览表格，包括名称、客户、金额、进度和状态。
 */
import GlassCard from '@/components/common/GlassCard.vue'
import StatusBadge from '@/components/common/StatusBadge.vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const goToDetail = (id: number) => router.push(`/projects/${id}`)

import type { Project } from '@/api/project'

const props = defineProps<{
  projects: Project[] // 项目列表数据
}>()

</script>

<template>
  <GlassCard class="project-list-card">
    <div class="glass-card-header">
      <div>
        <h3 class="glass-card-title">近期项目</h3>
        <p class="glass-card-subtitle">进行中的项目</p>
      </div>
      <router-link to="/projects" class="btn btn-ghost btn-sm">
        查看全部 <i class="ri-arrow-right-line"></i>
      </router-link>
    </div>

    <div class="table-scroll-container">
      <table class="data-table project-table">
        <thead>
          <tr>
            <th style="width: 200px; min-width: 200px; max-width: 200px;">项目名称</th>
            <th>客户</th>
            <th>合同金额</th>
            <th>收款进度</th>
            <th>状态</th>
          </tr>
        </thead>
        <tbody></tbody>
        <tbody>
          <tr
            v-for="p in props.projects"
            :key="p.name"
            class="project-row cursor-pointer hover:bg-white/5 transition-colors"
            @click="goToDetail(p.id)"
          >
            <td class="font-medium" style="max-width: 200px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;" :title="p.name">{{ p.name }}</td>
            <td class="text-secondary">{{ p.company }}</td>
            <td>¥{{ p.total_amount.toLocaleString() }}</td>
            <td>
              <div class="flex items-center gap-sm">
                <div class="progress-bar" style="width: 80px">
                  <div class="progress-bar-fill" :style="{ width: ((p.received_amount || 0) / p.total_amount * 100).toFixed(0) + '%' }"></div>
                </div>
                <span class="text-sm">{{ ((p.received_amount || 0) / p.total_amount * 100).toFixed(0) }}%</span>
              </div>
            </td>
            <td>
              <StatusBadge :status="p.status" />
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </GlassCard>
</template>

<style scoped>
.table-scroll-container {
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

.project-table {
  min-width: 500px;
}

.project-table th,
.project-table td {
  white-space: nowrap;
}

/* 响应式：小屏幕下缩小行高和字体 */
@media (max-width: 768px) {
  .project-table th,
  .project-table td {
    padding: var(--spacing-sm) var(--spacing-sm);
    font-size: 13px;
  }

  .project-table th {
    font-size: 12px;
  }

  .progress-bar {
    width: 60px !important;
  }
}
</style>
