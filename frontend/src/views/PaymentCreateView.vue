<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import GlassCard from '@/components/common/GlassCard.vue'
import DatePicker from '@/components/common/DatePicker.vue'
import { dictionaryApi, type DictionaryItem } from '@/api/dictionary'
import { projectApi, paymentApi, type Project, type PaymentRequest } from '@/api/project'
import { useToast } from '@/composables/useToast'
import dayjs from 'dayjs'

const router = useRouter()
const route = useRoute()
const toast = useToast()
const routeProjectId = route.params.id

// Data Interfaces
interface PaymentItem {
  id?: number
  stage: string
  amount: string
  planDate: string
  method: string
  status: string
  remark: string
}

// State
const projects = ref<Project[]>([])
const currentProject = ref<Project | null>(null)
const selectedProjectId = ref<number | ''>(routeProjectId ? Number(routeProjectId) : '')
const loading = ref(false)
const originalPaymentIds = ref<number[]>([]) // Track original IDs for deletion

const paymentItems = ref<PaymentItem[]>([])

// Options
const paymentStageOptions = ref<DictionaryItem[]>([])
const paymentMethodOptions = ref<DictionaryItem[]>([])
const paymentStatuses = [
  { id: 'paid', name: '已收款' },
  { id: 'pending', name: '待收款' },
]

const currentProjectName = computed(() => {
  if (currentProject.value) return currentProject.value.name
  const p = projects.value.find((p) => p.id === selectedProjectId.value)
  return p ? p.name : ''
})

// === Methods ===

const fetchDictionaries = async () => {
    try {
        const [stageRes, methodRes] = await Promise.all([
            dictionaryApi.getItems('payment_stage'),
            dictionaryApi.getItems('payment_method')
        ])
        if (stageRes.data.code === 0) paymentStageOptions.value = stageRes.data.data || []
        if (methodRes.data.code === 0) paymentMethodOptions.value = methodRes.data.data || []
    } catch (e) {
        console.error('Failed to fetch dictionaries', e)
        toast.error('字典数据加载失败')
    }
}

const fetchProjectList = async () => {
    if (!routeProjectId) {
         try {
            const { data } = await projectApi.list({ page: 1, page_size: 100 }) 
            if (data.code === 0 && data.data) {
                projects.value = data.data.list || [] 
            }
        } catch (e) {
            console.error(e)
            toast.error('获取项目列表失败')
        }
    }
}

// Load Project Details & Payments
const fetchProjectData = async (projectId: number) => {
    loading.value = true
    try {
        // 1. Get Project Details
        const pRes = await projectApi.get(projectId)
        if (pRes.data.code === 0 && pRes.data.data) {
            currentProject.value = pRes.data.data
        }

        // 2. Get Existing Payments
        const payRes = await projectApi.getPayments(projectId)
        if (payRes.data.code === 0 && payRes.data.data) {
            const payments = payRes.data.data
            originalPaymentIds.value = payments.map((p) => p.id)
            paymentItems.value = payments.map(p => ({
                id: p.id,
                stage: p.stage,
                amount: p.amount.toString(),
                planDate: p.plan_date ? dayjs(p.plan_date).format('YYYY-MM-DD') : '',
                method: p.method || '',
                status: p.status,
                remark: p.remark
            }))
        } else {
             paymentItems.value = []
             originalPaymentIds.value = []
        }
        
        ensurePaymentItems()
    } catch (e) {
        console.error(e)
        toast.error('获取项目数据失败')
    } finally {
        loading.value = false
    }
}

const ensurePaymentItems = () => {
    if (!currentProject.value) return

    // If One-time payment, ensure exactly one item exists
    if (currentProject.value.payment_method === '一次性付款') {
         if (paymentItems.value.length === 0) {
            // New item
            addPaymentItem()
         } else if (paymentItems.value.length > 1) {
             // Keep only the first one (edge case cleanup)
             const first = paymentItems.value[0]
             if (first) {
                 paymentItems.value = [first]
             }
         }
         
         // Force stage to 'all'
         if (paymentItems.value[0]) {
             paymentItems.value[0].stage = 'all'
         }
    } else {
        // Installment: If empty, maybe add one? Or leave empty. 
        // ProjectCreateView logic: if empty, adds one. Let's start with one if empty.
        if (paymentItems.value.length === 0) {
            addPaymentItem()
        }
    }
}

const addPaymentItem = () => {
    let defaultStage = ''
    if (paymentStageOptions.value.length > 0) defaultStage = paymentStageOptions.value[0]?.value || ''

    let defaultMethod = ''
    if (paymentMethodOptions.value.length > 0) defaultMethod = paymentMethodOptions.value[0]?.value || ''

    paymentItems.value.push({
        stage: defaultStage,
        amount: '',
        planDate: '',
        method: defaultMethod,
        status: 'pending',
        remark: ''
    })
}

const removePaymentItem = (index: number) => {
    if (currentProject.value?.payment_method === '一次性付款') return
    paymentItems.value.splice(index, 1)
}


const savePayments = async () => {
    const projectId = Number(selectedProjectId.value)
    if (!projectId) return

    loading.value = true
    try {
        // 1. Delete removed items
        const currentIds = new Set(paymentItems.value.map(i => i.id).filter(id => id !== undefined) as number[])
        const toDelete = originalPaymentIds.value.filter(id => !currentIds.has(id))

        for (const id of toDelete) {
            try {
                await paymentApi.delete(id)
            } catch (e) {
                console.error(`Failed to delete payment ${id}`, e)
            }
        }

        // 2. Create / Update items
        for (const item of paymentItems.value) {
            if (!item.amount || parseFloat(item.amount) <= 0) continue

             const payload: PaymentRequest = {
                project_id: projectId,
                stage: item.stage || '款项',
                amount: parseFloat(item.amount),
                plan_date: item.planDate,
                status: item.status,
                method: item.method,
                remark: item.remark
            }

            if (item.id) {
                await paymentApi.update(item.id, payload)
            } else {
                await paymentApi.create(payload)
            }
        }

        toast.success('收款计划保存成功')
        router.back()

    } catch (e) {
        console.error(e)
        toast.error('保存失败')
    } finally {
        loading.value = false
    }
}

const handleSubmit = async () => {
    // Validation
    if (!selectedProjectId.value) {
        toast.error('请选择项目')
        return
    }
    
    // Check if any valid items
    const validItems = paymentItems.value.filter(i => i.amount && parseFloat(i.amount) > 0 && i.planDate)
    if (validItems.length === 0) {
        toast.error('请至少填写一笔有效的收款记录（需包含金额和日期）')
        return
    }

    await savePayments()
}


const handleCancel = () => router.back()


// === Watchers ===

// Watch selected project ID change
watch(selectedProjectId, async (newId) => {
    if (newId) {
        await fetchProjectData(Number(newId))
    } else {
        currentProject.value = null
        paymentItems.value = []
    }
})

// Lifecycle
onMounted(async () => {
    await fetchDictionaries()
    await fetchProjectList()
    
    // If route has ID, verify valid and load data
    if (selectedProjectId.value) {
        await fetchProjectData(Number(selectedProjectId.value))
    }
})

</script>

<template>
  <div class="payment-create-view">
    <!-- Back Header -->
    <div style="margin-top: -12px; margin-bottom: 24px" class="flex items-center gap-md">
      <button class="btn btn-ghost btn-sm pl-0 hover:bg-transparent" @click="handleCancel">
        <i class="ri-arrow-left-line text-2xl text-primary"></i>
      </button>

      <!-- Context Display / Selection -->
      <div class="project-context">
        <div v-if="routeProjectId" class="text-xl font-bold">
          {{ currentProjectName }}
        </div>
        <div v-else class="w-64">
          <div class="input-wrapper">
            <select v-model="selectedProjectId" class="form-select bg-glass" required>
              <option value="" disabled>请选择项目</option>
              <option v-for="p in projects" :key="p.id" :value="p.id">
                {{ p.name }}
              </option>
            </select>
            <i class="ri-arrow-down-s-line select-arrow"></i>
          </div>
        </div>
      </div>
    </div>

    <GlassCard class="w-full">
      <form @submit.prevent="handleSubmit" class="payment-form">
        <!-- Header with Add Button -->
        <div class="flex justify-between items-center mb-md">
           <h3 class="section-title mb-0 border-none pb-0 pl-0">收款计划</h3>
            <button 
              v-if="currentProject?.payment_method === '分期付款'" 
              type="button" 
              class="btn btn-sm btn-primary" 
              @click="addPaymentItem"
            >
              <i class="ri-add-line mr-1"></i>添加分期
            </button>
        </div>

        <!-- Payment Items List -->
        <div 
          v-for="(item, index) in paymentItems"
          :key="index"
          class="glass-panel p-md mb-md rounded-lg"
          style="border: 1px solid var(--color-primary);"
        >
             <!-- Item Header (Only for Installments) -->
             <div class="flex justify-between items-center mb-4" v-if="currentProject?.payment_method === '分期付款'">
               <span class="text-sm font-bold text-primary">第 {{ index + 1 }} 期</span>
               <button type="button" class="text-xs text-danger hover:underline" @click="removePaymentItem(index)">
                 删除
               </button>
             </div>

            <div class="form-grid">
              <!-- Row 1: Amount & Date -->
              <div class="input-group">
                <label>收款金额 (¥) <span class="text-red-500">*</span></label>
                <div class="input-wrapper">
                  <input
                    v-model="item.amount"
                    type="number"
                    placeholder="0.00"
                    min="0"
                    step="0.01"
                    required
                    spellcheck="false"
                    autocomplete="off"
                    autocorrect="off"
                    autocapitalize="off"
                  />
                </div>
              </div>

              <div class="input-group">
                <label>收款日期 <span class="text-red-500">*</span></label>
                <div class="input-wrapper">
                  <DatePicker v-model="item.planDate" required placeholder="请选择收款日期" />
                </div>
              </div>

              <!-- Row 2: Stage & Method -->
              <div class="input-group">
                <label>款项阶段 <span class="text-red-500">*</span></label>
                <div class="input-wrapper">
                  <select 
                    v-model="item.stage" 
                    class="form-select"
                    :disabled="currentProject?.payment_method === '一次性付款'"
                  >
                    <option v-for="s in paymentStageOptions" :key="s.value" :value="s.value">
                      {{ s.label }}
                    </option>
                  </select>
                  <i class="ri-arrow-down-s-line select-arrow"></i>
                </div>
              </div>

              <div class="input-group">
                <label>收款方式 <span class="text-red-500">*</span></label>
                <div class="input-wrapper">
                  <select v-model="item.method" class="form-select">
                    <option v-for="m in paymentMethodOptions" :key="m.value" :value="m.value">
                      {{ m.label }}
                    </option>
                  </select>
                  <i class="ri-arrow-down-s-line select-arrow"></i>
                </div>
              </div>

              <!-- Row 3: Status -->
              <div class="input-group">
                <label>状态 <span class="text-red-500">*</span></label>
                <div class="input-wrapper">
                  <select v-model="item.status" class="form-select">
                    <option v-for="status in paymentStatuses" :key="status.id" :value="status.id">
                      {{ status.name }}
                    </option>
                  </select>
                  <i class="ri-arrow-down-s-line select-arrow"></i>
                </div>
              </div>

              <!-- Row 4: Remark -->
              <div class="input-group span-2">
                <label>备注</label>
                <div class="input-wrapper">
                  <textarea
                    v-model="item.remark"
                    rows="2"
                    placeholder="请输入备注信息（选填）"
                    class="form-textarea"
                    spellcheck="false"
                    autocomplete="off"
                    autocorrect="off"
                    autocapitalize="off"
                  ></textarea>
                </div>
              </div>
            </div>
        </div>
        
        <!-- Empty State warning if no items -->
         <div v-if="paymentItems.length === 0" class="text-center text-secondary py-lg">
             暂无收款计划，请添加。
         </div>

        <div class="form-actions mt-xl">
          <button type="button" class="btn btn-ghost" @click="handleCancel">取消</button>
          <button type="submit" class="btn btn-primary" :disabled="loading">
            {{ loading ? '保存中...' : '确认保存' }}
          </button>
        </div>
      </form>
    </GlassCard>
  </div>
</template>

<style scoped>
.bg-glass {
  background: rgba(255, 255, 255, 0.5) !important;
  backdrop-filter: blur(10px);
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: var(--spacing-md);
  padding-left: var(--spacing-sm);
  border-left: 3px solid var(--color-primary);
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--spacing-md);
}

.input-group {
  margin-bottom: var(--spacing-md);
}

.input-group.span-2 {
  grid-column: span 2;
}

.input-group label {
  display: block;
  font-size: 13px;
  font-weight: 500;
  color: var(--text-secondary);
  margin-bottom: 8px;
}

.input-wrapper {
  position: relative;
}

.input-wrapper input,
.input-wrapper .form-select,
.form-textarea {
  width: 100%;
  padding: 10px 14px;
  font-size: 14px;
  color: var(--text-primary);
  background: var(--bg-base);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-sm);
  outline: none;
  transition: all 0.2s;
  appearance: none;
}

.input-wrapper input:focus,
.input-wrapper .form-select:focus,
.form-textarea:focus {
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px rgba(255, 159, 10, 0.1);
  background: var(--bg-elevated);
}

.select-arrow {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  color: var(--text-secondary);
  pointer-events: none;
}

.form-textarea {
  resize: vertical;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--spacing-md);
  border-top: 1px solid var(--border-color);
  padding-top: var(--spacing-lg);
}

@media (max-width: 640px) {
  .form-grid {
    grid-template-columns: 1fr;
  }
  .input-group.span-2 {
    grid-column: span 1;
  }
}

.text-primary {
  color: var(--color-primary);
}
.text-danger {
  color: var(--color-danger);
}
</style>
