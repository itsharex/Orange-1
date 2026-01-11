<script setup lang="ts">
/**
 * @file IncomeChart.vue
 * @description 收入趋势图表组件
 * 使用 Chart.js 渲染折线图，支持周/月/年维度切换。
 * 自动适配亮色/暗色主题。
 */
import { ref, onMounted, watch, computed } from 'vue'
import Chart from 'chart.js/auto'
import type { ChartConfiguration, TooltipItem } from 'chart.js'
import { useThemeStore } from '@/stores/theme'
import GlassCard from '@/components/common/GlassCard.vue'

const props = defineProps<{
  labels: string[]        // X轴标签 (日期)
  values: number[]        // Y轴数据 (金额)
  modelValue?: string     // 当前选中的时间维度 (v-model)
}>()

const emit = defineEmits(['update:modelValue', 'change'])

const currentPeriod = ref(props.modelValue || 'month')

/**
 * 切换时间维度
 * @param p 维度标识 ('week' | 'month' | 'year')
 */
const setPeriod = (p: string) => {
  currentPeriod.value = p
  emit('update:modelValue', p)
  emit('change', p)
}

watch(() => props.modelValue, (val) => {
  if (val) currentPeriod.value = val
})

const canvasRef = ref<HTMLCanvasElement | null>(null)
let chartInstance: Chart | null = null
const themeStore = useThemeStore()

// 监听主题变化，重建图表以更新颜色
watch(() => themeStore.effectiveTheme, () => {
  if (chartInstance) {
    chartInstance.destroy()
    chartInstance = null
  }
  initChart()
})

/**
 * 初始化图表
 * 配置 Chart.js 实例，包括渐变背景、坐标轴样式、Tooltip 等
 */
function initChart() {
  if (!canvasRef.value) return

  const ctx = canvasRef.value.getContext('2d')
  if (!ctx) return
  
  const isDark = themeStore.effectiveTheme === 'dark'

  // Create gradient
  const gradient = ctx.createLinearGradient(0, 0, 0, 280)
  gradient.addColorStop(0, 'rgba(255, 159, 10, 0.3)')
  gradient.addColorStop(1, 'rgba(255, 159, 10, 0)')

  const config: ChartConfiguration<'line'> = {
    type: 'line',
    data: {
      labels: props.labels,
      datasets: [
        {
          label: "收款金额",
          data: props.values,
          borderColor: '#FF9F0A',
          backgroundColor: gradient,
          tension: 0.4,
          fill: true,
          pointBackgroundColor: '#fff',
          pointBorderColor: '#FF9F0A',
          pointBorderWidth: 2,
          pointRadius: 5,
          pointHoverRadius: 7,
        },
      ],
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: { display: false },
        tooltip: {
          backgroundColor: isDark
            ? "rgba(44, 44, 46, 0.95)"
            : "rgba(255, 255, 255, 0.95)",
          titleColor: isDark ? "#f5f5f7" : "#1d1d1f",
          bodyColor: isDark ? "#f5f5f7" : "#1d1d1f",
          borderColor: isDark ? "rgba(255,255,255,0.1)" : "rgba(0,0,0,0.1)",
          borderWidth: 1,
          padding: 12,
          cornerRadius: 8,
          displayColors: false,
          callbacks: {
            label: (context: TooltipItem<'line'>) => `¥${Number(context.raw).toLocaleString()}`,
          },
        },
      },
      scales: {
        y: {
          beginAtZero: true,
          grid: {
            color: isDark ? "rgba(255,255,255,0.06)" : "rgba(0,0,0,0.05)",
          },
          ticks: {
            color: isDark ? "rgba(235,235,245,0.6)" : "rgba(60,60,67,0.6)",
            callback: (value: string | number) => `¥${Number(value) / 1000}k`,
          },
        },
          x: {
            grid: { display: false },
            ticks: {
              color: isDark ? "rgba(235,235,245,0.6)" : "rgba(60,60,67,0.6)",
            },
          },
        },
        interaction: {
          intersect: false,
          mode: 'index',
        },
      },
    }

  chartInstance = new Chart(ctx, config)
}

onMounted(() => {
  initChart()
})

watch(() => [themeStore.theme, props.values], () => {
  if (chartInstance) {
    chartInstance.destroy()
    initChart()
  }
})
const subtitle = computed(() => {
  switch (currentPeriod.value) {
    case 'week': return '近7天收款数据'
    case 'month': return '近30天收款数据'
    case 'year': return '近12个月收款数据'
    default: return '收款数据'
  }
})
</script>

<template>
  <GlassCard>
    <div class="glass-card-header">
      <div>
        <h3 class="glass-card-title">收入趋势</h3>
        <p class="glass-card-subtitle">{{ subtitle }}</p>
      </div>
      <div class="flex gap-sm">
        <button 
          v-for="p in ['week', 'month', 'year']" 
          :key="p"
          class="btn btn-sm"
          :class="currentPeriod === p ? 'btn-secondary active' : 'btn-ghost'"
          @click="setPeriod(p)"
        >
          {{ p === 'week' ? '周' : p === 'month' ? '月' : '年' }}
        </button>
      </div>
    </div>
    <div class="chart-container">
      <canvas ref="canvasRef"></canvas>
    </div>
  </GlassCard>
</template>

<style scoped>
.chart-container {
  position: relative;
  height: 300px;
  padding: 8px; /* spacing-sm */
}
</style>
