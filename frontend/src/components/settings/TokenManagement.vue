<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { tokenApi, type PersonalAccessToken } from '@/api/token'
import { useToast } from '@/composables/useToast'
import { useConfirm } from '@/composables/useConfirm'

const toast = useToast()
const { confirm } = useConfirm()

const tokens = ref<PersonalAccessToken[]>([])
const loading = ref(false)

// Create Modal State
const showCreateModal = ref(false)
const createForm = ref({
  name: '',
  expires_in: 30 // Default 30 days
})
const creating = ref(false)

// New Token Success Modal State
const showSuccessModal = ref(false)
const newTokenValue = ref('')

const expiryOptions = [
  { label: '7天', value: 7 },
  { label: '30天', value: 30 },
  { label: '90天', value: 90 },
  { label: '1年', value: 365 },
  { label: '永不过期', value: 0 },
]

// Formatting
const formatDate = (dateStr: string | null) => {
  if (!dateStr) return '永不过期'
  return new Date(dateStr).toLocaleString()
}

const formatLastUsed = (dateStr: string | null) => {
  if (!dateStr) return '从未使用'
  return new Date(dateStr).toLocaleString()
}

// Actions
const fetchTokens = async () => {
  loading.value = true
  try {
    const res = await tokenApi.list()
    if (res.data.code === 0) {
      tokens.value = res.data.data
    }
  } catch (error) {
    console.error('Failed to fetch tokens:', error)
  } finally {
    loading.value = false
  }
}

const openCreateModal = () => {
  createForm.value = { name: '', expires_in: 30 }
  showCreateModal.value = true
}

const handleCreateToken = async () => {
  if (!createForm.value.name) {
    toast.warning('请填写令牌名称')
    return
  }
  
  creating.value = true
  try {
    const res = await tokenApi.create(createForm.value)
    if (res.data.code === 0) {
      newTokenValue.value = res.data.data.token
      showCreateModal.value = false
      showSuccessModal.value = true // Show the token once
      fetchTokens() // Refresh list
    } else {
      toast.error(res.data.message || '创建失败')
    }
  } catch (error) {
    console.error('Create token failed:', error)
    toast.error('创建失败')
  } finally {
    creating.value = false
  }
}

const handleRevoke = async (token: PersonalAccessToken) => {
  const confirmed = await confirm({
    title: '撤销令牌',
    message: `确定要撤销令牌 "${token.name}" 吗？撤销后将无法恢复使用。`
  })
  
  if (confirmed) {
    try {
      const res = await tokenApi.revoke(token.id)
      if (res.data.code === 0) {
        toast.success('撤销成功')
        fetchTokens()
      }
    } catch (error) {
      console.error('Revoke failed:', error)
       toast.error('撤销失败')
    }
  }
}

const handleDelete = async (token: PersonalAccessToken) => {
  const confirmed = await confirm({
    title: '删除令牌',
    message: `确定要彻底删除令牌 "${token.name}" 吗？`
  })
  
  if (confirmed) {
    try {
      const res = await tokenApi.delete(token.id)
      if (res.data.code === 0) {
        toast.success('删除成功')
        fetchTokens()
      }
    } catch (error) {
      console.error('Delete failed:', error)
       toast.error('删除失败')
    }
  }
}

const copyToken = () => {
  navigator.clipboard.writeText(newTokenValue.value)
  toast.success('复制成功')
}

onMounted(() => {
  fetchTokens()
})
</script>

<template>
  <div class="token-management h-full flex flex-col">
    <!-- Header -->
    <div class="glass-card-header border-b border-color-border p-md flex justify-between items-center">
      <div>
        <h3 class="glass-card-title">开发者设置 / 个人访问令牌</h3>
        <p class="text-sm text-secondary mt-1">创建和管理用于访问 API 的个人访问令牌 (Personal Access Token)。</p>
      </div>
      <button class="btn btn-primary btn-sm" @click="openCreateModal">
        <i class="ri-add-line mr-1"></i> 生成新令牌
      </button>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-auto p-md">
      <div v-if="loading" class="flex justify-center py-8">
        <div class="spinner"></div> <!-- Assuming global spinner class or replace with text -->
      </div>
      
      <div v-else-if="tokens.length === 0" class="text-center py-12 text-secondary">
        <i class="ri-key-2-line text-4xl mb-4 block opacity-50"></i>
        暂无任何访问令牌
      </div>

      <div v-else class="glass-table-container">
        <table class="glass-table w-full text-left border-collapse">
          <thead>
            <tr>
              <th class="p-3 border-b border-color-border text-sm font-medium text-secondary">名称</th>
              <th class="p-3 border-b border-color-border text-sm font-medium text-secondary">状态</th>
              <th class="p-3 border-b border-color-border text-sm font-medium text-secondary">最后使用</th>
              <th class="p-3 border-b border-color-border text-sm font-medium text-secondary">过期时间</th>
              <th class="p-3 border-b border-color-border text-sm font-medium text-secondary text-right">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="token in tokens" :key="token.id" class="hover:bg-white/5 transition-colors">
              <td class="p-3 border-b border-color-border">
                <div class="font-medium">{{ token.name }}</div>
                <div class="text-xs text-secondary font-mono mt-0.5">ID: {{ token.id }}</div>
              </td>
              <td class="p-3 border-b border-color-border">
                <span v-if="token.status === 1" class="badge badge-success text-xs px-2 py-0.5 rounded-full bg-green-500/10 text-green-500 border border-green-500/20">Active</span>
                <span v-else class="badge badge-error text-xs px-2 py-0.5 rounded-full bg-red-500/10 text-red-500 border border-red-500/20">Revoked</span>
              </td>
              <td class="p-3 border-b border-color-border text-sm text-secondary">
                {{ formatLastUsed(token.last_used_at) }}
              </td>
              <td class="p-3 border-b border-color-border text-sm text-secondary">
                {{ formatDate(token.expires_at) }}
              </td>
              <td class="p-3 border-b border-color-border text-right">
                <button 
                  v-if="token.status === 1"
                  class="btn btn-ghost btn-xs text-warning mr-2" 
                  title="撤销"
                  @click="handleRevoke(token)"
                >
                  撤销
                </button>
                <button 
                  class="btn btn-ghost btn-xs text-danger" 
                  title="删除"
                  @click="handleDelete(token)"
                >
                  删除
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Create Modal -->
    <div v-if="showCreateModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm">
      <div class="glass-card w-[400px] p-0 animate-in fade-in zoom-in duration-200">
        <div class="p-4 border-b border-color-border flex justify-between items-center">
          <h3 class="font-bold">生成新令牌</h3>
          <button @click="showCreateModal = false" class="btn btn-ghost btn-icon btn-sm"><i class="ri-close-line"></i></button>
        </div>
        <div class="p-6 space-y-4">
          <div>
            <label class="block text-sm font-medium mb-1">令牌名称 <span class="text-red-500">*</span></label>
            <input v-model="createForm.name" type="text" class="form-input w-full" placeholder="例如: CI/CD Scripts, Development" />
          </div>
          <div>
            <label class="block text-sm font-medium mb-1">过期时间</label>
            <select v-model="createForm.expires_in" class="form-select w-full bg-transparent border border-color-border rounded-md p-2">
              <option v-for="opt in expiryOptions" :key="opt.value" :value="opt.value">{{ opt.label }}</option>
            </select>
          </div>
        </div>
        <div class="p-4 border-t border-color-border flex justify-end gap-2">
          <button class="btn btn-ghost" @click="showCreateModal = false">取消</button>
          <button class="btn btn-primary" :disabled="creating" @click="handleCreateToken">
            {{ creating ? '生成中...' : '生成令牌' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Success Modal -->
    <div v-if="showSuccessModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm">
      <div class="glass-card w-[480px] p-0 animate-in fade-in zoom-in duration-200 border border-green-500/30 shadow-[0_0_50px_rgba(34,197,94,0.1)]">
        <div class="p-4 border-b border-color-border bg-green-500/5">
          <h3 class="font-bold text-green-500 flex items-center">
            <i class="ri-checkbox-circle-fill mr-2 text-xl"></i> 生成成功
          </h3>
        </div>
        <div class="p-6 space-y-4">
          <p class="text-secondary text-sm">
            这是您唯一一次能看到该令牌的机会。请立即复制并将其保存在安全的地方。
            <br>如果丢失，您将无法恢复，只能重新生成。
          </p>
          <div class="relative mt-2">
            <input type="text" readonly :value="newTokenValue" class="form-input w-full pr-12 font-mono text-sm bg-black/5 dark:bg-white/5 border-primary/30" />
            <button class="absolute right-1 top-1 btn btn-sm btn-ghost text-primary" @click="copyToken">
              <i class="ri-file-copy-line"></i>
            </button>
          </div>
        </div>
        <div class="p-4 border-t border-color-border flex justify-end">
          <button class="btn btn-primary" @click="showSuccessModal = false">我已复制</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Reuse existing glass table styles or define minimal ones if missing */
.glass-table th {
  font-weight: 600;
}
.glass-table td {
  vertical-align: middle;
}
</style>
